package aconns

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/jpfluger/alibs-slim/autils"
)

var ErrDatabaseIsEmpty = errors.New("database is empty")
var ErrUsernameIsEmpty = errors.New("username is empty")
var ErrPasswordIsEmpty = errors.New("password is empty")

// IAdapterDB interface extends IAdapter with additional database-related methods.
type IAdapterDB interface {
	IAdapter
	GetDatabase() string
	GetUsername() string
	GetPassword() string
	Validate() error
	Test() (bool, TestStatus, error)
}

// ADBAdapterBase holds the basic database connection information.
type ADBAdapterBase struct {
	Adapter

	Database     string `json:"database,omitempty"`
	Username     string `json:"username,omitempty"`
	Password     string `json:"password,omitempty"`
	PasswordFile string `json:"passwordFile,omitempty"` // Loaded once then deleted when the password is populated.

	mu sync.RWMutex // Protects access to the fields.
}

// Validate checks if the ADBAdapterBase object is valid.
func (cn *ADBAdapterBase) Validate() error {
	cn.mu.Lock()
	defer cn.mu.Unlock()
	return cn.validate()
}

// validate checks if the ADBAdapterBase object is valid.
func (cn *ADBAdapterBase) validate() error {
	if err := cn.Adapter.Validate(); err != nil {
		return err
	}

	// Trim spaces from the string fields to avoid common errors.
	cn.Username = strings.TrimSpace(cn.Username)
	cn.Password = strings.TrimSpace(cn.Password)
	cn.PasswordFile = strings.TrimSpace(cn.PasswordFile)

	// Load the password from the PasswordFile if necessary.
	if cn.Password == "" && cn.PasswordFile != "" {
		var err error
		cn.Password, err = autils.ReadFileTrimSpaceWithError(cn.PasswordFile)
		if err != nil {
			return fmt.Errorf("failed to read password file: %w", err)
		}
		cn.PasswordFile = ""
	}

	cn.Database = strings.TrimSpace(cn.Database)
	if cn.Database == "" {
		return ErrDatabaseIsEmpty
	}
	if cn.Username == "" {
		return ErrUsernameIsEmpty
	}
	if cn.Password == "" {
		return ErrPasswordIsEmpty
	}

	return nil
}

// Test attempts to validate the ADBAdapterBase and returns the test status and error if any.
func (cn *ADBAdapterBase) Test() (bool, TestStatus, error) {
	cn.mu.Lock()
	defer cn.mu.Unlock()

	if err := cn.validate(); err != nil {
		return false, TESTSTATUS_FAILED, err
	}

	// Simulate a test connection to the database.
	// This is where you would add your actual connection logic.
	// For now, we'll just return a successful status.
	return true, TESTSTATUS_INITIALIZED_SUCCESSFUL, nil
}

// GetDatabase returns the database name of the database connection.
func (cn *ADBAdapterBase) GetDatabase() string {
	cn.mu.RLock()
	defer cn.mu.RUnlock()
	return cn.Database
}

// GetUsername returns the username of the database connection.
func (cn *ADBAdapterBase) GetUsername() string {
	cn.mu.RLock()
	defer cn.mu.RUnlock()
	return cn.Username
}

// GetPassword returns the password of the database connection.
func (cn *ADBAdapterBase) GetPassword() string {
	cn.mu.RLock()
	defer cn.mu.RUnlock()
	return cn.Password
}
