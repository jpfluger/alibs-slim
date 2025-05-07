package aconns

import (
	"github.com/jpfluger/alibs-slim/auuids"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTenantManager_AuthPipeline_Build(t *testing.T) {
	// Generate IDs
	tenantId := auuids.NewUUID()
	connId1 := NewConnId()
	connId2 := NewConnId()

	// Create auth connections
	auth1 := &Conn{
		Id: connId1,
		Adapter: &DummyAdapter{Adapter: Adapter{
			Name: "ldap-primary",
			Type: "ldap",
			Host: "ldap.example.com",
			Port: 389,
		}},
		Roles: ConnRoles{CONNROLE_AUTH},
		TenantInfo: ConnTenantInfo{
			TenantId: tenantId,
			Priority: 1,
			Region:   "us-west",
		},
	}
	auth1.SetAuthMethods(AuthMethods{AUTHMETHOD_PRIMARY, AUTHMETHOD_SSPR})

	auth2 := &Conn{
		Id: connId2,
		Adapter: &DummyAdapter{Adapter: Adapter{
			Name: "duo-mfa",
			Type: "duo",
			Host: "duo.example.com",
			Port: 443,
		}},
		Roles: ConnRoles{CONNROLE_AUTH},
		TenantInfo: ConnTenantInfo{
			TenantId: tenantId,
			Priority: 0, // Should be sorted before auth1
			Region:   "us-west",
		},
	}
	auth2.SetAuthMethods(AuthMethods{AUTHMETHOD_MFA})

	// Assemble TenantManager
	tm := &TenantManager{
		Auths: IConns{auth1, auth2},
	}
	tm.AuthFlows = tm.BuildAuthPipeline()

	// Check that all methods are included
	assert.Contains(t, tm.AuthFlows, AUTHMETHOD_PRIMARY)
	assert.Contains(t, tm.AuthFlows, AUTHMETHOD_MFA)
	assert.Contains(t, tm.AuthFlows, AUTHMETHOD_SSPR)

	// Validate adapter order for MFA
	mfaEntries := tm.AuthFlows[AUTHMETHOD_MFA]
	assert.Len(t, mfaEntries, 1)
	assert.Equal(t, connId2, mfaEntries[0].ConnId)

	// Validate adapter order for PRIMARY and SSPR
	primaryEntries := tm.AuthFlows[AUTHMETHOD_PRIMARY]
	ssprEntries := tm.AuthFlows[AUTHMETHOD_SSPR]
	assert.Len(t, primaryEntries, 1)
	assert.Len(t, ssprEntries, 1)
	assert.Equal(t, connId1, primaryEntries[0].ConnId)
	assert.Equal(t, connId1, ssprEntries[0].ConnId)

	// Confirm sorting by priority
	assert.Equal(t, 0, tm.AuthFlows[AUTHMETHOD_MFA][0].Priority)
	assert.Equal(t, 1, tm.AuthFlows[AUTHMETHOD_PRIMARY][0].Priority)
}
