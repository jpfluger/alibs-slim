package acrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"golang.org/x/crypto/pbkdf2"
	"io"
	"os"
)

// Notes:
// - All AES key sizes (128/192/256 bits) are FIPS 140-3 compliant.
// - AES-256 is recommended (default) for 128-bit quantum-effective security (via Grover's algorithm).
// - PBKDF2 with HMAC-SHA256 is FIPS-approved for key derivation.
// - For FIPS compliance, build with GOEXPERIMENT=systemcrypto or use certified modules.
// - In Go 1.24+, crypto/rand.Read panics on failure (error checks for pre-1.24 compatibility).
// - Monitor NIST for post-quantum updates; extend EncryptionType as needed.

const (
	pbkdf2Iterations = 600000 // OWASP-recommended for PBKDF2-HMAC-SHA256.
	saltSize         = 16     // Standard salt size.
)

// AESGCMEncrypt encrypts data using AES-GCM with the specified EncryptionType.
// Defaults to ENCRYPTIONTYPE_AES256 if invalid/unspecified.
func AESGCMEncrypt(data []byte, passphrase string, encType EncryptionType) ([]byte, error) {
	keySize := encType.KeySize()
	if keySize == 0 {
		encType = ENCRYPTIONTYPE_AES256
		keySize = encType.KeySize()
	}

	salt := make([]byte, saltSize)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return nil, fmt.Errorf("failed to generate salt: %w", err)
	}

	key := pbkdf2.Key([]byte(passphrase), salt, pbkdf2Iterations, keySize, sha256.New)

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create AES cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("failed to generate nonce: %w", err)
	}

	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return append(salt, ciphertext...), nil
}

// AESGCMDecrypt decrypts ciphertext using AES-GCM with the specified EncryptionType.
func AESGCMDecrypt(ciphertext []byte, passphrase string, encType EncryptionType) ([]byte, error) {
	keySize := encType.KeySize()
	if keySize == 0 {
		return nil, fmt.Errorf("invalid encryption type")
	}

	if len(ciphertext) < saltSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	salt, ciphertext := ciphertext[:saltSize], ciphertext[saltSize:]

	key := pbkdf2.Key([]byte(passphrase), salt, pbkdf2Iterations, keySize, sha256.New)

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create AES cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt: %w", err)
	}

	return plaintext, nil
}

// AESCTREncryptFile encrypts a file using AES-CTR + HMAC-SHA256 with the specified EncryptionType.
// Defaults to ENCRYPTIONTYPE_AES256 if invalid/unspecified. HMAC key fixed at 32 bytes.
func AESCTREncryptFile(inputPath, outputPath, passphrase string, encType EncryptionType) error {
	keySize := encType.KeySize()
	if keySize == 0 {
		encType = ENCRYPTIONTYPE_AES256
		keySize = encType.KeySize()
	}

	salt := make([]byte, saltSize)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return fmt.Errorf("failed to generate salt: %w", err)
	}

	keyMaterialSize := keySize + 32 // AES key + HMAC key
	keyMaterial := pbkdf2.Key([]byte(passphrase), salt, pbkdf2Iterations, keyMaterialSize, sha256.New)
	aesKey := keyMaterial[:keySize]
	hmacKey := keyMaterial[keySize:]

	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return fmt.Errorf("failed to generate IV: %w", err)
	}

	in, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("failed to open input file: %w", err)
	}
	defer in.Close()

	out, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer out.Close()

	// Write salt and IV at start.
	if _, err := out.Write(salt); err != nil {
		return fmt.Errorf("failed to write salt: %w", err)
	}
	if _, err := out.Write(iv); err != nil {
		return fmt.Errorf("failed to write IV: %w", err)
	}

	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return fmt.Errorf("failed to create AES cipher: %w", err)
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
				return fmt.Errorf("failed to write ciphertext: %w", err)
			}
			mac.Write(cipherBuf)
		}
		if readErr == io.EOF {
			break
		} else if readErr != nil {
			return fmt.Errorf("failed to read input: %w", readErr)
		}
	}

	if _, err := out.Write(mac.Sum(nil)); err != nil {
		return fmt.Errorf("failed to write HMAC: %w", err)
	}

	return nil
}

// AESCTRDecryptFile decrypts a file using AES-CTR + HMAC-SHA256 with the specified EncryptionType.
// HMAC key fixed at 32 bytes.
func AESCTRDecryptFile(inputPath, outputPath, passphrase string, encType EncryptionType) error {
	keySize := encType.KeySize()
	if keySize == 0 {
		return fmt.Errorf("invalid encryption type")
	}

	in, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("failed to open input file: %w", err)
	}
	defer in.Close()

	out, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer out.Close()

	// Read salt.
	salt := make([]byte, saltSize)
	if _, err := io.ReadFull(in, salt); err != nil {
		return fmt.Errorf("failed to read salt: %w", err)
	}

	// Read IV.
	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(in, iv); err != nil {
		return fmt.Errorf("failed to read IV: %w", err)
	}

	keyMaterialSize := keySize + 32 // AES key + HMAC key
	keyMaterial := pbkdf2.Key([]byte(passphrase), salt, pbkdf2Iterations, keyMaterialSize, sha256.New)
	aesKey := keyMaterial[:keySize]
	hmacKey := keyMaterial[keySize:]

	// Get file size to locate HMAC.
	info, err := in.Stat()
	if err != nil {
		return fmt.Errorf("failed to stat input file: %w", err)
	}

	hmacLen := sha256.Size
	saltIvLen := len(salt) + len(iv)
	contentLen := info.Size() - int64(saltIvLen) - int64(hmacLen)

	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return fmt.Errorf("failed to create AES cipher: %w", err)
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
			return fmt.Errorf("failed to read ciphertext: %w", err)
		}

		cipherBuf := buf[:n]
		plainBuf := make([]byte, n)
		stream.XORKeyStream(plainBuf, cipherBuf)
		if _, err := out.Write(plainBuf); err != nil {
			return fmt.Errorf("failed to write plaintext: %w", err)
		}
		mac.Write(cipherBuf)

		total += int64(n)
	}

	// Read stored HMAC.
	storedMac := make([]byte, hmacLen)
	if _, err := io.ReadFull(in, storedMac); err != nil {
		return fmt.Errorf("failed to read HMAC: %w", err)
	}

	// Verify HMAC.
	expectedMac := mac.Sum(nil)
	if !hmac.Equal(storedMac, expectedMac) {
		return fmt.Errorf("HMAC does not match: file corrupted or wrong passphrase")
	}

	return nil
}

// AESGCM256Encrypt is a wrapper for AESGCMEncrypt using AES-256.
func AESGCM256Encrypt(data []byte, passphrase string) ([]byte, error) {
	return AESGCMEncrypt(data, passphrase, ENCRYPTIONTYPE_AES256)
}

// AESGCM256Decrypt is a wrapper for AESGCMDecrypt using AES-256.
func AESGCM256Decrypt(ciphertext []byte, passphrase string) ([]byte, error) {
	return AESGCMDecrypt(ciphertext, passphrase, ENCRYPTIONTYPE_AES256)
}

// AESCTR256EncryptFile is a wrapper for AESCTREncryptFile using AES-256.
func AESCTR256EncryptFile(inputPath, outputPath, passphrase string) error {
	return AESCTREncryptFile(inputPath, outputPath, passphrase, ENCRYPTIONTYPE_AES256)
}

// AESCTR256DecryptFile is a wrapper for AESCTRDecryptFile using AES-256.
func AESCTR256DecryptFile(inputPath, outputPath, passphrase string) error {
	return AESCTRDecryptFile(inputPath, outputPath, passphrase, ENCRYPTIONTYPE_AES256)
}
