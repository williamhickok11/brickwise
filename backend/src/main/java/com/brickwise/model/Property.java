package com.brickwise.model;

import java.time.Instant;
import java.time.LocalDateTime;

/**
 * A managed real-estate property.
 *
 * <p><b>C# analogy:</b> A plain DTO / entity class. Java {@code record} is similar to a C#
 * {@code record}, but here we use a mutable class so JDBC can populate fields after insert —
 * closer to a typical EF entity or Dapper POCO.
 */
public class Property {
    private int id;
    private String name;
    private String address;
    private String propertyType;
    private double defaultMileage;
    private double currentValue;
    private Instant valueUpdatedAt;
    private Integer mainImageId;
    private Instant createdAt;
    private Instant updatedAt;

    public int getId() { return id; }
    public void setId(int id) { this.id = id; }
    public String getName() { return name; }
    public void setName(String name) { this.name = name; }
    public String getAddress() { return address; }
    public void setAddress(String address) { this.address = address; }
    public String getPropertyType() { return propertyType; }
    public void setPropertyType(String propertyType) { this.propertyType = propertyType; }
    public double getDefaultMileage() { return defaultMileage; }
    public void setDefaultMileage(double defaultMileage) { this.defaultMileage = defaultMileage; }
    public double getCurrentValue() { return currentValue; }
    public void setCurrentValue(double currentValue) { this.currentValue = currentValue; }
    public Instant getValueUpdatedAt() { return valueUpdatedAt; }
    public void setValueUpdatedAt(Instant valueUpdatedAt) { this.valueUpdatedAt = valueUpdatedAt; }
    public Integer getMainImageId() { return mainImageId; }
    public void setMainImageId(Integer mainImageId) { this.mainImageId = mainImageId; }
    public Instant getCreatedAt() { return createdAt; }
    public void setCreatedAt(Instant createdAt) { this.createdAt = createdAt; }
    public Instant getUpdatedAt() { return updatedAt; }
    public void setUpdatedAt(Instant updatedAt) { this.updatedAt = updatedAt; }
}
