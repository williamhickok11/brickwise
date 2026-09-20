package service

import (
	"context"
	"time"

	"brickwise/backend/internal/errors"
	"brickwise/backend/internal/models"
	"brickwise/backend/internal/repository"
)

type DashboardService struct {
	propertyRepo  repository.PropertyRepository
	mortgageRepo  repository.MortgageRepository
	residenceRepo repository.ResidenceRepository
	qboService    *QBOService
}

func NewDashboardService(
	propertyRepo repository.PropertyRepository,
	mortgageRepo repository.MortgageRepository,
	residenceRepo repository.ResidenceRepository,
	qboService *QBOService,
) *DashboardService {
	return &DashboardService{
		propertyRepo:  propertyRepo,
		mortgageRepo:  mortgageRepo,
		residenceRepo: residenceRepo,
		qboService:    qboService,
	}
}

func (s *DashboardService) GetSummary() (*models.DashboardSummary, error) {
	properties, err := s.propertyRepo.GetAll()
	if err != nil {
		return nil, errors.Wrap(err, errors.ErrInternal)
	}

	mortgages, err := s.mortgageRepo.GetAll()
	if err != nil {
		return nil, errors.Wrap(err, errors.ErrInternal)
	}

	mortgageByProperty := make(map[int]*models.Mortgage, len(mortgages))
	for _, m := range mortgages {
		mortgageByProperty[m.PropertyID] = m
	}

	activeTrue := true
	activeResidences, err := s.residenceRepo.GetAll(&models.ResidenceFilter{IsActive: &activeTrue})
	if err != nil {
		return nil, errors.Wrap(err, errors.ErrInternal)
	}

	rentByProperty := make(map[int]float64)
	for _, r := range activeResidences {
		rentByProperty[r.PropertyID] += r.MonthlyRent
	}

	summary := &models.DashboardSummary{
		AsOf:       time.Now(),
		Properties: make([]models.PropertyBreakdown, 0, len(properties)),
	}

	actualByProperty := map[int]float64{}
	var qboSyncedAt *time.Time
	if s.qboService != nil {
		actuals, syncedAt, qboErr := s.qboService.CurrentMonthCashFlow(context.Background())
		if qboErr == nil {
			actualByProperty = actuals
			qboSyncedAt = syncedAt
			if len(actuals) > 0 || syncedAt != nil {
				summary.QBOConnected = true
			}
			// Connected even with zero mappings if connection exists — check via Status-ish path
		}
		if status, statusErr := s.qboService.Status(context.Background()); statusErr == nil {
			summary.QBOConnected = status.Connected
			if status.PeriodStart != nil {
				summary.CashFlowPeriodStart = status.PeriodStart
			}
			if status.PeriodEnd != nil {
				summary.CashFlowPeriodEnd = status.PeriodEnd
			}
			if qboSyncedAt == nil {
				qboSyncedAt = status.LastSyncedAt
			}
		}
	}
	summary.QBOSyncedAt = qboSyncedAt

	var actualTotal float64
	hasAnyActual := false

	for _, p := range properties {
		var mortgageBalance, monthlyPI float64
		if m, ok := mortgageByProperty[p.ID]; ok {
			mortgageBalance = m.CurrentBalance
			monthlyPI = m.MonthlyPI
		}

		monthlyRent := rentByProperty[p.ID]
		equity := p.CurrentValue - mortgageBalance
		netCashFlow := monthlyRent - monthlyPI

		summary.PortfolioValue += p.CurrentValue
		summary.TotalMortgageBalance += mortgageBalance
		summary.TotalEquity += equity
		summary.MonthlyRevenue += monthlyRent
		summary.MonthlyMortgagePI += monthlyPI
		summary.EstimatedCashFlow += netCashFlow

		breakdown := models.PropertyBreakdown{
			PropertyID:      p.ID,
			Name:            p.Name,
			CurrentValue:    p.CurrentValue,
			MortgageBalance: mortgageBalance,
			Equity:          equity,
			MonthlyRent:     monthlyRent,
			MonthlyPI:       monthlyPI,
			NetCashFlow:     netCashFlow,
			CashFlowSource:  "estimate",
		}

		if amount, ok := actualByProperty[p.ID]; ok {
			v := amount
			breakdown.ActualCashFlow = &v
			breakdown.CashFlowSource = "quickbooks"
			actualTotal += amount
			hasAnyActual = true
		}

		summary.Properties = append(summary.Properties, breakdown)
	}

	if hasAnyActual {
		summary.ActualCashFlow = &actualTotal
	}

	return summary, nil
}
