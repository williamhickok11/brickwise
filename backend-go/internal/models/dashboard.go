package models

import "time"

type DashboardSummary struct {
	AsOf                 time.Time           `json:"as_of"`
	PortfolioValue       float64             `json:"portfolio_value"`
	TotalMortgageBalance float64             `json:"total_mortgage_balance"`
	TotalEquity          float64             `json:"total_equity"`
	MonthlyRevenue       float64             `json:"monthly_revenue"`
	MonthlyMortgagePI    float64             `json:"monthly_mortgage_pi"`
	EstimatedCashFlow    float64             `json:"estimated_cash_flow"`
	ActualCashFlow       *float64            `json:"actual_cash_flow"`
	QBOConnected         bool                `json:"qbo_connected"`
	QBOSyncedAt          *time.Time          `json:"qbo_synced_at,omitempty"`
	CashFlowPeriodStart  *time.Time          `json:"cash_flow_period_start,omitempty"`
	CashFlowPeriodEnd    *time.Time          `json:"cash_flow_period_end,omitempty"`
	Properties           []PropertyBreakdown `json:"properties"`
}

type PropertyBreakdown struct {
	PropertyID      int      `json:"property_id"`
	Name            string   `json:"name"`
	CurrentValue    float64  `json:"current_value"`
	MortgageBalance float64  `json:"mortgage_balance"`
	Equity          float64  `json:"equity"`
	MonthlyRent     float64  `json:"monthly_rent"`
	MonthlyPI       float64  `json:"monthly_pi"`
	NetCashFlow     float64  `json:"net_cash_flow"`
	ActualCashFlow  *float64 `json:"actual_cash_flow"`
	CashFlowSource  string   `json:"cash_flow_source"` // "quickbooks" | "estimate"
}
