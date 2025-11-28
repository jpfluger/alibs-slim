package aemail

import (
	"net/mail"
	"reflect"
	"testing"
)

// TestAddress_Validate checks if Validate correctly trims whitespace, validates the email,
// and returns errors for nil or invalid addresses. It also verifies mutation.
func TestAddress_Validate(t *testing.T) {
	tests := []struct {
		name    string
		a       *Address
		wantErr bool
		want    *Address // Expected state after mutation
	}{
		{
			name:    "nil address",
			a:       nil,
			wantErr: true,
		},
		{
			name:    "invalid email",
			a:       &Address{Name: "Test", Address: "invalid@"},
			wantErr: true,
		},
		{
			name:    "valid with trimming",
			a:       &Address{Name: " Test User ", Address: " test@example.com "},
			wantErr: false,
			want:    &Address{Name: "Test User", Address: "test@example.com"},
		},
		{
			name:    "valid no trimming needed",
			a:       &Address{Name: "Test User", Address: "test@example.com"},
			wantErr: false,
			want:    &Address{Name: "Test User", Address: "test@example.com"},
		},
		{
			name:    "empty name is ok",
			a:       &Address{Name: "", Address: "test@example.com"},
			wantErr: false,
			want:    &Address{Name: "", Address: "test@example.com"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.a.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Address.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.want != nil && !reflect.DeepEqual(tt.a, tt.want) {
				t.Errorf("Address.Validate() mutated = %+v, want %+v", tt.a, tt.want)
			}
		})
	}
}

// TestAddress_ToMailAddress checks if ToMailAddress correctly converts valid addresses
// to net/mail.Address and returns false for nil or invalid.
func TestAddress_ToMailAddress(t *testing.T) {
	tests := []struct {
		name string
		a    *Address
		want mail.Address
		ok   bool
	}{
		{
			name: "nil address",
			a:    nil,
			want: mail.Address{},
			ok:   false,
		},
		{
			name: "invalid email",
			a:    &Address{Name: "Test", Address: "invalid@"},
			want: mail.Address{},
			ok:   false,
		},
		{
			name: "valid with name",
			a:    &Address{Name: "Test User", Address: "test@example.com"},
			want: mail.Address{Name: "Test User", Address: "test@example.com"},
			ok:   true,
		},
		{
			name: "valid without name",
			a:    &Address{Name: "", Address: "test@example.com"},
			want: mail.Address{Name: "", Address: "test@example.com"},
			ok:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := tt.a.ToMailAddress()
			if ok != tt.ok {
				t.Errorf("Address.ToMailAddress() ok = %v, want %v", ok, tt.ok)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Address.ToMailAddress() got = %+v, want %+v", got, tt.want)
			}
		})
	}
}

// TestAddress_String checks if String returns the correct RFC 5322 format or empty for nil.
func TestAddress_String(t *testing.T) {
	tests := []struct {
		name string
		a    *Address
		want string
	}{
		{
			name: "nil address",
			a:    nil,
			want: "",
		},
		{
			name: "with name",
			a:    &Address{Name: "Test User", Address: "test@example.com"},
			want: "\"Test User\" <test@example.com>",
		},
		{
			name: "without name",
			a:    &Address{Name: "", Address: "test@example.com"},
			want: "<test@example.com>",
		},
		{
			name: "with special characters in name",
			a:    &Address{Name: `Test "User"`, Address: "test@example.com"},
			want: `"Test \"User\"" <test@example.com>`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.a.String(); got != tt.want {
				t.Errorf("Address.String() = %q, want %q", got, tt.want)
			}
		})
	}
}

// TestAddresses_Len checks if Len correctly returns the length for nil, empty, and populated slices.
func TestAddresses_Len(t *testing.T) {
	tests := []struct {
		name string
		as   Addresses
		want int
	}{
		{
			name: "nil",
			as:   nil,
			want: 0,
		},
		{
			name: "empty",
			as:   Addresses{},
			want: 0,
		},
		{
			name: "one element",
			as:   Addresses{{}},
			want: 1,
		},
		{
			name: "multiple elements",
			as:   Addresses{{}, {}},
			want: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.as.Len(); got != tt.want {
				t.Errorf("Addresses.Len() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestAddresses_Validate checks if Validate returns errors for slices with invalid addresses.
func TestAddresses_Validate(t *testing.T) {
	tests := []struct {
		name    string
		as      Addresses
		wantErr bool
	}{
		{
			name:    "nil",
			as:      nil,
			wantErr: false,
		},
		{
			name:    "empty",
			as:      Addresses{},
			wantErr: false,
		},
		{
			name: "all valid",
			as: Addresses{
				{Name: "User1", Address: "user1@example.com"},
				{Name: "User2", Address: "user2@example.com"},
			},
			wantErr: false,
		},
		{
			name: "one invalid",
			as: Addresses{
				{Name: "User1", Address: "user1@example.com"},
				{Name: "User2", Address: "invalid@"},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.as.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Addresses.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestAddresses_ToMailAddresses checks if ToMailAddresses converts valid addresses and skips invalid ones.
func TestAddresses_ToMailAddresses(t *testing.T) {
	tests := []struct {
		name string
		as   Addresses
		want []mail.Address
	}{
		{
			name: "nil",
			as:   nil,
			want: []mail.Address{},
		},
		{
			name: "empty",
			as:   Addresses{},
			want: []mail.Address{},
		},
		{
			name: "all valid",
			as: Addresses{
				{Name: "User1", Address: "user1@example.com"},
				{Name: "User2", Address: "user2@example.com"},
			},
			want: []mail.Address{
				{Name: "User1", Address: "user1@example.com"},
				{Name: "User2", Address: "user2@example.com"},
			},
		},
		{
			name: "mixed valid and invalid",
			as: Addresses{
				{Name: "User1", Address: "user1@example.com"},
				{Name: "User2", Address: "invalid@"},
				{Name: "User3", Address: "user3@example.com"},
			},
			want: []mail.Address{
				{Name: "User1", Address: "user1@example.com"},
				{Name: "User3", Address: "user3@example.com"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.as.ToMailAddresses()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Addresses.ToMailAddresses() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

// TestAddresses_Clone checks if Clone returns a deep copy (new slice) and modifications don't affect original.
func TestAddresses_Clone(t *testing.T) {
	tests := []struct {
		name string
		as   Addresses
	}{
		{
			name: "nil",
			as:   nil,
		},
		{
			name: "empty",
			as:   Addresses{},
		},
		{
			name: "with elements",
			as: Addresses{
				{Name: "User1", Address: "user1@example.com"},
				{Name: "User2", Address: "user2@example.com"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cloned := tt.as.Clone()
			if tt.as == nil {
				if cloned != nil {
					t.Errorf("Addresses.Clone() = %+v, want nil", cloned)
				}
				return
			}
			if !reflect.DeepEqual(cloned, tt.as) {
				t.Errorf("Addresses.Clone() = %+v, want %+v", cloned, tt.as)
			}
			// Check if it's a new slice (different pointer)
			if len(tt.as) > 0 && &cloned[0] == &tt.as[0] {
				t.Errorf("Addresses.Clone() did not create a new backing array")
			}
			// Modify clone and ensure original unchanged
			if len(cloned) > 0 {
				cloned[0].Name = "Modified"
				if tt.as[0].Name == "Modified" {
					t.Errorf("Modification to clone affected original")
				}
			}
		})
	}
}
