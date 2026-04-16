package models

import "time"

// Residence represents a tenant/lease occupant history over time for a specific property.
// v1 uses `is_active` to separate current vs former residents.
type Residence struct {
	ID         int      `json:"id"`
	PropertyID int      `json:"property_id"`
	Name       string   `json:"name"`
	Phone      string   `json:"phone"`
	Email      string   `json:"email"`
	StartDate  Date     `json:"start_date"`
	EndDate    *Date    `json:"end_date"`
	IsActive   bool     `json:"is_active"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type CreateResidenceRequest struct {
	PropertyID int     `json:"property_id"`
	Name       string  `json:"name"`
	Phone      string  `json:"phone"`
	Email      string  `json:"email"`
	StartDate  string  `json:"start_date"`
	EndDate    *string `json:"end_date"`
	IsActive   bool    `json:"is_active"`
}

type UpdateResidenceRequest struct {
	PropertyID int     `json:"property_id"`
	Name       string  `json:"name"`
	Phone      string  `json:"phone"`
	Email      string  `json:"email"`
	StartDate  string  `json:"start_date"`
	EndDate    *string `json:"end_date"`
	IsActive   bool    `json:"is_active"`
}

type ResidenceFilter struct {
	PropertyID *int  `json:"property_id"`
	IsActive   *bool `json:"is_active"`
}

