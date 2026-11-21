package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/url"
	"strings"
	"time"

	"brickwise/backend/internal/errors"
	"brickwise/backend/internal/models"
	"brickwise/backend/internal/qbo"
	"brickwise/backend/internal/repository"
)

const qboCacheTTL = 6 * time.Hour

// QBOService manages QuickBooks OAuth, Class mappings, and cash-flow sync.
type QBOService struct {
	repo         repository.QBORepository
	propertyRepo repository.PropertyRepository
	client       *qbo.Client
	cfg          qbo.Config
}

func NewQBOService(
	repo repository.QBORepository,
	propertyRepo repository.PropertyRepository,
	cfg qbo.Config,
	client *qbo.Client,
) *QBOService {
	return &QBOService{
		repo:         repo,
		propertyRepo: propertyRepo,
		client:       client,
		cfg:          cfg,
	}
}

func (s *QBOService) Status(ctx context.Context) (*models.QBOStatus, error) {
	status := &models.QBOStatus{
		Configured:    s.cfg.Configured(),
		ClassRequired: true,
		Mappings:      []models.PropertyQBOLink{},
		Message:       "Map each property to a QuickBooks Class. Your accounting team must code income and expenses to that Class.",
	}

	if !status.Configured {
		status.Message = "Set QBO_CLIENT_ID and QBO_CLIENT_SECRET to enable QuickBooks Online."
		return status, nil
	}

	conn, err := s.repo.GetConnection()
	if err != nil {
		return nil, errors.Wrap(err, errors.ErrInternal)
	}
	if conn == nil {
		return status, nil
	}

	status.Connected = true
	status.RealmID = conn.RealmID
	status.ConnectedAt = &conn.ConnectedAt

	mappings, err := s.repo.ListMappings()
	if err != nil {
		return nil, errors.Wrap(err, errors.ErrInternal)
	}
	if mappings != nil {
		status.Mappings = mappings
	}

	syncedAt, err := s.repo.LatestSyncTime()
	if err != nil {
		return nil, errors.Wrap(err, errors.ErrInternal)
	}
	status.LastSyncedAt = syncedAt

	start, end := currentMonthBounds(time.Now())
	status.PeriodStart = &start
	status.PeriodEnd = &end

	return status, nil
}

// StartConnect returns the Intuit authorize URL after storing CSRF state.
func (s *QBOService) StartConnect(ctx context.Context) (string, error) {
	if !s.cfg.Configured() {
		return "", errors.New("QBO_NOT_CONFIGURED", "QuickBooks credentials are not configured", 400)
	}

	state, err := randomState(24)
	if err != nil {
		return "", errors.Wrap(err, errors.ErrInternal)
	}
	if err := s.repo.SaveOAuthState(state, time.Now().Add(15*time.Minute)); err != nil {
		return "", errors.Wrap(err, errors.ErrInternal)
	}
	return s.client.AuthorizeURL(state), nil
}

// HandleCallback exchanges the OAuth code and stores tokens.
func (s *QBOService) HandleCallback(ctx context.Context, code, state, realmID string) error {
	if !s.cfg.Configured() {
		return errors.New("QBO_NOT_CONFIGURED", "QuickBooks credentials are not configured", 400)
	}
	if code == "" || state == "" || realmID == "" {
		return errors.New("INVALID_INPUT", "Missing OAuth callback parameters", 400)
	}

	ok, err := s.repo.ConsumeOAuthState(state)
	if err != nil {
		return errors.Wrap(err, errors.ErrInternal)
	}
	if !ok {
		return errors.New("INVALID_STATE", "Invalid or expired OAuth state", 400)
	}

	tokens, err := s.client.ExchangeCode(ctx, code)
	if err != nil {
		return errors.Wrap(err, errors.ErrInternal)
	}

	now := time.Now()
	conn := &models.QBOConnection{
		RealmID:     realmID,
		ConnectedAt: now,
	}
	qbo.ApplyTokens(conn, tokens, now)

	if err := s.repo.UpsertConnection(conn); err != nil {
		return errors.Wrap(err, errors.ErrInternal)
	}
	return nil
}

func (s *QBOService) FrontendRedirectURL(success bool, errMsg string) string {
	base := s.cfg.FrontendRedirect
	u, err := url.Parse(base)
	if err != nil {
		return base
	}
	q := u.Query()
	if success {
		q.Set("connected", "1")
	} else {
		q.Set("error", errMsg)
	}
	u.RawQuery = q.Encode()
	return u.String()
}

func (s *QBOService) Disconnect(ctx context.Context) error {
	if err := s.repo.DeleteConnection(); err != nil {
		return errors.Wrap(err, errors.ErrInternal)
	}
	_ = s.repo.DeleteAllMappings()
	return nil
}

func (s *QBOService) ListClasses(ctx context.Context) ([]models.QBOClass, error) {
	conn, err := s.ensureFreshToken(ctx)
	if err != nil {
		return nil, err
	}
	classes, err := s.client.ListClasses(ctx, conn.RealmID, conn.AccessToken)
	if err != nil {
		return nil, errors.Wrap(err, errors.ErrInternal)
	}
	return classes, nil
}

func (s *QBOService) UpdateMappings(ctx context.Context, req *models.UpdateQBOMappingsRequest) error {
	if req == nil {
		return errors.ErrInvalidInput
	}

	properties, err := s.propertyRepo.GetAll()
	if err != nil {
		return errors.Wrap(err, errors.ErrInternal)
	}
	validIDs := make(map[int]bool, len(properties))
	for _, p := range properties {
		validIDs[p.ID] = true
	}

	keep := make(map[int]bool)
	for _, m := range req.Mappings {
		if !validIDs[m.PropertyID] {
			return errors.New("INVALID_INPUT", fmt.Sprintf("Unknown property_id %d", m.PropertyID), 400)
		}
		classID := strings.TrimSpace(m.QBOClassID)
		if classID == "" {
			continue
		}
		link := &models.PropertyQBOLink{
			PropertyID:   m.PropertyID,
			QBOClassID:   classID,
			QBOClassName: strings.TrimSpace(m.QBOClassName),
		}
		if err := s.repo.UpsertMapping(link); err != nil {
			return errors.Wrap(err, errors.ErrInternal)
		}
		keep[m.PropertyID] = true
	}

	existing, err := s.repo.ListMappings()
	if err != nil {
		return errors.Wrap(err, errors.ErrInternal)
	}
	for _, e := range existing {
		if !keep[e.PropertyID] {
			if err := s.repo.DeleteMapping(e.PropertyID); err != nil {
				return errors.Wrap(err, errors.ErrInternal)
			}
		}
	}

	return nil
}

// SyncCashFlow pulls cash-basis P&L by Class for the current month.
func (s *QBOService) SyncCashFlow(ctx context.Context, force bool) (*models.QBOSyncResult, error) {
	conn, err := s.ensureFreshToken(ctx)
	if err != nil {
		return nil, err
	}

	start, end := currentMonthBounds(time.Now())
	if !force {
		cached, err := s.repo.ListCashFlowForPeriod(start, end)
		if err != nil {
			return nil, errors.Wrap(err, errors.ErrInternal)
		}
		if len(cached) > 0 {
			newest := cached[0].SyncedAt
			for _, c := range cached[1:] {
				if c.SyncedAt.After(newest) {
					newest = c.SyncedAt
				}
			}
			if time.Since(newest) < qboCacheTTL {
				return &models.QBOSyncResult{
					SyncedAt:    newest,
					PeriodStart: start,
					PeriodEnd:   end,
					Updated:     len(cached),
				}, nil
			}
		}
	}

	byClassName, err := s.client.ClassNetIncome(ctx, conn.RealmID, conn.AccessToken, start, end)
	if err != nil {
		return nil, errors.Wrap(err, errors.ErrInternal)
	}

	mappings, err := s.repo.ListMappings()
	if err != nil {
		return nil, errors.Wrap(err, errors.ErrInternal)
	}

	now := time.Now()
	updated := 0
	for _, m := range mappings {
		amount := lookupClassAmount(byClassName, m.QBOClassName, m.QBOClassID)
		cache := &models.QBOCashFlowCache{
			PropertyID:  m.PropertyID,
			PeriodStart: start,
			PeriodEnd:   end,
			NetAmount:   amount,
			SyncedAt:    now,
		}
		if err := s.repo.UpsertCashFlow(cache); err != nil {
			return nil, errors.Wrap(err, errors.ErrInternal)
		}
		updated++
	}

	return &models.QBOSyncResult{
		SyncedAt:    now,
		PeriodStart: start,
		PeriodEnd:   end,
		Updated:     updated,
	}, nil
}

// CurrentMonthCashFlow returns cached nets keyed by property ID (refreshing if stale).
func (s *QBOService) CurrentMonthCashFlow(ctx context.Context) (map[int]float64, *time.Time, error) {
	conn, err := s.repo.GetConnection()
	if err != nil {
		return nil, nil, errors.Wrap(err, errors.ErrInternal)
	}
	if conn == nil {
		return map[int]float64{}, nil, nil
	}

	if _, err := s.SyncCashFlow(ctx, false); err != nil {
		// Fall back to whatever is cached rather than failing the dashboard.
		start, end := currentMonthBounds(time.Now())
		cached, cacheErr := s.repo.ListCashFlowForPeriod(start, end)
		if cacheErr != nil {
			return nil, nil, errors.Wrap(err, errors.ErrInternal)
		}
		out := make(map[int]float64, len(cached))
		var syncedAt *time.Time
		for _, c := range cached {
			out[c.PropertyID] = c.NetAmount
			t := c.SyncedAt
			if syncedAt == nil || t.After(*syncedAt) {
				syncedAt = &t
			}
		}
		return out, syncedAt, nil
	}

	start, end := currentMonthBounds(time.Now())
	cached, err := s.repo.ListCashFlowForPeriod(start, end)
	if err != nil {
		return nil, nil, errors.Wrap(err, errors.ErrInternal)
	}
	out := make(map[int]float64, len(cached))
	var syncedAt *time.Time
	for _, c := range cached {
		out[c.PropertyID] = c.NetAmount
		t := c.SyncedAt
		if syncedAt == nil || t.After(*syncedAt) {
			syncedAt = &t
		}
	}
	return out, syncedAt, nil
}

func (s *QBOService) ensureFreshToken(ctx context.Context) (*models.QBOConnection, error) {
	if !s.cfg.Configured() {
		return nil, errors.New("QBO_NOT_CONFIGURED", "QuickBooks credentials are not configured", 400)
	}
	conn, err := s.repo.GetConnection()
	if err != nil {
		return nil, errors.Wrap(err, errors.ErrInternal)
	}
	if conn == nil {
		return nil, errors.New("QBO_NOT_CONNECTED", "QuickBooks is not connected", 400)
	}

	// Refresh a few minutes before expiry.
	if time.Until(conn.AccessTokenExpiresAt) > 3*time.Minute {
		return conn, nil
	}

	tokens, err := s.client.RefreshAccessToken(ctx, conn.RefreshToken)
	if err != nil {
		return nil, errors.Wrap(err, errors.ErrInternal)
	}
	qbo.ApplyTokens(conn, tokens, time.Now())
	if err := s.repo.UpsertConnection(conn); err != nil {
		return nil, errors.Wrap(err, errors.ErrInternal)
	}
	return conn, nil
}

func lookupClassAmount(byName map[string]float64, className, classID string) float64 {
	if className != "" {
		if v, ok := byName[className]; ok {
			return v
		}
		// Case-insensitive fallback
		for k, v := range byName {
			if strings.EqualFold(k, className) {
				return v
			}
		}
	}
	if classID != "" {
		if v, ok := byName[classID]; ok {
			return v
		}
	}
	return 0
}

func currentMonthBounds(now time.Time) (time.Time, time.Time) {
	loc := now.Location()
	start := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, loc)
	end := start.AddDate(0, 1, -1)
	return start, end
}

func randomState(nBytes int) (string, error) {
	b := make([]byte, nBytes)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
