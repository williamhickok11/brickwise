package models

import "time"

// Document is metadata for a lease PDF attached to a residence.
// One document per residence (enforced by unique residence_id).
type Document struct {
	ID          int       `json:"id"`
	ResidenceID int       `json:"residence_id"`
	FileName    string    `json:"file_name"`
	ContentType string    `json:"content_type"`
	Size        int64     `json:"size"`
	CreatedAt   time.Time `json:"created_at"`
}
