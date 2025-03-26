package aemail

import (
	"fmt"
	"net/mail"
	"strings"
)

// EmailAddress represents an email address as a string.
type EmailAddress string

// IsEmpty checks if the EmailAddress is empty after trimming whitespace.
func (e EmailAddress) IsEmpty() bool {
	return strings.TrimSpace(string(e)) == ""
}

// TrimSpace returns a new EmailAddress with leading and trailing whitespace removed.
func (e EmailAddress) TrimSpace() EmailAddress {
	return EmailAddress(strings.TrimSpace(string(e)))
}

// Name extracts the username part of the EmailAddress before the "@".
func (e EmailAddress) Name() string {
	parts := strings.SplitN(string(e), "@", 2)
	if len(parts) == 2 {
		return parts[0]
	}
	return ""
}

// Domain extracts the domain part of the EmailAddress after the "@".
func (e EmailAddress) Domain() string {
	parts := strings.SplitN(string(e), "@", 2)
	if len(parts) == 2 {
		return parts[1]
	}
	return ""
}

// String returns the string representation of the EmailAddress.
func (e EmailAddress) String() string {
	return string(e)
}

// ToStringTrimLower converts the SecretCacheKey to a lowercase string with trimmed whitespace.
func (e EmailAddress) ToStringTrimLower() string {
	return strings.ToLower(strings.TrimSpace(string(e)))
}

// HasPrefix checks if the EmailAddress starts with the specified prefix.
func (e EmailAddress) HasPrefix(tPrefix string) bool {
	return strings.HasPrefix(e.String(), tPrefix)
}

// IsValid checks if the EmailAddress is valid by calling Validate.
func (e EmailAddress) IsValid() bool {
	return e.Validate() == nil
}

// Validate checks if the EmailAddress is valid according to the mail package.
func (e EmailAddress) Validate() error {
	_, err := e.ToRFC()
	return err
}

// ToRFC converts the EmailAddress to a *mail.Address, which validates it.
func (e EmailAddress) ToRFC() (*mail.Address, error) {
	if e.IsEmpty() {
		return nil, fmt.Errorf("undefined email address")
	}
	return mail.ParseAddress(e.String())
}

// EmailAddresses is a slice of EmailAddress, used to handle multiple email addresses.
type EmailAddresses []EmailAddress

// HasValues checks if the EmailAddresses slice contains any non-empty addresses.
func (es EmailAddresses) HasValues() bool {
	return len(es) > 0
}

// HasMatch checks if the EmailAddresses slice contains the specified EmailAddress.
func (es EmailAddresses) HasMatch(tEmail EmailAddress) bool {
	for _, e := range es {
		if e == tEmail {
			return true
		}
	}
	return false
}

// HasPrefix checks if any EmailAddress in the slice starts with the specified prefix.
func (es EmailAddresses) HasPrefix(tPrefix string) bool {
	for _, e := range es {
		if e.HasPrefix(tPrefix) {
			return true
		}
	}
	return false
}

// Clone creates a copy of the EmailAddresses slice.
func (es EmailAddresses) Clone() EmailAddresses {
	cloned := make(EmailAddresses, len(es))
	copy(cloned, es)
	return cloned
}

// ToArrStrings converts the EmailAddresses slice to a slice of strings.
func (es EmailAddresses) ToArrStrings() []string {
	strs := make([]string, len(es))
	for i, e := range es {
		strs[i] = e.String()
	}
	return strs
}

// IncludeIfInTargets returns a new EmailAddresses slice containing only the addresses that are in the 'targets' slice.
func (es EmailAddresses) IncludeIfInTargets(targets EmailAddresses) EmailAddresses {
	var included EmailAddresses
	for _, e := range es {
		if targets.HasMatch(e) {
			included = append(included, e)
		}
	}
	return included
}

// Clean returns a new EmailAddresses slice with empty addresses removed.
func (es EmailAddresses) Clean() EmailAddresses {
	var cleaned EmailAddresses
	for _, e := range es {
		if !e.IsEmpty() {
			cleaned = append(cleaned, e)
		}
	}
	return cleaned
}
