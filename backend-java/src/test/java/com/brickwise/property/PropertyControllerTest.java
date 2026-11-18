package com.brickwise.property;

import static org.hamcrest.Matchers.nullValue;
import static org.mockito.Mockito.when;
import static org.springframework.test.web.servlet.request.MockMvcRequestBuilders.get;
import static org.springframework.test.web.servlet.result.MockMvcResultMatchers.jsonPath;
import static org.springframework.test.web.servlet.result.MockMvcResultMatchers.status;

import java.time.OffsetDateTime;
import java.time.ZoneOffset;
import java.util.List;
import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.autoconfigure.web.servlet.WebMvcTest;
import org.springframework.test.context.bean.override.mockito.MockitoBean;
import org.springframework.test.web.servlet.MockMvc;

// Controller test with the service replaced. Same idea as backend/internal/handlers/property_test.go.
@WebMvcTest(PropertyController.class)
class PropertyControllerTest {

    @Autowired
    private MockMvc mockMvc;

    @MockitoBean
    private PropertyService propertyService;

    @Test
    void listPropertiesReturnsSnakeCaseJson() throws Exception {
        OffsetDateTime createdAt = OffsetDateTime.of(2026, 3, 30, 20, 17, 18, 0, ZoneOffset.UTC);
        when(propertyService.listProperties()).thenReturn(List.of(
                new Property(1, "Oak House", "1 Oak St", "single_family", 12.5, 250000, null, createdAt, createdAt)
        ));

        mockMvc.perform(get("/api/v1/properties"))
                .andExpect(status().isOk())
                .andExpect(jsonPath("$[0].id").value(1))
                .andExpect(jsonPath("$[0].name").value("Oak House"))
                .andExpect(jsonPath("$[0].address").value("1 Oak St"))
                .andExpect(jsonPath("$[0].property_type").value("single_family"))
                .andExpect(jsonPath("$[0].default_mileage").value(12.5))
                .andExpect(jsonPath("$[0].current_value").value(250000))
                .andExpect(jsonPath("$[0].value_updated_at").value(nullValue()))
                .andExpect(jsonPath("$[0].created_at").value("2026-03-30T20:17:18Z"))
                .andExpect(jsonPath("$[0].updated_at").value("2026-03-30T20:17:18Z"));
    }
}
