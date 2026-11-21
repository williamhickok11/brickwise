package com.brickwise.model;

import java.time.Instant;

/** Metadata for a lease PDF attached to a residence (one per residence). */
public class Document {
    private int id;
    private int residenceId;
    private String fileName;
    private String contentType;
    private long size;
    private Instant createdAt;

    public int getId() { return id; }
    public void setId(int id) { this.id = id; }
    public int getResidenceId() { return residenceId; }
    public void setResidenceId(int residenceId) { this.residenceId = residenceId; }
    public String getFileName() { return fileName; }
    public void setFileName(String fileName) { this.fileName = fileName; }
    public String getContentType() { return contentType; }
    public void setContentType(String contentType) { this.contentType = contentType; }
    public long getSize() { return size; }
    public void setSize(long size) { this.size = size; }
    public Instant getCreatedAt() { return createdAt; }
    public void setCreatedAt(Instant createdAt) { this.createdAt = createdAt; }
}
