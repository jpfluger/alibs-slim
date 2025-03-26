package anetwork

import (
	"crypto/tls"
	"fmt"
	"github.com/asaskevich/govalidator"
	"strings"
)

// ICertificateProvider defines the interface for certificate providers.
type ICertificateProvider interface {
	GetType() CertificateProviderType
	Validate(dirCache string) error
	GetCertificate(hello *tls.ClientHelloInfo) (*tls.Certificate, error)
	GetDomains() []string
	MatchDomain(serverName string) bool
	GetExplicitDomains() []string
	MatchExplicitDomain(serverName string) bool
}

// ICertificateProviders is a slice of ICertificateProvider.
type ICertificateProviders []ICertificateProvider

// CertificateProviderBase provides a base implementation for certificate providers.
type CertificateProviderBase struct {
	Type    CertificateProviderType `json:"type"`
	Domains []string                `json:"domains"`

	// Forces no wildcards.
	// The list for AddDomains should be satisfied by the list in Domains.
	// They are different because in a proxy situation, the Domain list
	// reflects the TLS cert but the AddDomains + Domains = explicitDomains reflects the proxy pathway.
	// For example:
	//   Domains = "example.com", "*.example.com"
	//   AddDomains = "www.example.com" -> which is satisfied by Domains["*.example.com"]
	// During validation, the AddDomains and Domains get the non-wildcard domains merged
	// into explicitDomains.
	AddDomains []string `json:"addDomains,omitempty"`

	// all lower-case, trimmed
	explicitDomains []string
}

// GetType returns the certificate provider type.
func (cpb *CertificateProviderBase) GetType() CertificateProviderType {
	return cpb.Type
}

// GetDomains returns the domains.
func (cpb *CertificateProviderBase) GetDomains() []string {
	return cpb.Domains
}

// GetExplicitDomains returns the Domains plus AddDomains minus wildcards.
func (cpb *CertificateProviderBase) GetExplicitDomains() []string {
	return cpb.explicitDomains
}

// ValidateDomains checks if the domains are valid and merges non-wildcard domains into explicitDomains.
func (cpb *CertificateProviderBase) ValidateDomains(allowIPs bool) error {
	if len(cpb.Domains) == 0 {
		return fmt.Errorf("domains is empty")
	}

	domainSet := make(map[string]struct{})
	for i, domain := range cpb.Domains {
		domain = strings.ToLower(strings.TrimSpace(domain))
		cpb.Domains[i] = domain
		if !IsValidDomain(domain, allowIPs) {
			return fmt.Errorf("invalid domain: %s", domain)
		}
		if !strings.HasPrefix(domain, "*.") {
			domainSet[domain] = struct{}{}
		}
	}

	for i, addDomain := range cpb.AddDomains {
		addDomain = strings.ToLower(strings.TrimSpace(addDomain))
		cpb.AddDomains[i] = addDomain
		if strings.HasPrefix(addDomain, "*.") {
			return fmt.Errorf("addDomains cannot contain wildcards: %s", addDomain)
		}
		if !cpb.MatchDomain(addDomain) {
			return fmt.Errorf("addDomain %s is not satisfied by domains", addDomain)
		}
		domainSet[addDomain] = struct{}{}
	}

	// explicitDomains are not required in Validate but could be within the
	// app at a higher-level.
	cpb.explicitDomains = make([]string, 0, len(domainSet))
	for domain := range domainSet {
		cpb.explicitDomains = append(cpb.explicitDomains, domain)
	}

	return nil
}

// MatchDomain checks if the server name matches any of the domains.
func (cpb *CertificateProviderBase) MatchDomain(serverName string) bool {
	serverName = strings.ToLower(strings.TrimSpace(serverName))
	for _, domain := range cpb.Domains {
		if strings.HasPrefix(domain, "*.") {
			if strings.HasSuffix(serverName, domain[1:]) {
				return true
			}
		} else if serverName == domain {
			return true
		}
	}
	return false
}

// MatchExplicitDomain checks if the given server name matches any of the explicit domains.
// It converts the server name to lowercase and trims any leading or trailing whitespace
// before performing the comparison.
func (cpb *CertificateProviderBase) MatchExplicitDomain(serverName string) bool {
	// Convert the server name to lowercase and trim whitespace
	serverName = strings.ToLower(strings.TrimSpace(serverName))

	// Iterate over the explicit domains and check for a match
	for _, domain := range cpb.explicitDomains {
		if serverName == domain {
			return true
		}
	}

	// Return false if no match is found
	return false
}

// CleanFirstDomain returns the first domain after cleaning it.
func (cpb *CertificateProviderBase) CleanFirstDomain() string {
	if len(cpb.Domains) == 0 {
		return ""
	}
	return CleanDomain(cpb.Domains[0])
}

// IsValidDomain checks if the domain is valid.
func IsValidDomain(domain string, allowIPs bool) bool {
	isValid, _ := IsValidDomainWithError(domain, allowIPs)
	return isValid
}

// IsValidDomainWithError checks if the domain is valid and provides an error description if it is not.
func IsValidDomainWithError(domain string, allowIPs bool) (bool, error) {
	// Check if the domain is empty
	if strings.TrimSpace(domain) == "" {
		return false, fmt.Errorf("domain is empty")
	}

	// Check if the domain is an IP address
	if govalidator.IsIP(domain) {
		if allowIPs {
			return true, nil
		}
		return false, fmt.Errorf("IP addresses are not allowed unless explicitly permitted")
	}

	// Handle wildcard domains
	if strings.HasPrefix(domain, "*.") {
		domain = domain[2:]
	}

	// Ensure the domain contains at least one "."
	isSingleLabelDomain := !strings.Contains(domain, ".")
	if isSingleLabelDomain {
		if !allowIPs {
			return false, fmt.Errorf("single-label domains not allowed")
		}
	}

	// Validate as a DNS name
	if !govalidator.IsDNSName(domain) {
		return false, fmt.Errorf("invalid DNS name format")
	}

	return true, nil
}

// CleanDomain removes wildcard prefix from the domain.
func CleanDomain(domain string) string {
	return strings.Replace(domain, "*.", "", 1)
}
