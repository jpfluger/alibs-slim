package acrypt

import (
	"strings"
)

/*
Comments on usage of SECRETSKEY_APPHOSTEK and SECRETSKEY_SVT in applications.

SECRETSKEY_APPHOSTEK is the "App Host Encryption Key") and is discovered in a few ways:
1. The root of the service file structure
2. Passed-in by an administrator
3. Returned by remote-provisioning service
It is used to decrypt the running-config or any config saved locally, in a db or remote-provisioning service.

SECRETSKEY_SVT is the "Secrets Vault Token", stored inside the config file.
It is the "master token" used to encrypt/decrypt database fields.
This is a default token but of course more secrets can be created for your purposes.
*/

// SecretsKey constants define specific keys used for encryption and configuration.
const (
	SECRETSKEY_APPHOSTEK SecretsKey = "apphostek" // App Host Encryption Key
	// Used for credential secrets and db passwords
	SECRETSKEY_SVT  SecretsKey = "svt"  // Secrets Vault Token
	SECRETSKEY_SALT SecretsKey = "salt" // Salt value
	SECRETSKEY_JWT  SecretsKey = "jwt"  // JWT signing key
)

// SecretsKey is a type that represents a key used in secret management.
type SecretsKey string

// IsEmpty checks if the SecretsKey is empty after trimming whitespace.
func (sk SecretsKey) IsEmpty() bool {
	return strings.TrimSpace(string(sk)) == ""
}

// TrimSpace trims whitespace from both ends of the SecretsKey.
func (sk SecretsKey) TrimSpace() SecretsKey {
	return SecretsKey(strings.TrimSpace(string(sk)))
}

// String converts the SecretsKey to a regular string.
func (sk SecretsKey) String() string {
	return string(sk)
}

// ToStringTrimLower converts the SecretsKey to a lowercase string with trimmed whitespace.
func (sk SecretsKey) ToStringTrimLower() string {
	return strings.ToLower(strings.TrimSpace(string(sk)))
}

// SecretsKeys is a slice of SecretsKey, representing a collection of keys.
type SecretsKeys []SecretsKey
