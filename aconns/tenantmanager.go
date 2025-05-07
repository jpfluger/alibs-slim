package aconns

// TenantManager is a multi-tenant orchestration layer that organizes connections by role.
// It is typically derived from a flat Manager instance and allows structured access to:
// - Authentication connections (Auths)
// - Master infrastructure (Masters)
// - Tenant-specific scoped connections (Tenants)
//
// Additionally, it builds and stores AuthFlows, a method-to-adapter pipeline
// derived from the Auths list, which dictates how authentication methods are resolved.
type TenantManager struct {
	// Auths are authentication-related connections (e.g., LDAP, OIDC, SSO).
	// These are scanned to construct the AuthFlows used during login or auth workflows.
	Auths IConns `json:"auths,omitempty"`

	// Masters are primary infrastructure connections such as default databases, cache layers, etc.
	Masters IConns `json:"masters,omitempty"`

	// Tenants are tenant-specific connections (e.g., one DB per client, region-sharded services).
	Tenants IConns `json:"tenants,omitempty"`

	// AuthFlows is the ordered authentication pipeline derived from Auths.
	// Each method (primary, MFA, SSPR) maps to an ordered list of adapters.
	AuthFlows AuthPipeline `json:"authFlows,omitempty"`
}

// BuildAuthPipeline scans the Auths list and assembles a map of AuthMethod → ordered AuthAdapterEntries.
// Each entry includes the ConnId, Adapter instance, and priority (from ConnTenantInfo).
// The resulting pipeline allows controlled and deterministic auth evaluation at runtime.
func (tm *TenantManager) BuildAuthPipeline() AuthPipeline {
	pipeline := make(AuthPipeline)

	for _, conn := range tm.Auths {
		methods := conn.GetAuthMethods()
		for _, method := range methods {
			entry := &AuthAdapterEntry{
				ConnId:   conn.GetId(),
				Adapter:  conn.GetAdapter(),
				Priority: conn.GetTenantInfo().Priority,
			}
			pipeline[method] = append(pipeline[method], entry)
		}
	}

	// Sort each method slice by Priority (ascending)
	for _, entries := range pipeline {
		entries.SortByPriority()
	}

	return pipeline
}
