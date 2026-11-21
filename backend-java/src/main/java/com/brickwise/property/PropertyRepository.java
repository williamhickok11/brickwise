package com.brickwise.property;

import java.time.LocalDateTime;
import java.time.OffsetDateTime;
import java.time.ZoneOffset;
import java.time.format.DateTimeFormatter;
import java.time.format.DateTimeFormatterBuilder;
import java.time.format.DateTimeParseException;
import java.time.temporal.ChronoField;
import java.time.temporal.TemporalAccessor;
import java.util.List;
import org.springframework.jdbc.core.JdbcTemplate;
import org.springframework.jdbc.core.RowMapper;
import org.springframework.stereotype.Repository;

// Dapper Query<T> / Go SQLPropertyRepository.GetAll in backend/internal/repository/repository.go.
@Repository
public class PropertyRepository {

    private static final DateTimeFormatter SQLITE_TIMESTAMP = new DateTimeFormatterBuilder()
            .appendPattern("yyyy-MM-dd HH:mm:ss")
            .optionalStart()
            .appendFraction(ChronoField.NANO_OF_SECOND, 1, 9, true)
            .optionalEnd()
            .optionalStart()
            .appendOffset("+HH:MM", "Z")
            .optionalEnd()
            .toFormatter();

    private static final String LIST_SQL = """
            SELECT id, name, address, property_type, default_mileage, current_value,
                   value_updated_at, created_at, updated_at
            FROM properties
            ORDER BY created_at DESC
            """;

    private final JdbcTemplate jdbc;

    public PropertyRepository(JdbcTemplate jdbc) {
        this.jdbc = jdbc;
    }

    public List<Property> findAll() {
        return jdbc.query(LIST_SQL, rowMapper());
    }

    private static RowMapper<Property> rowMapper() {
        return (rs, rowNum) -> new Property(
                rs.getInt("id"),
                rs.getString("name"),
                rs.getString("address"),
                rs.getString("property_type"),
                rs.getDouble("default_mileage"),
                rs.getDouble("current_value"),
                parseTimestamp(rs.getString("value_updated_at")),
                parseTimestamp(rs.getString("created_at")),
                parseTimestamp(rs.getString("updated_at"))
        );
    }

    // SQLite stores "2026-03-30 20:17:18" with no offset. Treat that as UTC so JSON matches Go's RFC3339.
    static OffsetDateTime parseTimestamp(String raw) {
        if (raw == null || raw.isBlank()) {
            return null;
        }
        String value = raw.trim();
        try {
            return OffsetDateTime.parse(value);
        } catch (DateTimeParseException ignored) {
            TemporalAccessor parsed = SQLITE_TIMESTAMP.parse(value);
            if (parsed.isSupported(ChronoField.OFFSET_SECONDS)) {
                return OffsetDateTime.from(parsed);
            }
            return LocalDateTime.from(parsed).atOffset(ZoneOffset.UTC);
        }
    }
}
