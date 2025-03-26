package aclient_ldap

import (
	"crypto/tls"
	"fmt"
	"math/rand"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/go-ldap/ldap/v3"
)

// ILdapConfig defines the interface for obtaining LDAP connection configuration parameters.
type ILdapConfig interface {
	GetURL() string                        // Returns the LDAP server URL.
	GetAdminDN() string                    // Returns the admin DN for authentication.
	GetAdminPass() string                  // Returns the password for the admin DN.
	GetMaxOpen() int                       // Returns the maximum allowed open connections.
	GetMaxDialerTimeout() int              // Returns the dialer timeout in seconds.
	GetUseSSL() bool                       // Indicates if SSL should be used.
	GetNetwork() string                    // Returns the network type, usually "tcp".
	GetAddress() string                    // Returns the LDAP server address.
	GetHost() string                       // Returns the LDAP server host.
	GetSkipTLS() bool                      // Indicates if TLS upgrade should be skipped.
	GetInsecureSkipVerify() bool           // Indicates whether to skip SSL certificate verification.
	GetTLSCertificates() []tls.Certificate // Returns the TLS certificates for SSL if applicable.
}

// LdapConnPool manages a pool of reusable LDAP connections.
type LdapConnPool struct {
	mu           sync.Mutex                                // Protects access to the connection pool.
	conns        []ILdapConn                               // Slice holding available LDAP connections.
	reqConns     map[uint64]chan ILdapConn                 // Map for holding requests waiting for a connection.
	openConn     int                                       // Current number of open connections.
	maxOpen      int                                       // Maximum allowed open connections.
	DsName       string                                    // Identifier for the data source.
	initConnFunc func(conf ILdapConfig) (ILdapConn, error) // Injected function for initializing LDAP connection.
}

type ILdapConnPool interface {
	GetConnection(conf ILdapConfig) (ILdapConn, error)
	PutConnection(conn ILdapConn)
	CloseAllConnections() error
}

// Default initializer function that can be replaced for testing.
func defaultInitConnFunc(conf ILdapConfig) (ILdapConn, error) {
	return initLDAPConn(conf)
}

// InitLdapConnPool initializes a new LdapConnPool with the default connection initializer.
func InitLdapConnPool(maxOpen int) *LdapConnPool {
	return &LdapConnPool{
		conns:        []ILdapConn{},
		reqConns:     make(map[uint64]chan ILdapConn),
		maxOpen:      maxOpen,
		initConnFunc: defaultInitConnFunc, // Default connection initializer.
	}
}

// GetConnection retrieves a connection from the pool or creates a new one if necessary.
func (lcp *LdapConnPool) GetConnection(conf ILdapConfig) (ILdapConn, error) {
	lcp.mu.Lock()
	connNum := len(lcp.conns)
	if connNum > 0 {
		// Retrieve a connection from the pool.
		conn := lcp.conns[0]
		copy(lcp.conns, lcp.conns[1:])
		lcp.conns = lcp.conns[:connNum-1]
		lcp.openConn++
		lcp.mu.Unlock()

		// Check if the connection is still valid.
		if conn.IsClosing() {
			if lcp.initConnFunc == nil {
				return initLDAPConn(conf)
			} else {
				return lcp.initConnFunc(conf) // Use injected function.
			}
		}
		return conn, nil
	}

	// Check if maximum open connections have been reached.
	if lcp.maxOpen != 0 && lcp.openConn >= lcp.maxOpen {
		req := make(chan ILdapConn, 1)
		reqKey := lcp.nextRequestKeyLocked()
		lcp.reqConns[reqKey] = req
		lcp.mu.Unlock()
		return <-req, nil // Wait for an available connection.
	}

	// Create a new connection if none are available.
	lcp.openConn++
	lcp.mu.Unlock()
	if lcp.initConnFunc == nil {
		return initLDAPConn(conf)
	}
	return lcp.initConnFunc(conf) // Use injected function.
}

// PutConnection returns a connection back to the pool or sends it to a waiting requester.
func (lcp *LdapConnPool) PutConnection(conn ILdapConn) {
	lcp.mu.Lock()
	defer lcp.mu.Unlock()

	// If there is a request in the queue, fulfill it.
	if len(lcp.reqConns) > 0 {
		var req chan ILdapConn
		var reqKey uint64
		for reqKey, req = range lcp.reqConns {
			break
		}
		delete(lcp.reqConns, reqKey)
		req <- conn
		return
	}

	// Add the connection back to the pool if it is still open.
	lcp.openConn--
	if !conn.IsClosing() {
		lcp.conns = append(lcp.conns, conn)
	}
}

// CloseAllConnections closes all active connections in the pool and resets the pool state.
func (lcp *LdapConnPool) CloseAllConnections() error {
	lcp.mu.Lock()
	defer lcp.mu.Unlock()

	var rerr error
	// Close all connections in the pool
	for _, conn := range lcp.conns {
		err := conn.Close()
		if err != nil && rerr == nil {
			rerr = err
		}
	}

	// Reset the pool's internal state.
	lcp.conns = nil
	lcp.openConn = 0
	lcp.reqConns = make(map[uint64]chan ILdapConn)

	return rerr
}

// nextRequestKeyLocked generates a unique key for connection requests.
func (lcp *LdapConnPool) nextRequestKeyLocked() uint64 {
	for {
		reqKey := rand.Uint64()
		if _, ok := lcp.reqConns[reqKey]; !ok {
			return reqKey
		}
	}
}

// ILdapConn defines an interface for the ldap.Conn methods that you want to mock.
type ILdapConn interface {
	Search(searchRequest *ldap.SearchRequest) (*ldap.SearchResult, error)
	Bind(username, password string) error
	StartTLS(config *tls.Config) error
	IsClosing() bool
	Close() error
}

// LdapConnWrapper wraps ldap.Conn and implements the ILdapConn interface.
type LdapConnWrapper struct {
	Conn *ldap.Conn
}

func (w *LdapConnWrapper) Bind(username, password string) error {
	return w.Conn.Bind(username, password)
}

func (w *LdapConnWrapper) StartTLS(config *tls.Config) error {
	return w.Conn.StartTLS(config)
}

func (w *LdapConnWrapper) Close() error {
	return w.Conn.Close()
}

func (w *LdapConnWrapper) IsClosing() bool {
	return w.Conn.IsClosing()
}

func initLDAPConn(conf ILdapConfig) (ILdapConn, error) { //*ldap.Conn
	var l *ldap.Conn
	var err error

	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic caught: %v", r)
		}
	}()

	// Create a net.Dialer with the configured timeout
	dialer := &net.Dialer{
		Timeout: time.Duration(conf.GetMaxDialerTimeout()) * time.Second,
	}

	// Use URL-based connection if provided
	if conf.GetURL() != "" {
		// Validate the URL format before dialing
		if !strings.HasPrefix(conf.GetURL(), "ldap://") && !strings.HasPrefix(conf.GetURL(), "ldaps://") {
			return nil, fmt.Errorf("invalid ldap URL: %s", conf.GetURL())
		}
		// Attempt to dial the LDAP server with the provided URL
		l, err = ldap.DialURL(conf.GetURL(), ldap.DialWithDialer(dialer))
		if err != nil {
			return nil, fmt.Errorf("could not dial ldap with DialURL where host=%s; %v", conf.GetHost(), err)
		}
	} else if !conf.GetUseSSL() {
		// Non-SSL connection using the net.Dialer (host:port only)
		l, err = ldap.DialURL(fmt.Sprintf("ldap://%s", conf.GetAddress()), ldap.DialWithDialer(dialer))
		if err != nil {
			return nil, fmt.Errorf("could not dial ldap with DialURL where host=%s; %v", conf.GetHost(), err)
		}

		// Upgrade to TLS if required
		if !conf.GetSkipTLS() {
			err = l.StartTLS(&tls.Config{InsecureSkipVerify: conf.GetInsecureSkipVerify()})
			if err != nil {
				return nil, fmt.Errorf("could not start TLS for ldap connection where host=%s; %v", conf.GetHost(), err)
			}
		}
	} else {
		// SSL connection using DialURL and net.Dialer
		ldapAddress := conf.GetAddress()
		if !strings.HasPrefix(ldapAddress, "ldaps://") {
			ldapAddress = fmt.Sprintf("ldaps://%s", ldapAddress)
		}
		l, err = ldap.DialURL(ldapAddress, ldap.DialWithDialer(dialer), ldap.DialWithTLSConfig(&tls.Config{
			InsecureSkipVerify: conf.GetInsecureSkipVerify(),
			ServerName:         conf.GetHost(),
			Certificates:       conf.GetTLSCertificates(),
		}))
		if err != nil {
			return nil, fmt.Errorf("could not dial ldap connection with SSL where host=%s; %v", conf.GetHost(), err)
		}
	}

	// Bind the connection using the provided admin credentials
	err = l.Bind(conf.GetAdminDN(), conf.GetAdminPass())
	if err != nil {
		return nil, err
	}

	//return &LdapConnWrapper{Conn: l}, nil
	return l, nil
}
