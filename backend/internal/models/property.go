package models

import "time"

type Property struct {
	ID            int       `json:"id"`
	Name          string    `json:"name"`
	Address       string    `json:"address"`
	PropertyType  string    `json:"property_type"`
	DefaultMileage float64  `json:"default_mileage"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type CreatePropertyRequest struct {
	Name          string  `json:"name"`
	Address       string  `json:"address"`
	PropertyType  string  `json:"property_type"`
	DefaultMileage float64 `json:"default_mileage"`
}
