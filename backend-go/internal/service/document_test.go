package service

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"brickwise/backend/internal/errors"
	"brickwise/backend/internal/models"
)

type mockResidenceRepo struct {
	byID map[int]*models.Residence
}

func (m *mockResidenceRepo) Create(residence *models.Residence) error { return nil }
func (m *mockResidenceRepo) Update(residence *models.Residence) error { return nil }
func (m *mockResidenceRepo) GetAll(filter *models.ResidenceFilter) ([]*models.Residence, error) {
	return nil, nil
}
func (m *mockResidenceRepo) GetByID(id int) (*models.Residence, error) {
	r, ok := m.byID[id]
	if !ok {
		return nil, fmt.Errorf("residence with id %d not found", id)
	}
	return r, nil
}

type mockDocumentRepo struct {
	byResidence map[int]*models.Document
	nextID      int
}

func (m *mockDocumentRepo) GetByResidenceID(residenceID int) (*models.Document, error) {
	doc, ok := m.byResidence[residenceID]
	if !ok {
		return nil, fmt.Errorf("document for residence %d not found", residenceID)
	}
	copy := *doc
	return &copy, nil
}

func (m *mockDocumentRepo) Upsert(doc *models.Document) error {
	if m.nextID == 0 {
		m.nextID = 1
	}
	if existing, ok := m.byResidence[doc.ResidenceID]; ok {
		doc.ID = existing.ID
	} else {
		doc.ID = m.nextID
		m.nextID++
	}
	doc.CreatedAt = time.Now()
	copy := *doc
	m.byResidence[doc.ResidenceID] = &copy
	return nil
}

func (m *mockDocumentRepo) DeleteByResidenceID(residenceID int) error {
	delete(m.byResidence, residenceID)
	return nil
}

func TestDocumentService_UploadAndDownloadPDF(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("DOCUMENTS_DIR", dir)

	residences := &mockResidenceRepo{byID: map[int]*models.Residence{
		7: {ID: 7, Name: "Tenant"},
	}}
	docs := &mockDocumentRepo{byResidence: map[int]*models.Document{}}
	svc := NewDocumentService(docs, residences)

	pdf := []byte("%PDF-1.4 lease content")
	doc, err := svc.UploadLeasePDF(7, "lease.pdf", "application/pdf", int64(len(pdf)), bytes.NewReader(pdf))
	if err != nil {
		t.Fatalf("upload failed: %v", err)
	}
	if doc.FileName != "lease.pdf" {
		t.Fatalf("expected file name lease.pdf, got %s", doc.FileName)
	}

	stored := filepath.Join(dir, "7.pdf")
	if _, err := os.Stat(stored); err != nil {
		t.Fatalf("expected file on disk: %v", err)
	}

	// Replace on re-upload
	pdf2 := []byte("%PDF-1.4 replaced")
	doc2, err := svc.UploadLeasePDF(7, "updated-lease.pdf", "application/pdf", int64(len(pdf2)), bytes.NewReader(pdf2))
	if err != nil {
		t.Fatalf("replace upload failed: %v", err)
	}
	if doc2.FileName != "updated-lease.pdf" {
		t.Fatalf("expected updated file name, got %s", doc2.FileName)
	}
	if len(docs.byResidence) != 1 {
		t.Fatalf("expected one document metadata row, got %d", len(docs.byResidence))
	}

	got, f, err := svc.OpenForDownload(7)
	if err != nil {
		t.Fatalf("download open failed: %v", err)
	}
	defer f.Close()
	if got.FileName != "updated-lease.pdf" {
		t.Fatalf("download metadata mismatch: %s", got.FileName)
	}
}

func TestDocumentService_RejectsNonPDF(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("DOCUMENTS_DIR", dir)

	residences := &mockResidenceRepo{byID: map[int]*models.Residence{
		1: {ID: 1, Name: "Tenant"},
	}}
	svc := NewDocumentService(&mockDocumentRepo{byResidence: map[int]*models.Document{}}, residences)

	_, err := svc.UploadLeasePDF(1, "notes.txt", "text/plain", 4, strings.NewReader("note"))
	if err == nil {
		t.Fatal("expected non-PDF upload to fail")
	}
	apiErr, ok := err.(*errors.APIError)
	if !ok || apiErr.Code != "INVALID_FILE_TYPE" {
		t.Fatalf("expected INVALID_FILE_TYPE, got %v", err)
	}
}
