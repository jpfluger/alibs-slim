package anetwork

import (
	"path/filepath"
	"testing"

	"github.com/jpfluger/alibs-slim/autils"
)

func TestCertificateProviderSelfSign_Validate(t *testing.T) {
	dirCache := t.TempDir()

	tests := []struct {
		name    string
		domains []string
		wantErr bool
	}{
		{"ValidDomains", []string{"example.com", "*.example.com"}, false},
		{"InvalidDomain", []string{"example..com"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ss := &CertificateProviderSelfSign{
				CertificateProviderBase: CertificateProviderBase{
					Domains: tt.domains,
				},
			}

			err := ss.Validate(dirCache)
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				// Check if the certificate files were created
				publicPath := filepath.Join(dirCache, ss.Keys.PublicFile)
				privatePath := filepath.Join(dirCache, ss.Keys.PrivateFile)
				if !autils.FileExists(publicPath) || !autils.FileExists(privatePath) {
					t.Errorf("Certificate files not created")
				}

				// Check if the certificate is loaded
				if ss.cert == nil {
					t.Errorf("Certificate not loaded")
				}
			}
		})
	}
}

func TestGenerateSelfSignedCertificate(t *testing.T) {
	domains := []string{"example.com", "*.example.com"}
	cert, err := GenerateSelfSignedCertificate(SelfSignCertOpts{Domains: domains})
	if err != nil {
		t.Fatalf("GenerateSelfSignedCertificate() error = %v", err)
	}

	if cert == nil {
		t.Fatalf("GenerateSelfSignedCertificate() returned nil certificate")
	}

	if len(cert.Certificate) == 0 {
		t.Fatalf("Generated certificate is empty")
	}
}
