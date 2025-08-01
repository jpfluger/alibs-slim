package acrypt

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/sha512"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"hash"
	"math/big"
)

// Notes:
// - ECDSA with NIST P-256/P-384 curves is FIPS 186-5 compliant for signing.
// - RSA with 2048/3072-bit keys, PSS for signing (SP 800-186), OAEP for encryption (SP 800-56B).
// - Use SHA-256 as default hash; SHA-512 optional.
// - Key generation uses crypto/rand (DRBG in FIPS mode).
// - For FIPS compliance, enable via GODEBUG=fips140=on in Go 1.24+.
// - Private keys returned as PEM-encoded for storage/security; parse as needed.
// - Monitor NIST for post-quantum transitions (e.g., ML-DSA in future).

// SigningType defines the asymmetric signing algorithm.
type SigningType int

const (
	SigningTypeECDSA256   SigningType = iota // ECDSA with P-256
	SigningTypeECDSA384                      // ECDSA with P-384
	SigningTypeRSAPSS2048                    // RSA-PSS with 2048-bit key
	SigningTypeRSAPSS3072                    // RSA-PSS with 3072-bit key
)

// GenerateECDSAKey generates an ECDSA private key for the specified curve.
func GenerateECDSAKey(curve elliptic.Curve) (*ecdsa.PrivateKey, error) {
	priv, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("failed to generate ECDSA key: %w", err)
	}
	return priv, nil
}

// GenerateECDSA256Key is a wrapper for P-256.
func GenerateECDSA256Key() (*ecdsa.PrivateKey, error) {
	return GenerateECDSAKey(elliptic.P256())
}

// GenerateECDSA384Key is a wrapper for P-384.
func GenerateECDSA384Key() (*ecdsa.PrivateKey, error) {
	return GenerateECDSAKey(elliptic.P384())
}

// GenerateRSAKey generates an RSA private key with the specified bit size (2048 or 3072 recommended).
func GenerateRSAKey(bits int) (*rsa.PrivateKey, error) {
	if bits < 2048 {
		return nil, fmt.Errorf("RSA key size too small; minimum 2048 bits for FIPS")
	}
	priv, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, fmt.Errorf("failed to generate RSA key: %w", err)
	}
	return priv, nil
}

// GenerateRSA2048Key is a wrapper for 2048-bit RSA.
func GenerateRSA2048Key() (*rsa.PrivateKey, error) {
	return GenerateRSAKey(2048)
}

// GenerateRSA3072Key is a wrapper for 3072-bit RSA.
func GenerateRSA3072Key() (*rsa.PrivateKey, error) {
	return GenerateRSAKey(3072)
}

// EncodePrivateKeyToPEM encodes a private key (ECDSA or RSA) to PEM format.
func EncodePrivateKeyToPEM(priv interface{}) ([]byte, error) {
	var keyBytes []byte
	var err error

	switch k := priv.(type) {
	case *ecdsa.PrivateKey:
		keyBytes, err = x509.MarshalECPrivateKey(k)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal ECDSA key: %w", err)
		}
	case *rsa.PrivateKey:
		keyBytes = x509.MarshalPKCS1PrivateKey(k)
	default:
		return nil, fmt.Errorf("unsupported private key type")
	}

	return pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: keyBytes}), nil
}

// ParsePEMPrivateKey parses a PEM-encoded private key (ECDSA or RSA).
func ParsePEMPrivateKey(pemBytes []byte) (interface{}, error) {
	block, _ := pem.Decode(pemBytes)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block")
	}

	if block.Type != "PRIVATE KEY" {
		return nil, fmt.Errorf("unsupported PEM block type: %s", block.Type)
	}

	// Try ECDSA first
	ecdsaPriv, err := x509.ParseECPrivateKey(block.Bytes)
	if err == nil {
		return ecdsaPriv, nil
	}

	// Then RSA
	rsaPriv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err == nil {
		return rsaPriv, nil
	}

	return nil, fmt.Errorf("failed to parse private key")
}

// SignWithECDSA signs data using ECDSA with SHA-256.
func SignWithECDSA(priv *ecdsa.PrivateKey, data []byte) ([]byte, error) {
	hash := sha256.Sum256(data)
	r, s, err := ecdsa.Sign(rand.Reader, priv, hash[:])
	if err != nil {
		return nil, fmt.Errorf("failed to sign with ECDSA: %w", err)
	}
	return append(r.Bytes(), s.Bytes()...), nil // Concatenated r|s (DER encoding optional)
}

// VerifyECDSASignature verifies an ECDSA signature with SHA-256.
func VerifyECDSASignature(pub *ecdsa.PublicKey, data, sig []byte) (bool, error) {
	hash := sha256.Sum256(data)
	sigLen := len(sig) / 2
	r := new(big.Int).SetBytes(sig[:sigLen])
	s := new(big.Int).SetBytes(sig[sigLen:])
	return ecdsa.Verify(pub, hash[:], r, s), nil
}

// SignWithRSAPSS signs data using RSA-PSS with the specified hash (defaults to SHA-256).
func SignWithRSAPSS(priv *rsa.PrivateKey, data []byte, hashID crypto.Hash) ([]byte, error) {
	if hashID == 0 {
		hashID = crypto.SHA256
	}

	hasher := newHasher(hashID)
	if hasher == nil {
		return nil, fmt.Errorf("unsupported hash: %v", hashID)
	}

	if _, err := hasher.Write(data); err != nil {
		return nil, fmt.Errorf("failed to hash data: %w", err)
	}
	digest := hasher.Sum(nil)

	return rsa.SignPSS(rand.Reader, priv, hashID, digest, nil)
}

// VerifyRSAPSSSignature verifies an RSA-PSS signature with the specified hash (defaults to SHA-256).
func VerifyRSAPSSSignature(pub *rsa.PublicKey, data, sig []byte, hashID crypto.Hash) error {
	if hashID == 0 {
		hashID = crypto.SHA256
	}

	hasher := newHasher(hashID)
	if hasher == nil {
		return fmt.Errorf("unsupported hash: %v", hashID)
	}

	if _, err := hasher.Write(data); err != nil {
		return fmt.Errorf("failed to hash data: %w", err)
	}
	digest := hasher.Sum(nil)

	return rsa.VerifyPSS(pub, hashID, digest, sig, nil)
}

// newHasher creates a hash.Hash instance for the given crypto.Hash ID (FIPS-approved only).
func newHasher(hashID crypto.Hash) hash.Hash {
	switch hashID {
	case crypto.SHA256:
		return sha256.New()
	case crypto.SHA384:
		return sha512.New384()
	case crypto.SHA512:
		return sha512.New()
	default:
		return nil // Unsupported
	}
}

// EncryptWithRSAOAEP encrypts data using RSA-OAEP with SHA-256.
func EncryptWithRSAOAEP(pub *rsa.PublicKey, data []byte) ([]byte, error) {
	return rsa.EncryptOAEP(sha256.New(), rand.Reader, pub, data, nil)
}

// DecryptWithRSAOAEP decrypts data using RSA-OAEP with SHA-256.
func DecryptWithRSAOAEP(priv *rsa.PrivateKey, ciphertext []byte) ([]byte, error) {
	return rsa.DecryptOAEP(sha256.New(), rand.Reader, priv, ciphertext, nil)
}

// AsymmetricSign signs data using the specified SigningType.
// priv must match the type (ECDSA or RSA PrivateKey).
func AsymmetricSign(priv interface{}, data []byte, sigType SigningType) ([]byte, error) {
	switch sigType {
	case SigningTypeECDSA256, SigningTypeECDSA384:
		ecdsaPriv, ok := priv.(*ecdsa.PrivateKey)
		if !ok {
			return nil, fmt.Errorf("invalid private key for ECDSA")
		}
		return SignWithECDSA(ecdsaPriv, data)
	case SigningTypeRSAPSS2048, SigningTypeRSAPSS3072:
		rsaPriv, ok := priv.(*rsa.PrivateKey)
		if !ok {
			return nil, fmt.Errorf("invalid private key for RSA")
		}
		return SignWithRSAPSS(rsaPriv, data, crypto.SHA256)
	default:
		return nil, fmt.Errorf("unsupported signing type")
	}
}

// AsymmetricVerify verifies a signature using the specified SigningType.
// pub must match the type (ECDSA or RSA PublicKey).
func AsymmetricVerify(pub interface{}, data, sig []byte, sigType SigningType) (bool, error) {
	switch sigType {
	case SigningTypeECDSA256, SigningTypeECDSA384:
		ecdsaPub, ok := pub.(*ecdsa.PublicKey)
		if !ok {
			return false, fmt.Errorf("invalid public key for ECDSA")
		}
		return VerifyECDSASignature(ecdsaPub, data, sig)
	case SigningTypeRSAPSS2048, SigningTypeRSAPSS3072:
		rsaPub, ok := pub.(*rsa.PublicKey)
		if !ok {
			return false, fmt.Errorf("invalid public key for RSA")
		}
		err := VerifyRSAPSSSignature(rsaPub, data, sig, crypto.SHA256)
		return err == nil, err
	default:
		return false, fmt.Errorf("unsupported signing type")
	}
}
