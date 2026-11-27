package com.brickwise.config;

import com.brickwise.model.CreatePropertyRequest;
import com.brickwise.model.CreateTimeEntryRequest;
import com.brickwise.model.Property;
import com.brickwise.repository.PropertyRepository;
import com.brickwise.repository.TimeEntryRepository;
import com.brickwise.service.PropertyService;
import com.brickwise.service.TimeEntryService;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.boot.ApplicationArguments;
import org.springframework.boot.ApplicationRunner;
import org.springframework.stereotype.Component;

import java.time.LocalDate;
import java.util.List;

/**
 * Inserts a couple of realistic sample REPS (time entries) when the DB is empty,
 * so local demos and first-run setups have data to export / email.
 */
@Component
public class SampleDataSeeder implements ApplicationRunner {
    private static final Logger log = LoggerFactory.getLogger(SampleDataSeeder.class);

    private final TimeEntryRepository timeEntryRepo;
    private final PropertyRepository propertyRepo;
    private final PropertyService propertyService;
    private final TimeEntryService timeEntryService;

    public SampleDataSeeder(
            TimeEntryRepository timeEntryRepo,
            PropertyRepository propertyRepo,
            PropertyService propertyService,
            TimeEntryService timeEntryService) {
        this.timeEntryRepo = timeEntryRepo;
        this.propertyRepo = propertyRepo;
        this.propertyService = propertyService;
        this.timeEntryService = timeEntryService;
    }

    @Override
    public void run(ApplicationArguments args) {
        if (!timeEntryRepo.findAll(null, null, null, null).isEmpty()) {
            return;
        }

        Integer propertyId = ensureSampleProperty();
        LocalDate today = LocalDate.now();

        timeEntryService.create(entry(
                propertyId,
                today.minusDays(5).toString(),
                "Property Management",
                "Walked Oak Street duplex; checked smoke detectors and noted porch light out",
                2.5,
                "Tenant mentioned porch light; ordered bulb",
                24.0,
                true));

        timeEntryService.create(entry(
                propertyId,
                today.minusDays(2).toString(),
                "Maintenance & Repairs",
                "Replaced kitchen faucet cartridge and re-caulked tub surround",
                1.75,
                "Parts from Ace Hardware ~$42",
                12.0,
                false));

        timeEntryService.create(entry(
                null,
                today.minusDays(1).toString(),
                "Accounting & Admin",
                "Reconciled October REPS hours and mileage for quarterly tax packet",
                3.0,
                "General admin — no property drive",
                0.0,
                false));

        log.info("Seeded sample REPS time entries (and a sample property if needed)");
    }

    private Integer ensureSampleProperty() {
        List<Property> existing = propertyRepo.findAll();
        if (!existing.isEmpty()) {
            return existing.get(0).getId();
        }
        CreatePropertyRequest req = new CreatePropertyRequest();
        req.setName("Oak Street Duplex");
        req.setAddress("412 Oak Street, Nashville, TN 37206");
        req.setPropertyType("residential");
        req.setDefaultMileage(24.0);
        return propertyService.create(req).getId();
    }

    private static CreateTimeEntryRequest entry(
            Integer propertyId,
            String date,
            String category,
            String description,
            double hours,
            String notes,
            double mileage,
            boolean fullDrive) {
        CreateTimeEntryRequest req = new CreateTimeEntryRequest();
        req.setPropertyId(propertyId);
        req.setDate(date);
        req.setCategory(category);
        req.setDescription(description);
        req.setHours(hours);
        req.setNotes(notes);
        req.setMileage(mileage);
        req.setFullDrive(fullDrive);
        return req;
    }
}
