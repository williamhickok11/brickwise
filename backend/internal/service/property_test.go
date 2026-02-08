package service

import (
	"fmt"
	"testing"

	"brickwise/backend/internal/models"
)

// mockPropertyRepository is a simple in-memory repository for testing
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

func TestPropertyService_CreateProperty(t *testing.T) {
	repo := newMockPropertyRepository()
	service := NewPropertyService(repo)

	req := &models.CreatePropertyRequest{
		Name:         "Test Property",
		Address:      "123 Test St",
		PropertyType: "residential",
	}

	property, err := service.CreateProperty(req)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if property.Name != req.Name {
		t.Errorf("Expected name %s, got %s", req.Name, property.Name)
	}

	if property.ID == 0 {
		t.Error("Expected property to have an ID")
	}
}

func TestPropertyService_CreateProperty_Validation(t *testing.T) {
	repo := newMockPropertyRepository()
	service := NewPropertyService(repo)

	req := &models.CreatePropertyRequest{
		Name:         "", // Invalid: empty name
		Address:      "123 Test St",
		PropertyType: "residential",
	}

	_, err := service.CreateProperty(req)
	if err == nil {
		t.Fatal("Expected validation error")
	}

	apiErr, ok := err.(*errors.APIError)
	if !ok {
		t.Fatal("Expected APIError")
	}

	if apiErr.Status != errors.ErrInvalidInput.Status {
		t.Errorf("Expected status %d, got %d", errors.ErrInvalidInput.Status, apiErr.Status)
	}
}

func TestPropertyService_GetProperty(t *testing.T) {
	repo := newMockPropertyRepository()
	service := NewPropertyService(repo)

	// Create a property first
	req := &models.CreatePropertyRequest{
		Name:         "Test Property",
		Address:      "123 Test St",
		PropertyType: "residential",
	}
	created, _ := service.CreateProperty(req)

	// Retrieve it
	property, err := service.GetProperty(created.ID)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if property.ID != created.ID {
		t.Errorf("Expected ID %d, got %d", created.ID, property.ID)
	}
}

func TestPropertyService_GetProperty_NotFound(t *testing.T) {
	repo := newMockPropertyRepository()
	service := NewPropertyService(repo)

	_, err := service.GetProperty(999)
	if err == nil {
		t.Fatal("Expected error")
	}

	apiErr, ok := err.(*errors.APIError)
	if !ok {
		t.Fatal("Expected APIError")
	}

	if apiErr.Status != errors.ErrNotFound.Status {
		t.Errorf("Expected status %d, got %d", errors.ErrNotFound.Status, apiErr.Status)
	}
}
