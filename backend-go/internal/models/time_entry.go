package models

import (
	"encoding/json"
	"time"
)

// Date is a custom type for date-only values (YYYY-MM-DD format)
type Date struct {
	time.Time
}

// UnmarshalJSON implements json.Unmarshaler
func (d *Date) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}
	t, err := time.Parse("2006-01-02", str)
	if err != nil {
		return err
	}
	d.Time = t
	return nil
}

// MarshalJSON implements json.Marshaler
func (d Date) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.Time.Format("2006-01-02"))
}

type TimeEntry struct {
	ID          int       `json:"id"`
	PropertyID  *int      `json:"property_id"` // NULL for "General" activities
	Date        Date      `json:"date"`
	Category    string    `json:"category"`
	Description string    `json:"description"`
	Hours       float64   `json:"hours"`
	Notes       string    `json:"notes"`
	Mileage     float64   `json:"mileage"`
	FullDrive   bool      `json:"full_drive"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CreateTimeEntryRequest struct {
	PropertyID  *int    `json:"property_id"`
	Date        string  `json:"date"`
	Category    string  `json:"category"`
	Description string  `json:"description"`
	Hours       float64 `json:"hours"`
	Notes       string  `json:"notes"`
	Mileage     float64 `json:"mileage"`
	FullDrive   bool    `json:"full_drive"`
}

type UpdateTimeEntryRequest struct {
	PropertyID  *int    `json:"property_id"`
	Date        string  `json:"date"`
	Category    string  `json:"category"`
	Description string  `json:"description"`
	Hours       float64 `json:"hours"`
	Notes       string  `json:"notes"`
	Mileage     float64 `json:"mileage"`
	FullDrive   bool    `json:"full_drive"`
}

type TimeEntryFilter struct {
	PropertyID *int    `json:"property_id"`
	StartDate  *string `json:"start_date"`
	EndDate    *string `json:"end_date"`
	Category   *string `json:"category"`
}
