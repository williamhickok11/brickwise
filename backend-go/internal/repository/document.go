package repository

import (
	"database/sql"
	"fmt"
	"os"

	"brickwise/backend/internal/models"
)

// DocumentRepository defines data access for residence documents.
type DocumentRepository interface {
	GetByResidenceID(residenceID int) (*models.Document, error)
	Upsert(doc *models.Document) error
	DeleteByResidenceID(residenceID int) error
}

type SQLDocumentRepository struct {
	db       *sql.DB
	isSQLite bool
}

func NewDocumentRepository(db *sql.DB) DocumentRepository {
	isSQLite := os.Getenv("DATABASE_URL") == ""
	return &SQLDocumentRepository{db: db, isSQLite: isSQLite}
}

func (r *SQLDocumentRepository) GetByResidenceID(residenceID int) (*models.Document, error) {
	query := `SELECT id, residence_id, file_name, content_type, size, created_at
	          FROM documents WHERE residence_id = $1`
	doc := &models.Document{}
	err := r.db.QueryRow(query, residenceID).Scan(
		&doc.ID,
		&doc.ResidenceID,
		&doc.FileName,
		&doc.ContentType,
		&doc.Size,
		&doc.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("document for residence %d not found", residenceID)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get document: %w", err)
	}
	return doc, nil
}

func (r *SQLDocumentRepository) Upsert(doc *models.Document) error {
	if r.isSQLite {
		// Replace any existing row for this residence.
		if _, err := r.db.Exec(`DELETE FROM documents WHERE residence_id = $1`, doc.ResidenceID); err != nil {
			return fmt.Errorf("failed to clear existing document: %w", err)
		}
		result, err := r.db.Exec(
			`INSERT INTO documents (residence_id, file_name, content_type, size)
			 VALUES ($1, $2, $3, $4)`,
			doc.ResidenceID, doc.FileName, doc.ContentType, doc.Size,
		)
		if err != nil {
			return fmt.Errorf("failed to insert document: %w", err)
		}
		id, err := result.LastInsertId()
		if err != nil {
			return fmt.Errorf("failed to get last insert id: %w", err)
		}
		doc.ID = int(id)
		return r.db.QueryRow(
			`SELECT created_at FROM documents WHERE id = $1`, doc.ID,
		).Scan(&doc.CreatedAt)
	}

	query := `INSERT INTO documents (residence_id, file_name, content_type, size)
	          VALUES ($1, $2, $3, $4)
	          ON CONFLICT (residence_id) DO UPDATE SET
	            file_name = EXCLUDED.file_name,
	            content_type = EXCLUDED.content_type,
	            size = EXCLUDED.size,
	            created_at = CURRENT_TIMESTAMP
	          RETURNING id, created_at`
	return r.db.QueryRow(query, doc.ResidenceID, doc.FileName, doc.ContentType, doc.Size).
		Scan(&doc.ID, &doc.CreatedAt)
}

func (r *SQLDocumentRepository) DeleteByResidenceID(residenceID int) error {
	_, err := r.db.Exec(`DELETE FROM documents WHERE residence_id = $1`, residenceID)
	if err != nil {
		return fmt.Errorf("failed to delete document: %w", err)
	}
	return nil
}
