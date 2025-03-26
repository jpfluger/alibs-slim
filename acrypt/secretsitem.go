package acrypt

import (
	"fmt"
	"sync"
)

// SecretsItem represents a single item in the secrets management system.
type SecretsItem struct {
	Key   SecretsKey   `json:"key"`   // The key associated with the secret item.
	Value SecretsValue `json:"value"` // The current value of the secret.
	mu    sync.RWMutex // Mutex to protect concurrent access to the SecretsItem.
}

// NewSecretsItem creates a new SecretsItem with the provided key and value.
func NewSecretsItem(key SecretsKey, value string, encoding EncodingType, encryption EncryptionType) *SecretsItem {
	item := &SecretsItem{
		Key: key,
	}

	// Validate and initialize the SecretsValue.
	item.Value.Value.Validate(fmt.Sprintf("%s;%s;%s;%s", CRYPTMODE_DECRYPTED, encoding, encryption, value))
	return item
}

// NewSecretsItemAuto creates a new SecretsItem with a random value of the specified length.
func NewSecretsItemAuto(key SecretsKey, length int) (*SecretsItem, error) {
	var randomValue string
	var err error

	switch length {
	case 20:
		randomValue, err = RandGenerate20()
	case 32:
		randomValue, err = RandGenerate32()
	case 64:
		randomValue, err = RandGenerate64()
	default:
		return nil, fmt.Errorf("unsupported random value length: %d", length)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to generate random value: %v", err)
	}

	return NewSecretsItem(key, randomValue, ENCODINGTYPE_PLAIN, ENCRYPTIONTYPE_AES256), nil
}

// GetKey safely retrieves the key of the SecretsItem.
func (si *SecretsItem) GetKey() SecretsKey {
	si.mu.RLock()
	defer si.mu.RUnlock()
	return si.Key
}

// SetKey safely sets the key of the SecretsItem.
func (si *SecretsItem) SetKey(key SecretsKey) {
	si.mu.Lock()
	defer si.mu.Unlock()
	si.Key = key
}

// GetDecodedValue safely retrieves the decoded value of the SecretsItem.
func (si *SecretsItem) GetDecodedValue(password string) ([]byte, error) {
	si.mu.RLock()
	defer si.mu.RUnlock()

	if si.Value.IsDecoded() {
		return si.Value.GetDecoded(), nil
	}

	return si.Value.Decode(password, false)
}

// Decode ensures the value is decoded and optionally cached.
func (si *SecretsItem) Decode(password string, cacheDecoded bool) ([]byte, error) {
	si.mu.RLock()
	defer si.mu.RUnlock()
	return si.Value.Decode(password, cacheDecoded)
}

// SetValueDecrypted safely sets the decrypted value of the SecretsItem.
func (si *SecretsItem) SetValueDecrypted(value string, encoding EncodingType, encryption EncryptionType) {
	si.mu.Lock()
	defer si.mu.Unlock()

	si.Value.Value.Validate(fmt.Sprintf("%s;%s;%s;%s", CRYPTMODE_DECRYPTED, encoding, encryption, value))
}

// IsExpired checks if the secret or its old value has expired.
func (si *SecretsItem) IsExpired() bool {
	si.mu.RLock()
	defer si.mu.RUnlock()
	return si.Value.HasExpired()
}

// SecretsItems represents a collection of SecretsItem pointers.
type SecretsItems []*SecretsItem

// Find searches for a SecretsItem by key within a slice of SecretsItems.
func (sis SecretsItems) Find(key SecretsKey) *SecretsItem {
	if len(sis) == 0 || key.IsEmpty() {
		return nil
	}
	for _, si := range sis {
		if si.GetKey() == key {
			return si
		}
	}
	return nil
}

// Set adds or updates a SecretsItem in the collection based on its key.
func (sis *SecretsItems) Set(item *SecretsItem) error {
	if item == nil {
		return fmt.Errorf("secrets item is nil")
	}
	if item.GetKey().IsEmpty() {
		return fmt.Errorf("secrets key is empty")
	}
	if *sis == nil {
		*sis = SecretsItems{}
	}
	for i, si := range *sis {
		if si.GetKey() == item.GetKey() {
			(*sis)[i] = item
			return nil
		}
	}
	*sis = append(*sis, item)
	return nil
}

// Remove deletes a SecretsItem from the collection based on its key.
func (sis *SecretsItems) Remove(key SecretsKey) {
	if key.IsEmpty() {
		return
	}
	var newItems SecretsItems
	for _, si := range *sis {
		if si.GetKey() != key {
			newItems = append(newItems, si)
		}
	}
	*sis = newItems
}

// SecretsItemsMap is a map of SecretsKey to SecretsItem for easy lookup.
type SecretsItemsMap map[SecretsKey]*SecretsItem
