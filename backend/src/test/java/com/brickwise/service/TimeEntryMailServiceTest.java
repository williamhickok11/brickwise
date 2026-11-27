package com.brickwise.service;

import com.brickwise.config.BrickwiseProperties;
import com.brickwise.error.ApiException;
import com.brickwise.model.EmailTimeEntriesRequest;
import org.junit.jupiter.api.Test;
import org.springframework.mail.javamail.JavaMailSender;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertThrows;
import static org.mockito.Mockito.mock;

class TimeEntryMailServiceTest {
    @Test
    void emailExport_requiresMailHost() {
        TimeEntryMailService mail = new TimeEntryMailService(
                mock(TimeEntryService.class),
                mock(JavaMailSender.class),
                new BrickwiseProperties(),
                "");
        EmailTimeEntriesRequest req = new EmailTimeEntriesRequest();
        req.setTo("owner@example.com");
        ApiException ex = assertThrows(ApiException.class, () -> mail.emailExport(req, null));
        assertEquals(400, ex.getStatus());
        assertTrueContains(ex.getMessage(), "MAIL_HOST");
    }

    @Test
    void emailExport_requiresTo() {
        TimeEntryMailService mail = new TimeEntryMailService(
                mock(TimeEntryService.class),
                mock(JavaMailSender.class),
                new BrickwiseProperties(),
                "smtp.example.com");
        ApiException ex = assertThrows(ApiException.class, () -> mail.emailExport(new EmailTimeEntriesRequest(), null));
        assertEquals(400, ex.getStatus());
    }

    private static void assertTrueContains(String haystack, String needle) {
        if (haystack == null || !haystack.contains(needle)) {
            throw new AssertionError("expected message to contain " + needle + " but was: " + haystack);
        }
    }
}
