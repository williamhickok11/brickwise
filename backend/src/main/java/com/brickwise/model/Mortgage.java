package com.brickwise.model;

import java.time.Instant;
import java.time.LocalDate;

public class Mortgage {
    private int id;
    private int propertyId;
    private double originalLoan;
    private double currentBalance;
    private double interestRate;
    private double monthlyPayment;
    private double monthlyPi;
    private LocalDate balanceAsOf;
    private String lender;
    private String notes;
    private Instant createdAt;
    private Instant updatedAt;

    public int getId() { return id; }
    public void setId(int id) { this.id = id; }
    public int getPropertyId() { return propertyId; }
    public void setPropertyId(int propertyId) { this.propertyId = propertyId; }
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
    public LocalDate getBalanceAsOf() { return balanceAsOf; }
    public void setBalanceAsOf(LocalDate balanceAsOf) { this.balanceAsOf = balanceAsOf; }
    public String getLender() { return lender; }
    public void setLender(String lender) { this.lender = lender; }
    public String getNotes() { return notes; }
    public void setNotes(String notes) { this.notes = notes; }
    public Instant getCreatedAt() { return createdAt; }
    public void setCreatedAt(Instant createdAt) { this.createdAt = createdAt; }
    public Instant getUpdatedAt() { return updatedAt; }
    public void setUpdatedAt(Instant updatedAt) { this.updatedAt = updatedAt; }
}
