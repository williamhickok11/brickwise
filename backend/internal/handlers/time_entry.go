package handlers

import (
	"encoding/csv"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"brickwise/backend/internal/errors"
	"brickwise/backend/internal/logger"
	"brickwise/backend/internal/models"
	"brickwise/backend/internal/service"

	"github.com/go-chi/chi/v5"
)

// TimeEntryHandler handles HTTP requests for time entries
type TimeEntryHandler struct {
	service *service.TimeEntryService
	logger  *logger.Logger
}

// NewTimeEntryHandler creates a new time entry handler
func NewTimeEntryHandler(service *service.TimeEntryService) *TimeEntryHandler {
	return &TimeEntryHandler{
		service: service,
		logger:  logger.Default(),
	}
}

// ListTimeEntries handles GET /time-entries
func (h *TimeEntryHandler) ListTimeEntries(w http.ResponseWriter, r *http.Request) {
	filter := &models.TimeEntryFilter{}

	// Parse query parameters
	if propertyIDStr := r.URL.Query().Get("property_id"); propertyIDStr != "" {
		if propertyIDStr == "null" || propertyIDStr == "0" {
			zero := 0
			filter.PropertyID = &zero
		} else {
			propertyID, err := strconv.Atoi(propertyIDStr)
			if err == nil {
				filter.PropertyID = &propertyID
			}
		}
	}

	if startDate := r.URL.Query().Get("start_date"); startDate != "" {
		filter.StartDate = &startDate
	}

	if endDate := r.URL.Query().Get("end_date"); endDate != "" {
		filter.EndDate = &endDate
	}

	if category := r.URL.Query().Get("category"); category != "" {
		filter.Category = &category
	}

	entries, err := h.service.ListTimeEntries(filter)
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	h.respondJSON(w, http.StatusOK, entries)
}

// GetTimeEntry handles GET /time-entries/{id}
func (h *TimeEntryHandler) GetTimeEntry(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		errors.WriteJSON(w, errors.ErrInvalidInput)
		return
	}

	entry, err := h.service.GetTimeEntry(id)
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	h.respondJSON(w, http.StatusOK, entry)
}

// CreateTimeEntry handles POST /time-entries
func (h *TimeEntryHandler) CreateTimeEntry(w http.ResponseWriter, r *http.Request) {
	var req models.CreateTimeEntryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Warn("Failed to decode request body", "error", err)
		errors.WriteJSON(w, errors.ErrInvalidInput)
		return
	}

	entry, err := h.service.CreateTimeEntry(&req)
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	h.respondJSON(w, http.StatusCreated, entry)
}

// UpdateTimeEntry handles PUT /time-entries/{id}
func (h *TimeEntryHandler) UpdateTimeEntry(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		errors.WriteJSON(w, errors.ErrInvalidInput)
		return
	}

	var req models.UpdateTimeEntryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Warn("Failed to decode request body", "error", err)
		errors.WriteJSON(w, errors.ErrInvalidInput)
		return
	}

	entry, err := h.service.UpdateTimeEntry(id, &req)
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	h.respondJSON(w, http.StatusOK, entry)
}

// DeleteTimeEntry handles DELETE /time-entries/{id}
func (h *TimeEntryHandler) DeleteTimeEntry(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		errors.WriteJSON(w, errors.ErrInvalidInput)
		return
	}

	if err := h.service.DeleteTimeEntry(id); err != nil {
		h.handleError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ExportTimeEntries handles GET /time-entries/export
func (h *TimeEntryHandler) ExportTimeEntries(w http.ResponseWriter, r *http.Request) {
	filter := &models.TimeEntryFilter{}

	// Parse query parameters
	if propertyIDStr := r.URL.Query().Get("property_id"); propertyIDStr != "" {
		if propertyIDStr == "null" || propertyIDStr == "0" {
			zero := 0
			filter.PropertyID = &zero
		} else {
			propertyID, err := strconv.Atoi(propertyIDStr)
			if err == nil {
				filter.PropertyID = &propertyID
			}
		}
	}

	if startDate := r.URL.Query().Get("start_date"); startDate != "" {
		filter.StartDate = &startDate
	}

	if endDate := r.URL.Query().Get("end_date"); endDate != "" {
		filter.EndDate = &endDate
	}

	if category := r.URL.Query().Get("category"); category != "" {
		filter.Category = &category
	}

	entries, err := h.service.ListTimeEntries(filter)
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	// Set headers for CSV download
	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", "attachment; filename=reps_time_entries.csv")

	writer := csv.NewWriter(w)
	defer writer.Flush()

	// Write header
	headers := []string{"Date", "Activity Category", "Description", "Time Spent (hrs)", "Notes", "Mileage"}
	if err := writer.Write(headers); err != nil {
		h.logger.Error("Failed to write CSV header", "error", err)
		return
	}

	// Write data rows
	for _, entry := range entries {
		propertyName := "General"
		if entry.PropertyID != nil {
			// In a real implementation, you'd fetch the property name
			// For now, we'll just use the ID
			propertyName = ""
		}

		row := []string{
			entry.Date.Time.Format("2006-01-02"),
			entry.Category,
			entry.Description,
			strconv.FormatFloat(entry.Hours, 'f', 2, 64),
			entry.Notes,
			strconv.FormatFloat(entry.Mileage, 'f', 2, 64),
		}
		if err := writer.Write(row); err != nil {
			h.logger.Error("Failed to write CSV row", "error", err)
			return
		}
	}
}

// handleError processes errors and writes appropriate HTTP responses
func (h *TimeEntryHandler) handleError(w http.ResponseWriter, r *http.Request, err error) {
	apiErr, ok := err.(*errors.APIError)
	if !ok {
		// Wrap unknown errors as internal server errors
		h.logger.Error("Unexpected error", "error", err, "path", r.URL.Path)
		apiErr = errors.ErrInternal
	} else {
		// Log non-internal errors at warn level
		if apiErr.Status >= http.StatusInternalServerError {
			h.logger.Error("Internal error", "error", err, "path", r.URL.Path)
		} else {
			h.logger.Warn("Request error", "error", err, "path", r.URL.Path, "status", apiErr.Status)
		}
	}

	errors.WriteJSON(w, apiErr)
}

// respondJSON writes a JSON response
func (h *TimeEntryHandler) respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		h.logger.Error("Failed to encode JSON response", "error", err)
	}
}
