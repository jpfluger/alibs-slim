package aconns

import (
	"encoding/json"
	"fmt"
	"github.com/jpfluger/alibs-slim/areflect"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func init() {
	_ = areflect.TypeManager().Register(TYPEMANAGER_CONNADAPTERS, "aconn-tests", returnTypeManagerConnAdapters)
}

func returnTypeManagerConnAdapters(typeName string) (reflect.Type, error) {
	var rtype reflect.Type // nil is the zero value for pointers, maps, slices, channels, and function types, interfaces, and other compound types.
	switch AdapterType(typeName) {
	case ADAPTERTYPE_DUMMY:
		// Return the type of NoteFlag if typeName is "flag".
		rtype = reflect.TypeOf(DummyAdapter{})
	case ADAPTERTYPE_DUMMYDB:
		// Return the type of NoteFlag if typeName is "flag".
		rtype = reflect.TypeOf(DummyAdapterDB{})
	}
	// Return the determined reflect.Type and no error.
	return rtype, nil
}

const ADAPTERTYPE_DUMMY AdapterType = "dummy"
const ADAPTERTYPE_DUMMYDB AdapterType = "dummydb"

// DummyAdapter satisfies IAdapter
type DummyAdapter struct {
	Adapter
	shouldFail bool
}

// Validate checks if the DummyAdapter is valid.
func (d DummyAdapter) Validate() error {
	return nil
}

// Test simulates testing the DummyAdapter.
func (d DummyAdapter) Test() (bool, TestStatus, error) {
	if d.shouldFail {
		return false, TESTSTATUS_FAILED, fmt.Errorf("dummy adapter test failed")
	}
	return true, TESTSTATUS_INITIALIZED_SUCCESSFUL, nil
}

func TestConn_GetId(t *testing.T) {
	id := NewConnId()
	conn := &Conn{Id: id}
	assert.Equal(t, id, conn.GetId())
}

func TestConn_GetAdapter(t *testing.T) {
	adapter := &DummyAdapter{}
	conn := &Conn{Adapter: adapter}
	assert.Equal(t, adapter, conn.GetAdapter())
}

func TestConn_Validate(t *testing.T) {
	tests := []struct {
		name    string
		conn    *Conn
		wantErr bool
	}{
		{"Valid Conn", &Conn{Adapter: &DummyAdapter{}}, false},
		{"Nil Adapter", &Conn{Adapter: nil}, true},
		{"Auto-assign ID", &Conn{Adapter: &DummyAdapter{}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.conn.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Conn.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && tt.conn.Id.IsNil() {
				t.Errorf("Conn.Validate() did not auto-assign ID")
			}
		})
	}
}

func TestConn_UnmarshalJSON(t *testing.T) {
	jsonData := `{
		"id": "550e8400-e29b-41d4-a716-446655440000",
		"adapter": {
			"type": "dummy"
		}
	}`
	conn := &Conn{}
	err := json.Unmarshal([]byte(jsonData), conn)
	assert.NoError(t, err)
	assert.Equal(t, ParseConnId("550e8400-e29b-41d4-a716-446655440000").String(), conn.GetId().String())
	assert.IsType(t, &DummyAdapter{}, conn.GetAdapter())
}

func TestConn_UnmarshalJSON_InvalidAdapter(t *testing.T) {
	jsonData := `{
		"id": "550e8400-e29b-41d4-a716-446655440000",
		"adapter": {
			"type": "invalid"
		}
	}`
	conn := &Conn{}
	err := json.Unmarshal([]byte(jsonData), conn)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot find type struct 'invalid'")
}

func TestConn_Test(t *testing.T) {
	tests := []struct {
		name       string
		conn       *Conn
		wantOk     bool
		wantStatus TestStatus
		wantErr    bool
	}{
		{
			name: "Successful Test",
			conn: &Conn{
				Adapter: &DummyAdapter{},
			},
			wantOk:     true,
			wantStatus: TESTSTATUS_INITIALIZED_SUCCESSFUL,
			wantErr:    false,
		},
		{
			name: "Failed Test",
			conn: &Conn{
				Adapter: &DummyAdapter{shouldFail: true},
			},
			wantOk:     false,
			wantStatus: TESTSTATUS_FAILED,
			wantErr:    true,
		},
		{
			name: "Ignored Connection",
			conn: &Conn{
				Ignore: true,
			},
			wantOk:     true,
			wantStatus: TESTSTATUS_INITIALIZED,
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ok, status, err := tt.conn.Test()
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

func TestConn_TestBootstrap(t *testing.T) {
	tests := []struct {
		name       string
		conn       *Conn
		wantOk     bool
		wantStatus TestStatus
		wantErr    bool
	}{
		{
			name: "Bootstrap Test Successful",
			conn: &Conn{
				IsBootstrap: true,
				Adapter:     &DummyAdapter{},
			},
			wantOk:     true,
			wantStatus: TESTSTATUS_INITIALIZED_SUCCESSFUL,
			wantErr:    false,
		},
		{
			name: "Bootstrap Test Failed",
			conn: &Conn{
				IsBootstrap: true,
				Adapter:     &DummyAdapter{shouldFail: true},
			},
			wantOk:     false,
			wantStatus: TESTSTATUS_FAILED,
			wantErr:    true,
		},
		{
			name: "Not Bootstrap",
			conn: &Conn{
				IsBootstrap: false,
				Adapter:     &DummyAdapter{},
			},
			wantOk:     true,
			wantStatus: TESTSTATUS_INITIALIZED,
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ok, status, err := tt.conn.TestBootstrap()
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

func TestConn_TestRequired(t *testing.T) {
	tests := []struct {
		name       string
		conn       *Conn
		wantOk     bool
		wantStatus TestStatus
		wantErr    bool
	}{
		{
			name: "Required Test Successful",
			conn: &Conn{
				IsRequired: true,
				Adapter:    &DummyAdapter{},
			},
			wantOk:     true,
			wantStatus: TESTSTATUS_INITIALIZED_SUCCESSFUL,
			wantErr:    false,
		},
		{
			name: "Required Test Failed",
			conn: &Conn{
				IsRequired: true,
				Adapter:    &DummyAdapter{shouldFail: true},
			},
			wantOk:     false,
			wantStatus: TESTSTATUS_FAILED,
			wantErr:    true,
		},
		{
			name: "Not Required",
			conn: &Conn{
				IsRequired: false,
				Adapter:    &DummyAdapter{},
			},
			wantOk:     true,
			wantStatus: TESTSTATUS_INITIALIZED,
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ok, status, err := tt.conn.TestRequired()
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
