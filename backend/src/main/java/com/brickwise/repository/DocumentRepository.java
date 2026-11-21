package com.brickwise.repository;

import com.brickwise.config.JdbcSupport;
import com.brickwise.model.Document;
import org.springframework.beans.factory.annotation.Qualifier;
import org.springframework.dao.EmptyResultDataAccessException;
import org.springframework.jdbc.core.JdbcTemplate;
import org.springframework.jdbc.core.RowMapper;
import org.springframework.stereotype.Repository;

import java.util.Optional;

@Repository
public class DocumentRepository {
    private final JdbcTemplate jdbc;
    private final boolean postgres;

    private static final RowMapper<Document> MAPPER = (rs, i) -> {
        Document d = new Document();
        d.setId(rs.getInt("id"));
        d.setResidenceId(rs.getInt("residence_id"));
        d.setFileName(rs.getString("file_name"));
        d.setContentType(rs.getString("content_type"));
        d.setSize(rs.getLong("size"));
        d.setCreatedAt(JdbcSupport.instant(rs, "created_at"));
        return d;
    };

    public DocumentRepository(JdbcTemplate jdbc, @Qualifier("usingPostgres") boolean usingPostgres) {
        this.jdbc = jdbc;
        this.postgres = usingPostgres;
    }

    public Optional<Document> findByResidenceId(int residenceId) {
        try {
            return Optional.ofNullable(jdbc.queryForObject(
                    "SELECT * FROM documents WHERE residence_id = ?", MAPPER, residenceId));
        } catch (EmptyResultDataAccessException e) {
            return Optional.empty();
        }
    }

    public Document upsert(Document doc) {
        Optional<Document> existing = findByResidenceId(doc.getResidenceId());
        if (existing.isPresent()) {
            jdbc.update("""
                    UPDATE documents SET file_name=?, content_type=?, size=?, created_at=CURRENT_TIMESTAMP
                    WHERE residence_id=?
                    """, doc.getFileName(), doc.getContentType(), doc.getSize(), doc.getResidenceId());
        } else {
            jdbc.update("""
                    INSERT INTO documents (residence_id, file_name, content_type, size)
                    VALUES (?, ?, ?, ?)
                    """, doc.getResidenceId(), doc.getFileName(), doc.getContentType(), doc.getSize());
        }
        return findByResidenceId(doc.getResidenceId()).orElseThrow();
    }

    public void deleteByResidenceId(int residenceId) {
        jdbc.update("DELETE FROM documents WHERE residence_id = ?", residenceId);
    }
}
