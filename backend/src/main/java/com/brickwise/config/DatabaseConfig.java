package com.brickwise.config;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.boot.context.properties.EnableConfigurationProperties;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.core.env.Environment;
import org.springframework.jdbc.core.JdbcTemplate;
import org.springframework.jdbc.datasource.DriverManagerDataSource;

import javax.sql.DataSource;
import java.net.URI;
import java.sql.Connection;
import java.sql.SQLException;
import java.sql.Statement;

/**
 * Builds the DataSource and applies schema on startup (same approach as the Go {@code db.InitDB}).
 *
 * <p><b>C# analogy:</b> Configuring a connection string in {@code appsettings.json} and running
 * EnsureCreated / migrations at startup. We use {@link JdbcTemplate} (think Dapper) rather than
 * a full ORM so the SQL stays visible — helpful while learning Java data access.
 *
 * <p>If {@code DATABASE_URL} is set → Postgres; otherwise → SQLite file at {@code brickwise.db-path}.
 */
@Configuration
@EnableConfigurationProperties(BrickwiseProperties.class)
public class DatabaseConfig {
    private static final Logger log = LoggerFactory.getLogger(DatabaseConfig.class);

    @Bean
    public DataSource dataSource(Environment env, BrickwiseProperties props) throws SQLException {
        String databaseUrl = env.getProperty("DATABASE_URL", "").trim();
        DriverManagerDataSource ds = new DriverManagerDataSource();

        if (!databaseUrl.isEmpty()) {
            String normalized = normalizePostgresUrl(databaseUrl);
            ds.setDriverClassName("org.postgresql.Driver");
            ds.setUrl(toJdbcPostgresUrl(normalized));
            log.info("Using Postgres database (DATABASE_URL)");
        } else {
            String path = props.getDbPath();
            ds.setDriverClassName("org.sqlite.JDBC");
            ds.setUrl("jdbc:sqlite:" + path);
            log.info("Using SQLite database path={}", path);
        }

        try (Connection conn = ds.getConnection(); Statement st = conn.createStatement()) {
            if (databaseUrl.isEmpty()) {
                st.execute("PRAGMA foreign_keys = ON");
            }
            SchemaBootstrap.createTables(st, !databaseUrl.isEmpty());
        }
        return ds;
    }

    @Bean
    public JdbcTemplate jdbcTemplate(DataSource dataSource) {
        return new JdbcTemplate(dataSource);
    }

    @Bean
    public boolean usingPostgres(Environment env) {
        return !env.getProperty("DATABASE_URL", "").trim().isEmpty();
    }

    static String normalizePostgresUrl(String raw) {
        String url = raw.trim();
        if (url.startsWith("postgresql://")) {
            url = "postgres://" + url.substring("postgresql://".length());
        }
        URI u = URI.create(url);
        if (!"postgres".equals(u.getScheme()) && !"postgresql".equals(u.getScheme())) {
            throw new IllegalArgumentException("unsupported scheme (expected postgres://)");
        }
        String host = u.getHost() == null ? "" : u.getHost().toLowerCase();
        String query = u.getRawQuery() == null ? "" : u.getRawQuery();
        boolean hasSsl = query.contains("sslmode=");
        if (!hasSsl && (host.contains("supabase.co") || host.contains("supabase.com"))) {
            query = query.isEmpty() ? "sslmode=require" : query + "&sslmode=require";
            int q = url.indexOf('?');
            String base = q >= 0 ? url.substring(0, q) : url;
            url = base + "?" + query;
        }
        return url;
    }

    static String toJdbcPostgresUrl(String postgresUrl) {
        // postgres://user:pass@host:5432/db?sslmode=require → jdbc:postgresql://host:5432/db?user=...&password=...
        URI u = URI.create(postgresUrl.replace("postgres://", "http://").replace("postgresql://", "http://"));
        String userInfo = u.getUserInfo();
        String user = null;
        String pass = null;
        if (userInfo != null) {
            int colon = userInfo.indexOf(':');
            if (colon >= 0) {
                user = userInfo.substring(0, colon);
                pass = userInfo.substring(colon + 1);
            } else {
                user = userInfo;
            }
        }
        StringBuilder jdbc = new StringBuilder("jdbc:postgresql://");
        jdbc.append(u.getHost());
        if (u.getPort() > 0) {
            jdbc.append(':').append(u.getPort());
        }
        jdbc.append(u.getPath());
        StringBuilder q = new StringBuilder();
        if (u.getRawQuery() != null) {
            q.append(u.getRawQuery());
        }
        if (user != null) {
            if (!q.isEmpty()) q.append('&');
            q.append("user=").append(user);
        }
        if (pass != null) {
            if (!q.isEmpty()) q.append('&');
            q.append("password=").append(pass);
        }
        if (!q.isEmpty()) {
            jdbc.append('?').append(q);
        }
        return jdbc.toString();
    }
}
