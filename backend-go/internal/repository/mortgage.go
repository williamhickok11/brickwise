package repository

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

	"brickwise/backend/internal/models"
)

type MortgageRepository interface {
	Upsert(mortgage *models.Mortgage) error
	GetByPropertyID(propertyID int) (*models.Mortgage, error)
	GetAll() ([]*models.Mortgage, error)
	DeleteByPropertyID(propertyID int) error
}

type SQLMortgageRepository struct {
	db       *sql.DB
	isSQLite bool
}

func NewMortgageRepository(db *sql.DB) MortgageRepository {
	isSQLite := os.Getenv("DATABASE_URL") == ""
	return &SQLMortgageRepository{db: db, isSQLite: isSQLite}
}

func (r *SQLMortgageRepository) Upsert(mortgage *models.Mortgage) error {
	var balanceAsOf any
	if mortgage.BalanceAsOf != nil {
		balanceAsOf = mortgage.BalanceAsOf.Time
	}

	existing, err := r.GetByPropertyID(mortgage.PropertyID)
	if err != nil && !isNotFound(err) {
		return err
	}

	if existing == nil {
		if r.isSQLite {
			query := `INSERT INTO mortgages (property_id, original_loan, current_balance, interest_rate,
				monthly_payment, monthly_pi, balance_as_of, lender, notes)
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
			result, err := r.db.Exec(query,
				mortgage.PropertyID, mortgage.OriginalLoan, mortgage.CurrentBalance,
				mortgage.InterestRate, mortgage.MonthlyPayment, mortgage.MonthlyPI,
				balanceAsOf, mortgage.Lender, mortgage.Notes,
			)
			if err != nil {
				return fmt.Errorf("failed to create mortgage: %w", err)
			}
			id, err := result.LastInsertId()
			if err != nil {
				return fmt.Errorf("failed to get last insert id: %w", err)
			}
			mortgage.ID = int(id)
			return r.db.QueryRow(
				`SELECT created_at, updated_at FROM mortgages WHERE id = $1`, mortgage.ID,
			).Scan(&mortgage.CreatedAt, &mortgage.UpdatedAt)
		}

		query := `INSERT INTO mortgages (property_id, original_loan, current_balance, interest_rate,
			monthly_payment, monthly_pi, balance_as_of, lender, notes)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
			RETURNING id, created_at, updated_at`
		return r.db.QueryRow(query,
			mortgage.PropertyID, mortgage.OriginalLoan, mortgage.CurrentBalance,
			mortgage.InterestRate, mortgage.MonthlyPayment, mortgage.MonthlyPI,
			balanceAsOf, mortgage.Lender, mortgage.Notes,
		).Scan(&mortgage.ID, &mortgage.CreatedAt, &mortgage.UpdatedAt)
	}

	mortgage.ID = existing.ID
	if r.isSQLite {
		query := `UPDATE mortgages SET original_loan = $1, current_balance = $2, interest_rate = $3,
			monthly_payment = $4, monthly_pi = $5, balance_as_of = $6, lender = $7, notes = $8,
			updated_at = CURRENT_TIMESTAMP WHERE property_id = $9`
		result, err := r.db.Exec(query,
			mortgage.OriginalLoan, mortgage.CurrentBalance, mortgage.InterestRate,
			mortgage.MonthlyPayment, mortgage.MonthlyPI, balanceAsOf,
			mortgage.Lender, mortgage.Notes, mortgage.PropertyID,
		)
		if err != nil {
			return fmt.Errorf("failed to update mortgage: %w", err)
		}
		rowsAffected, _ := result.RowsAffected()
		if rowsAffected == 0 {
			return fmt.Errorf("mortgage for property %d not found", mortgage.PropertyID)
		}
		return r.db.QueryRow(`SELECT updated_at FROM mortgages WHERE id = $1`, mortgage.ID).
			Scan(&mortgage.UpdatedAt)
	}

	query := `UPDATE mortgages SET original_loan = $1, current_balance = $2, interest_rate = $3,
		monthly_payment = $4, monthly_pi = $5, balance_as_of = $6, lender = $7, notes = $8,
		updated_at = CURRENT_TIMESTAMP WHERE property_id = $9
		RETURNING updated_at`
	return r.db.QueryRow(query,
		mortgage.OriginalLoan, mortgage.CurrentBalance, mortgage.InterestRate,
		mortgage.MonthlyPayment, mortgage.MonthlyPI, balanceAsOf,
		mortgage.Lender, mortgage.Notes, mortgage.PropertyID,
	).Scan(&mortgage.UpdatedAt)
}

func (r *SQLMortgageRepository) GetByPropertyID(propertyID int) (*models.Mortgage, error) {
	mortgage := &models.Mortgage{}
	query := `SELECT id, property_id, original_loan, current_balance, interest_rate,
		monthly_payment, monthly_pi, balance_as_of, lender, notes, created_at, updated_at
		FROM mortgages WHERE property_id = $1`

	var balanceAsOf sql.NullTime
	err := r.db.QueryRow(query, propertyID).Scan(
		&mortgage.ID, &mortgage.PropertyID, &mortgage.OriginalLoan, &mortgage.CurrentBalance,
		&mortgage.InterestRate, &mortgage.MonthlyPayment, &mortgage.MonthlyPI,
		&balanceAsOf, &mortgage.Lender, &mortgage.Notes, &mortgage.CreatedAt, &mortgage.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("mortgage for property %d not found", propertyID)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get mortgage: %w", err)
	}
	if balanceAsOf.Valid {
		mortgage.BalanceAsOf = &models.Date{Time: balanceAsOf.Time}
	}
	return mortgage, nil
}

func (r *SQLMortgageRepository) GetAll() ([]*models.Mortgage, error) {
	query := `SELECT id, property_id, original_loan, current_balance, interest_rate,
		monthly_payment, monthly_pi, balance_as_of, lender, notes, created_at, updated_at
		FROM mortgages ORDER BY property_id`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query mortgages: %w", err)
	}
	defer rows.Close()

	var mortgages []*models.Mortgage
	for rows.Next() {
		mortgage := &models.Mortgage{}
		var balanceAsOf sql.NullTime
		if err := rows.Scan(
			&mortgage.ID, &mortgage.PropertyID, &mortgage.OriginalLoan, &mortgage.CurrentBalance,
			&mortgage.InterestRate, &mortgage.MonthlyPayment, &mortgage.MonthlyPI,
			&balanceAsOf, &mortgage.Lender, &mortgage.Notes, &mortgage.CreatedAt, &mortgage.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan mortgage: %w", err)
		}
		if balanceAsOf.Valid {
			mortgage.BalanceAsOf = &models.Date{Time: balanceAsOf.Time}
		}
		mortgages = append(mortgages, mortgage)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating mortgages: %w", err)
	}
	return mortgages, nil
}

func (r *SQLMortgageRepository) DeleteByPropertyID(propertyID int) error {
	query := `DELETE FROM mortgages WHERE property_id = $1`
	result, err := r.db.Exec(query, propertyID)
	if err != nil {
		return fmt.Errorf("failed to delete mortgage: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("mortgage for property %d not found", propertyID)
	}
	return nil
}

func isNotFound(err error) bool {
	return err != nil && strings.Contains(err.Error(), "not found")
}
