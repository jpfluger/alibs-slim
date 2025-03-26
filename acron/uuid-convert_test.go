package acron

import (
	"github.com/gofrs/uuid/v5"
	googleUUID "github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestConvertGoogleUUIDToGOFRSID tests the ConvertGoogleUUIDToGOFRSID function.
func TestConvertGoogleUUIDToGOFRSID(t *testing.T) {
	// Generate a new Google UUID
	googleID := googleUUID.New()

	// Convert the Google UUID to a Gofrs UUID
	gofrsID, err := ConvertGoogleUUIDToGOFRSID(googleID)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Convert the Gofrs UUID back to a Google UUID
	convertedGoogleID := ConvertGOFRSIDToGoogleUUID(gofrsID)

	// Check if the original Google UUID and the converted Google UUID are the same
	if googleID != convertedGoogleID {
		t.Errorf("Expected %v, got %v", googleID, convertedGoogleID)
	}
}

// TestConvertGOFRSIDToGoogleUUID tests the ConvertGOFRSIDToGoogleUUID function.
func TestConvertGOFRSIDToGoogleUUID(t *testing.T) {
	// Generate a new Gofrs UUID
	gofrsID, err := uuid.NewV4()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Convert the Gofrs UUID to a Google UUID
	googleID := ConvertGOFRSIDToGoogleUUID(gofrsID)

	// Convert the Google UUID back to a Gofrs UUID
	convertedGofrsID, err := ConvertGoogleUUIDToGOFRSID(googleID)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Check if the original Gofrs UUID and the converted Gofrs UUID are the same
	if gofrsID != convertedGofrsID {
		t.Errorf("Expected %v, got %v", gofrsID, convertedGofrsID)
	}
}

// Test conversion failure case when trying to convert invalid UUID
func TestConvertGoogleUUIDToGOFRSID_InvalidUUID(t *testing.T) {
	// Create an invalid Google UUID (zeroed)
	invalidUUID := googleUUID.UUID{}

	// Try to convert it to a Gofrs UUID
	newuuid, err := ConvertGoogleUUIDToGOFRSID(invalidUUID)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, uuid.Nil, newuuid, "The resulting Gofrs UUID should be Nil")
}
