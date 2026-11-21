package com.brickwise.service;

import com.brickwise.error.ApiException;
import com.brickwise.model.CreatePropertyRequest;
import com.brickwise.model.Property;
import com.brickwise.model.UpdatePropertyValueRequest;
import com.brickwise.repository.PropertyRepository;
import com.brickwise.validation.DomainValidator;
import org.springframework.stereotype.Service;

import java.time.Instant;
import java.time.LocalDate;
import java.time.ZoneOffset;
import java.time.format.DateTimeParseException;
import java.util.List;

/**
 * Business logic for properties.
 *
 * <p><b>C# analogy:</b> An application service / domain service registered in DI with
 * {@code services.AddScoped<PropertyService>()}. {@code @Service} is Spring's stereotype for that.
 */
@Service
public class PropertyService {
    private final PropertyRepository repo;

    public PropertyService(PropertyRepository repo) {
        this.repo = repo;
    }

    public Property create(CreatePropertyRequest req) {
        DomainValidator.validateProperty(req);
        Property p = new Property();
        p.setName(req.getName());
        p.setAddress(req.getAddress());
        p.setPropertyType(req.getPropertyType());
        p.setDefaultMileage(req.getDefaultMileage());
        return repo.create(p);
    }

    public Property get(int id) {
        return repo.findById(id).orElseThrow(ApiException::notFound);
    }

    public List<Property> list() {
        return repo.findAll();
    }

    public Property update(int id, CreatePropertyRequest req) {
        DomainValidator.validateProperty(req);
        get(id);
        Property p = new Property();
        p.setId(id);
        p.setName(req.getName());
        p.setAddress(req.getAddress());
        p.setPropertyType(req.getPropertyType());
        p.setDefaultMileage(req.getDefaultMileage());
        try {
            repo.update(p);
        } catch (IllegalStateException e) {
            throw ApiException.notFound();
        }
        return get(id);
    }

    public void delete(int id) {
        try {
            repo.delete(id);
        } catch (IllegalStateException e) {
            throw ApiException.notFound();
        }
    }

    public Property updateValue(int id, UpdatePropertyValueRequest req) {
        if (req.getCurrentValue() < 0) {
            throw ApiException.invalidInput("current_value must be non-negative");
        }
        Instant when;
        if (req.getValueUpdatedAt() != null && !req.getValueUpdatedAt().isBlank()) {
            try {
                when = LocalDate.parse(req.getValueUpdatedAt()).atStartOfDay().toInstant(ZoneOffset.UTC);
            } catch (DateTimeParseException e) {
                throw new ApiException("INVALID_DATE", "value_updated_at must be in YYYY-MM-DD format", 400);
            }
        } else {
            when = Instant.now();
        }
        try {
            repo.updateValue(id, req.getCurrentValue(), when);
        } catch (IllegalStateException e) {
            throw ApiException.notFound();
        }
        return get(id);
    }
}
