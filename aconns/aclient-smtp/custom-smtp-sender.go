package aclient_smtp

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/smtp"
	"strconv"
	"strings"
	"time"
)

// CustomSMTPSender supports multiple SMTP connection modes and auto-detects the connection type.
type CustomSMTPSender struct {
	addr              string
	auth              smtp.Auth
	tlsConfig         *tls.Config
	connectionTimeout time.Duration
	dialMode          DialMode
}

// NewCustomSMTPSender creates a new CustomSMTPSender with optional connection mode.
func NewCustomSMTPSender(addr string, auth smtp.Auth, tlsConfig *tls.Config, timeout time.Duration, mode DialMode) *CustomSMTPSender {
	return &CustomSMTPSender{
		addr:              addr,
		auth:              auth,
		tlsConfig:         tlsConfig,
		connectionTimeout: timeout,
		dialMode:          mode,
	}
}

// Send sends a message using either implicit TLS or STARTTLS as configured.
func (s *CustomSMTPSender) Send(reversePath string, recipients []string, msg []byte) error {
	client, err := s.connect()
	if err != nil {
		return fmt.Errorf("failed to connect to SMTP server: %w", err)
	}
	defer client.Quit()

	// Authenticate if required
	if s.auth != nil {
		if err := client.Auth(s.auth); err != nil {
			return fmt.Errorf("authentication failed: %w", err)
		}
	}

	// Send MAIL FROM and RCPT TO commands, then send the message
	if err := client.Mail(reversePath); err != nil {
		return fmt.Errorf("MAIL FROM command failed: %w", err)
	}
	for _, recipient := range recipients {
		if err := client.Rcpt(recipient); err != nil {
			return fmt.Errorf("RCPT TO command failed for %s: %w", recipient, err)
		}
	}
	wc, err := client.Data()
	if err != nil {
		return fmt.Errorf("failed to start DATA command: %w", err)
	}
	defer wc.Close()

	if _, err := wc.Write(msg); err != nil {
		return fmt.Errorf("failed to write message body: %w", err)
	}
	return nil
}

// TestConnection verifies connectivity and authentication without sending an email.
func (s *CustomSMTPSender) TestConnection() error {
	client, err := s.connect()
	if err != nil {
		return fmt.Errorf("failed to connect to SMTP server: %w", err)
	}
	defer client.Quit()

	// Authenticate if required
	if s.auth != nil {
		if err := client.Auth(s.auth); err != nil {
			return fmt.Errorf("authentication failed: %w", err)
		}
	}

	// Issue NOOP to test the connection
	if err := client.Noop(); err != nil {
		return fmt.Errorf("NOOP command failed: %w", err)
	}
	return nil
}

// connect establishes an SMTP connection based on the specified or auto-detected connection mode.
func (s *CustomSMTPSender) connect() (*smtp.Client, error) {
	if s.dialMode.IsEmpty() || s.dialMode == DIALMODE_UNKNOWN {
		// Auto-detect mode if connectionMode is Unknown
		return s.autoDetectAndConnect()
	}

	// Proceed with the explicitly specified mode
	switch s.dialMode {
	case DIALMODE_NOTLS:
		return s.connectWithoutTLS()
	case DIALMODE_TLS:
		return s.connectWithTLS()
	case DIALMODE_STARTTLS:
		return s.connectWithStartTLS()
	default:
		return nil, fmt.Errorf("unsupported connection mode: %s", s.dialMode)
	}
}

// autoDetectAndConnect tries to auto-detect the appropriate connection mode based on the server and port.
func (s *CustomSMTPSender) autoDetectAndConnect() (*smtp.Client, error) {
	port := s.getPort()

	switch port {
	case 465:
		// Attempt implicit TLS
		return s.connectWithTLS()
	case 587, 25:
		// Attempt STARTTLS, then fall back to plain
		client, err := s.connectWithStartTLS()
		if err == nil {
			return client, nil
		}
		return s.connectWithoutTLS()
	default:
		// Unknown port: Try STARTTLS first, then fall back to plain
		client, err := s.connectWithStartTLS()
		if err == nil {
			return client, nil
		}
		return s.connectWithoutTLS()
	}
}

// connectWithTLS attempts a direct TLS connection (implicit TLS).
func (s *CustomSMTPSender) connectWithTLS() (*smtp.Client, error) {
	dialer := &net.Dialer{Timeout: s.connectionTimeout}
	conn, err := tls.DialWithDialer(dialer, "tcp", s.addr, s.tlsConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to establish TLS connection: %w", err)
	}
	return smtp.NewClient(conn, s.extractHost())
}

// connectWithStartTLS establishes a plain connection and upgrades with STARTTLS if supported.
func (s *CustomSMTPSender) connectWithStartTLS() (*smtp.Client, error) {
	dialer := &net.Dialer{Timeout: s.connectionTimeout}
	conn, err := dialer.Dial("tcp", s.addr)
	if err != nil {
		return nil, fmt.Errorf("failed to establish plain connection for STARTTLS: %w", err)
	}

	client, err := smtp.NewClient(conn, s.extractHost())
	if err != nil {
		return nil, fmt.Errorf("failed to create SMTP client for STARTTLS: %w", err)
	}

	// Check for STARTTLS support and attempt upgrade
	if ok, _ := client.Extension("STARTTLS"); ok {
		if err = client.StartTLS(s.tlsConfig); err != nil {
			client.Quit()
			return nil, fmt.Errorf("STARTTLS command failed: %w", err)
		}
		return client, nil
	}

	client.Quit()
	return nil, fmt.Errorf("SMTP server does not support STARTTLS")
}

// connectWithoutTLS establishes an unencrypted connection.
func (s *CustomSMTPSender) connectWithoutTLS() (*smtp.Client, error) {
	dialer := &net.Dialer{Timeout: s.connectionTimeout}
	conn, err := dialer.Dial("tcp", s.addr)
	if err != nil {
		return nil, fmt.Errorf("failed to establish unencrypted connection: %w", err)
	}
	return smtp.NewClient(conn, s.extractHost())
}

// extractHost parses the host part from the address for the SMTP handshake.
func (s *CustomSMTPSender) extractHost() string {
	return strings.Split(s.addr, ":")[0]
}

// getPort extracts the port number from the address, defaulting to 25 if parsing fails.
func (s *CustomSMTPSender) getPort() int {
	_, portStr, err := net.SplitHostPort(s.addr)
	if err != nil {
		return 25 // Default to port 25 if parsing fails
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return 25
	}
	return port
}
