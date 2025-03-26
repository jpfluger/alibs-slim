package g_aconns

import (
	"encoding/json"
	"github.com/jpfluger/alibs-slim/aconns"
	"github.com/jpfluger/alibs-slim/aconns/aclient-http"
	"github.com/jpfluger/alibs-slim/aconns/adb-pg"
	"github.com/jpfluger/alibs-slim/anetwork"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestManager_MarshalUnmarshalJSON(t *testing.T) {
	// Create a new manager and add connections
	manager := NewManager()
	httpConn := &aconns.Conn{
		Id:      aconns.NewConnId(),
		Adapter: &aclient_http.AClientHTTP{Type: aclient_http.ADAPTERTYPE_HTTP, Name: "http"},
	}
	pgConn := &aconns.Conn{
		Id: aconns.NewConnId(),
		Adapter: &adb_pg.ADBPG{
			ADBAdapterBase: aconns.ADBAdapterBase{
				Adapter: aconns.Adapter{Type: adb_pg.ADAPTERTYPE_PG, Name: "pg"},
			},
		},
	}

	manager.Conns = append(manager.Conns, httpConn, pgConn)

	// Marshal the manager to JSON
	data, err := json.Marshal(manager)
	assert.NoError(t, err)

	// Unmarshal the JSON back to a manager
	var newManager Manager
	err = json.Unmarshal(data, &newManager)
	assert.NoError(t, err)

	// Validate the unmarshaled manager
	assert.True(t, newManager.HasConns())
	assert.Equal(t, len(manager.Conns), len(newManager.Conns))
}

func TestManager_Validate(t *testing.T) {
	manager := NewManager()
	httpConn := &aconns.Conn{
		Id:      aconns.NewConnId(),
		Adapter: &aclient_http.AClientHTTP{Type: aclient_http.ADAPTERTYPE_HTTP, Name: "http", Url: *anetwork.MustParseNetURL("https://example.com")},
	}
	pgConn := &aconns.Conn{
		Id: aconns.NewConnId(),
		Adapter: &adb_pg.ADBPG{
			ADBAdapterBase: aconns.ADBAdapterBase{
				Adapter:  aconns.Adapter{Type: adb_pg.ADAPTERTYPE_PG, Name: "pg", Host: "localhost"},
				Database: "db",
				Username: "username",
				Password: "password",
			},
		},
	}

	manager.Conns = append(manager.Conns, httpConn, pgConn)

	err := manager.Validate()
	assert.NoError(t, err)
}

func TestManager_ValidateBootstrap(t *testing.T) {
	manager := NewManager()
	httpConn := &aconns.Conn{
		Id:          aconns.NewConnId(),
		Adapter:     &aclient_http.AClientHTTP{Type: aclient_http.ADAPTERTYPE_HTTP, Name: "http", Url: *anetwork.MustParseNetURL("https://example.com")},
		IsBootstrap: true,
	}
	pgConn := &aconns.Conn{
		Id: aconns.NewConnId(),
		Adapter: &adb_pg.ADBPG{
			ADBAdapterBase: aconns.ADBAdapterBase{
				Adapter: aconns.Adapter{Type: adb_pg.ADAPTERTYPE_PG, Name: "pg", Host: "localhost"},
			},
		},
	}

	manager.Conns = append(manager.Conns, httpConn, pgConn)

	err := manager.ValidateBootstrap()
	assert.NoError(t, err)
}

func TestManager_ValidateRequired(t *testing.T) {
	manager := NewManager()
	httpConn := &aconns.Conn{
		Id:         aconns.NewConnId(),
		Adapter:    &aclient_http.AClientHTTP{Type: aclient_http.ADAPTERTYPE_HTTP, Name: "http", Url: *anetwork.MustParseNetURL("https://example.com")},
		IsRequired: true,
	}
	pgConn := &aconns.Conn{
		Id: aconns.NewConnId(),
		Adapter: &adb_pg.ADBPG{
			ADBAdapterBase: aconns.ADBAdapterBase{
				Adapter: aconns.Adapter{Type: adb_pg.ADAPTERTYPE_PG, Name: "pg", Host: "localhost"},
			},
		},
	}

	manager.Conns = append(manager.Conns, httpConn, pgConn)

	err := manager.ValidateRequired()
	assert.NoError(t, err)
}

func TestManager_FindConn(t *testing.T) {
	manager := NewManager()
	conn := &aconns.Conn{
		Id:      aconns.NewConnId(),
		Adapter: &DummyAdapter{},
	}

	manager.Conns.Set(conn)
	retrievedConn := manager.FindConn(conn.GetId())
	assert.Equal(t, conn, retrievedConn)
}

func TestManager_FindAdapter(t *testing.T) {
	manager := NewManager()
	adapter := &DummyAdapter{}
	conn := &aconns.Conn{
		Id:      aconns.NewConnId(),
		Adapter: adapter,
	}

	manager.Conns.Set(conn)
	retrievedAdapter := manager.FindAdapter(adapter.GetName())
	assert.Equal(t, adapter, retrievedAdapter)
}
