package aclient_redis

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/alexedwards/scs/redisstore"
	"github.com/gomodule/redigo/redis"
	"github.com/jpfluger/alibs-slim/aconns"
	"github.com/jpfluger/alibs-slim/autils"
)

const (
	ADAPTERTYPE_REDIS         = aconns.AdapterType("redis")
	REDIS_DEFAULT_PORT        = 6379
	REDIS_MAXIDLE_CONNECTIONS = 50
	REDIS_IDLETIMEOUT         = time.Duration(240 * time.Minute)
)

// AClientRedis represents a Redis client adapter.
type AClientRedis struct {
	aconns.ADBAdapterBase

	MaxIdleConnections int           `json:"maxIdleConnections,omitempty"`
	IdleTimeout        time.Duration `json:"idleTimeout,omitempty"`

	pool  *redis.Pool
	rsMap *RedisStoreMap

	mu sync.RWMutex
}

// validate checks if the AClientRedis object is valid.
func (cn *AClientRedis) validate() error {
	if err := cn.ADBAdapterBase.Validate(); err != nil {
		if err != aconns.ErrDatabaseIsEmpty {
			return err
		}
	}

	if cn.MaxIdleConnections <= 0 {
		cn.MaxIdleConnections = REDIS_MAXIDLE_CONNECTIONS
	}

	if cn.Port <= 0 {
		cn.Port = REDIS_DEFAULT_PORT
	}

	cn.Database = autils.ToStringTrimLower(cn.Database)
	if cn.Database == "" {
		cn.Database = "0"
	}
	if err := cn.ValidateDb(cn.Database); err != nil {
		return fmt.Errorf("redis requires an int for the database; %v", err)
	}

	if cn.IdleTimeout <= 0 {
		cn.IdleTimeout = REDIS_IDLETIMEOUT
	}

	return nil
}

// Validate checks if the AClientRedis object is valid.
func (cn *AClientRedis) Validate() error {
	cn.mu.Lock()
	defer cn.mu.Unlock()
	return cn.validate()
}

// Test attempts to validate the AClientRedis, open a connection if necessary, and test the connection.
func (cn *AClientRedis) Test() (bool, aconns.TestStatus, error) {
	cn.mu.Lock()
	defer cn.mu.Unlock()
	return cn.test()
}

// Test attempts to validate the AClientRedis, open a connection if necessary, and test the connection.
func (cn *AClientRedis) test() (bool, aconns.TestStatus, error) {
	if err := cn.validate(); err != nil {
		cn.UpdateHealth(aconns.HEALTHSTATUS_VALIDATE_FAILED)
		return false, aconns.TESTSTATUS_FAILED, err
	}

	if cn.pool == nil {
		if err := cn.openConnection(); err != nil {
			cn.UpdateHealth(aconns.HEALTHSTATUS_OPEN_FAILED)
			return false, aconns.TESTSTATUS_FAILED, err
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) // Ping timeout
	defer cancel()
	if err := cn.testConnectionWithCtx(ctx, cn.pool); err != nil {
		status := aconns.HEALTHSTATUS_PING_FAILED
		if context.DeadlineExceeded == err {
			status = aconns.HEALTHSTATUS_TIMEOUT
		} else if strings.Contains(err.Error(), "network") || strings.Contains(err.Error(), "connection refused") {
			status = aconns.HEALTHSTATUS_NETWORK_ERROR
		}
		cn.UpdateHealth(status)
		//alog.LOGGER(alog.LOGGER_APP).Warn().Err(err).Msg("Redis test failed")
		return false, aconns.TESTSTATUS_FAILED, err
	}

	cn.UpdateHealth(aconns.HEALTHSTATUS_HEALTHY)
	return true, aconns.TESTSTATUS_INITIALIZED_SUCCESSFUL, nil
}

// OpenConnection opens a connection to the Redis server.
func (cn *AClientRedis) OpenConnection() error {
	cn.mu.Lock()
	defer cn.mu.Unlock()
	return cn.openConnection()
}

// openConnection opens a connection to the Redis server.
func (cn *AClientRedis) openConnection() error {
	pool := &redis.Pool{
		MaxIdle: cn.MaxIdleConnections,
		Dial: func() (redis.Conn, error) {
			return dialWithDB("tcp", cn.getAddress(), cn.GetPassword(), cn.GetDatabase())
		},
		IdleTimeout: cn.IdleTimeout,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) // Open timeout
	defer cancel()
	if err := cn.testConnectionWithCtx(ctx, pool); err != nil {
		return err
	}

	cn.pool = pool
	return nil
}

func (cn *AClientRedis) GetAddress() string {
	cn.mu.RLock()
	defer cn.mu.RUnlock()
	return cn.getAddress()
}

func (cn *AClientRedis) getAddress() string {
	port := cn.Port
	return fmt.Sprintf("%s:%s", cn.Host, strconv.Itoa(port))
}

func (cn *AClientRedis) testConnectionWithCtx(ctx context.Context, pool *redis.Pool) error {
	if pool == nil {
		return fmt.Errorf("no redis pool has been created where host=%s", cn.Host)
	}

	dbc := pool.Get()
	defer dbc.Close()

	done := make(chan error, 1)
	go func() {
		_, err := dbc.Do("PING")
		done <- err
	}()

	select {
	case err := <-done:
		if err != nil {
			return fmt.Errorf("test connection failed for redis where host=%s; %v", cn.Host, err)
		}
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (cn *AClientRedis) CloseConnection() error {
	cn.mu.Lock()
	defer cn.mu.Unlock()

	if cn.pool != nil {
		if err := cn.pool.Close(); err != nil {
			return fmt.Errorf("error in closing the redis client-pool; %v", err)
		}
		cn.pool = nil
		cn.UpdateHealth(aconns.HEALTHSTATUS_CLOSED)
	}

	return nil
}

func (cn *AClientRedis) Pool() *redis.Pool {
	cn.mu.RLock()
	if cn.IsHealthy() && !cn.GetHealth().IsStale(5*time.Minute) {
		defer cn.mu.RUnlock()
		return cn.pool
	}
	cn.mu.RUnlock()

	// Upgrade to write lock for refresh
	cn.mu.Lock()
	defer cn.mu.Unlock()
	cn.test() // Refresh and test
	return cn.pool
}

func (cn *AClientRedis) Refresh() error {
	cn.mu.Lock()
	defer cn.mu.Unlock()
	if cn.pool != nil {
		cn.pool.Close()
		cn.pool = nil
	}
	return cn.openConnection()
}

func (cn *AClientRedis) GetRedisStore(prefix string) *redisstore.RedisStore {
	cn.mu.Lock()
	defer cn.mu.Unlock()

	if cn.pool == nil {
		return nil
	}
	prefix = strings.TrimSpace(prefix)
	if prefix == "" {
		return nil
	}
	if cn.rsMap == nil {
		cn.rsMap = NewRedisStoreMap()
	}
	rsStore, exists := cn.rsMap.Get(prefix)
	if exists {
		return rsStore
	}
	newStore := redisstore.NewWithPrefix(cn.pool, prefix)
	cn.rsMap.Set(prefix, newStore)
	return newStore
}

func dial(network, address, password string) (redis.Conn, error) {
	c, err := redis.Dial(network, address)
	if err != nil {
		return nil, err
	}
	if password != "" {
		if _, err := c.Do("AUTH", password); err != nil {
			c.Close()
			return nil, err
		}
	}
	return c, err
}

func dialWithDB(network, address, password, DB string) (redis.Conn, error) {
	c, err := dial(network, address, password)
	if err != nil {
		return nil, err
	}
	if _, err := c.Do("SELECT", DB); err != nil {
		c.Close()
		return nil, err
	}
	return c, err
}

func (cn *AClientRedis) ValidateDb(dbName string) error {
	if dbName == "" {
		return fmt.Errorf("dbName parameter is empty")
	}
	ii, err := strconv.Atoi(dbName)
	if err != nil {
		return fmt.Errorf("dbName (=%s) must be a number between 0 to 15 (including 15); %v", dbName, err)
	}
	if ii >= 0 && ii < 16 {
		return nil
	}
	return fmt.Errorf("dbName (=%s) is out of range; dbName should be an integer between 0 to 15 (including 15)", dbName)
}

//package aclient_redis
//
//import (
//	"fmt"
//	"strconv"
//	"strings"
//	"sync"
//	"time"
//
//	"github.com/alexedwards/scs/redisstore"
//	"github.com/gomodule/redigo/redis"
//	"github.com/jpfluger/alibs-slim/aconns"
//	"github.com/jpfluger/alibs-slim/autils"
//)
//
//const (
//	ADAPTERTYPE_REDIS         = aconns.AdapterType("redis")
//	REDIS_DEFAULT_PORT        = 6379
//	REDIS_MAXIDLE_CONNECTIONS = 50
//	REDIS_IDLETIMEOUT         = time.Duration(240 * time.Minute)
//)
//
//type AClientRedis struct {
//	aconns.ADBAdapterBase
//
//	MaxIdleConnections int           `json:"maxIdleConnections,omitempty"`
//	IdleTimeout        time.Duration `json:"idleTimeout,omitempty"`
//
//	pool  *redis.Pool
//	rsMap *RedisStoreMap
//
//	mu sync.RWMutex
//}
//
//// validate checks if the AClientRedis object is valid.
//func (cn *AClientRedis) validate() error {
//	if err := cn.ADBAdapterBase.Validate(); err != nil {
//		if err != aconns.ErrDatabaseIsEmpty {
//			return err
//		}
//	}
//
//	if cn.MaxIdleConnections <= 0 {
//		cn.MaxIdleConnections = REDIS_MAXIDLE_CONNECTIONS
//	}
//
//	if cn.Port <= 0 {
//		cn.Port = REDIS_DEFAULT_PORT
//	}
//
//	cn.Database = autils.ToStringTrimLower(cn.Database)
//	if cn.Database == "" {
//		cn.Database = "0"
//	}
//	if err := cn.ValidateDb(cn.Database); err != nil {
//		return fmt.Errorf("redis requires an int for the database; %v", err)
//	}
//
//	if cn.IdleTimeout <= 0 {
//		cn.IdleTimeout = REDIS_IDLETIMEOUT
//	}
//
//	return nil
//}
//
//// Validate checks if the AClientRedis object is valid.
//func (cn *AClientRedis) Validate() error {
//	cn.mu.Lock()
//	defer cn.mu.Unlock()
//	return cn.validate()
//}
//
//// Test attempts to validate the AClientRedis, open a connection if necessary, and test the connection.
//func (cn *AClientRedis) Test() (bool, aconns.TestStatus, error) {
//	cn.mu.Lock()
//	defer cn.mu.Unlock()
//
//	if err := cn.validate(); err != nil {
//		return false, aconns.TESTSTATUS_FAILED, err
//	}
//
//	if cn.pool != nil {
//		if err := cn.testConnection(cn.pool); err == nil {
//			return true, aconns.TESTSTATUS_INITIALIZED_SUCCESSFUL, nil
//		}
//	}
//
//	if err := cn.openConnection(); err != nil {
//		return false, aconns.TESTSTATUS_FAILED, err
//	}
//
//	return true, aconns.TESTSTATUS_INITIALIZED_SUCCESSFUL, nil
//}
//
//// OpenConnection opens a connection to the Redis server.
//func (cn *AClientRedis) OpenConnection() error {
//	cn.mu.Lock()
//	defer cn.mu.Unlock()
//	return cn.openConnection()
//}
//
//// openConnection opens a connection to the Redis server.
//func (cn *AClientRedis) openConnection() error {
//	pool := &redis.Pool{
//		MaxIdle: cn.MaxIdleConnections,
//		Dial: func() (redis.Conn, error) {
//			return dialWithDB("tcp", cn.getAddress(), cn.GetPassword(), cn.GetDatabase())
//		},
//		IdleTimeout: cn.IdleTimeout,
//	}
//
//	if err := cn.testConnection(pool); err != nil {
//		return err
//	}
//
//	cn.pool = pool
//	return nil
//}
//
//func (cn *AClientRedis) GetAddress() string {
//	cn.mu.RLock()
//	defer cn.mu.RUnlock()
//	return cn.getAddress()
//}
//
//func (cn *AClientRedis) getAddress() string {
//	port := cn.Port
//	return fmt.Sprintf("%s:%s", cn.Host, strconv.Itoa(port))
//}
//
//func (cn *AClientRedis) testConnection(pool *redis.Pool) error {
//	if pool == nil {
//		return fmt.Errorf("no redis pool has been created where host=%s", cn.Host)
//	}
//
//	dbc := pool.Get()
//	defer dbc.Close()
//
//	_, err := dbc.Do("PING")
//	if err != nil {
//		return fmt.Errorf("test connection failed for redis where host=%s; %v", cn.Host, err)
//	}
//
//	return nil
//}
//
//func (cn *AClientRedis) CloseConnection() error {
//	cn.mu.Lock()
//	defer cn.mu.Unlock()
//
//	if cn.pool != nil {
//		if err := cn.pool.Close(); err != nil {
//			return fmt.Errorf("error in closing the redis client-pool; %v", err)
//		}
//		cn.pool = nil
//	}
//
//	return nil
//}
//
//func (cn *AClientRedis) Pool() *redis.Pool {
//	cn.mu.RLock()
//	defer cn.mu.RUnlock()
//	return cn.pool
//}
//
//func (cn *AClientRedis) GetRedisStore(prefix string) *redisstore.RedisStore {
//	cn.mu.Lock()
//	defer cn.mu.Unlock()
//
//	if cn.pool == nil {
//		return nil
//	}
//	prefix = strings.TrimSpace(prefix)
//	if prefix == "" {
//		return nil
//	}
//	if cn.rsMap == nil {
//		cn.rsMap = NewRedisStoreMap()
//	}
//	rsStore, exists := cn.rsMap.Get(prefix)
//	if exists {
//		return rsStore
//	}
//	newStore := redisstore.NewWithPrefix(cn.pool, prefix)
//	cn.rsMap.Set(prefix, newStore)
//	return newStore
//}
//
//func dial(network, address, password string) (redis.Conn, error) {
//	c, err := redis.Dial(network, address)
//	if err != nil {
//		return nil, err
//	}
//	if password != "" {
//		if _, err := c.Do("AUTH", password); err != nil {
//			c.Close()
//			return nil, err
//		}
//	}
//	return c, err
//}
//
//func dialWithDB(network, address, password, DB string) (redis.Conn, error) {
//	c, err := dial(network, address, password)
//	if err != nil {
//		return nil, err
//	}
//	if _, err := c.Do("SELECT", DB); err != nil {
//		c.Close()
//		return nil, err
//	}
//	return c, err
//}
//
//func (cn *AClientRedis) ValidateDb(dbName string) error {
//	if dbName == "" {
//		return fmt.Errorf("dbName parameter is empty")
//	}
//	ii, err := strconv.Atoi(dbName)
//	if err != nil {
//		return fmt.Errorf("dbName (=%s) must be a number between 0 to 15 (including 15); %v", dbName, err)
//	}
//	if ii >= 0 && ii < 16 {
//		return nil
//	}
//	return fmt.Errorf("dbName (=%s) is out of range; dbName should be an integer between 0 to 15 (including 15)", dbName)
//}
