package com.brickwise.model;

import java.time.Instant;
import java.time.LocalDate;
import java.util.ArrayList;
import java.util.List;

/** QuickBooks-related DTOs (connection, mappings, status, sync). */
public final class QboModels {
    private QboModels() {}

    public static class QBOConnection {
        private int id;
        private String realmId;
        private String accessToken;
        private String refreshToken;
        private Instant accessTokenExpiresAt;
        private Instant refreshTokenExpiresAt;
        private Instant connectedAt;
        private Instant updatedAt;

        public int getId() { return id; }
        public void setId(int id) { this.id = id; }
        public String getRealmId() { return realmId; }
        public void setRealmId(String realmId) { this.realmId = realmId; }
        public String getAccessToken() { return accessToken; }
        public void setAccessToken(String accessToken) { this.accessToken = accessToken; }
        public String getRefreshToken() { return refreshToken; }
        public void setRefreshToken(String refreshToken) { this.refreshToken = refreshToken; }
        public Instant getAccessTokenExpiresAt() { return accessTokenExpiresAt; }
        public void setAccessTokenExpiresAt(Instant accessTokenExpiresAt) { this.accessTokenExpiresAt = accessTokenExpiresAt; }
        public Instant getRefreshTokenExpiresAt() { return refreshTokenExpiresAt; }
        public void setRefreshTokenExpiresAt(Instant refreshTokenExpiresAt) { this.refreshTokenExpiresAt = refreshTokenExpiresAt; }
        public Instant getConnectedAt() { return connectedAt; }
        public void setConnectedAt(Instant connectedAt) { this.connectedAt = connectedAt; }
        public Instant getUpdatedAt() { return updatedAt; }
        public void setUpdatedAt(Instant updatedAt) { this.updatedAt = updatedAt; }
    }

    public static class PropertyQBOLink {
        private int propertyId;
        private String qboClassId;
        private String qboClassName;
        private Instant updatedAt;

        public int getPropertyId() { return propertyId; }
        public void setPropertyId(int propertyId) { this.propertyId = propertyId; }
        public String getQboClassId() { return qboClassId; }
        public void setQboClassId(String qboClassId) { this.qboClassId = qboClassId; }
        public String getQboClassName() { return qboClassName; }
        public void setQboClassName(String qboClassName) { this.qboClassName = qboClassName; }
        public Instant getUpdatedAt() { return updatedAt; }
        public void setUpdatedAt(Instant updatedAt) { this.updatedAt = updatedAt; }
    }

    public static class QBOCashFlowCache {
        private int propertyId;
        private LocalDate periodStart;
        private LocalDate periodEnd;
        private double netAmount;
        private Instant syncedAt;

        public int getPropertyId() { return propertyId; }
        public void setPropertyId(int propertyId) { this.propertyId = propertyId; }
        public LocalDate getPeriodStart() { return periodStart; }
        public void setPeriodStart(LocalDate periodStart) { this.periodStart = periodStart; }
        public LocalDate getPeriodEnd() { return periodEnd; }
        public void setPeriodEnd(LocalDate periodEnd) { this.periodEnd = periodEnd; }
        public double getNetAmount() { return netAmount; }
        public void setNetAmount(double netAmount) { this.netAmount = netAmount; }
        public Instant getSyncedAt() { return syncedAt; }
        public void setSyncedAt(Instant syncedAt) { this.syncedAt = syncedAt; }
    }

    public static class QBOClass {
        private String id;
        private String name;
        public QBOClass() {}
        public QBOClass(String id, String name) { this.id = id; this.name = name; }
        public String getId() { return id; }
        public void setId(String id) { this.id = id; }
        public String getName() { return name; }
        public void setName(String name) { this.name = name; }
    }

    public static class QBOStatus {
        private boolean configured;
        private boolean connected;
        private String realmId;
        private Instant connectedAt;
        private Instant lastSyncedAt;
        private Instant periodStart;
        private Instant periodEnd;
        private List<PropertyQBOLink> mappings = new ArrayList<>();
        private boolean classRequired;
        private String message;

        public boolean isConfigured() { return configured; }
        public void setConfigured(boolean configured) { this.configured = configured; }
        public boolean isConnected() { return connected; }
        public void setConnected(boolean connected) { this.connected = connected; }
        public String getRealmId() { return realmId; }
        public void setRealmId(String realmId) { this.realmId = realmId; }
        public Instant getConnectedAt() { return connectedAt; }
        public void setConnectedAt(Instant connectedAt) { this.connectedAt = connectedAt; }
        public Instant getLastSyncedAt() { return lastSyncedAt; }
        public void setLastSyncedAt(Instant lastSyncedAt) { this.lastSyncedAt = lastSyncedAt; }
        public Instant getPeriodStart() { return periodStart; }
        public void setPeriodStart(Instant periodStart) { this.periodStart = periodStart; }
        public Instant getPeriodEnd() { return periodEnd; }
        public void setPeriodEnd(Instant periodEnd) { this.periodEnd = periodEnd; }
        public List<PropertyQBOLink> getMappings() { return mappings; }
        public void setMappings(List<PropertyQBOLink> mappings) { this.mappings = mappings; }
        public boolean isClassRequired() { return classRequired; }
        public void setClassRequired(boolean classRequired) { this.classRequired = classRequired; }
        public String getMessage() { return message; }
        public void setMessage(String message) { this.message = message; }
    }

    public static class UpdateQBOMappingsRequest {
        private List<PropertyQBOMappingInput> mappings = new ArrayList<>();
        public List<PropertyQBOMappingInput> getMappings() { return mappings; }
        public void setMappings(List<PropertyQBOMappingInput> mappings) { this.mappings = mappings; }
    }

    public static class PropertyQBOMappingInput {
        private int propertyId;
        private String qboClassId;
        private String qboClassName;
        public int getPropertyId() { return propertyId; }
        public void setPropertyId(int propertyId) { this.propertyId = propertyId; }
        public String getQboClassId() { return qboClassId; }
        public void setQboClassId(String qboClassId) { this.qboClassId = qboClassId; }
        public String getQboClassName() { return qboClassName; }
        public void setQboClassName(String qboClassName) { this.qboClassName = qboClassName; }
    }

    public static class QBOSyncResult {
        private Instant syncedAt;
        private Instant periodStart;
        private Instant periodEnd;
        private int updated;
        public Instant getSyncedAt() { return syncedAt; }
        public void setSyncedAt(Instant syncedAt) { this.syncedAt = syncedAt; }
        public Instant getPeriodStart() { return periodStart; }
        public void setPeriodStart(Instant periodStart) { this.periodStart = periodStart; }
        public Instant getPeriodEnd() { return periodEnd; }
        public void setPeriodEnd(Instant periodEnd) { this.periodEnd = periodEnd; }
        public int getUpdated() { return updated; }
        public void setUpdated(int updated) { this.updated = updated; }
    }
}
