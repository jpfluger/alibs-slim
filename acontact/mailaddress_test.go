package acontact

import (
	"github.com/bojanz/address"
	"testing"
)

// TestMailAddress_Validate tests the Validate method of the MailAddress struct.
func TestMailAddress_Validate(t *testing.T) {
	tests := []struct {
		mail    *MailAddress
		wantErr bool
		errMsg  string
	}{
		{
			mail: &MailAddress{
				Address: address.Address{
					Line1:       "123 Main St",
					Locality:    "Anytown",
					Region:      "CA",
					PostalCode:  "12345",
					CountryCode: "US",
				},
				Raw: "",
			},
			wantErr: false,
		},
		{
			mail: &MailAddress{
				Address: address.Address{},
				Raw:     "123 Main St, Anytown, CA 12345, US",
			},
			wantErr: false,
		},
		{
			mail: &MailAddress{
				Address: address.Address{},
				Raw:     "",
			},
			wantErr: true,
			errMsg:  "either structured address or raw address must be provided",
		},
		{
			mail: &MailAddress{
				Address: address.Address{
					Line1:       "123 Main St",
					Locality:    "Anytown",
					Region:      "CA",
					PostalCode:  "12345",
					CountryCode: "US",
				},
				Raw: "123 Main St, Anytown, CA 12345, US",
			},
			wantErr: true,
			errMsg:  "only one of structured address or raw address should be provided",
		},
	}

	for _, tt := range tests {
		err := tt.mail.Validate()
		if (err != nil) != tt.wantErr {
			t.Errorf("MailAddress.Validate() error = %v, wantErr %v", err, tt.wantErr)
		}
		if err != nil && tt.errMsg != "" && err.Error() != tt.errMsg {
			t.Errorf("MailAddress.Validate() error = %v, wantErrMsg %v", err, tt.errMsg)
		}
	}
}

// TestMailAddress_IsEmpty tests the IsEmpty method of the MailAddress struct.
func TestMailAddress_IsEmpty(t *testing.T) {
	mail := &MailAddress{}
	if !mail.IsEmpty() {
		t.Errorf("MailAddress.IsEmpty() should be true for an empty address")
	}

	mail.Line1 = "123 Main St"
	if mail.IsEmpty() {
		t.Errorf("MailAddress.IsEmpty() should be false when an address line is set")
	}
}

// TestMailAddress_GetVerify tests the GetVerify method of the MailAddress struct.
func TestMailAddress_GetVerify(t *testing.T) {
	mail := &MailAddress{IsVerified: true}
	if !mail.GetVerify() {
		t.Errorf("MailAddress.GetVerify() should return true when IsVerified is true")
	}

	mail.IsVerified = false
	if mail.GetVerify() {
		t.Errorf("MailAddress.GetVerify() should return false when IsVerified is false")
	}
}

// TestMailAddress_VerifyWithFields tests the VerifyWithFields method of the MailAddress struct.
func TestMailAddress_VerifyWithFields(t *testing.T) {
	mail := &MailAddress{
		Address: address.Address{
			Line1:       "123 Main St",
			Locality:    "Anytown",
			Region:      "CA",
			PostalCode:  "12345",
			CountryCode: "US",
		},
	}

	// Test with all required fields present.
	err := mail.VerifyWithFields("en", address.FieldLine1, address.FieldLocality, address.FieldRegion)
	if err != nil {
		t.Errorf("MailAddress.VerifyWithFields() should not return an error when all fields are present")
	}

	// Test with a missing required field.
	mail.Line1 = ""
	err = mail.VerifyWithFields("en", address.FieldLine1)
	if err == nil {
		t.Errorf("MailAddress.VerifyWithFields() should return an error when a required field is missing")
	}
}

// TestMailAddress_ToHTML tests the ToHTML method of the MailAddress struct.
func TestMailAddress_ToHTML(t *testing.T) {
	mail := &MailAddress{
		Address: address.Address{
			Line1:       "123 Main St",
			Locality:    "Anytown",
			Region:      "CA",
			PostalCode:  "12345",
			CountryCode: "US",
		},
	}

	html := mail.ToHTML("en")
	if html == "" {
		t.Errorf("MailAddress.ToHTML() should return a non-empty string")
	}
}
