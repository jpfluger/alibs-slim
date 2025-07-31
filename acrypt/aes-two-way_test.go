package acrypt

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"github.com/stretchr/testify/assert"
	"io"
	"os"
	"path/filepath"
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

func TestAESCTR256EncryptDecrypt_SmallFile(t *testing.T) {
	tmp := t.TempDir()
	plainFile := filepath.Join(tmp, "plain.txt")
	encFile := filepath.Join(tmp, "plain.txt.enc")
	decFile := filepath.Join(tmp, "plain_decrypted.txt")
	passphrase := "supersecret"

	// Write small plain file.
	content := []byte("Hello AES-CTR world! This is a small test file.")
	assert.NoError(t, os.WriteFile(plainFile, content, 0644))

	// Encrypt
	assert.NoError(t, AESCTR256EncryptFile(plainFile, encFile, passphrase))
	assert.FileExists(t, encFile)

	// Decrypt
	assert.NoError(t, AESCTR256DecryptFile(encFile, decFile, passphrase))
	assert.FileExists(t, decFile)

	// Check content matches original.
	decrypted, err := os.ReadFile(decFile)
	assert.NoError(t, err)
	assert.Equal(t, content, decrypted)
}

func TestAESCTR256EncryptDecrypt_LargeFile(t *testing.T) {
	tmp := t.TempDir()
	plainFile := filepath.Join(tmp, "big.bin")
	encFile := filepath.Join(tmp, "big.bin.enc")
	decFile := filepath.Join(tmp, "big_decrypted.bin")
	passphrase := "larger_secret"

	// Generate ~20 MB random file.
	out, err := os.Create(plainFile)
	assert.NoError(t, err)
	defer out.Close()

	buf := make([]byte, 1024*1024)
	for i := 0; i < 20; i++ {
		_, err := rand.Read(buf)
		assert.NoError(t, err)
		_, err = out.Write(buf)
		assert.NoError(t, err)
	}

	// Encrypt
	assert.NoError(t, AESCTR256EncryptFile(plainFile, encFile, passphrase))
	assert.FileExists(t, encFile)

	// Decrypt
	assert.NoError(t, AESCTR256DecryptFile(encFile, decFile, passphrase))
	assert.FileExists(t, decFile)

	// Compare file digests (not full bytes to save RAM)
	origHash, err := fileSHA256(plainFile)
	assert.NoError(t, err)

	decHash, err := fileSHA256(decFile)
	assert.NoError(t, err)

	assert.Equal(t, origHash, decHash)
}

func TestAESCTR256Decrypt_WrongPassphrase(t *testing.T) {
	tmp := t.TempDir()
	plainFile := filepath.Join(tmp, "plain.txt")
	encFile := filepath.Join(tmp, "plain.txt.enc")
	decFile := filepath.Join(tmp, "plain_decrypted.txt")
	passphrase := "correctpass"
	badPassphrase := "wrongpass"

	assert.NoError(t, os.WriteFile(plainFile, []byte("Sensitive stuff!"), 0644))
	assert.NoError(t, AESCTR256EncryptFile(plainFile, encFile, passphrase))

	// Decrypt with wrong passphrase should fail HMAC check.
	err := AESCTR256DecryptFile(encFile, decFile, badPassphrase)
	assert.ErrorContains(t, err, "HMAC does not match")
}

func fileSHA256(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}
