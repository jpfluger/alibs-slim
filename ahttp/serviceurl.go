package ahttp

import (
	"fmt"
	"github.com/jpfluger/alibs-slim/anetwork"
	"github.com/jpfluger/alibs-slim/autils"
	"net/url"
	"path"
	"strings"
)

// ServiceUrl holds configuration details for a service, including its port, public URL,
// and TLS certificate information if applicable.
type ServiceUrl struct {
	ListenPort int    `json:"listenPort"` // Port on which the service listens, between 1 and 65535.
	PublicUrl  string `json:"publicUrl"`  // Public URL for the service.
	CertFile   string `json:"certFile"`   // Path to the TLS certificate file (required if using HTTPS).
	KeyFile    string `json:"keyFile"`    // Path to the TLS key file (required if using HTTPS).

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
// - TLS file requirements if using HTTPS.
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

	// If HTTPS is used, validate existence of certificate and key files.
	if su.isTLS {
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
