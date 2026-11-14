package models

import "time"

// PropertyImage is metadata for a photo attached to a property.
type PropertyImage struct {
	ID          int       `json:"id"`
	PropertyID  int       `json:"property_id"`
	FileName    string    `json:"file_name"`
	ContentType string    `json:"content_type"`
	Size        int64     `json:"size"`
	IsMain      bool      `json:"is_main"`
	CreatedAt   time.Time `json:"created_at"`
}
