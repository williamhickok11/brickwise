package com.brickwise.service;

import com.brickwise.config.BrickwiseProperties;
import com.brickwise.error.ApiException;
import com.brickwise.model.EmailTimeEntriesRequest;
import com.brickwise.model.TimeEntry;
import jakarta.mail.MessagingException;
import jakarta.mail.internet.MimeMessage;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.core.io.ByteArrayResource;
import org.springframework.mail.javamail.JavaMailSender;
import org.springframework.mail.javamail.MimeMessageHelper;
import org.springframework.stereotype.Service;

import java.nio.charset.StandardCharsets;
import java.util.List;

/**
 * Emails a REPS CSV export (or a short summary when the list is empty).
 *
 * <p>Requires {@code MAIL_HOST} (and usually username/password) to be configured.
 */
@Service
public class TimeEntryMailService {
    private final TimeEntryService timeEntryService;
    private final JavaMailSender mailSender;
    private final BrickwiseProperties props;
    private final String mailHost;

    public TimeEntryMailService(
            TimeEntryService timeEntryService,
            JavaMailSender mailSender,
            BrickwiseProperties props,
            @Value("${spring.mail.host:}") String mailHost) {
        this.timeEntryService = timeEntryService;
        this.mailSender = mailSender;
        this.props = props;
        this.mailHost = mailHost == null ? "" : mailHost.trim();
    }

    public void emailExport(EmailTimeEntriesRequest req, Integer propertyIdFilter) {
        if (req.getTo() == null || req.getTo().isBlank()) {
            throw ApiException.validation("to: to is required");
        }
        if (!isMailConfigured()) {
            throw ApiException.validation(
                    "email: mail is not configured; set MAIL_HOST (and MAIL_USERNAME / MAIL_PASSWORD) in the environment");
        }

        List<TimeEntry> entries = timeEntryService.list(
                propertyIdFilter, req.getStartDate(), req.getEndDate(), req.getCategory());
        String csv = timeEntryService.toCsv(entries);
        String subject = (req.getSubject() == null || req.getSubject().isBlank())
                ? "Brickwise REPS export (" + entries.size() + " entries)"
                : req.getSubject().trim();
        String summary = buildSummary(entries);

        try {
            MimeMessage message = mailSender.createMimeMessage();
            MimeMessageHelper helper = new MimeMessageHelper(message, true, StandardCharsets.UTF_8.name());
            helper.setFrom(props.getMailFrom());
            helper.setTo(req.getTo().trim());
            helper.setSubject(subject);
            helper.setText(summary, false);
            helper.addAttachment(
                    "time-entries.csv",
                    new ByteArrayResource(csv.getBytes(StandardCharsets.UTF_8)),
                    "text/csv");
            mailSender.send(message);
        } catch (MessagingException | RuntimeException e) {
            throw ApiException.internal("email: failed to send: " + e.getMessage());
        }
    }

    public boolean isMailConfigured() {
        return !mailHost.isEmpty();
    }

    private static String buildSummary(List<TimeEntry> entries) {
        if (entries.isEmpty()) {
            return "No REPS time entries matched the export filters. An empty CSV is attached.";
        }
        double hours = 0;
        double miles = 0;
        for (TimeEntry e : entries) {
            hours += e.getHours();
            miles += e.getMileage();
        }
        return "Brickwise REPS export\n"
                + "Entries: " + entries.size() + "\n"
                + "Total hours: " + String.format("%.2f", hours) + "\n"
                + "Total mileage: " + String.format("%.1f", miles) + "\n\n"
                + "The full CSV is attached as time-entries.csv.";
    }
}
