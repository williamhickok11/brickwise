#!/bin/bash
# Script to create properties for REPS V2
# Usage: ./create-properties.sh

DB_PATH="${DB_PATH:-./brickwise.db}"
DATABASE_URL="${DATABASE_URL:-}"

if [ -n "$DATABASE_URL" ]; then
    echo "Creating properties in Postgres database..."
    psql "$DATABASE_URL" -f "$(dirname "$0")/create-properties.sql"
else
    echo "Creating properties in SQLite database: $DB_PATH"
    sqlite3 "$DB_PATH" < "$(dirname "$0")/create-properties.sql"
fi

echo "Properties created successfully!"
