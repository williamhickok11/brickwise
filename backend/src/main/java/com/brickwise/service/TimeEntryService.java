package com.brickwise.service;

import com.brickwise.error.ApiException;
import com.brickwise.model.CreateTimeEntryRequest;
import com.brickwise.model.TimeEntry;
import com.brickwise.model.UpdateTimeEntryRequest;
import com.brickwise.repository.PropertyRepository;
import com.brickwise.repository.TimeEntryRepository;
import com.brickwise.validation.DomainValidator;
import org.springframework.stereotype.Service;

import java.time.LocalDate;
import java.util.List;

@Service
public class TimeEntryService {
    private final TimeEntryRepository repo;
    private final PropertyRepository propertyRepo;

    public TimeEntryService(TimeEntryRepository repo, PropertyRepository propertyRepo) {
        this.repo = repo;
        this.propertyRepo = propertyRepo;
    }

    public TimeEntry create(CreateTimeEntryRequest req) {
        DomainValidator.validateTimeEntry(req);
        if (req.getPropertyId() != null) {
            propertyRepo.findById(req.getPropertyId()).orElseThrow(ApiException::notFound);
        }
        return repo.create(fromRequest(req));
    }

    public TimeEntry get(int id) {
        return repo.findById(id).orElseThrow(ApiException::notFound);
    }

    public List<TimeEntry> list(Integer propertyId, String startDate, String endDate, String category) {
        return repo.findAll(propertyId, startDate, endDate, category);
    }

    public TimeEntry update(int id, UpdateTimeEntryRequest req) {
        DomainValidator.validateTimeEntry(req);
        if (req.getPropertyId() != null) {
            propertyRepo.findById(req.getPropertyId()).orElseThrow(ApiException::notFound);
        }
        TimeEntry t = fromRequest(req);
        t.setId(id);
        try {
            repo.update(t);
        } catch (IllegalStateException e) {
            throw ApiException.notFound();
        }
        return get(id);
    }

    public void delete(int id) {
        try {
            repo.delete(id);
        } catch (IllegalStateException e) {
            throw ApiException.notFound();
        }
    }

    private static TimeEntry fromRequest(CreateTimeEntryRequest req) {
        TimeEntry t = new TimeEntry();
        t.setPropertyId(req.getPropertyId());
        t.setDate(LocalDate.parse(req.getDate()));
        t.setCategory(req.getCategory());
        t.setDescription(req.getDescription());
        t.setHours(req.getHours());
        t.setNotes(req.getNotes() == null ? "" : req.getNotes());
        t.setMileage(req.getMileage());
        t.setFullDrive(req.isFullDrive());
        return t;
    }
}
