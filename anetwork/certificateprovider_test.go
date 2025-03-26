package anetwork

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestValidateDomains checks the ValidateDomains method.
func TestValidateDomains(t *testing.T) {
	tests := []struct {
		name       string
		domains    []string
		addDomains []string
		allowIPs   bool
		wantErr    bool
	}{
		{"ValidDomains", []string{"example.com", "*.example.com"}, []string{"www.example.com"}, false, false},
		{"InvalidDomain", []string{"example..com"}, []string{}, false, true},
		{"ValidIP", []string{"192.168.1.1"}, []string{}, true, false},
		{"InvalidIP", []string{"999.999.999.999"}, []string{}, true, false},
		{"AddDomainWithWildcard", []string{"example.com"}, []string{"*.example.com"}, false, true},
		{"AddDomainNotSatisfied", []string{"example.com"}, []string{"sub.example.com"}, false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpb := &CertificateProviderBase{Domains: tt.domains, AddDomains: tt.addDomains}
			err := cpb.ValidateDomains(tt.allowIPs)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateDomains() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestIsValidDomainWithError(t *testing.T) {
	tests := []struct {
		domain   string
		allowIPs bool
		expected bool
		errMsg   string
	}{
		// Test cases for empty domain
		{"", false, false, "domain is empty"},
		{"  ", false, false, "domain is empty"},

		// Test cases for IP addresses
		{"192.168.1.1", true, true, ""},
		{"192.168.1.1", false, false, "IP addresses are not allowed unless explicitly permitted"},
		{"::1", true, true, ""},
		{"::1", false, false, "IP addresses are not allowed unless explicitly permitted"},

		// Test cases for valid DNS names
		{"example.com", false, true, ""},
		{"sub.example.com", false, true, ""},
		{"*.example.com", false, true, ""},
		{"xn--d1acufc.xn--p1ai", false, true, ""}, // Punycode example
		{"example", true, true, ""},               // Technically valid as a single label
		{"invalid_domain", true, true, ""},
		{"localhost", true, true, ""},

		// Test cases for invalid domains
		{"-example.com", false, false, "invalid DNS name format"},
		{"example..com", false, false, "invalid DNS name format"},
		{"example", false, false, "single-label domains not allowed"},
		{"*.example..com", false, false, "invalid DNS name format"},
		{"invalid_domain", false, false, "single-label domains not allowed"},
		{"localhost", false, false, "single-label domains not allowed"},

		// Test cases for wildcard handling
		{"*.sub.example.com", false, true, ""},
		{"*example.com", false, false, "invalid DNS name format"}, // Wildcard must have `*.` prefix
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("Domain: '%s' AllowIPs: %v", test.domain, test.allowIPs), func(t *testing.T) {
			result, err := IsValidDomainWithError(test.domain, test.allowIPs)
			assert.Equal(t, test.expected, result)

			if test.errMsg == "" {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.EqualError(t, err, test.errMsg)
			}
		})
	}
}

// TestCleanFirstDomain checks the CleanFirstDomain method.
func TestCleanFirstDomain(t *testing.T) {
	tests := []struct {
		name    string
		domains []string
		want    string
	}{
		{"CleanWildcard", []string{"*.example.com"}, "example.com"},
		{"NoWildcard", []string{"example.com"}, "example.com"},
		{"EmptyDomains", []string{}, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpb := &CertificateProviderBase{Domains: tt.domains}
			if got := cpb.CleanFirstDomain(); got != tt.want {
				t.Errorf("CleanFirstDomain() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestGetDomains checks the GetDomains method.
func TestGetDomains(t *testing.T) {
	domains := []string{"example.com", "*.example.com"}
	cpb := &CertificateProviderBase{Domains: domains}
	if got := cpb.GetDomains(); !equal(got, domains) {
		t.Errorf("GetDomains() = %v, want %v", got, domains)
	}
}

// TestMatchDomain checks the MatchDomain method.
func TestMatchDomain(t *testing.T) {
	tests := []struct {
		name       string
		domains    []string
		serverName string
		want       bool
	}{
		{"ExactMatch", []string{"example.com"}, "example.com", true},
		{"WildcardMatch", []string{"*.example.com"}, "sub.example.com", true},
		{"NoMatch", []string{"example.com"}, "sub.example.com", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpb := &CertificateProviderBase{Domains: tt.domains}
			if got := cpb.MatchDomain(tt.serverName); got != tt.want {
				t.Errorf("MatchDomain() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Helper function to compare two slices.
func equal(a, b []string) bool {
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

// TestMatchExplicitDomain checks the MatchExplicitDomain method.
func TestMatchExplicitDomain(t *testing.T) {
	tests := []struct {
		name            string
		explicitDomains []string
		serverName      string
		want            bool
	}{
		{"ExactMatch", []string{"example.com"}, "example.com", true},
		{"CaseInsensitiveMatch", []string{"example.com"}, "EXAMPLE.COM", true},
		{"TrimmedMatch", []string{"example.com"}, " example.com ", true},
		{"NoMatch", []string{"example.com"}, "sub.example.com", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpb := &CertificateProviderBase{explicitDomains: tt.explicitDomains}
			if got := cpb.MatchExplicitDomain(tt.serverName); got != tt.want {
				t.Errorf("MatchExplicitDomain() = %v, want %v", got, tt.want)
			}
		})
	}

	testMerge := []struct {
		name       string
		domains    []string
		adddomains []string
		serverName string
		want       bool
	}{
		{"ExactMatch", []string{"example.com"}, []string{}, "example.com", true},
		{"CaseInsensitiveMatch", []string{"example.com"}, []string{}, "EXAMPLE.COM", true},
		{"TrimmedMatch", []string{"example.com"}, []string{}, " example.com ", true},
		{"NoMatch", []string{"example.com"}, []string{}, "sub.example.com", false},
		{"DupsExactMatch", []string{"example.com"}, []string{"example.com"}, "example.com", true},
		{"DupsCaseInsensitiveMatch", []string{"example.com"}, []string{"example.com"}, "EXAMPLE.COM", true},
		{"DupsTrimmedMatch", []string{"example.com"}, []string{"example.com"}, " example.com ", true},
		{"DupsNoMatch", []string{"example.com"}, []string{"example.com"}, "sub.example.com", false},
		{"MergeExactMatch", []string{"*.example.com"}, []string{"sub.example.com"}, "sub.example.com", true},
		{"MergeCaseInsensitiveMatch", []string{"*.example.com"}, []string{"sub.example.com"}, "SUB.EXAMPLE.COM", true},
		{"MergeTrimmedMatch", []string{"*.example.com"}, []string{"sub.example.com"}, " sub.example.com ", true},
		{"MergeNoMatch", []string{"*.example.com"}, []string{"sub.example.com"}, "sub.example.com", false},
	}

	for _, tt := range testMerge {
		t.Run(tt.name, func(t *testing.T) {
			cpb := &CertificateProviderBase{Domains: tt.domains, AddDomains: tt.adddomains}
			err := cpb.ValidateDomains(false)
			assert.NoError(t, err)
			if err != nil {
				if got := cpb.MatchExplicitDomain(tt.serverName); got != tt.want {
					t.Errorf("MatchExplicitDomain() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}
