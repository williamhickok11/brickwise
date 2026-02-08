package repository

import (
	"database/sql"
	"fmt"
	"os"

	"brickwise/backend/internal/models"
)

// PropertyRepository defines the interface for property data access
// This abstraction allows switching between SQLite and Postgres easily
type PropertyRepository interface {
	Create(property *models.Property) error
	GetByID(id int) (*models.Property, error)
	GetAll() ([]*models.Property, error)
	Update(property *models.Property) error
	Delete(id int) error
}

// SQLPropertyRepository implements PropertyRepository using SQL database
type SQLPropertyRepository struct {
	db       *sql.DB
	isSQLite bool
}

// NewPropertyRepository creates a new property repository
func NewPropertyRepository(db *sql.DB) PropertyRepository {
	isSQLite := os.Getenv("DATABASE_URL") == ""
	return &SQLPropertyRepository{db: db, isSQLite: isSQLite}
}

// Create inserts a new property into the database
func (r *SQLPropertyRepository) Create(property *models.Property) error {
	if r.isSQLite {
		// SQLite doesn't support RETURNING in older versions, use LastInsertId
		query := `INSERT INTO properties (name, address, property_type, default_mileage) 
		          VALUES ($1, $2, $3, $4)`
		result, err := r.db.Exec(query, property.Name, property.Address, property.PropertyType, property.DefaultMileage)
		if err != nil {
			return fmt.Errorf("failed to create property: %w", err)
		}
		id, err := result.LastInsertId()
		if err != nil {
			return fmt.Errorf("failed to get last insert id: %w", err)
		}
		property.ID = int(id)
		// Fetch the created record to get timestamps
		return r.db.QueryRow(
			`SELECT default_mileage, created_at, updated_at FROM properties WHERE id = $1`,
			property.ID,
		).Scan(&property.DefaultMileage, &property.CreatedAt, &property.UpdatedAt)
	}

	// Postgres supports RETURNING
	query := `INSERT INTO properties (name, address, property_type, default_mileage) 
	          VALUES ($1, $2, $3, $4) 
	          RETURNING id, default_mileage, created_at, updated_at`
	
	err := r.db.QueryRow(query, property.Name, property.Address, property.PropertyType, property.DefaultMileage).
		Scan(&property.ID, &property.DefaultMileage, &property.CreatedAt, &property.UpdatedAt)
	
	if err != nil {
		return fmt.Errorf("failed to create property: %w", err)
	}
	
	return nil
}

// GetByID retrieves a property by ID
func (r *SQLPropertyRepository) GetByID(id int) (*models.Property, error) {
	property := &models.Property{}
	query := `SELECT id, name, address, property_type, default_mileage, created_at, updated_at 
	          FROM properties 
	          WHERE id = $1`
	
	err := r.db.QueryRow(query, id).Scan(
		&property.ID,
		&property.Name,
		&property.Address,
		&property.PropertyType,
		&property.DefaultMileage,
		&property.CreatedAt,
		&property.UpdatedAt,
	)
	
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("property with id %d not found", id)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get property: %w", err)
	}
	
	return property, nil
}

// GetAll retrieves all properties
func (r *SQLPropertyRepository) GetAll() ([]*models.Property, error) {
	query := `SELECT id, name, address, property_type, default_mileage, created_at, updated_at 
	          FROM properties 
	          ORDER BY created_at DESC`
	
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query properties: %w", err)
	}
	defer rows.Close()
	
	var properties []*models.Property
	for rows.Next() {
		property := &models.Property{}
		if err := rows.Scan(
			&property.ID,
			&property.Name,
			&property.Address,
			&property.PropertyType,
			&property.DefaultMileage,
			&property.CreatedAt,
			&property.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan property: %w", err)
		}
		properties = append(properties, property)
	}
	
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating properties: %w", err)
	}
	
	return properties, nil
}

// Update updates an existing property
func (r *SQLPropertyRepository) Update(property *models.Property) error {
	if r.isSQLite {
		// SQLite: update and then fetch
		query := `UPDATE properties 
		          SET name = $1, address = $2, property_type = $3, default_mileage = $4, updated_at = CURRENT_TIMESTAMP 
		          WHERE id = $5`
		result, err := r.db.Exec(query, property.Name, property.Address, property.PropertyType, property.DefaultMileage, property.ID)
		if err != nil {
			return fmt.Errorf("failed to update property: %w", err)
		}
		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return fmt.Errorf("failed to get rows affected: %w", err)
		}
		if rowsAffected == 0 {
			return fmt.Errorf("property with id %d not found", property.ID)
		}
		// Fetch updated timestamp
		return r.db.QueryRow(`SELECT updated_at FROM properties WHERE id = $1`, property.ID).
			Scan(&property.UpdatedAt)
	}

	// Postgres: use RETURNING
	query := `UPDATE properties 
	          SET name = $1, address = $2, property_type = $3, default_mileage = $4, updated_at = CURRENT_TIMESTAMP 
	          WHERE id = $5 
	          RETURNING updated_at`
	
	err := r.db.QueryRow(
		query,
		property.Name,
		property.Address,
		property.PropertyType,
		property.DefaultMileage,
		property.ID,
	).Scan(&property.UpdatedAt)
	
	if err == sql.ErrNoRows {
		return fmt.Errorf("property with id %d not found", property.ID)
	}
	if err != nil {
		return fmt.Errorf("failed to update property: %w", err)
	}
	
	return nil
}

// Delete removes a property by ID
func (r *SQLPropertyRepository) Delete(id int) error {
	query := `DELETE FROM properties WHERE id = $1`
	
	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete property: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("property with id %d not found", id)
	}
	
	return nil
}
