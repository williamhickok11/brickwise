package com.brickwise.repository;

import com.brickwise.config.JdbcSupport;
import com.brickwise.model.Mortgage;
import org.springframework.beans.factory.annotation.Qualifier;
import org.springframework.dao.EmptyResultDataAccessException;
import org.springframework.jdbc.core.JdbcTemplate;
import org.springframework.jdbc.core.RowMapper;
import org.springframework.stereotype.Repository;

import java.sql.Date;
import java.util.List;
import java.util.Optional;

@Repository
public class MortgageRepository {
    private final JdbcTemplate jdbc;
    private final boolean postgres;

    private static final RowMapper<Mortgage> MAPPER = (rs, i) -> {
        Mortgage m = new Mortgage();
        m.setId(rs.getInt("id"));
        m.setPropertyId(rs.getInt("property_id"));
        m.setOriginalLoan(rs.getDouble("original_loan"));
        m.setCurrentBalance(rs.getDouble("current_balance"));
        m.setInterestRate(rs.getDouble("interest_rate"));
        m.setMonthlyPayment(rs.getDouble("monthly_payment"));
        m.setMonthlyPi(rs.getDouble("monthly_pi"));
        m.setBalanceAsOf(JdbcSupport.localDate(rs, "balance_as_of"));
        m.setLender(nullToEmpty(rs.getString("lender")));
        m.setNotes(nullToEmpty(rs.getString("notes")));
        m.setCreatedAt(JdbcSupport.instant(rs, "created_at"));
        m.setUpdatedAt(JdbcSupport.instant(rs, "updated_at"));
        return m;
    };

    private static String nullToEmpty(String s) { return s == null ? "" : s; }

    public MortgageRepository(JdbcTemplate jdbc, @Qualifier("usingPostgres") boolean usingPostgres) {
        this.jdbc = jdbc;
        this.postgres = usingPostgres;
    }

    public Optional<Mortgage> findByPropertyId(int propertyId) {
        try {
            return Optional.ofNullable(jdbc.queryForObject(
                    "SELECT * FROM mortgages WHERE property_id = ?", MAPPER, propertyId));
        } catch (EmptyResultDataAccessException e) {
            return Optional.empty();
        }
    }

    public List<Mortgage> findAll() {
        return jdbc.query("SELECT * FROM mortgages", MAPPER);
    }

    public Mortgage upsert(int propertyId, Mortgage m) {
        Optional<Mortgage> existing = findByPropertyId(propertyId);
        Date balanceAsOf = m.getBalanceAsOf() == null ? null : Date.valueOf(m.getBalanceAsOf());
        if (existing.isPresent()) {
            jdbc.update("""
                    UPDATE mortgages SET original_loan=?, current_balance=?, interest_rate=?, monthly_payment=?,
                    monthly_pi=?, balance_as_of=?, lender=?, notes=?, updated_at=CURRENT_TIMESTAMP
                    WHERE property_id=?
                    """, m.getOriginalLoan(), m.getCurrentBalance(), m.getInterestRate(), m.getMonthlyPayment(),
                    m.getMonthlyPi(), balanceAsOf, m.getLender(), m.getNotes(), propertyId);
        } else {
            jdbc.update("""
                    INSERT INTO mortgages (property_id, original_loan, current_balance, interest_rate, monthly_payment,
                    monthly_pi, balance_as_of, lender, notes)
                    VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
                    """, propertyId, m.getOriginalLoan(), m.getCurrentBalance(), m.getInterestRate(),
                    m.getMonthlyPayment(), m.getMonthlyPi(), balanceAsOf, m.getLender(), m.getNotes());
        }
        return findByPropertyId(propertyId).orElseThrow();
    }

    public void deleteByPropertyId(int propertyId) {
        int n = jdbc.update("DELETE FROM mortgages WHERE property_id = ?", propertyId);
        if (n == 0) throw new IllegalStateException("mortgage not found");
    }
}
