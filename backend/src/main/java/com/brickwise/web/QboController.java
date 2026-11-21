package com.brickwise.web;

import com.brickwise.model.QboModels.QBOClass;
import com.brickwise.model.QboModels.QBOStatus;
import com.brickwise.model.QboModels.QBOSyncResult;
import com.brickwise.model.QboModels.UpdateQBOMappingsRequest;
import com.brickwise.service.QboService;
import jakarta.servlet.http.HttpServletResponse;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.io.IOException;
import java.util.List;
import java.util.Map;

@RestController
@RequestMapping("/api/v1/integrations/qbo")
public class QboController {
    private final QboService service;

    public QboController(QboService service) {
        this.service = service;
    }

    @GetMapping("/status")
    public QBOStatus status() {
        return service.status();
    }

    @GetMapping("/connect")
    public Map<String, String> connect() {
        return Map.of("authorize_url", service.startConnect());
    }

    @GetMapping("/callback")
    public void callback(
            @RequestParam(value = "code", required = false) String code,
            @RequestParam(value = "state", required = false) String state,
            @RequestParam(value = "realmId", required = false) String realmId,
            @RequestParam(value = "realm_id", required = false) String realmIdAlt,
            @RequestParam(value = "error", required = false) String error,
            HttpServletResponse response) throws IOException {
        if (error != null && !error.isBlank()) {
            response.sendRedirect(service.frontendRedirect(false, error));
            return;
        }
        String realm = realmId != null && !realmId.isBlank() ? realmId : realmIdAlt;
        try {
            service.handleCallback(code, state, realm);
            response.sendRedirect(service.frontendRedirect(true, ""));
        } catch (Exception e) {
            response.sendRedirect(service.frontendRedirect(false, "oauth_failed"));
        }
    }

    @DeleteMapping("/disconnect")
    public ResponseEntity<Void> disconnect() {
        service.disconnect();
        return ResponseEntity.noContent().build();
    }

    @GetMapping("/classes")
    public List<QBOClass> classes() {
        return service.listClasses();
    }

    @PutMapping("/mappings")
    public QBOStatus mappings(@RequestBody UpdateQBOMappingsRequest req) {
        service.updateMappings(req);
        return service.status();
    }

    @PostMapping("/sync")
    public QBOSyncResult sync() {
        return service.sync();
    }
}
