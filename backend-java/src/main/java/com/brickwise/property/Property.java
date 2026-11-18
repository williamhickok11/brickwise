package com.brickwise.property;

import java.time.OffsetDateTime;

// A record is a C# record: immutable data, constructor, and accessors generated for you.
// spring.jackson.property-naming-strategy turns propertyType into the JSON field property_type.
public record Property(
        int id,
        String name,
        String address,
        String propertyType,
        double defaultMileage,
        double currentValue,
        OffsetDateTime valueUpdatedAt,
        OffsetDateTime createdAt,
        OffsetDateTime updatedAt
) {
}
