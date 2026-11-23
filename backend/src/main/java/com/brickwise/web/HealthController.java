package com.brickwise.web;

import org.springframework.beans.factory.annotation.Qualifier;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RestController;

import java.time.Instant;
import java.util.LinkedHashMap;
import java.util.Map;

@RestController
public class HealthController {
    private final boolean usingPostgres;

    public HealthController(@Qualifier("usingPostgres") boolean usingPostgres) {
        this.usingPostgres = usingPostgres;
    }

    @GetMapping("/health")
    public Map<String, String> health() {
        Map<String, String> body = new LinkedHashMap<>();
        body.put("status", "ok");
        body.put("timestamp", Instant.now().toString());
        // postgres = DATABASE_URL set (Supabase/Postgres); sqlite = local file fallback
        body.put("database", usingPostgres ? "postgres" : "sqlite");
        return body;
    }
}
