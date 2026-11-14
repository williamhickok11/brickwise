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

// PropertyImageHandler handles photo upload/list/main/delete for properties.
type PropertyImageHandler struct {
	service *service.PropertyImageService
	logger  *logger.Logger
}

func NewPropertyImageHandler(service *service.PropertyImageService) *PropertyImageHandler {
	return &PropertyImageHandler{
		service: service,
		logger:  logger.Default(),
	}
}

// ListImages handles GET /properties/{id}/images
func (h *PropertyImageHandler) ListImages(w http.ResponseWriter, r *http.Request) {
	propertyID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		errors.WriteJSON(w, errors.ErrInvalidInput)
		return
	}

	images, err := h.service.List(propertyID)
	if err != nil {
		h.handleError(w, r, err)
		return
	}
	h.respondJSON(w, http.StatusOK, images)
}

// UploadImage handles POST /properties/{id}/images (multipart field "file").
func (h *PropertyImageHandler) UploadImage(w http.ResponseWriter, r *http.Request) {
	propertyID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		errors.WriteJSON(w, errors.ErrInvalidInput)
		return
	}

	if err := r.ParseMultipartForm(9 << 20); err != nil {
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
	img, err := h.service.Upload(propertyID, header.Filename, contentType, header.Size, file)
	if err != nil {
		h.handleError(w, r, err)
		return
	}
	h.respondJSON(w, http.StatusCreated, img)
}

// SetMainImage handles PUT /properties/{id}/images/{imageId}/main
func (h *PropertyImageHandler) SetMainImage(w http.ResponseWriter, r *http.Request) {
	propertyID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		errors.WriteJSON(w, errors.ErrInvalidInput)
		return
	}
	imageID, err := strconv.Atoi(chi.URLParam(r, "imageId"))
	if err != nil {
		errors.WriteJSON(w, errors.ErrInvalidInput)
		return
	}

	img, err := h.service.SetMain(propertyID, imageID)
	if err != nil {
		h.handleError(w, r, err)
		return
	}
	h.respondJSON(w, http.StatusOK, img)
}

// DeleteImage handles DELETE /properties/{id}/images/{imageId}
func (h *PropertyImageHandler) DeleteImage(w http.ResponseWriter, r *http.Request) {
	propertyID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		errors.WriteJSON(w, errors.ErrInvalidInput)
		return
	}
	imageID, err := strconv.Atoi(chi.URLParam(r, "imageId"))
	if err != nil {
		errors.WriteJSON(w, errors.ErrInvalidInput)
		return
	}

	if err := h.service.Delete(propertyID, imageID); err != nil {
		h.handleError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// ServeImage handles GET /properties/{id}/images/{imageId}/file
func (h *PropertyImageHandler) ServeImage(w http.ResponseWriter, r *http.Request) {
	propertyID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		errors.WriteJSON(w, errors.ErrInvalidInput)
		return
	}
	imageID, err := strconv.Atoi(chi.URLParam(r, "imageId"))
	if err != nil {
		errors.WriteJSON(w, errors.ErrInvalidInput)
		return
	}

	img, f, err := h.service.OpenForDownload(propertyID, imageID)
	if err != nil {
		h.handleError(w, r, err)
		return
	}
	defer f.Close()

	w.Header().Set("Content-Type", img.ContentType)
	w.Header().Set("Content-Disposition", "inline; filename=\""+filepath.Base(img.FileName)+"\"")
	if img.Size > 0 {
		w.Header().Set("Content-Length", strconv.FormatInt(img.Size, 10))
	}
	w.WriteHeader(http.StatusOK)
	if _, err := io.Copy(w, f); err != nil {
		h.logger.Error("Failed to write property image", "error", err, "path", r.URL.Path)
	}
}

func (h *PropertyImageHandler) handleError(w http.ResponseWriter, r *http.Request, err error) {
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

func (h *PropertyImageHandler) respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		h.logger.Error("Failed to encode JSON response", "error", err)
	}
}
