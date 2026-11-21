package com.brickwise.repository;

import com.brickwise.config.JdbcSupport;
import com.brickwise.model.QboModels.PropertyQBOLink;
import com.brickwise.model.QboModels.QBOCashFlowCache;
import com.brickwise.model.QboModels.QBOConnection;
import org.springframework.dao.EmptyResultDataAccessException;
import org.springframework.jdbc.core.JdbcTemplate;
import org.springframework.jdbc.core.RowMapper;
import org.springframework.stereotype.Repository;

import java.sql.Date;
import java.sql.Timestamp;
import java.time.Instant;
import java.time.LocalDate;
import java.util.List;
import java.util.Optional;

@Repository
public class QboRepository {
    private final JdbcTemplate jdbc;

    private static final RowMapper<QBOConnection> CONN = (rs, i) -> {
        QBOConnection c = new QBOConnection();
        c.setId(rs.getInt("id"));
        c.setRealmId(rs.getString("realm_id"));
        c.setAccessToken(rs.getString("access_token"));
        c.setRefreshToken(rs.getString("refresh_token"));
        c.setAccessTokenExpiresAt(JdbcSupport.instant(rs, "access_token_expires_at"));
        c.setRefreshTokenExpiresAt(JdbcSupport.instant(rs, "refresh_token_expires_at"));
        c.setConnectedAt(JdbcSupport.instant(rs, "connected_at"));
        c.setUpdatedAt(JdbcSupport.instant(rs, "updated_at"));
        return c;
    };

    private static final RowMapper<PropertyQBOLink> LINK = (rs, i) -> {
        PropertyQBOLink l = new PropertyQBOLink();
        l.setPropertyId(rs.getInt("property_id"));
        l.setQboClassId(rs.getString("qbo_class_id"));
        l.setQboClassName(rs.getString("qbo_class_name"));
        l.setUpdatedAt(JdbcSupport.instant(rs, "updated_at"));
        return l;
    };

    private static final RowMapper<QBOCashFlowCache> CACHE = (rs, i) -> {
        QBOCashFlowCache c = new QBOCashFlowCache();
        c.setPropertyId(rs.getInt("property_id"));
        c.setPeriodStart(JdbcSupport.localDate(rs, "period_start"));
        c.setPeriodEnd(JdbcSupport.localDate(rs, "period_end"));
        c.setNetAmount(rs.getDouble("net_amount"));
        c.setSyncedAt(JdbcSupport.instant(rs, "synced_at"));
        return c;
    };

    public QboRepository(JdbcTemplate jdbc) {
        this.jdbc = jdbc;
    }

    public Optional<QBOConnection> findConnection() {
        try {
            return Optional.ofNullable(jdbc.queryForObject(
                    "SELECT * FROM qbo_connections ORDER BY id DESC LIMIT 1", CONN));
        } catch (EmptyResultDataAccessException e) {
            return Optional.empty();
        }
    }

    public void saveConnection(QBOConnection c) {
        jdbc.update("DELETE FROM qbo_connections");
        jdbc.update("""
                INSERT INTO qbo_connections (realm_id, access_token, refresh_token, access_token_expires_at,
                refresh_token_expires_at, connected_at, updated_at)
                VALUES (?, ?, ?, ?, ?, ?, ?)
                """, c.getRealmId(), c.getAccessToken(), c.getRefreshToken(),
                Timestamp.from(c.getAccessTokenExpiresAt()), Timestamp.from(c.getRefreshTokenExpiresAt()),
                Timestamp.from(c.getConnectedAt()), Timestamp.from(c.getUpdatedAt()));
    }

    public void deleteConnection() {
        jdbc.update("DELETE FROM qbo_connections");
        jdbc.update("DELETE FROM qbo_oauth_states");
    }

    public void saveOauthState(String state, Instant expiresAt) {
        jdbc.update("DELETE FROM qbo_oauth_states WHERE expires_at < CURRENT_TIMESTAMP");
        jdbc.update("INSERT INTO qbo_oauth_states (state, expires_at) VALUES (?, ?)",
                state, Timestamp.from(expiresAt));
    }

    public boolean consumeOauthState(String state) {
        int n = jdbc.update("DELETE FROM qbo_oauth_states WHERE state = ? AND expires_at >= CURRENT_TIMESTAMP", state);
        return n > 0;
    }

    public List<PropertyQBOLink> listLinks() {
        return jdbc.query("SELECT * FROM property_qbo_links ORDER BY property_id", LINK);
    }

    public void replaceLinks(List<PropertyQBOLink> links) {
        jdbc.update("DELETE FROM property_qbo_links");
        for (PropertyQBOLink l : links) {
            jdbc.update("""
                    INSERT INTO property_qbo_links (property_id, qbo_class_id, qbo_class_name, updated_at)
                    VALUES (?, ?, ?, CURRENT_TIMESTAMP)
                    """, l.getPropertyId(), l.getQboClassId(), l.getQboClassName());
        }
    }

    public List<QBOCashFlowCache> listCashFlowCache() {
        return jdbc.query("SELECT * FROM qbo_cash_flow_cache", CACHE);
    }

    public void upsertCashFlow(int propertyId, LocalDate start, LocalDate end, double net, Instant syncedAt) {
        jdbc.update("DELETE FROM qbo_cash_flow_cache WHERE property_id = ?", propertyId);
        jdbc.update("""
                INSERT INTO qbo_cash_flow_cache (property_id, period_start, period_end, net_amount, synced_at)
                VALUES (?, ?, ?, ?, ?)
                """, propertyId, Date.valueOf(start), Date.valueOf(end), net, Timestamp.from(syncedAt));
    }

    public Optional<Instant> latestSyncedAt() {
        try {
            Instant i = jdbc.query("SELECT synced_at FROM qbo_cash_flow_cache ORDER BY synced_at DESC LIMIT 1",
                    rs -> rs.next() ? JdbcSupport.instant(rs, "synced_at") : null);
            return Optional.ofNullable(i);
        } catch (Exception e) {
            return Optional.empty();
        }
    }
}
