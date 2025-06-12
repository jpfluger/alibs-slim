package aconns

import (
	"encoding/json"
	"fmt"
	"github.com/jpfluger/alibs-slim/auuids"
)

const TYPEMANAGER_CONNADAPTERS = "connadapters"

// IConn interface defines the methods that a connection should implement.
type IConn interface {
	DoIgnore() bool
	GetId() ConnId
	GetAdapter() IAdapter
	GetIsRequired() bool
	GetIsBootstrap() bool
	Validate() error
	GetTenantInfo() *ConnTenantInfo
	GetAuthScopes() AuthScopes
	GetAuthUsages() AuthUsages
}

type IConns []IConn

func (conns IConns) GetConnCounts() (cTotal int, cIgnore int, cBootstrap int, cIsRequired int) {
	if conns == nil || len(conns) == 0 {
		return
	}
	for _, conn := range conns {
		cTotal++
		if conn.DoIgnore() {
			cIgnore++
		}
		if conn.GetIsBootstrap() {
			cBootstrap++
		}
		if conn.GetIsRequired() {
			cIsRequired++
		}
	}
	return
}

// UnmarshalJSON unmarshals JSON data into a slice of Conn structs.
func (conns *IConns) UnmarshalJSON(data []byte) error {
	var rawConns []json.RawMessage
	if err := json.Unmarshal(data, &rawConns); err != nil {
		return fmt.Errorf("failed to unmarshal IConns: %v", err)
	}

	for _, rawConn := range rawConns {
		var conn Conn
		if err := json.Unmarshal(rawConn, &conn); err != nil {
			return fmt.Errorf("failed to unmarshal Conn: %v", err)
		}
		*conns = append(*conns, &conn)
	}

	return nil
}

// ToMap converts the IConns to an IConnMap.
func (conns IConns) ToMap() IConnMap {
	connMap := make(IConnMap)
	for _, conn := range conns {
		if conn != nil {
			if conn.GetId().IsNil() || conn.GetAdapter() == nil {
				continue
			}
			connMap[conn.GetId()] = conn
		}
	}
	return connMap
}

// ToAdapterArray converts the IConns to an array of IAdapters.
func (conns IConns) ToAdapterArray() IAdapters {
	adapters := IAdapters{}
	for _, conn := range conns {
		if conn != nil && conn.GetAdapter() != nil {
			adapters = append(adapters, conn.GetAdapter())
		}
	}
	return adapters
}

// ToAdapterMap converts the IConns to an IAdapterMap.
func (conns IConns) ToAdapterMap() IAdapterMap {
	adapters := IAdapterMap{}
	for _, conn := range conns {
		if conn != nil && conn.GetAdapter() != nil {
			adapters[conn.GetAdapter().GetName()] = conn.GetAdapter()
		}
	}
	return adapters
}

// Get retrieves an IConn by its UUID.
func (conns IConns) Get(id ConnId) (IConn, bool) {
	for _, conn := range conns {
		if conn != nil && conn.GetId() == id {
			return conn, true
		}
	}
	return nil, false
}

// Set adds or updates an IConn in the slice.
func (conns *IConns) Set(conn IConn) {
	if conn != nil && !conn.GetId().IsNil() && conn.GetAdapter() != nil {
		for i, existingConn := range *conns {
			if existingConn != nil && existingConn.GetId() == conn.GetId() {
				(*conns)[i] = conn
				return
			}
		}
		*conns = append(*conns, conn)
	}
}

// Remove deletes an IConn from the slice by its UUID.
func (conns *IConns) Remove(id ConnId) {
	for i, conn := range *conns {
		if conn != nil && conn.GetId() == id {
			*conns = append((*conns)[:i], (*conns)[i+1:]...)
			return
		}
	}
}

// FindByConnId finds an IConn by its UUID.
func (conns IConns) FindByConnId(id ConnId) (IConn, bool) {
	return conns.Get(id)
}

// FindConnByAdapterName finds an IConn by the adapter name.
func (conns IConns) FindConnByAdapterName(name AdapterName) (IConn, bool) {
	for _, conn := range conns {
		if conn != nil && conn.GetAdapter() != nil && conn.GetAdapter().GetName() == name {
			return conn, true
		}
	}
	return nil, false
}

// FindByAdapterName finds an IAdapter by its name.
func (conns IConns) FindByAdapterName(name AdapterName) (IAdapter, bool) {
	for _, conn := range conns {
		if conn != nil && conn.GetAdapter() != nil && conn.GetAdapter().GetName() == name {
			return conn.GetAdapter(), true
		}
	}
	return nil, false
}

// IConnMap is a map of IConn with UUID as the key.
type IConnMap map[ConnId]IConn

// Get retrieves an IConn by its UUID.
func (m IConnMap) Get(id ConnId) (IConn, bool) {
	conn, exists := m[id]
	return conn, exists
}

// Set adds or updates an IConn in the map.
func (m IConnMap) Set(id ConnId, conn IConn) {
	if conn != nil {
		if conn.GetId().IsNil() || conn.GetAdapter() == nil {
			return
		}
		m[id] = conn
	}
}

// Remove deletes an IConn from the map by its UUID.
func (m IConnMap) Remove(id ConnId) {
	delete(m, id)
}

// ToArray converts the IConnMap to an array of IConns.
func (m IConnMap) ToArray() IConns {
	conns := IConns{}
	for _, conn := range m {
		if conn != nil {
			if conn.GetId().IsNil() || conn.GetAdapter() == nil {
				continue
			}
			conns = append(conns, conn)
		}
	}
	return conns
}

// ToAdapterArray converts the IConnMap to an array of IAdapters.
func (m IConnMap) ToAdapterArray() IAdapters {
	adapters := IAdapters{}
	for _, conn := range m {
		if conn != nil && conn.GetAdapter() != nil {
			adapters = append(adapters, conn.GetAdapter())
		}
	}
	return adapters
}

// ToAdapterMap converts the IConnMap to an IAdapterMap.
func (m IConnMap) ToAdapterMap() IAdapterMap {
	adapters := IAdapterMap{}
	for _, conn := range m {
		if conn != nil && conn.GetAdapter() != nil {
			adapters[conn.GetAdapter().GetName()] = conn.GetAdapter()
		}
	}
	return adapters
}

// FilterByAuthScope returns all connections that have the specified auth scope.
func (conns IConns) FilterByAuthScope(authScope AuthScope) IConns {
	var filtered IConns
	for _, conn := range conns {
		if conn.GetAuthScopes().Has(authScope) {
			filtered = append(filtered, conn)
		}
	}
	return filtered
}

// FilterByAnyAuthScope returns connections that match *any* of the given auth scopes.
func (conns IConns) FilterByAnyRole(authScopes AuthScopes) IConns {
	var filtered IConns
	for _, conn := range conns {
		for _, a := range conn.GetAuthScopes() {
			if authScopes.Has(a) {
				filtered = append(filtered, conn)
				break
			}
		}
	}
	return filtered
}

// GetTenantInfos returns a slice of ConnTenantInfo from all connections.
func (conns IConns) GetTenantInfos() ConnTenantInfos {
	var infos ConnTenantInfos
	for _, conn := range conns {
		info := conn.GetTenantInfo()
		if !info.TenantId.IsNil() {
			infos = append(infos, info)
		}
	}
	return infos
}

// FindByTenantId returns the first connection that matches the given tenant ID.
func (conns IConns) FindByTenantId(tenantId auuids.UUID) (IConn, bool) {
	for _, conn := range conns {
		if conn.GetTenantInfo().TenantId == tenantId {
			return conn, true
		}
	}
	return nil, false
}
