package com.brickwise.repository;

import com.brickwise.config.JdbcSupport;
import com.brickwise.model.TimeEntry;
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
public class TimeEntryRepository {
    private final JdbcTemplate jdbc;
    private final boolean postgres;

    private static final RowMapper<TimeEntry> MAPPER = (rs, i) -> {
        TimeEntry t = new TimeEntry();
        t.setId(rs.getInt("id"));
        Object pid = rs.getObject("property_id");
        if (pid != null) t.setPropertyId(((Number) pid).intValue());
        t.setDate(JdbcSupport.localDate(rs, "date"));
        t.setCategory(rs.getString("category"));
        t.setDescription(rs.getString("description"));
        t.setHours(rs.getDouble("hours"));
        t.setNotes(rs.getString("notes") == null ? "" : rs.getString("notes"));
        t.setMileage(rs.getDouble("mileage"));
        t.setFullDrive(JdbcSupport.boolVal(rs, "full_drive"));
        t.setCreatedAt(JdbcSupport.instant(rs, "created_at"));
        t.setUpdatedAt(JdbcSupport.instant(rs, "updated_at"));
        return t;
    };

    public TimeEntryRepository(JdbcTemplate jdbc, @Qualifier("usingPostgres") boolean usingPostgres) {
        this.jdbc = jdbc;
        this.postgres = usingPostgres;
    }

    public TimeEntry create(TimeEntry entry) {
        if (postgres) {
            return jdbc.queryForObject("""
                    INSERT INTO time_entries (property_id, date, category, description, hours, notes, mileage, full_drive)
                    VALUES (?, ?, ?, ?, ?, ?, ?, ?)
                    RETURNING *
                    """, MAPPER, entry.getPropertyId(), Date.valueOf(entry.getDate()), entry.getCategory(),
                    entry.getDescription(), entry.getHours(), entry.getNotes(), entry.getMileage(), entry.isFullDrive());
        }
        KeyHolder keys = new GeneratedKeyHolder();
        jdbc.update(con -> {
            PreparedStatement ps = con.prepareStatement("""
                    INSERT INTO time_entries (property_id, date, category, description, hours, notes, mileage, full_drive)
                    VALUES (?, ?, ?, ?, ?, ?, ?, ?)
                    """, Statement.RETURN_GENERATED_KEYS);
            if (entry.getPropertyId() == null) ps.setObject(1, null);
            else ps.setInt(1, entry.getPropertyId());
            ps.setDate(2, Date.valueOf(entry.getDate()));
            ps.setString(3, entry.getCategory());
            ps.setString(4, entry.getDescription());
            ps.setDouble(5, entry.getHours());
            ps.setString(6, entry.getNotes());
            ps.setDouble(7, entry.getMileage());
            ps.setInt(8, entry.isFullDrive() ? 1 : 0);
            return ps;
        }, keys);
        return findById(keys.getKey().intValue()).orElseThrow();
    }

    public Optional<TimeEntry> findById(int id) {
        try {
            return Optional.ofNullable(jdbc.queryForObject("SELECT * FROM time_entries WHERE id = ?", MAPPER, id));
        } catch (EmptyResultDataAccessException e) {
            return Optional.empty();
        }
    }

    public List<TimeEntry> findAll(Integer propertyId, String startDate, String endDate, String category) {
        StringBuilder sql = new StringBuilder("SELECT * FROM time_entries WHERE 1=1");
        List<Object> args = new ArrayList<>();
        if (propertyId != null) { sql.append(" AND property_id = ?"); args.add(propertyId); }
        if (startDate != null && !startDate.isBlank()) { sql.append(" AND date >= ?"); args.add(startDate); }
        if (endDate != null && !endDate.isBlank()) { sql.append(" AND date <= ?"); args.add(endDate); }
        if (category != null && !category.isBlank()) { sql.append(" AND category = ?"); args.add(category); }
        sql.append(" ORDER BY date DESC, id DESC");
        return jdbc.query(sql.toString(), MAPPER, args.toArray());
    }

    public void update(TimeEntry entry) {
        int n = jdbc.update("""
                UPDATE time_entries SET property_id=?, date=?, category=?, description=?, hours=?, notes=?,
                mileage=?, full_drive=?, updated_at=CURRENT_TIMESTAMP WHERE id=?
                """,
                entry.getPropertyId(), Date.valueOf(entry.getDate()), entry.getCategory(), entry.getDescription(),
                entry.getHours(), entry.getNotes(), entry.getMileage(),
                postgres ? entry.isFullDrive() : (entry.isFullDrive() ? 1 : 0), entry.getId());
        if (n == 0) throw new IllegalStateException("time entry not found");
    }

    public void delete(int id) {
        int n = jdbc.update("DELETE FROM time_entries WHERE id = ?", id);
        if (n == 0) throw new IllegalStateException("time entry not found");
    }
}
