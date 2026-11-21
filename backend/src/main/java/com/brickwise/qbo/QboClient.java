package com.brickwise.qbo;

import com.brickwise.config.BrickwiseProperties;
import com.brickwise.model.QboModels.QBOClass;
import com.brickwise.model.QboModels.QBOConnection;
import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.ObjectMapper;
import org.springframework.stereotype.Component;

import java.net.URI;
import java.net.URLEncoder;
import java.net.http.HttpClient;
import java.net.http.HttpRequest;
import java.net.http.HttpResponse;
import java.nio.charset.StandardCharsets;
import java.time.Duration;
import java.time.Instant;
import java.time.LocalDate;
import java.util.ArrayList;
import java.util.Base64;
import java.util.HashMap;
import java.util.Iterator;
import java.util.List;
import java.util.Map;

/**
 * Intuit OAuth + QuickBooks Accounting API client (hand-rolled HTTP, like the Go version).
 *
 * <p><b>C# analogy:</b> An {@code HttpClient}-based typed client registered as a singleton.
 * Java 11+ {@link HttpClient} is the standard-library equivalent of .NET's {@code HttpClient}.
 */
@Component
public class QboClient {
    private static final String OAUTH_AUTHORIZE = "https://appcenter.intuit.com/connect/oauth2";
    private static final String OAUTH_TOKEN = "https://oauth.platform.intuit.com/oauth2/v1/tokens/bearer";
    private static final String SCOPE = "com.intuit.quickbooks.accounting";

    private final BrickwiseProperties.Qbo cfg;
    private final HttpClient http = HttpClient.newBuilder().connectTimeout(Duration.ofSeconds(20)).build();
    private final ObjectMapper mapper = new ObjectMapper();

    public QboClient(BrickwiseProperties props) {
        this.cfg = props.getQbo();
    }

    public BrickwiseProperties.Qbo config() {
        return cfg;
    }

    public String authorizeUrl(String state) {
        return OAUTH_AUTHORIZE
                + "?client_id=" + enc(cfg.getClientId())
                + "&redirect_uri=" + enc(cfg.getRedirectUri())
                + "&response_type=code"
                + "&scope=" + enc(SCOPE)
                + "&state=" + enc(state);
    }

    public TokenResponse exchangeCode(String code) throws Exception {
        return postToken("grant_type=authorization_code&code=" + enc(code)
                + "&redirect_uri=" + enc(cfg.getRedirectUri()));
    }

    public TokenResponse refresh(String refreshToken) throws Exception {
        return postToken("grant_type=refresh_token&refresh_token=" + enc(refreshToken));
    }

    public List<QBOClass> listClasses(String realmId, String accessToken) throws Exception {
        String query = "SELECT Id, Name FROM Class MAXRESULTS 1000";
        String endpoint = apiBase() + "/v3/company/" + encPath(realmId)
                + "/query?query=" + enc(query) + "&minorversion=73";
        JsonNode root = getJson(endpoint, accessToken);
        List<QBOClass> out = new ArrayList<>();
        JsonNode classes = root.path("QueryResponse").path("Class");
        if (classes.isArray()) {
            for (JsonNode n : classes) {
                out.add(new QBOClass(n.path("Id").asText(), n.path("Name").asText()));
            }
        }
        return out;
    }

    public Map<String, Double> classNetIncome(String realmId, String accessToken, LocalDate start, LocalDate end)
            throws Exception {
        String endpoint = apiBase() + "/v3/company/" + encPath(realmId)
                + "/reports/ProfitAndLoss?start_date=" + start
                + "&end_date=" + end
                + "&accounting_method=Cash&summarize_column_by=Class&minorversion=73";
        JsonNode report = getJson(endpoint, accessToken);
        return parseClassNetIncome(report);
    }

    public static void applyTokens(QBOConnection conn, TokenResponse tokens, Instant now) {
        conn.setAccessToken(tokens.accessToken);
        conn.setRefreshToken(tokens.refreshToken);
        conn.setAccessTokenExpiresAt(now.plusSeconds(tokens.expiresIn));
        long refreshSecs = tokens.xRefreshTokenExpiresIn > 0 ? tokens.xRefreshTokenExpiresIn : 100L * 24 * 3600;
        conn.setRefreshTokenExpiresAt(now.plusSeconds(refreshSecs));
        conn.setUpdatedAt(now);
    }

    private TokenResponse postToken(String form) throws Exception {
        String basic = Base64.getEncoder().encodeToString(
                (cfg.getClientId() + ":" + cfg.getClientSecret()).getBytes(StandardCharsets.UTF_8));
        HttpRequest req = HttpRequest.newBuilder(URI.create(OAUTH_TOKEN))
                .timeout(Duration.ofSeconds(45))
                .header("Content-Type", "application/x-www-form-urlencoded")
                .header("Accept", "application/json")
                .header("Authorization", "Basic " + basic)
                .POST(HttpRequest.BodyPublishers.ofString(form))
                .build();
        HttpResponse<String> resp = http.send(req, HttpResponse.BodyHandlers.ofString());
        if (resp.statusCode() >= 300) {
            throw new IllegalStateException("token endpoint returned " + resp.statusCode());
        }
        JsonNode n = mapper.readTree(resp.body());
        TokenResponse t = new TokenResponse();
        t.accessToken = n.path("access_token").asText();
        t.refreshToken = n.path("refresh_token").asText();
        t.expiresIn = n.path("expires_in").asInt(3600);
        t.xRefreshTokenExpiresIn = n.path("x_refresh_token_expires_in").asInt(0);
        return t;
    }

    private JsonNode getJson(String endpoint, String accessToken) throws Exception {
        HttpRequest req = HttpRequest.newBuilder(URI.create(endpoint))
                .timeout(Duration.ofSeconds(45))
                .header("Accept", "application/json")
                .header("Authorization", "Bearer " + accessToken)
                .GET()
                .build();
        HttpResponse<String> resp = http.send(req, HttpResponse.BodyHandlers.ofString());
        if (resp.statusCode() >= 300) {
            throw new IllegalStateException("qbo api returned " + resp.statusCode());
        }
        return mapper.readTree(resp.body());
    }

    private String apiBase() {
        return "production".equalsIgnoreCase(cfg.getEnvironment())
                ? "https://quickbooks.api.intuit.com"
                : "https://sandbox-quickbooks.api.intuit.com";
    }

    private static Map<String, Double> parseClassNetIncome(JsonNode report) {
        List<String> classNames = new ArrayList<>();
        JsonNode cols = report.path("Columns").path("Column");
        if (cols.isArray()) {
            for (int i = 0; i < cols.size(); i++) {
                if (i == 0) continue;
                String title = cols.get(i).path("ColTitle").asText("").trim();
                if (title.isEmpty() || "Total".equalsIgnoreCase(title)) {
                    classNames.add("");
                } else {
                    classNames.add(title);
                }
            }
        }
        JsonNode netCols = findNetIncome(report.path("Rows").path("Row"));
        Map<String, Double> result = new HashMap<>();
        if (netCols == null || !netCols.isArray()) return result;
        for (int i = 0; i < classNames.size(); i++) {
            String name = classNames.get(i);
            if (name.isEmpty()) continue;
            int idx = i + 1;
            if (idx >= netCols.size()) break;
            String raw = netCols.get(idx).path("value").asText("").replace(",", "").trim();
            if (raw.isEmpty()) continue;
            try {
                result.put(name, Double.parseDouble(raw));
            } catch (NumberFormatException ignored) {
            }
        }
        return result;
    }

    private static JsonNode findNetIncome(JsonNode rows) {
        if (!rows.isArray()) return null;
        for (JsonNode row : rows) {
            if ("NetIncome".equalsIgnoreCase(row.path("group").asText())) {
                JsonNode summary = row.path("Summary").path("ColData");
                if (summary.isArray() && !summary.isEmpty()) return summary;
                JsonNode col = row.path("ColData");
                if (col.isArray() && !col.isEmpty()) return col;
            }
            JsonNode summary = row.path("Summary").path("ColData");
            if (summary.isArray() && !summary.isEmpty()
                    && "Net Income".equalsIgnoreCase(summary.get(0).path("value").asText("").trim())) {
                return summary;
            }
            JsonNode nested = findNetIncome(row.path("Rows").path("Row"));
            if (nested != null) return nested;
        }
        return null;
    }

    private static String enc(String s) {
        return URLEncoder.encode(s, StandardCharsets.UTF_8);
    }

    private static String encPath(String s) {
        return URLEncoder.encode(s, StandardCharsets.UTF_8).replace("+", "%20");
    }

    public static class TokenResponse {
        public String accessToken;
        public String refreshToken;
        public int expiresIn;
        public int xRefreshTokenExpiresIn;
    }
}
