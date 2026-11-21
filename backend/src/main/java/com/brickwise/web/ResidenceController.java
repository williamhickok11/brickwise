package com.brickwise.web;

import com.brickwise.model.CreateResidenceRequest;
import com.brickwise.model.Residence;
import com.brickwise.model.UpdateResidenceRequest;
import com.brickwise.service.ResidenceService;
import org.springframework.http.HttpStatus;
import org.springframework.web.bind.annotation.*;

import java.util.List;

@RestController
@RequestMapping("/api/v1/residences")
public class ResidenceController {
    private final ResidenceService service;

    public ResidenceController(ResidenceService service) {
        this.service = service;
    }

    @GetMapping({"", "/"})
    public List<Residence> list(
            @RequestParam(value = "property_id", required = false) Integer propertyId,
            @RequestParam(value = "is_active", required = false) Boolean isActive) {
        return service.list(propertyId, isActive);
    }

    @GetMapping("/{id}")
    public Residence get(@PathVariable int id) {
        return service.get(id);
    }

    @PostMapping({"", "/"})
    @ResponseStatus(HttpStatus.CREATED)
    public Residence create(@RequestBody CreateResidenceRequest req) {
        return service.create(req);
    }

    @PutMapping("/{id}")
    public Residence update(@PathVariable int id, @RequestBody UpdateResidenceRequest req) {
        return service.update(id, req);
    }
}
