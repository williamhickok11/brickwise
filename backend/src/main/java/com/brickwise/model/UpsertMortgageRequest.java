package com.brickwise.model;

public class UpsertMortgageRequest {
    private double originalLoan;
    private double currentBalance;
    private double interestRate;
    private double monthlyPayment;
    private double monthlyPi;
    private String balanceAsOf;
    private String lender;
    private String notes;

    public double getOriginalLoan() { return originalLoan; }
    public void setOriginalLoan(double originalLoan) { this.originalLoan = originalLoan; }
    public double getCurrentBalance() { return currentBalance; }
    public void setCurrentBalance(double currentBalance) { this.currentBalance = currentBalance; }
    public double getInterestRate() { return interestRate; }
    public void setInterestRate(double interestRate) { this.interestRate = interestRate; }
    public double getMonthlyPayment() { return monthlyPayment; }
    public void setMonthlyPayment(double monthlyPayment) { this.monthlyPayment = monthlyPayment; }
    public double getMonthlyPi() { return monthlyPi; }
    public void setMonthlyPi(double monthlyPi) { this.monthlyPi = monthlyPi; }
    public String getBalanceAsOf() { return balanceAsOf; }
    public void setBalanceAsOf(String balanceAsOf) { this.balanceAsOf = balanceAsOf; }
    public String getLender() { return lender; }
    public void setLender(String lender) { this.lender = lender; }
    public String getNotes() { return notes; }
    public void setNotes(String notes) { this.notes = notes; }
}
