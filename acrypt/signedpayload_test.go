package acrypt

import (
	"crypto/ed25519"
	"crypto/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSignPayload(t *testing.T) {
	// Generate key pair
	pubKey, privKey, err := ed25519.GenerateKey(rand.Reader)
	assert.NoError(t, err)
	assert.NotEmpty(t, pubKey)
	assert.NotEmpty(t, privKey)

	payload := "example-payload"
	sp := &SignedPayload{}

	// Sign the payload
	signedPayload, err := sp.Sign(payload, privKey)
	assert.NoError(t, err)
	assert.NotEmpty(t, signedPayload.Signature)
	assert.NotEmpty(t, signedPayload.Nonce)
	assert.NotZero(t, signedPayload.Timestamp)

	// Ensure signature is valid
	isValid, err := signedPayload.Verify(payload, pubKey)
	assert.NoError(t, err)
	assert.True(t, isValid)
}

func TestSignPayloadWithInvalidPrivateKey(t *testing.T) {
	payload := "example-payload"
	invalidPrivKey := []byte("invalid-key")

	sp := &SignedPayload{}
	_, err := sp.Sign(payload, invalidPrivKey)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid private key size")
}

func TestVerifyPayloadSignatureWithInvalidPublicKey(t *testing.T) {
	// Generate valid key pair
	_, privKey, _ := ed25519.GenerateKey(rand.Reader)
	payload := "example-payload"

	sp := &SignedPayload{}
	signedPayload, err := sp.Sign(payload, privKey)
	assert.NoError(t, err)

	invalidPubKey := []byte("invalid-public-key")

	// Verify with invalid public key
	isValid, err := signedPayload.Verify(payload, invalidPubKey)
	assert.Error(t, err)
	assert.False(t, isValid)
	assert.Contains(t, err.Error(), "invalid public key size")
}

func TestVerifyPayloadSignatureWithAlteredPayload(t *testing.T) {
	// Generate valid key pair
	pubKey, privKey, _ := ed25519.GenerateKey(rand.Reader)
	originalPayload := "example-payload"
	alteredPayload := "tampered-payload"

	sp := &SignedPayload{}
	signedPayload, err := sp.Sign(originalPayload, privKey)
	assert.NoError(t, err)

	// Verify with altered payload
	isValid, err := signedPayload.Verify(alteredPayload, pubKey)
	assert.NoError(t, err)
	assert.False(t, isValid)
}

func TestValidateTimestamp(t *testing.T) {
	sp := &SignedPayload{
		Timestamp: time.Now().Unix(),
	}

	// Test with valid timestamp
	err := sp.ValidateTimestamp(5 * time.Minute)
	assert.NoError(t, err)

	// Test with old timestamp
	sp.Timestamp = time.Now().Add(-10 * time.Minute).Unix()
	err = sp.ValidateTimestamp(5 * time.Minute)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "signed payload is too old")

	// Test with future timestamp
	sp.Timestamp = time.Now().Add(10 * time.Minute).Unix()
	err = sp.ValidateTimestamp(5 * time.Minute)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "timestamp is in the future")
}

func TestGenerateNonce(t *testing.T) {
	sp := &SignedPayload{
		Timestamp: time.Now().Unix(),
	}
	nonce1 := sp.GenerateNonce()
	nonce2 := sp.GenerateNonce()

	assert.NotEmpty(t, nonce1)
	assert.NotEmpty(t, nonce2)
	assert.NotEqual(t, nonce1, nonce2)
}
