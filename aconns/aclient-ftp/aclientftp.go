package aclient_ftp

import (
	"context"
	"crypto/tls"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/jlaffaye/ftp"
	"github.com/jpfluger/alibs-slim/aconns"
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
		cn.UpdateHealth(aconns.HEALTHSTATUS_VALIDATE_FAILED)
		return false, aconns.TESTSTATUS_FAILED, err
	}

	if cn.connFTP == nil {
		if err := cn.openConnection(); err != nil {
			cn.UpdateHealth(aconns.HEALTHSTATUS_OPEN_FAILED)
			return false, aconns.TESTSTATUS_FAILED, err
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cn.ConnectionTimeout)*time.Second)
	defer cancel()
	if err := cn.testConnectionWithCtx(ctx); err != nil {
		status := aconns.HEALTHSTATUS_PING_FAILED
		if context.DeadlineExceeded == err {
			status = aconns.HEALTHSTATUS_TIMEOUT
		} else if strings.Contains(err.Error(), "network") || strings.Contains(err.Error(), "connection refused") {
			status = aconns.HEALTHSTATUS_NETWORK_ERROR
		}
		cn.UpdateHealth(status)
		//alog.LOGGER(alog.LOGGER_APP).Warn().Err(err).Msg("FTP test failed")
		return false, aconns.TESTSTATUS_FAILED, err
	}

	cn.UpdateHealth(aconns.HEALTHSTATUS_HEALTHY)
	return true, aconns.TESTSTATUS_INITIALIZED_SUCCESSFUL, nil
}

// Refresh refreshes the FTP connection by closing the existing one (if any) and opening a new one.
func (cn *AClientFTP) Refresh() error {
	cn.mu.Lock()
	defer cn.mu.Unlock()
	if cn.connFTP != nil {
		cn.connFTP.Quit()
		cn.connFTP = nil
	}
	return cn.openConnection()
}

// OpenConnection opens a connection to the FTP server.
func (cn *AClientFTP) OpenConnection() error {
	cn.mu.Lock()
	defer cn.mu.Unlock()
	return cn.openConnection()
}

// openConnection opens a connection to the FTP server.
func (cn *AClientFTP) openConnection() error {
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
		cn.UpdateHealth(aconns.HEALTHSTATUS_CLOSED)
	}
	return nil
}

// testConnectionWithCtx tests the FTP connection using a provided context.
func (cn *AClientFTP) testConnectionWithCtx(ctx context.Context) error {
	if cn.connFTP == nil {
		return fmt.Errorf("no FTP connection has been created where host=%s", cn.address)
	}

	done := make(chan error, 1)
	go func() {
		err := cn.connFTP.NoOp()
		done <- err
	}()

	select {
	case err := <-done:
		if err != nil {
			return fmt.Errorf("test connection failed for FTP where host=%s; %v", cn.address, err)
		}
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// FTPClient returns the FTP connection.
func (cn *AClientFTP) FTPClient() *ftp.ServerConn {
	cn.mu.RLock()
	if cn.IsHealthy() && !cn.GetHealth().IsStale(5*time.Minute) {
		defer cn.mu.RUnlock()
		return cn.connFTP
	}
	cn.mu.RUnlock()

	// Upgrade to write lock for refresh
	cn.mu.Lock()
	defer cn.mu.Unlock()
	if _, _, err := cn.Test(); err != nil {
		return nil
	}
	return cn.connFTP
}

//package aclient_ftp
//
//import (
//	"crypto/tls"
//	"fmt"
//	"github.com/jlaffaye/ftp"
//	"github.com/jpfluger/alibs-slim/aconns"
//	"strings"
//	"sync"
//	"time"
//)
//
//// ADAPTERTYPE_FTP defines the adapter type for FTP.
//const (
//	ADAPTERTYPE_FTP   = aconns.AdapterType("ftp")
//	FTP_DEFAULT_PORT  = 21
//	FTPS_DEFAULT_PORT = 990
//)
//
//// AClientFTP represents an FTP client with connection details.
//type AClientFTP struct {
//	aconns.ADBAdapterBase
//
//	ConnectionTimeout  int  `json:"connectionTimeout,omitempty"`
//	InsecureSkipVerify bool `json:"insecureSkipVerify,omitempty"`
//	UseFTPS            bool `json:"useFTPS,omitempty"`
//
//	// CDWorkingDir changes to a working directory upon
//	// opening of the connection. If CDWorkingDir is empty
//	// then Database will be used, if populated.
//	CDWorkingDir string `json:"cdWorkingDir,omitempty"`
//
//	address string
//	connFTP *ftp.ServerConn
//
//	mu sync.RWMutex
//}
//
//// validate checks if the AClientFTP object is valid.
//func (cn *AClientFTP) validate() error {
//	if err := cn.ADBAdapterBase.Validate(); err != nil {
//		if err != aconns.ErrDatabaseIsEmpty {
//			return err
//		}
//	}
//
//	cn.CDWorkingDir = strings.TrimSpace(cn.CDWorkingDir)
//	if cn.CDWorkingDir == "" && cn.Database != "" {
//		cn.CDWorkingDir = cn.Database
//	}
//
//	if cn.Port <= 0 {
//		if cn.UseFTPS {
//			cn.Port = FTPS_DEFAULT_PORT
//		} else {
//			cn.Port = FTP_DEFAULT_PORT
//		}
//	}
//
//	if cn.ConnectionTimeout <= 0 {
//		cn.ConnectionTimeout = 30
//	}
//
//	cn.address = fmt.Sprintf("%s:%d", cn.Host, cn.Port)
//
//	return nil
//}
//
//// Validate checks if the AClientFTP object is valid.
//func (cn *AClientFTP) Validate() error {
//	cn.mu.Lock()
//	defer cn.mu.Unlock()
//	return cn.validate()
//}
//
//// Test attempts to validate the AClientFTP, open a connection if necessary, and test the connection.
//func (cn *AClientFTP) Test() (bool, aconns.TestStatus, error) {
//	cn.mu.Lock()
//	defer cn.mu.Unlock()
//
//	if err := cn.validate(); err != nil {
//		return false, aconns.TESTSTATUS_FAILED, err
//	}
//
//	if cn.connFTP != nil {
//		if err := cn.testConnection(); err == nil {
//			return true, aconns.TESTSTATUS_INITIALIZED_SUCCESSFUL, nil
//		}
//	}
//
//	if err := cn.openConnection(); err != nil {
//		return false, aconns.TESTSTATUS_FAILED, err
//	}
//
//	return true, aconns.TESTSTATUS_INITIALIZED_SUCCESSFUL, nil
//}
//
//// OpenConnection opens a connection to the FTP server.
//func (cn *AClientFTP) OpenConnection() error {
//	cn.mu.Lock()
//	defer cn.mu.Unlock()
//	return cn.openConnection()
//}
//
//// openConnection opens a connection to the FTP server.
//func (cn *AClientFTP) openConnection() error {
//	if err := cn.validate(); err != nil {
//		return err
//	}
//
//	var connFTP *ftp.ServerConn
//	var err error
//
//	if cn.UseFTPS {
//		connFTP, err = ftp.Dial(cn.address,
//			ftp.DialWithTimeout(time.Duration(cn.ConnectionTimeout)*time.Second),
//			ftp.DialWithTLS(&tls.Config{InsecureSkipVerify: cn.InsecureSkipVerify}),
//		)
//	} else {
//		connFTP, err = ftp.Dial(cn.address,
//			ftp.DialWithTimeout(time.Duration(cn.ConnectionTimeout)*time.Second),
//		)
//	}
//
//	if err != nil {
//		return err
//	}
//
//	err = connFTP.Login(cn.Username, cn.Password)
//	if err != nil {
//		connFTP.Quit()
//		return err
//	}
//
//	if cn.CDWorkingDir != "" {
//		if err := connFTP.ChangeDir(cn.CDWorkingDir); err != nil {
//			connFTP.Quit()
//			return fmt.Errorf("failed to change directory to %s: %v", cn.CDWorkingDir, err)
//		}
//	}
//
//	cn.connFTP = connFTP
//	return nil
//}
//
//// CloseConnection closes the FTP connection.
//func (cn *AClientFTP) CloseConnection() error {
//	cn.mu.Lock()
//	defer cn.mu.Unlock()
//
//	if cn.connFTP != nil {
//		if err := cn.connFTP.Quit(); err != nil {
//			return fmt.Errorf("forced quit of FTP connection had errors where host=%s; %v", cn.address, err)
//		}
//		cn.connFTP = nil
//	}
//	return nil
//}
//
//// testConnection tests the FTP connection.
//func (cn *AClientFTP) testConnection() error {
//	if cn.connFTP == nil {
//		return fmt.Errorf("no FTP connection has been created where host=%s", cn.Host)
//	}
//	return nil
//}
//
//// FTPClient returns the FTP connection.
//func (cn *AClientFTP) FTPClient() *ftp.ServerConn {
//	cn.mu.RLock()
//	defer cn.mu.RUnlock()
//	return cn.connFTP
//}
