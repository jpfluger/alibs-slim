package aclient_ldap

import (
	"crypto/tls"
	"testing"

	"github.com/go-ldap/ldap/v3"
	"github.com/jpfluger/alibs-slim/aconns"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

//// MockLdapConnPool mocks LdapConnPool for testing purposes.
//type MockLdapConnPool struct {
//	mock.Mock
//}
//
//func (m *MockLdapConnPool) CloseAllConnections() error {
//	return nil
//}
//
//func (m *MockLdapConnPool) GetConnection(conf ILdapConfig) (ILdapConn, error) {
//	args := m.Called(conf)
//	return args.Get(0).(*ldap.Conn), args.Error(1)
//}
//
//func (m *MockLdapConnPool) PutConnection(conn ILdapConn) {
//	m.Called(conn)
//}
//
//// MockLdapConn is a mock for the ILdapConn interface.
//type MockLdapConn struct {
//	mock.Mock
//}
//
//func (m *MockLdapConn) StartTLS(config *tls.Config) error {
//	args := m.Called(config)
//	return args.Error(0)
//}
//
//func (m *MockLdapConn) Bind(username, password string) error {
//	args := m.Called(username, password)
//	return args.Error(0)
//}
//
//func (m *MockLdapConn) Search(searchRequest *ldap.SearchRequest) (*ldap.SearchResult, error) {
//	args := m.Called(searchRequest)
//	return args.Get(0).(*ldap.SearchResult), args.Error(1)
//}
//
//func (m *MockLdapConn) Close() error {
//	args := m.Called()
//	return args.Error(0)
//}
//
//func (m *MockLdapConn) IsClosing() bool {
//	args := m.Called()
//	return args.Bool(0)
//}

// MockLdapConnPool mocks LdapConnPool for testing purposes.
type MockLdapConnPool struct {
	mock.Mock
}

func (m *MockLdapConnPool) CloseAllConnections() error {
	return nil
}

// GetConnection now returns ILdapConn (the mock we created).
func (m *MockLdapConnPool) GetConnection(conf ILdapConfig) (ILdapConn, error) {
	args := m.Called(conf)
	return args.Get(0).(ILdapConn), args.Error(1)
}

func (m *MockLdapConnPool) PutConnection(conn ILdapConn) {
	m.Called(conn)
}

// MockLdapConn is a mock for the ILdapConn interface.
type MockLdapConn struct {
	mock.Mock
}

func (m *MockLdapConn) StartTLS(config *tls.Config) error {
	args := m.Called(config)
	return args.Error(0)
}

func (m *MockLdapConn) Bind(username, password string) error {
	args := m.Called(username, password)
	return args.Error(0)
}

func (m *MockLdapConn) Search(searchRequest *ldap.SearchRequest) (*ldap.SearchResult, error) {
	args := m.Called(searchRequest)
	return args.Get(0).(*ldap.SearchResult), args.Error(1)
}

func (m *MockLdapConn) Close() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockLdapConn) IsClosing() bool {
	args := m.Called()
	return args.Bool(0)
}

// Unit Test: Validate Function
//func TestAClientLDAP_Validate(t *testing.T) {
//	client := &AClientLDAP{
//		Base:   "dc=example,dc=com",
//		BindDN: "cn=admin,dc=example,dc=com",
//		UseSSL: true,
//		ADBAdapterBase: aconns.ADBAdapterBase{
//			Adapter: aconns.Adapter{
//				Type: ADAPTERTYPE_LDAP,
//				Name: "ldap",
//				Host: "ldap.example.com",
//				Port: 0,
//			},
//		},
//	}
//
//	err := client.Validate()
//	assert.Nil(t, err)                               // Validate should pass.
//	assert.Equal(t, 636, client.Port)                // Port should default to 636 when SSL is enabled.
//	assert.Equal(t, "ldap.example.com", client.Host) // ServerName should not be changed.
//}

// Unit Test: InitPool
func TestAClientLDAP_InitPool(t *testing.T) {
	client := &AClientLDAP{
		Base:   "dc=example,dc=com",
		BindDN: "cn=admin,dc=example,dc=com",
	}

	mockPool := new(MockLdapConnPool)
	client.ldapPool = mockPool

	mockLdapConn := &ldap.Conn{}
	mockLdapConfig := new(MockLdapConfig)

	mockPool.On("GetConnection", mockLdapConfig).Return(mockLdapConn, nil)
	mockPool.On("PutConnection", mockLdapConn).Return(nil)

	err := client.InitPool()
	assert.Nil(t, err) // Pool should be initialized successfully.
}

// Unit Test: Test Connection
func TestAClientLDAP_TestConnection(t *testing.T) {
	client := &AClientLDAP{
		Base:              "dc=example,dc=com",
		BindDN:            "cn=admin,dc=example,dc=com",
		ConnectionTimeout: 5,
		ADBAdapterBase: aconns.ADBAdapterBase{
			Password: "adminpass", // needed for the admin bind
			Adapter: aconns.Adapter{
				Type: ADAPTERTYPE_LDAP,
				Name: "ldap",
				Host: "ldap.example.com",
				Port: 0,
			},
		},
	}

	// ---- Mock the pool ------------------------------------------------
	mockPool := new(MockLdapConnPool)
	client.ldapPool = mockPool

	// ---- Mock the connection -----------------------------------------
	mockConn := new(MockLdapConn)

	// Pool will return our mock connection
	mockPool.On("GetConnection", mock.AnythingOfType("*aclient_ldap.clientLdapConfig")).
		Return(mockConn, nil)
	mockPool.On("PutConnection", mockConn).Return()

	// ---- Admin bind (BindDN is set) -----------------------------------
	mockConn.On("Bind", "cn=admin,dc=example,dc=com", "adminpass").
		Return(nil)

	// ---- Base-object search -------------------------------------------
	mockConn.On("Search", mock.MatchedBy(func(sr *ldap.SearchRequest) bool {
		return sr.BaseDN == "dc=example,dc=com" &&
			sr.Scope == ldap.ScopeBaseObject &&
			sr.Filter == "(objectClass=*)" &&
			len(sr.Attributes) == 1 && sr.Attributes[0] == "dn"
	})).Return(&ldap.SearchResult{}, nil)

	// ---- Run the method under test ------------------------------------
	success, status, err := client.Test()

	// ---- Assertions ---------------------------------------------------
	assert.True(t, success)
	assert.Equal(t, aconns.TESTSTATUS_INITIALIZED_SUCCESSFUL, status)
	assert.Nil(t, err)

	// Verify that every mocked call was made
	mockPool.AssertExpectations(t)
	mockConn.AssertExpectations(t)
}

// Unit Test: OpenConnection
func TestAClientLDAP_OpenConnection(t *testing.T) {
	client := &AClientLDAP{
		Base:   "dc=example,dc=com",
		BindDN: "cn=admin,dc=example,dc=com",
		ADBAdapterBase: aconns.ADBAdapterBase{
			Adapter: aconns.Adapter{
				Host: "ldap.example.com",
				Port: 0,
			},
		},
	}

	mockPool := new(MockLdapConnPool)
	client.ldapPool = mockPool

	mockConn := new(MockLdapConn)

	mockPool.On("GetConnection", mock.AnythingOfType("*aclient_ldap.clientLdapConfig")).Return(mockConn, nil)
	mockPool.On("PutConnection", mockConn).Return(nil)

	err := client.OpenConnection()
	assert.Nil(t, err) // Open connection should succeed.

	mockPool.AssertExpectations(t)
}

// Unit Test: CloseConnection
func TestAClientLDAP_CloseConnection(t *testing.T) {
	client := &AClientLDAP{
		Base:   "dc=example,dc=com",
		BindDN: "cn=admin,dc=example,dc=com",
		ADBAdapterBase: aconns.ADBAdapterBase{
			Adapter: aconns.Adapter{
				Host: "ldap.example.com",
				Port: 0,
			},
		},
	}

	mockPool := new(MockLdapConnPool)
	client.ldapPool = mockPool

	mockLdapConn := &ldap.Conn{}

	mockPool.On("PutConnection", mockLdapConn).Return(nil)

	err := client.CloseConnection()
	assert.Nil(t, err) // Close connection should succeed.
}

func TestGetConnection_ConnectionIsClosing(t *testing.T) {
	// Initialize LdapConnPool with one connection that is closing.
	mockConf := new(MockLdapConfig)
	mockConf.On("GetMaxDialerTimeout").Return(5)

	mockLdapConn := new(MockLdapConn)
	mockLdapConn.On("IsClosing").Return(true) // Simulate that the connection is closing.

	pool := &LdapConnPool{
		conns:    []ILdapConn{mockLdapConn},
		openConn: 0,
		maxOpen:  10,
	}

	// Mock initLDAPConn to return a new connection when the current one is closing.
	mockLdapConnNew := new(MockLdapConn)
	mockLdapConnNew.On("Bind", "cn=admin,dc=example,dc=com", "password").Return(nil)
	pool.initConnFunc = func(conf ILdapConfig) (ILdapConn, error) {
		return mockLdapConnNew, nil
	}

	// Call GetConnection.
	conn, err := pool.GetConnection(mockConf)

	// Assertions
	assert.NotNil(t, conn)                    // A new connection should be returned.
	assert.Nil(t, err)                        // No error should occur.
	mockLdapConn.AssertCalled(t, "IsClosing") // Ensure IsClosing was called on the old connection.
}

//func TestAClientLDAP_Authenticate(t *testing.T) {
//	client := &AClientLDAP{
//		Base:   "dc=example,dc=com",
//		BindDN: "cn=admin,dc=example,dc=com",
//		ADBAdapterBase: aconns.ADBAdapterBase{
//			Adapter: aconns.Adapter{
//				Type: ADAPTERTYPE_LDAP,
//				Name: "ldap",
//				Host: "ldap.example.com",
//				Port: 0,
//			},
//		},
//	}
//
//	mockPool := new(MockLdapConnPool)
//	client.ldapPool = mockPool
//
//	mockLdapConn := new(MockLdapConn)
//
//	// Mock the connection pool behavior, returning a *LdapConnWrapper that wraps the *ldap.Conn
//	mockPool.On("GetConnection", mock.AnythingOfType("*aclient_ldap.clientLdapConfig")).Return(mockLdapConn, nil)
//	//mockPool.On("PutConnection", mock.AnythingOfType("*ldap.Conn")).Return(nil)
//	mockPool.On("PutConnection", mockLdapConn).Return(nil)
//
//	// Mock the Bind and Search methods for the connection
//	mockLdapConn.On("Bind", "username", "password").Return(nil)
//	mockLdapConn.On("Search", mock.AnythingOfType("*ldap.SearchRequest")).Return(&ldap.SearchResult{
//		Entries: []*ldap.Entry{
//			{
//				DN: "cn=username,dc=example,dc=com",
//				Attributes: []*ldap.EntryAttribute{
//					{Name: "cn", Values: []string{"username"}},
//				},
//			},
//		},
//	}, nil)
//	mockLdapConn.On("Close").Return(nil)
//
//	// Call the method under test
//	success, attrs, err := client.Authenticate("username", "password")
//
//	// Assertions
//	assert.True(t, success)                  // Authentication should succeed.
//	assert.Nil(t, err)                       // No error should be returned.
//	assert.NotNil(t, attrs)                  // Attributes should not be nil.
//	assert.Equal(t, "username", attrs["cn"]) // Verify the attributes.
//}
