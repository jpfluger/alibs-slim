package aclient_sftp

import (
	"fmt"
	"github.com/jpfluger/alibs-slim/aconns"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"io"
	"os"
	"strings"
	"sync"
	"time"
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

	if err := cn.validate(); err != nil {
		return false, aconns.TESTSTATUS_FAILED, err
	}

	if cn.connSFTP != nil {
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
func (cn *AClientSFTP) OpenConnection() error {
	cn.mu.Lock()
	defer cn.mu.Unlock()
	return cn.openConnection()
}

// openConnection opens a persistent connection to the SFTP server.
func (cn *AClientSFTP) openConnection() error {
	if err := cn.validate(); err != nil {
		return err
	}

	// SSH client configuration
	config := &ssh.ClientConfig{
		User: cn.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(cn.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // Note: Use proper host key validation in production
		Timeout:         time.Duration(cn.ConnectionTimeout) * time.Second,
	}

	// Connect to SSH server
	sshConn, err := ssh.Dial("tcp", cn.address, config)
	if err != nil {
		return fmt.Errorf("failed to dial: %v", err)
	}

	// Create SFTP client
	sftpClient, err := sftp.NewClient(sshConn)
	if err != nil {
		_ = sshConn.Close() // Close SSH connection if SFTP client creation fails
		return fmt.Errorf("failed to create SFTP client: %v", err)
	}

	// Store connections in struct for later use and explicit closure
	cn.sshConn = sshConn
	cn.connSFTP = sftpClient

	// Test connection immediately after opening
	if err := cn.testConnection(); err != nil {
		_ = cn.closeConnection()
		return fmt.Errorf("connection test failed: %v", err)
	}

	return nil
}

// CloseConnection closes the SFTP and SSH connections.
func (cn *AClientSFTP) CloseConnection() error {
	cn.mu.Lock()
	defer cn.mu.Unlock()
	return cn.closeConnection()
}

// closeConnection closes the SFTP and SSH connections.
func (cn *AClientSFTP) closeConnection() error {
	var closeErrors []error

	// Close SFTP connection if it exists
	if cn.connSFTP != nil {
		if err := cn.connSFTP.Close(); err != nil {
			closeErrors = append(closeErrors, fmt.Errorf("failed to close SFTP connection to host=%s: %v", cn.address, err))
		}
		cn.connSFTP = nil
	}

	// Close SSH connection if it exists
	if cn.sshConn != nil {
		if err := cn.sshConn.Close(); err != nil {
			closeErrors = append(closeErrors, fmt.Errorf("failed to close SSH connection to host=%s: %v", cn.address, err))
		}
		cn.sshConn = nil
	}

	// Aggregate and return errors if any occurred
	if len(closeErrors) > 0 {
		return fmt.Errorf("errors occurred while closing connections: %v", closeErrors)
	}

	return nil
}

// TestConnection tests if the SFTP connection is active.
func (cn *AClientSFTP) TestConnection() error {
	cn.mu.RLock()
	defer cn.mu.RUnlock()
	return cn.testConnection()
}

// testConnection tests if the SFTP connection is active and validates the working directory.
func (cn *AClientSFTP) testConnection() error {
	if cn.connSFTP == nil {
		return fmt.Errorf("no active SFTP connection to host=%s", cn.address)
	}

	// Perform a basic operation to validate the connection
	wd, err := cn.connSFTP.Getwd()
	if err != nil {
		return fmt.Errorf("SFTP connection is invalid where host=%s: %v", cn.address, err)
	}

	fmt.Println(wd)

	//// Check if the working directory exists
	//if cn.CDWorkingDir != "" {
	//	if _, err := cn.connSFTP.Stat(cn.CDWorkingDir); err != nil {
	//		return fmt.Errorf("working directory %s is not accessible: %v", cn.CDWorkingDir, err)
	//	}
	//}
	//fmt.Printf("SFTP connection to %s is valid. Current directory: %s\n", cn.address, wd)

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

// GetFileBytes retrieves a file from the SFTP server and returns its contents as a byte slice.
func (cn *AClientSFTP) GetFileBytes(remoteFilePath string) ([]byte, error) {
	cn.mu.RLock()
	defer cn.mu.RUnlock()

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
	defer cn.mu.RUnlock()

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
