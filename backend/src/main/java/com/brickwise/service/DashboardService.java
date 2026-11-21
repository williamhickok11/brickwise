package com.brickwise.service;

import com.brickwise.model.DashboardSummary;
import com.brickwise.model.Mortgage;
import com.brickwise.model.Property;
import com.brickwise.model.PropertyBreakdown;
import com.brickwise.model.QboModels.QBOStatus;
import com.brickwise.model.Residence;
import com.brickwise.repository.MortgageRepository;
import com.brickwise.repository.PropertyRepository;
import com.brickwise.repository.ResidenceRepository;
import org.springframework.stereotype.Service;

import java.time.Instant;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

@Service
public class DashboardService {
    private final PropertyRepository propertyRepo;
    private final MortgageRepository mortgageRepo;
    private final ResidenceRepository residenceRepo;
    private final QboService qboService;

    public DashboardService(
            PropertyRepository propertyRepo,
            MortgageRepository mortgageRepo,
            ResidenceRepository residenceRepo,
            QboService qboService) {
        this.propertyRepo = propertyRepo;
        this.mortgageRepo = mortgageRepo;
        this.residenceRepo = residenceRepo;
        this.qboService = qboService;
    }

    public DashboardSummary getSummary() {
        List<Property> properties = propertyRepo.findAll();
        Map<Integer, Mortgage> mortgageByProperty = new HashMap<>();
        for (Mortgage m : mortgageRepo.findAll()) {
            mortgageByProperty.put(m.getPropertyId(), m);
        }
        Map<Integer, Double> rentByProperty = new HashMap<>();
        for (Residence r : residenceRepo.findAll(null, true)) {
            rentByProperty.merge(r.getPropertyId(), r.getMonthlyRent(), Double::sum);
        }

        DashboardSummary summary = new DashboardSummary();
        summary.setAsOf(Instant.now());

        Map<Integer, Double> actualByProperty = new HashMap<>();
        Instant qboSyncedAt = null;
        try {
            QBOStatus status = qboService.status();
            summary.setQboConnected(status.isConnected());
            summary.setCashFlowPeriodStart(status.getPeriodStart());
            summary.setCashFlowPeriodEnd(status.getPeriodEnd());
            qboSyncedAt = status.getLastSyncedAt();
            actualByProperty.putAll(qboService.currentMonthCashFlow());
        } catch (Exception ignored) {
            // QBO optional
        }
        summary.setQboSyncedAt(qboSyncedAt);

        double actualTotal = 0;
        boolean hasActual = false;

        for (Property p : properties) {
            Mortgage m = mortgageByProperty.get(p.getId());
            double mortgageBalance = m == null ? 0 : m.getCurrentBalance();
            double monthlyPi = m == null ? 0 : m.getMonthlyPi();
            double monthlyRent = rentByProperty.getOrDefault(p.getId(), 0.0);
            double equity = p.getCurrentValue() - mortgageBalance;
            double net = monthlyRent - monthlyPi;

            summary.setPortfolioValue(summary.getPortfolioValue() + p.getCurrentValue());
            summary.setTotalMortgageBalance(summary.getTotalMortgageBalance() + mortgageBalance);
            summary.setTotalEquity(summary.getTotalEquity() + equity);
            summary.setMonthlyRevenue(summary.getMonthlyRevenue() + monthlyRent);
            summary.setMonthlyMortgagePi(summary.getMonthlyMortgagePi() + monthlyPi);
            summary.setEstimatedCashFlow(summary.getEstimatedCashFlow() + net);

            PropertyBreakdown b = new PropertyBreakdown();
            b.setPropertyId(p.getId());
            b.setName(p.getName());
            b.setCurrentValue(p.getCurrentValue());
            b.setMortgageBalance(mortgageBalance);
            b.setEquity(equity);
            b.setMonthlyRent(monthlyRent);
            b.setMonthlyPi(monthlyPi);
            b.setNetCashFlow(net);
            if (actualByProperty.containsKey(p.getId())) {
                double actual = actualByProperty.get(p.getId());
                b.setActualCashFlow(actual);
                b.setCashFlowSource("quickbooks");
                actualTotal += actual;
                hasActual = true;
            } else {
                b.setCashFlowSource("estimate");
            }
            summary.getProperties().add(b);
        }
        if (hasActual) {
            summary.setActualCashFlow(actualTotal);
        }
        return summary;
    }
}
