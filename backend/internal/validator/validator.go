package validator

import (
	"fmt"
	"strings"
	"time"

	"brickwise/backend/internal/errors"
	"brickwise/backend/internal/models"
)

// ValidationError represents a field validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ValidationErrors is a collection of validation errors
type ValidationErrors struct {
	Errors []ValidationError `json:"errors"`
}

func (v *ValidationErrors) Error() string {
	var msgs []string
	for _, err := range v.Errors {
		msgs = append(msgs, fmt.Sprintf("%s: %s", err.Field, err.Message))
	}
	return strings.Join(msgs, ", ")
}

// ValidateProperty validates property creation/update requests
func ValidateProperty(name, address, propertyType string) error {
	var errs ValidationErrors

	if strings.TrimSpace(name) == "" {
		errs.Errors = append(errs.Errors, ValidationError{
			Field:   "name",
			Message: "name is required",
		})
	}

	if len(name) > 255 {
		errs.Errors = append(errs.Errors, ValidationError{
			Field:   "name",
			Message: "name must be 255 characters or less",
		})
	}

	if strings.TrimSpace(address) == "" {
		errs.Errors = append(errs.Errors, ValidationError{
			Field:   "address",
			Message: "address is required",
		})
	}

	if len(address) > 500 {
		errs.Errors = append(errs.Errors, ValidationError{
			Field:   "address",
			Message: "address must be 500 characters or less",
		})
	}

	validTypes := map[string]bool{
		"residential": true,
		"commercial":  true,
		"industrial":  true,
		"land":        true,
	}

	if strings.TrimSpace(propertyType) == "" {
		errs.Errors = append(errs.Errors, ValidationError{
			Field:   "property_type",
			Message: "property_type is required",
		})
	} else if !validTypes[strings.ToLower(propertyType)] {
		errs.Errors = append(errs.Errors, ValidationError{
			Field:   "property_type",
			Message: fmt.Sprintf("property_type must be one of: %s", strings.Join([]string{"residential", "commercial", "industrial", "land"}, ", ")),
		})
	}

	if len(errs.Errors) > 0 {
		return &errs
	}

	return nil
}

// ValidateResidenceCreate validates residence creation/update requests.
func ValidateResidenceCreate(req *models.CreateResidenceRequest) error {
	var errs ValidationErrors

	if req.PropertyID <= 0 {
		errs.Errors = append(errs.Errors, ValidationError{
			Field:   "property_id",
			Message: "property_id is required",
		})
	}

	if strings.TrimSpace(req.Name) == "" {
		errs.Errors = append(errs.Errors, ValidationError{
			Field:   "name",
			Message: "name is required",
		})
	} else if len(req.Name) > 255 {
		errs.Errors = append(errs.Errors, ValidationError{
			Field:   "name",
			Message: "name must be 255 characters or less",
		})
	}

	if strings.TrimSpace(req.StartDate) == "" {
		errs.Errors = append(errs.Errors, ValidationError{
			Field:   "start_date",
			Message: "start_date is required",
		})
	} else if _, err := time.Parse("2006-01-02", req.StartDate); err != nil {
		errs.Errors = append(errs.Errors, ValidationError{
			Field:   "start_date",
			Message: "start_date must be in YYYY-MM-DD format",
		})
	}

	if req.EndDate != nil {
		if strings.TrimSpace(*req.EndDate) == "" {
			// Allow explicit empty string as "not set"
		} else if _, err := time.Parse("2006-01-02", *req.EndDate); err != nil {
			errs.Errors = append(errs.Errors, ValidationError{
				Field:   "end_date",
				Message: "end_date must be in YYYY-MM-DD format",
			})
		}
	}

	if len(req.Phone) > 255 {
		errs.Errors = append(errs.Errors, ValidationError{
			Field:   "phone",
			Message: "phone must be 255 characters or less",
		})
	}

	if len(req.Email) > 255 {
		errs.Errors = append(errs.Errors, ValidationError{
			Field:   "email",
			Message: "email must be 255 characters or less",
		})
	}

	if len(errs.Errors) > 0 {
		return &errs
	}
	return nil
}

func ValidateResidenceUpdate(req *models.UpdateResidenceRequest) error {
	// Same constraints for v1; move-out is modeled as update with `is_active=false` and optional `end_date`.
	return ValidateResidenceCreate(&models.CreateResidenceRequest{
		PropertyID: req.PropertyID,
		Name:       req.Name,
		Phone:      req.Phone,
		Email:      req.Email,
		StartDate:  req.StartDate,
		EndDate:    req.EndDate,
		IsActive:   req.IsActive,
	})
}

// ToAPIError converts ValidationErrors to an APIError
func ToAPIError(err error) *errors.APIError {
	if ve, ok := err.(*ValidationErrors); ok {
		return errors.New("VALIDATION_ERROR", ve.Error(), errors.ErrInvalidInput.Status)
	}
	return errors.ErrInvalidInput
}

// Valid categories for time entries
var validCategories = map[string]bool{
	"Property Management":   true,
	"Maintenance & Repairs": true,
	"Contractor Oversight":  true,
	"Accounting & Admin":    true,
	"Deal Sourcing":         true,
	"Construction Oversight": true,
	"Software Management":   true,
}

// ValidateTimeEntry validates time entry creation requests
func ValidateTimeEntry(req interface{}) error {
	var errs ValidationErrors
	var date, category, description string
	var hours float64

	// Type assertion based on request type
	switch v := req.(type) {
	case *models.CreateTimeEntryRequest:
		date = v.Date
		category = v.Category
		description = v.Description
		hours = v.Hours
	case *models.UpdateTimeEntryRequest:
		date = v.Date
		category = v.Category
		description = v.Description
		hours = v.Hours
	default:
		return &ValidationErrors{Errors: []ValidationError{
			{Field: "request", Message: "invalid request type"},
		}}
	}

	if strings.TrimSpace(date) == "" {
		errs.Errors = append(errs.Errors, ValidationError{
			Field:   "date",
			Message: "date is required",
		})
	}

	if strings.TrimSpace(category) == "" {
		errs.Errors = append(errs.Errors, ValidationError{
			Field:   "category",
			Message: "category is required",
		})
	} else if !validCategories[category] {
		categories := []string{}
		for cat := range validCategories {
			categories = append(categories, cat)
		}
		errs.Errors = append(errs.Errors, ValidationError{
			Field:   "category",
			Message: fmt.Sprintf("category must be one of: %s", strings.Join(categories, ", ")),
		})
	}

	if strings.TrimSpace(description) == "" {
		errs.Errors = append(errs.Errors, ValidationError{
			Field:   "description",
			Message: "description is required",
		})
	}

	if len(description) > 2000 {
		errs.Errors = append(errs.Errors, ValidationError{
			Field:   "description",
			Message: "description must be 2000 characters or less",
		})
	}

	if hours <= 0 {
		errs.Errors = append(errs.Errors, ValidationError{
			Field:   "hours",
			Message: "hours must be greater than 0",
		})
	}

	if hours > 24 {
		errs.Errors = append(errs.Errors, ValidationError{
			Field:   "hours",
			Message: "hours must be 24 or less",
		})
	}

	if len(errs.Errors) > 0 {
		return &errs
	}

	return nil
}

// ValidateTimeEntryUpdate validates time entry update requests
func ValidateTimeEntryUpdate(req *models.UpdateTimeEntryRequest) error {
	return ValidateTimeEntry(req)
}
