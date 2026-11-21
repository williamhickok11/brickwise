package com.brickwise.service;

import com.brickwise.error.ApiException;
import com.brickwise.model.CreatePropertyRequest;
import com.brickwise.model.Property;
import com.brickwise.repository.PropertyRepository;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.InjectMocks;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;

import java.util.Optional;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertThrows;
import static org.mockito.ArgumentMatchers.any;
import static org.mockito.Mockito.when;

/**
 * Service-level tests (Mockito).
 *
 * <p><b>C# analogy:</b> xUnit/NUnit + Moq. {@code @ExtendWith(MockitoExtension.class)} is like
 * adding a Moq fixture; {@code @Mock}/{@code @InjectMocks} mirror mocked dependencies.
 */
@ExtendWith(MockitoExtension.class)
class PropertyServiceTest {

    @Mock
    PropertyRepository repo;

    @InjectMocks
    PropertyService service;

    @Test
    void createValidProperty() {
        CreatePropertyRequest req = new CreatePropertyRequest();
        req.setName("Wood Bridge");
        req.setAddress("123 Main");
        req.setPropertyType("residential");
        req.setDefaultMileage(22);

        when(repo.create(any(Property.class))).thenAnswer(inv -> {
            Property p = inv.getArgument(0);
            p.setId(1);
            return p;
        });

        Property created = service.create(req);
        assertEquals(1, created.getId());
        assertEquals("Wood Bridge", created.getName());
    }

    @Test
    void getMissingThrowsNotFound() {
        when(repo.findById(99)).thenReturn(Optional.empty());
        ApiException ex = assertThrows(ApiException.class, () -> service.get(99));
        assertEquals("NOT_FOUND", ex.getCode());
    }
}
