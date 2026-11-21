package com.brickwise.web;

import com.brickwise.model.Mortgage;
import com.brickwise.model.UpsertMortgageRequest;
import com.brickwise.service.MortgageService;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

@RestController
@RequestMapping("/api/v1/properties/{propertyId}/mortgage")
public class MortgageController {
    private final MortgageService service;

    public MortgageController(MortgageService service) {
        this.service = service;
    }

    @GetMapping
    public Mortgage get(@PathVariable int propertyId) {
        return service.get(propertyId);
    }

    @PutMapping
    public Mortgage upsert(@PathVariable int propertyId, @RequestBody UpsertMortgageRequest req) {
        return service.upsert(propertyId, req);
    }

    @DeleteMapping
    public ResponseEntity<Void> delete(@PathVariable int propertyId) {
        service.delete(propertyId);
        return ResponseEntity.noContent().build();
    }
}
