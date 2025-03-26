package acontact

import (
	"fmt"
	"strings"
	"testing"
)

// TestPhone_Validate tests the Validate method of the Phone type
func TestPhone_Validate(t *testing.T) {
	tests := []struct {
		phone *Phone
		want  error
	}{
		{&Phone{Type: PHONETYPE_MOBILE, Number: "+12025550100", Country: "US"}, nil},
		{&Phone{Type: "", Number: "+12025550100", Country: "US"}, fmt.Errorf("phone type is empty")},
		{&Phone{Type: PHONETYPE_MOBILE, Number: "2025550100", Country: "US"}, nil},
		{&Phone{Type: PHONETYPE_MOBILE, Number: "INVALID", Country: "US"}, fmt.Errorf("invalid phone number:")},
		{&Phone{Type: PHONETYPE_MOBILE, Number: "+12025550100", Country: ""}, nil},
	}

	for _, tt := range tests {
		err := tt.phone.Validate()
		if (err != nil) != (tt.want != nil) {
			t.Errorf("Phone.Validate() error = %v, wantErr %v", err, tt.want)
		}
		if err != nil && tt.want != nil && !strings.Contains(err.Error(), tt.want.Error()) {
			t.Errorf("Phone.Validate() error = %v, wantErr %v", err, tt.want)
		}
	}
}

// TestPhones_FindByType tests the FindByType method of the Phones type
func TestPhones_FindByType(t *testing.T) {
	phones := Phones{
		&Phone{Type: PHONETYPE_MOBILE, Number: "+12025550100", Country: "US"},
		&Phone{Type: PHONETYPE_HOME, Number: "+12025550101", Country: "US"},
	}

	if got := phones.FindByType(PHONETYPE_MOBILE); got.Number != "+12025550100" {
		t.Errorf("Phones.FindByType() = %v, want %v", got.Number, "+12025550100")
	}
	if got := phones.FindByType(PHONETYPE_HOME); got.Number != "+12025550101" {
		t.Errorf("Phones.FindByType() = %v, want %v", got.Number, "+12025550101")
	}
	if got := phones.FindByType(PHONETYPE_WORK); got != nil {
		t.Errorf("Phones.FindByType() = %v, want %v", got, nil)
	}
}

// TestPhone_Validate_UnknownCountry tests the Validate method of the Phone type when the country is unknown
func TestPhone_Validate_UnknownCountry(t *testing.T) {
	tests := []struct {
		phone *Phone
		want  error
	}{
		// Assuming "+12025550100" is a valid number in the default region set in the library
		{&Phone{Type: PHONETYPE_MOBILE, Number: "+12025550100", Country: ""}, nil},
		// Assuming "2025550100" requires a country code to be valid
		{&Phone{Type: PHONETYPE_MOBILE, Number: "2025550100", Country: ""}, fmt.Errorf("invalid phone number:")},
		// Assuming "5550100" is invalid number
		{&Phone{Type: PHONETYPE_MOBILE, Number: "5550100", Country: ""}, fmt.Errorf("invalid phone number:")},
		// Invalid number without a country code should still result in an error
		{&Phone{Type: PHONETYPE_MOBILE, Number: "INVALID", Country: ""}, fmt.Errorf("invalid phone number:")},
	}

	for _, tt := range tests {
		err := tt.phone.Validate()
		if (err != nil) != (tt.want != nil) {
			t.Errorf("Phone.Validate() error = %v, wantErr %v", err, tt.want)
		}
		if err != nil && tt.want != nil && !strings.Contains(err.Error(), tt.want.Error()) {
			t.Errorf("Phone.Validate() error = %v, wantErr %v", err, tt.want)
		}
	}
}
