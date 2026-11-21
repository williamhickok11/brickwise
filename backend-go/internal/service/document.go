package service

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"brickwise/backend/internal/errors"
	"brickwise/backend/internal/models"
	"brickwise/backend/internal/repository"
)

const (
	defaultDocumentsDir = "./data/documents"
	maxLeasePDFBytes    = 10 << 20 // 10 MiB
)

// DocumentService stores one lease PDF per residence on local disk.
type DocumentService struct {
	repo          repository.DocumentRepository
	residenceRepo repository.ResidenceRepository
	dir           string
}

func NewDocumentService(repo repository.DocumentRepository, residenceRepo repository.ResidenceRepository) *DocumentService {
	dir := os.Getenv("DOCUMENTS_DIR")
	if strings.TrimSpace(dir) == "" {
		dir = defaultDocumentsDir
	}
	return &DocumentService{
		repo:          repo,
		residenceRepo: residenceRepo,
		dir:           dir,
	}
}

func (s *DocumentService) GetByResidenceID(residenceID int) (*models.Document, error) {
	if _, err := s.residenceRepo.GetByID(residenceID); err != nil {
		return nil, errors.ErrNotFound
	}
	doc, err := s.repo.GetByResidenceID(residenceID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return nil, errors.ErrNotFound
		}
		return nil, errors.Wrap(err, errors.ErrInternal)
	}
	return doc, nil
}

// UploadLeasePDF saves a PDF for the residence, replacing any existing file.
func (s *DocumentService) UploadLeasePDF(residenceID int, fileName string, contentType string, size int64, body io.Reader) (*models.Document, error) {
	if _, err := s.residenceRepo.GetByID(residenceID); err != nil {
		return nil, errors.ErrNotFound
	}

	if size > maxLeasePDFBytes {
		return nil, errors.New("FILE_TOO_LARGE", "lease PDF must be 10MB or smaller", errors.ErrInvalidInput.Status)
	}

	lowerName := strings.ToLower(fileName)
	if !strings.HasSuffix(lowerName, ".pdf") {
		return nil, errors.New("INVALID_FILE_TYPE", "only PDF files are allowed", errors.ErrInvalidInput.Status)
	}
	ct := strings.ToLower(strings.TrimSpace(contentType))
	if ct != "" && ct != "application/pdf" && !strings.HasPrefix(ct, "application/pdf;") {
		return nil, errors.New("INVALID_FILE_TYPE", "content type must be application/pdf", errors.ErrInvalidInput.Status)
	}
	if ct == "" {
		ct = "application/pdf"
	}

	if err := os.MkdirAll(s.dir, 0o755); err != nil {
		return nil, errors.Wrap(err, errors.ErrInternal)
	}

	storagePath := s.storagePath(residenceID)
	tmpPath := storagePath + ".tmp"
	out, err := os.Create(tmpPath)
	if err != nil {
		return nil, errors.Wrap(err, errors.ErrInternal)
	}

	written, copyErr := io.Copy(out, io.LimitReader(body, maxLeasePDFBytes+1))
	closeErr := out.Close()
	if copyErr != nil {
		os.Remove(tmpPath)
		return nil, errors.Wrap(copyErr, errors.ErrInternal)
	}
	if closeErr != nil {
		os.Remove(tmpPath)
		return nil, errors.Wrap(closeErr, errors.ErrInternal)
	}
	if written == 0 {
		os.Remove(tmpPath)
		return nil, errors.New("INVALID_FILE", "file is empty", errors.ErrInvalidInput.Status)
	}
	if written > maxLeasePDFBytes {
		os.Remove(tmpPath)
		return nil, errors.New("FILE_TOO_LARGE", "lease PDF must be 10MB or smaller", errors.ErrInvalidInput.Status)
	}

	if err := os.Rename(tmpPath, storagePath); err != nil {
		os.Remove(tmpPath)
		return nil, errors.Wrap(err, errors.ErrInternal)
	}

	safeName := filepath.Base(fileName)
	if safeName == "." || safeName == "/" || safeName == "" {
		safeName = fmt.Sprintf("lease-%d.pdf", residenceID)
	}

	doc := &models.Document{
		ResidenceID: residenceID,
		FileName:    safeName,
		ContentType: "application/pdf",
		Size:        written,
	}
	if err := s.repo.Upsert(doc); err != nil {
		return nil, errors.Wrap(err, errors.ErrInternal)
	}
	return doc, nil
}

// OpenForDownload returns metadata and an open file handle for the lease PDF.
// Caller must close the returned file.
func (s *DocumentService) OpenForDownload(residenceID int) (*models.Document, *os.File, error) {
	doc, err := s.GetByResidenceID(residenceID)
	if err != nil {
		return nil, nil, err
	}
	f, err := os.Open(s.storagePath(residenceID))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil, errors.ErrNotFound
		}
		return nil, nil, errors.Wrap(err, errors.ErrInternal)
	}
	return doc, f, nil
}

func (s *DocumentService) storagePath(residenceID int) string {
	return filepath.Join(s.dir, fmt.Sprintf("%d.pdf", residenceID))
}
