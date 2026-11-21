package validator

import "testing"

func TestValidateProperty_AcceptsMixed(t *testing.T) {
	err := ValidateProperty("Mixed Plaza", "100 Main St", "mixed")
	if err != nil {
		t.Fatalf("expected mixed property type to be accepted, got %v", err)
	}
}

func TestValidateProperty_AcceptsAllKnownTypes(t *testing.T) {
	types := []string{"residential", "commercial", "industrial", "land", "mixed", "MIXED"}
	for _, pt := range types {
		if err := ValidateProperty("Test", "123 St", pt); err != nil {
			t.Errorf("expected property_type %q to be accepted, got %v", pt, err)
		}
	}
}

func TestValidateProperty_RejectsUnknownType(t *testing.T) {
	err := ValidateProperty("Test", "123 St", "warehouse")
	if err == nil {
		t.Fatal("expected unknown property_type to be rejected")
	}
}
