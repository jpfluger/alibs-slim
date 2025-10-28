package acrypt

import (
	"bytes"
	"errors"
	"strings"
	"sync"
	"testing"
	"time"
)

// Mock GenerateSecretKey for controlled testing.
var mockGenerateSecretKey = func() ([]byte, error) {
	return []byte("test-key-32-bytes-long-for-mock"), nil // 32 bytes example
}

func TestCryptValue_GetDecoded(t *testing.T) {
	cv := &CryptValue{
		valueDecoded: []byte("decoded"),
	}
	got := cv.GetDecoded()
	if string(got) != "decoded" {
		t.Errorf("expected %q, got %q", "decoded", got)
	}

	// Concurrency test: Multiple readers.
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			cv.GetDecoded()
		}()
	}
	wg.Wait() // No panic means mutex works.
}

func TestCryptValue_HasValue(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  bool
	}{
		{"Empty", "", false},
		{"Whitespace", "   ", false},
		{"Set", "base64;YWJj", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cv := &CryptValue{Value: tt.value}
			if got := cv.HasValue(); got != tt.want {
				t.Errorf("expected %v, got %v", tt.want, got)
			}
		})
	}
}

func TestCryptValue_IsValid(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  bool
	}{
		{"Invalid Format", "invalid", false},
		{"Wrong Prefix", "hex;YWJj", false},
		{"Invalid Base64", "base64;invalid@", false},
		{"Valid", "base64;YWJj", true}, // "abc" base64
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cv := &CryptValue{Value: tt.value}
			if got := cv.IsValid(); got != tt.want {
				t.Errorf("expected %v, got %v", tt.want, got)
			}
		})
	}
}

func TestCryptValue_Decode(t *testing.T) {
	tests := []struct {
		name      string
		value     string
		want      []byte
		wantErr   bool
		cacheTest bool // Test caching
	}{
		{"Invalid Format", "invalid", nil, true, false},
		{"Wrong Prefix", "hex;YWJj", nil, true, false},
		{"Invalid Base64", "base64;invalid@", nil, true, false},
		{"Valid", "base64;YWJj", []byte("abc"), false, false},
		{"Cached", "base64;YWJj", []byte("abc"), false, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cv := &CryptValue{Value: tt.value}
			if tt.cacheTest {
				_, _ = cv.Decode() // Prime cache
			}
			got, err := cv.Decode()
			if (err != nil) != tt.wantErr {
				t.Errorf("expected error %v, got %v", tt.wantErr, err)
			}
			if !bytes.Equal(got, tt.want) {
				t.Errorf("expected %v, got %v", tt.want, got)
			}
			// Check cache set
			if !tt.wantErr && !bytes.Equal(cv.valueDecoded, tt.want) {
				t.Errorf("cache not set: expected %v, got %v", tt.want, cv.valueDecoded)
			}
		})
	}

	// Concurrency: Writer and readers.
	cv := &CryptValue{Value: "base64;YWJj"}
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		cv.Decode() // Write to cache
	}()
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			cv.Decode() // Read
		}()
	}
	wg.Wait()
}

func TestCryptValue_Rotate(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name         string
		initialValue string
		newValue     string
		duration     time.Duration
		wantOldValue string
		wantExpires  bool // Check if expires set
		wantErr      bool
	}{
		{"Basic Rotate", "base64;b2xk", "base64;bmV3", 0, "base64;b2xk", false, false},
		{"With Duration", "base64;b2xk", "base64;bmV3", time.Minute, "base64;b2xk", true, false},
		{"Invalid New", "base64;b2xk", "invalid", 0, "", false, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cv := &CryptValue{Value: tt.initialValue}
			err := cv.Rotate(tt.newValue, tt.duration)
			if (err != nil) != tt.wantErr {
				t.Errorf("expected error %v, got %v", tt.wantErr, err)
			}
			if !tt.wantErr {
				if cv.Value != tt.newValue {
					t.Errorf("expected Value %q, got %q", tt.newValue, cv.Value)
				}
				if cv.OldValue != tt.initialValue {
					t.Errorf("expected OldValue %q, got %q", tt.initialValue, cv.OldValue)
				}
				if tt.wantExpires {
					if cv.OldValueExpiresAt == nil || cv.OldValueExpiresAt.Before(now.Add(tt.duration)) {
						t.Errorf("expected expires around %v, got %v", now.Add(tt.duration), cv.OldValueExpiresAt)
					}
				} else if cv.OldValueExpiresAt != nil {
					t.Error("expected no expires, but set")
				}
				// Check decode cached
				if len(cv.valueDecoded) == 0 {
					t.Error("decode not cached after rotate")
				}
			}
		})
	}
}

func TestCryptValue_HasExpired(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name         string
		maxDuration  int
		oldExpiresAt *time.Time
		want         bool
	}{
		{"No Duration", 0, nil, false},
		{"MaxDuration Expired", 1, nil, true}, // Assuming time.Since(time.Time{}) > 1 min
		{"Old Expired", 0, &now, true},        // Past time
		{"Old Not Expired", 0, func() *time.Time { future := now.Add(time.Hour); return &future }(), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cv := &CryptValue{MaxDuration: tt.maxDuration, OldValueExpiresAt: tt.oldExpiresAt}
			if got := cv.HasExpired(); got != tt.want {
				t.Errorf("expected %v, got %v", tt.want, got)
			}
		})
	}
}

func TestCryptValueMap_Initialize(t *testing.T) {
	oldGen := GenerateSecretKey
	defer func() { GenerateSecretKey = oldGen }()
	GenerateSecretKey = mockGenerateSecretKey

	cvm := make(CryptValueMap)
	required := []SecretsKey{"test"}
	err := cvm.Initialize(required)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(cvm) != 1 {
		t.Errorf("expected 1 item, got %d", len(cvm))
	}
	cv := cvm["test"]
	if !cv.HasValue() || !cv.IsValid() || len(cv.GetDecoded()) == 0 {
		t.Error("initialized value invalid or not cached")
	}

	// Error case: Mock gen error
	GenerateSecretKey = func() ([]byte, error) { return nil, errors.New("gen error") }
	err = cvm.Initialize([]SecretsKey{"error"})
	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestCryptValueMap_SetCryptValueClearBytes(t *testing.T) {
	tests := []struct {
		name        string
		key         SecretsKey
		clearBytes  []byte
		wantErr     bool
		wantErrMsg  string
		wantValue   string // Expected formatted Value after set
		wantDecoded []byte // Expected decoded value
	}{
		{
			name:        "Successful Set",
			key:         "test",
			clearBytes:  []byte("abc"),
			wantErr:     false,
			wantValue:   "base64;YWJj",
			wantDecoded: []byte("abc"),
		},
		{
			name:       "Empty Key",
			key:        "", // Assuming SecretsKey("") is empty
			clearBytes: []byte("abc"),
			wantErr:    true,
			wantErrMsg: "key is empty",
		},
		{
			name:       "Nil ClearBytes",
			key:        "test",
			clearBytes: nil,
			wantErr:    true,
			wantErrMsg: "clear bytes is empty",
		},
		{
			name:       "Empty ClearBytes",
			key:        "test",
			clearBytes: []byte{},
			wantErr:    true,
			wantErrMsg: "clear bytes is empty",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cvm := make(CryptValueMap)
			err := cvm.SetCryptValueClearBytes(tt.key, tt.clearBytes)
			if (err != nil) != tt.wantErr {
				t.Errorf("expected error %v, got %v", tt.wantErr, err)
			}
			if tt.wantErr && !strings.Contains(err.Error(), tt.wantErrMsg) {
				t.Errorf("expected error message containing %q, got %v", tt.wantErrMsg, err)
			}
			if !tt.wantErr {
				cv, exists := cvm[tt.key]
				if !exists {
					t.Error("key not set in map")
				}
				if cv.Value != tt.wantValue {
					t.Errorf("expected Value %q, got %q", tt.wantValue, cv.Value)
				}
				decoded := cv.GetDecoded()
				if !bytes.Equal(decoded, tt.wantDecoded) {
					t.Errorf("expected decoded %v, got %v", tt.wantDecoded, decoded)
				}
			}
		})
	}
}

func TestCryptValueMap_Validate(t *testing.T) {
	tests := []struct {
		name     string
		cvm      CryptValueMap
		required []SecretsKey
		wantErr  string // Expected error substring
	}{
		{"No Errors", CryptValueMap{"valid": &CryptValue{Value: "base64;YWJj"}}, []SecretsKey{"valid"}, ""},
		{"Missing Key", CryptValueMap{"valid": &CryptValue{Value: "base64;YWJj"}}, []SecretsKey{"missing"}, "missing missing"},
		{"Invalid Value", CryptValueMap{"invalid": &CryptValue{Value: "invalid"}}, []SecretsKey{"invalid"}, "invalid invalid"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cvm.Validate(tt.required)
			if tt.wantErr == "" {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			} else {
				if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
					t.Errorf("expected error containing %q, got %v", tt.wantErr, err)
				}
			}
		})
	}
}

func TestCryptValueMap_Rotate(t *testing.T) {
	oldGen := GenerateSecretKey
	defer func() { GenerateSecretKey = oldGen }()
	GenerateSecretKey = mockGenerateSecretKey

	cvm := CryptValueMap{"test": &CryptValue{Value: "base64;b2xk"}}
	required := []SecretsKey{"test"}
	err := cvm.Rotate(required, time.Minute)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	cv := cvm["test"]
	if cv.OldValue != "base64;b2xk" || cv.OldValueExpiresAt == nil {
		t.Error("old value or expires not set")
	}
	if !strings.Contains(cv.Value, "base64;") {
		t.Errorf("new value invalid: %s", cv.Value)
	}

	// Error: Missing key
	err = cvm.Rotate([]SecretsKey{"missing"}, time.Minute)
	if err == nil {
		t.Error("expected error")
	}
}

func TestCryptValueMap_HasAnyExpired(t *testing.T) {
	past := time.Now().Add(-time.Hour)
	cvm := CryptValueMap{
		"expired": &CryptValue{OldValueExpiresAt: &past},
		"not":     &CryptValue{},
	}
	if !cvm.HasAnyExpired([]SecretsKey{"expired"}) {
		t.Error("expected true for expired")
	}
	if cvm.HasAnyExpired([]SecretsKey{"not"}) {
		t.Error("expected false for not expired")
	}
}

func TestCryptValueMap_GetDecoded(t *testing.T) {
	cvm := CryptValueMap{"test": &CryptValue{Value: "base64;YWJj"}}
	got, err := cvm.GetDecoded("test")
	if err != nil || string(got) != "abc" {
		t.Errorf("expected 'abc', got %q err %v", got, err)
	}
	_, err = cvm.GetDecoded("missing")
	if err == nil {
		t.Error("expected error for missing")
	}
}

func TestCryptValueMap_Set(t *testing.T) {
	cvm := make(CryptValueMap)
	err := cvm.Set("test", "base64;YWJj")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(cvm) != 1 || !cvm["test"].IsValid() {
		t.Error("set failed")
	}
	err = cvm.Set("", "base64;YWJj")
	if err == nil {
		t.Error("expected error for empty key")
	}
	err = cvm.Set("invalid", "invalid")
	if err == nil {
		t.Error("expected error for invalid value")
	}
}

func TestCryptValueMap_Delete(t *testing.T) {
	cvm := CryptValueMap{"test": &CryptValue{}}
	cvm.Delete("test")
	if len(cvm) != 0 {
		t.Error("delete failed")
	}
	cvm.Delete("missing") // No-op
}
