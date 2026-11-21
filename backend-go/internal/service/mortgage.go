package service

import (
	"strings"
	"time"

	"brickwise/backend/internal/errors"
	"brickwise/backend/internal/models"
	"brickwise/backend/internal/repository"
)

type MortgageService struct {
	repo         repository.MortgageRepository
	propertyRepo repository.PropertyRepository
}

func NewMortgageService(repo repository.MortgageRepository, propertyRepo repository.PropertyRepository) *MortgageService {
	return &MortgageService{repo: repo, propertyRepo: propertyRepo}
}

func (s *MortgageService) UpsertMortgage(propertyID int, req *models.UpsertMortgageRequest) (*models.Mortgage, error) {
	if _, err := s.propertyRepo.GetByID(propertyID); err != nil {
		return nil, errors.ErrNotFound
	}

	if req.CurrentBalance < 0 || req.OriginalLoan < 0 || req.MonthlyPayment < 0 || req.MonthlyPI < 0 {
		return nil, errors.New("INVALID_VALUE", "financial amounts must be non-negative", errors.ErrInvalidInput.Status)
	}

	var balanceAsOf *models.Date
	if req.BalanceAsOf != nil && strings.TrimSpace(*req.BalanceAsOf) != "" {
		parsed, err := time.Parse("2006-01-02", *req.BalanceAsOf)
		if err != nil {
			return nil, errors.New("INVALID_DATE", "balance_as_of must be in YYYY-MM-DD format", errors.ErrInvalidInput.Status)
		}
		balanceAsOf = &models.Date{Time: parsed}
	}

	mortgage := &models.Mortgage{
		PropertyID:     propertyID,
		OriginalLoan:   req.OriginalLoan,
		CurrentBalance: req.CurrentBalance,
		InterestRate:   req.InterestRate,
		MonthlyPayment: req.MonthlyPayment,
		MonthlyPI:      req.MonthlyPI,
		BalanceAsOf:    balanceAsOf,
		Lender:         req.Lender,
		Notes:          req.Notes,
	}

	if err := s.repo.Upsert(mortgage); err != nil {
		return nil, errors.Wrap(err, errors.ErrInternal)
	}

	return mortgage, nil
}

func (s *MortgageService) GetMortgage(propertyID int) (*models.Mortgage, error) {
	mortgage, err := s.repo.GetByPropertyID(propertyID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return nil, errors.ErrNotFound
		}
		return nil, errors.Wrap(err, errors.ErrInternal)
	}
	return mortgage, nil
}

func (s *MortgageService) DeleteMortgage(propertyID int) error {
	if err := s.repo.DeleteByPropertyID(propertyID); err != nil {
		if strings.Contains(err.Error(), "not found") {
			return errors.ErrNotFound
		}
		return errors.Wrap(err, errors.ErrInternal)
	}
	return nil
}
