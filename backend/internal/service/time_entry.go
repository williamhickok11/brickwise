package service

import (
	"strings"
	"time"

	"brickwise/backend/internal/errors"
	"brickwise/backend/internal/models"
	"brickwise/backend/internal/repository"
	"brickwise/backend/internal/validator"
)

// TimeEntryService handles business logic for time entries
type TimeEntryService struct {
	repo            repository.TimeEntryRepository
	propertyRepo    repository.PropertyRepository
}

// NewTimeEntryService creates a new time entry service
func NewTimeEntryService(repo repository.TimeEntryRepository, propertyRepo repository.PropertyRepository) *TimeEntryService {
	return &TimeEntryService{
		repo:         repo,
		propertyRepo: propertyRepo,
	}
}

// CreateTimeEntry creates a new time entry with validation
func (s *TimeEntryService) CreateTimeEntry(req *models.CreateTimeEntryRequest) (*models.TimeEntry, error) {
	// Validate input
	if err := validator.ValidateTimeEntry(req); err != nil {
		return nil, validator.ToAPIError(err)
	}

	// Parse date
	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return nil, errors.New("INVALID_DATE", "date must be in YYYY-MM-DD format", errors.ErrInvalidInput.Status)
	}

	entry := &models.TimeEntry{
		PropertyID:  req.PropertyID,
		Date:        models.Date{Time: date},
		Category:    req.Category,
		Description: req.Description,
		Hours:       req.Hours,
		Notes:       req.Notes,
		Mileage:     req.Mileage,
		FullDrive:   req.FullDrive,
	}

	// Auto-calculate mileage if full_drive is true and property_id is set
	if entry.FullDrive && entry.PropertyID != nil {
		property, err := s.propertyRepo.GetByID(*entry.PropertyID)
		if err == nil && property.DefaultMileage > 0 {
			entry.Mileage = property.DefaultMileage
		}
	}

	if err := s.repo.Create(entry); err != nil {
		return nil, errors.Wrap(err, errors.ErrInternal)
	}

	return entry, nil
}

// GetTimeEntry retrieves a time entry by ID
func (s *TimeEntryService) GetTimeEntry(id int) (*models.TimeEntry, error) {
	entry, err := s.repo.GetByID(id)
	if err != nil {
		return nil, errors.Wrap(err, errors.ErrNotFound)
	}
	return entry, nil
}

// ListTimeEntries retrieves all time entries with optional filters
func (s *TimeEntryService) ListTimeEntries(filter *models.TimeEntryFilter) ([]*models.TimeEntry, error) {
	entries, err := s.repo.GetAll(filter)
	if err != nil {
		return nil, errors.Wrap(err, errors.ErrInternal)
	}
	return entries, nil
}

// UpdateTimeEntry updates an existing time entry
func (s *TimeEntryService) UpdateTimeEntry(id int, req *models.UpdateTimeEntryRequest) (*models.TimeEntry, error) {
	// Validate input
	if err := validator.ValidateTimeEntryUpdate(req); err != nil {
		return nil, validator.ToAPIError(err)
	}

	// Parse date
	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return nil, errors.New("INVALID_DATE", "date must be in YYYY-MM-DD format", errors.ErrInvalidInput.Status)
	}

	entry := &models.TimeEntry{
		ID:          id,
		PropertyID:  req.PropertyID,
		Date:        models.Date{Time: date},
		Category:    req.Category,
		Description: req.Description,
		Hours:       req.Hours,
		Notes:       req.Notes,
		Mileage:     req.Mileage,
		FullDrive:   req.FullDrive,
	}

	// Auto-calculate mileage if full_drive is true and property_id is set
	if entry.FullDrive && entry.PropertyID != nil {
		property, err := s.propertyRepo.GetByID(*entry.PropertyID)
		if err == nil && property.DefaultMileage > 0 {
			entry.Mileage = property.DefaultMileage
		}
	}

	if err := s.repo.Update(entry); err != nil {
		if strings.Contains(err.Error(), "not found") {
			return nil, errors.ErrNotFound
		}
		return nil, errors.Wrap(err, errors.ErrInternal)
	}

	// Fetch updated entry to return complete data
	return s.GetTimeEntry(id)
}

// DeleteTimeEntry deletes a time entry by ID
func (s *TimeEntryService) DeleteTimeEntry(id int) error {
	if err := s.repo.Delete(id); err != nil {
		if strings.Contains(err.Error(), "not found") {
			return errors.ErrNotFound
		}
		return errors.Wrap(err, errors.ErrInternal)
	}
	return nil
}
