package aclient_badger

import (
	"os"
	"testing"

	"github.com/jpfluger/alibs-slim/aconns"
	"github.com/jpfluger/alibs-slim/acrypt"
	"github.com/stretchr/testify/assert"
)

func TestAClientBadger_Validate(t *testing.T) {
	cn := &AClientBadger{
		ADBAdapterBase: aconns.ADBAdapterBase{
			Adapter: aconns.Adapter{
				Type: ADAPTERTYPE_BADGER,
				Name: aconns.AdapterName("test_badger"),
				Host: "local",
				Port: 0,
			},
			Database: "/tmp/test_db",
			Username: "badger",
			Password: "",
		},
	}

	err := cn.Validate()
	assert.NoError(t, err)
}

func TestAClientBadger_Validate_Invalid(t *testing.T) {
	cn := &AClientBadger{
		ADBAdapterBase: aconns.ADBAdapterBase{
			Adapter: aconns.Adapter{
				Type: ADAPTERTYPE_BADGER,
				Name: aconns.AdapterName("test_badger"),
			},
			Database: "", // Empty database
		},
	}

	err := cn.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "database (directory path) is empty")
}

func TestAClientBadger_Test_Success(t *testing.T) {
	tempDir := t.TempDir()
	cn := &AClientBadger{
		ADBAdapterBase: aconns.ADBAdapterBase{
			Adapter: aconns.Adapter{
				Type: ADAPTERTYPE_BADGER,
				Name: aconns.AdapterName("test_badger"),
				Host: "local",
				Port: 0,
			},
			Database: tempDir,
			Username: "badger",
			Password: "",
		},
	}

	ok, status, err := cn.Test()
	assert.True(t, ok)
	assert.Equal(t, aconns.TESTSTATUS_INITIALIZED_SUCCESSFUL, status)
	assert.NoError(t, err)

	// Clean up
	cn.CloseConnection()
}

func TestAClientBadger_Test_Success_WithEncryption(t *testing.T) {
	pass, err := acrypt.GenerateEncryptionKeyWithLengthBase64(16)
	if err != nil {
		t.Error(err)
		return
	}

	tempDir := t.TempDir()
	cn := &AClientBadger{
		ADBAdapterBase: aconns.ADBAdapterBase{
			Adapter: aconns.Adapter{
				Type: ADAPTERTYPE_BADGER,
				Name: aconns.AdapterName("test_badger"),
				Host: "local",
				Port: 0,
			},
			Database: tempDir,
			Username: "badger",
			Password: pass, // "thisis16byteskey", // 16-byte password for AES-128 encryption
		},
	}

	ok, status, err := cn.Test()
	assert.True(t, ok)
	assert.Equal(t, aconns.TESTSTATUS_INITIALIZED_SUCCESSFUL, status)
	assert.NoError(t, err)

	// Clean up
	cn.CloseConnection()
}

func TestAClientBadger_Test_Failure(t *testing.T) {
	tempFile, err := os.CreateTemp("", "badger_test_file")
	assert.NoError(t, err)
	defer os.Remove(tempFile.Name())

	cn := &AClientBadger{
		ADBAdapterBase: aconns.ADBAdapterBase{
			Adapter: aconns.Adapter{
				Type: ADAPTERTYPE_BADGER,
				Name: aconns.AdapterName("test_badger"),
			},
			Database: tempFile.Name(), // Path to a file, not a directory
		},
	}

	ok, status, err := cn.Test()
	assert.False(t, ok)
	assert.Equal(t, aconns.TESTSTATUS_FAILED, status)
	assert.Error(t, err)
}

func TestAClientBadger_OpenConnection_Success(t *testing.T) {
	tempDir := t.TempDir()
	cn := &AClientBadger{
		ADBAdapterBase: aconns.ADBAdapterBase{
			Adapter: aconns.Adapter{
				Type: ADAPTERTYPE_BADGER,
				Name: aconns.AdapterName("test_badger"),
				Host: "local",
				Port: 0,
			},
			Database: tempDir,
			Username: "badger",
			Password: "",
		},
	}

	err := cn.OpenConnection()
	assert.NoError(t, err)
	assert.NotNil(t, cn.DB())

	// Clean up
	cn.CloseConnection()
}

func TestAClientBadger_OpenConnection_Failure(t *testing.T) {
	tempFile, err := os.CreateTemp("", "badger_test_file")
	assert.NoError(t, err)
	defer os.Remove(tempFile.Name())

	cn := &AClientBadger{
		ADBAdapterBase: aconns.ADBAdapterBase{
			Adapter: aconns.Adapter{
				Type: ADAPTERTYPE_BADGER,
				Name: aconns.AdapterName("test_badger"),
			},
			Database: tempFile.Name(), // Path to a file, not a directory
		},
	}

	err = cn.OpenConnection()
	assert.Error(t, err)
	assert.Nil(t, cn.db)
}

func TestAClientBadger_CloseConnection(t *testing.T) {
	tempDir := t.TempDir()
	cn := &AClientBadger{
		ADBAdapterBase: aconns.ADBAdapterBase{
			Adapter: aconns.Adapter{
				Type: ADAPTERTYPE_BADGER,
				Name: aconns.AdapterName("test_badger"),
				Host: "local",
				Port: 0,
			},
			Database: tempDir,
			Username: "badger",
			Password: "",
		},
	}

	// Open first
	err := cn.OpenConnection()
	assert.NoError(t, err)

	err = cn.CloseConnection()
	assert.NoError(t, err)
	assert.Nil(t, cn.db)
}

func TestAClientBadger_CloseConnection_NoDB(t *testing.T) {
	cn := &AClientBadger{}

	err := cn.CloseConnection()
	assert.NoError(t, err)
	assert.Nil(t, cn.db)
}

func TestAClientBadger_GetAddress(t *testing.T) {
	cn := &AClientBadger{
		ADBAdapterBase: aconns.ADBAdapterBase{
			Adapter: aconns.Adapter{
				Type: ADAPTERTYPE_BADGER,
				Name: aconns.AdapterName("test_badger"),
				Host: "local",
				Port: 0,
			},
		},
	}

	address := cn.GetAddress()
	assert.Equal(t, "local:0", address)
}

func TestAClientBadger_Refresh(t *testing.T) {
	tempDir := t.TempDir()
	cn := &AClientBadger{
		ADBAdapterBase: aconns.ADBAdapterBase{
			Adapter: aconns.Adapter{
				Type: ADAPTERTYPE_BADGER,
				Name: aconns.AdapterName("test_badger"),
				Host: "local",
				Port: 0,
			},
			Database: tempDir,
			Username: "badger",
			Password: "",
		},
	}

	// Open first
	err := cn.OpenConnection()
	assert.NoError(t, err)
	oldDB := cn.db
	assert.NotNil(t, oldDB)

	// Refresh
	err = cn.Refresh()
	assert.NoError(t, err)
	newDB := cn.db
	assert.NotNil(t, newDB)
	assert.NotEqual(t, oldDB, newDB) // Different pointer after reopen

	// Clean up
	cn.CloseConnection()
}

func TestAClientBadger_GetBadgerStore(t *testing.T) {
	tempDir := t.TempDir()
	cn := &AClientBadger{
		ADBAdapterBase: aconns.ADBAdapterBase{
			Adapter: aconns.Adapter{
				Type: ADAPTERTYPE_BADGER,
				Name: aconns.AdapterName("test_badger"),
				Host: "local",
				Port: 0,
			},
			Database: tempDir,
			Username: "badger",
			Password: "",
		},
	}

	// Open first
	err := cn.OpenConnection()
	assert.NoError(t, err)

	store := cn.GetBadgerStore("test_prefix")
	assert.NotNil(t, store)

	// Clean up
	cn.CloseConnection()
}

func TestAClientBadger_GetBadgerStore_NoDB(t *testing.T) {
	cn := &AClientBadger{}

	store := cn.GetBadgerStore("test_prefix")
	assert.Nil(t, store)
}
