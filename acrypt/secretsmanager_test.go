package acrypt

import (
	"os"
	"testing"
)

// Unit tests for SecretsManager
func TestSecretsManager(t *testing.T) {
	// Test NewSecretsManager
	masterPassword := "test_master_password"
	sm := NewSecretsManager(masterPassword)
	if sm.GetMasterPassword() != masterPassword {
		t.Errorf("expected master password %s, got %s", masterPassword, sm.GetMasterPassword())
	}

	// Test SetSecret and FindSecret
	key := SecretsKey("test_key")
	item := NewSecretsItem(key, "test_value", ENCODINGTYPE_PLAIN, ENCRYPTIONTYPE_AES256)
	err := sm.SetSecret(item)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	found := sm.FindSecret(key)
	if found == nil || found.GetKey() != key {
		t.Errorf("expected to find secret with key %s", key)
	}

	// Test RemoveSecret
	sm.RemoveSecret(key)
	found = sm.FindSecret(key)
	if found != nil {
		t.Errorf("did not expect to find secret with key %s after removal", key)
	}

	// Test EnsureCryptMode
	item2 := NewSecretsItem(SecretsKey("key2"), "value2", ENCODINGTYPE_PLAIN, ENCRYPTIONTYPE_AES256)
	sm.SetSecret(item2)
	err = sm.EnsureCryptMode(CRYPTMODE_DECRYPTED, "test_password")
	if err != nil {
		t.Fatalf("unexpected error ensuring CryptMode: %v", err)
	}

	// Test SaveToFile and LoadFromFile
	filePath := "test_secrets_manager.json"
	err = SaveSecretsManagerToFile(filePath, "test_password", sm)
	if err != nil {
		t.Fatalf("unexpected error saving to file: %v", err)
	}

	loadedSm := &SecretsManager{}
	_, err = LoadSecretsManagerFromFile(filePath, "test_password", loadedSm)
	if err != nil {
		t.Fatalf("unexpected error loading from file: %v", err)
	}
	if len(loadedSm.Secrets) != 1 {
		t.Errorf("expected 1 secret in loaded manager, got %d", len(loadedSm.Secrets))
	}

	// Clean up
	os.Remove(filePath)
}
