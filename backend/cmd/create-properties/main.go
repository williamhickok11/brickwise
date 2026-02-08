package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Connect to database
	var db *sql.DB
	var err error

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL != "" {
		db, err = sql.Open("postgres", dbURL)
	} else {
		dbPath := os.Getenv("DB_PATH")
		if dbPath == "" {
			dbPath = "./brickwise.db"
		}
		db, err = sql.Open("sqlite3", dbPath+"?_foreign_keys=1")
	}

	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	// Properties to create
	properties := []struct {
		name          string
		address       string
		propertyType  string
		defaultMileage float64
	}{
		{
			name:          "3522 Wood Bridge Dr",
			address:       "3522 Wood Bridge Dr",
			propertyType:  "residential",
			defaultMileage: 0,
		},
		{
			name:          "3500 Saindon St",
			address:       "3500 Saindon St",
			propertyType:  "residential",
			defaultMileage: 0,
		},
	}

	// Insert properties
	for _, prop := range properties {
		// Check if property already exists
		var count int
		err := db.QueryRow("SELECT COUNT(*) FROM properties WHERE address = $1", prop.address).Scan(&count)
		if err != nil {
			log.Fatalf("Failed to check existing property: %v", err)
		}

		if count > 0 {
			fmt.Printf("Property '%s' already exists, skipping...\n", prop.address)
			continue
		}

		// Insert property
		var id int64
		if dbURL != "" {
			// Postgres
			err = db.QueryRow(
				"INSERT INTO properties (name, address, property_type, default_mileage) VALUES ($1, $2, $3, $4) RETURNING id",
				prop.name, prop.address, prop.propertyType, prop.defaultMileage,
			).Scan(&id)
		} else {
			// SQLite
			result, err := db.Exec(
				"INSERT INTO properties (name, address, property_type, default_mileage) VALUES ($1, $2, $3, $4)",
				prop.name, prop.address, prop.propertyType, prop.defaultMileage,
			)
			if err != nil {
				log.Fatalf("Failed to insert property: %v", err)
			}
			id, err = result.LastInsertId()
		}

		if err != nil {
			log.Fatalf("Failed to insert property '%s': %v", prop.address, err)
		}

		fmt.Printf("✓ Created property: %s (ID: %d)\n", prop.address, id)
	}

	fmt.Println("\nAll properties created successfully!")
}
