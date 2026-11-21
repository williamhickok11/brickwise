package com.brickwise.repository;

import com.brickwise.config.JdbcSupport;
import com.brickwise.model.Residence;
import org.springframework.beans.factory.annotation.Qualifier;
import org.springframework.dao.EmptyResultDataAccessException;
import org.springframework.jdbc.core.JdbcTemplate;
import org.springframework.jdbc.core.RowMapper;
import org.springframework.jdbc.support.GeneratedKeyHolder;
import org.springframework.jdbc.support.KeyHolder;
import org.springframework.stereotype.Repository;

import java.sql.Date;
import java.sql.PreparedStatement;
import java.sql.Statement;
import java.util.ArrayList;
import java.util.List;
import java.util.Optional;

@Repository
public class ResidenceRepository {
    private final JdbcTemplate jdbc;
    private final boolean postgres;

    private static final RowMapper<Residence> MAPPER = (rs, i) -> {
        Residence r = new Residence();
        r.setId(rs.getInt("id"));
        r.setPropertyId(rs.getInt("property_id"));
        r.setName(rs.getString("name"));
        r.setPhone(rs.getString("phone") == null ? "" : rs.getString("phone"));
        r.setEmail(rs.getString("email") == null ? "" : rs.getString("email"));
        r.setStartDate(JdbcSupport.localDate(rs, "start_date"));
        r.setEndDate(JdbcSupport.localDate(rs, "end_date"));
        r.setActive(JdbcSupport.boolVal(rs, "is_active"));
        r.setMonthlyRent(rs.getDouble("monthly_rent"));
        r.setCreatedAt(JdbcSupport.instant(rs, "created_at"));
        r.setUpdatedAt(JdbcSupport.instant(rs, "updated_at"));
        return r;
    };

    private static final String SELECT = """
            SELECT id, property_id, name, phone, email, start_date, end_date, is_active, monthly_rent, created_at, updated_at
            FROM residences
            """;

    public ResidenceRepository(JdbcTemplate jdbc, @Qualifier("usingPostgres") boolean usingPostgres) {
        this.jdbc = jdbc;
        this.postgres = usingPostgres;
    }

    public Residence create(Residence residence) {
        if (postgres) {
            return jdbc.queryForObject("""
                    INSERT INTO residences (property_id, name, phone, email, start_date, end_date, is_active, monthly_rent)
                    VALUES (?, ?, ?, ?, ?, ?, ?, ?)
                    RETURNING id, property_id, name, phone, email, start_date, end_date, is_active, monthly_rent, created_at, updated_at
                    """, MAPPER,
                    residence.getPropertyId(), residence.getName(), residence.getPhone(), residence.getEmail(),
                    Date.valueOf(residence.getStartDate()),
                    residence.getEndDate() == null ? null : Date.valueOf(residence.getEndDate()),
                    residence.isActive(), residence.getMonthlyRent());
        }
        KeyHolder keys = new GeneratedKeyHolder();
        jdbc.update(con -> {
            PreparedStatement ps = con.prepareStatement("""
                    INSERT INTO residences (property_id, name, phone, email, start_date, end_date, is_active, monthly_rent)
                    VALUES (?, ?, ?, ?, ?, ?, ?, ?)
                    """, Statement.RETURN_GENERATED_KEYS);
            ps.setInt(1, residence.getPropertyId());
            ps.setString(2, residence.getName());
            ps.setString(3, residence.getPhone());
            ps.setString(4, residence.getEmail());
            ps.setDate(5, Date.valueOf(residence.getStartDate()));
            if (residence.getEndDate() == null) ps.setObject(6, null);
            else ps.setDate(6, Date.valueOf(residence.getEndDate()));
            ps.setInt(7, residence.isActive() ? 1 : 0);
            ps.setDouble(8, residence.getMonthlyRent());
            return ps;
        }, keys);
        return findById(keys.getKey().intValue()).orElseThrow();
    }

    public Optional<Residence> findById(int id) {
        try {
            return Optional.ofNullable(jdbc.queryForObject(SELECT + " WHERE id = ?", MAPPER, id));
        } catch (EmptyResultDataAccessException e) {
            return Optional.empty();
        }
    }

    public List<Residence> findAll(Integer propertyId, Boolean isActive) {
        StringBuilder sql = new StringBuilder(SELECT + " WHERE 1=1");
        List<Object> args = new ArrayList<>();
        if (propertyId != null) {
            sql.append(" AND property_id = ?");
            args.add(propertyId);
        }
        if (isActive != null) {
            sql.append(" AND is_active = ?");
            args.add(postgres ? isActive : (isActive ? 1 : 0));
        }
        sql.append(" ORDER BY start_date DESC, id DESC");
        return jdbc.query(sql.toString(), MAPPER, args.toArray());
    }

    public void update(Residence residence) {
        int n = jdbc.update("""
                UPDATE residences
                SET property_id = ?, name = ?, phone = ?, email = ?, start_date = ?, end_date = ?,
                    is_active = ?, monthly_rent = ?, updated_at = CURRENT_TIMESTAMP
                WHERE id = ?
                """,
                residence.getPropertyId(), residence.getName(), residence.getPhone(), residence.getEmail(),
                Date.valueOf(residence.getStartDate()),
                residence.getEndDate() == null ? null : Date.valueOf(residence.getEndDate()),
                postgres ? residence.isActive() : (residence.isActive() ? 1 : 0),
                residence.getMonthlyRent(), residence.getId());
        if (n == 0) throw new IllegalStateException("residence not found");
    }
}
