package aclient_ftp

import (
	"crypto/tls"
	"fmt"
	"github.com/jlaffaye/ftp"
	"github.com/jpfluger/alibs-slim/aconns"
	"strings"
	"sync"
	"time"
)

// ADAPTERTYPE_FTP defines the adapter type for FTP.
const (
	ADAPTERTYPE_FTP   = aconns.AdapterType("ftp")
	FTP_DEFAULT_PORT  = 21
	FTPS_DEFAULT_PORT = 990
)

// AClientFTP represents an FTP client with connection details.
type AClientFTP struct {
	aconns.ADBAdapterBase

	ConnectionTimeout  int  `json:"connectionTimeout,omitempty"`
	InsecureSkipVerify bool `json:"insecureSkipVerify,omitempty"`
	UseFTPS            bool `json:"useFTPS,omitempty"`

	// CDWorkingDir changes to a working directory upon
	// opening of the connection. If CDWorkingDir is empty
	// then Database will be used, if populated.
	CDWorkingDir string `json:"cdWorkingDir,omitempty"`

	address string
	connFTP *ftp.ServerConn

	mu sync.RWMutex
}

// validate checks if the AClientFTP object is valid.
func (cn *AClientFTP) validate() error {
	if err := cn.ADBAdapterBase.Validate(); err != nil {
		if err != aconns.ErrDatabaseIsEmpty {
			return err
		}
	}

	cn.CDWorkingDir = strings.TrimSpace(cn.CDWorkingDir)
	if cn.CDWorkingDir == "" && cn.Database != "" {
		cn.CDWorkingDir = cn.Database
	}

	if cn.Port <= 0 {
		if cn.UseFTPS {
			cn.Port = FTPS_DEFAULT_PORT
		} else {
			cn.Port = FTP_DEFAULT_PORT
		}
	}

	if cn.ConnectionTimeout <= 0 {
		cn.ConnectionTimeout = 30
	}

	cn.address = fmt.Sprintf("%s:%d", cn.Host, cn.Port)

	return nil
}

// Validate checks if the AClientFTP object is valid.
func (cn *AClientFTP) Validate() error {
	cn.mu.Lock()
	defer cn.mu.Unlock()
	return cn.validate()
}

// Test attempts to validate the AClientFTP, open a connection if necessary, and test the connection.
func (cn *AClientFTP) Test() (bool, aconns.TestStatus, error) {
	cn.mu.Lock()
	defer cn.mu.Unlock()

	if err := cn.validate(); err != nil {
		return false, aconns.TESTSTATUS_FAILED, err
	}

	if cn.connFTP != nil {
		if err := cn.testConnection(); err == nil {
			return true, aconns.TESTSTATUS_INITIALIZED_SUCCESSFUL, nil
		}
	}

	if err := cn.openConnection(); err != nil {
		return false, aconns.TESTSTATUS_FAILED, err
	}

	return true, aconns.TESTSTATUS_INITIALIZED_SUCCESSFUL, nil
}

// OpenConnection opens a connection to the FTP server.
func (cn *AClientFTP) OpenConnection() error {
	cn.mu.Lock()
	defer cn.mu.Unlock()
	return cn.openConnection()
}

// openConnection opens a connection to the FTP server.
func (cn *AClientFTP) openConnection() error {
	if err := cn.validate(); err != nil {
		return err
	}

	var connFTP *ftp.ServerConn
	var err error

	if cn.UseFTPS {
		connFTP, err = ftp.Dial(cn.address,
			ftp.DialWithTimeout(time.Duration(cn.ConnectionTimeout)*time.Second),
			ftp.DialWithTLS(&tls.Config{InsecureSkipVerify: cn.InsecureSkipVerify}),
		)
	} else {
		connFTP, err = ftp.Dial(cn.address,
			ftp.DialWithTimeout(time.Duration(cn.ConnectionTimeout)*time.Second),
		)
	}

	if err != nil {
		return err
	}

	err = connFTP.Login(cn.Username, cn.Password)
	if err != nil {
		connFTP.Quit()
		return err
	}

	if cn.CDWorkingDir != "" {
		if err := connFTP.ChangeDir(cn.CDWorkingDir); err != nil {
			connFTP.Quit()
			return fmt.Errorf("failed to change directory to %s: %v", cn.CDWorkingDir, err)
		}
	}

	cn.connFTP = connFTP
	return nil
}

// CloseConnection closes the FTP connection.
func (cn *AClientFTP) CloseConnection() error {
	cn.mu.Lock()
	defer cn.mu.Unlock()

	if cn.connFTP != nil {
		if err := cn.connFTP.Quit(); err != nil {
			return fmt.Errorf("forced quit of FTP connection had errors where host=%s; %v", cn.address, err)
		}
		cn.connFTP = nil
	}
	return nil
}

// testConnection tests the FTP connection.
func (cn *AClientFTP) testConnection() error {
	if cn.connFTP == nil {
		return fmt.Errorf("no FTP connection has been created where host=%s", cn.Host)
	}
	return nil
}

// FTPClient returns the FTP connection.
func (cn *AClientFTP) FTPClient() *ftp.ServerConn {
	cn.mu.RLock()
	defer cn.mu.RUnlock()
	return cn.connFTP
}
