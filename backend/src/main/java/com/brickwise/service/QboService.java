package com.brickwise.service;

import com.brickwise.error.ApiException;
import com.brickwise.model.Property;
import com.brickwise.model.QboModels.PropertyQBOLink;
import com.brickwise.model.QboModels.PropertyQBOMappingInput;
import com.brickwise.model.QboModels.QBOCashFlowCache;
import com.brickwise.model.QboModels.QBOClass;
import com.brickwise.model.QboModels.QBOConnection;
import com.brickwise.model.QboModels.QBOStatus;
import com.brickwise.model.QboModels.QBOSyncResult;
import com.brickwise.model.QboModels.UpdateQBOMappingsRequest;
import com.brickwise.qbo.QboClient;
import com.brickwise.repository.PropertyRepository;
import com.brickwise.repository.QboRepository;
import org.springframework.stereotype.Service;

import java.net.URLEncoder;
import java.nio.charset.StandardCharsets;
import java.security.SecureRandom;
import java.time.Instant;
import java.time.LocalDate;
import java.time.ZoneOffset;
import java.time.temporal.TemporalAdjusters;
import java.util.ArrayList;
import java.util.HashMap;
import java.util.HexFormat;
import java.util.List;
import java.util.Map;

@Service
public class QboService {
    private final QboRepository repo;
    private final PropertyRepository propertyRepo;
    private final QboClient client;

    public QboService(QboRepository repo, PropertyRepository propertyRepo, QboClient client) {
        this.repo = repo;
        this.propertyRepo = propertyRepo;
        this.client = client;
    }

    public QBOStatus status() {
        QBOStatus status = new QBOStatus();
        status.setConfigured(client.config().isConfigured());
        status.setClassRequired(true);
        status.setMessage("Map each property to a QuickBooks Class. Your accounting team must code income and expenses to that Class.");
        if (!status.isConfigured()) {
            status.setMessage("Set QBO_CLIENT_ID and QBO_CLIENT_SECRET to enable QuickBooks Online.");
            return status;
        }
        repo.findConnection().ifPresent(conn -> {
            status.setConnected(true);
            status.setRealmId(conn.getRealmId());
            status.setConnectedAt(conn.getConnectedAt());
        });
        status.setMappings(repo.listLinks());
        repo.latestSyncedAt().ifPresent(status::setLastSyncedAt);
        LocalDate[] bounds = currentMonthBounds();
        status.setPeriodStart(bounds[0].atStartOfDay().toInstant(ZoneOffset.UTC));
        status.setPeriodEnd(bounds[1].atStartOfDay().toInstant(ZoneOffset.UTC));
        return status;
    }

    public String startConnect() {
        if (!client.config().isConfigured()) {
            throw new ApiException("QBO_NOT_CONFIGURED", "QuickBooks credentials are not configured", 400);
        }
        byte[] bytes = new byte[24];
        new SecureRandom().nextBytes(bytes);
        String state = HexFormat.of().formatHex(bytes);
        repo.saveOauthState(state, Instant.now().plusSeconds(15 * 60));
        return client.authorizeUrl(state);
    }

    public void handleCallback(String code, String state, String realmId) {
        if (!client.config().isConfigured()) {
            throw new ApiException("QBO_NOT_CONFIGURED", "QuickBooks credentials are not configured", 400);
        }
        if (code == null || code.isBlank() || state == null || state.isBlank() || realmId == null || realmId.isBlank()) {
            throw ApiException.invalidInput("Missing OAuth callback parameters");
        }
        if (!repo.consumeOauthState(state)) {
            throw new ApiException("INVALID_STATE", "Invalid or expired OAuth state", 400);
        }
        try {
            var tokens = client.exchangeCode(code);
            Instant now = Instant.now();
            QBOConnection conn = new QBOConnection();
            conn.setRealmId(realmId);
            conn.setConnectedAt(now);
            QboClient.applyTokens(conn, tokens, now);
            repo.saveConnection(conn);
        } catch (Exception e) {
            throw ApiException.internal("OAuth token exchange failed: " + e.getMessage());
        }
    }

    public String frontendRedirect(boolean success, String errMsg) {
        String base = client.config().getFrontendRedirect();
        String sep = base.contains("?") ? "&" : "?";
        if (success) return base + sep + "connected=1";
        return base + sep + "error=" + URLEncoder.encode(errMsg == null ? "error" : errMsg, StandardCharsets.UTF_8);
    }

    public void disconnect() {
        repo.deleteConnection();
        repo.replaceLinks(List.of());
    }

    public List<QBOClass> listClasses() {
        QBOConnection conn = ensureFreshToken();
        try {
            return client.listClasses(conn.getRealmId(), conn.getAccessToken());
        } catch (Exception e) {
            throw ApiException.internal("Failed to list QBO classes: " + e.getMessage());
        }
    }

    public void updateMappings(UpdateQBOMappingsRequest req) {
        if (req == null) throw ApiException.invalidInput("Invalid input provided");
        Map<Integer, Boolean> valid = new HashMap<>();
        for (Property p : propertyRepo.findAll()) valid.put(p.getId(), true);
        List<PropertyQBOLink> keep = new ArrayList<>();
        for (PropertyQBOMappingInput m : req.getMappings()) {
            if (!valid.containsKey(m.getPropertyId())) {
                throw ApiException.invalidInput("Unknown property_id " + m.getPropertyId());
            }
            if (m.getQboClassId() == null || m.getQboClassId().isBlank()) continue;
            PropertyQBOLink link = new PropertyQBOLink();
            link.setPropertyId(m.getPropertyId());
            link.setQboClassId(m.getQboClassId().trim());
            link.setQboClassName(m.getQboClassName() == null ? "" : m.getQboClassName());
            keep.add(link);
        }
        repo.replaceLinks(keep);
    }

    public QBOSyncResult sync() {
        QBOConnection conn = ensureFreshToken();
        List<PropertyQBOLink> links = repo.listLinks();
        LocalDate[] bounds = currentMonthBounds();
        try {
            Map<String, Double> byClass = client.classNetIncome(
                    conn.getRealmId(), conn.getAccessToken(), bounds[0], bounds[1]);
            Instant now = Instant.now();
            int updated = 0;
            for (PropertyQBOLink link : links) {
                Double amount = byClass.get(link.getQboClassName());
                if (amount == null) amount = 0.0;
                repo.upsertCashFlow(link.getPropertyId(), bounds[0], bounds[1], amount, now);
                updated++;
            }
            QBOSyncResult result = new QBOSyncResult();
            result.setSyncedAt(now);
            result.setPeriodStart(bounds[0].atStartOfDay().toInstant(ZoneOffset.UTC));
            result.setPeriodEnd(bounds[1].atStartOfDay().toInstant(ZoneOffset.UTC));
            result.setUpdated(updated);
            return result;
        } catch (Exception e) {
            throw ApiException.internal("QBO sync failed: " + e.getMessage());
        }
    }

    public Map<Integer, Double> currentMonthCashFlow() {
        Map<Integer, Double> out = new HashMap<>();
        for (QBOCashFlowCache c : repo.listCashFlowCache()) {
            out.put(c.getPropertyId(), c.getNetAmount());
        }
        return out;
    }

    private QBOConnection ensureFreshToken() {
        QBOConnection conn = repo.findConnection().orElseThrow(() ->
                new ApiException("QBO_NOT_CONNECTED", "QuickBooks is not connected", 400));
        if (conn.getAccessTokenExpiresAt().isAfter(Instant.now().plusSeconds(120))) {
            return conn;
        }
        try {
            var tokens = client.refresh(conn.getRefreshToken());
            Instant now = Instant.now();
            QboClient.applyTokens(conn, tokens, now);
            repo.saveConnection(conn);
            return conn;
        } catch (Exception e) {
            throw ApiException.internal("Failed to refresh QBO token: " + e.getMessage());
        }
    }

    private static LocalDate[] currentMonthBounds() {
        LocalDate now = LocalDate.now();
        return new LocalDate[]{now.withDayOfMonth(1), now.with(TemporalAdjusters.lastDayOfMonth())};
    }
}
