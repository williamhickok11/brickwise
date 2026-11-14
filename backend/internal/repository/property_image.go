package repository

import (
	"database/sql"
	"fmt"
	"os"

	"brickwise/backend/internal/models"
)

// PropertyImageRepository defines data access for property photos.
type PropertyImageRepository interface {
	ListByPropertyID(propertyID int) ([]*models.PropertyImage, error)
	GetByID(id int) (*models.PropertyImage, error)
	Create(img *models.PropertyImage) error
	UpdateSize(id int, size int64) error
	Delete(id int) error
	ClearMain(propertyID int) error
	SetMain(propertyID, imageID int) error
	CountByPropertyID(propertyID int) (int, error)
}

type SQLPropertyImageRepository struct {
	db       *sql.DB
	isSQLite bool
}

func NewPropertyImageRepository(db *sql.DB) PropertyImageRepository {
	isSQLite := os.Getenv("DATABASE_URL") == ""
	return &SQLPropertyImageRepository{db: db, isSQLite: isSQLite}
}

func (r *SQLPropertyImageRepository) ListByPropertyID(propertyID int) ([]*models.PropertyImage, error) {
	query := `SELECT id, property_id, file_name, content_type, size, is_main, created_at
	          FROM property_images
	          WHERE property_id = $1
	          ORDER BY is_main DESC, created_at ASC`
	rows, err := r.db.Query(query, propertyID)
	if err != nil {
		return nil, fmt.Errorf("failed to list property images: %w", err)
	}
	defer rows.Close()

	var images []*models.PropertyImage
	for rows.Next() {
		img, err := scanPropertyImage(rows)
		if err != nil {
			return nil, err
		}
		images = append(images, img)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating property images: %w", err)
	}
	return images, nil
}

func (r *SQLPropertyImageRepository) GetByID(id int) (*models.PropertyImage, error) {
	query := `SELECT id, property_id, file_name, content_type, size, is_main, created_at
	          FROM property_images WHERE id = $1`
	row := r.db.QueryRow(query, id)
	img, err := scanPropertyImage(row)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("property image with id %d not found", id)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get property image: %w", err)
	}
	return img, nil
}

func (r *SQLPropertyImageRepository) Create(img *models.PropertyImage) error {
	isMain := 0
	if img.IsMain {
		isMain = 1
	}

	if r.isSQLite {
		result, err := r.db.Exec(
			`INSERT INTO property_images (property_id, file_name, content_type, size, is_main)
			 VALUES ($1, $2, $3, $4, $5)`,
			img.PropertyID, img.FileName, img.ContentType, img.Size, isMain,
		)
		if err != nil {
			return fmt.Errorf("failed to insert property image: %w", err)
		}
		id, err := result.LastInsertId()
		if err != nil {
			return fmt.Errorf("failed to get last insert id: %w", err)
		}
		img.ID = int(id)
		return r.db.QueryRow(
			`SELECT created_at FROM property_images WHERE id = $1`, img.ID,
		).Scan(&img.CreatedAt)
	}

	query := `INSERT INTO property_images (property_id, file_name, content_type, size, is_main)
	          VALUES ($1, $2, $3, $4, $5)
	          RETURNING id, created_at`
	return r.db.QueryRow(query, img.PropertyID, img.FileName, img.ContentType, img.Size, isMain).
		Scan(&img.ID, &img.CreatedAt)
}

func (r *SQLPropertyImageRepository) UpdateSize(id int, size int64) error {
	result, err := r.db.Exec(`UPDATE property_images SET size = $1 WHERE id = $2`, size, id)
	if err != nil {
		return fmt.Errorf("failed to update property image size: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("property image with id %d not found", id)
	}
	return nil
}

func (r *SQLPropertyImageRepository) Delete(id int) error {
	result, err := r.db.Exec(`DELETE FROM property_images WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to delete property image: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("property image with id %d not found", id)
	}
	return nil
}

func (r *SQLPropertyImageRepository) ClearMain(propertyID int) error {
	_, err := r.db.Exec(`UPDATE property_images SET is_main = 0 WHERE property_id = $1`, propertyID)
	if err != nil {
		return fmt.Errorf("failed to clear main image: %w", err)
	}
	return nil
}

func (r *SQLPropertyImageRepository) SetMain(propertyID, imageID int) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	if _, err := tx.Exec(`UPDATE property_images SET is_main = 0 WHERE property_id = $1`, propertyID); err != nil {
		return fmt.Errorf("failed to clear main image: %w", err)
	}
	result, err := tx.Exec(
		`UPDATE property_images SET is_main = 1 WHERE id = $1 AND property_id = $2`,
		imageID, propertyID,
	)
	if err != nil {
		return fmt.Errorf("failed to set main image: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("property image with id %d not found", imageID)
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit main image change: %w", err)
	}
	return nil
}

func (r *SQLPropertyImageRepository) CountByPropertyID(propertyID int) (int, error) {
	var count int
	err := r.db.QueryRow(`SELECT COUNT(*) FROM property_images WHERE property_id = $1`, propertyID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count property images: %w", err)
	}
	return count, nil
}

type rowScanner interface {
	Scan(dest ...any) error
}

func scanPropertyImage(row rowScanner) (*models.PropertyImage, error) {
	img := &models.PropertyImage{}
	if err := row.Scan(
		&img.ID,
		&img.PropertyID,
		&img.FileName,
		&img.ContentType,
		&img.Size,
		&img.IsMain,
		&img.CreatedAt,
	); err != nil {
		return nil, err
	}
	return img, nil
}
