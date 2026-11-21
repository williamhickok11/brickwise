package handlers

import (
	"encoding/json"
	"net/http"

	"brickwise/backend/internal/errors"
	"brickwise/backend/internal/logger"
	"brickwise/backend/internal/service"
)

type DashboardHandler struct {
	service *service.DashboardService
	logger  *logger.Logger
}

func NewDashboardHandler(service *service.DashboardService) *DashboardHandler {
	return &DashboardHandler{
		service: service,
		logger:  logger.Default(),
	}
}

func (h *DashboardHandler) GetSummary(w http.ResponseWriter, r *http.Request) {
	summary, err := h.service.GetSummary()
	if err != nil {
		h.handleError(w, r, err)
		return
	}
	h.respondJSON(w, http.StatusOK, summary)
}

func (h *DashboardHandler) handleError(w http.ResponseWriter, r *http.Request, err error) {
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

func (h *DashboardHandler) respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		h.logger.Error("Failed to encode JSON response", "error", err)
	}
}
