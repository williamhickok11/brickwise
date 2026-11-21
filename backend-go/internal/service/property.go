package service

import (
	"strings"
	"time"

	"brickwise/backend/internal/errors"
	"brickwise/backend/internal/models"
	"brickwise/backend/internal/repository"
	"brickwise/backend/internal/validator"
)

// PropertyService handles business logic for properties
type PropertyService struct {
	repo repository.PropertyRepository
}

// NewPropertyService creates a new property service
func NewPropertyService(repo repository.PropertyRepository) *PropertyService {
	return &PropertyService{repo: repo}
}

// CreateProperty creates a new property with validation
func (s *PropertyService) CreateProperty(req *models.CreatePropertyRequest) (*models.Property, error) {
	// Validate input
	if err := validator.ValidateProperty(req.Name, req.Address, req.PropertyType); err != nil {
		return nil, validator.ToAPIError(err)
	}

	property := &models.Property{
		Name:          req.Name,
		Address:       req.Address,
		PropertyType:  req.PropertyType,
		DefaultMileage: req.DefaultMileage,
	}

	if err := s.repo.Create(property); err != nil {
		return nil, errors.Wrap(err, errors.ErrInternal)
	}

	return property, nil
}

// GetProperty retrieves a property by ID
func (s *PropertyService) GetProperty(id int) (*models.Property, error) {
	property, err := s.repo.GetByID(id)
	if err != nil {
		return nil, errors.Wrap(err, errors.ErrNotFound)
	}
	return property, nil
}

// ListProperties retrieves all properties
func (s *PropertyService) ListProperties() ([]*models.Property, error) {
	properties, err := s.repo.GetAll()
	if err != nil {
		return nil, errors.Wrap(err, errors.ErrInternal)
	}
	return properties, nil
}

// UpdateProperty updates an existing property
func (s *PropertyService) UpdateProperty(id int, req *models.CreatePropertyRequest) (*models.Property, error) {
	// Validate input
	if err := validator.ValidateProperty(req.Name, req.Address, req.PropertyType); err != nil {
		return nil, validator.ToAPIError(err)
	}

	property := &models.Property{
		ID:            id,
		Name:          req.Name,
		Address:       req.Address,
		PropertyType:  req.PropertyType,
		DefaultMileage: req.DefaultMileage,
	}

	if err := s.repo.Update(property); err != nil {
		// Check if it's a not found error
		if strings.Contains(err.Error(), "not found") {
			return nil, errors.ErrNotFound
		}
		return nil, errors.Wrap(err, errors.ErrInternal)
	}

	// Fetch updated property to return complete data
	return s.GetProperty(id)
}

// DeleteProperty deletes a property by ID
func (s *PropertyService) DeleteProperty(id int) error {
	if err := s.repo.Delete(id); err != nil {
		if strings.Contains(err.Error(), "not found") {
			return errors.ErrNotFound
		}
		return errors.Wrap(err, errors.ErrInternal)
	}
	return nil
}

// UpdatePropertyValue updates only the current value fields on a property.
func (s *PropertyService) UpdatePropertyValue(id int, req *models.UpdatePropertyValueRequest) (*models.Property, error) {
	if req.CurrentValue < 0 {
		return nil, errors.New("INVALID_VALUE", "current_value must be non-negative", errors.ErrInvalidInput.Status)
	}

	var valueUpdatedAt *time.Time
	if req.ValueUpdatedAt != nil && strings.TrimSpace(*req.ValueUpdatedAt) != "" {
		parsed, err := time.Parse("2006-01-02", *req.ValueUpdatedAt)
		if err != nil {
			return nil, errors.New("INVALID_DATE", "value_updated_at must be in YYYY-MM-DD format", errors.ErrInvalidInput.Status)
		}
		valueUpdatedAt = &parsed
	} else {
		now := time.Now()
		valueUpdatedAt = &now
	}

	if err := s.repo.UpdateValue(id, req.CurrentValue, valueUpdatedAt); err != nil {
		if strings.Contains(err.Error(), "not found") {
			return nil, errors.ErrNotFound
		}
		return nil, errors.Wrap(err, errors.ErrInternal)
	}

	return s.GetProperty(id)
}
