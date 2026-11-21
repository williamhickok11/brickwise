package db

import (
	"strings"
	"testing"
)

func TestNormalizePostgresURL_SupabaseAddsSSL(t *testing.T) {
	in := "postgresql://postgres:secret@db.abcdefghijklmnop.supabase.co:5432/postgres"
	out, err := normalizePostgresURL(in)
	if err != nil {
		t.Fatalf("normalize: %v", err)
	}
	if !strings.HasPrefix(out, "postgres://") {
		t.Fatalf("expected postgres:// scheme, got %q", redact(out))
	}
	if !strings.Contains(out, "sslmode=require") {
		t.Fatalf("expected sslmode=require, got %q", redact(out))
	}
}

func TestNormalizePostgresURL_PreservesExistingSSL(t *testing.T) {
	in := "postgres://user:pass@localhost:5432/brickwise?sslmode=disable"
	out, err := normalizePostgresURL(in)
	if err != nil {
		t.Fatalf("normalize: %v", err)
	}
	if !strings.Contains(out, "sslmode=disable") {
		t.Fatalf("expected sslmode=disable preserved, got %q", redact(out))
	}
}

func TestNormalizePostgresURL_RejectsBadScheme(t *testing.T) {
	_, err := normalizePostgresURL("mysql://localhost/db")
	if err == nil {
		t.Fatal("expected error for mysql scheme")
	}
}

func redact(s string) string {
	if i := strings.Index(s, "@"); i > 0 {
		if j := strings.Index(s, "://"); j >= 0 && j+3 < i {
			return s[:j+3] + "***" + s[i:]
		}
	}
	return s
}
