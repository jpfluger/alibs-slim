package acontact

import (
	"strings"
	// Importing custom utilities package for string manipulation
	"github.com/jpfluger/alibs-slim/autils"
)

// EntityType defines a custom type for contact entities
type EntityType string

const (
	ENTITYTYPE_ENTITY      = EntityType("entity")      // Represents a general entity
	ENTITYTYPE_COMPANY     = EntityType("company")     // Represents a company contact
	ENTITYTYPE_SUPPLIER    = EntityType("supplier")    // Represents a supplier
	ENTITYTYPE_PARTNER     = EntityType("partner")     // Represents a business partner
	ENTITYTYPE_SERVICE     = EntityType("service")     // Represents a service provider
	ENTITYTYPE_SUPPORT     = EntityType("support")     // Represents a support contact
	ENTITYTYPE_DISTRIBUTOR = EntityType("distributor") // Represents a distributor
	ENTITYTYPE_CONTRACTOR  = EntityType("contractor")  // Represents a contractor
	ENTITYTYPE_CONSULTANT  = EntityType("consultant")  // Represents a consultant
	ENTITYTYPE_STAKEHOLDER = EntityType("stakeholder") // Represents a stakeholder
)

// IsEmpty checks if the EntityType is empty after trimming spaces
func (ct EntityType) IsEmpty() bool {
	return strings.TrimSpace(string(ct)) == ""
}

// TrimSpace returns a new EntityType with leading and trailing spaces removed
func (ct EntityType) TrimSpace() EntityType {
	return EntityType(strings.TrimSpace(string(ct)))
}

// HasMatch checks if the target EntityType matches the current EntityType
func (ct EntityType) HasMatch(target EntityType) bool {
	return ct == target
}

// String returns the string representation of EntityType in trimmed and lowercase format
func (ct EntityType) String() string {
	return ct.ToStringTrimLower()
}

// ToStringTrimLower returns the EntityType as a trimmed and lowercase string
func (ct EntityType) ToStringTrimLower() string {
	return autils.ToStringTrimLower(string(ct))
}

// EntityTypes defines a slice of EntityType
type EntityTypes []EntityType
