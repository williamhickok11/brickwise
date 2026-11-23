package com.brickwise.config;

import java.sql.ResultSet;
import java.sql.SQLException;
import java.sql.Timestamp;
import java.time.Instant;
import java.time.LocalDate;
import java.time.LocalDateTime;
import java.time.ZoneOffset;
import java.time.format.DateTimeFormatter;
import java.time.format.DateTimeParseException;

/**
 * Helpers for reading SQLite/Postgres timestamps via JDBC.
 *
 * <p><b>C# analogy:</b> Small extension methods you might write around {@code SqlDataReader}.
 */
public final class JdbcSupport {
    private JdbcSupport() {}

    public static Instant instant(ResultSet rs, String col) throws SQLException {
        Object v = rs.getObject(col);
        if (v == null) return null;
        if (v instanceof Timestamp ts) return ts.toInstant();
        if (v instanceof Number n) return epochToInstant(n.longValue());
        String s = v.toString().trim();
        if (s.isEmpty()) return null;
        if (isAllDigits(s)) return epochToInstant(Long.parseLong(s));
        try {
            return Instant.parse(s.contains("T") ? s.replace(" ", "T") + (s.endsWith("Z") ? "" : "Z") : s + "T00:00:00Z");
        } catch (DateTimeParseException e) {
            try {
                LocalDateTime ldt = LocalDateTime.parse(s.replace(" ", "T"), DateTimeFormatter.ISO_LOCAL_DATE_TIME);
                return ldt.toInstant(ZoneOffset.UTC);
            } catch (DateTimeParseException e2) {
                LocalDate d = LocalDate.parse(s.substring(0, Math.min(10, s.length())));
                return d.atStartOfDay().toInstant(ZoneOffset.UTC);
            }
        }
    }

    public static LocalDate localDate(ResultSet rs, String col) throws SQLException {
        Object v = rs.getObject(col);
        if (v == null) return null;
        if (v instanceof java.sql.Date d) return d.toLocalDate();
        if (v instanceof LocalDate ld) return ld;
        if (v instanceof Number n) {
            return epochToInstant(n.longValue()).atZone(ZoneOffset.UTC).toLocalDate();
        }
        String s = v.toString().trim();
        if (s.isEmpty()) return null;
        // SQLite JDBC often returns DATE columns as unix-seconds strings (e.g. "1795305600")
        if (isAllDigits(s)) {
            return epochToInstant(Long.parseLong(s)).atZone(ZoneOffset.UTC).toLocalDate();
        }
        return LocalDate.parse(s.substring(0, Math.min(10, s.length())));
    }

    public static boolean boolVal(ResultSet rs, String col) throws SQLException {
        Object v = rs.getObject(col);
        if (v == null) return false;
        if (v instanceof Boolean b) return b;
        if (v instanceof Number n) return n.intValue() != 0;
        return "true".equalsIgnoreCase(v.toString()) || "1".equals(v.toString());
    }

    private static Instant epochToInstant(long epoch) {
        // Heuristic: values past year ~2286 in seconds are treated as millis
        if (epoch > 10_000_000_000L) {
            return Instant.ofEpochMilli(epoch);
        }
        return Instant.ofEpochSecond(epoch);
    }

    private static boolean isAllDigits(String s) {
        if (s.isEmpty()) return false;
        for (int i = 0; i < s.length(); i++) {
            if (!Character.isDigit(s.charAt(i))) return false;
        }
        return true;
    }
}
