package acrypt

import (
	"fmt"
	"sync"
)

var (
	globalsASM ISecretsManager
	muASM      sync.RWMutex // Ensures thread-safe access to globalsASM.
)

// GetAppSecretsManager returns the global instance of ISecretsManager.
func GetAppSecretsManager() ISecretsManager {
	muASM.RLock()
	defer muASM.RUnlock()
	return globalsASM
}

// APPSECRETS is a shortcut to GetAppSecretsManager().
// Prior to using, set the manager using SetAppSecretsManager.
func APPSECRETS() ISecretsManager {
	return GetAppSecretsManager()
}

// SetAppSecretsManager sets or updates the global instance of ISecretsManager.
// Use force=true to overwrite an existing instance.
func SetAppSecretsManager(target ISecretsManager, force bool) error {
	if target == nil {
		return fmt.Errorf("cannot set a nil ISecretsManager")
	}

	muASM.Lock()
	defer muASM.Unlock()

	if globalsASM != nil && !force {
		return fmt.Errorf("ISecretsManager is already set; use force=true to overwrite")
	}

	globalsASM = target

	return nil
}

func JWTSECRETKEY() []byte {
	return GetAppSecretsManager().GetSecret(SECRETSKEY_JWT)
}
