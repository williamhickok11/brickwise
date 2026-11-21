package com.brickwise.validation;

import com.brickwise.error.ApiException;
import com.brickwise.model.CreatePropertyRequest;
import com.brickwise.model.CreateResidenceRequest;
import com.brickwise.model.CreateTimeEntryRequest;

import java.time.LocalDate;
import java.time.format.DateTimeParseException;
import java.util.ArrayList;
import java.util.List;
import java.util.Locale;
import java.util.Set;

/**
 * Hand-rolled domain validation (parity with Go {@code internal/validator}).
 *
 * <p><b>C# analogy:</b> FluentValidation or DataAnnotations checks collected into one error message.
 * Spring also offers Bean Validation ({@code @NotBlank}), but we keep explicit checks so the
 * learning path mirrors the Go code and the error messages stay identical for the frontend.
 */
public final class DomainValidator {
    private static final Set<String> PROPERTY_TYPES = Set.of(
            "residential", "commercial", "industrial", "land", "mixed");
    private static final Set<String> TIME_CATEGORIES = Set.of(
            "Property Management",
            "Maintenance & Repairs",
            "Contractor Oversight",
            "Accounting & Admin",
            "Deal Sourcing",
            "Construction Oversight",
            "Software Management");

    private DomainValidator() {}

    public static void validateProperty(CreatePropertyRequest req) {
        List<String> errs = new ArrayList<>();
        if (req.getName() == null || req.getName().isBlank()) errs.add("name: name is required");
        else if (req.getName().length() > 255) errs.add("name: name must be 255 characters or less");
        if (req.getAddress() == null || req.getAddress().isBlank()) errs.add("address: address is required");
        else if (req.getAddress().length() > 500) errs.add("address: address must be 500 characters or less");
        if (req.getPropertyType() == null || req.getPropertyType().isBlank()) {
            errs.add("property_type: property_type is required");
        } else if (!PROPERTY_TYPES.contains(req.getPropertyType().toLowerCase(Locale.ROOT))) {
            errs.add("property_type: property_type must be one of: residential, commercial, industrial, land, mixed");
        }
        throwIf(errs);
    }

    public static void validateResidence(CreateResidenceRequest req) {
        List<String> errs = new ArrayList<>();
        if (req.getPropertyId() <= 0) errs.add("property_id: property_id is required");
        if (req.getName() == null || req.getName().isBlank()) errs.add("name: name is required");
        else if (req.getName().length() > 255) errs.add("name: name must be 255 characters or less");
        if (req.getStartDate() == null || req.getStartDate().isBlank()) {
            errs.add("start_date: start_date is required");
        } else {
            try { LocalDate.parse(req.getStartDate()); }
            catch (DateTimeParseException e) { errs.add("start_date: start_date must be in YYYY-MM-DD format"); }
        }
        if (req.getEndDate() != null && !req.getEndDate().isBlank()) {
            try { LocalDate.parse(req.getEndDate()); }
            catch (DateTimeParseException e) { errs.add("end_date: end_date must be in YYYY-MM-DD format"); }
        }
        if (req.getPhone() != null && req.getPhone().length() > 255) errs.add("phone: phone must be 255 characters or less");
        if (req.getEmail() != null && req.getEmail().length() > 255) errs.add("email: email must be 255 characters or less");
        throwIf(errs);
    }

    public static void validateTimeEntry(CreateTimeEntryRequest req) {
        List<String> errs = new ArrayList<>();
        if (req.getDate() == null || req.getDate().isBlank()) errs.add("date: date is required");
        if (req.getCategory() == null || req.getCategory().isBlank()) errs.add("category: category is required");
        else if (!TIME_CATEGORIES.contains(req.getCategory())) {
            errs.add("category: category must be one of: " + String.join(", ", TIME_CATEGORIES));
        }
        if (req.getDescription() == null || req.getDescription().isBlank()) errs.add("description: description is required");
        else if (req.getDescription().length() > 2000) errs.add("description: description must be 2000 characters or less");
        if (req.getHours() <= 0) errs.add("hours: hours must be greater than 0");
        if (req.getHours() > 24) errs.add("hours: hours must be 24 or less");
        throwIf(errs);
    }

    private static void throwIf(List<String> errs) {
        if (!errs.isEmpty()) {
            throw ApiException.validation(String.join(", ", errs));
        }
    }
}
