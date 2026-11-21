package com.brickwise.web;

import com.brickwise.model.CreateTimeEntryRequest;
import com.brickwise.model.TimeEntry;
import com.brickwise.model.UpdateTimeEntryRequest;
import com.brickwise.service.TimeEntryService;
import org.springframework.http.HttpHeaders;
import org.springframework.http.HttpStatus;
import org.springframework.http.MediaType;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.nio.charset.StandardCharsets;
import java.util.List;

@RestController
@RequestMapping("/api/v1/time-entries")
public class TimeEntryController {
    private final TimeEntryService service;

    public TimeEntryController(TimeEntryService service) {
        this.service = service;
    }

    @GetMapping({"", "/"})
    public List<TimeEntry> list(
            @RequestParam(value = "property_id", required = false) Integer propertyId,
            @RequestParam(value = "start_date", required = false) String startDate,
            @RequestParam(value = "end_date", required = false) String endDate,
            @RequestParam(value = "category", required = false) String category) {
        return service.list(propertyId, startDate, endDate, category);
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

    @GetMapping("/export")
    public ResponseEntity<byte[]> export(
            @RequestParam(value = "property_id", required = false) Integer propertyId,
            @RequestParam(value = "start_date", required = false) String startDate,
            @RequestParam(value = "end_date", required = false) String endDate,
            @RequestParam(value = "category", required = false) String category) {
        List<TimeEntry> entries = service.list(propertyId, startDate, endDate, category);
        StringBuilder csv = new StringBuilder("Date,Activity Category,Description,Time Spent (hrs),Notes,Mileage\n");
        for (TimeEntry e : entries) {
            csv.append(e.getDate()).append(',')
                    .append(quote(e.getCategory())).append(',')
                    .append(quote(e.getDescription())).append(',')
                    .append(e.getHours()).append(',')
                    .append(quote(e.getNotes())).append(',')
                    .append(e.getMileage()).append('\n');
        }
        byte[] body = csv.toString().getBytes(StandardCharsets.UTF_8);
        return ResponseEntity.ok()
                .header(HttpHeaders.CONTENT_DISPOSITION, "attachment; filename=\"time-entries.csv\"")
                .contentType(MediaType.parseMediaType("text/csv"))
                .body(body);
    }

    private static String quote(String s) {
        if (s == null) {
            return "\"\"";
        }
        return "\"" + s.replace("\"", "\"\"") + "\"";
    }
}
