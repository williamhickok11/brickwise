-- Property photos (multiple per property; one may be marked main).
-- Schema is applied automatically by db.createTables on server start.
-- This file documents the DDL for manual runs.

-- SQLite
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

-- Postgres / Supabase
-- CREATE TABLE IF NOT EXISTS property_images (
-- 	id SERIAL PRIMARY KEY,
-- 	property_id INTEGER NOT NULL REFERENCES properties(id) ON DELETE CASCADE,
-- 	file_name VARCHAR(255) NOT NULL,
-- 	content_type VARCHAR(100) NOT NULL,
-- 	size BIGINT NOT NULL,
-- 	is_main BOOLEAN NOT NULL DEFAULT FALSE,
-- 	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
-- );
-- CREATE INDEX IF NOT EXISTS idx_property_images_property_id ON property_images(property_id);
