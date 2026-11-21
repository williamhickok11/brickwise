package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"path/filepath"
	"strconv"

	"brickwise/backend/internal/errors"
	"brickwise/backend/internal/logger"
	"brickwise/backend/internal/service"

	"github.com/go-chi/chi/v5"
)

// DocumentHandler handles lease PDF upload/download for residences.
type DocumentHandler struct {
	service *service.DocumentService
	logger  *logger.Logger
}

func NewDocumentHandler(service *service.DocumentService) *DocumentHandler {
	return &DocumentHandler{
		service: service,
		logger:  logger.Default(),
	}
}

// GetDocument handles GET /residences/{id}/document
func (h *DocumentHandler) GetDocument(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		errors.WriteJSON(w, errors.ErrInvalidInput)
		return
	}

	doc, err := h.service.GetByResidenceID(id)
	if err != nil {
		h.handleError(w, r, err)
		return
	}
	h.respondJSON(w, http.StatusOK, doc)
}

// UploadDocument handles POST /residences/{id}/document (multipart field "file").
func (h *DocumentHandler) UploadDocument(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		errors.WriteJSON(w, errors.ErrInvalidInput)
		return
	}

	// 10MB + small overhead for multipart boundaries
	if err := r.ParseMultipartForm(11 << 20); err != nil {
		h.logger.Warn("Failed to parse multipart form", "error", err)
		errors.WriteJSON(w, errors.New("INVALID_UPLOAD", "invalid multipart upload", errors.ErrInvalidInput.Status))
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		errors.WriteJSON(w, errors.New("MISSING_FILE", "multipart field \"file\" is required", errors.ErrInvalidInput.Status))
		return
	}
	defer file.Close()

	contentType := header.Header.Get("Content-Type")
	doc, err := h.service.UploadLeasePDF(id, header.Filename, contentType, header.Size, file)
	if err != nil {
		h.handleError(w, r, err)
		return
	}
	h.respondJSON(w, http.StatusOK, doc)
}

// DownloadDocument handles GET /residences/{id}/document/download
func (h *DocumentHandler) DownloadDocument(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		errors.WriteJSON(w, errors.ErrInvalidInput)
		return
	}

	doc, f, err := h.service.OpenForDownload(id)
	if err != nil {
		h.handleError(w, r, err)
		return
	}
	defer f.Close()

	w.Header().Set("Content-Type", doc.ContentType)
	w.Header().Set("Content-Disposition", "attachment; filename=\""+filepath.Base(doc.FileName)+"\"")
	w.Header().Set("Content-Length", strconv.FormatInt(doc.Size, 10))
	w.WriteHeader(http.StatusOK)
	if _, err := io.Copy(w, f); err != nil {
		h.logger.Error("Failed to write document download", "error", err, "path", r.URL.Path)
	}
}

func (h *DocumentHandler) handleError(w http.ResponseWriter, r *http.Request, err error) {
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

func (h *DocumentHandler) respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		h.logger.Error("Failed to encode JSON response", "error", err)
	}
}
