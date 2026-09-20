package repository

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"brickwise/backend/internal/models"
)

// QBORepository persists QuickBooks connection, mappings, and cash-flow cache.
type QBORepository interface {
	GetConnection() (*models.QBOConnection, error)
	UpsertConnection(conn *models.QBOConnection) error
	DeleteConnection() error

	SaveOAuthState(state string, expiresAt time.Time) error
	ConsumeOAuthState(state string) (bool, error)

	ListMappings() ([]models.PropertyQBOLink, error)
	UpsertMapping(link *models.PropertyQBOLink) error
	DeleteMapping(propertyID int) error
	DeleteAllMappings() error

	UpsertCashFlow(cache *models.QBOCashFlowCache) error
	ListCashFlowForPeriod(periodStart, periodEnd time.Time) ([]models.QBOCashFlowCache, error)
	LatestSyncTime() (*time.Time, error)
}

type SQLQBORepository struct {
	db       *sql.DB
	isSQLite bool
}

func NewQBORepository(db *sql.DB) QBORepository {
	return &SQLQBORepository{
		db:       db,
		isSQLite: os.Getenv("DATABASE_URL") == "",
	}
}

func (r *SQLQBORepository) GetConnection() (*models.QBOConnection, error) {
	conn := &models.QBOConnection{}
	err := r.db.QueryRow(`
		SELECT id, realm_id, access_token, refresh_token,
			access_token_expires_at, refresh_token_expires_at, connected_at, updated_at
		FROM qbo_connections
		ORDER BY id ASC
		LIMIT 1`).Scan(
		&conn.ID,
		&conn.RealmID,
		&conn.AccessToken,
		&conn.RefreshToken,
		&conn.AccessTokenExpiresAt,
		&conn.RefreshTokenExpiresAt,
		&conn.ConnectedAt,
		&conn.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get qbo connection: %w", err)
	}
	return conn, nil
}

func (r *SQLQBORepository) UpsertConnection(conn *models.QBOConnection) error {
	existing, err := r.GetConnection()
	if err != nil {
		return err
	}

	if existing == nil {
		if r.isSQLite {
			res, err := r.db.Exec(`
				INSERT INTO qbo_connections (
					realm_id, access_token, refresh_token,
					access_token_expires_at, refresh_token_expires_at, connected_at, updated_at
				) VALUES ($1, $2, $3, $4, $5, $6, $7)`,
				conn.RealmID, conn.AccessToken, conn.RefreshToken,
				conn.AccessTokenExpiresAt, conn.RefreshTokenExpiresAt,
				conn.ConnectedAt, conn.UpdatedAt,
			)
			if err != nil {
				return fmt.Errorf("insert qbo connection: %w", err)
			}
			id, err := res.LastInsertId()
			if err != nil {
				return err
			}
			conn.ID = int(id)
			return nil
		}
		return r.db.QueryRow(`
			INSERT INTO qbo_connections (
				realm_id, access_token, refresh_token,
				access_token_expires_at, refresh_token_expires_at, connected_at, updated_at
			) VALUES ($1, $2, $3, $4, $5, $6, $7)
			RETURNING id`,
			conn.RealmID, conn.AccessToken, conn.RefreshToken,
			conn.AccessTokenExpiresAt, conn.RefreshTokenExpiresAt,
			conn.ConnectedAt, conn.UpdatedAt,
		).Scan(&conn.ID)
	}

	conn.ID = existing.ID
	_, err = r.db.Exec(`
		UPDATE qbo_connections SET
			realm_id = $1,
			access_token = $2,
			refresh_token = $3,
			access_token_expires_at = $4,
			refresh_token_expires_at = $5,
			connected_at = $6,
			updated_at = $7
		WHERE id = $8`,
		conn.RealmID, conn.AccessToken, conn.RefreshToken,
		conn.AccessTokenExpiresAt, conn.RefreshTokenExpiresAt,
		conn.ConnectedAt, conn.UpdatedAt, conn.ID,
	)
	if err != nil {
		return fmt.Errorf("update qbo connection: %w", err)
	}
	return nil
}

func (r *SQLQBORepository) DeleteConnection() error {
	_, err := r.db.Exec(`DELETE FROM qbo_connections`)
	return err
}

func (r *SQLQBORepository) SaveOAuthState(state string, expiresAt time.Time) error {
	_, err := r.db.Exec(
		`INSERT INTO qbo_oauth_states (state, expires_at) VALUES ($1, $2)`,
		state, expiresAt,
	)
	if err != nil {
		return fmt.Errorf("save oauth state: %w", err)
	}
	return nil
}

func (r *SQLQBORepository) ConsumeOAuthState(state string) (bool, error) {
	var expiresAt time.Time
	err := r.db.QueryRow(
		`SELECT expires_at FROM qbo_oauth_states WHERE state = $1`, state,
	).Scan(&expiresAt)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("lookup oauth state: %w", err)
	}
	_, _ = r.db.Exec(`DELETE FROM qbo_oauth_states WHERE state = $1`, state)
	_, _ = r.db.Exec(`DELETE FROM qbo_oauth_states WHERE expires_at < CURRENT_TIMESTAMP`)
	if time.Now().After(expiresAt) {
		return false, nil
	}
	return true, nil
}

func (r *SQLQBORepository) ListMappings() ([]models.PropertyQBOLink, error) {
	rows, err := r.db.Query(`
		SELECT property_id, qbo_class_id, qbo_class_name, updated_at
		FROM property_qbo_links
		ORDER BY property_id ASC`)
	if err != nil {
		return nil, fmt.Errorf("list qbo mappings: %w", err)
	}
	defer rows.Close()

	var out []models.PropertyQBOLink
	for rows.Next() {
		var link models.PropertyQBOLink
		if err := rows.Scan(&link.PropertyID, &link.QBOClassID, &link.QBOClassName, &link.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, link)
	}
	return out, rows.Err()
}

func (r *SQLQBORepository) UpsertMapping(link *models.PropertyQBOLink) error {
	now := time.Now()
	link.UpdatedAt = now
	if r.isSQLite {
		_, err := r.db.Exec(`
			INSERT INTO property_qbo_links (property_id, qbo_class_id, qbo_class_name, updated_at)
			VALUES ($1, $2, $3, $4)
			ON CONFLICT(property_id) DO UPDATE SET
				qbo_class_id = excluded.qbo_class_id,
				qbo_class_name = excluded.qbo_class_name,
				updated_at = excluded.updated_at`,
			link.PropertyID, link.QBOClassID, link.QBOClassName, link.UpdatedAt,
		)
		return err
	}
	_, err := r.db.Exec(`
		INSERT INTO property_qbo_links (property_id, qbo_class_id, qbo_class_name, updated_at)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (property_id) DO UPDATE SET
			qbo_class_id = EXCLUDED.qbo_class_id,
			qbo_class_name = EXCLUDED.qbo_class_name,
			updated_at = EXCLUDED.updated_at`,
		link.PropertyID, link.QBOClassID, link.QBOClassName, link.UpdatedAt,
	)
	return err
}

func (r *SQLQBORepository) DeleteMapping(propertyID int) error {
	_, err := r.db.Exec(`DELETE FROM property_qbo_links WHERE property_id = $1`, propertyID)
	return err
}

func (r *SQLQBORepository) DeleteAllMappings() error {
	_, err := r.db.Exec(`DELETE FROM property_qbo_links`)
	return err
}

func (r *SQLQBORepository) UpsertCashFlow(cache *models.QBOCashFlowCache) error {
	if r.isSQLite {
		_, err := r.db.Exec(`
			INSERT INTO qbo_cash_flow_cache (property_id, period_start, period_end, net_amount, synced_at)
			VALUES ($1, $2, $3, $4, $5)
			ON CONFLICT(property_id) DO UPDATE SET
				period_start = excluded.period_start,
				period_end = excluded.period_end,
				net_amount = excluded.net_amount,
				synced_at = excluded.synced_at`,
			cache.PropertyID, cache.PeriodStart, cache.PeriodEnd, cache.NetAmount, cache.SyncedAt,
		)
		return err
	}
	_, err := r.db.Exec(`
		INSERT INTO qbo_cash_flow_cache (property_id, period_start, period_end, net_amount, synced_at)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (property_id) DO UPDATE SET
			period_start = EXCLUDED.period_start,
			period_end = EXCLUDED.period_end,
			net_amount = EXCLUDED.net_amount,
			synced_at = EXCLUDED.synced_at`,
		cache.PropertyID, cache.PeriodStart, cache.PeriodEnd, cache.NetAmount, cache.SyncedAt,
	)
	return err
}

func (r *SQLQBORepository) ListCashFlowForPeriod(periodStart, periodEnd time.Time) ([]models.QBOCashFlowCache, error) {
	rows, err := r.db.Query(`
		SELECT property_id, period_start, period_end, net_amount, synced_at
		FROM qbo_cash_flow_cache
		WHERE period_start = $1 AND period_end = $2`,
		periodStart, periodEnd,
	)
	if err != nil {
		return nil, fmt.Errorf("list cash flow cache: %w", err)
	}
	defer rows.Close()

	var out []models.QBOCashFlowCache
	for rows.Next() {
		var c models.QBOCashFlowCache
		if err := rows.Scan(&c.PropertyID, &c.PeriodStart, &c.PeriodEnd, &c.NetAmount, &c.SyncedAt); err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, rows.Err()
}

func (r *SQLQBORepository) LatestSyncTime() (*time.Time, error) {
	var t sql.NullTime
	err := r.db.QueryRow(`SELECT MAX(synced_at) FROM qbo_cash_flow_cache`).Scan(&t)
	if err != nil {
		return nil, err
	}
	if !t.Valid {
		return nil, nil
	}
	return &t.Time, nil
}
