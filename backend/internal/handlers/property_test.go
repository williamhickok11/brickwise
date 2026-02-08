package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"brickwise/backend/internal/models"
	"brickwise/backend/internal/repository"
	"brickwise/backend/internal/service"

	"github.com/go-chi/chi/v5"
)

// mockPropertyRepository implements repository.PropertyRepository for testing
type mockPropertyRepository struct {
	properties map[int]*models.Property
	nextID     int
}

func newMockPropertyRepository() *mockPropertyRepository {
	return &mockPropertyRepository{
		properties: make(map[int]*models.Property),
		nextID:     1,
	}
}

func (m *mockPropertyRepository) Create(property *models.Property) error {
	property.ID = m.nextID
	m.properties[m.nextID] = property
	m.nextID++
	return nil
}

func (m *mockPropertyRepository) GetByID(id int) (*models.Property, error) {
	property, ok := m.properties[id]
	if !ok {
		return nil, fmt.Errorf("property with id %d not found", id)
	}
	return property, nil
}

func (m *mockPropertyRepository) GetAll() ([]*models.Property, error) {
	var properties []*models.Property
	for _, p := range m.properties {
		properties = append(properties, p)
	}
	return properties, nil
}

func (m *mockPropertyRepository) Update(property *models.Property) error {
	if _, ok := m.properties[property.ID]; !ok {
		return fmt.Errorf("property with id %d not found", property.ID)
	}
	m.properties[property.ID] = property
	return nil
}

func (m *mockPropertyRepository) Delete(id int) error {
	if _, ok := m.properties[id]; !ok {
		return fmt.Errorf("property with id %d not found", id)
	}
	delete(m.properties, id)
	return nil
}

func TestCreateProperty(t *testing.T) {
	// Create a service with a mock repository for integration-style testing
	// In a real scenario, you might want to use a service interface for easier mocking
	mockRepo := &mockPropertyRepository{}
	svc := service.NewPropertyService(mockRepo)
	handler := NewPropertyHandler(svc)

	reqBody := models.CreatePropertyRequest{
		Name:         "Test Property",
		Address:      "123 Test St",
		PropertyType: "residential",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/properties", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.CreateProperty(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, w.Code)
	}

	var property models.Property
	if err := json.NewDecoder(w.Body).Decode(&property); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if property.Name != reqBody.Name {
		t.Errorf("Expected name %s, got %s", reqBody.Name, property.Name)
	}
}

func TestGetProperty(t *testing.T) {
	mockRepo := &mockPropertyRepository{}
	svc := service.NewPropertyService(mockRepo)
	
	// Create a property first
	req := &models.CreatePropertyRequest{
		Name:         "Test Property",
		Address:      "123 Test St",
		PropertyType: "residential",
	}
	property, _ := svc.CreateProperty(req)

	handler := NewPropertyHandler(svc)

	r := chi.NewRouter()
	r.Get("/properties/{id}", handler.GetProperty)

	req := httptest.NewRequest("GET", "/properties/1", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var result models.Property
	if err := json.NewDecoder(w.Body).Decode(&result); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if result.ID != property.ID {
		t.Errorf("Expected ID %d, got %d", property.ID, result.ID)
	}
	
	if result.Name != property.Name {
		t.Errorf("Expected name %s, got %s", property.Name, result.Name)
	}
}

func TestGetPropertyNotFound(t *testing.T) {
	mockRepo := &mockPropertyRepository{}
	svc := service.NewPropertyService(mockRepo)
	handler := NewPropertyHandler(svc)

	r := chi.NewRouter()
	r.Get("/properties/{id}", handler.GetProperty)

	req := httptest.NewRequest("GET", "/properties/999", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}
