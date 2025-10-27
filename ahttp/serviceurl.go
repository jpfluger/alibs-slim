package ahttp

import (
	"fmt"
	"net/url"
	"path"
	"path/filepath"
	"strings"

	"github.com/jpfluger/alibs-slim/anetwork"
	"github.com/jpfluger/alibs-slim/autils"
)

// ServiceUrl holds configuration details for a service, including its port, public URL,
// and TLS certificate information if applicable.
type ServiceUrl struct {
	ListenPort int    `json:"listenPort"` // Port on which the service listens, between 1 and 65535.
	PublicUrl  string `json:"publicUrl"`  // Public URL for the service.
	CertFile   string `json:"certFile"`   // Path to the TLS certificate file (required if using HTTPS).
	KeyFile    string `json:"keyFile"`    // Path to the TLS key file (required if using HTTPS).

	// CertStorage is an optional edge-case scenario where its expected
	// the certs are embedded into the ServiceUrl.
	CertStorage *anetwork.CertStorage `json:"certStorage"`

	u        *url.URL // Parsed URL object for internal use.
	isTLS    bool     // Indicates if the service uses HTTPS (TLS).
	dirRoot  string
	certFile string
	keyFile  string
}

// GetPublicUrl returns the parsed public URL as a *url.URL type.
// It is assumed that Validate() has been called before this method,
// so 'u' should be non-nil after successful validation.
func (su *ServiceUrl) GetPublicUrl() *url.URL {
	return su.u
}

// GetIsTLS indicates whether the service is using HTTPS (TLS).
func (su *ServiceUrl) GetIsTLS() bool {
	return su.isTLS
}

// Validate checks the ServiceUrl configuration for errors, including:
// - Valid listen port within the range [1, 65535]
// - Properly formatted PublicUrl
// - TLS file or storage requirements if using HTTPS.
// Returns an error if validation fails, or nil if successful.
func (su *ServiceUrl) Validate() error {
	if su == nil {
		return fmt.Errorf("serviceUrl struct is nil")
	}

	su.PublicUrl = strings.TrimSpace(su.PublicUrl)
	if su.ListenPort == 0 {
		if strings.HasPrefix(su.PublicUrl, "https://") {
			su.ListenPort = 443
		} else {
			su.ListenPort = 80
		}
	}

	// Validate the listen port range.
	if anetwork.IsOutsidePortRange(su.ListenPort) {
		return fmt.Errorf("listenPort %d is out of range (1-65535)", su.ListenPort)
	}

	// Parse and store the public URL, and check for errors.
	parsedUrl, err := url.Parse(su.PublicUrl)
	if err != nil {
		return fmt.Errorf("failed to parse publicUrl: %v", err)
	}
	su.u = parsedUrl

	// Check if the service URL scheme indicates HTTPS.
	su.isTLS = strings.HasPrefix(su.u.String(), "https://")

	// If HTTPS is used, validate certificate and key (files or storage).
	if su.isTLS {
		if su.CertStorage != nil {
			if _, err = su.CertStorage.ToTLSCertificate(); err != nil {
				return fmt.Errorf("invalid cert storage: %v", err)
			}
			// CertStorage is valid, so skip file checks.
			return nil
		}
		// Only check files if CertStorage is not used.
		su.CertFile = strings.TrimSpace(su.CertFile)
		su.KeyFile = strings.TrimSpace(su.KeyFile)
		if su.dirRoot != "" {
			if !path.IsAbs(su.CertFile) {
				su.certFile = path.Join(su.dirRoot, su.CertFile)
			}
			if !path.IsAbs(su.KeyFile) {
				su.keyFile = path.Join(su.dirRoot, su.KeyFile)
			}
		}
		if _, err = autils.ResolveFile(su.CertFile); err != nil {
			return fmt.Errorf("cert file is required when using HTTPS; %v", err)
		}
		if _, err = autils.ResolveFile(su.KeyFile); err != nil {
			return fmt.Errorf("key file is required when using HTTPS; %v", err)
		}
	}

	return nil
}

// SetRootDir sets the root directory. If any certificate is a
// relative path, then the rootDir is applied, if valid.
func (su *ServiceUrl) SetRootDir(dirRoot string) error {
	dirRoot = strings.TrimSpace(dirRoot)
	if _, err := autils.ResolveDirectory(dirRoot); err != nil {
		return fmt.Errorf("failed to resolve root dir: %v", err)
	}
	su.dirRoot = dirRoot
	return nil
}

// GetCertFile returns the certificate file.
func (su *ServiceUrl) GetCertFile() string {
	if su.certFile != "" {
		return su.certFile
	}
	return su.CertFile
}

// GetKeyFile returns the key file.
func (su *ServiceUrl) GetKeyFile() string {
	if su.keyFile != "" {
		return su.keyFile
	}
	return su.KeyFile
}

// ServiceUrlOpts defines options for creating a new ServiceUrl.
type ServiceUrlOpts struct {
	Default        *ServiceUrl
	FreePortSearch int
	RequireCerts   bool
	SSOpts         *anetwork.SelfSignCertOpts // If RequireCerts is true and SSOpts is nil, use QNAP/Synology-like defaults.
	DirCerts       string                     // If set, save certs to files in this dir; otherwise, use CertStorage if enabled.
	UseCertStorage bool                       // If true, load certs into memory (CertStorage).
}

// NewServiceUrl creates a new ServiceUrl based on the provided options.
func NewServiceUrl(opts *ServiceUrlOpts) (*ServiceUrl, error) {
	if opts == nil {
		opts = &ServiceUrlOpts{
			Default:        nil,
			FreePortSearch: 0,
			RequireCerts:   true,
			SSOpts: &anetwork.SelfSignCertOpts{
				AutoDetectIPs:   true,
				IncludeLoopback: true,
				IncludeLocalSAN: true,
				IncludeIPv4:     true,
			},
			UseCertStorage: true, // Default to in-memory for security.
		}
	}

	var su *ServiceUrl
	if opts.Default != nil {
		su = opts.Default
		if err := su.Validate(); err != nil {
			return nil, err
		}
	} else {
		su = &ServiceUrl{
			ListenPort: 0, // Will be set below.
			PublicUrl:  "",
		}
	}

	// Handle port: Use FreePortSearch to find a free port if set.
	if opts.FreePortSearch > 0 {
		port, err := anetwork.FindNextOpenPort(opts.FreePortSearch)
		if err != nil {
			return nil, fmt.Errorf("failed to find free port: %v", err)
		}
		su.ListenPort = port
	} else if su.ListenPort == 0 {
		// Fallback default port based on cert requirement.
		if opts.RequireCerts {
			su.ListenPort = 443
		} else {
			su.ListenPort = 80
		}
	}

	// Set PublicUrl if not provided.
	if su.PublicUrl == "" {
		scheme := "http"
		if opts.RequireCerts {
			scheme = "https"
		}
		su.PublicUrl = fmt.Sprintf("%s://localhost:%d", scheme, su.ListenPort)
	}

	// Handle certs if required and not already set.
	if opts.RequireCerts && (su.CertFile == "" && su.KeyFile == "" && su.CertStorage == nil) {
		// Use provided SSOpts or fallback to defaults.
		ssOpts := opts.SSOpts
		if ssOpts == nil {
			// QNAP/Synology-like defaults: local SAN, loopback, IPv4.
			ssOpts = &anetwork.SelfSignCertOpts{
				AutoDetectIPs:   true,
				IncludeLoopback: true,
				IncludeLocalSAN: true,
				IncludeIPv4:     true,
			}
		}

		// Generate self-signed cert.
		tlsCert, err := anetwork.GenerateSelfSignedCertificate(*ssOpts)
		if err != nil {
			return nil, fmt.Errorf("failed to generate self-signed cert: %v", err)
		}

		if opts.DirCerts != "" {
			// Save to files if DirCerts is set.
			if _, err := autils.ResolveDirectory(opts.DirCerts); err != nil {
				return nil, fmt.Errorf("invalid cert directory: %v", err)
			}
			certPath := filepath.Join(opts.DirCerts, "selfsigned-cert.pem")
			keyPath := filepath.Join(opts.DirCerts, "selfsigned-key.pem")
			if err := anetwork.SaveCertificateToFile(tlsCert, certPath, keyPath); err != nil {
				return nil, fmt.Errorf("failed to save cert files: %v", err)
			}
			su.CertFile = certPath
			su.KeyFile = keyPath
		}

		if opts.UseCertStorage {
			// Load into CertStorage.
			certStorage, err := anetwork.NewCertStorage(tlsCert)
			if err != nil {
				return nil, fmt.Errorf("failed to load cert into storage: %v", err)
			}
			su.CertStorage = certStorage
		}
	}

	// New check: If UseCertStorage is true but CertStorage is nil, load from files if available.
	if opts.UseCertStorage && su.CertStorage == nil && su.CertFile != "" && su.KeyFile != "" {
		certStorage, err := anetwork.NewCertStorageByFile(su.CertFile, su.KeyFile)
		if err != nil {
			return nil, fmt.Errorf("failed to load cert from files into storage: %v", err)
		}
		su.CertStorage = certStorage
	}

	// Final validation.
	if err := su.Validate(); err != nil {
		return nil, err
	}

	return su, nil
}
