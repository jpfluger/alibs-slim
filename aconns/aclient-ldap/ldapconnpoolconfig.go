package aclient_ldap

import (
	"crypto/tls"
	"fmt"
)

// Helper struct to implement the LdapConfig interface for AClientLDAP
type clientLdapConfig struct {
	client *AClientLDAP
}

func (c *clientLdapConfig) GetURL() string {
	if c.client.UseSSL {
		return fmt.Sprintf("ldaps://%s", c.client.getAddress())
	}
	return fmt.Sprintf("ldap://%s", c.client.getAddress())
}

func (c *clientLdapConfig) GetAdminDN() string {
	return c.client.BindDN
}

func (c *clientLdapConfig) GetAdminPass() string {
	return c.client.GetPassword()
}

func (c *clientLdapConfig) GetMaxOpen() int {
	return 10 // Return configurable max open connections if needed.
}

func (c *clientLdapConfig) GetMaxDialerTimeout() int {
	return c.client.ConnectionTimeout
}

func (c *clientLdapConfig) GetUseSSL() bool {
	return c.client.UseSSL
}

func (c *clientLdapConfig) GetNetwork() string {
	return "tcp"
}

func (c *clientLdapConfig) GetAddress() string {
	return c.client.getAddress()
}

func (c *clientLdapConfig) GetHost() string {
	return c.client.Host
}

func (c *clientLdapConfig) GetSkipTLS() bool {
	return c.client.SkipTLS
}

func (c *clientLdapConfig) GetInsecureSkipVerify() bool {
	return c.client.InsecureSkipVerify
}

func (c *clientLdapConfig) GetTLSCertificates() []tls.Certificate {
	return c.client.ClientCertificates
}
