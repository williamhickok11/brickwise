package com.brickwise.model;

import com.fasterxml.jackson.annotation.JsonProperty;

public class CreateTimeEntryRequest {
    private Integer propertyId;
    private String date;
    private String category;
    private String description;
    private double hours;
    private String notes;
    private double mileage;
    @JsonProperty("full_drive")
    private boolean fullDrive;

    public Integer getPropertyId() { return propertyId; }
    public void setPropertyId(Integer propertyId) { this.propertyId = propertyId; }
    public String getDate() { return date; }
    public void setDate(String date) { this.date = date; }
    public String getCategory() { return category; }
    public void setCategory(String category) { this.category = category; }
    public String getDescription() { return description; }
    public void setDescription(String description) { this.description = description; }
    public double getHours() { return hours; }
    public void setHours(double hours) { this.hours = hours; }
    public String getNotes() { return notes; }
    public void setNotes(String notes) { this.notes = notes; }
    public double getMileage() { return mileage; }
    public void setMileage(double mileage) { this.mileage = mileage; }
    public boolean isFullDrive() { return fullDrive; }
    public void setFullDrive(boolean fullDrive) { this.fullDrive = fullDrive; }
}
