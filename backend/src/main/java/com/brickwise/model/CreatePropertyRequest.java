package com.brickwise.model;

public class CreatePropertyRequest {
    private String name;
    private String address;
    private String propertyType;
    private double defaultMileage;

    public String getName() { return name; }
    public void setName(String name) { this.name = name; }
    public String getAddress() { return address; }
    public void setAddress(String address) { this.address = address; }
    public String getPropertyType() { return propertyType; }
    public void setPropertyType(String propertyType) { this.propertyType = propertyType; }
    public double getDefaultMileage() { return defaultMileage; }
    public void setDefaultMileage(double defaultMileage) { this.defaultMileage = defaultMileage; }
}
