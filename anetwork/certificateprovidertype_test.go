package anetwork

import (
	"testing"
)

func TestCertificateProviderType_IsEmpty(t *testing.T) {
	tests := []struct {
		name string
		cpt  CertificateProviderType
		want bool
	}{
		{"EmptyString", CertificateProviderType(""), true},
		{"SpacesOnly", CertificateProviderType("   "), true},
		{"NonEmptyString", CertificateProviderType("example"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.cpt.IsEmpty(); got != tt.want {
				t.Errorf("IsEmpty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCertificateProviderType_TrimSpace(t *testing.T) {
	tests := []struct {
		name string
		cpt  CertificateProviderType
		want CertificateProviderType
	}{
		{"NoSpaces", CertificateProviderType("example"), CertificateProviderType("example")},
		{"LeadingSpaces", CertificateProviderType("  example"), CertificateProviderType("example")},
		{"TrailingSpaces", CertificateProviderType("example  "), CertificateProviderType("example")},
		{"LeadingAndTrailingSpaces", CertificateProviderType("  example  "), CertificateProviderType("example")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.cpt.TrimSpace(); got != tt.want {
				t.Errorf("TrimSpace() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCertificateProviderType_String(t *testing.T) {
	tests := []struct {
		name string
		cpt  CertificateProviderType
		want string
	}{
		{"NonEmptyString", CertificateProviderType("example"), "example"},
		{"EmptyString", CertificateProviderType(""), ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.cpt.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCertificateProviderType_ToStringTrimLower(t *testing.T) {
	tests := []struct {
		name string
		cpt  CertificateProviderType
		want string
	}{
		{"MixedCaseWithSpaces", CertificateProviderType("  ExAmPlE  "), "example"},
		{"LowerCase", CertificateProviderType("example"), "example"},
		{"UpperCase", CertificateProviderType("EXAMPLE"), "example"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.cpt.ToStringTrimLower(); got != tt.want {
				t.Errorf("ToStringTrimLower() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCertificateProviderTypes_Contains(t *testing.T) {
	tests := []struct {
		name string
		cpts CertificateProviderTypes
		cpt  CertificateProviderType
		want bool
	}{
		{"Contains", CertificateProviderTypes{"example", "test"}, CertificateProviderType("test"), true},
		{"DoesNotContain", CertificateProviderTypes{"example", "test"}, CertificateProviderType("notfound"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.cpts.Contains(tt.cpt); got != tt.want {
				t.Errorf("Contains() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCertificateProviderTypes_Add(t *testing.T) {
	tests := []struct {
		name string
		cpts CertificateProviderTypes
		cpt  CertificateProviderType
		want CertificateProviderTypes
	}{
		{"AddNew", CertificateProviderTypes{"example"}, CertificateProviderType("test"), CertificateProviderTypes{"example", "test"}},
		{"AddExisting", CertificateProviderTypes{"example", "test"}, CertificateProviderType("test"), CertificateProviderTypes{"example", "test"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpts := tt.cpts
			cpts.Add(tt.cpt)
			if !equalCertificateProviderTypes(cpts, tt.want) {
				t.Errorf("Add() = %v, want %v", cpts, tt.want)
			}
		})
	}
}

func TestCertificateProviderTypes_Remove(t *testing.T) {
	tests := []struct {
		name string
		cpts CertificateProviderTypes
		cpt  CertificateProviderType
		want CertificateProviderTypes
	}{
		{"RemoveExisting", CertificateProviderTypes{"example", "test"}, CertificateProviderType("test"), CertificateProviderTypes{"example"}},
		{"RemoveNonExisting", CertificateProviderTypes{"example"}, CertificateProviderType("test"), CertificateProviderTypes{"example"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpts := tt.cpts
			cpts.Remove(tt.cpt)
			if !equalCertificateProviderTypes(cpts, tt.want) {
				t.Errorf("Remove() = %v, want %v", cpts, tt.want)
			}
		})
	}
}

// Helper function to compare two CertificateProviderTypes slices
func equalCertificateProviderTypes(a, b CertificateProviderTypes) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
