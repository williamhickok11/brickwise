package com.brickwise.repository;

import com.brickwise.config.JdbcSupport;
import com.brickwise.model.PropertyImage;
import org.springframework.beans.factory.annotation.Qualifier;
import org.springframework.dao.EmptyResultDataAccessException;
import org.springframework.jdbc.core.JdbcTemplate;
import org.springframework.jdbc.core.RowMapper;
import org.springframework.jdbc.support.GeneratedKeyHolder;
import org.springframework.jdbc.support.KeyHolder;
import org.springframework.stereotype.Repository;

import java.sql.PreparedStatement;
import java.sql.Statement;
import java.util.List;
import java.util.Optional;

@Repository
public class PropertyImageRepository {
    private final JdbcTemplate jdbc;
    private final boolean postgres;

    private static final RowMapper<PropertyImage> MAPPER = (rs, i) -> {
        PropertyImage img = new PropertyImage();
        img.setId(rs.getInt("id"));
        img.setPropertyId(rs.getInt("property_id"));
        img.setFileName(rs.getString("file_name"));
        img.setContentType(rs.getString("content_type"));
        img.setSize(rs.getLong("size"));
        img.setMain(JdbcSupport.boolVal(rs, "is_main"));
        img.setCreatedAt(JdbcSupport.instant(rs, "created_at"));
        return img;
    };

    public PropertyImageRepository(JdbcTemplate jdbc, @Qualifier("usingPostgres") boolean usingPostgres) {
        this.jdbc = jdbc;
        this.postgres = usingPostgres;
    }

    public List<PropertyImage> listByPropertyId(int propertyId) {
        return jdbc.query("SELECT * FROM property_images WHERE property_id = ? ORDER BY is_main DESC, id ASC",
                MAPPER, propertyId);
    }

    public Optional<PropertyImage> findById(int id) {
        try {
            return Optional.ofNullable(jdbc.queryForObject("SELECT * FROM property_images WHERE id = ?", MAPPER, id));
        } catch (EmptyResultDataAccessException e) {
            return Optional.empty();
        }
    }

    public PropertyImage create(PropertyImage image) {
        if (postgres) {
            return jdbc.queryForObject("""
                    INSERT INTO property_images (property_id, file_name, content_type, size, is_main)
                    VALUES (?, ?, ?, ?, ?) RETURNING *
                    """, MAPPER, image.getPropertyId(), image.getFileName(), image.getContentType(),
                    image.getSize(), image.isMain());
        }
        KeyHolder keys = new GeneratedKeyHolder();
        jdbc.update(con -> {
            PreparedStatement ps = con.prepareStatement("""
                    INSERT INTO property_images (property_id, file_name, content_type, size, is_main)
                    VALUES (?, ?, ?, ?, ?)
                    """, Statement.RETURN_GENERATED_KEYS);
            ps.setInt(1, image.getPropertyId());
            ps.setString(2, image.getFileName());
            ps.setString(3, image.getContentType());
            ps.setLong(4, image.getSize());
            ps.setInt(5, image.isMain() ? 1 : 0);
            return ps;
        }, keys);
        return findById(keys.getKey().intValue()).orElseThrow();
    }

    public void updateSize(int id, long size) {
        jdbc.update("UPDATE property_images SET size = ? WHERE id = ?", size, id);
    }

    public void clearMain(int propertyId) {
        jdbc.update("UPDATE property_images SET is_main = ? WHERE property_id = ?", postgres ? false : 0, propertyId);
    }

    public void setMain(int imageId) {
        jdbc.update("UPDATE property_images SET is_main = ? WHERE id = ?", postgres ? true : 1, imageId);
    }

    public void delete(int id) {
        int n = jdbc.update("DELETE FROM property_images WHERE id = ?", id);
        if (n == 0) throw new IllegalStateException("image not found");
    }

    public int countByPropertyId(int propertyId) {
        Integer n = jdbc.queryForObject("SELECT COUNT(*) FROM property_images WHERE property_id = ?", Integer.class, propertyId);
        return n == null ? 0 : n;
    }

    public List<PropertyImage> listIdsByProperty(int propertyId) {
        return listByPropertyId(propertyId);
    }
}
