package alegal

import (
	"fmt"
	"strings"
	"sync"

	"github.com/jpfluger/alibs-slim/acontact"
)

// LegalOperator represents an operator in a legal context and embeds acontact.Company.
type LegalOperator struct {
	acontact.ContactCore // Embedding acontact.Entity to inherit its methods and properties.
	isValid              bool
}

func (lo *LegalOperator) IsValid() bool {
	return lo != nil && lo.isValid
}

// Validate the LegalOperator's required fields.
func (lo *LegalOperator) Validate() error {
	if lo == nil {
		return fmt.Errorf("legal operator is nil")
	}

	// Check required name
	name := strings.TrimSpace(lo.GetName())
	if name == "" {
		return fmt.Errorf("legal operator name is empty")
	}

	// Check required business URL
	if !lo.Urls.HasType(acontact.URLTYPE_BUSINESS) {
		return fmt.Errorf("legal operator missing 'business' URL")
	}
	business := lo.Urls.FindByType(acontact.URLTYPE_BUSINESS)
	if business.Link == nil || !business.Link.IsUrl() {
		return fmt.Errorf("business URL is missing or invalid")
	}

	// Optionally validate 'legal' URL if present
	if lo.Urls.HasType(acontact.URLTYPE_LEGAL) {
		legal := lo.Urls.FindByType(acontact.URLTYPE_LEGAL)
		if legal.Link == nil || !legal.Link.IsUrl() {
			return fmt.Errorf("legal URL is present but invalid")
		}
	}

	lo.isValid = true
	return nil
}

var appLegalOperator *LegalOperator
var appLegalOperatorView *LegalOperatorView
var muLegalOperator sync.RWMutex

func LEGALOPERATOR() *LegalOperator {
	muLegalOperator.RLock()
	defer muLegalOperator.RUnlock()
	return appLegalOperator
}

func SetLegalOperator(legalOperator *LegalOperator) {
	muLegalOperator.Lock()
	defer muLegalOperator.Unlock()
	appLegalOperator = legalOperator
	appLegalOperatorView = NewLegalOperatorView(legalOperator)
}

func GetLegalOperatorView() *LegalOperatorView {
	muLegalOperator.RLock()
	defer muLegalOperator.RUnlock()
	return appLegalOperatorView
}
