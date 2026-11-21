package com.brickwise.model;

import com.fasterxml.jackson.annotation.JsonProperty;

import java.time.Instant;

public class PropertyImage {
    private int id;
    private int propertyId;
    private String fileName;
    private String contentType;
    private long size;
    @JsonProperty("is_main")
    private boolean isMain;
    private Instant createdAt;

    public int getId() { return id; }
    public void setId(int id) { this.id = id; }
    public int getPropertyId() { return propertyId; }
    public void setPropertyId(int propertyId) { this.propertyId = propertyId; }
    public String getFileName() { return fileName; }
    public void setFileName(String fileName) { this.fileName = fileName; }
    public String getContentType() { return contentType; }
    public void setContentType(String contentType) { this.contentType = contentType; }
    public long getSize() { return size; }
    public void setSize(long size) { this.size = size; }
    public boolean isMain() { return isMain; }
    public void setMain(boolean main) { isMain = main; }
    public Instant getCreatedAt() { return createdAt; }
    public void setCreatedAt(Instant createdAt) { this.createdAt = createdAt; }
}
