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
        TimeEntry entry = fromRequest(req);
        applyFullDriveMileage(entry);
        return repo.create(entry);
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
        applyFullDriveMileage(t);
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

    /** Builds a CSV download body for the given REPS time entries. */
    public String toCsv(List<TimeEntry> entries) {
        StringBuilder csv = new StringBuilder(
                "Date,Activity Category,Description,Time Spent (hrs),Notes,Mileage,Property ID\n");
        for (TimeEntry e : entries) {
            csv.append(e.getDate()).append(',')
                    .append(quote(e.getCategory())).append(',')
                    .append(quote(e.getDescription())).append(',')
                    .append(e.getHours()).append(',')
                    .append(quote(e.getNotes())).append(',')
                    .append(e.getMileage()).append(',')
                    .append(e.getPropertyId() == null ? "" : e.getPropertyId())
                    .append('\n');
        }
        return csv.toString();
    }

    private static String quote(String s) {
        if (s == null) {
            return "\"\"";
        }
        return "\"" + s.replace("\"", "\"\"") + "\"";
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

    /** When Full Drive is checked, fill mileage from the property's default round-trip miles. */
    private void applyFullDriveMileage(TimeEntry entry) {
        if (!entry.isFullDrive() || entry.getPropertyId() == null) {
            return;
        }
        propertyRepo.findById(entry.getPropertyId()).ifPresent(property -> {
            if (property.getDefaultMileage() > 0) {
                entry.setMileage(property.getDefaultMileage());
            }
        });
    }
}
