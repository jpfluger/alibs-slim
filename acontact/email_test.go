package acontact

import (
	"fmt"
	"github.com/jpfluger/alibs-slim/aemail"
	"testing"
)

// TestEmail_Validate tests the Validate method of the Email type
func TestEmail_Validate(t *testing.T) {
	tests := []struct {
		email *Email
		want  error
	}{
		{&Email{Type: EMAILTYPE_WORK, Address: aemail.EmailAddress("work@example.com")}, nil},
		{&Email{Type: "", Address: aemail.EmailAddress("invalid@example")}, fmt.Errorf("email type is empty")},
		{&Email{Type: EMAILTYPE_PERSONAL, Address: aemail.EmailAddress("")}, fmt.Errorf("email address is empty")},
	}

	for _, tt := range tests {
		err := tt.email.Validate()
		if (err != nil) != (tt.want != nil) {
			t.Errorf("Email.Validate() error = %v, wantErr %v", err, tt.want)
		}
		if err != nil && tt.want != nil && err.Error() != tt.want.Error() {
			t.Errorf("Email.Validate() error = %v, wantErr %v", err, tt.want)
		}
	}
}

// TestEmails_FindByType tests the FindByType method of the Emails type
func TestEmails_FindByType(t *testing.T) {
	emails := Emails{
		&Email{Type: EMAILTYPE_WORK, Address: aemail.EmailAddress("work@example.com")},
		&Email{Type: EMAILTYPE_PERSONAL, Address: aemail.EmailAddress("home@example.com")},
	}

	if got := emails.FindByType(EMAILTYPE_WORK); got.Address != "work@example.com" {
		t.Errorf("Emails.FindByType() = %v, want %v", got.Address, "work@example.com")
	}
	if got := emails.FindByType(EMAILTYPE_PERSONAL); got.Address != "home@example.com" {
		t.Errorf("Emails.FindByType() = %v, want %v", got.Address, "home@example.com")
	}
	if got := emails.FindByType(EMAILTYPE_NO_REPLY); got != nil {
		t.Errorf("Emails.FindByType() = %v, want %v", got, nil)
	}
}

// TestEmails_FindByTypeOrDefault tests the FindByTypeOrDefault method of the Emails type
func TestEmails_FindByTypeOrDefault(t *testing.T) {
	emails := Emails{
		&Email{Type: EMAILTYPE_WORK, Address: aemail.EmailAddress("work@example.com"), IsDefault: true},
		&Email{Type: EMAILTYPE_PERSONAL, Address: aemail.EmailAddress("home@example.com")},
	}

	if got := emails.FindByTypeOrDefault(EMAILTYPE_NO_REPLY); got.Address != "work@example.com" {
		t.Errorf("Emails.FindByTypeOrDefault() = %v, want %v", got.Address, "work@example.com")
	}
}

// TestEmails_HasType tests the HasType method of the Emails type
func TestEmails_HasType(t *testing.T) {
	emails := Emails{
		&Email{Type: EMAILTYPE_WORK, Address: aemail.EmailAddress("work@example.com")},
	}

	if !emails.HasType(EMAILTYPE_WORK) {
		t.Errorf("Emails.HasType() should be true, got false")
	}
}

// TestEmails_Clone tests the Clone method of the Emails type
func TestEmails_Clone(t *testing.T) {
	emails := Emails{
		&Email{Type: EMAILTYPE_WORK, Address: aemail.EmailAddress("work@example.com")},
	}
	clonedEmails := emails.Clone()

	if len(clonedEmails) != 1 || clonedEmails[0].Address != "work@example.com" {
		t.Errorf("Emails.Clone() failed to clone emails")
	}
}
