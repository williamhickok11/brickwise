package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

// InitDB initializes the database connection
// Supports both SQLite (default) and Postgres (via DATABASE_URL env var)
// Tradeoff: Simple env-based detection. For production, consider explicit config.
func InitDB() (*sql.DB, error) {
	var db *sql.DB
	var err error

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL != "" {
		// Postgres connection
		db, err = sql.Open("postgres", dbURL)
		if err != nil {
			return nil, fmt.Errorf("failed to open postgres connection: %w", err)
		}
	} else {
		// SQLite connection (default for development)
		dbPath := os.Getenv("DB_PATH")
		if dbPath == "" {
			dbPath = "./brickwise.db"
		}
		db, err = sql.Open("sqlite3", dbPath+"?_foreign_keys=1")
		if err != nil {
			return nil, fmt.Errorf("failed to open sqlite connection: %w", err)
		}
	}

	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Create tables
	if err := createTables(db); err != nil {
		return nil, fmt.Errorf("failed to create tables: %w", err)
	}

	return db, nil
}

// createTables creates the database schema
// Uses SQL that works for both SQLite and Postgres
// Tradeoff: Detects database type via connection string/env var
func createTables(db *sql.DB) error {
	dbURL := os.Getenv("DATABASE_URL")
	var createTableSQL string

	if dbURL != "" {
		// Postgres
		createTableSQL = `CREATE TABLE IF NOT EXISTS properties (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			address VARCHAR(500) NOT NULL,
			property_type VARCHAR(50) NOT NULL,
			default_mileage REAL DEFAULT 0,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`
	} else {
		// SQLite
		createTableSQL = `CREATE TABLE IF NOT EXISTS properties (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			address TEXT NOT NULL,
			property_type TEXT NOT NULL,
			default_mileage REAL DEFAULT 0,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`
	}

	if _, err := db.Exec(createTableSQL); err != nil {
		return fmt.Errorf("failed to create properties table: %w", err)
	}

	// Add default_mileage column if it doesn't exist (for existing databases)
	if dbURL != "" {
		// Postgres: check if column exists and add if not
		alterSQL := `DO $$ 
		BEGIN 
			IF NOT EXISTS (SELECT 1 FROM information_schema.columns 
				WHERE table_name='properties' AND column_name='default_mileage') THEN
				ALTER TABLE properties ADD COLUMN default_mileage REAL DEFAULT 0;
			END IF;
		END $$;`
		db.Exec(alterSQL)
	} else {
		// SQLite: add column if it doesn't exist
		alterSQL := `ALTER TABLE properties ADD COLUMN default_mileage REAL DEFAULT 0`
		db.Exec(alterSQL) // Ignore error if column already exists
	}

	// Add financial columns to properties
	if dbURL != "" {
		alterFinancialSQL := `DO $$
		BEGIN
			IF NOT EXISTS (SELECT 1 FROM information_schema.columns
				WHERE table_name='properties' AND column_name='current_value') THEN
				ALTER TABLE properties ADD COLUMN current_value REAL DEFAULT 0;
			END IF;
			IF NOT EXISTS (SELECT 1 FROM information_schema.columns
				WHERE table_name='properties' AND column_name='value_updated_at') THEN
				ALTER TABLE properties ADD COLUMN value_updated_at TIMESTAMP NULL;
			END IF;
		END $$;`
		db.Exec(alterFinancialSQL)
	} else {
		db.Exec(`ALTER TABLE properties ADD COLUMN current_value REAL DEFAULT 0`)
		db.Exec(`ALTER TABLE properties ADD COLUMN value_updated_at DATETIME NULL`)
	}

	// Create residences table (tenant/lease history)
	if dbURL != "" {
		// Postgres
		createResidencesSQL := `CREATE TABLE IF NOT EXISTS residences (
			id SERIAL PRIMARY KEY,
			property_id INTEGER NOT NULL REFERENCES properties(id) ON DELETE CASCADE,
			name VARCHAR(255) NOT NULL,
			phone TEXT,
			email TEXT,
			start_date DATE NOT NULL,
			end_date DATE NULL,
			is_active BOOLEAN NOT NULL DEFAULT TRUE,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`
		if _, err := db.Exec(createResidencesSQL); err != nil {
			return fmt.Errorf("failed to create residences table: %w", err)
		}
	} else {
		// SQLite
		createResidencesSQL := `CREATE TABLE IF NOT EXISTS residences (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			property_id INTEGER NOT NULL REFERENCES properties(id) ON DELETE CASCADE,
			name TEXT NOT NULL,
			phone TEXT,
			email TEXT,
			start_date DATE NOT NULL,
			end_date DATE NULL,
			is_active INTEGER NOT NULL DEFAULT 1,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`
		if _, err := db.Exec(createResidencesSQL); err != nil {
			return fmt.Errorf("failed to create residences table: %w", err)
		}
	}

	// Add monthly_rent to residences
	if dbURL != "" {
		alterRentSQL := `DO $$
		BEGIN
			IF NOT EXISTS (SELECT 1 FROM information_schema.columns
				WHERE table_name='residences' AND column_name='monthly_rent') THEN
				ALTER TABLE residences ADD COLUMN monthly_rent REAL DEFAULT 0;
			END IF;
		END $$;`
		db.Exec(alterRentSQL)
	} else {
		db.Exec(`ALTER TABLE residences ADD COLUMN monthly_rent REAL DEFAULT 0`)
	}

	// Create mortgages table (1:1 optional with property)
	if dbURL != "" {
		createMortgagesSQL := `CREATE TABLE IF NOT EXISTS mortgages (
			id SERIAL PRIMARY KEY,
			property_id INTEGER NOT NULL UNIQUE REFERENCES properties(id) ON DELETE CASCADE,
			original_loan REAL DEFAULT 0,
			current_balance REAL DEFAULT 0,
			interest_rate REAL DEFAULT 0,
			monthly_payment REAL DEFAULT 0,
			monthly_pi REAL DEFAULT 0,
			balance_as_of DATE NULL,
			lender TEXT,
			notes TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`
		if _, err := db.Exec(createMortgagesSQL); err != nil {
			return fmt.Errorf("failed to create mortgages table: %w", err)
		}
	} else {
		createMortgagesSQL := `CREATE TABLE IF NOT EXISTS mortgages (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			property_id INTEGER NOT NULL UNIQUE REFERENCES properties(id) ON DELETE CASCADE,
			original_loan REAL DEFAULT 0,
			current_balance REAL DEFAULT 0,
			interest_rate REAL DEFAULT 0,
			monthly_payment REAL DEFAULT 0,
			monthly_pi REAL DEFAULT 0,
			balance_as_of DATE NULL,
			lender TEXT,
			notes TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`
		if _, err := db.Exec(createMortgagesSQL); err != nil {
			return fmt.Errorf("failed to create mortgages table: %w", err)
		}
	}

	// Create time_entries table
	if dbURL != "" {
		// Postgres
		createTimeEntriesSQL := `CREATE TABLE IF NOT EXISTS time_entries (
			id SERIAL PRIMARY KEY,
			property_id INTEGER REFERENCES properties(id) ON DELETE SET NULL,
			date DATE NOT NULL,
			category VARCHAR(100) NOT NULL,
			description TEXT NOT NULL,
			hours REAL NOT NULL,
			notes TEXT,
			mileage REAL DEFAULT 0,
			full_drive BOOLEAN DEFAULT FALSE,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`
		if _, err := db.Exec(createTimeEntriesSQL); err != nil {
			return fmt.Errorf("failed to create time_entries table: %w", err)
		}
	} else {
		// SQLite
		createTimeEntriesSQL := `CREATE TABLE IF NOT EXISTS time_entries (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			property_id INTEGER REFERENCES properties(id) ON DELETE SET NULL,
			date DATE NOT NULL,
			category TEXT NOT NULL,
			description TEXT NOT NULL,
			hours REAL NOT NULL,
			notes TEXT,
			mileage REAL DEFAULT 0,
			full_drive INTEGER DEFAULT 0,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`
		if _, err := db.Exec(createTimeEntriesSQL); err != nil {
			return fmt.Errorf("failed to create time_entries table: %w", err)
		}
	}

	if err := createQBOTables(db, dbURL != ""); err != nil {
		return err
	}

	if err := createDocumentsTable(db, dbURL != ""); err != nil {
		return err
	}

	return nil
}

func createDocumentsTable(db *sql.DB, isPostgres bool) error {
	var stmt string
	if isPostgres {
		stmt = `CREATE TABLE IF NOT EXISTS documents (
			id SERIAL PRIMARY KEY,
			residence_id INTEGER NOT NULL UNIQUE REFERENCES residences(id) ON DELETE CASCADE,
			file_name VARCHAR(255) NOT NULL,
			content_type VARCHAR(100) NOT NULL,
			size BIGINT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`
	} else {
		stmt = `CREATE TABLE IF NOT EXISTS documents (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			residence_id INTEGER NOT NULL UNIQUE REFERENCES residences(id) ON DELETE CASCADE,
			file_name TEXT NOT NULL,
			content_type TEXT NOT NULL,
			size INTEGER NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`
	}
	if _, err := db.Exec(stmt); err != nil {
		return fmt.Errorf("failed to create documents table: %w", err)
	}
	return nil
}

func createQBOTables(db *sql.DB, isPostgres bool) error {
	if isPostgres {
		stmts := []string{
			`CREATE TABLE IF NOT EXISTS qbo_connections (
				id SERIAL PRIMARY KEY,
				realm_id TEXT NOT NULL,
				access_token TEXT NOT NULL,
				refresh_token TEXT NOT NULL,
				access_token_expires_at TIMESTAMP NOT NULL,
				refresh_token_expires_at TIMESTAMP NOT NULL,
				connected_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
				updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
			)`,
			`CREATE TABLE IF NOT EXISTS qbo_oauth_states (
				state TEXT PRIMARY KEY,
				expires_at TIMESTAMP NOT NULL
			)`,
			`CREATE TABLE IF NOT EXISTS property_qbo_links (
				property_id INTEGER PRIMARY KEY REFERENCES properties(id) ON DELETE CASCADE,
				qbo_class_id TEXT NOT NULL,
				qbo_class_name TEXT NOT NULL DEFAULT '',
				updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
			)`,
			`CREATE TABLE IF NOT EXISTS qbo_cash_flow_cache (
				property_id INTEGER PRIMARY KEY REFERENCES properties(id) ON DELETE CASCADE,
				period_start DATE NOT NULL,
				period_end DATE NOT NULL,
				net_amount REAL NOT NULL DEFAULT 0,
				synced_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
			)`,
		}
		for _, stmt := range stmts {
			if _, err := db.Exec(stmt); err != nil {
				return fmt.Errorf("failed to create qbo tables: %w", err)
			}
		}
		return nil
	}

	stmts := []string{
		`CREATE TABLE IF NOT EXISTS qbo_connections (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			realm_id TEXT NOT NULL,
			access_token TEXT NOT NULL,
			refresh_token TEXT NOT NULL,
			access_token_expires_at DATETIME NOT NULL,
			refresh_token_expires_at DATETIME NOT NULL,
			connected_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS qbo_oauth_states (
			state TEXT PRIMARY KEY,
			expires_at DATETIME NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS property_qbo_links (
			property_id INTEGER PRIMARY KEY REFERENCES properties(id) ON DELETE CASCADE,
			qbo_class_id TEXT NOT NULL,
			qbo_class_name TEXT NOT NULL DEFAULT '',
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS qbo_cash_flow_cache (
			property_id INTEGER PRIMARY KEY REFERENCES properties(id) ON DELETE CASCADE,
			period_start DATE NOT NULL,
			period_end DATE NOT NULL,
			net_amount REAL NOT NULL DEFAULT 0,
			synced_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		)`,
	}
	for _, stmt := range stmts {
		if _, err := db.Exec(stmt); err != nil {
			return fmt.Errorf("failed to create qbo tables: %w", err)
		}
	}
	return nil
}
