package com.brickwise.config;

import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.assertTrue;

class DatabaseConfigTest {
    @Test
    void normalizePostgresUrl_addsSslForSupabase() {
        String raw = "postgresql://postgres:secret@db.abcdefgh.supabase.co:5432/postgres";
        String normalized = DatabaseConfig.normalizePostgresUrl(raw);
        assertTrue(normalized.contains("sslmode=require"));
        assertTrue(normalized.startsWith("postgres://"));
    }

    @Test
    void toJdbcPostgresUrl_includesCredentials() {
        String jdbc = DatabaseConfig.toJdbcPostgresUrl(
                "postgres://postgres:s3cret@db.abcdefgh.supabase.co:5432/postgres?sslmode=require");
        assertTrue(jdbc.startsWith("jdbc:postgresql://db.abcdefgh.supabase.co:5432/postgres?"));
        assertTrue(jdbc.contains("user=postgres"));
        assertTrue(jdbc.contains("password=s3cret"));
        assertTrue(jdbc.contains("sslmode=require"));
    }
}
