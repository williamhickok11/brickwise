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

type MortgageHandler struct {
	service *service.MortgageService
	logger  *logger.Logger
}

func NewMortgageHandler(service *service.MortgageService) *MortgageHandler {
	return &MortgageHandler{
		service: service,
		logger:  logger.Default(),
	}
}

func (h *MortgageHandler) GetMortgage(w http.ResponseWriter, r *http.Request) {
	propertyID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		errors.WriteJSON(w, errors.ErrInvalidInput)
		return
	}

	mortgage, err := h.service.GetMortgage(propertyID)
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	h.respondJSON(w, http.StatusOK, mortgage)
}

func (h *MortgageHandler) UpsertMortgage(w http.ResponseWriter, r *http.Request) {
	propertyID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		errors.WriteJSON(w, errors.ErrInvalidInput)
		return
	}

	var req models.UpsertMortgageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Warn("Failed to decode request body", "error", err)
		errors.WriteJSON(w, errors.ErrInvalidInput)
		return
	}

	mortgage, err := h.service.UpsertMortgage(propertyID, &req)
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	h.respondJSON(w, http.StatusOK, mortgage)
}

func (h *MortgageHandler) DeleteMortgage(w http.ResponseWriter, r *http.Request) {
	propertyID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		errors.WriteJSON(w, errors.ErrInvalidInput)
		return
	}

	if err := h.service.DeleteMortgage(propertyID); err != nil {
		h.handleError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *MortgageHandler) handleError(w http.ResponseWriter, r *http.Request, err error) {
	apiErr, ok := err.(*errors.APIError)
	if !ok {
		h.logger.Error("Unexpected error", "error", err, "path", r.URL.Path)
		apiErr = errors.ErrInternal
	} else if apiErr.Status >= http.StatusInternalServerError {
		h.logger.Error("Internal error", "error", err, "path", r.URL.Path)
	} else {
		h.logger.Warn("Request error", "error", err, "path", r.URL.Path, "status", apiErr.Status)
	}
	errors.WriteJSON(w, apiErr)
}

func (h *MortgageHandler) respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		h.logger.Error("Failed to encode JSON response", "error", err)
	}
}
