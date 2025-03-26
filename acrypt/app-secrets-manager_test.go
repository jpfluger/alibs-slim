package acrypt

import (
	"testing"
)

func TestSetAppSecretsManager(t *testing.T) {
	mgr := NewSecretsManager("test_password")

	// Test setting the global manager
	err := SetAppSecretsManager(mgr, false)
	if err != nil {
		t.Fatalf("unexpected error setting global manager: %v", err)
	}

	// Ensure the global manager is set
	if GetAppSecretsManager() != mgr {
		t.Fatalf("expected global manager to be set, got: %v", GetAppSecretsManager())
	}

	// Test attempting to overwrite without force
	err = SetAppSecretsManager(NewSecretsManager("new_password"), false)
	if err == nil || err.Error() != "ISecretsManager is already set; use force=true to overwrite" {
		t.Fatalf("expected overwrite error, got: %v", err)
	}

	// Test overwriting with force
	err = SetAppSecretsManager(NewSecretsManager("new_password"), true)
	if err != nil {
		t.Fatalf("unexpected error overwriting global manager: %v", err)
	}
}

//func TestLoadAppSecretsManagerFromFile(t *testing.T) {
//	// Create a temporary file for testing
//	tempFile, err := os.CreateTemp("", "secrets_manager_*.json")
//	if err != nil {
//		t.Fatalf("failed to create temp file: %v", err)
//	}
//	tempFilePath := tempFile.Name()
//	defer os.Remove(tempFilePath)
//
//	// Write mock data to the file
//	mockSecretsManager := NewSecretsManager("test_password")
//	mockSecretsManager.Secrets = SecretsItems{{Key: "key1", Value: SecretsValue{}}}
//	mockData, err := json.Marshal(mockSecretsManager)
//	if err != nil {
//		t.Fatalf("failed to marshal mock secrets manager: %v", err)
//	}
//
//	if _, err := tempFile.Write(mockData); err != nil {
//		t.Fatalf("failed to write to temp file: %v", err)
//	}
//	tempFile.Close()
//
//	mgr := NewSecretsManager("")
//
//	// Test loading secrets manager from file
//	err = LoadAppSecretsManagerFromFile(tempFilePath, "", mgr, false)
//	if err != nil {
//		t.Fatalf("unexpected error loading secrets manager from file: %v", err)
//	}
//
//	if len(mgr.Secrets) != 1 || mgr.Secrets[0].Key != "key1" {
//		t.Errorf("expected secrets manager to load secrets, got: %v", mgr.Secrets)
//	}
//}
//
//func TestSaveAppSecretsManagerFile(t *testing.T) {
//	// Create a temporary file path for testing
//	tempFilePath := "test_secrets_manager.json"
//	defer os.Remove(tempFilePath)
//
//	mgr := NewSecretsManager("test_password")
//	mgr.Secrets = SecretsItems{{Key: "key1", Value: SecretsValue{}}}
//
//	// Test saving secrets manager to file
//	err := SaveAppSecretsManagerFile(tempFilePath, "", mgr)
//	if err != nil {
//		t.Fatalf("unexpected error saving secrets manager to file: %v", err)
//	}
//
//	// Verify file contents
//	data, err := os.ReadFile(tempFilePath)
//	if err != nil {
//		t.Fatalf("failed to read saved file: %v", err)
//	}
//
//	var loadedMgr SecretsManager
//	if err := json.Unmarshal(data, &loadedMgr); err != nil {
//		t.Fatalf("failed to unmarshal saved secrets manager: %v", err)
//	}
//
//	if len(loadedMgr.Secrets) != 1 || loadedMgr.Secrets[0].Key != "key1" {
//		t.Errorf("expected saved secrets to contain key1, got: %v", loadedMgr.Secrets)
//	}
//}
