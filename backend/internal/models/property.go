package models

import "time"

type Property struct {
	ID             int        `json:"id"`
	Name           string     `json:"name"`
	Address        string     `json:"address"`
	PropertyType   string     `json:"property_type"`
	DefaultMileage float64    `json:"default_mileage"`
	CurrentValue   float64    `json:"current_value"`
	ValueUpdatedAt *time.Time `json:"value_updated_at"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

type UpdatePropertyValueRequest struct {
	CurrentValue   float64  `json:"current_value"`
	ValueUpdatedAt *string  `json:"value_updated_at"`
}

type CreatePropertyRequest struct {
	Name          string  `json:"name"`
	Address       string  `json:"address"`
	PropertyType  string  `json:"property_type"`
	DefaultMileage float64 `json:"default_mileage"`
}
