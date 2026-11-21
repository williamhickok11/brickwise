package com.brickwise.config;

import java.sql.SQLException;
import java.sql.Statement;

/** Creates tables on startup (parity with Go {@code createTables}). */
final class SchemaBootstrap {
    private SchemaBootstrap() {}

    static void createTables(Statement st, boolean postgres) throws SQLException {
        if (postgres) {
            st.execute("""
                CREATE TABLE IF NOT EXISTS properties (
                  id SERIAL PRIMARY KEY,
                  name VARCHAR(255) NOT NULL,
                  address VARCHAR(500) NOT NULL,
                  property_type VARCHAR(50) NOT NULL,
                  default_mileage REAL DEFAULT 0,
                  current_value REAL DEFAULT 0,
                  value_updated_at TIMESTAMP NULL,
                  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
                )
                """);
            st.execute("""
                CREATE TABLE IF NOT EXISTS residences (
                  id SERIAL PRIMARY KEY,
                  property_id INTEGER NOT NULL REFERENCES properties(id) ON DELETE CASCADE,
                  name VARCHAR(255) NOT NULL,
                  phone TEXT,
                  email TEXT,
                  start_date DATE NOT NULL,
                  end_date DATE NULL,
                  is_active BOOLEAN NOT NULL DEFAULT TRUE,
                  monthly_rent REAL DEFAULT 0,
                  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
                )
                """);
            st.execute("""
                CREATE TABLE IF NOT EXISTS mortgages (
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
                )
                """);
            st.execute("""
                CREATE TABLE IF NOT EXISTS time_entries (
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
                )
                """);
            st.execute("""
                CREATE TABLE IF NOT EXISTS documents (
                  id SERIAL PRIMARY KEY,
                  residence_id INTEGER NOT NULL UNIQUE REFERENCES residences(id) ON DELETE CASCADE,
                  file_name VARCHAR(255) NOT NULL,
                  content_type VARCHAR(100) NOT NULL,
                  size BIGINT NOT NULL,
                  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
                )
                """);
            st.execute("""
                CREATE TABLE IF NOT EXISTS property_images (
                  id SERIAL PRIMARY KEY,
                  property_id INTEGER NOT NULL REFERENCES properties(id) ON DELETE CASCADE,
                  file_name VARCHAR(255) NOT NULL,
                  content_type VARCHAR(100) NOT NULL,
                  size BIGINT NOT NULL,
                  is_main BOOLEAN NOT NULL DEFAULT FALSE,
                  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
                )
                """);
            st.execute("""
                CREATE TABLE IF NOT EXISTS qbo_connections (
                  id SERIAL PRIMARY KEY,
                  realm_id TEXT NOT NULL,
                  access_token TEXT NOT NULL,
                  refresh_token TEXT NOT NULL,
                  access_token_expires_at TIMESTAMP NOT NULL,
                  refresh_token_expires_at TIMESTAMP NOT NULL,
                  connected_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
                )
                """);
            st.execute("""
                CREATE TABLE IF NOT EXISTS qbo_oauth_states (
                  state TEXT PRIMARY KEY,
                  expires_at TIMESTAMP NOT NULL
                )
                """);
            st.execute("""
                CREATE TABLE IF NOT EXISTS property_qbo_links (
                  property_id INTEGER PRIMARY KEY REFERENCES properties(id) ON DELETE CASCADE,
                  qbo_class_id TEXT NOT NULL,
                  qbo_class_name TEXT NOT NULL DEFAULT '',
                  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
                )
                """);
            st.execute("""
                CREATE TABLE IF NOT EXISTS qbo_cash_flow_cache (
                  property_id INTEGER PRIMARY KEY REFERENCES properties(id) ON DELETE CASCADE,
                  period_start DATE NOT NULL,
                  period_end DATE NOT NULL,
                  net_amount REAL NOT NULL DEFAULT 0,
                  synced_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
                )
                """);
        } else {
            st.execute("""
                CREATE TABLE IF NOT EXISTS properties (
                  id INTEGER PRIMARY KEY AUTOINCREMENT,
                  name TEXT NOT NULL,
                  address TEXT NOT NULL,
                  property_type TEXT NOT NULL,
                  default_mileage REAL DEFAULT 0,
                  current_value REAL DEFAULT 0,
                  value_updated_at DATETIME NULL,
                  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
                  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
                )
                """);
            st.execute("""
                CREATE TABLE IF NOT EXISTS residences (
                  id INTEGER PRIMARY KEY AUTOINCREMENT,
                  property_id INTEGER NOT NULL REFERENCES properties(id) ON DELETE CASCADE,
                  name TEXT NOT NULL,
                  phone TEXT,
                  email TEXT,
                  start_date DATE NOT NULL,
                  end_date DATE NULL,
                  is_active INTEGER NOT NULL DEFAULT 1,
                  monthly_rent REAL DEFAULT 0,
                  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
                  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
                )
                """);
            st.execute("""
                CREATE TABLE IF NOT EXISTS mortgages (
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
                )
                """);
            st.execute("""
                CREATE TABLE IF NOT EXISTS time_entries (
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
                )
                """);
            st.execute("""
                CREATE TABLE IF NOT EXISTS documents (
                  id INTEGER PRIMARY KEY AUTOINCREMENT,
                  residence_id INTEGER NOT NULL UNIQUE REFERENCES residences(id) ON DELETE CASCADE,
                  file_name TEXT NOT NULL,
                  content_type TEXT NOT NULL,
                  size INTEGER NOT NULL,
                  created_at DATETIME DEFAULT CURRENT_TIMESTAMP
                )
                """);
            st.execute("""
                CREATE TABLE IF NOT EXISTS property_images (
                  id INTEGER PRIMARY KEY AUTOINCREMENT,
                  property_id INTEGER NOT NULL REFERENCES properties(id) ON DELETE CASCADE,
                  file_name TEXT NOT NULL,
                  content_type TEXT NOT NULL,
                  size INTEGER NOT NULL,
                  is_main INTEGER NOT NULL DEFAULT 0,
                  created_at DATETIME DEFAULT CURRENT_TIMESTAMP
                )
                """);
            st.execute("""
                CREATE TABLE IF NOT EXISTS qbo_connections (
                  id INTEGER PRIMARY KEY AUTOINCREMENT,
                  realm_id TEXT NOT NULL,
                  access_token TEXT NOT NULL,
                  refresh_token TEXT NOT NULL,
                  access_token_expires_at DATETIME NOT NULL,
                  refresh_token_expires_at DATETIME NOT NULL,
                  connected_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
                  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
                )
                """);
            st.execute("""
                CREATE TABLE IF NOT EXISTS qbo_oauth_states (
                  state TEXT PRIMARY KEY,
                  expires_at DATETIME NOT NULL
                )
                """);
            st.execute("""
                CREATE TABLE IF NOT EXISTS property_qbo_links (
                  property_id INTEGER PRIMARY KEY REFERENCES properties(id) ON DELETE CASCADE,
                  qbo_class_id TEXT NOT NULL,
                  qbo_class_name TEXT NOT NULL DEFAULT '',
                  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
                )
                """);
            st.execute("""
                CREATE TABLE IF NOT EXISTS qbo_cash_flow_cache (
                  property_id INTEGER PRIMARY KEY REFERENCES properties(id) ON DELETE CASCADE,
                  period_start DATE NOT NULL,
                  period_end DATE NOT NULL,
                  net_amount REAL NOT NULL DEFAULT 0,
                  synced_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
                )
                """);
        }
        st.execute("CREATE INDEX IF NOT EXISTS idx_property_images_property_id ON property_images(property_id)");
    }
}
