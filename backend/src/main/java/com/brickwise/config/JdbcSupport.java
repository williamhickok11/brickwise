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
        String s = v.toString().trim();
        if (s.isEmpty()) return null;
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
        String s = v.toString().trim();
        if (s.isEmpty()) return null;
        return LocalDate.parse(s.substring(0, Math.min(10, s.length())));
    }

    public static boolean boolVal(ResultSet rs, String col) throws SQLException {
        Object v = rs.getObject(col);
        if (v == null) return false;
        if (v instanceof Boolean b) return b;
        if (v instanceof Number n) return n.intValue() != 0;
        return "true".equalsIgnoreCase(v.toString()) || "1".equals(v.toString());
    }
}
