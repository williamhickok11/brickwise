package com.brickwise.service;

import com.brickwise.model.TimeEntry;
import com.brickwise.repository.PropertyRepository;
import com.brickwise.repository.TimeEntryRepository;
import org.junit.jupiter.api.Test;

import java.time.LocalDate;
import java.util.List;

import static org.junit.jupiter.api.Assertions.assertTrue;
import static org.mockito.Mockito.mock;

class TimeEntryServiceCsvTest {
    @Test
    void toCsv_includesHeaderAndQuotedFields() {
        TimeEntryService service = new TimeEntryService(mock(TimeEntryRepository.class), mock(PropertyRepository.class));
        TimeEntry e = new TimeEntry();
        e.setDate(LocalDate.of(2026, 3, 15));
        e.setCategory("Property Management");
        e.setDescription("Checked unit \"A\" and porch");
        e.setHours(2.5);
        e.setNotes("ok");
        e.setMileage(10.0);
        e.setPropertyId(7);

        String csv = service.toCsv(List.of(e));
        assertTrue(csv.startsWith("Date,Activity Category,Description,Time Spent (hrs),Notes,Mileage,Property ID\n"));
        assertTrue(csv.contains("\"Checked unit \"\"A\"\" and porch\""));
        assertTrue(csv.contains(",7\n"));
    }
}
