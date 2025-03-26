package anetwork

import (
	"crypto/tls"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

// Mock implementation of ICertificateProvider for testing
type MockCertificateProvider struct {
	CertificateProviderBase
}

func (m *MockCertificateProvider) GetType() CertificateProviderType {
	return "mock"
}

func (m *MockCertificateProvider) Validate(dirCache string) error {
	return nil
}

func (m *MockCertificateProvider) GetCertificate(hello *tls.ClientHelloInfo) (*tls.Certificate, error) {
	return nil, nil
}

func TestDomainDetail_GetName(t *testing.T) {
	domain := &DomainDetail{Name: "example.com"}
	assert.Equal(t, "example.com", domain.GetName())
}

func TestDomainDetail_GetCertProvider(t *testing.T) {
	provider := &MockCertificateProvider{}
	domain := &DomainDetail{CertProvider: provider}
	assert.Equal(t, provider, domain.GetCertProvider())
}

func TestDomainDetail_UnmarshalJSON(t *testing.T) {
	jsonData := []byte(`{
		"name": "example.com",
		"certProvider": {
			"type": "self-sign",
			"domains": ["example.com"]
		}
	}`)

	var domain DomainDetail
	err := json.Unmarshal(jsonData, &domain)
	assert.NoError(t, err)
	assert.Equal(t, "example.com", domain.GetName())
	assert.NotNil(t, domain.GetCertProvider())
	assert.Equal(t, CERTIFICATEPROVIDERTYPE_SELFSIGN, domain.GetCertProvider().GetType())
}

func TestDomainDetail_UnmarshalJSON_InvalidType(t *testing.T) {
	jsonData := []byte(`{
		"name": "example.com",
		"certProvider": {
			"type": "invalid-type",
			"domains": ["example.com"]
		}
	}`)

	var domain DomainDetail
	err := json.Unmarshal(jsonData, &domain)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot find type struct 'invalid-type'")
}

func TestDomainDetail_UnmarshalJSON_EmptyCertProvider(t *testing.T) {
	jsonData := []byte(`{
		"name": "example.com",
		"certProvider": {}
	}`)

	var domain DomainDetail
	err := json.Unmarshal(jsonData, &domain)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "type field not found or is not a string in certProvider")
}
