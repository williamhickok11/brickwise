package com.brickwise.service;

import com.brickwise.error.ApiException;
import com.brickwise.model.Mortgage;
import com.brickwise.model.UpsertMortgageRequest;
import com.brickwise.repository.MortgageRepository;
import com.brickwise.repository.PropertyRepository;
import org.springframework.stereotype.Service;

import java.time.LocalDate;

@Service
public class MortgageService {
    private final MortgageRepository repo;
    private final PropertyRepository propertyRepo;

    public MortgageService(MortgageRepository repo, PropertyRepository propertyRepo) {
        this.repo = repo;
        this.propertyRepo = propertyRepo;
    }

    public Mortgage get(int propertyId) {
        propertyRepo.findById(propertyId).orElseThrow(ApiException::notFound);
        return repo.findByPropertyId(propertyId).orElseThrow(ApiException::notFound);
    }

    public Mortgage upsert(int propertyId, UpsertMortgageRequest req) {
        propertyRepo.findById(propertyId).orElseThrow(ApiException::notFound);
        Mortgage m = new Mortgage();
        m.setPropertyId(propertyId);
        m.setOriginalLoan(req.getOriginalLoan());
        m.setCurrentBalance(req.getCurrentBalance());
        m.setInterestRate(req.getInterestRate());
        m.setMonthlyPayment(req.getMonthlyPayment());
        m.setMonthlyPi(req.getMonthlyPi());
        if (req.getBalanceAsOf() != null && !req.getBalanceAsOf().isBlank()) {
            m.setBalanceAsOf(LocalDate.parse(req.getBalanceAsOf()));
        }
        m.setLender(req.getLender() == null ? "" : req.getLender());
        m.setNotes(req.getNotes() == null ? "" : req.getNotes());
        return repo.upsert(propertyId, m);
    }

    public void delete(int propertyId) {
        propertyRepo.findById(propertyId).orElseThrow(ApiException::notFound);
        try {
            repo.deleteByPropertyId(propertyId);
        } catch (IllegalStateException e) {
            throw ApiException.notFound();
        }
    }
}
