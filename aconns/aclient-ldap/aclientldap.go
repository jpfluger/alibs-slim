package aclient_ldap

import (
	"context"
	"crypto/tls"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-ldap/ldap/v3"
	"github.com/jpfluger/alibs-slim/aconns"
	"github.com/jpfluger/alibs-slim/asessions"
	"github.com/jpfluger/alibs-slim/autils"
)

const (
	ADAPTERTYPE_LDAP        = aconns.AdapterType("ldap")
	LDAP_CONNECTION_TIMEOUT = 5
	LDAP_DEFAULT_PORT       = 389
	LDAP_DEFAULT_PORT_SSL   = 636
)

type AClientLDAP struct {
	aconns.ADBAdapterBase

	Attributes autils.StringsArray `json:"attributes,omitempty"`
	Base       string              `json:"base,omitempty"`
	BindDN     string              `json:"bindDN,omitempty"`

	GroupFilter string `json:"groupFilter,omitempty"`

	//ServerName string `json:"serverName,omitempty"`

	UserFilter         string                        `json:"userFilter,omitempty"`
	InsecureSkipVerify bool                          `json:"insecureSkipVerify,omitempty"`
	UseSSL             bool                          `json:"useSSL,omitempty"`
	SkipTLS            bool                          `json:"skipTLS,omitempty"`
	ClientCertificates []tls.Certificate             `json:"clientCertificates,omitempty"`
	PermGroups         map[string]asessions.RoleName `json:"permGroups,omitempty"`

	ConnectionTimeout int `json:"connectionTimeout,omitempty"`

	ldapPool ILdapConnPool // LDAP connection pool instance

	mu sync.RWMutex
}

// validate checks if the AClientLDAP object is valid, including essential configurations like
// the server name, bind DN, and port. It sets defaults if certain values are not provided.
func (cn *AClientLDAP) validate() error {
	if err := cn.ADBAdapterBase.Validate(); err != nil {
		if err != aconns.ErrDatabaseIsEmpty {
			return err
		}
	}

	//cn.ServerName = strings.TrimSpace(cn.ServerName)
	//if cn.ServerName == "" {
	//	cn.ServerName = cn.Host
	//}

	if cn.ConnectionTimeout <= 0 {
		cn.ConnectionTimeout = LDAP_CONNECTION_TIMEOUT
	}

	cn.BindDN = strings.TrimSpace(cn.BindDN)

	if cn.Port <= 0 {
		if cn.UseSSL {
			cn.Port = LDAP_DEFAULT_PORT_SSL
		} else {
			cn.Port = LDAP_DEFAULT_PORT
		}
	}

	return nil // Note: Moved initPool out of validate to Test/Open to keep validate config-only
}

// Validate locks the AClientLDAP instance, validates its configuration,
// and ensures that essential fields are correctly set before use.
func (cn *AClientLDAP) Validate() error {
	cn.mu.Lock()
	defer cn.mu.Unlock()
	return cn.validate()
}

// InitPool initializes the LDAP connection pool for the client if it hasn't been initialized yet.
// It ensures that the pool can be used to manage multiple LDAP connections efficiently.
func (cn *AClientLDAP) InitPool() error {
	// Initialize the connection pool using the current configuration.
	cn.mu.Lock()
	defer cn.mu.Unlock()
	return cn.initPool()
}

func (cn *AClientLDAP) initPool() error {
	if cn.ldapPool != nil {
		return nil // Pool already initialized.
	}

	ldapConfig := &clientLdapConfig{
		client: cn,
	}

	cn.ldapPool = InitLdapConnPool(10)

	// Open initial connections if needed.
	conn, err := cn.ldapPool.GetConnection(ldapConfig)
	if err != nil {
		return err
	}
	defer cn.ldapPool.PutConnection(conn)

	return nil
}

// Test validates the AClientLDAP object and checks if a connection can be successfully
// established with the LDAP server, returning a status indicating success or failure.
func (cn *AClientLDAP) Test() (bool, aconns.TestStatus, error) {
	cn.mu.Lock()
	defer cn.mu.Unlock()
	return cn.test()
}

func (cn *AClientLDAP) test() (bool, aconns.TestStatus, error) {
	if err := cn.validate(); err != nil {
		cn.UpdateHealth(aconns.HEALTHSTATUS_VALIDATE_FAILED)
		return false, aconns.TESTSTATUS_FAILED, err
	}

	if cn.ldapPool == nil {
		if err := cn.initPool(); err != nil {
			cn.UpdateHealth(aconns.HEALTHSTATUS_OPEN_FAILED)
			return false, aconns.TESTSTATUS_FAILED, err
		}
	}

	conn, err := cn.ldapPool.GetConnection(&clientLdapConfig{client: cn})
	if err != nil {
		cn.UpdateHealth(aconns.HEALTHSTATUS_OPEN_FAILED)
		return false, aconns.TESTSTATUS_FAILED, err
	}
	defer cn.ldapPool.PutConnection(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cn.ConnectionTimeout)*time.Second)
	defer cancel()
	if err = cn.testConnectionWithCtx(ctx, conn); err != nil {
		status := aconns.HEALTHSTATUS_PING_FAILED
		if context.DeadlineExceeded == err {
			status = aconns.HEALTHSTATUS_TIMEOUT
		} else if strings.Contains(err.Error(), "network") || strings.Contains(err.Error(), "connection refused") {
			status = aconns.HEALTHSTATUS_NETWORK_ERROR
		} else if strings.Contains(err.Error(), "invalid credentials") || strings.Contains(err.Error(), "bind failed") {
			status = aconns.HEALTHSTATUS_AUTH_FAILED
		}
		cn.UpdateHealth(status)
		//alog.LOGGER(alog.LOGGER_APP).Warn().Err(err).Msg("LDAP test failed")
		return false, aconns.TESTSTATUS_FAILED, err
	}

	cn.UpdateHealth(aconns.HEALTHSTATUS_HEALTHY)
	return true, aconns.TESTSTATUS_INITIALIZED_SUCCESSFUL, nil
}

// OpenConnection opens a connection to the LDAP server by initializing the connection pool
// (if not already initialized) and fetching a connection from it.
func (cn *AClientLDAP) OpenConnection() error {
	cn.mu.Lock()
	defer cn.mu.Unlock()
	return cn.openConnection()
}

// openConnection ensures that the connection pool is initialized and retrieves a connection
// from the pool. This function doesn't manage the connection directly but through the pool.
func (cn *AClientLDAP) openConnection() error {
	// Ensure the pool is initialized before attempting to open a connection.
	if cn.ldapPool == nil {
		if err := cn.initPool(); err != nil {
			return fmt.Errorf("failed to initialize LDAP pool: %v", err)
		}
	}

	// Test a connection to verify
	conn, err := cn.ldapPool.GetConnection(&clientLdapConfig{client: cn})
	if err != nil {
		return err
	}
	defer cn.ldapPool.PutConnection(conn)

	return nil
}

// CloseConnection closes all connections in the LDAP pool and resets the pool.
func (cn *AClientLDAP) CloseConnection() error {
	cn.mu.Lock()
	defer cn.mu.Unlock()

	if cn.ldapPool != nil {
		cn.ldapPool.CloseAllConnections()
		cn.ldapPool = nil
		cn.UpdateHealth(aconns.HEALTHSTATUS_CLOSED)
	}
	return nil
}

// Refresh refreshes the LDAP connection pool by closing the existing one (if any) and initializing a new one.
func (cn *AClientLDAP) Refresh() error {
	cn.mu.Lock()
	defer cn.mu.Unlock()
	if cn.ldapPool != nil {
		cn.ldapPool.CloseAllConnections()
		cn.ldapPool = nil
	}
	return cn.initPool()
}

// GetAddress returns the address of the LDAP server.
func (cn *AClientLDAP) GetAddress() string {
	cn.mu.RLock()
	defer cn.mu.RUnlock()
	return cn.getAddress()
}

// getAddress returns the address of the LDAP server.
func (cn *AClientLDAP) getAddress() string {
	port := cn.GetPort()
	return fmt.Sprintf("%s:%s", cn.GetHost(), strconv.Itoa(port))
}

// GetConnectionTimeout returns the connection timeout.
func (cn *AClientLDAP) GetConnectionTimeout() int {
	cn.mu.RLock()
	defer cn.mu.RUnlock()
	return cn.ConnectionTimeout
}

// testConnectionWithCtx tests the connection to the LDAP server using a provided context.
// Since the library does not provide a native SearchWithContext, this uses a goroutine to simulate timeout behavior.
func (cn *AClientLDAP) testConnectionWithCtx(ctx context.Context, conn ILdapConn) error {
	if conn == nil {
		return fmt.Errorf("no ldap conn")
	}

	done := make(chan error, 1)
	go func() {
		// Real test: Bind admin (if set) and search base
		if cn.BindDN != "" && cn.Password != "" {
			if err := conn.Bind(cn.BindDN, cn.Password); err != nil {
				done <- fmt.Errorf("admin bind failed: %v", err)
				return
			}
		}
		searchReq := ldap.NewSearchRequest(cn.Base, ldap.ScopeBaseObject, ldap.NeverDerefAliases, 0, 0, false, "(objectClass=*)", []string{"dn"}, nil)
		if _, err := conn.Search(searchReq); err != nil {
			done <- fmt.Errorf("base search failed: %v", err)
			return
		}
		done <- nil
	}()

	select {
	case err := <-done:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}

// DB returns an LDAP connection from the pool. It checks health first and refreshes if stale or unhealthy.
func (cn *AClientLDAP) DB() (ILdapConn, error) {
	cn.mu.RLock()
	if cn.IsHealthy() && !cn.GetHealth().IsStale(5*time.Minute) {
		defer cn.mu.RUnlock()
		return cn.ldapPool.GetConnection(&clientLdapConfig{client: cn})
	}
	cn.mu.RUnlock()

	// Upgrade to write lock for refresh
	cn.mu.Lock()
	defer cn.mu.Unlock()
	if _, _, err := cn.test(); err != nil { // Refresh and test
		return nil, err
	}
	return cn.ldapPool.GetConnection(&clientLdapConfig{client: cn})
}

// GetGroupsForUser retrieves the groups for a given user from the LDAP directory.
// It uses the provided connection to perform the search.
func (cn *AClientLDAP) GetGroupsForUser(db ILdapConn, username string) ([]string, error) {
	cn.mu.RLock()
	defer cn.mu.RUnlock()

	ldapSearchUser := ldap.EscapeFilter(username)
	myBase := cn.Base
	myGroupFilter := cn.GroupFilter

	searchRequest := ldap.NewSearchRequest(
		myBase,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf(myGroupFilter, ldapSearchUser),
		[]string{"cn"},
		nil,
	)

	sr, err := db.Search(searchRequest)
	if err != nil {
		return nil, err
	}

	groups := []string{}
	for _, entry := range sr.Entries {
		groups = append(groups, entry.GetAttributeValue("cn"))
	}

	return groups, nil
}

// AuthenticateWithGroups verifies the user's credentials and retrieves the user's group information
// from the LDAP directory. It binds with the admin account (if needed) and ensures that
// the connection is returned to the pool after use.
func (cn *AClientLDAP) AuthenticateWithGroups(username, password string) (bool, *LDAPUserResult, error) {
	var db ILdapConn
	var err error

	// Ensure the connection is returned to the pool after usage
	defer func() {
		if db != nil {
			cn.ldapPool.PutConnection(db)
		}
	}()

	// Get a connection from the pool (using DB() which handles health)
	db, err = cn.DB()
	if err != nil {
		return false, nil, err
	}

	// Capture shared values
	cn.mu.RLock()
	myBindDN := cn.BindDN
	myPassword := cn.GetPassword()
	myBase := cn.Base
	myAttributes := cn.Attributes
	myUserFilter := cn.UserFilter
	cn.mu.RUnlock()

	// Bind with the admin DN if necessary
	if myBindDN != "" && myPassword != "" {
		err = db.Bind(myBindDN, myPassword)
		if err != nil {
			return false, nil, err
		}
	}

	// Prepare search request to find the user
	attributes := append(myAttributes, "dn")
	searchRequest := ldap.NewSearchRequest(
		myBase,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf(myUserFilter, username),
		attributes,
		nil,
	)

	// Execute the search request
	sr, err := db.Search(searchRequest)
	if err != nil {
		return false, nil, err
	}

	// Handle search results
	if len(sr.Entries) < 1 {
		return false, nil, fmt.Errorf("user does not exist")
	}
	if len(sr.Entries) > 1 {
		return false, nil, fmt.Errorf("too many entries returned")
	}

	// Populate the user result object
	user := &LDAPUserResult{
		Username:       username,
		Attributes:     map[string]string{},
		Groups:         []string{},
		IsLoginSuccess: false,
	}

	userDN := sr.Entries[0].DN

	// Extract user attributes and groups
	for _, attr := range sr.Entries[0].Attributes {
		if attr.Name == "primaryGroupID" {
			for _, value := range attr.Values {
				switch value {
				case "512":
					user.Groups = append(user.Groups, "Domain Admins")
				case "513":
					user.Groups = append(user.Groups, "Domain Users")
				case "519":
					user.Groups = append(user.Groups, "Enterprise Admins")
				case "544":
					user.Groups = append(user.Groups, "Administrators")
				case "548":
					user.Groups = append(user.Groups, "Account Operators")
				case "549":
					user.Groups = append(user.Groups, "Server Operators")
				case "551":
					user.Groups = append(user.Groups, "Backup Operators")
				case "550":
					user.Groups = append(user.Groups, "Print Operators")
				case "518":
					user.Groups = append(user.Groups, "Schema Admins")
				case "517":
					user.Groups = append(user.Groups, "Cert Publishers")
				case "514":
					user.Groups = append(user.Groups, "Guests")
				}
			}
		}
		if attr.Name == "memberOf" {
			for _, value := range attr.Values {
				dn, err := ldap.ParseDN(value)
				if err != nil {
					break
				}
				for _, rdn := range dn.RDNs {
					for _, rdnAttr := range rdn.Attributes {
						user.Groups = append(user.Groups, rdnAttr.Value)
						break // just want the first one
					}
					break // just want the first one
				}
			}
		} else {
			user.Attributes[attr.Name] = sr.Entries[0].GetAttributeValue(attr.Name)
		}
	}

	// Bind with the user's credentials
	err = db.Bind(userDN, password)
	if err != nil {
		return false, user, err
	}

	user.IsLoginSuccess = true

	// Re-bind with the admin DN if necessary
	if myBindDN != "" && myPassword != "" {
		err = db.Bind(myBindDN, myPassword)
		if err != nil {
			return true, user, err
		}
	}

	return true, user, nil
}

// LDAPUserResult represents the result of a user authentication and group retrieval operation
// from the LDAP server. It includes the username, a map of user attributes, a list of groups,
// and a flag indicating whether the login was successful.
type LDAPUserResult struct {
	Username       string            `json:"username,omitempty"`
	Attributes     map[string]string `json:"attributes,omitempty"`
	Groups         []string          `json:"groups,omitempty"`
	IsLoginSuccess bool              `json:"isLoginSuccess,omitempty"`
}

//package aclient_ldap
//
//import (
//	"crypto/tls"
//	"fmt"
//	"strconv"
//	"strings"
//	"sync"
//	"time"
//
//	"github.com/go-ldap/ldap/v3"
//	"github.com/jpfluger/alibs-slim/aconns"
//	"github.com/jpfluger/alibs-slim/asessions"
//	"github.com/jpfluger/alibs-slim/autils"
//)
//
//const (
//	ADAPTERTYPE_LDAP        = aconns.AdapterType("ldap")
//	LDAP_CONNECTION_TIMEOUT = 5
//	LDAP_DEFAULT_PORT       = 389
//	LDAP_DEFAULT_PORT_SSL   = 636
//)
//
//type AClientLDAP struct {
//	aconns.ADBAdapterBase
//
//	Attributes autils.StringsArray `json:"attributes,omitempty"`
//	Base       string              `json:"base,omitempty"`
//	BindDN     string              `json:"bindDN,omitempty"`
//
//	GroupFilter string `json:"groupFilter,omitempty"`
//
//	//ServerName string `json:"serverName,omitempty"`
//
//	UserFilter         string                        `json:"userFilter,omitempty"`
//	InsecureSkipVerify bool                          `json:"insecureSkipVerify,omitempty"`
//	UseSSL             bool                          `json:"useSSL,omitempty"`
//	SkipTLS            bool                          `json:"skipTLS,omitempty"`
//	ClientCertificates []tls.Certificate             `json:"clientCertificates,omitempty"`
//	PermGroups         map[string]asessions.RoleName `json:"permGroups,omitempty"`
//
//	ConnectionTimeout int `json:"connectionTimeout,omitempty"`
//
//	IsHealthy       bool      `json:"-"`
//	LastHealthCheck time.Time `json:"-"`
//
//	ldapPool ILdapConnPool // *LdapConnPool // LDAP connection pool instance
//
//	mu sync.RWMutex
//}
//
//// validate checks if the AClientLDAP object is valid, including essential configurations like
//// the server name, bind DN, and port. It sets defaults if certain values are not provided.
//func (cn *AClientLDAP) validate() error {
//	if err := cn.ADBAdapterBase.Validate(); err != nil {
//		if err != aconns.ErrDatabaseIsEmpty {
//			return err
//		}
//	}
//
//	//cn.ServerName = strings.TrimSpace(cn.ServerName)
//	//if cn.ServerName == "" {
//	//	cn.ServerName = cn.Host
//	//}
//
//	if cn.ConnectionTimeout <= 0 {
//		cn.ConnectionTimeout = LDAP_CONNECTION_TIMEOUT
//	}
//
//	cn.BindDN = strings.TrimSpace(cn.BindDN)
//
//	if cn.Port <= 0 {
//		if cn.UseSSL {
//			cn.Port = LDAP_DEFAULT_PORT_SSL
//		} else {
//			cn.Port = LDAP_DEFAULT_PORT
//		}
//	}
//
//	if err := cn.initPool(); err != nil {
//		return fmt.Errorf("aconns.AClientLDAP init pool failed: %v", err)
//	}
//
//	return nil
//}
//
//// Validate locks the AClientLDAP instance, validates its configuration,
//// and ensures that essential fields are correctly set before use.
//func (cn *AClientLDAP) Validate() error {
//	cn.mu.Lock()
//	defer cn.mu.Unlock()
//	return cn.validate()
//}
//
//// InitPool initializes the LDAP connection pool for the client if it hasn't been initialized yet.
//// It ensures that the pool can be used to manage multiple LDAP connections efficiently.
//func (cn *AClientLDAP) InitPool() error {
//	// Initialize the connection pool using the current configuration.
//	cn.mu.Lock()
//	defer cn.mu.Unlock()
//	return cn.initPool()
//}
//
//func (cn *AClientLDAP) initPool() error {
//	if cn.ldapPool != nil {
//		return nil // Pool already initialized.
//	}
//
//	ldapConfig := &clientLdapConfig{
//		client: cn,
//	}
//
//	cn.ldapPool = InitLdapConnPool(10)
//
//	// Open initial connections if needed.
//	conn, err := cn.ldapPool.GetConnection(ldapConfig)
//	if err != nil {
//		return err
//	}
//	defer cn.ldapPool.PutConnection(conn)
//
//	return nil
//}
//
//// Test validates the AClientLDAP object and checks if a connection can be successfully
//// established with the LDAP server, returning a status indicating success or failure.
//func (cn *AClientLDAP) Test() (bool, aconns.TestStatus, error) {
//	cn.mu.Lock()
//	defer cn.mu.Unlock()
//
//	if err := cn.validate(); err != nil {
//		return false, aconns.TESTSTATUS_FAILED, err
//	}
//
//	if cn.ldapPool == nil {
//		if err := cn.initPool(); err != nil {
//			return false, aconns.TESTSTATUS_FAILED, err
//		}
//	}
//
//	conn, err := cn.ldapPool.GetConnection(&clientLdapConfig{client: cn})
//	if err != nil {
//		return false, aconns.TESTSTATUS_FAILED, err
//	}
//	defer cn.ldapPool.PutConnection(conn)
//
//	// Test the connection.
//	if err = cn.testConnection(conn); err != nil {
//		return false, aconns.TESTSTATUS_FAILED, err
//	}
//
//	return true, aconns.TESTSTATUS_INITIALIZED_SUCCESSFUL, nil
//}
//
//// OpenConnection opens a connection to the LDAP server by initializing the connection pool
//// (if not already initialized) and fetching a connection from it.
//func (cn *AClientLDAP) OpenConnection() error {
//	cn.mu.Lock()
//	defer cn.mu.Unlock()
//	return cn.openConnection()
//}
//
//// openConnection ensures that the connection pool is initialized and retrieves a connection
//// from the pool. This function doesn't manage the connection directly but through the pool.
//func (cn *AClientLDAP) openConnection() error {
//	// Ensure the pool is initialized before attempting to open a connection.
//	if cn.ldapPool == nil {
//		if err := cn.InitPool(); err != nil {
//			return fmt.Errorf("failed to initialize LDAP pool: %v", err)
//		}
//	}
//
//	// Get a connection from the pool.
//	_, err := cn.ldapPool.GetConnection(&clientLdapConfig{client: cn})
//	if err != nil {
//		return fmt.Errorf("failed to get connection from LDAP pool: %v", err)
//	}
//
//	return nil
//}
//
//// GetAddress returns the address of the LDAP server in the format "host:port" by locking
//// access to the AClientLDAP instance to ensure thread-safe operations.
//func (cn *AClientLDAP) GetAddress() string {
//	cn.mu.RLock()
//	defer cn.mu.RUnlock()
//	return cn.getAddress()
//}
//
//// getAddress constructs the server's address in the format "host:port" based on the
//// host and port information provided in the AClientLDAP object.
//func (cn *AClientLDAP) getAddress() string {
//	port := cn.GetPort()
//	return fmt.Sprintf("%s:%s", cn.GetHost(), strconv.Itoa(port))
//}
//
//// testConnection checks if the provided LDAP connection is valid and functional.
//// It returns an error if the connection is invalid or nil.
//func (cn *AClientLDAP) testConnection(db ILdapConn) error {
//	if db == nil {
//		return fmt.Errorf("no ldap conn has been created where host=%s", cn.GetHost())
//	}
//	return nil
//}
//
//// CloseConnection closes all active connections in the LDAP connection pool and resets the pool.
//// It ensures that no connections remain open after calling this function.
//func (cn *AClientLDAP) CloseConnection() error {
//	cn.mu.Lock()
//	defer cn.mu.Unlock()
//
//	if cn.ldapPool == nil {
//		return nil // Pool is already closed.
//	}
//
//	// Delegate the pool closing to the pool's own method.
//	err := cn.ldapPool.CloseAllConnections()
//	if err != nil {
//		return err
//	}
//
//	// Set the pool to nil after it's closed.
//	cn.ldapPool = nil
//	return nil
//}
//
//// GetConnectionTimeout returns the connection timeout value for the LDAP client,
//// ensuring thread-safe access to the ConnectionTimeout field.
//func (cn *AClientLDAP) GetConnectionTimeout() int {
//	cn.mu.RLock()
//	defer cn.mu.RUnlock()
//	return cn.ConnectionTimeout
//}
//
//// DB retrieves an LDAP connection from the connection pool, ensuring thread safety.
//// If the pool is not initialized, it returns an error. The connection should be released after use.
//func (cn *AClientLDAP) DB() (ILdapConn, error) {
//	cn.mu.RLock()
//	defer cn.mu.RUnlock()
//	return cn.db()
//}
//
//// db retrieves an LDAP connection from the connection pool. If the pool is not initialized,
//// it returns an error. This function is responsible for obtaining a connection for operations.
//// Since the pool is handling the connections now, always ensure to release the connection back to the pool after usage. For example:
////
////	db, err := cn.DB()
////	if err != nil {
////	// Handle error
////	}
////	defer cn.ldapPool.PutConnection(db)
//func (cn *AClientLDAP) db() (ILdapConn, error) {
//	// If the connection pool is not initialized, return an error.
//	if cn.ldapPool == nil {
//		return nil, fmt.Errorf("LDAP connection pool is not initialized")
//	}
//
//	// Get a connection from the pool.
//	conn, err := cn.ldapPool.GetConnection(&clientLdapConfig{client: cn})
//	if err != nil {
//		return nil, fmt.Errorf("failed to get connection from LDAP pool: %v", err)
//	}
//
//	return conn, nil
//}
//
//// Authenticate verifies the provided username and password against the LDAP server.
//// It retrieves a connection from the pool, binds with the admin user (if needed), and checks
//// the user's credentials. It also releases the connection back to the pool after use.
//func (cn *AClientLDAP) Authenticate(username, password string) (bool, map[string]string, error) {
//	var db ILdapConn
//	var err error
//
//	// Ensure the connection is returned to the pool after usage
//	defer func() {
//		if db != nil {
//			cn.ldapPool.PutConnection(db) // No need to lock here if the pool manages locking internally
//		}
//	}()
//
//	// Acquire a read lock to safely access shared data in cn
//	cn.mu.RLock()
//
//	// Get a connection from the pool
//	db, err = cn.db()
//	if err != nil {
//		cn.mu.RUnlock() // Unlock early if there's an error
//		return false, nil, fmt.Errorf("failed to get ldap instance: %v", err)
//	}
//
//	// Capture necessary fields to avoid holding the lock for longer than necessary
//	myConnectionTimeout := cn.ConnectionTimeout
//	myAttributes := cn.Attributes
//	myBindDN := cn.BindDN
//	myUserFilter := cn.UserFilter
//
//	// Unlock as soon as we're done accessing shared data
//	cn.mu.RUnlock()
//
//	// Bind with the admin DN if necessary
//	bindPassword := cn.GetPassword()
//	if myBindDN != "" && bindPassword != "" {
//		err = db.Bind(myBindDN, bindPassword)
//		if err != nil {
//			return false, nil, fmt.Errorf("failed to bind with admin credentials: %v", err)
//		}
//	}
//
//	// Prepare search request to find the user
//	attributes := append(myAttributes, "dn")
//	searchRequest := ldap.NewSearchRequest(
//		cn.Base,
//		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, myConnectionTimeout, false,
//		fmt.Sprintf(myUserFilter, username),
//		attributes,
//		nil,
//	)
//
//	// Execute the search request
//	sr, err := db.Search(searchRequest)
//	if err != nil {
//		return false, nil, fmt.Errorf("search failed: %v", err)
//	}
//
//	// Handle search results
//	if len(sr.Entries) < 1 {
//		return false, nil, fmt.Errorf("user does not exist")
//	}
//	if len(sr.Entries) > 1 {
//		return false, nil, fmt.Errorf("too many entries returned")
//	}
//
//	// Get the user's DN and attributes
//	userDN := sr.Entries[0].DN
//	user := map[string]string{}
//	for _, attr := range myAttributes { // Use captured myAttributes
//		user[attr] = sr.Entries[0].GetAttributeValue(attr)
//	}
//
//	// Bind with the user's credentials
//	err = db.Bind(userDN, password)
//	if err != nil {
//		return false, user, fmt.Errorf("failed to bind user: %v", err)
//	}
//
//	// Re-bind with the admin DN if necessary
//	if myBindDN != "" && bindPassword != "" {
//		err = db.Bind(myBindDN, bindPassword)
//		if err != nil {
//			return true, user, fmt.Errorf("failed to re-bind with admin credentials: %v", err)
//		}
//	}
//
//	return true, user, nil
//}
//
//// GetGroupsOfUser searches for the specified user's groups in the LDAP directory
//// by retrieving a connection from the pool, performing a search, and collecting the group names.
//// It releases the connection back to the pool after use.
//func (cn *AClientLDAP) GetGroupsOfUser(ldapSearchUser string) ([]string, error) {
//	var db ILdapConn
//	var err error
//
//	// Ensure the connection is returned to the pool after usage
//	defer func() {
//		if db != nil {
//			cn.ldapPool.PutConnection(db)
//		}
//	}()
//
//	// Acquire a read lock to safely access shared data in cn
//	cn.mu.RLock()
//
//	// Get a connection from the pool
//	db, err = cn.DB()
//	if err != nil {
//		cn.mu.RUnlock() // Unlock early if there's an error
//		return nil, err
//	}
//
//	// Capture shared values
//	myBase := cn.Base
//	myGroupFilter := cn.GroupFilter
//
//	// Unlock as soon as we're done accessing shared data
//	cn.mu.RUnlock()
//
//	// Prepare search request to find the groups for the user
//	searchRequest := ldap.NewSearchRequest(
//		myBase,
//		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
//		fmt.Sprintf(myGroupFilter, ldapSearchUser),
//		[]string{"cn"},
//		nil,
//	)
//
//	// Execute the search request
//	sr, err := db.Search(searchRequest)
//	if err != nil {
//		return nil, err
//	}
//
//	// Collect the groups from the search results
//	groups := []string{}
//	for _, entry := range sr.Entries {
//		groups = append(groups, entry.GetAttributeValue("cn"))
//	}
//
//	return groups, nil
//}
//
//// AuthenticateWithGroups verifies the user's credentials and retrieves the user's group information
//// from the LDAP directory. It binds with the admin account (if needed) and ensures that
//// the connection is returned to the pool after use.
//func (cn *AClientLDAP) AuthenticateWithGroups(username, password string) (bool, *LDAPUserResult, error) {
//	var db ILdapConn
//	var err error
//
//	// Ensure the connection is returned to the pool after usage
//	defer func() {
//		if db != nil {
//			cn.ldapPool.PutConnection(db)
//		}
//	}()
//
//	// Acquire a read lock to safely access shared data in cn
//	cn.mu.RLock()
//
//	// Get a connection from the pool
//	db, err = cn.DB()
//	if err != nil {
//		cn.mu.RUnlock() // Unlock early if there's an error
//		return false, nil, err
//	}
//
//	// Capture shared values
//	myBindDN := cn.BindDN
//	myPassword := cn.GetPassword()
//	myBase := cn.Base
//	myAttributes := cn.Attributes
//	myUserFilter := cn.UserFilter
//
//	// Unlock as soon as we're done accessing shared data
//	cn.mu.RUnlock()
//
//	// Bind with the admin DN if necessary
//	if myBindDN != "" && myPassword != "" {
//		err = db.Bind(myBindDN, myPassword)
//		if err != nil {
//			return false, nil, err
//		}
//	}
//
//	// Prepare search request to find the user
//	attributes := append(myAttributes, "dn")
//	searchRequest := ldap.NewSearchRequest(
//		myBase,
//		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
//		fmt.Sprintf(myUserFilter, username),
//		attributes,
//		nil,
//	)
//
//	// Execute the search request
//	sr, err := db.Search(searchRequest)
//	if err != nil {
//		return false, nil, err
//	}
//
//	// Handle search results
//	if len(sr.Entries) < 1 {
//		return false, nil, fmt.Errorf("user does not exist")
//	}
//	if len(sr.Entries) > 1 {
//		return false, nil, fmt.Errorf("too many entries returned")
//	}
//
//	// Populate the user result object
//	user := &LDAPUserResult{
//		Username:       username,
//		Attributes:     map[string]string{},
//		Groups:         []string{},
//		IsLoginSuccess: false,
//	}
//
//	userDN := sr.Entries[0].DN
//
//	// Extract user attributes and groups
//	for _, attr := range sr.Entries[0].Attributes {
//		if attr.Name == "primaryGroupID" {
//			for _, value := range attr.Values {
//				switch value {
//				case "512":
//					user.Groups = append(user.Groups, "Domain Admins")
//				case "513":
//					user.Groups = append(user.Groups, "Domain Users")
//				case "519":
//					user.Groups = append(user.Groups, "Enterprise Admins")
//				case "544":
//					user.Groups = append(user.Groups, "Administrators")
//				case "548":
//					user.Groups = append(user.Groups, "Account Operators")
//				case "549":
//					user.Groups = append(user.Groups, "Server Operators")
//				case "551":
//					user.Groups = append(user.Groups, "Backup Operators")
//				case "550":
//					user.Groups = append(user.Groups, "Print Operators")
//				case "518":
//					user.Groups = append(user.Groups, "Schema Admins")
//				case "517":
//					user.Groups = append(user.Groups, "Cert Publishers")
//				case "514":
//					user.Groups = append(user.Groups, "Guests")
//				}
//			}
//		}
//		if attr.Name == "memberOf" {
//			for _, value := range attr.Values {
//				dn, err := ldap.ParseDN(value)
//				if err != nil {
//					break
//				}
//				for _, rdn := range dn.RDNs {
//					for _, rdnAttr := range rdn.Attributes {
//						user.Groups = append(user.Groups, rdnAttr.Value)
//						break // just want the first one
//					}
//					break // just want the first one
//				}
//			}
//		} else {
//			user.Attributes[attr.Name] = sr.Entries[0].GetAttributeValue(attr.Name)
//		}
//	}
//
//	// Bind with the user's credentials
//	err = db.Bind(userDN, password)
//	if err != nil {
//		return false, user, err
//	}
//
//	user.IsLoginSuccess = true
//
//	// Re-bind with the admin DN if necessary
//	if myBindDN != "" && myPassword != "" {
//		err = db.Bind(myBindDN, myPassword)
//		if err != nil {
//			return true, user, err
//		}
//	}
//
//	return true, user, nil
//}
//
//// LDAPUserResult represents the result of a user authentication and group retrieval operation
//// from the LDAP server. It includes the username, a map of user attributes, a list of groups,
//// and a flag indicating whether the login was successful.
//type LDAPUserResult struct {
//	Username       string
//	Attributes     map[string]string
//	Groups         []string
//	IsLoginSuccess bool
//}
