package aclient_sftp

import (
	"github.com/jpfluger/alibs-slim/aconns"
	"github.com/pkg/sftp"
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestAClientSFTP_Validate checks validation logic
func TestAClientSFTP_Validate(t *testing.T) {
	client := &AClientSFTP{
		ADBAdapterBase: aconns.ADBAdapterBase{
			Adapter: aconns.Adapter{
				Type: aconns.AdapterType("sftp"),
				Name: aconns.AdapterName("test_sftp"),
				Host: "localhost",
				Port: 22,
			},
			Username: "user",
			Password: "pass",
		},
		ConnectionTimeout:  30,
		InsecureSkipVerify: true,
		CDWorkingDir:       "/testdir",
	}

	err := client.Validate()
	assert.NoError(t, err)
	assert.Equal(t, "localhost:22", client.address)
}

// TestAClientSFTP_OpenConnection tests opening an SFTP connection
func TestAClientSFTP_OpenConnection(t *testing.T) {
	client := &AClientSFTP{
		ADBAdapterBase: aconns.ADBAdapterBase{
			Adapter: aconns.Adapter{
				Type: aconns.AdapterType("sftp"),
				Name: aconns.AdapterName("test_sftp"),
				Host: "localhost",
				Port: 22,
			},
			Username: "user",
			Password: "pass",
		},
		ConnectionTimeout:  30,
		InsecureSkipVerify: true,
		CDWorkingDir:       "/testdir",
	}

	err := client.OpenConnection()
	assert.Error(t, err) // Expecting an error because there's no actual SFTP server running
}

// TestAClientSFTP_TestConnection tests if an active connection is correctly detected
func TestAClientSFTP_TestConnection(t *testing.T) {
	client := &AClientSFTP{
		ADBAdapterBase: aconns.ADBAdapterBase{
			Adapter: aconns.Adapter{
				Type: aconns.AdapterType("sftp"),
				Name: aconns.AdapterName("test_sftp"),
				Host: "localhost",
				Port: 22,
			},
			Username: "user",
			Password: "pass",
		},
		ConnectionTimeout: 30,
	}

	err := client.TestConnection()
	assert.Error(t, err) // Should return an error because no SFTP connection exists
}

// TestAClientSFTP_CloseConnection verifies connection closure logic
func TestAClientSFTP_CloseConnection(t *testing.T) {
	client := &AClientSFTP{
		ADBAdapterBase: aconns.ADBAdapterBase{
			Adapter: aconns.Adapter{
				Type: aconns.AdapterType("sftp"),
				Name: aconns.AdapterName("test_sftp"),
				Host: "localhost",
				Port: 22,
			},
			Username: "user",
			Password: "pass",
		},
		ConnectionTimeout: 30,
	}

	// Simulate open connections
	client.sshConn = nil
	client.connSFTP = nil

	err := client.CloseConnection()
	assert.NoError(t, err) // Should close without errors

	assert.Nil(t, client.sshConn)
	assert.Nil(t, client.connSFTP)
}

// TestAClientSFTP_SFTPClient tests retrieval of the active SFTP connection
func TestAClientSFTP_SFTPClient(t *testing.T) {
	client := &AClientSFTP{
		ADBAdapterBase: aconns.ADBAdapterBase{
			Adapter: aconns.Adapter{
				Type: aconns.AdapterType("sftp"),
				Name: aconns.AdapterName("test_sftp"),
				Host: "localhost",
				Port: 22,
			},
			Username: "user",
			Password: "pass",
		},
		ConnectionTimeout: 30,
	}

	assert.Nil(t, client.SFTPClient()) // Should be nil initially

	// Simulate a connection
	client.connSFTP = &sftp.Client{}

	assert.NotNil(t, client.SFTPClient()) // Should return a non-nil SFTP client
}
