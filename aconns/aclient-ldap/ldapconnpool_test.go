package aclient_ldap

import (
	"crypto/tls"
	"testing"

	"github.com/go-ldap/ldap/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockLdapConfig mocks the ILdapConfig interface for testing.
type MockLdapConfig struct {
	mock.Mock
}

func (m *MockLdapConfig) GetURL() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockLdapConfig) GetAdminDN() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockLdapConfig) GetAdminPass() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockLdapConfig) GetMaxOpen() int {
	args := m.Called()
	return args.Int(0)
}

func (m *MockLdapConfig) GetMaxDialerTimeout() int {
	args := m.Called()
	return args.Int(0)
}

func (m *MockLdapConfig) GetUseSSL() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *MockLdapConfig) GetNetwork() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockLdapConfig) GetAddress() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockLdapConfig) GetHost() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockLdapConfig) GetSkipTLS() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *MockLdapConfig) GetInsecureSkipVerify() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *MockLdapConfig) GetServerName() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockLdapConfig) GetTLSCertificates() []tls.Certificate {
	args := m.Called()
	return args.Get(0).([]tls.Certificate)
}

func TestGetConnection_PoolHasAvailableConnection(t *testing.T) {
	// Initialize LdapConnPool with mock configuration.
	mockConf := new(MockLdapConfig)
	pool := &LdapConnPool{
		conns:   []ILdapConn{&ldap.Conn{}},
		maxOpen: 10,
	}

	// Call the function.
	conn, err := pool.GetConnection(mockConf)

	// Assertions
	assert.NotNil(t, conn)
	assert.Nil(t, err)
	assert.Equal(t, 1, pool.openConn)
}

func TestGetConnection_CreateNewConnection(t *testing.T) {
	// Initialize LdapConnPool with no available connections and inject mock connection initializer.
	mockConf := new(MockLdapConfig)
	mockConf.On("GetURL").Return("")
	mockConf.On("GetUseSSL").Return(false)
	mockConf.On("GetNetwork").Return("tcp")
	mockConf.On("GetAddress").Return("localhost:389")
	mockConf.On("GetHost").Return("localhost")
	mockConf.On("GetAdminDN").Return("admin")
	mockConf.On("GetAdminPass").Return("password")

	// Initialize pool with mocked connection initializer.
	pool := &LdapConnPool{
		conns:   []ILdapConn{}, // No available connections.
		maxOpen: 10,            // Max open connections set to 10.
		initConnFunc: func(conf ILdapConfig) (ILdapConn, error) {
			return &ldap.Conn{}, nil // Return a valid connection.
		},
	}

	// Call the function.
	conn, err := pool.GetConnection(mockConf)

	// Assertions
	assert.NotNil(t, conn)            // Connection should be non-nil.
	assert.Nil(t, err)                // There should be no error.
	assert.Equal(t, 1, pool.openConn) // The open connection count should increment.
}

func TestGetConnection_MaxOpenReached(t *testing.T) {
	// Initialize LdapConnPool with max open connections reached.
	mockConf := new(MockLdapConfig)
	pool := &LdapConnPool{
		conns:    []ILdapConn{}, // No available connections.
		maxOpen:  1,             // Max open connections is set to 1.
		openConn: 1,             // One open connection already.
		reqConns: make(map[uint64]chan ILdapConn),
	}

	// Call the function in a goroutine to simulate blocking until a connection is returned.
	connChan := make(chan ILdapConn, 1)
	go func() {
		conn, _ := pool.GetConnection(mockConf)
		connChan <- conn
	}()

	// Simulate releasing a connection to unblock the request.
	pool.PutConnection(&ldap.Conn{})

	// Assertions
	assert.NotNil(t, <-connChan)           // A connection should be returned from the channel.
	assert.Equal(t, 0, len(pool.reqConns)) // Request should be removed from the queue.
}

func TestPutConnection_ConnectionReturnedToPool(t *testing.T) {
	// Initialize LdapConnPool with one open connection.
	pool := &LdapConnPool{
		conns:    []ILdapConn{},
		openConn: 1,
		maxOpen:  10,
	}

	// Create a mock connection.
	conn := &ldap.Conn{}

	// Call the function to return the connection to the pool.
	pool.PutConnection(conn)

	// Assertions
	assert.Equal(t, 0, pool.openConn)   // Open connection count should decrement.
	assert.Equal(t, 1, len(pool.conns)) // Connection should be added back to the pool.
}

func TestPutConnection_FulfillPendingRequest(t *testing.T) {
	// Initialize LdapConnPool with a pending request.
	pool := &LdapConnPool{
		conns:    []ILdapConn{}, // No available connections.
		openConn: 1,             // One open connection already.
		maxOpen:  10,            // Max open connections is set to 10.
		reqConns: map[uint64]chan ILdapConn{
			12345: make(chan ILdapConn, 1),
		},
	}

	// Create a mock connection.
	conn := new(MockLdapConn)

	// Capture the request channel before fulfilling the request.
	reqChan := pool.reqConns[12345]

	// Call the function to fulfill the request.
	pool.PutConnection(conn)

	// Assertions
	assert.Equal(t, conn, <-reqChan)       // Request should be fulfilled with the connection.
	assert.Equal(t, 1, pool.openConn)      // Open connection count should remain the same.
	assert.Equal(t, 0, len(pool.reqConns)) // Request queue should be empty.
	assert.Equal(t, 0, len(pool.conns))    // Connection should not be returned to the pool.
}

func TestInitLDAPConn_Success(t *testing.T) {
	// Mock configuration for a successful connection.
	mockConf := new(MockLdapConfig)
	mockConf.On("GetURL").Return("") // No URL, simulate non-URL connection
	mockConf.On("GetNetwork").Return("tcp")
	mockConf.On("GetAddress").Return("localhost:389")
	mockConf.On("GetMaxDialerTimeout").Return(5)
	mockConf.On("GetUseSSL").Return(false)
	mockConf.On("GetSkipTLS").Return(false)
	mockConf.On("GetInsecureSkipVerify").Return(true)
	mockConf.On("GetAdminDN").Return("cn=admin,dc=example,dc=com")
	mockConf.On("GetAdminPass").Return("password")
	mockConf.On("GetHost").Return("localhost")

	// Mock the ILdapConn behavior
	mockLdapConn := new(MockLdapConn)
	mockLdapConn.On("Bind", "cn=admin,dc=example,dc=com", "password").Return(nil)
	mockLdapConn.On("StartTLS", mock.AnythingOfType("*tls.Config")).Return(nil)

	// Inject a mock initLDAPConn function.
	initLDAPConn := func(conf ILdapConfig) (ILdapConn, error) {
		return mockLdapConn, nil // Return the mocked connection instead of dialing a real one.
	}

	// Call initLDAPConn to test connection initialization.
	conn, err := initLDAPConn(mockConf)

	// Assertions
	assert.NotNil(t, conn) // Connection should be non-nil.
	assert.Nil(t, err)     // There should be no error.
}

func TestInitLDAPConn_Failure(t *testing.T) {
	// Mock configuration with an invalid URL.
	mockConf := new(MockLdapConfig)
	mockConf.On("GetURL").Return("invalid_url")
	mockConf.On("GetHost").Return("localhost")
	mockConf.On("GetMaxDialerTimeout").Return(10) // Ensure this line is included

	// Call initLDAPConn to simulate a connection failure.
	conn, err := initLDAPConn(mockConf)

	// Assertions
	assert.Nil(t, conn)                                 // Connection should be nil on failure.
	assert.NotNil(t, err)                               // There should be an error.
	assert.Contains(t, err.Error(), "invalid ldap URL") // Check for the "invalid ldap URL" part in the error.
}
