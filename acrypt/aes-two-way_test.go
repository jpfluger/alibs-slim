package acrypt

import (
	"bytes"
	"testing"
)

func TestAESGCM128EncryptDecrypt(t *testing.T) {
	// Define a passphrase for encryption/decryption.
	passphrase := "my-secret-passphrase"

	// Define the plaintext data to encrypt.
	plaintext := []byte("Hello, World!")

	// Encrypt the plaintext.
	ciphertext, err := AESGCM128Encrypt(plaintext, passphrase)
	if err != nil {
		t.Fatalf("Failed to encrypt: %v", err)
	}

	// Ensure that the ciphertext is not nil and not equal to the plaintext.
	if ciphertext == nil || bytes.Equal(ciphertext, plaintext) {
		t.Errorf("Encryption failed: ciphertext is nil or equal to plaintext")
	}

	// Decrypt the ciphertext.
	decryptedText, err := AESGCM128Decrypt(ciphertext, passphrase)
	if err != nil {
		t.Fatalf("Failed to decrypt: %v", err)
	}

	// Ensure that the decrypted text matches the original plaintext.
	if !bytes.Equal(decryptedText, plaintext) {
		t.Errorf("Decryption failed: decrypted text does not match original plaintext")
	}
}

func TestAESGCM128DecryptWithWrongPassphrase(t *testing.T) {
	// Define a passphrase for encryption.
	passphrase := "my-secret-passphrase"

	// Define a wrong passphrase for decryption.
	wrongPassphrase := "wrong-passphrase"

	// Define the plaintext data to encrypt.
	plaintext := []byte("Hello, World!")

	// Encrypt the plaintext with the correct passphrase.
	ciphertext, err := AESGCM128Encrypt(plaintext, passphrase)
	if err != nil {
		t.Fatalf("Failed to encrypt: %v", err)
	}

	// Attempt to decrypt the ciphertext with the wrong passphrase.
	_, err = AESGCM128Decrypt(ciphertext, wrongPassphrase)
	if err == nil {
		t.Errorf("Expected an error when decrypting with the wrong passphrase, but got nil")
	}
}

func TestAESGCM256EncryptDecrypt(t *testing.T) {
	passphrase := "strongpassphrase"
	plaintext := []byte("This is a secret message.")

	// Encrypt the plaintext
	ciphertext, err := AESGCM256Encrypt(plaintext, passphrase)
	if err != nil {
		t.Fatalf("Encryption failed: %v", err)
	}

	// Decrypt the ciphertext
	decryptedText, err := AESGCM256Decrypt(ciphertext, passphrase)
	if err != nil {
		t.Fatalf("Decryption failed: %v", err)
	}

	// Check if the decrypted text matches the original plaintext
	if !bytes.Equal(decryptedText, plaintext) {
		t.Errorf("Decrypted text does not match the original plaintext.\nGot: %s\nWant: %s", decryptedText, plaintext)
	}
}

func TestAESGCM256EncryptDecryptWithDifferentPassphrase(t *testing.T) {
	passphrase := "strongpassphrase"
	wrongPassphrase := "wrongpassphrase"
	plaintext := []byte("This is a secret message.")

	// Encrypt the plaintext
	ciphertext, err := AESGCM256Encrypt(plaintext, passphrase)
	if err != nil {
		t.Fatalf("Encryption failed: %v", err)
	}

	// Attempt to decrypt the ciphertext with a wrong passphrase
	_, err = AESGCM256Decrypt(ciphertext, wrongPassphrase)
	if err == nil {
		t.Fatal("Decryption should have failed with a wrong passphrase")
	}
}

func TestAESGCM256EncryptEmptyData(t *testing.T) {
	passphrase := "strongpassphrase"
	plaintext := []byte("")

	// Encrypt the empty plaintext
	ciphertext, err := AESGCM256Encrypt(plaintext, passphrase)
	if err != nil {
		t.Fatalf("Encryption failed: %v", err)
	}

	// Decrypt the ciphertext
	decryptedText, err := AESGCM256Decrypt(ciphertext, passphrase)
	if err != nil {
		t.Fatalf("Decryption failed: %v", err)
	}

	// Check if the decrypted text matches the original plaintext
	if !bytes.Equal(decryptedText, plaintext) {
		t.Errorf("Decrypted text does not match the original plaintext.\nGot: %s\nWant: %s", decryptedText, plaintext)
	}
}
