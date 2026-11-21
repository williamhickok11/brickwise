package com.brickwise.config;

import org.springframework.boot.context.properties.ConfigurationProperties;

/**
 * Typed settings bound from {@code application.yml} / environment variables.
 *
 * <p><b>C# analogy:</b> Like the Options pattern ({@code IOptions<T>}) with
 * {@code builder.Services.Configure<T>(configuration.GetSection(...))}.
 */
@ConfigurationProperties(prefix = "brickwise")
public class BrickwiseProperties {
    private String dbPath = "./brickwise.db";
    private String documentsDir = "./data/documents";
    private String propertyImagesDir = "./data/property-images";
    private final Qbo qbo = new Qbo();

    public String getDbPath() { return dbPath; }
    public void setDbPath(String dbPath) { this.dbPath = dbPath; }
    public String getDocumentsDir() { return documentsDir; }
    public void setDocumentsDir(String documentsDir) { this.documentsDir = documentsDir; }
    public String getPropertyImagesDir() { return propertyImagesDir; }
    public void setPropertyImagesDir(String propertyImagesDir) { this.propertyImagesDir = propertyImagesDir; }
    public Qbo getQbo() { return qbo; }

    public static class Qbo {
        private String clientId = "";
        private String clientSecret = "";
        private String redirectUri = "http://localhost:8080/api/v1/integrations/qbo/callback";
        private String frontendRedirect = "http://localhost:5173/settings/quickbooks";
        private String environment = "sandbox";

        public String getClientId() { return clientId; }
        public void setClientId(String clientId) { this.clientId = clientId; }
        public String getClientSecret() { return clientSecret; }
        public void setClientSecret(String clientSecret) { this.clientSecret = clientSecret; }
        public String getRedirectUri() { return redirectUri; }
        public void setRedirectUri(String redirectUri) { this.redirectUri = redirectUri; }
        public String getFrontendRedirect() { return frontendRedirect; }
        public void setFrontendRedirect(String frontendRedirect) { this.frontendRedirect = frontendRedirect; }
        public String getEnvironment() { return environment; }
        public void setEnvironment(String environment) { this.environment = environment; }

        public boolean isConfigured() {
            return clientId != null && !clientId.isBlank() && clientSecret != null && !clientSecret.isBlank();
        }
    }
}
