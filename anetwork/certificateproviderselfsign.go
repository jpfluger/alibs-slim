package anetwork

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"fmt"
	"math/big"
	"net"
	"os"
	"path/filepath"
	"time"

	"github.com/jpfluger/alibs-slim/autils"
)

// CERTIFICATEPROVIDERTYPE_SELFSIGN defines the type for self-signed certificates.
const CERTIFICATEPROVIDERTYPE_SELFSIGN CertificateProviderType = "self-sign"

// CertificateProviderSelfSign provides self-signed certificates.
type CertificateProviderSelfSign struct {
	CertificateProviderBase
	Keys struct {
		PublicFile  string `json:"publicFile"`
		PrivateFile string `json:"privateFile"`
	} `json:"keys"`
	IPs             []string `json:"ips"`
	AutoDetectIPs   bool     `json:"autoDetectIPs"`
	IncludeLoopback bool     `json:"includeLoopback"`
	IncludeLocalSAN bool     `json:"includeLocalSAN"`
	IncludeIPv4     bool     `json:"includeIPv4"`
	IncludeIPv6     bool     `json:"includeIPv6"`
	cert            *tls.Certificate
}

// Validate checks the validity of the self-signed certificate and generates it if necessary.
func (ss *CertificateProviderSelfSign) Validate(dirCache string) error {
	if err := ss.CertificateProviderBase.ValidateDomains(true); err != nil {
		return err
	}

	// Generate default filenames if not provided
	if ss.Keys.PublicFile == "" || ss.Keys.PrivateFile == "" {
		firstDomain := ss.CleanFirstDomain()
		ss.Keys.PublicFile = "selfsigned-cert." + firstDomain + ".pem"
		ss.Keys.PrivateFile = "selfsigned-key." + firstDomain + ".pem"
	}

	publicPath := filepath.Join(dirCache, ss.Keys.PublicFile)
	privatePath := filepath.Join(dirCache, ss.Keys.PrivateFile)

	if autils.FileExists(publicPath) && autils.FileExists(privatePath) {
		cert, err := LoadCertificateFromFile(publicPath, privatePath)
		if err != nil {
			return err
		}
		ss.cert = cert
	} else {
		cert, err := GenerateSelfSignedCertificate(SelfSignCertOpts{
			Domains:         ss.Domains,
			IPs:             ss.IPs,
			AutoDetectIPs:   ss.AutoDetectIPs,
			IncludeLoopback: ss.IncludeLoopback,
			IncludeLocalSAN: ss.IncludeLocalSAN,
			IncludeIPv4:     ss.IncludeIPv4,
			IncludeIPv6:     false,
		})
		if err != nil {
			return err
		}
		ss.cert = cert

		if err := SaveCertificateToFile(cert, publicPath, privatePath); err != nil {
			return err
		}
	}

	return nil
}

// GetCertificate returns the self-signed certificate.
func (ss *CertificateProviderSelfSign) GetCertificate(hello *tls.ClientHelloInfo) (*tls.Certificate, error) {
	return ss.cert, nil
}

// LoadCertificateFromFile loads a certificate from the specified files.
func LoadCertificateFromFile(publicPath, privatePath string) (*tls.Certificate, error) {
	certPEMBlock, err := os.ReadFile(publicPath)
	if err != nil {
		return nil, err
	}

	keyPEMBlock, err := os.ReadFile(privatePath)
	if err != nil {
		return nil, err
	}

	cert, err := tls.X509KeyPair(certPEMBlock, keyPEMBlock)
	if err != nil {
		return nil, err
	}

	return &cert, nil
}

// SelfSignCertOpts defines options for generating a self-signed certificate.
type SelfSignCertOpts struct {
	Domains         []string // User-provided domains/FQDNs for SAN
	IPs             []string // Exact IPs to include in SAN
	AutoDetectIPs   bool     // Auto-detect non-loopback IPs from interfaces
	IncludeLoopback bool     // Include loopback IPs (127.0.0.1 and ::1, filtered by IP version)
	IncludeLocalSAN bool     // Include localhost, hostname, and hostname.local in SAN
	IncludeIPv4     bool     // Include IPv4 addresses (default: true)
	IncludeIPv6     bool     // Include IPv6 addresses (default: true)
}

// GetIPs converts the string IPs in the opts to net.IP, filtering by IncludeIPv4 and IncludeIPv6.
func (sso *SelfSignCertOpts) GetIPs() ([]net.IP, error) {
	var validIPs []net.IP
	var invalidIPs []string

	for _, ipStr := range sso.IPs {
		ip := net.ParseIP(ipStr)
		if ip == nil {
			invalidIPs = append(invalidIPs, ipStr)
			continue
		}
		// Apply IPv4/IPv6 filters
		if (sso.IncludeIPv4 && ip.To4() != nil) || (sso.IncludeIPv6 && ip.To16() != nil && ip.To4() == nil) {
			validIPs = appendUniqueIP(validIPs, ip)
		}
	}

	if len(invalidIPs) > 0 {
		return validIPs, fmt.Errorf("invalid IP addresses provided: %v", invalidIPs)
	}

	return validIPs, nil
}

// SaveCertificateToFile saves the certificate to the specified files.
func SaveCertificateToFile(cert *tls.Certificate, publicPath, privatePath string) error {
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: cert.Certificate[0]})

	keyPEM, err := x509.MarshalECPrivateKey(cert.PrivateKey.(*ecdsa.PrivateKey))
	if err != nil {
		return err
	}
	keyPEMBlock := pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: keyPEM})

	if err := os.WriteFile(publicPath, certPEM, 0644); err != nil {
		return err
	}

	if err := os.WriteFile(privatePath, keyPEMBlock, 0600); err != nil {
		return err
	}

	return nil
}

// GenerateSelfSignedCertificate generates a self-signed certificate based on the provided options.
func GenerateSelfSignedCertificate(opts SelfSignCertOpts) (*tls.Certificate, error) {
	// Set defaults for IP version inclusion
	if !opts.IncludeIPv4 && !opts.IncludeIPv6 {
		// If both false, treat as including both to avoid unexpected empty IPs
		opts.IncludeIPv4 = true
		opts.IncludeIPv6 = true
	}

	// Get the machine's hostname
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "localhost" // Fallback
	}

	// Determine Common Name (CN): first domain if provided, else hostname or localhost
	commonName := hostname
	if len(opts.Domains) > 0 {
		commonName = opts.Domains[0]
	}

	// Build SAN DNS names: start with provided domains
	sanDomains := make([]string, 0, len(opts.Domains))
	sanDomains = append(sanDomains, opts.Domains...)

	// If IncludeLocalSAN, add local entries and avoid duplicates
	if opts.IncludeLocalSAN {
		localEntries := []string{"localhost", hostname, hostname + ".local"}
		for _, entry := range localEntries {
			duplicate := false
			for _, d := range sanDomains {
				if d == entry {
					duplicate = true
					break
				}
			}
			if !duplicate {
				sanDomains = append(sanDomains, entry)
			}
		}
	}

	// Build SAN IP addresses: start with user-provided IPs
	sanIPs, err := opts.GetIPs()
	if err != nil {
		return nil, fmt.Errorf("failed to parse provided IPs: %w", err)
	}

	// If IncludeLoopback, add loopback IPs filtered by version
	if opts.IncludeLoopback {
		if opts.IncludeIPv4 {
			sanIPs = appendUniqueIP(sanIPs, net.ParseIP("127.0.0.1"))
		}
		if opts.IncludeIPv6 {
			sanIPs = appendUniqueIP(sanIPs, net.ParseIP("::1"))
		}
	}

	// If AutoDetectIPs, detect and add non-loopback, non-link-local IPs filtered by version
	if opts.AutoDetectIPs {
		addrs, err := net.InterfaceAddrs()
		if err != nil {
			return nil, fmt.Errorf("failed to detect network interfaces: %w", err)
		}
		for _, addr := range addrs {
			ipNet, ok := addr.(*net.IPNet)
			if !ok {
				continue
			}
			ip := ipNet.IP
			if ip.IsLoopback() || ip.IsLinkLocalUnicast() {
				continue // Skip loopback and link-local
			}
			if (opts.IncludeIPv4 && ip.To4() != nil) || (opts.IncludeIPv6 && ip.To16() != nil && ip.To4() == nil) {
				sanIPs = appendUniqueIP(sanIPs, ip)
			}
		}
	}

	// Error if no SAN entries and no CN (prevent invalid cert)
	if len(sanDomains) == 0 && len(sanIPs) == 0 && commonName == "" {
		return nil, errors.New("no domains or IPs provided; cannot generate certificate")
	}

	// Generate ECDSA private key (FIPS-compliant with P-256)
	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, err
	}

	// Create certificate template
	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			CommonName: commonName,
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(365 * 24 * time.Hour),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		DNSNames:              sanDomains,
		IPAddresses:           sanIPs,
	}

	// Create self-signed certificate
	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		return nil, err
	}

	// Encode certificate and private key
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	keyPEM, err := x509.MarshalECPrivateKey(priv)
	if err != nil {
		return nil, err
	}

	keyPEMBlock := pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: keyPEM})
	cert, err := tls.X509KeyPair(certPEM, keyPEMBlock)
	if err != nil {
		return nil, err
	}

	return &cert, nil
}

// appendUniqueIP appends an IP to the list if it's not already present.
func appendUniqueIP(ips []net.IP, newIP net.IP) []net.IP {
	if newIP == nil {
		return ips
	}
	for _, ip := range ips {
		if ip.Equal(newIP) {
			return ips
		}
	}
	return append(ips, newIP)
}

// CertStorage holds PEM-encoded certificate and key for storage.
type CertStorage struct {
	CertPEM []byte `json:"cert_pem"`
	KeyPEM  []byte `json:"key_pem"`
}

// ToTLSCertificate parses the stored PEM into a tls.Certificate.
func (cs *CertStorage) ToTLSCertificate() (*tls.Certificate, error) {
	if len(cs.CertPEM) == 0 || len(cs.KeyPEM) == 0 {
		return nil, errors.New("missing cert or key PEM")
	}
	cert, err := tls.X509KeyPair(cs.CertPEM, cs.KeyPEM)
	if err != nil {
		return nil, err
	}
	return &cert, nil
}

// NewCertStorage creates a new CertStorage from the provided TLS certificate
// by encoding it to PEM format.
func NewCertStorage(cert *tls.Certificate) (*CertStorage, error) {
	if cert == nil {
		return nil, fmt.Errorf("provided certificate is nil")
	}
	if len(cert.Certificate) == 0 {
		return nil, fmt.Errorf("certificate has no public key data")
	}
	if cert.PrivateKey == nil {
		return nil, fmt.Errorf("certificate has no private key")
	}

	// Encode the certificate (public part) to PEM
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: cert.Certificate[0]})

	// Encode the private key to PEM (assuming ECDSA; adjust if other types are used)
	ecdsaKey, ok := cert.PrivateKey.(*ecdsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("private key is not an ECDSA key")
	}
	keyBytes, err := x509.MarshalECPrivateKey(ecdsaKey)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal private key: %v", err)
	}
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: keyBytes})

	// Create and return CertStorage
	return &CertStorage{
		CertPEM: certPEM,
		KeyPEM:  keyPEM,
	}, nil
}

// NewCertStorageByFile creates a new CertStorage by reading PEM-encoded certificate
// and key from the provided file paths.
func NewCertStorageByFile(pubCert string, privCert string) (*CertStorage, error) {
	// Read public certificate file
	certBytes, err := os.ReadFile(pubCert)
	if err != nil {
		return nil, fmt.Errorf("failed to read public cert file %s: %v", pubCert, err)
	}
	if len(certBytes) == 0 {
		return nil, fmt.Errorf("public cert file %s is empty", pubCert)
	}

	// Read private key file
	keyBytes, err := os.ReadFile(privCert)
	if err != nil {
		return nil, fmt.Errorf("failed to read private key file %s: %v", privCert, err)
	}
	if len(keyBytes) == 0 {
		return nil, fmt.Errorf("private key file %s is empty", privCert)
	}

	// Create and return CertStorage (assumes files are already PEM-encoded)
	return &CertStorage{
		CertPEM: certBytes,
		KeyPEM:  keyBytes,
	}, nil
}
