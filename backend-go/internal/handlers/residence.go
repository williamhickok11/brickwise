package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"brickwise/backend/internal/errors"
	"brickwise/backend/internal/logger"
	"brickwise/backend/internal/models"
	"brickwise/backend/internal/service"

	"github.com/go-chi/chi/v5"
)

// ResidenceHandler handles HTTP requests for residences/tenants.
type ResidenceHandler struct {
	service *service.ResidenceService
	logger  *logger.Logger
}

func NewResidenceHandler(service *service.ResidenceService) *ResidenceHandler {
	return &ResidenceHandler{
		service: service,
		logger:  logger.Default(),
	}
}

// ListResidences handles GET /residences
func (h *ResidenceHandler) ListResidences(w http.ResponseWriter, r *http.Request) {
	filter := &models.ResidenceFilter{}

	// Optional query parameters
	if propertyIDStr := r.URL.Query().Get("property_id"); propertyIDStr != "" {
		propertyID, err := strconv.Atoi(propertyIDStr)
		if err != nil {
			errors.WriteJSON(w, errors.ErrInvalidInput)
			return
		}
		filter.PropertyID = &propertyID
	}

	if isActiveStr := r.URL.Query().Get("is_active"); isActiveStr != "" {
		// strconv.ParseBool accepts "1","t","T","TRUE","true" and "0","f","F","FALSE","false"
		isActive, err := strconv.ParseBool(isActiveStr)
		if err != nil {
			errors.WriteJSON(w, errors.ErrInvalidInput)
			return
		}
		filter.IsActive = &isActive
	}

	residences, err := h.service.ListResidences(filter)
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	h.respondJSON(w, http.StatusOK, residences)
}

// GetResidence handles GET /residences/{id}
func (h *ResidenceHandler) GetResidence(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		errors.WriteJSON(w, errors.ErrInvalidInput)
		return
	}

	residence, err := h.service.GetResidence(id)
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	h.respondJSON(w, http.StatusOK, residence)
}

// CreateResidence handles POST /residences
func (h *ResidenceHandler) CreateResidence(w http.ResponseWriter, r *http.Request) {
	var req models.CreateResidenceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Warn("Failed to decode request body", "error", err)
		errors.WriteJSON(w, errors.ErrInvalidInput)
		return
	}

	residence, err := h.service.CreateResidence(&req)
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	h.respondJSON(w, http.StatusCreated, residence)
}

// UpdateResidence handles PUT /residences/{id}
func (h *ResidenceHandler) UpdateResidence(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		errors.WriteJSON(w, errors.ErrInvalidInput)
		return
	}

	var req models.UpdateResidenceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Warn("Failed to decode request body", "error", err)
		errors.WriteJSON(w, errors.ErrInvalidInput)
		return
	}

	residence, err := h.service.UpdateResidence(id, &req)
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	h.respondJSON(w, http.StatusOK, residence)
}

// handleError processes errors and writes appropriate HTTP responses
func (h *ResidenceHandler) handleError(w http.ResponseWriter, r *http.Request, err error) {
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
func (h *ResidenceHandler) respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		h.logger.Error("Failed to encode JSON response", "error", err)
	}
}

