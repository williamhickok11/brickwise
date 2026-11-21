package com.brickwise.model;

public class UpdatePropertyValueRequest {
    private double currentValue;
    private String valueUpdatedAt;

    public double getCurrentValue() { return currentValue; }
    public void setCurrentValue(double currentValue) { this.currentValue = currentValue; }
    public String getValueUpdatedAt() { return valueUpdatedAt; }
    public void setValueUpdatedAt(String valueUpdatedAt) { this.valueUpdatedAt = valueUpdatedAt; }
}
