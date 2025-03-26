package acron

import (
	"fmt"
	"github.com/gofrs/uuid/v5"
	googleUUID "github.com/google/uuid"
)

// ConvertGoogleUUIDToGOFRSID converts a Google UUID (github.com/google/uuid)
// to a Gofrs UUID (github.com/gofrs/uuid). It returns the converted Gofrs UUID
// or an error if the conversion fails.
func ConvertGoogleUUIDToGOFRSID(target googleUUID.UUID) (uuid.UUID, error) {
	// Convert the Google UUID to a byte array and then to a Gofrs UUID
	newid, err := uuid.FromBytes(target[:])
	if err != nil {
		// Return a Nil Gofrs UUID and an error if the conversion fails
		return uuid.Nil, fmt.Errorf("error converting Google UUID to Gofrs UUID: %v", err)
	}
	// Return the successfully converted Gofrs UUID
	return newid, nil
}

// ConvertGOFRSIDToGoogleUUID converts a Gofrs UUID (github.com/gofrs/uuid)
// to a Google UUID (github.com/google/uuid). It returns the converted Google UUID.
func ConvertGOFRSIDToGoogleUUID(target uuid.UUID) googleUUID.UUID {
	// Directly cast the Gofrs UUID to a Google UUID
	googleUUID := googleUUID.UUID(target)
	return googleUUID
}
