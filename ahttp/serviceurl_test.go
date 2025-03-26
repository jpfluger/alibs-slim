package ahttp

import (
	"github.com/jpfluger/alibs-slim/autils"
	"os"
	"path/filepath"
	"testing"
)

// createServiceUrlTestFiles creates a temporary directory and optional cert/key files.
func createServiceUrlTestFiles(createCert bool, createKey bool) (dirRoot string, err error) {
	// Create a temporary root directory
	dirRoot, err = autils.CreateTempDir()
	if err != nil {
		return "", err
	}

	// Create fake cert and key files
	certFilePath := filepath.Join(dirRoot, "cert.pem")
	keyFilePath := filepath.Join(dirRoot, "key.pem")

	if createCert {
		if err = os.WriteFile(certFilePath, []byte("fake cert"), 0644); err != nil {
			return "", err
		}
	}
	if createKey {
		if err = os.WriteFile(keyFilePath, []byte("fake key"), 0644); err != nil {
			return "", err
		}
	}
	return dirRoot, nil
}

// TestServiceUrl_Validate_SuccessHTTP tests validation of an HTTP service.
func TestServiceUrl_Validate_SuccessHTTP(t *testing.T) {
	dirRoot, err := createServiceUrlTestFiles(false, false)
	if dirRoot != "" {
		defer os.RemoveAll(dirRoot)
	}
	if err != nil {
		t.Fatalf("failed to create test files: %s", err)
	}

	su := &ServiceUrl{
		PublicUrl:  "http://example.com",
		ListenPort: 8080,
	}

	err = su.Validate()
	if err != nil {
		t.Fatalf("Expected validation to pass, but got error: %v", err)
	}

	if su.GetIsTLS() {
		t.Errorf("Expected GetIsTLS() to return false for HTTP, but got true")
	}

	if su.ListenPort != 8080 {
		t.Errorf("Expected listen port to be 8080, but got %d", su.ListenPort)
	}
}

// TestServiceUrl_Validate_SuccessHTTPS tests validation of an HTTPS service.
func TestServiceUrl_Validate_SuccessHTTPS(t *testing.T) {
	dirRoot, err := createServiceUrlTestFiles(true, true)
	if dirRoot != "" {
		defer os.RemoveAll(dirRoot)
	}
	if err != nil {
		t.Fatalf("failed to create test files: %s", err)
	}

	su := &ServiceUrl{
		PublicUrl:  "https://example.com",
		ListenPort: 443,
		CertFile:   filepath.Join(dirRoot, "cert.pem"),
		KeyFile:    filepath.Join(dirRoot, "key.pem"),
	}

	err = su.Validate()
	if err != nil {
		t.Fatalf("Expected validation to pass, but got error: %v", err)
	}

	if !su.GetIsTLS() {
		t.Errorf("Expected GetIsTLS() to return true for HTTPS, but got false")
	}
}

// TestServiceUrl_Validate_FailInvalidPort tests validation failure for an out-of-range port.
func TestServiceUrl_Validate_FailInvalidPort(t *testing.T) {
	su := &ServiceUrl{
		PublicUrl:  "http://example.com",
		ListenPort: 70000, // Invalid port
	}

	err := su.Validate()
	if err == nil {
		t.Fatal("Expected validation to fail due to invalid port, but it passed")
	}
}

// TestServiceUrl_Validate_FailMissingCert tests validation failure when HTTPS is used but no cert file exists.
func TestServiceUrl_Validate_FailMissingCert(t *testing.T) {
	dirRoot, err := createServiceUrlTestFiles(false, true) // Cert missing
	if dirRoot != "" {
		defer os.RemoveAll(dirRoot)
	}
	if err != nil {
		t.Fatalf("failed to create test files: %s", err)
	}

	su := &ServiceUrl{
		PublicUrl:  "https://example.com",
		ListenPort: 443,
		CertFile:   filepath.Join(dirRoot, "missing.crt"),
		KeyFile:    filepath.Join(dirRoot, "key.pem"),
	}

	err = su.Validate()
	if err == nil {
		t.Fatal("Expected validation to fail due to missing cert file, but it passed")
	}
}

// TestServiceUrl_Validate_FailMissingKey tests validation failure when HTTPS is used but no key file exists.
func TestServiceUrl_Validate_FailMissingKey(t *testing.T) {
	dirRoot, err := createServiceUrlTestFiles(true, false) // Key missing
	if dirRoot != "" {
		defer os.RemoveAll(dirRoot)
	}
	if err != nil {
		t.Fatalf("failed to create test files: %s", err)
	}

	su := &ServiceUrl{
		PublicUrl:  "https://example.com",
		ListenPort: 443,
		CertFile:   filepath.Join(dirRoot, "cert.pem"),
		KeyFile:    filepath.Join(dirRoot, "missing.key"),
	}

	err = su.Validate()
	if err == nil {
		t.Fatal("Expected validation to fail due to missing key file, but it passed")
	}
}

// TestServiceUrl_GetPublicUrl tests GetPublicUrl after validation.
func TestServiceUrl_GetPublicUrl(t *testing.T) {
	su := &ServiceUrl{
		PublicUrl: "http://example.com",
	}

	err := su.Validate()
	if err != nil {
		t.Fatalf("Validation failed: %v", err)
	}

	if su.GetPublicUrl().String() != "http://example.com" {
		t.Errorf("Expected GetPublicUrl() to return 'http://example.com', but got %s", su.GetPublicUrl())
	}
}

// TestServiceUrl_SetRootDir tests setting a valid root directory.
func TestServiceUrl_SetRootDir(t *testing.T) {
	dirRoot, err := createServiceUrlTestFiles(false, false)
	if dirRoot != "" {
		defer os.RemoveAll(dirRoot)
	}
	if err != nil {
		t.Fatalf("failed to create test files: %s", err)
	}

	su := &ServiceUrl{}
	err = su.SetRootDir(dirRoot)
	if err != nil {
		t.Fatalf("Expected SetRootDir to succeed, but got error: %v", err)
	}

	if su.dirRoot != dirRoot {
		t.Errorf("Expected dirRoot to be '%s', but got '%s'", dirRoot, su.dirRoot)
	}
}

// TestServiceUrl_SetRootDir_Fail tests setting an invalid root directory.
func TestServiceUrl_SetRootDir_Fail(t *testing.T) {
	su := &ServiceUrl{}
	err := su.SetRootDir("")
	if err == nil {
		t.Fatal("Expected SetRootDir to fail for an empty directory, but it passed")
	}
}

// TestServiceUrl_GetCertFile tests GetCertFile behavior.
func TestServiceUrl_GetCertFile(t *testing.T) {
	su := &ServiceUrl{
		CertFile: "cert.pem",
	}
	if su.GetCertFile() != "cert.pem" {
		t.Errorf("Expected GetCertFile() to return 'cert.pem', but got '%s'", su.GetCertFile())
	}
}

// TestServiceUrl_GetKeyFile tests GetKeyFile behavior.
func TestServiceUrl_GetKeyFile(t *testing.T) {
	su := &ServiceUrl{
		KeyFile: "key.pem",
	}
	if su.GetKeyFile() != "key.pem" {
		t.Errorf("Expected GetKeyFile() to return 'key.pem', but got '%s'", su.GetKeyFile())
	}
}
