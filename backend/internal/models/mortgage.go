package models

import "time"

type Mortgage struct {
	ID             int        `json:"id"`
	PropertyID     int        `json:"property_id"`
	OriginalLoan   float64    `json:"original_loan"`
	CurrentBalance float64    `json:"current_balance"`
	InterestRate   float64    `json:"interest_rate"`
	MonthlyPayment float64    `json:"monthly_payment"`
	MonthlyPI      float64    `json:"monthly_pi"`
	BalanceAsOf    *Date      `json:"balance_as_of"`
	Lender         string     `json:"lender"`
	Notes          string     `json:"notes"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

type UpsertMortgageRequest struct {
	OriginalLoan   float64 `json:"original_loan"`
	CurrentBalance float64 `json:"current_balance"`
	InterestRate   float64 `json:"interest_rate"`
	MonthlyPayment float64 `json:"monthly_payment"`
	MonthlyPI      float64 `json:"monthly_pi"`
	BalanceAsOf    *string `json:"balance_as_of"`
	Lender         string  `json:"lender"`
	Notes          string  `json:"notes"`
}
