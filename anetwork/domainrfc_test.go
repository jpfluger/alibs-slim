package anetwork

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDomainRFC_IsEmpty(t *testing.T) {
	tests := []struct {
		domain   DomainRFC
		expected bool
	}{
		{"", true},
		{"   ", true},
		{"example.com", false},
	}

	for _, test := range tests {
		assert.Equal(t, test.expected, test.domain.IsEmpty(), "Domain: '%s'", test.domain)
	}
}

func TestDomainRFC_IsValid(t *testing.T) {
	tests := []struct {
		domain   DomainRFC
		expected bool
		errMsg   string
	}{
		{"example.com", true, ""},
		{"sub.example.com", true, ""},
		{"*.example.com", true, ""},
		{"localhost", false, "single-label domains not allowed"},
		{"invalid_domain", false, "single-label domains not allowed"},
		{"192.168.1.1", false, "IP addresses are not allowed unless explicitly permitted"},
		{"", false, "domain is empty"},
	}

	for _, test := range tests {
		ok, err := test.domain.IsValid()
		assert.Equal(t, test.expected, ok, "Domain: '%s'", test.domain)
		if !test.expected {
			assert.Error(t, err)
			assert.EqualError(t, err, test.errMsg)
		} else {
			assert.NoError(t, err)
		}
	}
}

func TestDomainRFC_IsValidWithOptions(t *testing.T) {
	tests := []struct {
		domain   DomainRFC
		allowIPs bool
		expected bool
		errMsg   string
	}{
		{"example.com", false, true, ""},
		{"192.168.1.1", true, true, ""},
		{"192.168.1.1", false, false, "IP addresses are not allowed unless explicitly permitted"},
		{"localhost", false, false, "single-label domains not allowed"},
		{"localhost", true, true, ""},
		{"*.example.com", false, true, ""},
	}

	for _, test := range tests {
		ok, err := test.domain.IsValidWithOptions(test.allowIPs)
		assert.Equal(t, test.expected, ok, "Domain: '%s'", test.domain)
		if !test.expected {
			assert.Error(t, err)
			assert.EqualError(t, err, test.errMsg)
		} else {
			assert.NoError(t, err)
		}
	}
}

func TestDomainRFC_Normalize(t *testing.T) {
	tests := []struct {
		domain   DomainRFC
		expected DomainRFC
	}{
		{" example.com ", "example.com"},
		{"EXAMPLE.COM", "example.com"},
		{"xn--d1acufc.xn--p1ai", "xn--d1acufc.xn--p1ai"}, // Already ASCII
	}

	for _, test := range tests {
		normalized, err := test.domain.Normalize()
		assert.NoError(t, err)
		assert.Equal(t, test.expected, normalized, "Domain: '%s'", test.domain)
	}
}

func TestDomainRFCs_FilterInvalid(t *testing.T) {
	domains := DomainRFCs{
		"example.com",
		"localhost",
		"192.168.1.1",
		"invalid_domain",
		"*.example.com",
	}

	valid, invalid := domains.FilterInvalid()
	assert.ElementsMatch(t, valid, DomainRFCs{"example.com", "*.example.com"})
	assert.Len(t, invalid, 3)
	assert.Contains(t, invalid, DomainRFC("localhost"))
	assert.Contains(t, invalid, DomainRFC("192.168.1.1"))
	assert.Contains(t, invalid, DomainRFC("invalid_domain"))
}

func TestDomainRFCs_FilterInvalidWithErrors(t *testing.T) {
	domains := DomainRFCs{
		"example.com",
		"localhost",
		"192.168.1.1",
		"invalid_domain",
		"*.example.com",
		"bad.****.example.com",
	}

	valid, invalid := domains.FilterInvalidWithErrors(false)
	assert.ElementsMatch(t, valid, DomainRFCs{"example.com", "*.example.com"})
	assert.Len(t, invalid, 4)
	assert.EqualError(t, invalid[DomainRFC("localhost")], "single-label domains not allowed")
	assert.EqualError(t, invalid[DomainRFC("192.168.1.1")], "IP addresses are not allowed unless explicitly permitted")
	assert.EqualError(t, invalid[DomainRFC("invalid_domain")], "single-label domains not allowed")
	assert.EqualError(t, invalid[DomainRFC("bad.****.example.com")], "invalid DNS name format")
}
