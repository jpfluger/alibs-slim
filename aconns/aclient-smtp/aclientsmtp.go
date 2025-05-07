package aclient_smtp

import (
	"crypto/tls"
	"fmt"
	"github.com/jhillyerd/enmime/v2"
	"github.com/jpfluger/alibs-slim/autils"
	"net/smtp"
	"strings"
	"sync"
	"time"

	"github.com/jpfluger/alibs-slim/aconns"
)

const (
	ADAPTERTYPE_SMTP        = aconns.AdapterType("smtp")
	SMTP_DEFAULT_PORT       = 587
	SMTP_CONNECTION_TIMEOUT = 30

	AUTHTYPE_SENDER_NONE     = "none"
	AUTHTYPE_SENDER_PLAIN    = "plain"
	AUTHTYPE_SENDER_IDENTITY = "identity"
)

// AClientSMTP represents an SMTP client with connection details.
type AClientSMTP struct {
	aconns.Adapter

	Username     string `json:"username,omitempty"`
	Password     string `json:"password,omitempty"`
	PasswordFile string `json:"passwordFile,omitempty"` // Loaded once then deleted when the password is populated.

	Identity string `json:"identity,omitempty"`

	ConnectionTimeout  int      `json:"connectionTimeout,omitempty"`
	InsecureSkipVerify bool     `json:"insecureSkipVerify,omitempty"`
	DialMode           DialMode `json:"dialMode,omitempty"`

	AuthType AuthType `json:"authType,omitempty"`

	AllowEmptySubject bool `json:"allowEmptySubject,omitempty"`
	AllowNoTextHMTL   bool `json:"allowNoTextHMTL,omitempty"`
	AllowAttachments  bool `json:"allowAttachments,omitempty"`

	address        string
	auth           smtp.Auth
	hasInitialized bool

	mu sync.RWMutex
}

// validate checks if the AClientSMTP object is valid.
func (cn *AClientSMTP) validate() error {
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

	// Check if AuthType is valid.
	if cn.AuthType != AUTHTYPE_SENDER_NONE && cn.AuthType != AUTHTYPE_SENDER_PLAIN && cn.AuthType != AUTHTYPE_SENDER_IDENTITY {
		return fmt.Errorf("invalid auth type: %s", cn.AuthType)
	}

	if cn.AuthType != AUTHTYPE_SENDER_NONE {
		if cn.Username == "" {
			return aconns.ErrUsernameIsEmpty
		}
		if cn.Password == "" {
			return aconns.ErrPasswordIsEmpty
		}
	}

	cn.Identity = strings.TrimSpace(cn.Identity)
	if cn.AuthType == AUTHTYPE_SENDER_IDENTITY && cn.Identity == "" {
		return fmt.Errorf("identity is empty")
	}

	if cn.Port <= 0 {
		cn.Port = SMTP_DEFAULT_PORT
	}

	if cn.ConnectionTimeout <= 0 {
		cn.ConnectionTimeout = SMTP_CONNECTION_TIMEOUT
	}

	cn.address = fmt.Sprintf("%s:%d", cn.Host, cn.Port)

	return nil
}

// Validate checks if the AClientSMTP object is valid.
func (cn *AClientSMTP) Validate() error {
	cn.mu.Lock()
	defer cn.mu.Unlock()
	return cn.validate()
}

// GetUsername returns the username of the database connection.
func (cn *AClientSMTP) GetUsername() string {
	cn.mu.RLock()
	defer cn.mu.RUnlock()
	return cn.Username
}

// GetPassword returns the password of the database connection.
func (cn *AClientSMTP) GetPassword() string {
	cn.mu.RLock()
	defer cn.mu.RUnlock()
	return cn.Password
}

// GetAddress returns the full address plus port of the SMTP server.
func (cn *AClientSMTP) GetAddress() string {
	cn.mu.RLock()
	defer cn.mu.RUnlock()
	return cn.getAddress()
}

// getAddress returns the full address plus port of the SMTP server.
func (cn *AClientSMTP) getAddress() string {
	return fmt.Sprintf("%s:%d", cn.Host, cn.Port)
}

// Test attempts to validate the AClientSMTP, open a connection if necessary, and test the connection.
func (cn *AClientSMTP) Test() (bool, aconns.TestStatus, error) {
	if err := cn.OpenConnection(); err != nil {
		return false, aconns.TESTSTATUS_FAILED, err
	}
	return true, aconns.TESTSTATUS_INITIALIZED_SUCCESSFUL, nil
}

// OpenConnection verifies and initializes the SMTP connection.
func (cn *AClientSMTP) OpenConnection() error {
	cn.mu.Lock()
	defer cn.mu.Unlock()

	// Reset initialization state
	cn.hasInitialized = false

	// Create and test the sender
	sender, err := cn.getSenderNoLock()
	if err != nil {
		return fmt.Errorf("failed to initialize SMTP sender: %w", err)
	}

	// Test connection
	if err := sender.TestConnection(); err != nil {
		return fmt.Errorf("connection test failed: %w", err)
	}
	cn.hasInitialized = true
	return nil
}

// CloseConnection is a no-op method to satisfy the interface requirements.
func (cn *AClientSMTP) CloseConnection() error {
	cn.mu.Lock()
	defer cn.mu.Unlock()

	// Mark as uninitialized when closing
	cn.hasInitialized = false
	return nil
}

// HasInitialized determines if the SMTP struct and auth have been validated.
func (cn *AClientSMTP) HasInitialized() bool {
	cn.mu.RLock()
	defer cn.mu.RUnlock()
	return cn.hasInitialized
}

// getSenderNoLock creates a CustomSMTPSender with the current configuration.
func (cn *AClientSMTP) getSenderNoLock() (*CustomSMTPSender, error) {

	// Set up the SMTP authentication based on cn.AuthType
	var auth smtp.Auth
	switch cn.AuthType {
	case AUTHTYPE_SENDER_NONE:
		auth = nil
	case AUTHTYPE_SENDER_PLAIN:
		auth = smtp.PlainAuth("", cn.Username, cn.Password, cn.Host)
	case AUTHTYPE_SENDER_IDENTITY:
		auth = smtp.PlainAuth(cn.Identity, cn.Username, cn.Password, cn.Host)
	}

	// Set up the TLS configuration
	tlsConfig := &tls.Config{
		InsecureSkipVerify: cn.InsecureSkipVerify,
		ServerName:         cn.Host,
	}

	// Convert ConnectionTimeout to time.Duration
	timeout := time.Duration(cn.ConnectionTimeout) * time.Second

	// Return a new CustomSMTPSender
	return NewCustomSMTPSender(cn.getAddress(), auth, tlsConfig, timeout, cn.DialMode), nil
}

// GetSender safely acquires the necessary lock and calls getSenderNoLock.
func (cn *AClientSMTP) GetSender() (enmime.Sender, error) {
	cn.mu.RLock()
	defer cn.mu.RUnlock()
	return cn.getSenderNoLock()
}

func (cn *AClientSMTP) SendMail(mailPiece *MailPiece) error {
	if mailPiece == nil {
		return fmt.Errorf("mail piece is nil")
	}

	cn.mu.RLock()
	defer cn.mu.RUnlock()

	// Apply overrides
	mailPiece.AllowEmptySubject = cn.AllowEmptySubject
	mailPiece.AllowNoTextHMTL = cn.AllowNoTextHMTL
	mailPiece.AllowAttachments = cn.AllowAttachments

	// Create the MIME message
	builder, err := mailPiece.validateWithBuilder()
	if err != nil {
		return fmt.Errorf("failed to build MIME message: %w", err)
	}

	// Get the CustomSMTPSender with TLS configuration without additional locking
	sender, err := cn.getSenderNoLock()
	if err != nil {
		return fmt.Errorf("failed to create SMTP sender: %w", err)
	}

	// Send the email using enmime, which in turn calls the CustomSMTPSender
	return builder.Send(sender)
}

type AClientSMTPs []*AClientSMTP

func (cns AClientSMTPs) FindByName(name aconns.AdapterName) *AClientSMTP {
	if cns == nil || len(cns) == 0 {
		return nil
	}
	for _, cn := range cns {
		if cn.Name == name {
			return cn
		}
	}
	return nil
}
