package anetwork

import (
	"crypto/tls"
	"fmt"
	"github.com/jpfluger/alibs-slim/autils"
)

// CERTIFICATEPROVIDERTYPE_SIGNED defines the type for signed certificates.
const CERTIFICATEPROVIDERTYPE_SIGNED CertificateProviderType = "signed"

// CertificateProviderSigned provides signed certificates.
type CertificateProviderSigned struct {
	CertificateProviderBase
	Keys struct {
		PublicFile  string `json:"publicFile"`
		PrivateFile string `json:"privateFile"`
	} `json:"keys"`
	cert *tls.Certificate
}

// GetType returns the type of the certificate provider.
func (ss *CertificateProviderSigned) GetType() CertificateProviderType {
	return CERTIFICATEPROVIDERTYPE_SIGNED
}

// Validate checks the validity of the signed certificate.
func (ss *CertificateProviderSigned) Validate(dirCache string) error {
	if err := ss.CertificateProviderBase.ValidateDomains(false); err != nil {
		return err
	}

	var publicPath, privatePath string
	if cleanedPath, err := autils.CleanFilePathWithDirOption(ss.Keys.PublicFile, dirCache); err != nil {
		return fmt.Errorf("invalid public key file; %v", err)
	} else {
		publicPath = cleanedPath
	}
	if cleanedPath, err := autils.CleanFilePathWithDirOption(ss.Keys.PrivateFile, dirCache); err != nil {
		return fmt.Errorf("invalid public key file; %v", err)
	} else {
		privatePath = cleanedPath
	}

	if autils.FileExists(publicPath) && autils.FileExists(privatePath) {
		cert, err := LoadCertificateFromFile(publicPath, privatePath)
		if err != nil {
			return err
		}
		ss.cert = cert
	} else {
		return fmt.Errorf("certificate '%s' does not exist", ss.Keys.PublicFile)
	}

	return nil
}

// GetCertificate returns the signed certificate.
func (ss *CertificateProviderSigned) GetCertificate(hello *tls.ClientHelloInfo) (*tls.Certificate, error) {
	return ss.cert, nil
}
