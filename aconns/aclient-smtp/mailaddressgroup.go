package aclient_smtp

import (
	"fmt"
	"strings"

	"github.com/jpfluger/alibs-slim/aemail"
)

// MailAddressGroup represents a group of email addresses, including From, To, CC, and BCC fields.
type MailAddressGroup struct {
	From aemail.Address   `json:"from,omitempty"`
	To   aemail.Addresses `json:"to,omitempty"`
	CC   aemail.Addresses `json:"cc,omitempty"`
	BCC  aemail.Addresses `json:"bcc,omitempty"`
}

// HasFrom checks if the MailAddressGroup has a non-empty From address.
func (mag *MailAddressGroup) HasFrom() bool {
	if mag == nil {
		return false
	}
	return !mag.From.Address.IsEmpty()
}

// HasTo checks if the MailAddressGroup has any To addresses.
func (mag *MailAddressGroup) HasTo() bool {
	if mag == nil {
		return false
	}
	return len(mag.To) > 0
}

// HasCC checks if the MailAddressGroup has any CC addresses.
func (mag *MailAddressGroup) HasCC() bool {
	if mag == nil {
		return false
	}
	return len(mag.CC) > 0
}

// HasBCC checks if the MailAddressGroup has any BCC addresses.
func (mag *MailAddressGroup) HasBCC() bool {
	if mag == nil {
		return false
	}
	return len(mag.BCC) > 0
}

// HasAnyRecipients checks if the MailAddressGroup has any To, CC, or BCC addresses.
func (mag *MailAddressGroup) HasAnyRecipients() bool {
	if mag == nil {
		return false
	}
	return mag.HasTo() || mag.HasCC() || mag.HasBCC()
}

// Validate validates the MailAddressGroup.
// It requires a valid From address and validates To, CC, and BCC if present.
// It mutates the receiver by trimming fields via underlying Validate calls.
func (mag *MailAddressGroup) Validate() error {
	if mag == nil {
		return fmt.Errorf("mail address group is nil")
	}
	// Always require a "From" address.
	if err := mag.From.Validate(); err != nil {
		return fmt.Errorf("invalid from address: %w", err)
	}
	if mag.HasTo() {
		if err := mag.To.Validate(); err != nil {
			return fmt.Errorf("invalid to addresses: %w", err)
		}
	}
	if mag.HasCC() {
		if err := mag.CC.Validate(); err != nil {
			return fmt.Errorf("invalid cc addresses: %w", err)
		}
	}
	if mag.HasBCC() {
		if err := mag.BCC.Validate(); err != nil {
			return fmt.Errorf("invalid bcc addresses: %w", err)
		}
	}
	return nil
}

// MergeInto merges the non-empty properties of the source (mag) into the target (mergeTo)
// only if the target properties are empty.
func (mag *MailAddressGroup) MergeInto(mergeTo *MailAddressGroup) {
	// Merge From
	if !mag.From.Address.IsEmpty() && mergeTo.From.Address.IsEmpty() {
		mergeTo.From = mag.From
	}

	// Merge To
	if len(mergeTo.To) == 0 && len(mag.To) > 0 {
		mergeTo.To = mag.To
	}

	// Merge CC
	if len(mergeTo.CC) == 0 && len(mag.CC) > 0 {
		mergeTo.CC = mag.CC
	}

	// Merge BCC
	if len(mergeTo.BCC) == 0 && len(mag.BCC) > 0 {
		mergeTo.BCC = mag.BCC
	}
}

// MergeIntoTemplate merges non-empty values from a MailAddressGroup into a MailPiece.
// This ensures that only fields with data in the MailAddressGroup overwrite corresponding fields in MailPiece.
// Assumes MailPiece uses net/mail.Address types; conversions are applied.
func (mag *MailAddressGroup) MergeIntoTemplate(template *MailPiece) {
	if mag == nil {
		return
	}
	if !mag.From.Address.IsEmpty() {
		template.From = mag.From
	}

	var mustMerge bool
	if mag.To != nil && len(mag.To) > 0 {
		mustMerge = true
	} else if mag.CC != nil && len(mag.CC) > 0 {
		mustMerge = true
	} else if mag.BCC != nil && len(mag.BCC) > 0 {
		mustMerge = true
	}

	if mustMerge {
		template.To = mag.To.Clone()
		template.CC = mag.CC.Clone()
		template.BCC = mag.BCC.Clone()
	}
}

// FromMAG creates a new MailAddressGroup by merging defaultMAG into a deep copy of mergeMAG (or an empty one if nil).
// This avoids mutating the input and prevents shared slice issues via deep copying.
func FromMAG(defaultMAG *MailAddressGroup, mergeMAG *MailAddressGroup) *MailAddressGroup {
	var newMAG MailAddressGroup // Use value type for copy, then take address later.
	if mergeMAG != nil {
		newMAG = *mergeMAG // Shallow struct copy.
		// Deep copy slices to avoid sharing underlying arrays.
		newMAG.To = mergeMAG.To.Clone()
		newMAG.CC = mergeMAG.CC.Clone()
		newMAG.BCC = mergeMAG.BCC.Clone()
	}
	// Merge defaults into the copy (may reassign slices if empty).
	defaultMAG.MergeInto(&newMAG)
	return &newMAG
}

// NewMAG creates a new MailAddressGroup with all fields provided.
func NewMAG(from aemail.Address, to aemail.Addresses, cc aemail.Addresses, bcc aemail.Addresses) *MailAddressGroup {
	return &MailAddressGroup{
		From: from,
		To:   to,
		CC:   cc,
		BCC:  bcc,
	}
}

// NewMAGWithTo creates a new MailAddressGroup with a single "To" address.
func NewMAGWithTo(from aemail.Address, to aemail.Address) *MailAddressGroup {
	return &MailAddressGroup{
		From: from,
		To:   aemail.Addresses{to},
	}
}

// NewMAGWithCC creates a new MailAddressGroup with "From" and "CC" addresses.
func NewMAGWithCC(from aemail.Address, cc aemail.Addresses) *MailAddressGroup {
	return &MailAddressGroup{
		From: from,
		CC:   cc,
	}
}

// NewMAGWithBCC creates a new MailAddressGroup with "From" and "BCC" addresses.
func NewMAGWithBCC(from aemail.Address, bcc aemail.Addresses) *MailAddressGroup {
	return &MailAddressGroup{
		From: from,
		BCC:  bcc,
	}
}

// NewMAGWithToAndCC creates a new MailAddressGroup with a single "To" address and "CC" addresses.
func NewMAGWithToAndCC(from aemail.Address, to aemail.Address, cc aemail.Addresses) *MailAddressGroup {
	return &MailAddressGroup{
		From: from,
		To:   aemail.Addresses{to},
		CC:   cc,
	}
}

// NewMAGWithToAndBCC creates a new MailAddressGroup with a single "To" address and "BCC" addresses.
func NewMAGWithToAndBCC(from aemail.Address, to aemail.Address, bcc aemail.Addresses) *MailAddressGroup {
	return &MailAddressGroup{
		From: from,
		To:   aemail.Addresses{to},
		BCC:  bcc,
	}
}

// NewMAGWithSingleToAndCCBCC creates a new MailAddressGroup with a single "To", "CC", and "BCC" addresses.
func NewMAGWithSingleToAndCCBCC(from aemail.Address, to aemail.Address, cc aemail.Addresses, bcc aemail.Addresses) *MailAddressGroup {
	return &MailAddressGroup{
		From: from,
		To:   aemail.Addresses{to},
		CC:   cc,
		BCC:  bcc,
	}
}

// NewAddress creates a new aemail.Address with only the Address string.
// The Name field remains empty. Trims whitespace for consistency.
func NewAddress(email string) aemail.Address {
	return aemail.Address{
		Address: aemail.EmailAddress(strings.TrimSpace(email)),
	}
}

// NewMAGWithToString creates a MailAddressGroup with a single "To" email address string.
func NewMAGWithToString(from string, to string) *MailAddressGroup {
	return &MailAddressGroup{
		From: NewAddress(from),
		To:   aemail.Addresses{NewAddress(to)},
	}
}

// NewMAGWithCCString creates a MailAddressGroup with "From" and "CC" email address strings.
func NewMAGWithCCString(from string, cc []string) *MailAddressGroup {
	return &MailAddressGroup{
		From: NewAddress(from),
		CC:   convertStringsToAddresses(cc),
	}
}

// NewMAGWithBCCString creates a MailAddressGroup with "From" and "BCC" email address strings.
func NewMAGWithBCCString(from string, bcc []string) *MailAddressGroup {
	return &MailAddressGroup{
		From: NewAddress(from),
		BCC:  convertStringsToAddresses(bcc),
	}
}

// NewMAGWithToCCBCCString creates a MailAddressGroup with "To", "CC", and "BCC" email address strings.
func NewMAGWithToCCBCCString(from string, to string, cc []string, bcc []string) *MailAddressGroup {
	return &MailAddressGroup{
		From: NewAddress(from),
		To:   aemail.Addresses{NewAddress(to)},
		CC:   convertStringsToAddresses(cc),
		BCC:  convertStringsToAddresses(bcc),
	}
}

// Helper function to convert a slice of email strings into aemail.Addresses.
func convertStringsToAddresses(emails []string) aemail.Addresses {
	addresses := make(aemail.Addresses, len(emails))
	for i, email := range emails {
		addresses[i] = NewAddress(email)
	}
	return addresses
}

// MailAddressGroupMap is a map of MailAddressGroupKey to *MailAddressGroup.
// It provides methods for retrieval and validation.
type MailAddressGroupMap map[MailAddressGroupKey]*MailAddressGroup

// Get retrieves the MailAddressGroup for the given key, or nil if not found or map is nil.
func (mm MailAddressGroupMap) Get(key MailAddressGroupKey) *MailAddressGroup {
	if len(mm) == 0 {
		return nil
	}
	return mm[key.TrimSpace()]
}

// Validate validates each MailAddressGroup in the map.
// Returns nil if the map is nil or empty; otherwise, returns the first validation error encountered.
func (mm MailAddressGroupMap) Validate() error {
	if len(mm) == 0 {
		return nil
	}
	for key, mag := range mm {
		if err := mag.Validate(); err != nil {
			return fmt.Errorf("invalid mail address group for key %q: %w", key, err)
		}
	}
	return nil
}
