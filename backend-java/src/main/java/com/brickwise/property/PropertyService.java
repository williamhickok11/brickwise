package com.brickwise.property;

import java.util.List;
import org.springframework.stereotype.Service;

// Business logic lives here. Same role as PropertyService in backend/internal/service/property.go.
@Service
public class PropertyService {

    private final PropertyRepository properties;

    public PropertyService(PropertyRepository properties) {
        this.properties = properties;
    }

    public List<Property> listProperties() {
        return properties.findAll();
    }
}
