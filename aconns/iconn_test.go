package aconns

import (
	"encoding/json"
	"github.com/jpfluger/alibs-slim/auuids"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIConns_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		jsonStr string
		wantErr bool
	}{
		{
			name: "Valid JSON with Dummy adapters",
			jsonStr: `[
				{
					"id": "550e8400-e29b-41d4-a716-446655440000",
					"adapter": {
						"type": "dummy"
					}
				},
				{
					"id": "550e8400-e29b-41d4-a716-446655440001",
					"adapter": {
						"type": "dummydb"
					}
				}
			]`,
			wantErr: false,
		},
		{
			name: "Invalid JSON",
			jsonStr: `[
				{
					"id": "550e8400-e29b-41d4-a716-446655440000",
					"adapter": {
						"type": "dummy"
					}
				},
				{
					"id": "550e8400-e29b-41d4-a716-446655440001",
					"adapter": {
						"type": "dummydb"
					}
				`,
			wantErr: true,
		},
		{
			name: "Empty Adapter",
			jsonStr: `[
				{
					"id": "550e8400-e29b-41d4-a716-446655440000",
					"adapter": {}
				}
			]`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var conns IConns
			err := json.Unmarshal([]byte(tt.jsonStr), &conns)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, conns)
				for _, conn := range conns {
					assert.NotNil(t, conn.GetAdapter())
				}
			}
		})
	}
}

func TestIConns_ToMap(t *testing.T) {
	conns := IConns{
		&Conn{
			Id:      NewConnId(),
			Adapter: &DummyAdapter{},
		},
		&Conn{
			Id:      NewConnId(),
			Adapter: &DummyAdapter{},
		},
	}

	connMap := conns.ToMap()
	assert.Equal(t, len(conns), len(connMap))
	for _, conn := range conns {
		assert.Equal(t, conn, connMap[conn.GetId()])
	}
}

func TestIConns_GetSetRemove(t *testing.T) {
	conns := IConns{}
	conn := &Conn{
		Id:      NewConnId(),
		Adapter: &DummyAdapter{},
	}

	// Test Set
	conns.Set(conn)
	retrievedConn, exists := conns.Get(conn.GetId())
	assert.True(t, exists)
	assert.Equal(t, conn, retrievedConn)

	// Test Remove
	conns.Remove(conn.GetId())
	_, exists = conns.Get(conn.GetId())
	assert.False(t, exists)
}

func TestIConns_FindByConnId(t *testing.T) {
	conns := IConns{}
	conn := &Conn{
		Id:      NewConnId(),
		Adapter: &DummyAdapter{},
	}

	conns.Set(conn)
	retrievedConn, exists := conns.FindByConnId(conn.GetId())
	assert.True(t, exists)
	assert.Equal(t, conn, retrievedConn)
}

func TestIConns_FindByAdapterName(t *testing.T) {
	conns := IConns{}
	adapter := &DummyAdapter{}
	conn := &Conn{
		Id:      NewConnId(),
		Adapter: adapter,
	}

	conns.Set(conn)
	retrievedAdapter, exists := conns.FindByAdapterName(adapter.GetName())
	assert.True(t, exists)
	assert.Equal(t, adapter, retrievedAdapter)
}

func TestIConns_ToAdapterArray(t *testing.T) {
	conns := IConns{
		&Conn{
			Id:      NewConnId(),
			Adapter: &DummyAdapter{},
		},
		&Conn{
			Id:      NewConnId(),
			Adapter: &DummyAdapter{},
		},
	}

	adapters := conns.ToAdapterArray()
	assert.Equal(t, len(conns), len(adapters))
	for i, conn := range conns {
		assert.Equal(t, conn.GetAdapter(), adapters[i])
	}
}

func TestIConns_ToAdapterMap(t *testing.T) {
	conns := IConns{
		&Conn{
			Id:      NewConnId(),
			Adapter: &DummyAdapter{Adapter: Adapter{Name: "one"}},
		},
		&Conn{
			Id:      NewConnId(),
			Adapter: &DummyAdapter{Adapter: Adapter{Name: "two"}},
		},
	}

	adapterMap := conns.ToAdapterMap()
	assert.Equal(t, len(conns), len(adapterMap))
	for _, conn := range conns {
		assert.Equal(t, conn.GetAdapter(), adapterMap[conn.GetAdapter().GetName()])
	}
}

func TestIConnMap_GetSetRemove(t *testing.T) {
	connMap := IConnMap{}
	conn := &Conn{
		Id:      NewConnId(),
		Adapter: &DummyAdapter{},
	}

	// Test Set
	connMap.Set(conn.GetId(), conn)
	retrievedConn, exists := connMap.Get(conn.GetId())
	assert.True(t, exists)
	assert.Equal(t, conn, retrievedConn)

	// Test Remove
	connMap.Remove(conn.GetId())
	_, exists = connMap.Get(conn.GetId())
	assert.False(t, exists)
}

func TestIConnMap_ToArray(t *testing.T) {
	id1 := NewConnId()
	id2 := NewConnId()
	connMap := IConnMap{
		id1: &Conn{
			Id:      id1,
			Adapter: &DummyAdapter{Adapter: Adapter{Name: "one"}},
		},
		id2: &Conn{
			Id:      id2,
			Adapter: &DummyAdapter{Adapter: Adapter{Name: "two"}},
		},
	}

	conns := connMap.ToArray()
	assert.Equal(t, len(connMap), len(conns))
	for _, conn := range conns {
		assert.Equal(t, conn, connMap[conn.GetId()])
	}
}

func TestIConnMap_ToAdapterArray(t *testing.T) {
	connMap := IConnMap{
		NewConnId(): &Conn{
			Id:      NewConnId(),
			Adapter: &DummyAdapter{},
		},
		NewConnId(): &Conn{
			Id:      NewConnId(),
			Adapter: &DummyAdapter{},
		},
	}

	adapters := connMap.ToAdapterArray()
	assert.Equal(t, len(connMap), len(adapters))
	for _, conn := range connMap {
		assert.Contains(t, adapters, conn.GetAdapter())
	}
}

func TestIConnMap_ToAdapterMap(t *testing.T) {
	connMap := IConnMap{
		NewConnId(): &Conn{
			Id:      NewConnId(),
			Adapter: &DummyAdapter{Adapter: Adapter{Name: "one"}},
		},
		NewConnId(): &Conn{
			Id:      NewConnId(),
			Adapter: &DummyAdapter{Adapter: Adapter{Name: "two"}},
		},
	}

	adapterMap := connMap.ToAdapterMap()
	assert.Equal(t, len(connMap), len(adapterMap))
	for _, conn := range connMap {
		assert.Equal(t, conn.GetAdapter(), adapterMap[conn.GetAdapter().GetName()])
	}
}

func TestIConns_FilterByRole(t *testing.T) {
	masterId := NewConnId()
	authId := NewConnId()

	conns := IConns{
		&Conn{
			Id:    masterId,
			Roles: ConnRoles{CONNROLE_MASTER},
		},
		&Conn{
			Id:    authId,
			Roles: ConnRoles{CONNROLE_AUTH},
		},
		&Conn{
			Id:    NewConnId(),
			Roles: ConnRoles{CONNROLE_MASTER, CONNROLE_AUTH},
		},
	}

	filtered := conns.FilterByRole(CONNROLE_AUTH)
	for _, conn := range filtered {
		assert.Contains(t, conn.GetRoles(), CONNROLE_AUTH)
	}
	assert.GreaterOrEqual(t, len(filtered), 2)
}

func TestIConns_FindByTenantId(t *testing.T) {
	targetTenant := auuids.NewUUID()
	otherTenant := auuids.NewUUID()

	conns := IConns{
		&Conn{
			Id: NewConnId(),
			TenantInfo: ConnTenantInfo{
				TenantId: otherTenant,
			},
		},
		&Conn{
			Id: NewConnId(),
			TenantInfo: ConnTenantInfo{
				TenantId: targetTenant,
			},
		},
	}

	conn, found := conns.FindByTenantId(targetTenant)
	assert.True(t, found)
	assert.Equal(t, targetTenant.String(), conn.GetTenantInfo().TenantId.String())

	_, notFound := conns.FindByTenantId(auuids.NewUUID())
	assert.False(t, notFound)
}

func TestIConns_GetTenantInfos(t *testing.T) {
	id1 := auuids.NewUUID()
	id2 := auuids.NewUUID()

	conns := IConns{
		&Conn{
			TenantInfo: ConnTenantInfo{
				TenantId: id1,
				Region:   "us-west",
				Priority: 0,
			},
		},
		&Conn{
			TenantInfo: ConnTenantInfo{
				TenantId: id2,
				Region:   "eu-central",
				Priority: 1,
			},
		},
		&Conn{
			TenantInfo: ConnTenantInfo{}, // Empty TenantId, should be excluded
		},
	}

	infos := conns.GetTenantInfos()
	assert.Len(t, infos, 2)
	assert.True(t, infos.HasTenant(id1))
	assert.True(t, infos.HasTenant(id2))
}
