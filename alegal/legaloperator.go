package alegal

import (
	"fmt"
	"github.com/jpfluger/alibs-slim/acontact"
)

// LegalOperator represents an operator in a legal context and embeds acontact.Company.
type LegalOperator struct {
	acontact.ContactCore // Embedding acontact.Entity to inherit its methods and properties.
}

// Validate the LegalOperator's required fields.
func (lo *LegalOperator) Validate() error {
	// Check if the LegalOperator is nil.
	if lo == nil {
		return fmt.Errorf("legal operator is nil")
	}

	// Check if the LegalOperator has a business URL.
	if !lo.Urls.HasType(acontact.URLTYPE_BUSINESS) {
		return fmt.Errorf("no url for 'business'")
	}

	// Check if the LegalOperator has a name set.
	if lo.GetName() == "" {
		return fmt.Errorf("name is empty")
	}

	// If all checks pass, return nil indicating successful initialization.
	return nil
}

var appLegalOperator *LegalOperator

func LEGALOPERATOR() *LegalOperator {
	return appLegalOperator
}

func SetLegalOperator(legalOperator *LegalOperator) {
	if appLegalOperator != nil {
		panic("appLegalOperator already initialized")
	}
	appLegalOperator = legalOperator
}
