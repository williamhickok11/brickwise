package com.brickwise.web;

import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertNull;
import static org.junit.jupiter.api.Assertions.assertThrows;

class TimeEntryControllerTest {
    @Test
    void parsePropertyIdFilter_omittedMeansAll() {
        assertNull(TimeEntryController.parsePropertyIdFilter(null));
        assertNull(TimeEntryController.parsePropertyIdFilter(""));
        assertNull(TimeEntryController.parsePropertyIdFilter("  "));
    }

    @Test
    void parsePropertyIdFilter_nullOrZeroMeansGeneral() {
        assertEquals(0, TimeEntryController.parsePropertyIdFilter("null"));
        assertEquals(0, TimeEntryController.parsePropertyIdFilter("NULL"));
        assertEquals(0, TimeEntryController.parsePropertyIdFilter("0"));
    }

    @Test
    void parsePropertyIdFilter_numericProperty() {
        assertEquals(12, TimeEntryController.parsePropertyIdFilter("12"));
    }

    @Test
    void parsePropertyIdFilter_rejectsGarbage() {
        assertThrows(NumberFormatException.class, () -> TimeEntryController.parsePropertyIdFilter("abc"));
    }
}
