package acrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
)

// Notes
//The potential of quantum computers to crack encryption keys is a significant concern for cybersecurity. Here’s a breakdown of the estimated time it would take for quantum computers to crack AES-128 and AES-256 encryption keys:
//
//AES-128: With the advancements in quantum computing, it is estimated that a sufficiently powerful quantum computer could crack an AES-128 key in about six months1.
//
//AES-256: Due to its longer key length, AES-256 is more resistant to quantum attacks. It is estimated that cracking an AES-256 key with a quantum computer would take about the same time as a classical computer would take to crack an AES-128 key, which is millions of years12.
//
//These estimates highlight the importance of transitioning to quantum-resistant encryption methods as quantum computing technology continues to advance.

// AESGCM128Encrypt encrypts the given data using AES-GCM with a 128-bit key derived from the passphrase.
// It returns the ciphertext, which includes the nonce prepended to the actual encrypted data.
func AESGCM128Encrypt(data []byte, passphrase string) ([]byte, error) {
	// Derive a 128-bit key from the passphrase using SHA-256.
	key := FromStringToSHA256CheckSum(passphrase)

	// Create a new cipher block from the key.
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// Create a new GCM instance which will be used for encryption.
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// Create a nonce with a length recommended by GCM.
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	// Encrypt the data using GCM seal, which will return the combined nonce and ciphertext.
	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return ciphertext, nil
}

// AESGCM128Decrypt decrypts the given ciphertext using AES-GCM with a 128-bit key derived from the passphrase.
// The ciphertext should have the nonce prepended to the actual encrypted data.
func AESGCM128Decrypt(ciphertext []byte, passphrase string) ([]byte, error) {
	// Derive a 128-bit key from the passphrase using SHA-256.
	key := FromStringToSHA256CheckSum(passphrase)

	// Create a new cipher block from the key.
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// Create a new GCM instance which will be used for decryption.
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// Separate the nonce and the actual ciphertext.
	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, err
	}
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	// Decrypt the data using GCM open, which will return the plaintext.
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}
	return plaintext, nil
}

// AESGCM256Encrypt encrypts the given data using AES-GCM with a 256-bit key derived from the passphrase.
// It returns the ciphertext, which includes the nonce prepended to the actual encrypted data.
func AESGCM256Encrypt(data []byte, passphrase string) ([]byte, error) {
	// Derive a 256-bit key from the passphrase using SHA-256.
	key := FromStringToSHA256CheckSum(passphrase)

	// Create a new cipher block from the key.
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// Create a new GCM instance which will be used for encryption.
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// Create a nonce with a length recommended by GCM.
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	// Encrypt the data using GCM seal, which will return the combined nonce and ciphertext.
	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return ciphertext, nil
}

// AESGCM256Decrypt decrypts the given ciphertext using AES-GCM with a 256-bit key derived from the passphrase.
// The ciphertext should have the nonce prepended to the actual encrypted data.
func AESGCM256Decrypt(ciphertext []byte, passphrase string) ([]byte, error) {
	// Derive a 256-bit key from the passphrase using SHA-256.
	key := FromStringToSHA256CheckSum(passphrase)

	// Create a new cipher block from the key.
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// Create a new GCM instance which will be used for decryption.
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// Separate the nonce and the actual ciphertext.
	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, err
	}
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	// Decrypt the data using GCM open, which will return the plaintext.
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}
	return plaintext, nil
}
