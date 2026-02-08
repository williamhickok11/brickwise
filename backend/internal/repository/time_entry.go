package repository

import (
	"database/sql"
	"fmt"
	"os"
	"strings"
	"time"

	"brickwise/backend/internal/models"
)

// TimeEntryRepository defines the interface for time entry data access
type TimeEntryRepository interface {
	Create(entry *models.TimeEntry) error
	GetByID(id int) (*models.TimeEntry, error)
	GetAll(filter *models.TimeEntryFilter) ([]*models.TimeEntry, error)
	Update(entry *models.TimeEntry) error
	Delete(id int) error
}

// SQLTimeEntryRepository implements TimeEntryRepository using SQL database
type SQLTimeEntryRepository struct {
	db       *sql.DB
	isSQLite bool
}

// NewTimeEntryRepository creates a new time entry repository
func NewTimeEntryRepository(db *sql.DB) TimeEntryRepository {
	isSQLite := os.Getenv("DATABASE_URL") == ""
	return &SQLTimeEntryRepository{db: db, isSQLite: isSQLite}
}

// Create inserts a new time entry into the database
func (r *SQLTimeEntryRepository) Create(entry *models.TimeEntry) error {
	if r.isSQLite {
		// SQLite: use LastInsertId
		query := `INSERT INTO time_entries (property_id, date, category, description, hours, notes, mileage, full_drive) 
		          VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
		var fullDriveInt int
		if entry.FullDrive {
			fullDriveInt = 1
		}
		result, err := r.db.Exec(query, entry.PropertyID, entry.Date.Time, entry.Category, entry.Description, entry.Hours, entry.Notes, entry.Mileage, fullDriveInt)
		if err != nil {
			return fmt.Errorf("failed to create time entry: %w", err)
		}
		id, err := result.LastInsertId()
		if err != nil {
			return fmt.Errorf("failed to get last insert id: %w", err)
		}
		entry.ID = int(id)
		// Fetch the created record to get timestamps
		return r.db.QueryRow(
			`SELECT created_at, updated_at FROM time_entries WHERE id = $1`,
			entry.ID,
		).Scan(&entry.CreatedAt, &entry.UpdatedAt)
	}

	// Postgres: use RETURNING
	query := `INSERT INTO time_entries (property_id, date, category, description, hours, notes, mileage, full_drive) 
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8) 
	          RETURNING id, created_at, updated_at`
	
	err := r.db.QueryRow(
		query,
		entry.PropertyID,
		entry.Date.Time,
		entry.Category,
		entry.Description,
		entry.Hours,
		entry.Notes,
		entry.Mileage,
		entry.FullDrive,
	).Scan(&entry.ID, &entry.CreatedAt, &entry.UpdatedAt)
	
	if err != nil {
		return fmt.Errorf("failed to create time entry: %w", err)
	}
	
	return nil
}

// GetByID retrieves a time entry by ID
func (r *SQLTimeEntryRepository) GetByID(id int) (*models.TimeEntry, error) {
	entry := &models.TimeEntry{}
	query := `SELECT id, property_id, date, category, description, hours, notes, mileage, full_drive, created_at, updated_at 
	          FROM time_entries 
	          WHERE id = $1`
	
	var err error
	if r.isSQLite {
		var fullDriveInt int
		var dateTime time.Time
		err = r.db.QueryRow(query, id).Scan(
			&entry.ID,
			&entry.PropertyID,
			&dateTime,
			&entry.Category,
			&entry.Description,
			&entry.Hours,
			&entry.Notes,
			&entry.Mileage,
			&fullDriveInt,
			&entry.CreatedAt,
			&entry.UpdatedAt,
		)
		entry.Date.Time = dateTime
		entry.FullDrive = fullDriveInt == 1
	} else {
		var dateTime time.Time
		err = r.db.QueryRow(query, id).Scan(
			&entry.ID,
			&entry.PropertyID,
			&dateTime,
			&entry.Category,
			&entry.Description,
			&entry.Hours,
			&entry.Notes,
			&entry.Mileage,
			&entry.FullDrive,
			&entry.CreatedAt,
			&entry.UpdatedAt,
		)
		entry.Date.Time = dateTime
	}
	
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("time entry with id %d not found", id)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get time entry: %w", err)
	}
	
	return entry, nil
}

// GetAll retrieves all time entries with optional filters
func (r *SQLTimeEntryRepository) GetAll(filter *models.TimeEntryFilter) ([]*models.TimeEntry, error) {
	var conditions []string
	var args []interface{}
	argIndex := 1

	if filter != nil {
		if filter.PropertyID != nil {
			if *filter.PropertyID == 0 {
				// NULL property_id means "General"
				conditions = append(conditions, fmt.Sprintf("property_id IS NULL"))
			} else {
				conditions = append(conditions, fmt.Sprintf("property_id = $%d", argIndex))
				args = append(args, *filter.PropertyID)
				argIndex++
			}
		}
		if filter.StartDate != nil {
			conditions = append(conditions, fmt.Sprintf("date >= $%d", argIndex))
			args = append(args, *filter.StartDate)
			argIndex++
		}
		if filter.EndDate != nil {
			conditions = append(conditions, fmt.Sprintf("date <= $%d", argIndex))
			args = append(args, *filter.EndDate)
			argIndex++
		}
		if filter.Category != nil && *filter.Category != "" {
			conditions = append(conditions, fmt.Sprintf("category = $%d", argIndex))
			args = append(args, *filter.Category)
			argIndex++
		}
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	query := fmt.Sprintf(`SELECT id, property_id, date, category, description, hours, notes, mileage, full_drive, created_at, updated_at 
	          FROM time_entries 
	          %s
	          ORDER BY date DESC, created_at DESC`, whereClause)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query time entries: %w", err)
	}
	defer rows.Close()

	var entries []*models.TimeEntry
	for rows.Next() {
		entry := &models.TimeEntry{}
		if r.isSQLite {
			var fullDriveInt int
				var dateTime time.Time
			if err := rows.Scan(
				&entry.ID,
				&entry.PropertyID,
				&dateTime,
				&entry.Category,
				&entry.Description,
				&entry.Hours,
				&entry.Notes,
				&entry.Mileage,
				&fullDriveInt,
				&entry.CreatedAt,
				&entry.UpdatedAt,
			); err != nil {
				return nil, fmt.Errorf("failed to scan time entry: %w", err)
			}
			entry.Date.Time = dateTime
			entry.FullDrive = fullDriveInt == 1
		} else {
			var dateTime time.Time
			if err := rows.Scan(
				&entry.ID,
				&entry.PropertyID,
				&dateTime,
				&entry.Category,
				&entry.Description,
				&entry.Hours,
				&entry.Notes,
				&entry.Mileage,
				&entry.FullDrive,
				&entry.CreatedAt,
				&entry.UpdatedAt,
			); err != nil {
				return nil, fmt.Errorf("failed to scan time entry: %w", err)
			}
			entry.Date.Time = dateTime
		}
		entries = append(entries, entry)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating time entries: %w", err)
	}

	return entries, nil
}

// Update updates an existing time entry
func (r *SQLTimeEntryRepository) Update(entry *models.TimeEntry) error {
	if r.isSQLite {
		// SQLite: update and then fetch
		var fullDriveInt int
		if entry.FullDrive {
			fullDriveInt = 1
		}
		query := `UPDATE time_entries 
		          SET property_id = $1, date = $2, category = $3, description = $4, hours = $5, notes = $6, mileage = $7, full_drive = $8, updated_at = CURRENT_TIMESTAMP 
		          WHERE id = $9`
		result, err := r.db.Exec(query, entry.PropertyID, entry.Date.Time, entry.Category, entry.Description, entry.Hours, entry.Notes, entry.Mileage, fullDriveInt, entry.ID)
		if err != nil {
			return fmt.Errorf("failed to update time entry: %w", err)
		}
		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return fmt.Errorf("failed to get rows affected: %w", err)
		}
		if rowsAffected == 0 {
			return fmt.Errorf("time entry with id %d not found", entry.ID)
		}
		// Fetch updated timestamp
		return r.db.QueryRow(`SELECT updated_at FROM time_entries WHERE id = $1`, entry.ID).
			Scan(&entry.UpdatedAt)
	}

	// Postgres: use RETURNING
	query := `UPDATE time_entries 
	          SET property_id = $1, date = $2, category = $3, description = $4, hours = $5, notes = $6, mileage = $7, full_drive = $8, updated_at = CURRENT_TIMESTAMP 
	          WHERE id = $9 
	          RETURNING updated_at`

	err := r.db.QueryRow(
		query,
		entry.PropertyID,
		entry.Date.Time,
		entry.Category,
		entry.Description,
		entry.Hours,
		entry.Notes,
		entry.Mileage,
		entry.FullDrive,
		entry.ID,
	).Scan(&entry.UpdatedAt)

	if err == sql.ErrNoRows {
		return fmt.Errorf("time entry with id %d not found", entry.ID)
	}
	if err != nil {
		return fmt.Errorf("failed to update time entry: %w", err)
	}

	return nil
}

// Delete removes a time entry by ID
func (r *SQLTimeEntryRepository) Delete(id int) error {
	query := `DELETE FROM time_entries WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete time entry: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("time entry with id %d not found", id)
	}

	return nil
}
