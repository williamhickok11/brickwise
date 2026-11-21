package com.brickwise.service;

import com.brickwise.error.ApiException;
import com.brickwise.model.CreateResidenceRequest;
import com.brickwise.model.Residence;
import com.brickwise.model.UpdateResidenceRequest;
import com.brickwise.repository.PropertyRepository;
import com.brickwise.repository.ResidenceRepository;
import com.brickwise.validation.DomainValidator;
import org.springframework.stereotype.Service;

import java.time.LocalDate;
import java.util.List;

/**
 * Residences are the tenant/occupant history for a property.
 *
 * <p>There is no separate Tenant entity yet — this mirrors the Go API. Soft move-out uses
 * {@code is_active=false} plus an optional {@code end_date}.
 */
@Service
public class ResidenceService {
    private final ResidenceRepository repo;
    private final PropertyRepository propertyRepo;

    public ResidenceService(ResidenceRepository repo, PropertyRepository propertyRepo) {
        this.repo = repo;
        this.propertyRepo = propertyRepo;
    }

    public Residence create(CreateResidenceRequest req) {
        DomainValidator.validateResidence(req);
        propertyRepo.findById(req.getPropertyId()).orElseThrow(ApiException::notFound);
        return repo.create(fromRequest(req));
    }

    public Residence get(int id) {
        return repo.findById(id).orElseThrow(ApiException::notFound);
    }

    public List<Residence> list(Integer propertyId, Boolean isActive) {
        return repo.findAll(propertyId, isActive);
    }

    public Residence update(int id, UpdateResidenceRequest req) {
        DomainValidator.validateResidence(req);
        propertyRepo.findById(req.getPropertyId()).orElseThrow(ApiException::notFound);
        Residence r = fromRequest(req);
        r.setId(id);
        try {
            repo.update(r);
        } catch (IllegalStateException e) {
            throw ApiException.notFound();
        }
        return get(id);
    }

    private static Residence fromRequest(CreateResidenceRequest req) {
        Residence r = new Residence();
        r.setPropertyId(req.getPropertyId());
        r.setName(req.getName());
        r.setPhone(req.getPhone() == null ? "" : req.getPhone());
        r.setEmail(req.getEmail() == null ? "" : req.getEmail());
        r.setStartDate(LocalDate.parse(req.getStartDate()));
        if (req.getEndDate() != null && !req.getEndDate().isBlank()) {
            r.setEndDate(LocalDate.parse(req.getEndDate()));
        }
        r.setActive(req.isActive());
        r.setMonthlyRent(req.getMonthlyRent());
        return r;
    }
}
