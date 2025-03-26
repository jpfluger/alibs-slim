package anetwork

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"github.com/jpfluger/alibs-slim/autils"
	"math/big"
	"os"
	"path/filepath"
	"time"
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
	cert *tls.Certificate
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
		cert, err := GenerateSelfSignedCertificate(ss.Domains)
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

// SaveCertificateToFile saves the certificate to the specified files.
func SaveCertificateToFile(cert *tls.Certificate, publicPath, privatePath string) error {
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: cert.Certificate[0]})
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: cert.PrivateKey.(*ecdsa.PrivateKey).D.Bytes()})

	if err := os.WriteFile(publicPath, certPEM, 0644); err != nil {
		return err
	}

	if err := os.WriteFile(privatePath, keyPEM, 0600); err != nil {
		return err
	}

	return nil
}

// GenerateSelfSignedCertificate generates a self-signed certificate for the given domains.
func GenerateSelfSignedCertificate(domains []string) (*tls.Certificate, error) {
	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, err
	}

	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			CommonName: domains[0],
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(365 * 24 * time.Hour),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		DNSNames:              domains,
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		return nil, err
	}

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
