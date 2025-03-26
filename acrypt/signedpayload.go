package acrypt

import (
	"crypto/ed25519"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"
)

// SignedPayload represents the payload signed by the subscriber for identity verification.
type SignedPayload struct {
	Timestamp int64  `json:"timestamp"` // Unix timestamp to prevent replay attacks
	Nonce     string `json:"nonce"`     // Random value to ensure uniqueness
	Signature string `json:"signature"` // Signature of the payload
}

// Sign signs the payload using the private key.
func (sp *SignedPayload) Sign(payload string, privKey ed25519.PrivateKey) (SignedPayload, error) {
	if len(privKey) != ed25519.PrivateKeySize {
		return SignedPayload{}, fmt.Errorf("invalid private key size: expected %d bytes, got %d bytes", ed25519.PrivateKeySize, len(privKey))
	}
	if strings.TrimSpace(payload) == "" {
		return SignedPayload{}, errors.New("payload cannot be empty")
	}

	// Generate the signature
	signature := ed25519.Sign(privKey, []byte(payload))
	return SignedPayload{
		Timestamp: time.Now().Unix(),
		Nonce:     sp.GenerateNonce(),
		Signature: hex.EncodeToString(signature),
	}, nil
}

// Verify verifies the signature of a signed payload using the public key.
func (sp *SignedPayload) Verify(payload string, pubKey ed25519.PublicKey) (bool, error) {
	if len(pubKey) != ed25519.PublicKeySize {
		return false, fmt.Errorf("invalid public key size: expected %d bytes, got %d bytes", ed25519.PublicKeySize, len(pubKey))
	}
	if strings.TrimSpace(payload) == "" {
		return false, errors.New("payload cannot be empty")
	}

	// Decode the signature
	sigBytes, err := hex.DecodeString(sp.Signature)
	if err != nil {
		return false, fmt.Errorf("invalid signature format: %v", err)
	}

	// Verify the signature
	isValid := ed25519.Verify(pubKey, []byte(payload), sigBytes)
	return isValid, nil
}

// ValidateTimestamp checks whether the timestamp is within an acceptable range to prevent replay attacks.
func (sp *SignedPayload) ValidateTimestamp(maxAge time.Duration) error {
	timestampTime := time.Unix(sp.Timestamp, 0)
	if time.Since(timestampTime) > maxAge {
		return errors.New("signed payload is too old")
	}
	if time.Until(timestampTime) > 0 {
		return errors.New("signed payload timestamp is in the future")
	}
	return nil
}

// GenerateNonce generates a unique nonce for payloads.
func (sp *SignedPayload) GenerateNonce() string {
	return fmt.Sprintf("%d-%d", time.Now().UnixNano(), time.Now().Unix())
}
