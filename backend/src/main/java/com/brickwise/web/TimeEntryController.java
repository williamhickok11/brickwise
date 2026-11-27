package com.brickwise.web;

import com.brickwise.model.CreateTimeEntryRequest;
import com.brickwise.model.EmailTimeEntriesRequest;
import com.brickwise.model.TimeEntry;
import com.brickwise.model.UpdateTimeEntryRequest;
import com.brickwise.service.TimeEntryMailService;
import com.brickwise.service.TimeEntryService;
import org.springframework.http.HttpHeaders;
import org.springframework.http.HttpStatus;
import org.springframework.http.MediaType;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.nio.charset.StandardCharsets;
import java.util.LinkedHashMap;
import java.util.List;
import java.util.Map;

@RestController
@RequestMapping("/api/v1/time-entries")
public class TimeEntryController {
    private final TimeEntryService service;
    private final TimeEntryMailService mailService;

    public TimeEntryController(TimeEntryService service, TimeEntryMailService mailService) {
        this.service = service;
        this.mailService = mailService;
    }

    @GetMapping({"", "/"})
    public List<TimeEntry> list(
            @RequestParam(value = "property_id", required = false) String propertyId,
            @RequestParam(value = "start_date", required = false) String startDate,
            @RequestParam(value = "end_date", required = false) String endDate,
            @RequestParam(value = "category", required = false) String category) {
        return service.list(parsePropertyIdFilter(propertyId), startDate, endDate, category);
    }

    /** CSV download of REPS time entries (optional filters match the list endpoint). */
    @GetMapping("/export")
    public ResponseEntity<byte[]> export(
            @RequestParam(value = "property_id", required = false) String propertyId,
            @RequestParam(value = "start_date", required = false) String startDate,
            @RequestParam(value = "end_date", required = false) String endDate,
            @RequestParam(value = "category", required = false) String category) {
        List<TimeEntry> entries = service.list(parsePropertyIdFilter(propertyId), startDate, endDate, category);
        byte[] body = service.toCsv(entries).getBytes(StandardCharsets.UTF_8);
        return ResponseEntity.ok()
                .header(HttpHeaders.CONTENT_DISPOSITION, "attachment; filename=\"time-entries.csv\"")
                .contentType(MediaType.parseMediaType("text/csv"))
                .body(body);
    }

    /**
     * Emails a REPS CSV attachment (plus a short hour/mileage summary in the body).
     * Requires {@code MAIL_HOST} to be configured.
     */
    @PostMapping("/export/email")
    public Map<String, Object> emailExport(@RequestBody EmailTimeEntriesRequest req) {
        mailService.emailExport(req, parsePropertyIdFilter(req.getPropertyId()));
        Map<String, Object> body = new LinkedHashMap<>();
        body.put("status", "sent");
        body.put("to", req.getTo());
        return body;
    }

    @GetMapping("/{id}")
    public TimeEntry get(@PathVariable int id) {
        return service.get(id);
    }

    @PostMapping({"", "/"})
    @ResponseStatus(HttpStatus.CREATED)
    public TimeEntry create(@RequestBody CreateTimeEntryRequest req) {
        return service.create(req);
    }

    @PutMapping("/{id}")
    public TimeEntry update(@PathVariable int id, @RequestBody UpdateTimeEntryRequest req) {
        return service.update(id, req);
    }

    @DeleteMapping("/{id}")
    public ResponseEntity<Void> delete(@PathVariable int id) {
        service.delete(id);
        return ResponseEntity.noContent().build();
    }

    /**
     * Matches the Go handler: {@code property_id=null} or {@code 0} means "General"
     * (DB {@code property_id IS NULL}). Omit the param for all properties.
     */
    static Integer parsePropertyIdFilter(String raw) {
        if (raw == null || raw.isBlank()) {
            return null;
        }
        String v = raw.trim();
        if ("null".equalsIgnoreCase(v) || "0".equals(v)) {
            return 0;
        }
        return Integer.parseInt(v);
    }
}
