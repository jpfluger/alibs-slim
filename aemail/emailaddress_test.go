package aemail

import (
	"testing"
)

// TestEmailAddress_IsEmpty checks if IsEmpty correctly identifies empty email addresses.
func TestEmailAddress_IsEmpty(t *testing.T) {
	tests := []struct {
		email EmailAddress
		want  bool
	}{
		{"", true},
		{" ", true},
		{"test@example.com", false},
	}

	for _, tt := range tests {
		if got := tt.email.IsEmpty(); got != tt.want {
			t.Errorf("EmailAddress.IsEmpty() = %v, want %v", got, tt.want)
		}
	}
}

// TestEmailAddress_TrimSpace checks if TrimSpace correctly trims whitespace from email addresses.
func TestEmailAddress_TrimSpace(t *testing.T) {
	email := EmailAddress(" test@example.com ")
	want := EmailAddress("test@example.com")
	if got := email.TrimSpace(); got != want {
		t.Errorf("EmailAddress.TrimSpace() = %v, want %v", got, want)
	}
}

// TestEmailAddress_Name checks if Name correctly extracts the username part of the email address.
func TestEmailAddress_Name(t *testing.T) {
	email := EmailAddress("test@example.com")
	want := "test"
	if got := email.Name(); got != want {
		t.Errorf("EmailAddress.Name() = %v, want %v", got, want)
	}
}

// TestEmailAddress_Domain checks if Domain correctly extracts the domain part of the email address.
func TestEmailAddress_Domain(t *testing.T) {
	email := EmailAddress("test@example.com")
	want := "example.com"
	if got := email.Domain(); got != want {
		t.Errorf("EmailAddress.Domain() = %v, want %v", got, want)
	}
}

// TestEmailAddress_HasPrefix checks if HasPrefix correctly identifies email addresses with a specific prefix.
func TestEmailAddress_HasPrefix(t *testing.T) {
	email := EmailAddress("test@example.com")
	prefix := "test"
	if !email.HasPrefix(prefix) {
		t.Errorf("EmailAddress.HasPrefix() = false, want true")
	}
}

// TestEmailAddress_IsValid checks if IsValid correctly identifies valid email addresses.
func TestEmailAddress_IsValid(t *testing.T) {
	email := EmailAddress("test@example.com")
	if !email.IsValid() {
		t.Errorf("EmailAddress.IsValid() = false, want true")
	}
}

// TestEmailAddress_Validate checks if Validate correctly identifies invalid email addresses.
func TestEmailAddress_Validate(t *testing.T) {
	email := EmailAddress("test@")
	if email.Validate() == nil {
		t.Errorf("EmailAddress.Validate() = nil, want error")
	}
}

// TestEmailAddresses_HasValues checks if HasValues correctly identifies non-empty slices of email addresses.
func TestEmailAddresses_HasValues(t *testing.T) {
	emails := EmailAddresses{EmailAddress("test@example.com")}
	if !emails.HasValues() {
		t.Errorf("EmailAddresses.HasValues() = false, want true")
	}
}

// TestEmailAddresses_HasMatch checks if HasMatch correctly identifies matching email addresses.
func TestEmailAddresses_HasMatch(t *testing.T) {
	emails := EmailAddresses{EmailAddress("test@example.com")}
	match := EmailAddress("test@example.com")
	if !emails.HasMatch(match) {
		t.Errorf("EmailAddresses.HasMatch() = false, want true")
	}
}

// TestEmailAddresses_HasPrefix checks if HasPrefix correctly identifies email addresses with a specific prefix.
func TestEmailAddresses_HasPrefix(t *testing.T) {
	emails := EmailAddresses{EmailAddress("test@example.com")}
	prefix := "test"
	if !emails.HasPrefix(prefix) {
		t.Errorf("EmailAddresses.HasPrefix() = false, want true")
	}
}

// TestEmailAddresses_Clean checks if Clean correctly removes empty email addresses from the slice.
func TestEmailAddresses_Clean(t *testing.T) {
	emails := EmailAddresses{EmailAddress("test@example.com"), EmailAddress("")}
	cleaned := emails.Clean()
	if len(cleaned) != 1 {
		t.Errorf("EmailAddresses.Clean() = %v, want %v", len(cleaned), 1)
	}
}

// TestEmailAddresses_ToArrStrings checks if ToArrStrings correctly converts email addresses to a slice of strings.
func TestEmailAddresses_ToArrStrings(t *testing.T) {
	emails := EmailAddresses{EmailAddress("test@example.com")}
	want := []string{"test@example.com"}
	if got := emails.ToArrStrings(); len(got) != 1 || got[0] != want[0] {
		t.Errorf("EmailAddresses.ToArrStrings() = %v, want %v", got, want)
	}
}

// TestEmailAddresses_IncludeIfInTargets checks if IncludeIfInTargets correctly filters email addresses.
func TestEmailAddresses_IncludeIfInTargets(t *testing.T) {
	emails := EmailAddresses{EmailAddress("test@example.com"), EmailAddress("user@example.com")}
	targets := EmailAddresses{EmailAddress("test@example.com")}
	included := emails.IncludeIfInTargets(targets)
	if len(included) != 1 || included[0] != "test@example.com" {
		t.Errorf("EmailAddresses.IncludeIfInTargets() = %v, want %v", included, targets)
	}
}
