package aclient_ftp

import (
	"github.com/jpfluger/alibs-slim/aconns"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAClientFTP_Validate(t *testing.T) {
	client := &AClientFTP{
		ADBAdapterBase: aconns.ADBAdapterBase{
			Adapter: aconns.Adapter{
				Type: aconns.AdapterType("ftp"),
				Name: aconns.AdapterName("test_ftp"),
				Host: "localhost",
				Port: 21,
			},
			Username: "user",
			Password: "pass",
		},
		ConnectionTimeout:  30,
		InsecureSkipVerify: true,
		UseFTPS:            false,
		CDWorkingDir:       "/testdir",
	}

	err := client.Validate()
	assert.NoError(t, err)
	assert.Equal(t, "localhost:21", client.address)
}

func TestAClientFTP_OpenConnection(t *testing.T) {
	client := &AClientFTP{
		ADBAdapterBase: aconns.ADBAdapterBase{
			Adapter: aconns.Adapter{
				Type: aconns.AdapterType("ftp"),
				Name: aconns.AdapterName("test_ftp"),
				Host: "localhost",
				Port: 21,
			},
			Username: "user",
			Password: "pass",
		},
		ConnectionTimeout:  30,
		InsecureSkipVerify: true,
		UseFTPS:            false,
		CDWorkingDir:       "/testdir",
	}

	err := client.OpenConnection()
	assert.Error(t, err) // Expecting an error because there's no actual FTP server running
}

func TestAClientFTP_OpenConnection_FTPS(t *testing.T) {
	client := &AClientFTP{
		ADBAdapterBase: aconns.ADBAdapterBase{
			Adapter: aconns.Adapter{
				Type: aconns.AdapterType("ftp"),
				Name: aconns.AdapterName("test_ftp"),
				Host: "localhost",
				Port: 21,
			},
			Username: "user",
			Password: "pass",
		},
		ConnectionTimeout:  30,
		InsecureSkipVerify: true,
		UseFTPS:            true,
		CDWorkingDir:       "/testdir",
	}

	err := client.OpenConnection()
	assert.Error(t, err) // Expecting an error because there's no actual FTPS server running
}
