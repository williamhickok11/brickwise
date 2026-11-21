package db

import (
	"database/sql"
	"fmt"
	"net/url"
	"os"
	"strings"

	"brickwise/backend/internal/logger"

	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

// InitDB initializes the database connection.
//
// - If DATABASE_URL is set, connects to Postgres (including Supabase).
// - Otherwise falls back to local SQLite (DB_PATH or ./brickwise.db).
//
// Schema is created/updated automatically via createTables on startup.
// Do not put credentials in code — set DATABASE_URL in the environment.
func InitDB() (*sql.DB, error) {
	log := logger.Default()
	dbURL := strings.TrimSpace(os.Getenv("DATABASE_URL"))

	var (
		db  *sql.DB
		err error
	)

	if dbURL != "" {
		normalized, errNorm := normalizePostgresURL(dbURL)
		if errNorm != nil {
			return nil, fmt.Errorf("invalid DATABASE_URL: %w", errNorm)
		}
		db, err = sql.Open("postgres", normalized)
		if err != nil {
			return nil, fmt.Errorf("failed to open postgres connection: %w", err)
		}
		// Supabase pooler / free-tier friendly defaults
		db.SetMaxOpenConns(10)
		db.SetMaxIdleConns(5)
		log.Info("Using Postgres database (DATABASE_URL)")
	} else {
		dbPath := os.Getenv("DB_PATH")
		if dbPath == "" {
			dbPath = "./brickwise.db"
		}
		db, err = sql.Open("sqlite3", dbPath+"?_foreign_keys=1")
		if err != nil {
			return nil, fmt.Errorf("failed to open sqlite connection: %w", err)
		}
		db.SetMaxOpenConns(25)
		db.SetMaxIdleConns(5)
		log.Info("Using SQLite database", "path", dbPath)
	}

	if err := db.Ping(); err != nil {
		_ = db.Close()
		if dbURL != "" {
			return nil, fmt.Errorf("failed to ping postgres (check DATABASE_URL, password, and sslmode): %w", err)
		}
		return nil, fmt.Errorf("failed to ping sqlite: %w", err)
	}

	if err := createTables(db); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("failed to create tables: %w", err)
	}

	return db, nil
}

// UsingPostgres reports whether DATABASE_URL is set (Postgres / Supabase mode).
func UsingPostgres() bool {
	return strings.TrimSpace(os.Getenv("DATABASE_URL")) != ""
}

// normalizePostgresURL accepts postgres:// or postgresql:// URLs and ensures
// Supabase hosts have sslmode=require when unset. Never logs the password.
func normalizePostgresURL(raw string) (string, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return "", fmt.Errorf("empty connection string")
	}

	// lib/pq accepts both; normalize scheme for consistent parsing.
	if strings.HasPrefix(raw, "postgresql://") {
		raw = "postgres://" + strings.TrimPrefix(raw, "postgresql://")
	}

	u, err := url.Parse(raw)
	if err != nil {
		return "", err
	}
	if u.Scheme != "postgres" && u.Scheme != "postgresql" {
		return "", fmt.Errorf("unsupported scheme %q (expected postgres:// or postgresql://)", u.Scheme)
	}

	q := u.Query()
	host := strings.ToLower(u.Hostname())
	if q.Get("sslmode") == "" && (strings.Contains(host, "supabase.co") || strings.Contains(host, "supabase.com")) {
		q.Set("sslmode", "require")
	}
	u.RawQuery = q.Encode()
	return u.String(), nil
}

// createTables creates the database schema for SQLite or Postgres.
func createTables(db *sql.DB) error {
	isPostgres := UsingPostgres()
	var createTableSQL string

	if isPostgres {
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

	if isPostgres {
		alterSQL := `DO $$ 
		BEGIN 
			IF NOT EXISTS (SELECT 1 FROM information_schema.columns 
				WHERE table_name='properties' AND column_name='default_mileage') THEN
				ALTER TABLE properties ADD COLUMN default_mileage REAL DEFAULT 0;
			END IF;
		END $$;`
		db.Exec(alterSQL)
	} else {
		db.Exec(`ALTER TABLE properties ADD COLUMN default_mileage REAL DEFAULT 0`)
	}

	if isPostgres {
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

	if isPostgres {
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

	if isPostgres {
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

	if isPostgres {
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

	if isPostgres {
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

	if err := createQBOTables(db, isPostgres); err != nil {
		return err
	}
	if err := createDocumentsTable(db, isPostgres); err != nil {
		return err
	}
	if err := createPropertyImagesTable(db, isPostgres); err != nil {
		return err
	}

	return nil
}

func createPropertyImagesTable(db *sql.DB, isPostgres bool) error {
	var stmt string
	if isPostgres {
		stmt = `CREATE TABLE IF NOT EXISTS property_images (
			id SERIAL PRIMARY KEY,
			property_id INTEGER NOT NULL REFERENCES properties(id) ON DELETE CASCADE,
			file_name VARCHAR(255) NOT NULL,
			content_type VARCHAR(100) NOT NULL,
			size BIGINT NOT NULL,
			is_main BOOLEAN NOT NULL DEFAULT FALSE,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`
	} else {
		stmt = `CREATE TABLE IF NOT EXISTS property_images (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			property_id INTEGER NOT NULL REFERENCES properties(id) ON DELETE CASCADE,
			file_name TEXT NOT NULL,
			content_type TEXT NOT NULL,
			size INTEGER NOT NULL,
			is_main INTEGER NOT NULL DEFAULT 0,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`
	}
	if _, err := db.Exec(stmt); err != nil {
		return fmt.Errorf("failed to create property_images table: %w", err)
	}
	if _, err := db.Exec(`CREATE INDEX IF NOT EXISTS idx_property_images_property_id ON property_images(property_id)`); err != nil {
		return fmt.Errorf("failed to create property_images index: %w", err)
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
