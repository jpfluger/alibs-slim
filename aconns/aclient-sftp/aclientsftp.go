package aclient_sftp

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/jpfluger/alibs-slim/aconns"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

// ADAPTERTYPE_SFTP defines the adapter type for FTP.
const (
	ADAPTERTYPE_FTP   = aconns.AdapterType("sftp")
	SFTP_DEFAULT_PORT = 22
)

// AClientSFTP represents an FTP client with connection details.
type AClientSFTP struct {
	aconns.ADBAdapterBase

	ConnectionTimeout  int  `json:"connectionTimeout,omitempty"`
	InsecureSkipVerify bool `json:"insecureSkipVerify,omitempty"`

	// CDWorkingDir changes to a working directory upon
	// opening of the connection. If CDWorkingDir is empty
	// then Database will be used, if populated.
	CDWorkingDir string `json:"cdWorkingDir,omitempty"`

	address  string
	sshConn  *ssh.Client
	connSFTP *sftp.Client

	mu sync.RWMutex
}

// validate checks if the AClientSFTP object is valid.
func (cn *AClientSFTP) validate() error {
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
		cn.Port = SFTP_DEFAULT_PORT
	}

	if cn.ConnectionTimeout <= 0 {
		cn.ConnectionTimeout = 30
	}

	cn.address = fmt.Sprintf("%s:%d", cn.Host, cn.Port)

	return nil
}

// Validate checks if the AClientSFTP object is valid.
func (cn *AClientSFTP) Validate() error {
	cn.mu.Lock()
	defer cn.mu.Unlock()
	return cn.validate()
}

// Test attempts to validate the AClientSFTP, open a connection if necessary, and test the connection.
func (cn *AClientSFTP) Test() (bool, aconns.TestStatus, error) {
	cn.mu.Lock()
	defer cn.mu.Unlock()
	return cn.test()
}

// Test attempts to validate the AClientSFTP, open a connection if necessary, and test the connection.
func (cn *AClientSFTP) test() (bool, aconns.TestStatus, error) {
	if err := cn.validate(); err != nil {
		cn.UpdateHealth(aconns.HEALTHSTATUS_VALIDATE_FAILED)
		return false, aconns.TESTSTATUS_FAILED, err
	}

	if cn.connSFTP == nil {
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
		//alog.LOGGER(alog.LOGGER_APP).Warn().Err(err).Msg("SFTP test failed")
		return false, aconns.TESTSTATUS_FAILED, err
	}

	cn.UpdateHealth(aconns.HEALTHSTATUS_HEALTHY)
	return true, aconns.TESTSTATUS_INITIALIZED_SUCCESSFUL, nil
}

// Refresh refreshes the SFTP connection by closing the existing one (if any) and opening a new one.
func (cn *AClientSFTP) Refresh() error {
	cn.mu.Lock()
	defer cn.mu.Unlock()
	if cn.connSFTP != nil {
		cn.connSFTP.Close()
		cn.connSFTP = nil
	}
	if cn.sshConn != nil {
		cn.sshConn.Close()
		cn.sshConn = nil
	}
	return cn.openConnection()
}

// openConnection opens a connection to the SFTP server.
func (cn *AClientSFTP) openConnection() error {
	sshConfig := &ssh.ClientConfig{
		User: cn.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(cn.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         time.Duration(cn.ConnectionTimeout) * time.Second,
	}

	sshClient, err := ssh.Dial("tcp", cn.address, sshConfig)
	if err != nil {
		return fmt.Errorf("failed to dial: %v", err)
	}

	sftpClient, err := sftp.NewClient(sshClient)
	if err != nil {
		sshClient.Close()
		return fmt.Errorf("failed to create SFTP client: %v", err)
	}

	cn.sshConn = sshClient
	cn.connSFTP = sftpClient

	if cn.CDWorkingDir != "" {
		// Check if the directory exists using Stat before attempting operations
		if _, err := cn.connSFTP.Stat(cn.CDWorkingDir); err != nil {
			cn.connSFTP.Close()
			cn.sshConn.Close()
			cn.sshConn = nil
			cn.connSFTP = nil
			return fmt.Errorf("working directory %s does not exist: %v", cn.CDWorkingDir, err)
		}
		// Note: sftp.Client does not support Chdir; use absolute paths for operations instead
	}

	return nil
}

// testConnectionWithCtx tests the connection to the SFTP server using a provided context.
func (cn *AClientSFTP) testConnectionWithCtx(ctx context.Context) error {
	if cn.connSFTP == nil {
		return fmt.Errorf("no active SFTP connection to host=%s", cn.address)
	}

	done := make(chan error, 1)
	go func() {
		// Test by stating the root directory (simple noop-like operation)
		_, err := cn.connSFTP.Stat(".")
		done <- err
	}()

	select {
	case err := <-done:
		if err != nil {
			return fmt.Errorf("test connection failed for SFTP where host=%s; %v", cn.address, err)
		}
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// CloseConnection closes the connection to the SFTP server.
func (cn *AClientSFTP) CloseConnection() error {
	cn.mu.Lock()
	defer cn.mu.Unlock()

	var err error

	if cn.connSFTP != nil {
		if closeErr := cn.connSFTP.Close(); closeErr != nil {
			err = fmt.Errorf("error when closing SFTP client where host=%s; %v", cn.address, closeErr)
		}
		cn.connSFTP = nil
	}

	if cn.sshConn != nil {
		if closeErr := cn.sshConn.Close(); closeErr != nil {
			if err != nil {
				err = fmt.Errorf("%v; %v", err, closeErr)
			} else {
				err = fmt.Errorf("error when closing SSH connection where host=%s; %v", cn.address, closeErr)
			}
		}
		cn.sshConn = nil
	}
	cn.UpdateHealth(aconns.HEALTHSTATUS_CLOSED)

	return err
}

// GetFileBytes retrieves a file from the SFTP server and returns its contents as a byte slice.
func (cn *AClientSFTP) GetFileBytes(remoteFilePath string) ([]byte, error) {
	cn.mu.RLock()
	if cn.IsHealthy() && !cn.GetHealth().IsStale(5*time.Minute) {
		defer cn.mu.RUnlock()
	} else {
		cn.mu.RUnlock()
		cn.mu.Lock()
		defer cn.mu.Unlock()
		if _, _, err := cn.test(); err != nil {
			return nil, err
		}
	}

	if cn.connSFTP == nil {
		return nil, fmt.Errorf("no active SFTP connection to host=%s", cn.address)
	}

	// Open remote file
	remoteFile, err := cn.connSFTP.Open(remoteFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open remote file %s: %v", remoteFilePath, err)
	}
	defer remoteFile.Close()

	// Read file content into memory
	fileBytes, err := io.ReadAll(remoteFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read remote file %s: %v", remoteFilePath, err)
	}

	return fileBytes, nil
}

// DownloadFile retrieves a file from the SFTP server and saves it to a local path.
func (cn *AClientSFTP) DownloadFile(remoteFilePath, localFilePath string) error {
	cn.mu.RLock()
	if cn.IsHealthy() && !cn.GetHealth().IsStale(5*time.Minute) {
		defer cn.mu.RUnlock()
	} else {
		cn.mu.RUnlock()
		cn.mu.Lock()
		defer cn.mu.Unlock()
		if _, _, err := cn.test(); err != nil {
			return err
		}
	}

	if cn.connSFTP == nil {
		return fmt.Errorf("no active SFTP connection to host=%s", cn.address)
	}

	// Open remote file
	remoteFile, err := cn.connSFTP.Open(remoteFilePath)
	if err != nil {
		return fmt.Errorf("failed to open remote file %s: %v", remoteFilePath, err)
	}
	defer remoteFile.Close()

	// Create local file
	localFile, err := os.Create(localFilePath)
	if err != nil {
		return fmt.Errorf("failed to create local file %s: %v", localFilePath, err)
	}
	defer localFile.Close()

	// Copy contents from remote file to local file
	_, err = io.Copy(localFile, remoteFile)
	if err != nil {
		return fmt.Errorf("failed to copy remote file %s to local path %s: %v", remoteFilePath, localFilePath, err)
	}

	return nil
}

// SFTPClient returns the active SFTP connection.
func (cn *AClientSFTP) SFTPClient() *sftp.Client {
	cn.mu.RLock()
	defer cn.mu.RUnlock()
	return cn.connSFTP
}

// CountFilesInDir returns the number of files in the specified directory.
func (cn *AClientSFTP) CountFilesInDir(dir string) (int, error) {
	cn.mu.RLock()
	defer cn.mu.RUnlock()

	if cn.connSFTP == nil {
		return 0, fmt.Errorf("no active SFTP connection to host=%s", cn.address)
	}

	// Default to CDWorkingDir if dir is empty
	if dir == "" {
		dir = cn.CDWorkingDir
		if dir == "" {
			return 0, fmt.Errorf("no directory specified and CDWorkingDir is not set")
		}
	}

	// Verify if the directory exists
	fileInfo, err := cn.connSFTP.Stat(dir)
	if err != nil {
		return 0, fmt.Errorf("failed to access directory %s: %v", dir, err)
	}

	// Ensure it's a directory
	if !fileInfo.IsDir() {
		return 0, fmt.Errorf("%s is not a directory", dir)
	}

	// List files in the directory
	files, err := cn.connSFTP.ReadDir(dir)
	if err != nil {
		return 0, fmt.Errorf("failed to read directory %s: %v", dir, err)
	}

	// Return the count of files
	return len(files), nil
}

// ListFilesAfterDate returns a list of files modified after the given date in the specified directory.
func (cn *AClientSFTP) ListFilesAfterDate(dir string, afterDate time.Time) ([]string, error) {
	cn.mu.RLock()
	defer cn.mu.RUnlock()

	if cn.connSFTP == nil {
		return nil, fmt.Errorf("no active SFTP connection to host=%s", cn.address)
	}

	// Default to CDWorkingDir if dir is empty
	if dir == "" {
		dir = cn.CDWorkingDir
		if dir == "" {
			return nil, fmt.Errorf("no directory specified and CDWorkingDir is not set")
		}
	}

	// Verify if the directory exists
	fileInfo, err := cn.connSFTP.Stat(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to access directory %s: %v", dir, err)
	}

	// Ensure it's a directory
	if !fileInfo.IsDir() {
		return nil, fmt.Errorf("%s is not a directory", dir)
	}

	// List files in the directory
	files, err := cn.connSFTP.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory %s: %v", dir, err)
	}

	// Filter files modified after the given date
	var recentFiles []string
	for _, file := range files {
		if file.ModTime().After(afterDate) {
			recentFiles = append(recentFiles, file.Name())
		}
	}

	return recentFiles, nil
}

//package aclient_sftp
//
//import (
//	"fmt"
//	"io"
//	"os"
//	"strings"
//	"sync"
//	"time"
//
//	"github.com/jpfluger/alibs-slim/aconns"
//	"github.com/pkg/sftp"
//	"golang.org/x/crypto/ssh"
//)
//
//// ADAPTERTYPE_SFTP defines the adapter type for FTP.
//const (
//	ADAPTERTYPE_FTP   = aconns.AdapterType("sftp")
//	SFTP_DEFAULT_PORT = 22
//)
//
//// AClientSFTP represents an FTP client with connection details.
//type AClientSFTP struct {
//	aconns.ADBAdapterBase
//
//	ConnectionTimeout  int  `json:"connectionTimeout,omitempty"`
//	InsecureSkipVerify bool `json:"insecureSkipVerify,omitempty"`
//
//	// CDWorkingDir changes to a working directory upon
//	// opening of the connection. If CDWorkingDir is empty
//	// then Database will be used, if populated.
//	CDWorkingDir string `json:"cdWorkingDir,omitempty"`
//
//	address  string
//	sshConn  *ssh.Client
//	connSFTP *sftp.Client
//
//	mu sync.RWMutex
//}
//
//// validate checks if the AClientSFTP object is valid.
//func (cn *AClientSFTP) validate() error {
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
//		cn.Port = SFTP_DEFAULT_PORT
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
//// Validate checks if the AClientSFTP object is valid.
//func (cn *AClientSFTP) Validate() error {
//	cn.mu.Lock()
//	defer cn.mu.Unlock()
//	return cn.validate()
//}
//
//// Test attempts to validate the AClientSFTP, open a connection if necessary, and test the connection.
//func (cn *AClientSFTP) Test() (bool, aconns.TestStatus, error) {
//	cn.mu.Lock()
//	defer cn.mu.Unlock()
//
//	if err := cn.validate(); err != nil {
//		return false, aconns.TESTSTATUS_FAILED, err
//	}
//
//	if cn.connSFTP != nil {
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
//func (cn *AClientSFTP) OpenConnection() error {
//	cn.mu.Lock()
//	defer cn.mu.Unlock()
//	return cn.openConnection()
//}
//
//// openConnection opens a persistent connection to the SFTP server.
//func (cn *AClientSFTP) openConnection() error {
//	if err := cn.validate(); err != nil {
//		return err
//	}
//
//	// SSH client configuration
//	config := &ssh.ClientConfig{
//		User: cn.Username,
//		Auth: []ssh.AuthMethod{
//			ssh.Password(cn.Password),
//		},
//		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // Note: Use proper host key validation in production
//		Timeout:         time.Duration(cn.ConnectionTimeout) * time.Second,
//	}
//
//	// Connect to SSH server
//	sshConn, err := ssh.Dial("tcp", cn.address, config)
//	if err != nil {
//		return fmt.Errorf("failed to dial: %v", err)
//	}
//
//	// Create SFTP client
//	sftpClient, err := sftp.NewClient(sshConn)
//	if err != nil {
//		_ = sshConn.Close() // Close SSH connection if SFTP client creation fails
//		return fmt.Errorf("failed to create SFTP client: %v", err)
//	}
//
//	// Store connections in struct for later use and explicit closure
//	cn.sshConn = sshConn
//	cn.connSFTP = sftpClient
//
//	// Test connection immediately after opening
//	if err := cn.testConnection(); err != nil {
//		_ = cn.closeConnection()
//		return fmt.Errorf("connection test failed: %v", err)
//	}
//
//	return nil
//}
//
//// CloseConnection closes the SFTP and SSH connections.
//func (cn *AClientSFTP) CloseConnection() error {
//	cn.mu.Lock()
//	defer cn.mu.Unlock()
//	return cn.closeConnection()
//}
//
//// closeConnection closes the SFTP and SSH connections.
//func (cn *AClientSFTP) closeConnection() error {
//	var closeErrors []error
//
//	// Close SFTP connection if it exists
//	if cn.connSFTP != nil {
//		if err := cn.connSFTP.Close(); err != nil {
//			closeErrors = append(closeErrors, fmt.Errorf("failed to close SFTP connection to host=%s: %v", cn.address, err))
//		}
//		cn.connSFTP = nil
//	}
//
//	// Close SSH connection if it exists
//	if cn.sshConn != nil {
//		if err := cn.sshConn.Close(); err != nil {
//			closeErrors = append(closeErrors, fmt.Errorf("failed to close SSH connection to host=%s: %v", cn.address, err))
//		}
//		cn.sshConn = nil
//	}
//
//	// Aggregate and return errors if any occurred
//	if len(closeErrors) > 0 {
//		return fmt.Errorf("errors occurred while closing connections: %v", closeErrors)
//	}
//
//	return nil
//}
//
//// TestConnection tests if the SFTP connection is active.
//func (cn *AClientSFTP) TestConnection() error {
//	cn.mu.RLock()
//	defer cn.mu.RUnlock()
//	return cn.testConnection()
//}
//
//// testConnection tests if the SFTP connection is active and validates the working directory.
//func (cn *AClientSFTP) testConnection() error {
//	if cn.connSFTP == nil {
//		return fmt.Errorf("no active SFTP connection to host=%s", cn.address)
//	}
//
//	// Perform a basic operation to validate the connection
//	wd, err := cn.connSFTP.Getwd()
//	if err != nil {
//		return fmt.Errorf("SFTP connection is invalid where host=%s: %v", cn.address, err)
//	}
//
//	fmt.Println(wd)
//
//	//// Check if the working directory exists
//	//if cn.CDWorkingDir != "" {
//	//	if _, err := cn.connSFTP.Stat(cn.CDWorkingDir); err != nil {
//	//		return fmt.Errorf("working directory %s is not accessible: %v", cn.CDWorkingDir, err)
//	//	}
//	//}
//	//fmt.Printf("SFTP connection to %s is valid. Current directory: %s\n", cn.address, wd)
//
//	return nil
//}
//
//// SFTPClient returns the active SFTP connection.
//func (cn *AClientSFTP) SFTPClient() *sftp.Client {
//	cn.mu.RLock()
//	defer cn.mu.RUnlock()
//	return cn.connSFTP
//}
//
//// CountFilesInDir returns the number of files in the specified directory.
//func (cn *AClientSFTP) CountFilesInDir(dir string) (int, error) {
//	cn.mu.RLock()
//	defer cn.mu.RUnlock()
//
//	if cn.connSFTP == nil {
//		return 0, fmt.Errorf("no active SFTP connection to host=%s", cn.address)
//	}
//
//	// Default to CDWorkingDir if dir is empty
//	if dir == "" {
//		dir = cn.CDWorkingDir
//		if dir == "" {
//			return 0, fmt.Errorf("no directory specified and CDWorkingDir is not set")
//		}
//	}
//
//	// Verify if the directory exists
//	fileInfo, err := cn.connSFTP.Stat(dir)
//	if err != nil {
//		return 0, fmt.Errorf("failed to access directory %s: %v", dir, err)
//	}
//
//	// Ensure it's a directory
//	if !fileInfo.IsDir() {
//		return 0, fmt.Errorf("%s is not a directory", dir)
//	}
//
//	// List files in the directory
//	files, err := cn.connSFTP.ReadDir(dir)
//	if err != nil {
//		return 0, fmt.Errorf("failed to read directory %s: %v", dir, err)
//	}
//
//	// Return the count of files
//	return len(files), nil
//}
//
//// ListFilesAfterDate returns a list of files modified after the given date in the specified directory.
//func (cn *AClientSFTP) ListFilesAfterDate(dir string, afterDate time.Time) ([]string, error) {
//	cn.mu.RLock()
//	defer cn.mu.RUnlock()
//
//	if cn.connSFTP == nil {
//		return nil, fmt.Errorf("no active SFTP connection to host=%s", cn.address)
//	}
//
//	// Default to CDWorkingDir if dir is empty
//	if dir == "" {
//		dir = cn.CDWorkingDir
//		if dir == "" {
//			return nil, fmt.Errorf("no directory specified and CDWorkingDir is not set")
//		}
//	}
//
//	// Verify if the directory exists
//	fileInfo, err := cn.connSFTP.Stat(dir)
//	if err != nil {
//		return nil, fmt.Errorf("failed to access directory %s: %v", dir, err)
//	}
//
//	// Ensure it's a directory
//	if !fileInfo.IsDir() {
//		return nil, fmt.Errorf("%s is not a directory", dir)
//	}
//
//	// List files in the directory
//	files, err := cn.connSFTP.ReadDir(dir)
//	if err != nil {
//		return nil, fmt.Errorf("failed to read directory %s: %v", dir, err)
//	}
//
//	// Filter files modified after the given date
//	var recentFiles []string
//	for _, file := range files {
//		if file.ModTime().After(afterDate) {
//			recentFiles = append(recentFiles, file.Name())
//		}
//	}
//
//	return recentFiles, nil
//}
//
//// GetFileBytes retrieves a file from the SFTP server and returns its contents as a byte slice.
//func (cn *AClientSFTP) GetFileBytes(remoteFilePath string) ([]byte, error) {
//	cn.mu.RLock()
//	defer cn.mu.RUnlock()
//
//	if cn.connSFTP == nil {
//		return nil, fmt.Errorf("no active SFTP connection to host=%s", cn.address)
//	}
//
//	// Open remote file
//	remoteFile, err := cn.connSFTP.Open(remoteFilePath)
//	if err != nil {
//		return nil, fmt.Errorf("failed to open remote file %s: %v", remoteFilePath, err)
//	}
//	defer remoteFile.Close()
//
//	// Read file content into memory
//	fileBytes, err := io.ReadAll(remoteFile)
//	if err != nil {
//		return nil, fmt.Errorf("failed to read remote file %s: %v", remoteFilePath, err)
//	}
//
//	return fileBytes, nil
//}
//
//// DownloadFile retrieves a file from the SFTP server and saves it to a local path.
//func (cn *AClientSFTP) DownloadFile(remoteFilePath, localFilePath string) error {
//	cn.mu.RLock()
//	defer cn.mu.RUnlock()
//
//	if cn.connSFTP == nil {
//		return fmt.Errorf("no active SFTP connection to host=%s", cn.address)
//	}
//
//	// Open remote file
//	remoteFile, err := cn.connSFTP.Open(remoteFilePath)
//	if err != nil {
//		return fmt.Errorf("failed to open remote file %s: %v", remoteFilePath, err)
//	}
//	defer remoteFile.Close()
//
//	// Create local file
//	localFile, err := os.Create(localFilePath)
//	if err != nil {
//		return fmt.Errorf("failed to create local file %s: %v", localFilePath, err)
//	}
//	defer localFile.Close()
//
//	// Copy contents from remote file to local file
//	_, err = io.Copy(localFile, remoteFile)
//	if err != nil {
//		return fmt.Errorf("failed to copy remote file %s to local path %s: %v", remoteFilePath, localFilePath, err)
//	}
//
//	return nil
//}
