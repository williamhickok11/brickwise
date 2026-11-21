package com.brickwise.model;

import java.time.Instant;
import java.util.ArrayList;
import java.util.List;

public class DashboardSummary {
    private Instant asOf;
    private double portfolioValue;
    private double totalMortgageBalance;
    private double totalEquity;
    private double monthlyRevenue;
    private double monthlyMortgagePi;
    private double estimatedCashFlow;
    private Double actualCashFlow;
    private boolean qboConnected;
    private Instant qboSyncedAt;
    private Instant cashFlowPeriodStart;
    private Instant cashFlowPeriodEnd;
    private List<PropertyBreakdown> properties = new ArrayList<>();

    public Instant getAsOf() { return asOf; }
    public void setAsOf(Instant asOf) { this.asOf = asOf; }
    public double getPortfolioValue() { return portfolioValue; }
    public void setPortfolioValue(double portfolioValue) { this.portfolioValue = portfolioValue; }
    public double getTotalMortgageBalance() { return totalMortgageBalance; }
    public void setTotalMortgageBalance(double totalMortgageBalance) { this.totalMortgageBalance = totalMortgageBalance; }
    public double getTotalEquity() { return totalEquity; }
    public void setTotalEquity(double totalEquity) { this.totalEquity = totalEquity; }
    public double getMonthlyRevenue() { return monthlyRevenue; }
    public void setMonthlyRevenue(double monthlyRevenue) { this.monthlyRevenue = monthlyRevenue; }
    public double getMonthlyMortgagePi() { return monthlyMortgagePi; }
    public void setMonthlyMortgagePi(double monthlyMortgagePi) { this.monthlyMortgagePi = monthlyMortgagePi; }
    public double getEstimatedCashFlow() { return estimatedCashFlow; }
    public void setEstimatedCashFlow(double estimatedCashFlow) { this.estimatedCashFlow = estimatedCashFlow; }
    public Double getActualCashFlow() { return actualCashFlow; }
    public void setActualCashFlow(Double actualCashFlow) { this.actualCashFlow = actualCashFlow; }
    public boolean isQboConnected() { return qboConnected; }
    public void setQboConnected(boolean qboConnected) { this.qboConnected = qboConnected; }
    public Instant getQboSyncedAt() { return qboSyncedAt; }
    public void setQboSyncedAt(Instant qboSyncedAt) { this.qboSyncedAt = qboSyncedAt; }
    public Instant getCashFlowPeriodStart() { return cashFlowPeriodStart; }
    public void setCashFlowPeriodStart(Instant cashFlowPeriodStart) { this.cashFlowPeriodStart = cashFlowPeriodStart; }
    public Instant getCashFlowPeriodEnd() { return cashFlowPeriodEnd; }
    public void setCashFlowPeriodEnd(Instant cashFlowPeriodEnd) { this.cashFlowPeriodEnd = cashFlowPeriodEnd; }
    public List<PropertyBreakdown> getProperties() { return properties; }
    public void setProperties(List<PropertyBreakdown> properties) { this.properties = properties; }
}
