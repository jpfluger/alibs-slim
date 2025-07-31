package acrypt

import (
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"strings"
	"testing"
)

func TestHashBCrypt(t *testing.T) {
	tests := []struct {
		name     string
		password string
		cost     int
		wantErr  bool
	}{
		{"valid password default cost", "password123", 0, false},
		{"valid password min cost", "password123", bcrypt.MinCost, false},
		{"valid password max cost", "password123", 12, false},
		{"empty password", "", 0, true},
		{"invalid cost too low", "password123", bcrypt.MinCost - 1, true},
		{"invalid cost too high", "password123", bcrypt.MaxCost + 1, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hashed, err := HashBCrypt(tt.password, tt.cost)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Empty(t, hashed)
			} else {
				assert.NoError(t, err)
				assert.True(t, IsBCryptHash(hashed))
			}
		})
	}
}

func TestVerifyBCrypt(t *testing.T) {
	hashed, _ := HashBCrypt("password123", 0)

	tests := []struct {
		name    string
		hashed  string
		plain   string
		want    bool
		wantErr bool
	}{
		{"match", hashed, "password123", true, false},
		{"mismatch", hashed, "wrongpassword", false, false},
		{"empty hashed", "", "password123", false, true},
		{"empty plain", hashed, "", false, true},
		{"invalid hash format", "invalid", "password123", false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			match, err := VerifyBCrypt(tt.hashed, tt.plain)
			assert.Equal(t, tt.want, match)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestIsBCryptHash(t *testing.T) {
	tests := []struct {
		name string
		hash string
		want bool
	}{
		{"valid 2a", "$2a$10$somehash", true},
		{"valid 2b", "$2b$10$somehash", true},
		{"valid 2y", "$2y$10$somehash", true},
		{"invalid prefix", "$3a$10$somehash", false},
		{"empty", "", false},
		{"no prefix", "somehash", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsBCryptHash(tt.hash)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestHashSCrypt(t *testing.T) {
	tests := []struct {
		name     string
		password string
		presets  *ScryptPresets
		wantErr  bool
	}{
		{"valid default", "password123", nil, false},
		{"valid custom", "password123", &ScryptPresets{N: 65536, R: 8, P: 1, KeyLen: 32}, false},
		{"empty password", "", nil, true},
		{"invalid N not power of 2", "password123", &ScryptPresets{N: 100, R: 8, P: 1, KeyLen: 32}, true},
		{"invalid negative R", "password123", &ScryptPresets{N: 32768, R: -1, P: 1, KeyLen: 32}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hashed, err := HashSCrypt(tt.password, tt.presets)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Empty(t, hashed)
			} else {
				assert.NoError(t, err)
				assert.True(t, IsSCryptHash(hashed))
				// Verify format has ln=
				parts := strings.Split(hashed, "$")
				assert.Contains(t, parts[2], "ln=")
			}
		})
	}
}

func TestVerifySCrypt(t *testing.T) {
	hashed, _ := HashSCrypt("password123", nil)

	tests := []struct {
		name    string
		hashed  string
		plain   string
		want    bool
		wantErr bool
	}{
		{"match", hashed, "password123", true, false},
		{"mismatch", hashed, "wrongpassword", false, false},
		{"empty hashed", "", "password123", false, true},
		{"empty plain", hashed, "", false, true},
		{"invalid format too few parts", "$scrypt$ln=15,r=8,p=1$salt", "password123", false, true},
		{"invalid identifier", "$wrong$ln=15,r=8,p=1$salt$hash", "password123", false, true},
		{"invalid params", "$scrypt$invalid$salt$hash", "password123", false, true},
		{"invalid salt decode", "$scrypt$ln=15,r=8,p=1$invalid-Salt$hash", "password123", false, true},
		{"invalid hash decode", "$scrypt$ln=15,r=8,p=1$salt$invalid-Hash", "password123", false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			match, err := VerifySCrypt(tt.hashed, tt.plain)
			assert.Equal(t, tt.want, match)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestIsSCryptHash(t *testing.T) {
	tests := []struct {
		name string
		hash string
		want bool
	}{
		{"valid", "$scrypt$ln=17,r=8,p=1$someSalt$someHash", true},
		{"invalid identifier", "$scrypt2$ln=17,r=8,p=1$someSalt$someHash", false},
		{"too few parts", "$scrypt$ln=17,r=8,p=1$someSalt", false},
		{"empty", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsSCryptHash(tt.hash)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestHashArgon2id(t *testing.T) {
	tests := []struct {
		name     string
		password string
		presets  *Argon2Presets
		wantErr  bool
	}{
		{"valid default", "password123", nil, false},
		{"valid custom", "password123", &Argon2Presets{Time: 1, Memory: 65536, Threads: 4, KeyLen: 32}, false},
		{"empty password", "", nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hashed, err := HashArgon2id(tt.password, tt.presets)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Empty(t, hashed)
			} else {
				assert.NoError(t, err)
				assert.True(t, IsArgon2idHash(hashed))
			}
		})
	}
}

func TestVerifyArgon2id(t *testing.T) {
	hashed, _ := HashArgon2id("password123", nil)

	tests := []struct {
		name    string
		hashed  string
		plain   string
		want    bool
		wantErr bool
	}{
		{"match", hashed, "password123", true, false},
		{"mismatch", hashed, "wrongpassword", false, false},
		{"empty hashed", "", "password123", false, true},
		{"empty plain", hashed, "", false, true},
		{"invalid format too few parts", "$argon2id$v=19$m=65536,t=1,p=4$salt", "password123", false, true},
		{"invalid identifier", "$argon2i$v=19$m=65536,t=1,p=4$salt$hash", "password123", false, true},
		{"invalid version", "$argon2id$v=18$m=65536,t=1,p=4$salt$hash", "password123", false, true},
		{"invalid params", "$argon2id$v=19$invalid$salt$hash", "password123", false, true},
		{"unknown param", "$argon2id$v=19$m=65536,t=1,p=4,x=1$salt$hash", "password123", false, true},
		{"invalid salt decode", "$argon2id$v=19$m=65536,t=1,p=4$invalid-Salt$hash", "password123", false, true},
		{"invalid hash decode", "$argon2id$v=19$m=65536,t=1,p=4$salt$invalid-Hash", "password123", false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			match, err := VerifyArgon2id(tt.hashed, tt.plain)
			assert.Equal(t, tt.want, match)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestIsArgon2idHash(t *testing.T) {
	tests := []struct {
		name string
		hash string
		want bool
	}{
		{"valid", "$argon2id$v=19$m=65536,t=1,p=4$someSalt$someHash", true},
		{"invalid identifier", "$argon2i$v=19$m=65536,t=1,p=4$someSalt$someHash", false},
		{"invalid version", "$argon2id$v=18$m=65536,t=1,p=4$someSalt$someHash", false},
		{"too few parts", "$argon2id$v=19$m=65536,t=1,p=4$someSalt", false},
		{"extra parts", "$argon2id$v=19$m=65536,t=1,p=4$someSalt$someHash$extra", false},
		{"empty", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsArgon2idHash(tt.hash)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestArgon2idVerifyStaticHash(t *testing.T) {
	// Static hash generated with known parameters (using RawStdEncoding)
	staticHash := "$argon2id$v=19$m=65536,t=1,p=4$W0slJwWEDTWj14RWKx73QQ$VjG08hh5d4Lj9CQyrd7vaHeOYGXZm1TEGXyYYsIGl9g"
	password := "password123"
	wrongPassword := "wrongpassword"

	// Test match
	match, err := VerifyArgon2id(staticHash, password)
	assert.NoError(t, err)
	assert.True(t, match)

	// Test mismatch
	match, err = VerifyArgon2id(staticHash, wrongPassword)
	assert.NoError(t, err)
	assert.False(t, match)
}
