package com.brickwise.repository;

import com.brickwise.config.JdbcSupport;
import com.brickwise.model.Property;
import org.springframework.beans.factory.annotation.Qualifier;
import org.springframework.dao.EmptyResultDataAccessException;
import org.springframework.jdbc.core.JdbcTemplate;
import org.springframework.jdbc.core.RowMapper;
import org.springframework.jdbc.support.GeneratedKeyHolder;
import org.springframework.jdbc.support.KeyHolder;
import org.springframework.stereotype.Repository;

import java.sql.PreparedStatement;
import java.sql.Statement;
import java.sql.Timestamp;
import java.time.Instant;
import java.util.List;
import java.util.Optional;

/**
 * Property persistence with {@link JdbcTemplate}.
 *
 * <p><b>C# analogy:</b> A Dapper repository class injected via DI.
 * {@code @Repository} is like registering the class in the DI container as a scoped/singleton service.
 * {@code Optional<T>} is Java's answer to nullable reference types / {@code T?} for absence of a value.
 */
@Repository
public class PropertyRepository {
    private final JdbcTemplate jdbc;
    private final boolean postgres;

    private static final RowMapper<Property> MAPPER = (rs, rowNum) -> {
        Property p = new Property();
        p.setId(rs.getInt("id"));
        p.setName(rs.getString("name"));
        p.setAddress(rs.getString("address"));
        p.setPropertyType(rs.getString("property_type"));
        p.setDefaultMileage(rs.getDouble("default_mileage"));
        p.setCurrentValue(rs.getDouble("current_value"));
        p.setValueUpdatedAt(JdbcSupport.instant(rs, "value_updated_at"));
        Object mid = rs.getObject("main_image_id");
        if (mid != null) p.setMainImageId(((Number) mid).intValue());
        p.setCreatedAt(JdbcSupport.instant(rs, "created_at"));
        p.setUpdatedAt(JdbcSupport.instant(rs, "updated_at"));
        return p;
    };

    private static final String SELECT = """
            SELECT id, name, address, property_type, default_mileage, current_value, value_updated_at,
                   (SELECT id FROM property_images WHERE property_id = properties.id AND is_main = TRUE LIMIT 1) AS main_image_id,
                   created_at, updated_at
            FROM properties
            """;

    public PropertyRepository(JdbcTemplate jdbc, @Qualifier("usingPostgres") boolean usingPostgres) {
        this.jdbc = jdbc;
        this.postgres = usingPostgres;
    }

    public Property create(Property property) {
        if (postgres) {
            return jdbc.queryForObject("""
                    INSERT INTO properties (name, address, property_type, default_mileage)
                    VALUES (?, ?, ?, ?)
                    RETURNING id, name, address, property_type, default_mileage, current_value, value_updated_at,
                      (SELECT id FROM property_images WHERE property_id = properties.id AND is_main = TRUE LIMIT 1) AS main_image_id,
                      created_at, updated_at
                    """, MAPPER, property.getName(), property.getAddress(), property.getPropertyType(), property.getDefaultMileage());
        }
        KeyHolder keys = new GeneratedKeyHolder();
        jdbc.update(con -> {
            PreparedStatement ps = con.prepareStatement(
                    "INSERT INTO properties (name, address, property_type, default_mileage) VALUES (?, ?, ?, ?)",
                    Statement.RETURN_GENERATED_KEYS);
            ps.setString(1, property.getName());
            ps.setString(2, property.getAddress());
            ps.setString(3, property.getPropertyType());
            ps.setDouble(4, property.getDefaultMileage());
            return ps;
        }, keys);
        int id = keys.getKey().intValue();
        return findById(id).orElseThrow();
    }

    public Optional<Property> findById(int id) {
        try {
            return Optional.ofNullable(jdbc.queryForObject(SELECT + " WHERE id = ?", MAPPER, id));
        } catch (EmptyResultDataAccessException e) {
            return Optional.empty();
        }
    }

    public List<Property> findAll() {
        return jdbc.query(SELECT + " ORDER BY created_at DESC", MAPPER);
    }

    public void update(Property property) {
        int n = jdbc.update("""
                UPDATE properties
                SET name = ?, address = ?, property_type = ?, default_mileage = ?, updated_at = CURRENT_TIMESTAMP
                WHERE id = ?
                """, property.getName(), property.getAddress(), property.getPropertyType(),
                property.getDefaultMileage(), property.getId());
        if (n == 0) throw new IllegalStateException("property not found");
    }

    public void updateValue(int id, double currentValue, Instant valueUpdatedAt) {
        int n = jdbc.update("""
                UPDATE properties SET current_value = ?, value_updated_at = ?, updated_at = CURRENT_TIMESTAMP
                WHERE id = ?
                """, currentValue, valueUpdatedAt == null ? null : Timestamp.from(valueUpdatedAt), id);
        if (n == 0) throw new IllegalStateException("property not found");
    }

    public void delete(int id) {
        int n = jdbc.update("DELETE FROM properties WHERE id = ?", id);
        if (n == 0) throw new IllegalStateException("property not found");
    }
}
