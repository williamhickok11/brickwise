package com.brickwise.web;

import com.brickwise.model.CreatePropertyRequest;
import com.brickwise.model.Property;
import com.brickwise.model.UpdatePropertyValueRequest;
import com.brickwise.service.PropertyImageService;
import com.brickwise.service.PropertyService;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.util.List;

/**
 * Property HTTP API under {@code /api/v1/properties}.
 *
 * <p><b>C# analogy:</b> An ASP.NET {@code [ApiController]} / {@code ControllerBase} with
 * attribute routing. {@code @RestController} = controller + JSON serialization by default.
 * Constructor injection replaces {@code IPropertyService} constructor DI.
 */
@RestController
@RequestMapping("/api/v1/properties")
public class PropertyController {
    private final PropertyService service;
    private final PropertyImageService imageService;

    public PropertyController(PropertyService service, PropertyImageService imageService) {
        this.service = service;
        this.imageService = imageService;
    }

    @GetMapping({"", "/"})
    public List<Property> list() {
        return service.list();
    }

    @GetMapping("/{id}")
    public Property get(@PathVariable int id) {
        return service.get(id);
    }

    @PostMapping({"", "/"})
    @ResponseStatus(HttpStatus.CREATED)
    public Property create(@RequestBody CreatePropertyRequest req) {
        return service.create(req);
    }

    @PutMapping("/{id}")
    public Property update(@PathVariable int id, @RequestBody CreatePropertyRequest req) {
        return service.update(id, req);
    }

    @PutMapping("/{id}/value")
    public Property updateValue(@PathVariable int id, @RequestBody UpdatePropertyValueRequest req) {
        return service.updateValue(id, req);
    }

    @DeleteMapping("/{id}")
    public ResponseEntity<Void> delete(@PathVariable int id) {
        service.delete(id);
        imageService.cleanupForProperty(id);
        return ResponseEntity.noContent().build();
    }
}
