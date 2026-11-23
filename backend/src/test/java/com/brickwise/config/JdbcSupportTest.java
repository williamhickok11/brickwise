package com.brickwise.config;

import org.junit.jupiter.api.Test;

import java.sql.ResultSet;
import java.time.LocalDate;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.mockito.Mockito.mock;
import static org.mockito.Mockito.when;

class JdbcSupportTest {
    @Test
    void localDate_parsesSqliteUnixSecondsString() throws Exception {
        ResultSet rs = mock(ResultSet.class);
        // 2026-11-22T00:00:00Z
        when(rs.getObject("date")).thenReturn("1795305600");
        assertEquals(LocalDate.of(2026, 11, 22), JdbcSupport.localDate(rs, "date"));
    }

    @Test
    void localDate_parsesIsoDate() throws Exception {
        ResultSet rs = mock(ResultSet.class);
        when(rs.getObject("date")).thenReturn("2026-11-22");
        assertEquals(LocalDate.of(2026, 11, 22), JdbcSupport.localDate(rs, "date"));
    }
}
