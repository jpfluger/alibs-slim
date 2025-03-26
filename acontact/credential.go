package acontact

import (
	"fmt"
	"strings"
	"time"

	"github.com/jpfluger/alibs-slim/acrypt"
	"github.com/jpfluger/alibs-slim/autils"
	"github.com/jpfluger/alibs-slim/azb"
)

// Credential represents a single credential entry.
type Credential struct {
	// Type identifies the credential's category.
	Type azb.ZBType `json:"type,omitempty"`

	// Label identifies this credential individually rather than relying on "Type".
	Label string `json:"label,omitempty"`

	// Secret holds the encrypted value.
	Secret acrypt.SecretsValue `json:"secret,omitempty"`

	// SecretFile is the path to the secret value, either plaintext or encrypted.
	SecretFile string `json:"secretFile,omitempty"`

	// Username or identifier for the secret. Optional.
	Username string `json:"username"`

	// URL associated with the secret. Optional.
	URL string `json:"url"`

	// Last access timestamp.
	LastAccess *time.Time `json:"lastAccess"`

	// Last modified timestamp.
	Modified time.Time `json:"modified"`
}

// Validate ensures the Credential is correctly initialized and loads secrets if required.
func (ac *Credential) Validate() error {
	if ac == nil {
		return fmt.Errorf("credential is nil")
	}

	// Validate the type field.
	if ac.Type.IsEmpty() {
		return fmt.Errorf("credential type is empty")
	}

	// Trim whitespace from the label.
	ac.Label = strings.TrimSpace(ac.Label)

	// Set the modified timestamp if it's not already set.
	if ac.Modified.IsZero() {
		ac.Modified = time.Now().UTC()
	}

	// Load and validate the secret if `Secret` is empty but `SecretFile` is provided.
	if !ac.Secret.HasValue() && strings.TrimSpace(ac.SecretFile) != "" {
		if err := ac.loadSecretFromFile(); err != nil {
			return err
		}
	}

	if !ac.Secret.HasValue() {
		return fmt.Errorf("credential secret is empty")
	}

	return nil
}

// loadSecretFromFile reads and validates the secret from the provided file path.
func (ac *Credential) loadSecretFromFile() error {
	secretData := autils.ReadFileTrimSpace(ac.SecretFile)
	if secretData == "" {
		return fmt.Errorf("failed to read secret from file %s", ac.SecretFile)
	}

	// Validate the secret with the loaded data.
	ac.Secret.Value.Validate(secretData)
	ac.SecretFile = "" // Clear SecretFile after loading.
	return nil
}

// GetSecret returns the decoded secret.
func (ac *Credential) GetSecret(masterPassword string) ([]byte, error) {
	// Decode the secret and cache the result.
	b, err := ac.Secret.Decode(masterPassword, false)
	if err != nil {
		return nil, fmt.Errorf("failed to decode secret: %v", err)
	}
	return b, nil
}

// SetSecret sets the new secret, encrypting it.
func (ac *Credential) SetSecret(masterPassword string, secret []byte) error {
	if masterPassword == "" {
		return fmt.Errorf("master password is empty")
	}
	// Decode the secret and cache the result.
	if err := ac.Secret.Value.Encode(secret, masterPassword); err != nil {
		return fmt.Errorf("failed to encode secret: %v", err)
	}
	return nil
}

// Credentials represents a list of Credential objects.
type Credentials []*Credential

// Clean validates and filters out invalid credentials.
func (acs Credentials) Clean() (Credentials, error) {
	var validCredentials Credentials
	if acs == nil || len(acs) == 0 {
		return validCredentials, nil
	}

	for _, ac := range acs {
		if err := ac.Validate(); err == nil {
			validCredentials = append(validCredentials, ac)
		}
	}

	return validCredentials, nil
}
