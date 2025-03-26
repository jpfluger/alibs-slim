package g_aconns

import (
	"fmt"
	"github.com/jpfluger/alibs-slim/aconns"
	"github.com/stretchr/testify/assert"
	"testing"
)

// DummyAdapter satisfies IAdapter for testing purposes.
type DummyAdapter struct {
	aconns.Adapter
	shouldFail bool
}

// Validate checks if the DummyAdapter is valid.
func (d DummyAdapter) Validate() error {
	return nil
}

// Test simulates testing the DummyAdapter.
func (d DummyAdapter) Test() (bool, aconns.TestStatus, error) {
	if d.shouldFail {
		return false, aconns.TESTSTATUS_FAILED, fmt.Errorf("dummy adapter test failed")
	}
	return true, aconns.TESTSTATUS_INITIALIZED_SUCCESSFUL, nil
}

func TestConnMapGlobal_GetSetRemove(t *testing.T) {
	connMap := CONNS()
	conn := &aconns.Conn{
		Id:      aconns.NewConnId(),
		Adapter: &DummyAdapter{},
	}

	// Test Set
	connMap.Set(conn)
	retrievedConn, exists := connMap.Get(conn.GetId())
	assert.True(t, exists)
	assert.Equal(t, conn, retrievedConn)
	assert.Equal(t, conn.GetId(), connMap.Index[conn.GetAdapter().GetName()])

	// Test Remove
	connMap.Remove(conn.GetId())
	_, exists = connMap.Get(conn.GetId())
	assert.False(t, exists)
	_, exists = connMap.Index[conn.GetAdapter().GetName()]
	assert.False(t, exists)
}

func TestConnMapGlobal_FindByAdapterName(t *testing.T) {
	connMap := CONNS()
	adapter := &DummyAdapter{}
	conn := &aconns.Conn{
		Id:      aconns.NewConnId(),
		Adapter: adapter,
	}

	connMap.Set(conn)
	retrievedAdapter, exists := connMap.FindByAdapterName(adapter.GetName())
	assert.True(t, exists)
	assert.Equal(t, adapter, retrievedAdapter)
}

func TestConnMapGlobal_Reset(t *testing.T) {
	connMap := CONNS()
	conn := &aconns.Conn{
		Id:      aconns.NewConnId(),
		Adapter: &DummyAdapter{},
	}

	connMap.Set(conn)
	assert.NotEmpty(t, connMap.Map)
	assert.NotEmpty(t, connMap.Index)

	connMap.Reset()
	assert.Empty(t, connMap.Map)
	assert.Empty(t, connMap.Index)
}

func TestConnMapGlobal_SetByIConns(t *testing.T) {
	connMap.Reset()
	connMap := CONNS()
	conns := aconns.IConns{
		&aconns.Conn{
			Id:      aconns.NewConnId(),
			Adapter: &DummyAdapter{Adapter: aconns.Adapter{Name: "one"}},
		},
		&aconns.Conn{
			Id:      aconns.NewConnId(),
			Adapter: &DummyAdapter{Adapter: aconns.Adapter{Name: "two"}},
		},
	}

	connMap.SetByIConns(conns)
	assert.Equal(t, len(conns), len(connMap.Map))
	assert.Equal(t, len(conns), len(connMap.Index))
	for _, conn := range conns {
		retrievedConn, exists := connMap.Get(conn.GetId())
		assert.True(t, exists)
		assert.Equal(t, conn, retrievedConn)
		assert.Equal(t, conn.GetId(), connMap.Index[conn.GetAdapter().GetName()])
	}
}
