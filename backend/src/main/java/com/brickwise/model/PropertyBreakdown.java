package com.brickwise.model;

public class PropertyBreakdown {
    private int propertyId;
    private String name;
    private double currentValue;
    private double mortgageBalance;
    private double equity;
    private double monthlyRent;
    private double monthlyPi;
    private double netCashFlow;
    private Double actualCashFlow;
    private String cashFlowSource;

    public int getPropertyId() { return propertyId; }
    public void setPropertyId(int propertyId) { this.propertyId = propertyId; }
    public String getName() { return name; }
    public void setName(String name) { this.name = name; }
    public double getCurrentValue() { return currentValue; }
    public void setCurrentValue(double currentValue) { this.currentValue = currentValue; }
    public double getMortgageBalance() { return mortgageBalance; }
    public void setMortgageBalance(double mortgageBalance) { this.mortgageBalance = mortgageBalance; }
    public double getEquity() { return equity; }
    public void setEquity(double equity) { this.equity = equity; }
    public double getMonthlyRent() { return monthlyRent; }
    public void setMonthlyRent(double monthlyRent) { this.monthlyRent = monthlyRent; }
    public double getMonthlyPi() { return monthlyPi; }
    public void setMonthlyPi(double monthlyPi) { this.monthlyPi = monthlyPi; }
    public double getNetCashFlow() { return netCashFlow; }
    public void setNetCashFlow(double netCashFlow) { this.netCashFlow = netCashFlow; }
    public Double getActualCashFlow() { return actualCashFlow; }
    public void setActualCashFlow(Double actualCashFlow) { this.actualCashFlow = actualCashFlow; }
    public String getCashFlowSource() { return cashFlowSource; }
    public void setCashFlowSource(String cashFlowSource) { this.cashFlowSource = cashFlowSource; }
}
