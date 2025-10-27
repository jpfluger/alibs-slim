package acrypt

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/hex"
	"encoding/pem"
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGenerateECDSAKey(t *testing.T) {
	tests := []struct {
		name  string
		curve elliptic.Curve
	}{
		{"P256", elliptic.P256()},
		{"P384", elliptic.P384()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			priv, err := GenerateECDSAKey(tt.curve)
			if err != nil {
				t.Fatalf("GenerateECDSAKey() error = %v", err)
			}
			if priv == nil {
				t.Fatal("GenerateECDSAKey() got nil private key")
			}
			if priv.Curve != tt.curve {
				t.Errorf("GenerateECDSAKey() curve mismatch: got %v, want %v", priv.Curve, tt.curve)
			}
		})
	}
}

func TestGenerateECDSA256Key(t *testing.T) {
	priv, err := GenerateECDSA256Key()
	if err != nil {
		t.Fatalf("GenerateECDSA256Key() error = %v", err)
	}
	if priv == nil {
		t.Fatal("GenerateECDSA256Key() got nil private key")
	}
	if priv.Curve != elliptic.P256() {
		t.Errorf("GenerateECDSA256Key() curve mismatch: got %v, want P256", priv.Curve)
	}
}

func TestGenerateECDSA384Key(t *testing.T) {
	priv, err := GenerateECDSA384Key()
	if err != nil {
		t.Fatalf("GenerateECDSA384Key() error = %v", err)
	}
	if priv == nil {
		t.Fatal("GenerateECDSA384Key() got nil private key")
	}
	if priv.Curve != elliptic.P384() {
		t.Errorf("GenerateECDSA384Key() curve mismatch: got %v, want P384", priv.Curve)
	}
}

func TestGenerateRSAKey(t *testing.T) {
	tests := []struct {
		name    string
		bits    int
		wantErr bool
	}{
		{"2048 bits", 2048, false},
		{"3072 bits", 3072, false},
		{"Too small", 1024, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			priv, err := GenerateRSAKey(tt.bits)
			if (err != nil) != tt.wantErr {
				t.Fatalf("GenerateRSAKey() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && priv == nil {
				t.Fatal("GenerateRSAKey() got nil private key")
			}
			if !tt.wantErr && priv.N.BitLen() != tt.bits {
				t.Errorf("GenerateRSAKey() bit length mismatch: got %d, want %d", priv.N.BitLen(), tt.bits)
			}
		})
	}
}

func TestGenerateRSA2048Key(t *testing.T) {
	priv, err := GenerateRSA2048Key()
	if err != nil {
		t.Fatalf("GenerateRSA2048Key() error = %v", err)
	}
	if priv == nil {
		t.Fatal("GenerateRSA2048Key() got nil private key")
	}
	if priv.N.BitLen() != 2048 {
		t.Errorf("GenerateRSA2048Key() bit length mismatch: got %d, want 2048", priv.N.BitLen())
	}
}

func TestGenerateRSA3072Key(t *testing.T) {
	priv, err := GenerateRSA3072Key()
	if err != nil {
		t.Fatalf("GenerateRSA3072Key() error = %v", err)
	}
	if priv == nil {
		t.Fatal("GenerateRSA3072Key() got nil private key")
	}
	if priv.N.BitLen() != 3072 {
		t.Errorf("GenerateRSA3072Key() bit length mismatch: got %d, want 3072", priv.N.BitLen())
	}
}

func TestEncodePrivateKeyToPEM(t *testing.T) {
	ecdsaPriv, _ := GenerateECDSA256Key()
	rsaPriv, _ := GenerateRSA2048Key()

	tests := []struct {
		name    string
		priv    interface{}
		wantErr bool
	}{
		{"ECDSA", ecdsaPriv, false},
		{"RSA", rsaPriv, false},
		{"Unsupported", struct{}{}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pemBytes, err := EncodePrivateKeyToPEM(tt.priv)
			if (err != nil) != tt.wantErr {
				t.Fatalf("EncodePrivateKeyToPEM() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && len(pemBytes) == 0 {
				t.Fatal("EncodePrivateKeyToPEM() got empty PEM")
			}
		})
	}
}

func TestParsePEMPrivateKey(t *testing.T) {
	ecdsaPriv, _ := GenerateECDSA256Key()
	ecdsaPem, _ := EncodePrivateKeyToPEM(ecdsaPriv)
	rsaPriv, _ := GenerateRSA2048Key()
	rsaPem, _ := EncodePrivateKeyToPEM(rsaPriv)

	tests := []struct {
		name    string
		pem     []byte
		wantErr bool
	}{
		{"ECDSA", ecdsaPem, false},
		{"RSA", rsaPem, false},
		{"Invalid PEM", []byte("invalid"), true},
		{"Wrong Type", pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: []byte{}}), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			priv, err := ParsePEMPrivateKey(tt.pem)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ParsePEMPrivateKey() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && priv == nil {
				t.Fatal("ParsePEMPrivateKey() got nil private key")
			}
		})
	}
}

func TestSignAndVerifyWithECDSA(t *testing.T) {
	priv, _ := GenerateECDSA256Key()
	pub := &priv.PublicKey
	data := []byte("test data")

	sig, err := SignWithECDSA(priv, data)
	if err != nil {
		t.Fatalf("SignWithECDSA() error = %v", err)
	}

	valid, err := VerifyECDSASignature(pub, data, sig)
	if err != nil {
		t.Fatalf("VerifyECDSASignature() error = %v", err)
	}
	if !valid {
		t.Error("VerifyECDSASignature() signature invalid")
	}

	// Invalid signature
	invalidSig := make([]byte, len(sig))
	valid, _ = VerifyECDSASignature(pub, data, invalidSig)
	if valid {
		t.Error("VerifyECDSASignature() accepted invalid signature")
	}
}

func TestSignAndVerifyWithRSAPSS(t *testing.T) {
	priv, _ := GenerateRSA2048Key()
	pub := &priv.PublicKey
	data := []byte("test data")

	tests := []struct {
		name   string
		hashID crypto.Hash
	}{
		{"SHA256", crypto.SHA256},
		{"SHA384", crypto.SHA384},
		{"SHA512", crypto.SHA512},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sig, err := SignWithRSAPSS(priv, data, tt.hashID)
			if err != nil {
				t.Fatalf("SignWithRSAPSS() error = %v", err)
			}

			err = VerifyRSAPSSSignature(pub, data, sig, tt.hashID)
			if err != nil {
				t.Errorf("VerifyRSAPSSSignature() error = %v", err)
			}

			// Invalid hash
			invalidErr := VerifyRSAPSSSignature(pub, data, sig, crypto.MD5SHA1)
			if invalidErr == nil {
				t.Error("VerifyRSAPSSSignature() accepted invalid hash")
			}
		})
	}

	// Unsupported hash
	_, err := SignWithRSAPSS(priv, data, crypto.SHA1)
	if err == nil {
		t.Error("SignWithRSAPSS() accepted unsupported hash")
	}
}

func TestEncryptAndDecryptWithRSAOAEP(t *testing.T) {
	priv, _ := GenerateRSA2048Key()
	pub := &priv.PublicKey
	data := []byte("test data")

	enc, err := EncryptWithRSAOAEP(pub, data)
	if err != nil {
		t.Fatalf("EncryptWithRSAOAEP() error = %v", err)
	}

	dec, err := DecryptWithRSAOAEP(priv, enc)
	if err != nil {
		t.Fatalf("DecryptWithRSAOAEP() error = %v", err)
	}
	if string(dec) != string(data) {
		t.Errorf("DecryptWithRSAOAEP() mismatch: got %s, want %s", dec, data)
	}

	// Invalid ciphertext
	_, err = DecryptWithRSAOAEP(priv, []byte("invalid"))
	if err == nil {
		t.Error("DecryptWithRSAOAEP() accepted invalid ciphertext")
	}
}

func TestAsymmetricSignAndVerify(t *testing.T) {
	data := []byte("test data")

	tests := []struct {
		name    string
		sigType SigningType
		genKey  func() (interface{}, error)
		getPub  func(interface{}) interface{}
		wantErr bool
	}{
		{"ECDSA256", SigningTypeECDSA256, func() (interface{}, error) { return GenerateECDSA256Key() }, func(p interface{}) interface{} { return &p.(*ecdsa.PrivateKey).PublicKey }, false},
		{"ECDSA384", SigningTypeECDSA384, func() (interface{}, error) { return GenerateECDSA384Key() }, func(p interface{}) interface{} { return &p.(*ecdsa.PrivateKey).PublicKey }, false},
		{"RSAPSS2048", SigningTypeRSAPSS2048, func() (interface{}, error) { return GenerateRSA2048Key() }, func(p interface{}) interface{} { return &p.(*rsa.PrivateKey).PublicKey }, false},
		{"RSAPSS3072", SigningTypeRSAPSS3072, func() (interface{}, error) { return GenerateRSA3072Key() }, func(p interface{}) interface{} { return &p.(*rsa.PrivateKey).PublicKey }, false},
		{"Invalid Type", 999, nil, nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.genKey == nil {
				_, err := AsymmetricSign(nil, data, tt.sigType)
				if err == nil {
					t.Error("AsymmetricSign() accepted invalid type")
				}
				return
			}

			priv, _ := tt.genKey()
			sig, err := AsymmetricSign(priv, data, tt.sigType)
			if err != nil {
				t.Fatalf("AsymmetricSign() error = %v", err)
			}

			pub := tt.getPub(priv)
			valid, err := AsymmetricVerify(pub, data, sig, tt.sigType)
			if err != nil {
				t.Fatalf("AsymmetricVerify() error = %v", err)
			}
			if !valid {
				t.Error("AsymmetricVerify() signature invalid")
			}

			// Wrong key type
			_, err = AsymmetricSign(&rsa.PrivateKey{}, data, SigningTypeECDSA256)
			if err == nil {
				t.Error("AsymmetricSign() accepted wrong key type")
			}
		})
	}
}

func TestNewHasher(t *testing.T) {
	tests := []struct {
		name    string
		hashID  crypto.Hash
		wantNil bool
	}{
		{"SHA256", crypto.SHA256, false},
		{"SHA384", crypto.SHA384, false},
		{"SHA512", crypto.SHA512, false},
		{"Unsupported", crypto.SHA1, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hasher := newHasher(tt.hashID)
			if (hasher == nil) != tt.wantNil {
				t.Errorf("newHasher() = %v, wantNil %v", hasher, tt.wantNil)
			}
		})
	}
}

func TestComputeSigningCertFingerprint(t *testing.T) {
	// Generate a valid self-signed certificate for testing
	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	assert.NoError(t, err, "failed to generate test key")

	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject:      pkix.Name{CommonName: "Test Cert"},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().Add(24 * time.Hour),
		KeyUsage:     x509.KeyUsageDigitalSignature,
	}
	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	assert.NoError(t, err, "failed to create test cert")

	// Encode to PEM
	validPEM := string(pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: derBytes}))

	// Compute expected fingerprint (SHA-256 of DER)
	h := sha256.Sum256(derBytes)
	expectedFingerprint := hex.EncodeToString(h[:])

	t.Run("Valid PEM Cert", func(t *testing.T) {
		fingerprint, err := ComputeSigningCertFingerprint(validPEM)
		assert.NoError(t, err)
		assert.Equal(t, expectedFingerprint, fingerprint)
	})

	t.Run("Invalid PEM Not Cert Block", func(t *testing.T) {
		invalidPEM := `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA...
-----END PUBLIC KEY-----`
		_, err := ComputeSigningCertFingerprint(invalidPEM)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not a CERTIFICATE block")
	})

	t.Run("Empty Input", func(t *testing.T) {
		_, err := ComputeSigningCertFingerprint("")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid PEM")
	})

	t.Run("Malformed Cert DER", func(t *testing.T) {
		malformedPEM := `-----BEGIN CERTIFICATE-----
invalidderbytes
-----END CERTIFICATE-----`
		_, err := ComputeSigningCertFingerprint(malformedPEM)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid PEM") // Updated to match actual error
	})

	// Subtest for invalid DER with valid PEM structure
	t.Run("Invalid Cert DER Valid PEM", func(t *testing.T) {
		// Minimal valid DER structure that's malformed for X.509 (e.g., bad sequence)
		badDER := []byte{0x30, 0x03, 0x00, 0x00, 0x00} // Invalid ASN.1
		malformedPEM := string(pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: badDER}))
		_, err := ComputeSigningCertFingerprint(malformedPEM)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid certificate")
	})
}
