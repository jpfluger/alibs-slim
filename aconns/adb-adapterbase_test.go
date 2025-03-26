package aconns

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// DummyAdapterDB satisfies IAdapterDB
type DummyAdapterDB struct {
	ADBAdapterBase
	shouldFail bool
}

// Validate checks if the DummyAdapterDB is valid.
func (d *DummyAdapterDB) Validate() error {
	if d.shouldFail {
		return fmt.Errorf("dummy adapter validation failed")
	}
	return d.ADBAdapterBase.Validate()
}

// Test simulates testing the DummyAdapterDB.
func (d *DummyAdapterDB) Test() (bool, TestStatus, error) {
	if d.shouldFail {
		return false, TESTSTATUS_FAILED, fmt.Errorf("dummy adapter test failed")
	}
	return d.ADBAdapterBase.Test()
}

func TestADBAdapterBase_Validate(t *testing.T) {
	tests := []struct {
		name    string
		adapter *DummyAdapterDB
		wantErr bool
	}{
		{
			name: "Valid Adapter",
			adapter: &DummyAdapterDB{
				ADBAdapterBase: ADBAdapterBase{
					Adapter: Adapter{
						Type: AdapterType("type1"),
						Name: AdapterName("name1"),
						Host: "localhost",
						Port: 3306,
					},
					Database: "testdb",
					Username: "user",
					Password: "password",
				},
			},
			wantErr: false,
		},
		{
			name: "Invalid Adapter - Missing Fields",
			adapter: &DummyAdapterDB{
				ADBAdapterBase: ADBAdapterBase{
					Adapter: Adapter{
						Type: AdapterType(""),
						Name: AdapterName(""),
						Host: "",
						Port: 0,
					},
					Database: "",
					Username: "",
					Password: "",
				},
			},
			wantErr: true,
		},
		{
			name: "Validation Failure",
			adapter: &DummyAdapterDB{
				shouldFail: true,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.adapter.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestADBAdapterBase_Test(t *testing.T) {
	tests := []struct {
		name       string
		adapter    *DummyAdapterDB
		wantOk     bool
		wantStatus TestStatus
		wantErr    bool
	}{
		{
			name: "Successful Test",
			adapter: &DummyAdapterDB{
				ADBAdapterBase: ADBAdapterBase{
					Adapter: Adapter{
						Type: AdapterType("type1"),
						Name: AdapterName("name1"),
						Host: "localhost",
						Port: 3306,
					},
					Database: "testdb",
					Username: "user",
					Password: "password",
				},
			},
			wantOk:     true,
			wantStatus: TESTSTATUS_INITIALIZED_SUCCESSFUL,
			wantErr:    false,
		},
		{
			name: "Failed Test",
			adapter: &DummyAdapterDB{
				shouldFail: true,
			},
			wantOk:     false,
			wantStatus: TESTSTATUS_FAILED,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ok, status, err := tt.adapter.Test()
			assert.Equal(t, tt.wantOk, ok)
			assert.Equal(t, tt.wantStatus, status)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
