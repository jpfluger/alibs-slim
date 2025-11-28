package aemail

import (
	"fmt"
	"net/mail"
	"strings"
)

// Address represents a single mail address.
// An address such as "Barry Gibbs <bg@example.com>" is represented
// as Address{Name: "Barry Gibbs", Address: "bg@example.com"}.
type Address struct {
	Name    string       `json:"name,omitempty"`    // Proper name; may be empty.
	Address EmailAddress `json:"address,omitempty"` // user@domain
}

// Validate validates the Address by trimming whitespace from Name and Address,
// and then checking if the Address is a valid email. It mutates the receiver
// to apply the trimming. Returns an error if the Address is nil or invalid.
func (a *Address) Validate() error {
	if a == nil {
		return fmt.Errorf("address is nil")
	}
	a.Name = strings.TrimSpace(a.Name)
	a.Address = a.Address.TrimSpace()
	if err := a.Address.Validate(); err != nil {
		return err
	}
	return nil
}

// ToMailAddress converts the Address to a net/mail.Address if valid.
// It first checks if the Address is non-nil and valid (without mutating).
// Returns the converted mail.Address and a boolean indicating success.
func (a *Address) ToMailAddress() (addr mail.Address, ok bool) {
	if a == nil {
		return mail.Address{}, false
	}
	if err := a.Address.Validate(); err != nil {
		return mail.Address{}, false
	}
	return mail.Address{
		Name:    a.Name,
		Address: a.Address.String(),
	}, true
}

// String returns the string representation of the Address in RFC 5322 format.
// If the Address is nil, it returns an empty string.
func (a *Address) String() string {
	if a == nil {
		return ""
	}
	addr := mail.Address{
		Name:    a.Name,
		Address: a.Address.String(),
	}
	return addr.String()
}

// Addresses is a slice of Address, providing methods for validation and conversion.
type Addresses []Address

// Len returns the number of addresses in the slice.
func (as Addresses) Len() int {
	if as == nil {
		return 0
	}
	return len(as)
}

// Validate checks each Address in the slice for validity.
func (as Addresses) Validate() error {
	if as.Len() == 0 {
		return nil
	}
	for ii, a := range as {
		if err := a.Validate(); err != nil {
			return fmt.Errorf("invalid address at index %d: %v", ii, err)
		}
	}
	return nil
}

// ToMailAddresses converts the Addresses slice to a slice of net/mail.Address.
func (as Addresses) ToMailAddresses() []mail.Address {
	ms := []mail.Address{}
	if as.Len() == 0 {
		return ms
	}
	for _, a := range as {
		m, ok := a.ToMailAddress()
		if !ok {
			continue
		}
		ms = append(ms, m)
	}
	return ms
}

// Clone creates a deep copy of the Addresses slice.
// It returns nil if the original is nil; otherwise, it allocates a new slice
// with the same length and copies the elements (which are value types, so no further nesting is needed).
// This prevents shared mutation issues by ensuring a new backing array.
func (as Addresses) Clone() Addresses {
	if as == nil {
		return nil
	}
	cloned := make(Addresses, len(as))
	copy(cloned, as)
	return cloned
}
