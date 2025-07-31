package acrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"golang.org/x/crypto/argon2"
	"io"
	"os"
)

// Notes
// The potential threat of quantum computers to symmetric encryption like AES is through Grover's algorithm, which provides a quadratic speedup.
// For AES-128, this reduces the effective security to 64 bits, making it potentially vulnerable to future quantum computers.
// For AES-256, it reduces to 128 bits of security, which is still considered secure against quantum attacks for the foreseeable future.
// Recent claims in 2024 about quantum attacks on encryption highlight the need for caution, but AES-256 remains robust.
// These insights emphasize the importance of using stronger keys and considering quantum-resistant algorithms as technology advances.

// AESGCM128Encrypt encrypts the given data using AES-GCM with a 128-bit key derived from the passphrase using Argon2.
// It returns the ciphertext, which includes the salt prepended to the nonce and encrypted data.
func AESGCM128Encrypt(data []byte, passphrase string) ([]byte, error) {
	salt := make([]byte, 16)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return nil, err
	}

	key := argon2.IDKey([]byte(passphrase), salt, 1, 64*1024, 4, 16)

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return append(salt, ciphertext...), nil
}

// AESGCM128Decrypt decrypts the given ciphertext using AES-GCM with a 128-bit key derived from the passphrase using Argon2.
// The ciphertext should have the salt prepended to the nonce and encrypted data.
func AESGCM128Decrypt(ciphertext []byte, passphrase string) ([]byte, error) {
	saltSize := 16
	if len(ciphertext) < saltSize {
		return nil, fmt.Errorf("ciphertext too short")
	}
	salt, ciphertext := ciphertext[:saltSize], ciphertext[saltSize:]

	key := argon2.IDKey([]byte(passphrase), salt, 1, 64*1024, 4, 16)

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}
	return plaintext, nil
}

// AESGCM256Encrypt encrypts the given data using AES-GCM with a 256-bit key derived from the passphrase using Argon2.
// It returns the ciphertext, which includes the salt prepended to the nonce and encrypted data.
func AESGCM256Encrypt(data []byte, passphrase string) ([]byte, error) {
	salt := make([]byte, 16)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return nil, err
	}

	key := argon2.IDKey([]byte(passphrase), salt, 1, 64*1024, 4, 32)

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return append(salt, ciphertext...), nil
}

// AESGCM256Decrypt decrypts the given ciphertext using AES-GCM with a 256-bit key derived from the passphrase using Argon2.
// The ciphertext should have the salt prepended to the nonce and encrypted data.
func AESGCM256Decrypt(ciphertext []byte, passphrase string) ([]byte, error) {
	saltSize := 16
	if len(ciphertext) < saltSize {
		return nil, fmt.Errorf("ciphertext too short")
	}
	salt, ciphertext := ciphertext[:saltSize], ciphertext[saltSize:]

	key := argon2.IDKey([]byte(passphrase), salt, 1, 64*1024, 4, 32)

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}
	return plaintext, nil
}

// AESCTR256EncryptFile encrypts an input file using AES-CTR + HMAC for streaming with keys derived using Argon2.
func AESCTR256EncryptFile(inputPath, outputPath, passphrase string) error {
	salt := make([]byte, 16)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return fmt.Errorf("failed to generate salt: %w", err)
	}

	keyMaterial := argon2.IDKey([]byte(passphrase), salt, 1, 64*1024, 4, 64)
	aesKey := keyMaterial[:32]
	hmacKey := keyMaterial[32:]

	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return fmt.Errorf("failed to generate IV: %w", err)
	}

	in, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write salt and IV at start.
	if _, err := out.Write(salt); err != nil {
		return err
	}
	if _, err := out.Write(iv); err != nil {
		return err
	}

	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return err
	}
	stream := cipher.NewCTR(block, iv)

	mac := hmac.New(sha256.New, hmacKey)
	mac.Write(salt)
	mac.Write(iv)

	buf := make([]byte, 32*1024)
	for {
		n, readErr := in.Read(buf)
		if n > 0 {
			plain := buf[:n]
			cipherBuf := make([]byte, n)
			stream.XORKeyStream(cipherBuf, plain)

			if _, err := out.Write(cipherBuf); err != nil {
				return err
			}

			mac.Write(cipherBuf)
		}

		if readErr == io.EOF {
			break
		} else if readErr != nil {
			return readErr
		}
	}

	if _, err := out.Write(mac.Sum(nil)); err != nil {
		return err
	}

	return nil
}

// AESCTR256DecryptFile decrypts an input file using AES-CTR + HMAC for streaming with keys derived using Argon2.
func AESCTR256DecryptFile(inputPath, outputPath, passphrase string) error {
	in, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Read salt.
	salt := make([]byte, 16)
	if _, err := io.ReadFull(in, salt); err != nil {
		return fmt.Errorf("failed to read salt: %w", err)
	}

	// Read IV.
	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(in, iv); err != nil {
		return fmt.Errorf("failed to read IV: %w", err)
	}

	keyMaterial := argon2.IDKey([]byte(passphrase), salt, 1, 64*1024, 4, 64)
	aesKey := keyMaterial[:32]
	hmacKey := keyMaterial[32:]

	// Get file size to locate HMAC.
	info, err := in.Stat()
	if err != nil {
		return err
	}
	hmacLen := sha256.Size
	saltIvLen := len(salt) + len(iv)
	contentLen := info.Size() - int64(saltIvLen) - int64(hmacLen)

	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return err
	}
	stream := cipher.NewCTR(block, iv)

	mac := hmac.New(sha256.New, hmacKey)
	mac.Write(salt)
	mac.Write(iv)

	buf := make([]byte, 32*1024)
	var total int64
	for total < contentLen {
		toRead := int64(len(buf))
		if contentLen-total < toRead {
			toRead = contentLen - total
		}

		n, err := in.Read(buf[:toRead])
		if err != nil && err != io.EOF {
			return err
		}

		cipherBuf := buf[:n]
		plainBuf := make([]byte, n)
		stream.XORKeyStream(plainBuf, cipherBuf)

		if _, err := out.Write(plainBuf); err != nil {
			return err
		}

		mac.Write(cipherBuf)
		total += int64(n)
	}

	// Read stored HMAC.
	storedMac := make([]byte, hmacLen)
	if _, err := io.ReadFull(in, storedMac); err != nil {
		return err
	}

	// Verify.
	expectedMac := mac.Sum(nil)
	if !hmac.Equal(storedMac, expectedMac) {
		return fmt.Errorf("HMAC does not match: file corrupted or wrong passphrase")
	}

	return nil
}
