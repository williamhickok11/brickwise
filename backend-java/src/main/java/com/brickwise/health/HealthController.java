package com.brickwise.health;

import java.time.OffsetDateTime;
import java.time.format.DateTimeFormatter;
import java.time.temporal.ChronoUnit;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RestController;

// Same response as the Go /health route in backend/cmd/server/main.go.
@RestController
public class HealthController {

    @GetMapping("/health")
    public HealthStatus health() {
        String timestamp = OffsetDateTime.now()
                .truncatedTo(ChronoUnit.SECONDS)
                .format(DateTimeFormatter.ISO_OFFSET_DATE_TIME);
        return new HealthStatus("ok", timestamp);
    }

    public record HealthStatus(String status, String timestamp) {
    }
}
