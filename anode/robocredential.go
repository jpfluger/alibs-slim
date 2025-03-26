package anode

import (
	"crypto/ed25519"
	"errors"
	"fmt"
	"github.com/jpfluger/alibs-slim/acrypt"
	"sync"
	"time"
)

// RoboCredential securely manages credentials for hybrid Pub-Priv Key and JWT authentication.
type RoboCredential struct {
	PublicKey        acrypt.CryptKeyBase64 `json:"publicKey,omitempty"`    // Subscriber's public key
	PrivateKey       acrypt.SecretsValue   `json:"privateKey"`             // Secure storage for private key
	AccessToken      string                `json:"accessToken,omitempty"`  // JWT Access Token
	RefreshToken     string                `json:"refreshToken,omitempty"` // JWT Refresh Token
	TokenExpiresAt   time.Time             `json:"-"`                      // Expiration time for the access token
	RefreshExpiresAt time.Time             `json:"-"`                      // Expiration time for the refresh token
	IsTokenValid     bool                  `json:"-"`                      // Tracks token validity
	mu               sync.RWMutex          // Thread-safe access
}

// Validate ensures that all required fields are present and valid.
// Required fields are those needed for two peers to make a connection.
// The PrivateKey check only ensures a value is present and doesn't try
// to decode the key.
func (rc *RoboCredential) Validate() error {
	if rc == nil {
		return errors.New("credential is nil")
	}

	rc.mu.RLock()
	defer rc.mu.RUnlock()

	// Check if PublicKey is present
	if len(rc.PublicKey) == 0 {
		return fmt.Errorf("public key is missing")
	}
	if rc.PublicKey.IsBase64() {
		return fmt.Errorf("public key is note base64")
	}

	// Check if PrivateKey is present
	if !rc.PrivateKey.HasValue() {
		return fmt.Errorf("private key is missing")
	}

	return nil
}

func (rc *RoboCredential) GenerateKeyPair(masterPassword string, durationMinutes int) error {
	pubKey, privKey, err := ed25519.GenerateKey(nil)
	if err != nil {
		return fmt.Errorf("failed to generate key pair: %v", err)
	}

	rc.mu.Lock()
	defer rc.mu.Unlock()

	// Store the public key as base64
	rc.PublicKey = acrypt.NewCryptKeyBase64(pubKey)

	// Set a default duration if the input is invalid
	if durationMinutes < 0 {
		durationMinutes = 0
	}

	// Use NewSecretsValueRawBase64Decrypted to create the SecretsValueRaw
	rawValue := acrypt.NewSecretsValueRawBase64Decrypted(acrypt.ENCRYPTIONTYPE_AES256, privKey)
	rc.PrivateKey = acrypt.SecretsValue{
		Value: rawValue,
	}

	// Set expiration duration
	if durationMinutes > 0 {
		rc.PrivateKey.Rotate(string(privKey), time.Duration(durationMinutes)*time.Minute)
	}

	// Encrypt the private key immediately
	if err := rc.PrivateKey.EnsureCryptMode(masterPassword, acrypt.CRYPTMODE_ENCRYPTED); err != nil {
		return fmt.Errorf("failed to encrypt private key: %v", err)
	}
	return nil
}

// GetDecodedPrivateKey retrieves the decoded private key from secure storage.
func (rc *RoboCredential) GetDecodedPrivateKey(masterPassword string) ([]byte, error) {
	rc.mu.RLock()
	defer rc.mu.RUnlock()

	// Ensure the key is decrypted before retrieval.
	if !rc.PrivateKey.IsDecoded() {
		_, err := rc.PrivateKey.Decode(masterPassword, true)
		if err != nil {
			return []byte{}, fmt.Errorf("failed to decode private key: %v", err)
		}
	}
	return rc.PrivateKey.GetDecoded(), nil
}

func (rc *RoboCredential) RotateKeys(masterPassword string) error {
	rc.mu.Lock()
	defer rc.mu.Unlock()

	// Generate new key pair
	pubKey, privKey, err := ed25519.GenerateKey(nil)
	if err != nil {
		return fmt.Errorf("failed to rotate keys: %v", err)
	}

	// Store new public key as base64
	rc.PublicKey = acrypt.NewCryptKeyBase64(pubKey)

	// Encode private key in SecretsValueRaw with base64
	rawValue := acrypt.NewSecretsValueRawBase64Decrypted(acrypt.ENCRYPTIONTYPE_AES256, privKey)
	rc.PrivateKey = acrypt.SecretsValue{
		Value: rawValue,
	}

	// Encrypt the private key
	if err = rc.PrivateKey.EnsureCryptMode(masterPassword, acrypt.CRYPTMODE_ENCRYPTED); err != nil {
		return fmt.Errorf("failed to encrypt new private key: %v", err)
	}

	// Invalidate tokens
	rc.AccessToken = ""
	rc.RefreshToken = ""
	rc.TokenExpiresAt = time.Time{}
	rc.RefreshExpiresAt = time.Time{}
	rc.IsTokenValid = false

	return nil
}

// ValidateToken checks whether the access token is still valid.
func (rc *RoboCredential) ValidateToken() bool {
	rc.mu.RLock()
	defer rc.mu.RUnlock()

	return rc.IsTokenValid && time.Now().Before(rc.TokenExpiresAt)
}

// InvalidateTokens explicitly invalidates all JWT tokens.
func (rc *RoboCredential) InvalidateTokens() {
	rc.mu.Lock()
	defer rc.mu.Unlock()

	rc.AccessToken = ""
	rc.RefreshToken = ""
	rc.TokenExpiresAt = time.Time{}
	rc.RefreshExpiresAt = time.Time{}
	rc.IsTokenValid = false
}

// RefreshAccessToken uses the refresh token to obtain a new access token.
func (rc *RoboCredential) RefreshAccessToken(refreshFunc func(string) (string, time.Time, error)) error {
	rc.mu.Lock()
	defer rc.mu.Unlock()

	if time.Now().After(rc.RefreshExpiresAt) {
		return fmt.Errorf("refresh token is expired")
	}

	newAccessToken, newExpiresAt, err := refreshFunc(rc.RefreshToken)
	if err != nil {
		return fmt.Errorf("failed to refresh access token: %v", err)
	}

	rc.AccessToken = newAccessToken
	rc.TokenExpiresAt = newExpiresAt
	rc.IsTokenValid = true
	return nil
}
