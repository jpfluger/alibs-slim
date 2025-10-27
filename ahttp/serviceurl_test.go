package ahttp

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/jpfluger/alibs-slim/anetwork"
	"github.com/jpfluger/alibs-slim/autils"
)

// createServiceUrlTestFiles creates a temporary directory and optional cert/key files.
func createServiceUrlTestFiles(createCert bool, createKey bool) (dirRoot string, err error) {
	// Create a temporary root directory
	dirRoot, err = autils.CreateTempDir()
	if err != nil {
		return "", err
	}

	if createCert || createKey {
		// Generate a valid self-signed certificate for tests requiring valid PEM.
		opts := anetwork.SelfSignCertOpts{
			Domains:         []string{"example.com"},
			IncludeLocalSAN: true,
			IncludeLoopback: true,
			IncludeIPv4:     true,
		}
		tlsCert, err := anetwork.GenerateSelfSignedCertificate(opts)
		if err != nil {
			os.RemoveAll(dirRoot)
			return "", fmt.Errorf("failed to generate test certificate: %v", err)
		}

		// Write certificate file if requested
		if createCert {
			certFilePath := filepath.Join(dirRoot, "cert.pem")
			if err := anetwork.SaveCertificateToFile(tlsCert, certFilePath, certFilePath+".key"); err != nil {
				os.RemoveAll(dirRoot)
				return "", fmt.Errorf("failed to write test cert file: %v", err)
			}
		}

		// Write key file if requested
		if createKey {
			keyFilePath := filepath.Join(dirRoot, "key.pem")
			if err := anetwork.SaveCertificateToFile(tlsCert, keyFilePath+".cert", keyFilePath); err != nil {
				os.RemoveAll(dirRoot)
				return "", fmt.Errorf("failed to write test key file: %v", err)
			}
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

// TestNewServiceUrl_NilOpts tests NewServiceUrl with nil options (defaults).
func TestNewServiceUrl_NilOpts(t *testing.T) {
	su, err := NewServiceUrl(nil)
	if err != nil {
		t.Fatalf("Expected NewServiceUrl to succeed with nil opts, but got error: %v", err)
	}

	if su == nil {
		t.Fatal("Expected ServiceUrl to be created, but got nil")
	}

	if !su.GetIsTLS() {
		t.Errorf("Expected IsTLS to be true (default RequireCerts), but got false")
	}

	if su.ListenPort != 443 {
		t.Errorf("Expected default ListenPort to be 443 for HTTPS, but got %d", su.ListenPort)
	}

	if su.PublicUrl != fmt.Sprintf("https://localhost:%d", su.ListenPort) {
		t.Errorf("Expected default PublicUrl to be 'https://localhost:%d', but got %s", su.ListenPort, su.PublicUrl)
	}

	if su.CertStorage == nil {
		t.Errorf("Expected CertStorage to be set (default UseCertStorage), but got nil")
	}
}

// TestNewServiceUrl_WithDefault tests NewServiceUrl with a provided default ServiceUrl.
func TestNewServiceUrl_WithDefault(t *testing.T) {
	defaultSU := &ServiceUrl{
		PublicUrl:  "http://test.com",
		ListenPort: 8080,
	}

	opts := &ServiceUrlOpts{
		Default: defaultSU,
	}

	su, err := NewServiceUrl(opts)
	if err != nil {
		t.Fatalf("Expected NewServiceUrl to succeed with default, but got error: %v", err)
	}

	if su != defaultSU {
		t.Errorf("Expected returned ServiceUrl to be the default one")
	}

	if su.PublicUrl != "http://test.com" {
		t.Errorf("Expected PublicUrl to remain 'http://test.com', but got %s", su.PublicUrl)
	}
}

// TestNewServiceUrl_FreePortSearch tests NewServiceUrl with FreePortSearch.
func TestNewServiceUrl_FreePortSearch(t *testing.T) {
	opts := &ServiceUrlOpts{
		FreePortSearch: 50000, // Start from a high ephemeral port to avoid conflicts in tests.
		RequireCerts:   false,
	}

	su, err := NewServiceUrl(opts)
	if err != nil {
		t.Fatalf("Expected NewServiceUrl to succeed, but got error: %v", err)
	}

	if su.ListenPort < 50000 {
		t.Errorf("Expected ListenPort >= 50000, but got %d", su.ListenPort)
	}
}

// TestNewServiceUrl_NoCerts tests NewServiceUrl without requiring certs.
func TestNewServiceUrl_NoCerts(t *testing.T) {
	opts := &ServiceUrlOpts{
		RequireCerts: false,
	}

	su, err := NewServiceUrl(opts)
	if err != nil {
		t.Fatalf("Expected NewServiceUrl to succeed, but got error: %v", err)
	}

	if su.GetIsTLS() {
		t.Errorf("Expected IsTLS to be false, but got true")
	}

	if su.ListenPort != 80 {
		t.Errorf("Expected default ListenPort to be 80 for HTTP, but got %d", su.ListenPort)
	}

	if su.CertStorage != nil {
		t.Errorf("Expected CertStorage to be nil, but got non-nil")
	}
}

// TestNewServiceUrl_SaveToFiles tests NewServiceUrl with DirCerts (save to files).
func TestNewServiceUrl_SaveToFiles(t *testing.T) {
	dirCerts, err := autils.CreateTempDir()
	if dirCerts != "" {
		defer os.RemoveAll(dirCerts)
	}
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}

	opts := &ServiceUrlOpts{
		RequireCerts:   true,
		DirCerts:       dirCerts,
		UseCertStorage: false, // Disable storage to test files only.
	}

	su, err := NewServiceUrl(opts)
	if err != nil {
		t.Fatalf("Expected NewServiceUrl to succeed, but got error: %v", err)
	}

	if su.CertFile == "" || su.KeyFile == "" {
		t.Errorf("Expected CertFile and KeyFile to be set, but got empty")
	}

	if _, err := os.Stat(su.CertFile); err != nil {
		t.Errorf("Expected cert file to exist, but got error: %v", err)
	}

	if _, err := os.Stat(su.KeyFile); err != nil {
		t.Errorf("Expected key file to exist, but got error: %v", err)
	}
}

// TestNewServiceUrl_UseCertStorage tests NewServiceUrl with UseCertStorage.
func TestNewServiceUrl_UseCertStorage(t *testing.T) {
	opts := &ServiceUrlOpts{
		RequireCerts:   true,
		UseCertStorage: true,
		DirCerts:       "", // No file saving.
	}

	su, err := NewServiceUrl(opts)
	if err != nil {
		t.Fatalf("Expected NewServiceUrl to succeed, but got error: %v", err)
	}

	if su.CertStorage == nil {
		t.Errorf("Expected CertStorage to be set, but got nil")
	}

	if su.CertFile != "" || su.KeyFile != "" {
		t.Errorf("Expected no file paths when using CertStorage only")
	}
}

// TestNewServiceUrl_InvalidDefault tests NewServiceUrl with invalid default.
func TestNewServiceUrl_InvalidDefault(t *testing.T) {
	invalidDefault := &ServiceUrl{
		ListenPort: 70000, // Invalid port.
	}

	opts := &ServiceUrlOpts{
		Default: invalidDefault,
	}

	_, err := NewServiceUrl(opts)
	if err == nil {
		t.Fatal("Expected NewServiceUrl to fail with invalid default, but it passed")
	}
}

// TestNewServiceUrl_InvalidDirCerts tests NewServiceUrl with invalid DirCerts.
func TestNewServiceUrl_InvalidDirCerts(t *testing.T) {
	opts := &ServiceUrlOpts{
		RequireCerts: true,
		DirCerts:     "/invalid/dir",
	}

	_, err := NewServiceUrl(opts)
	if err == nil {
		t.Fatal("Expected NewServiceUrl to fail with invalid DirCerts, but it passed")
	}
}

// TestNewServiceUrl_LoadFromFilesToStorage tests loading files into CertStorage when UseCertStorage is true.
func TestNewServiceUrl_LoadFromFilesToStorage(t *testing.T) {
	dirRoot, err := createServiceUrlTestFiles(true, true) // Create fake cert/key files.
	if dirRoot != "" {
		defer os.RemoveAll(dirRoot)
	}
	if err != nil {
		t.Fatalf("failed to create test files: %v", err)
	}

	defaultSU := &ServiceUrl{
		PublicUrl:  "https://example.com",
		ListenPort: 443,
		CertFile:   filepath.Join(dirRoot, "cert.pem"),
		KeyFile:    filepath.Join(dirRoot, "key.pem"),
	}

	opts := &ServiceUrlOpts{
		Default:        defaultSU,
		UseCertStorage: true,
	}

	su, err := NewServiceUrl(opts)
	if err != nil {
		t.Fatalf("Expected NewServiceUrl to succeed, but got error: %v", err)
	}

	if su.CertStorage == nil {
		t.Errorf("Expected CertStorage to be set after loading from files, but got nil")
	}
}
