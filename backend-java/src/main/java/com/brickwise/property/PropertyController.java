package com.brickwise.property;

import java.util.List;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

// [ApiController] / Go PropertyHandler in backend/internal/handlers/property.go.
@RestController
@RequestMapping("/api/v1/properties")
public class PropertyController {

    private final PropertyService propertyService;

    public PropertyController(PropertyService propertyService) {
        this.propertyService = propertyService;
    }

    // [HttpGet] / Go PropertyHandler.ListProperties.
    @GetMapping
    public List<Property> listProperties() {
        return propertyService.listProperties();
    }
}
