package handlers

import (
	"encoding/json"
	"net/http"

	"brickwise/backend/internal/errors"
	"brickwise/backend/internal/logger"
	"brickwise/backend/internal/models"
	"brickwise/backend/internal/service"
)

type QBOHandler struct {
	service *service.QBOService
	logger  *logger.Logger
}

func NewQBOHandler(service *service.QBOService) *QBOHandler {
	return &QBOHandler{
		service: service,
		logger:  logger.Default(),
	}
}

func (h *QBOHandler) GetStatus(w http.ResponseWriter, r *http.Request) {
	status, err := h.service.Status(r.Context())
	if err != nil {
		h.handleError(w, r, err)
		return
	}
	h.respondJSON(w, http.StatusOK, status)
}

func (h *QBOHandler) Connect(w http.ResponseWriter, r *http.Request) {
	authURL, err := h.service.StartConnect(r.Context())
	if err != nil {
		h.handleError(w, r, err)
		return
	}
	h.respondJSON(w, http.StatusOK, map[string]string{"authorize_url": authURL})
}

func (h *QBOHandler) Callback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")
	realmID := r.URL.Query().Get("realmId")
	if realmID == "" {
		realmID = r.URL.Query().Get("realm_id")
	}

	if errParam := r.URL.Query().Get("error"); errParam != "" {
		http.Redirect(w, r, h.service.FrontendRedirectURL(false, errParam), http.StatusFound)
		return
	}

	if err := h.service.HandleCallback(r.Context(), code, state, realmID); err != nil {
		h.logger.Error("QBO OAuth callback failed", "error", err)
		http.Redirect(w, r, h.service.FrontendRedirectURL(false, "oauth_failed"), http.StatusFound)
		return
	}

	http.Redirect(w, r, h.service.FrontendRedirectURL(true, ""), http.StatusFound)
}

func (h *QBOHandler) Disconnect(w http.ResponseWriter, r *http.Request) {
	if err := h.service.Disconnect(r.Context()); err != nil {
		h.handleError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *QBOHandler) ListClasses(w http.ResponseWriter, r *http.Request) {
	classes, err := h.service.ListClasses(r.Context())
	if err != nil {
		h.handleError(w, r, err)
		return
	}
	h.respondJSON(w, http.StatusOK, classes)
}

func (h *QBOHandler) UpdateMappings(w http.ResponseWriter, r *http.Request) {
	var req models.UpdateQBOMappingsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errors.WriteJSON(w, errors.ErrInvalidInput)
		return
	}
	if err := h.service.UpdateMappings(r.Context(), &req); err != nil {
		h.handleError(w, r, err)
		return
	}
	status, err := h.service.Status(r.Context())
	if err != nil {
		h.handleError(w, r, err)
		return
	}
	h.respondJSON(w, http.StatusOK, status)
}

func (h *QBOHandler) Sync(w http.ResponseWriter, r *http.Request) {
	result, err := h.service.SyncCashFlow(r.Context(), true)
	if err != nil {
		h.handleError(w, r, err)
		return
	}
	h.respondJSON(w, http.StatusOK, result)
}

func (h *QBOHandler) handleError(w http.ResponseWriter, r *http.Request, err error) {
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

func (h *QBOHandler) respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		h.logger.Error("Failed to encode JSON response", "error", err)
	}
}
