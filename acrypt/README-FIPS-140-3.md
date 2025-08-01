# FIPS 140-3 Compliance

## Overview

This package (`acrypt`) provides cryptographic utilities designed to be compliant with FIPS 140-3 standards when built and configured appropriately. FIPS 140-3 is a U.S. government standard for cryptographic modules, specifying requirements for secure key generation, encryption, hashing, and other operations. The provided functions use FIPS-approved algorithms such as AES (in GCM and CTR modes), PBKDF2 (with HMAC-SHA-256), SHA-256, SHA-512, ECDSA (with P-256/P-384), RSA-PSS (2048/3072 bits), RSA-OAEP, and secure random number generation via `crypto/rand`.

These utilities are suitable for applications requiring FIPS compliance, such as government systems or regulated industries. Note that full compliance depends on the build environment, runtime configuration, and certification of the underlying cryptographic modules (e.g., via NIST-validated implementations).

As of August 1, 2025, the Go Cryptographic Module (by Geomys LLC) for FIPS 140-3 is in "Review Pending" status with NIST's CMVP (submitted May 8, 2025; CAVP cert #A6650). For details, see the Go blog post on the FIPS 140-3 Go Cryptographic Module (July 15, 2025): https://go.dev/blog/fips140.

Key files:
- `fips3-aes.go`: AES-based encryption and decryption.
- `fips3-hashes.go`: SHA-256 and SHA-512 hashing functions.
- `fips3-passwords.go`: Password hashing with automatic FIPS mode detection (switches to PBKDF2).
- `fips3-secrets.go`: Secure secret key generation.
- `fips3-signing.go`: Asymmetric signing (ECDSA, RSA-PSS) and encryption (RSA-OAEP).

## Features

- **FIPS-Approved Algorithms**:
    - Encryption: AES-128/192/256 (GCM for in-memory data, CTR with HMAC-SHA-256 for files).
    - Asymmetric Signing: ECDSA with P-256/P-384 curves (FIPS 186-5), RSA-PSS with 2048/3072-bit keys and SHA-256/384/512 (SP 800-186).
    - Asymmetric Encryption: RSA-OAEP with SHA-256 (SP 800-56B).
    - Key Derivation: PBKDF2 with HMAC-SHA-256 (600,000 iterations minimum, per OWASP recommendations).
    - Hashing: SHA-256 and SHA-512 for checksums and integrity checks.
    - Randomness: High-entropy generation using `crypto/rand`.

- **Automatic FIPS Mode Switching**:
    - In password hashing, detects FIPS mode (via `crypto/fips140.Enabled()`) and uses PBKDF2 instead of Argon2id.

- **Output Formats**:
    - Hashes in hex, base64, or PHC format.
    - Encrypted data includes salts, nonces/IVs, and authentication tags.
    - Keys encoded in PEM format for asymmetric operations.

- **Error Handling**:
    - Functions return errors for invalid inputs, failed randomness, or decryption failures (e.g., wrong passphrase).

- **Compatibility Notes**:
    - Supports Go 1.23+ (with considerations for `crypto/rand` behavior in Go 1.24+).
    - For post-quantum security, AES-256 is recommended (quantum-effective against Grover's algorithm). Monitor NIST for post-quantum algorithms like ML-DSA.

## Enabling FIPS Compliance

To achieve FIPS 140-3 compliance:
1. **Go 1.24+ Native Module (Recommended)**: Enable the built-in FIPS 140-3 validated cryptographic module via `GODEBUG=fips140=on`. This uses pure-Go implementations with no cgo or performance overhead, as detailed in the Go blog (https://go.dev/blog/fips140).
   ```
   GODEBUG=fips140=on go run yourprogram.go
   ```
   For stricter enforcement, use `GODEBUG=fips140=only` to panic on non-FIPS usage.

2. **Build with System Crypto (Alternative)**: Use the build flag `GOEXPERIMENT=systemcrypto` to leverage OS-provided certified cryptographic modules (e.g., on platforms with FIPS-validated libraries like OpenSSL or BoringSSL).
   ```
   GOEXPERIMENT=systemcrypto go build
   ```

3. **Runtime Detection**: The package uses `crypto/fips140.Enabled()` to check if FIPS mode is active. This can be overridden for testing via the `IsFIPSMode` global variable.

Ensure your environment (OS, hardware) supports FIPS mode. Consult NIST guidelines for full certification.

## Dependencies

- Standard Go crypto libraries: `crypto/aes`, `crypto/cipher`, `crypto/ecdsa`, `crypto/hmac`, `crypto/rand`, `crypto/rsa`, `crypto/sha256`, `crypto/sha512`, `crypto/x509`, `crypto`, `math/big`.
- External: `golang.org/x/crypto/pbkdf2`, `crypto/fips140`.
- No additional installations required beyond Go standard library and x/crypto.

## Usage Examples

### Password Hashing
```go
import "acrypt"

// Hash a password (uses PBKDF2 in FIPS mode)
hashed, err := acrypt.HashPassword("securepassword", nil)
if err != nil {
    // Handle error
}

// Verify password
match, err := acrypt.MatchPassword(hashed, "securepassword")
if err != nil {
    // Handle error
}
if match {
    fmt.Println("Password matches")
}
```

### AES Encryption (In-Memory)
```go
data := []byte("sensitive data")
passphrase := "strongpassphrase"

encrypted, err := acrypt.AESGCM256Encrypt(data, passphrase)
if err != nil {
    // Handle error
}

decrypted, err := acrypt.AESGCM256Decrypt(encrypted, passphrase)
if err != nil {
    // Handle error
}
```

### File Encryption
```go
err := acrypt.AESCTR256EncryptFile("input.txt", "output.enc", "passphrase")
if err != nil {
    // Handle error
}

err = acrypt.AESCTR256DecryptFile("output.enc", "decrypted.txt", "passphrase")
if err != nil {
    // Handle error
}
```

### Hashing
```go
checksum := acrypt.FromStringToSHA256CheckSum("data")
hexChecksum := acrypt.FormatSHA256ChecksumHex(checksum)

base64Hash := acrypt.ToHashSHA256Base64("data", true) // Prepends "{sha256}"
```

### Secret Key Generation
```go
key, err := acrypt.GenerateSecretKey()
if err != nil {
    // Handle error
}

encoded := acrypt.EncodeJWTSecretKey(key)
```

### Asymmetric Signing (ECDSA)
```go
priv, err := acrypt.GenerateECDSA256Key()
if err != nil {
    // Handle error
}

data := []byte("data to sign")
sig, err := acrypt.SignWithECDSA(priv, data)
if err != nil {
    // Handle error
}

// Verify
valid, err := acrypt.VerifyECDSASignature(&priv.PublicKey, data, sig)
if err != nil {
    // Handle error
}
if valid {
    fmt.Println("Signature valid")
}
```

### Asymmetric Signing (RSA-PSS)
```go
priv, err := acrypt.GenerateRSA2048Key()
if err != nil {
    // Handle error
}

data := []byte("data to sign")
sig, err := acrypt.SignWithRSAPSS(priv, data, crypto.SHA256)
if err != nil {
    // Handle error
}

// Verify
err = acrypt.VerifyRSAPSSSignature(&priv.PublicKey, data, sig, crypto.SHA256)
if err == nil {
    fmt.Println("Signature valid")
}
```

### Asymmetric Encryption (RSA-OAEP)
```go
priv, err := acrypt.GenerateRSA2048Key()
if err != nil {
    // Handle error
}

pub := &priv.PublicKey
data := []byte("secret")

enc, err := acrypt.EncryptWithRSAOAEP(pub, data)
if err != nil {
    // Handle error
}

dec, err := acrypt.DecryptWithRSAOAEP(priv, enc)
if err != nil {
    // Handle error
}
```

## Configuration Options

- **PBKDF2 Presets** (for passwords in FIPS mode):
  ```go
  presets := acrypt.NewPBKDF2Presets()
  presets.Iterations = 1000000 // Increase for higher security
  presets.HashFunc = "sha512"
  hashed, _ := acrypt.HashPBKDF2("password", presets)
  ```

- **Encryption Types**:
    - `ENCRYPTIONTYPE_AES128`, `ENCRYPTIONTYPE_AES192`, `ENCRYPTIONTYPE_AES256` (default).

- **Signing Types**:
    - `SigningTypeECDSA256`, `SigningTypeECDSA384`, `SigningTypeRSAPSS2048`, `SigningTypeRSAPSS3072`.

## Limitations and Notes

- **Not for Password Storage Alone**: Use PBKDF2 for passwords, but combine with salting and peppering for best practices.
- **No Argon2 in FIPS Mode**: Automatically falls back to PBKDF2; Argon2id is used otherwise.
- **File Operations**: AES-CTR with HMAC is streaming-friendly but requires full file reads for integrity checks.
- **Asymmetric Choices**: Prefer ECDSA over RSA for better performance and security in new applications; RSA is included for legacy compatibility.
- **Quantum Considerations**: Monitor NIST for post-quantum algorithms; extend `EncryptionType` as needed.
- **Testing**: Override `acrypt.IsFIPSMode` for unit tests to simulate FIPS environments.
- **Compliance Disclaimer**: This package aids compliance but does not guarantee it. Obtain formal validation for your application if required.

For issues or contributions, refer to the repository guidelines.