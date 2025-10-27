package acrypt

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"
)

// ISecretsManager defines an interface for managing application secrets.
type ISecretsManager interface {
	SetSecret(item *SecretsItem) error
	FindSecret(key SecretsKey) *SecretsItem
	RemoveSecret(key SecretsKey)
	GetSecret(key SecretsKey) []byte
	EnsureCryptMode(targetMode CryptMode, password string) error
	GetMasterPassword() string
	SetMasterPassword(oldPassword, newPassword string) error
}

// SecretsManager manages secrets within an application.
type SecretsManager struct {
	Secrets SecretsItems `json:"secrets,omitempty"` // Collection of secrets.
	//mappedSecrets  SecretsItemsMap // Cached secrets map for efficient lookup.
	decodedSecrets DecodedSecretsMap // Cache decoded secrets for efficient lookup
	masterPassword string            // Required master password for encoding/decoding.
	mu             sync.RWMutex      // Mutex for thread-safe access.
}

type DecodedSecretsMap map[SecretsKey][]byte

// NewSecretsManager creates a new SecretsManager with an optional master password.
func NewSecretsManager(masterPassword string) *SecretsManager {
	sm, err := NewSecretsManagerWithOptions(masterPassword, false, nil)
	if err != nil {
		panic(err)
	}
	return sm
}

// NewSecretsManagerWithOptions creates a new SecretsManager with an optional master password and secretItems.
func NewSecretsManagerWithOptions(masterPassword string, newMasterIfEmpty bool, secretItems SecretsItems) (*SecretsManager, error) {
	if strings.TrimSpace(masterPassword) == "" && newMasterIfEmpty {
		if pass, err := RandGenerate20(); err != nil {
			return nil, err
		} else {
			masterPassword = pass
		}
	}
	if secretItems == nil {
		secretItems = SecretsItems{}
	}
	return &SecretsManager{
		Secrets:        secretItems,
		decodedSecrets: make(DecodedSecretsMap),
		masterPassword: masterPassword,
	}, nil
}

// SetSecret adds a new SecretsItem to the manager.
func (sm *SecretsManager) SetSecret(item *SecretsItem) error {
	if item == nil {
		return fmt.Errorf("cannot add a nil SecretsItem")
	}
	if item.GetKey().IsEmpty() {
		return fmt.Errorf("SecretsItem key is empty")
	}

	sm.mu.Lock()
	defer sm.mu.Unlock()

	if b, err := item.GetDecodedValue(sm.masterPassword); err != nil {
		return fmt.Errorf("failed to get decoded value; %v", err)
	} else {
		if err = sm.Secrets.Set(item); err != nil {
			return fmt.Errorf("failed to set secret; %v", err)
		}
		sm.decodedSecrets[item.GetKey()] = b
	}

	return nil
}

// Validate ensures the secrets can be decrypted using the
// masterPassword and are applied to the map for faster access.
func (sm *SecretsManager) Validate() error {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	if sm.decodedSecrets == nil {
		sm.decodedSecrets = make(DecodedSecretsMap)
	}
	if sm.Secrets == nil {
		sm.Secrets = SecretsItems{}
	}
	for _, secret := range sm.Secrets {
		if secret.GetKey().IsEmpty() {
			return fmt.Errorf("secret key is empty")
		}
		if b, err := secret.Decode(sm.masterPassword, false); err != nil {
			return fmt.Errorf("failed to decode secret; %v", err)
		} else {
			sm.decodedSecrets[secret.GetKey()] = b
		}
	}
	return nil
}

// GetSecret retrieves a fully decoded and decrypted SecretsItem by key.
func (sm *SecretsManager) GetSecret(key SecretsKey) []byte {
	if key.IsEmpty() {
		return nil
	}

	sm.mu.RLock()

	if len(sm.decodedSecrets) != len(sm.Secrets) {
		sm.mu.RUnlock()
		if err := sm.Validate(); err != nil {
			return nil
		}
		sm.mu.RLock()
	}
	defer sm.mu.RUnlock()

	if b, exists := sm.decodedSecrets[key]; exists {
		return b
	}

	return nil
}

// HasSecret returns true if it can fully decode and decrypt a SecretsItem by key.
func (sm *SecretsManager) HasSecret(key SecretsKey) bool {
	b := sm.GetSecret(key)
	return b != nil && len(b) > 0
}

// FindSecret retrieves a SecretsItem by key.
func (sm *SecretsManager) FindSecret(key SecretsKey) *SecretsItem {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	return sm.Secrets.Find(key)
}

// RemoveSecret deletes a SecretsItem by key.
func (sm *SecretsManager) RemoveSecret(key SecretsKey) {
	if key.IsEmpty() {
		return
	}

	sm.mu.Lock()
	defer sm.mu.Unlock()

	for i, item := range sm.Secrets {
		if item.GetKey() == key {
			sm.Secrets = append(sm.Secrets[:i], sm.Secrets[i+1:]...)
			break
		}
	}
	delete(sm.decodedSecrets, key)
}

// EnsureCryptMode ensures all secrets are in the specified CryptMode.
func (sm *SecretsManager) EnsureCryptMode(targetMode CryptMode, password string) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	for _, item := range sm.Secrets {
		if err := item.Value.EnsureCryptMode(password, targetMode); err != nil {
			return fmt.Errorf("failed to set CryptMode for key %s: %v", item.GetKey(), err)
		}
	}
	return nil
}

// GetMasterPassword retrieves the master password.
func (sm *SecretsManager) GetMasterPassword() string {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	return sm.masterPassword
}

// SetMasterPassword sets the master password and updates all secrets.
func (sm *SecretsManager) SetMasterPassword(oldPassword, newPassword string) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if strings.TrimSpace(oldPassword) == "" {
		oldPassword = sm.masterPassword
	}

	for _, item := range sm.Secrets {
		decodedValue, err := item.Value.Decode(oldPassword, false)
		if err != nil {
			return fmt.Errorf("failed to decode secret for key %s with old password: %v", item.GetKey(), err)
		}

		// Re-encrypt the secret with the new password.
		item.Value.Value.Validate(fmt.Sprintf("d;plain;aes256;%s", decodedValue))
		if err := item.Value.EnsureCryptMode(newPassword, CRYPTMODE_ENCRYPTED); err != nil {
			return fmt.Errorf("failed to re-encrypt secret for key %s with new password: %v", item.GetKey(), err)
		}
	}

	sm.masterPassword = newPassword
	return nil
}

// SaveSecretsManagerToFile saves the given SecretsManager to a file, optionally encrypting it with a password.
func SaveSecretsManagerToFile(filePath, password string, mgr ISecretsManager) error {
	if mgr == nil {
		return fmt.Errorf("mgr is nil")
	}
	data, err := json.Marshal(mgr)
	if err != nil {
		return fmt.Errorf("failed to marshal SecretsManager: %v", err)
	}

	if password != "" {
		data, err = AESGCM256Encrypt(data, password)
		if err != nil {
			return fmt.Errorf("failed to encrypt SecretsManager: %v", err)
		}
	}

	return os.WriteFile(filePath, data, 0600)
}

// LoadSecretsManagerFromFile loads a SecretsManager from a file, optionally decrypting it with a password.
// If no manager is provided, a default *SecretsManager is created.
func LoadSecretsManagerFromFile(filePath, password string, mgr ISecretsManager) (ISecretsManager, error) {
	if filePath == "" {
		return nil, fmt.Errorf("filePath is empty")
	}
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %v", err)
	}

	if password != "" {
		data, err = AESGCM256Decrypt(data, password)
		if err != nil {
			return nil, fmt.Errorf("failed to decrypt file: %v", err)
		}
	}

	if mgr == nil {
		mgr = NewSecretsManager(password)
	}

	if err := json.Unmarshal(data, mgr); err != nil {
		return nil, fmt.Errorf("failed to unmarshal SecretsManager: %v", err)
	}

	return mgr, nil
}
