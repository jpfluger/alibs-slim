package acontact

import (
	"fmt"
	"github.com/jpfluger/alibs-slim/acrypt"
)

type CredentialOptions struct {
	Key                acrypt.SecretsKey
	IgnoreValuesVerify bool
	SecretsManager     acrypt.ISecretsManager
	hasValidated       bool
}

func (opts *CredentialOptions) Validate() error {
	if opts == nil {
		return fmt.Errorf("options are nil")
	}
	if opts.Key.IsEmpty() {
		opts.Key = acrypt.SECRETSKEY_SVT
	}
	if !opts.IgnoreValuesVerify {
		if !opts.HasSecretsManager() {
			return fmt.Errorf("secrets manager is nil")
		}
	}
	opts.hasValidated = true
	return nil
}

func (opts *CredentialOptions) HasValidated() bool {
	return opts.hasValidated
}

func (opts *CredentialOptions) HasSecretsManager() bool {
	if opts.SecretsManager == nil {
		if acrypt.APPSECRETS() == nil {
			return false
		}
	}
	return true
}

func (opts *CredentialOptions) GetAppSecretsManager() acrypt.ISecretsManager {
	if opts.SecretsManager != nil {
		return opts.SecretsManager
	}
	return acrypt.APPSECRETS()
}
