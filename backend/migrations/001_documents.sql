-- Lease PDF documents attached to residences (one per residence).
-- Applied automatically by db.createTables on server start.

-- SQLite
CREATE TABLE IF NOT EXISTS documents (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	residence_id INTEGER NOT NULL UNIQUE REFERENCES residences(id) ON DELETE CASCADE,
	file_name TEXT NOT NULL,
	content_type TEXT NOT NULL,
	size INTEGER NOT NULL,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Postgres equivalent:
-- CREATE TABLE IF NOT EXISTS documents (
-- 	id SERIAL PRIMARY KEY,
-- 	residence_id INTEGER NOT NULL UNIQUE REFERENCES residences(id) ON DELETE CASCADE,
-- 	file_name VARCHAR(255) NOT NULL,
-- 	content_type VARCHAR(100) NOT NULL,
-- 	size BIGINT NOT NULL,
-- 	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
-- );
