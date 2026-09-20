package models

import "time"

// QBOConnection stores OAuth tokens for the linked QuickBooks company.
type QBOConnection struct {
	ID                    int       `json:"id"`
	RealmID               string    `json:"realm_id"`
	AccessToken           string    `json:"-"`
	RefreshToken          string    `json:"-"`
	AccessTokenExpiresAt  time.Time `json:"access_token_expires_at"`
	RefreshTokenExpiresAt time.Time `json:"refresh_token_expires_at"`
	ConnectedAt           time.Time `json:"connected_at"`
	UpdatedAt             time.Time `json:"updated_at"`
}

// PropertyQBOLink maps a Brickwise property to a QuickBooks Class.
type PropertyQBOLink struct {
	PropertyID   int       `json:"property_id"`
	QBOClassID   string    `json:"qbo_class_id"`
	QBOClassName string    `json:"qbo_class_name"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// QBOCashFlowCache stores the last synced net cash flow for a property/period.
type QBOCashFlowCache struct {
	PropertyID  int       `json:"property_id"`
	PeriodStart time.Time `json:"period_start"`
	PeriodEnd   time.Time `json:"period_end"`
	NetAmount   float64   `json:"net_amount"`
	SyncedAt    time.Time `json:"synced_at"`
}

// QBOClass is a Class from QuickBooks Online.
type QBOClass struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// QBOStatus is returned by GET /integrations/qbo/status.
type QBOStatus struct {
	Configured    bool               `json:"configured"`
	Connected     bool               `json:"connected"`
	RealmID       string             `json:"realm_id,omitempty"`
	ConnectedAt   *time.Time         `json:"connected_at,omitempty"`
	LastSyncedAt  *time.Time         `json:"last_synced_at,omitempty"`
	PeriodStart   *time.Time         `json:"period_start,omitempty"`
	PeriodEnd     *time.Time         `json:"period_end,omitempty"`
	Mappings      []PropertyQBOLink  `json:"mappings"`
	ClassRequired bool               `json:"class_required"`
	Message       string             `json:"message,omitempty"`
}

// UpdateQBOMappingsRequest updates property→Class mappings.
type UpdateQBOMappingsRequest struct {
	Mappings []PropertyQBOMappingInput `json:"mappings"`
}

// PropertyQBOMappingInput is one mapping row from the settings UI.
type PropertyQBOMappingInput struct {
	PropertyID   int    `json:"property_id"`
	QBOClassID   string `json:"qbo_class_id"`
	QBOClassName string `json:"qbo_class_name"`
}

// QBOSyncResult summarizes a cash-flow sync.
type QBOSyncResult struct {
	SyncedAt    time.Time `json:"synced_at"`
	PeriodStart time.Time `json:"period_start"`
	PeriodEnd   time.Time `json:"period_end"`
	Updated     int       `json:"updated"`
}
