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

	return nil
}
