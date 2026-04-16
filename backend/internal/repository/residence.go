package repository

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"brickwise/backend/internal/models"
)

// ResidenceRepository defines the interface for residence/tenant data access.
type ResidenceRepository interface {
	Create(residence *models.Residence) error
	GetByID(id int) (*models.Residence, error)
	GetAll(filter *models.ResidenceFilter) ([]*models.Residence, error)
	Update(residence *models.Residence) error
}

// SQLResidenceRepository implements ResidenceRepository using SQL database.
type SQLResidenceRepository struct {
	db       *sql.DB
	isSQLite bool
}

func NewResidenceRepository(db *sql.DB) ResidenceRepository {
	isSQLite := os.Getenv("DATABASE_URL") == ""
	return &SQLResidenceRepository{db: db, isSQLite: isSQLite}
}

func (r *SQLResidenceRepository) Create(residence *models.Residence) error {
	var endDateVal any = nil
	if residence.EndDate != nil {
		endDateVal = residence.EndDate.Time
	}

	if r.isSQLite {
		isActiveInt := 0
		if residence.IsActive {
			isActiveInt = 1
		}

		query := `INSERT INTO residences (property_id, name, phone, email, start_date, end_date, is_active)
		          VALUES ($1, $2, $3, $4, $5, $6, $7)`
		result, err := r.db.Exec(
			query,
			residence.PropertyID,
			residence.Name,
			residence.Phone,
			residence.Email,
			residence.StartDate.Time,
			endDateVal,
			isActiveInt,
		)
		if err != nil {
			return fmt.Errorf("failed to create residence: %w", err)
		}

		id, err := result.LastInsertId()
		if err != nil {
			return fmt.Errorf("failed to get last insert id: %w", err)
		}
		residence.ID = int(id)

		return r.db.QueryRow(
			`SELECT created_at, updated_at FROM residences WHERE id = $1`,
			residence.ID,
		).Scan(&residence.CreatedAt, &residence.UpdatedAt)
	}

	query := `INSERT INTO residences (property_id, name, phone, email, start_date, end_date, is_active)
	          VALUES ($1, $2, $3, $4, $5, $6, $7)
	          RETURNING id, created_at, updated_at`

	return r.db.QueryRow(
		query,
		residence.PropertyID,
		residence.Name,
		residence.Phone,
		residence.Email,
		residence.StartDate.Time,
		endDateVal,
		residence.IsActive,
	).Scan(&residence.ID, &residence.CreatedAt, &residence.UpdatedAt)
}

func (r *SQLResidenceRepository) GetByID(id int) (*models.Residence, error) {
	residence := &models.Residence{}
	query := `SELECT id, property_id, name, phone, email, start_date, end_date, is_active, created_at, updated_at
	          FROM residences
	          WHERE id = $1`

	if r.isSQLite {
		var startDate time.Time
		var endDate sql.NullTime
		var isActiveInt int
		if err := r.db.QueryRow(query, id).Scan(
			&residence.ID,
			&residence.PropertyID,
			&residence.Name,
			&residence.Phone,
			&residence.Email,
			&startDate,
			&endDate,
			&isActiveInt,
			&residence.CreatedAt,
			&residence.UpdatedAt,
		); err != nil {
			if err == sql.ErrNoRows {
				return nil, fmt.Errorf("residence with id %d not found", id)
			}
			return nil, fmt.Errorf("failed to get residence: %w", err)
		}

		residence.StartDate.Time = startDate
		if endDate.Valid {
			d := models.Date{Time: endDate.Time}
			residence.EndDate = &d
		} else {
			residence.EndDate = nil
		}
		residence.IsActive = isActiveInt == 1
		return residence, nil
	}

	var startDate time.Time
	var endDate sql.NullTime
	if err := r.db.QueryRow(query, id).Scan(
		&residence.ID,
		&residence.PropertyID,
		&residence.Name,
		&residence.Phone,
		&residence.Email,
		&startDate,
		&endDate,
		&residence.IsActive,
		&residence.CreatedAt,
		&residence.UpdatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("residence with id %d not found", id)
		}
		return nil, fmt.Errorf("failed to get residence: %w", err)
	}

	residence.StartDate.Time = startDate
	if endDate.Valid {
		d := models.Date{Time: endDate.Time}
		residence.EndDate = &d
	} else {
		residence.EndDate = nil
	}

	return residence, nil
}

func (r *SQLResidenceRepository) GetAll(filter *models.ResidenceFilter) ([]*models.Residence, error) {
	var conditions []string
	var args []interface{}
	argIndex := 1

	if filter != nil {
		if filter.PropertyID != nil {
			conditions = append(conditions, fmt.Sprintf("property_id = $%d", argIndex))
			args = append(args, *filter.PropertyID)
			argIndex++
		}
		if filter.IsActive != nil {
			conditions = append(conditions, fmt.Sprintf("is_active = $%d", argIndex))
			args = append(args, *filter.IsActive)
			argIndex++
		}
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + joinConditions(conditions, " AND ")
	}

	query := fmt.Sprintf(`SELECT id, property_id, name, phone, email, start_date, end_date, is_active, created_at, updated_at
	          FROM residences
	          %s
	          ORDER BY start_date DESC, created_at DESC`, whereClause)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query residences: %w", err)
	}
	defer rows.Close()

	residences := []*models.Residence{}
	for rows.Next() {
		residence := &models.Residence{}

		if r.isSQLite {
			var startDate time.Time
			var endDate sql.NullTime
			var isActiveInt int
			if err := rows.Scan(
				&residence.ID,
				&residence.PropertyID,
				&residence.Name,
				&residence.Phone,
				&residence.Email,
				&startDate,
				&endDate,
				&isActiveInt,
				&residence.CreatedAt,
				&residence.UpdatedAt,
			); err != nil {
				return nil, fmt.Errorf("failed to scan residence: %w", err)
			}

			residence.StartDate.Time = startDate
			residence.IsActive = isActiveInt == 1
			if endDate.Valid {
				d := models.Date{Time: endDate.Time}
				residence.EndDate = &d
			} else {
				residence.EndDate = nil
			}
		} else {
			var startDate time.Time
			var endDate sql.NullTime
			if err := rows.Scan(
				&residence.ID,
				&residence.PropertyID,
				&residence.Name,
				&residence.Phone,
				&residence.Email,
				&startDate,
				&endDate,
				&residence.IsActive,
				&residence.CreatedAt,
				&residence.UpdatedAt,
			); err != nil {
				return nil, fmt.Errorf("failed to scan residence: %w", err)
			}

			residence.StartDate.Time = startDate
			if endDate.Valid {
				d := models.Date{Time: endDate.Time}
				residence.EndDate = &d
			} else {
				residence.EndDate = nil
			}
		}

		residences = append(residences, residence)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating residences: %w", err)
	}

	return residences, nil
}

func (r *SQLResidenceRepository) Update(residence *models.Residence) error {
	var endDateVal any = nil
	if residence.EndDate != nil {
		endDateVal = residence.EndDate.Time
	}

	if r.isSQLite {
		isActiveInt := 0
		if residence.IsActive {
			isActiveInt = 1
		}

		query := `UPDATE residences
		          SET property_id = $1, name = $2, phone = $3, email = $4,
		              start_date = $5, end_date = $6, is_active = $7,
		              updated_at = CURRENT_TIMESTAMP
		          WHERE id = $8`

		result, err := r.db.Exec(
			query,
			residence.PropertyID,
			residence.Name,
			residence.Phone,
			residence.Email,
			residence.StartDate.Time,
			endDateVal,
			isActiveInt,
			residence.ID,
		)
		if err != nil {
			return fmt.Errorf("failed to update residence: %w", err)
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return fmt.Errorf("failed to get rows affected: %w", err)
		}
		if rowsAffected == 0 {
			return fmt.Errorf("residence with id %d not found", residence.ID)
		}

		return r.db.QueryRow(`SELECT updated_at FROM residences WHERE id = $1`, residence.ID).
			Scan(&residence.UpdatedAt)
	}

	query := `UPDATE residences
	          SET property_id = $1, name = $2, phone = $3, email = $4,
	              start_date = $5, end_date = $6, is_active = $7,
	              updated_at = CURRENT_TIMESTAMP
	          WHERE id = $8
	          RETURNING updated_at`

	if err := r.db.QueryRow(
		query,
		residence.PropertyID,
		residence.Name,
		residence.Phone,
		residence.Email,
		residence.StartDate.Time,
		endDateVal,
		residence.IsActive,
		residence.ID,
	).Scan(&residence.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("residence with id %d not found", residence.ID)
		}
		return fmt.Errorf("failed to update residence: %w", err)
	}

	return nil
}

func joinConditions(conditions []string, sep string) string {
	// Small helper to avoid importing strings in this file.
	if len(conditions) == 0 {
		return ""
	}
	out := conditions[0]
	for i := 1; i < len(conditions); i++ {
		out += sep + conditions[i]
	}
	return out
}

