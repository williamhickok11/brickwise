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

// PropertyHandler handles HTTP requests for properties
type PropertyHandler struct {
	service *service.PropertyService
	logger  *logger.Logger
}

// NewPropertyHandler creates a new property handler
func NewPropertyHandler(service *service.PropertyService) *PropertyHandler {
	return &PropertyHandler{
		service: service,
		logger:  logger.Default(),
	}
}

// ListProperties handles GET /properties
func (h *PropertyHandler) ListProperties(w http.ResponseWriter, r *http.Request) {
	properties, err := h.service.ListProperties()
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	h.respondJSON(w, http.StatusOK, properties)
}

// GetProperty handles GET /properties/{id}
func (h *PropertyHandler) GetProperty(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		errors.WriteJSON(w, errors.ErrInvalidInput)
		return
	}

	property, err := h.service.GetProperty(id)
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	h.respondJSON(w, http.StatusOK, property)
}

// CreateProperty handles POST /properties
func (h *PropertyHandler) CreateProperty(w http.ResponseWriter, r *http.Request) {
	var req models.CreatePropertyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Warn("Failed to decode request body", "error", err)
		errors.WriteJSON(w, errors.ErrInvalidInput)
		return
	}

	property, err := h.service.CreateProperty(&req)
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	h.respondJSON(w, http.StatusCreated, property)
}

// UpdateProperty handles PUT /properties/{id}
func (h *PropertyHandler) UpdateProperty(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		errors.WriteJSON(w, errors.ErrInvalidInput)
		return
	}

	var req models.CreatePropertyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Warn("Failed to decode request body", "error", err)
		errors.WriteJSON(w, errors.ErrInvalidInput)
		return
	}

	property, err := h.service.UpdateProperty(id, &req)
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	h.respondJSON(w, http.StatusOK, property)
}

// DeleteProperty handles DELETE /properties/{id}
func (h *PropertyHandler) DeleteProperty(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		errors.WriteJSON(w, errors.ErrInvalidInput)
		return
	}

	if err := h.service.DeleteProperty(id); err != nil {
		h.handleError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// handleError processes errors and writes appropriate HTTP responses
func (h *PropertyHandler) handleError(w http.ResponseWriter, r *http.Request, err error) {
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
func (h *PropertyHandler) respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		h.logger.Error("Failed to encode JSON response", "error", err)
	}
}
