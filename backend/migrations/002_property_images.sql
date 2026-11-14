-- Property photos (multiple per property; one may be marked main).
-- Applied automatically via createTables in internal/db for new installs;
-- this file documents the schema for manual Postgres/SQLite runs.

CREATE TABLE IF NOT EXISTS property_images (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	property_id INTEGER NOT NULL REFERENCES properties(id) ON DELETE CASCADE,
	file_name TEXT NOT NULL,
	content_type TEXT NOT NULL,
	size INTEGER NOT NULL,
	is_main INTEGER NOT NULL DEFAULT 0,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_property_images_property_id ON property_images(property_id);
