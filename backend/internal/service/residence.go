package service

import (
	"strings"
	"time"

	"brickwise/backend/internal/errors"
	"brickwise/backend/internal/models"
	"brickwise/backend/internal/repository"
	"brickwise/backend/internal/validator"
)

// ResidenceService handles business logic for residences/tenants.
type ResidenceService struct {
	repo         repository.ResidenceRepository
	propertyRepo repository.PropertyRepository
}

func NewResidenceService(repo repository.ResidenceRepository, propertyRepo repository.PropertyRepository) *ResidenceService {
	return &ResidenceService{
		repo:          repo,
		propertyRepo: propertyRepo,
	}
}

func (s *ResidenceService) CreateResidence(req *models.CreateResidenceRequest) (*models.Residence, error) {
	if err := validator.ValidateResidenceCreate(req); err != nil {
		return nil, validator.ToAPIError(err)
	}

	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		return nil, errors.New("INVALID_DATE", "start_date must be in YYYY-MM-DD format", errors.ErrInvalidInput.Status)
	}

	var endDate *models.Date
	if req.EndDate != nil && strings.TrimSpace(*req.EndDate) != "" {
		parsedEnd, err := time.Parse("2006-01-02", *req.EndDate)
		if err != nil {
			return nil, errors.New("INVALID_DATE", "end_date must be in YYYY-MM-DD format", errors.ErrInvalidInput.Status)
		}
		d := models.Date{Time: parsedEnd}
		endDate = &d
	}

	// Validate referenced property exists.
	if _, err := s.propertyRepo.GetByID(req.PropertyID); err != nil {
		return nil, errors.ErrNotFound
	}

	residence := &models.Residence{
		PropertyID: req.PropertyID,
		Name:       req.Name,
		Phone:      req.Phone,
		Email:      req.Email,
		StartDate:  models.Date{Time: startDate},
		EndDate:    endDate,
		IsActive:   req.IsActive,
	}

	if err := s.repo.Create(residence); err != nil {
		return nil, errors.Wrap(err, errors.ErrInternal)
	}

	return residence, nil
}

func (s *ResidenceService) GetResidence(id int) (*models.Residence, error) {
	residence, err := s.repo.GetByID(id)
	if err != nil {
		return nil, errors.Wrap(err, errors.ErrNotFound)
	}
	return residence, nil
}

func (s *ResidenceService) ListResidences(filter *models.ResidenceFilter) ([]*models.Residence, error) {
	residences, err := s.repo.GetAll(filter)
	if err != nil {
		return nil, errors.Wrap(err, errors.ErrInternal)
	}
	return residences, nil
}

func (s *ResidenceService) UpdateResidence(id int, req *models.UpdateResidenceRequest) (*models.Residence, error) {
	if err := validator.ValidateResidenceUpdate(req); err != nil {
		return nil, validator.ToAPIError(err)
	}

	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		return nil, errors.New("INVALID_DATE", "start_date must be in YYYY-MM-DD format", errors.ErrInvalidInput.Status)
	}

	var endDate *models.Date
	if req.EndDate != nil && strings.TrimSpace(*req.EndDate) != "" {
		parsedEnd, err := time.Parse("2006-01-02", *req.EndDate)
		if err != nil {
			return nil, errors.New("INVALID_DATE", "end_date must be in YYYY-MM-DD format", errors.ErrInvalidInput.Status)
		}
		d := models.Date{Time: parsedEnd}
		endDate = &d
	}

	if _, err := s.propertyRepo.GetByID(req.PropertyID); err != nil {
		return nil, errors.ErrNotFound
	}

	residence := &models.Residence{
		ID:         id,
		PropertyID: req.PropertyID,
		Name:       req.Name,
		Phone:      req.Phone,
		Email:      req.Email,
		StartDate:  models.Date{Time: startDate},
		EndDate:    endDate,
		IsActive:   req.IsActive,
	}

	if err := s.repo.Update(residence); err != nil {
		if strings.Contains(err.Error(), "not found") {
			return nil, errors.ErrNotFound
		}
		return nil, errors.Wrap(err, errors.ErrInternal)
	}

	return s.GetResidence(id)
}

