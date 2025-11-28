package aclient_badger

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/alexedwards/scs/badgerstore"
	"github.com/dgraph-io/badger/v4"
	"github.com/jpfluger/alibs-slim/aconns"
)

const (
	ADAPTERTYPE_BADGER = aconns.AdapterType("badger")
)

// AClientBadger represents a Badger client adapter.
type AClientBadger struct {
	aconns.ADBAdapterBase

	db    *badger.DB
	bsMap *BadgerStoreMap

	mu sync.RWMutex
}

// validate checks if the AClientBadger object is valid.
func (cn *AClientBadger) validate() error {
	if cn.Host == "" {
		cn.Host = "local"
	}

	if cn.Port <= 0 {
		cn.Port = 0
	}

	if cn.Username == "" {
		cn.Username = "badger"
	}

	cn.Database = strings.TrimSpace(cn.Database)
	if cn.Database == "" {
		return fmt.Errorf("database (directory path) is empty")
	}

	err := cn.ADBAdapterBase.Validate()
	if err == aconns.ErrPasswordIsEmpty {
		err = nil // Password is optional for Badger (used for encryption)
	}
	if err != nil && err != aconns.ErrDatabaseIsEmpty {
		return err
	}

	if cn.Password != "" {
		keyLen := len(cn.Password)
		if keyLen != 16 && keyLen != 24 && keyLen != 32 {
			return fmt.Errorf("encryption key must be 16, 24, or 32 bytes long; got %d", keyLen)
		}
	}

	return nil
}

// Validate checks if the AClientBadger object is valid.
func (cn *AClientBadger) Validate() error {
	cn.mu.Lock()
	defer cn.mu.Unlock()
	return cn.validate()
}

// Test attempts to validate the AClientBadger, open a connection if necessary, and test the connection.
func (cn *AClientBadger) Test() (bool, aconns.TestStatus, error) {
	cn.mu.Lock()
	defer cn.mu.Unlock()
	return cn.test()
}

// Test attempts to validate the AClientBadger, open a connection if necessary, and test the connection.
func (cn *AClientBadger) test() (bool, aconns.TestStatus, error) {
	if err := cn.validate(); err != nil {
		cn.UpdateHealth(aconns.HEALTHSTATUS_VALIDATE_FAILED)
		return false, aconns.TESTSTATUS_FAILED, err
	}

	if cn.db == nil {
		if err := cn.openConnection(); err != nil {
			cn.UpdateHealth(aconns.HEALTHSTATUS_OPEN_FAILED)
			return false, aconns.TESTSTATUS_FAILED, err
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) // Test timeout
	defer cancel()
	if err := cn.testConnectionWithCtx(ctx, cn.db); err != nil {
		status := aconns.HEALTHSTATUS_PING_FAILED
		if context.DeadlineExceeded == err {
			status = aconns.HEALTHSTATUS_TIMEOUT
		} else if strings.Contains(err.Error(), "network") || strings.Contains(err.Error(), "connection refused") {
			status = aconns.HEALTHSTATUS_NETWORK_ERROR
		}
		cn.UpdateHealth(status)
		//alog.LOGGER(alog.LOGGER_APP).Warn().Err(err).Msg("Badger test failed")
		return false, aconns.TESTSTATUS_FAILED, err
	}

	cn.UpdateHealth(aconns.HEALTHSTATUS_HEALTHY)
	return true, aconns.TESTSTATUS_INITIALIZED_SUCCESSFUL, nil
}

// OpenConnection opens a connection to the Badger database.
func (cn *AClientBadger) OpenConnection() error {
	cn.mu.Lock()
	defer cn.mu.Unlock()
	return cn.openConnection()
}

// openConnection opens a connection to the Badger database.
func (cn *AClientBadger) openConnection() error {
	opts := badger.DefaultOptions(cn.Database)
	if cn.GetPassword() != "" {
		opts = opts.WithEncryptionKey([]byte(cn.GetPassword()))
	}

	db, err := badger.Open(opts)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) // Open timeout
	defer cancel()
	if err := cn.testConnectionWithCtx(ctx, db); err != nil {
		db.Close()
		return err
	}

	cn.db = db
	return nil
}

func (cn *AClientBadger) GetAddress() string {
	cn.mu.RLock()
	defer cn.mu.RUnlock()
	return cn.getAddress()
}

func (cn *AClientBadger) getAddress() string {
	port := cn.Port
	return fmt.Sprintf("%s:%s", cn.Host, strconv.Itoa(port))
}

func (cn *AClientBadger) testConnectionWithCtx(ctx context.Context, db *badger.DB) error {
	if db == nil {
		return fmt.Errorf("no badger db has been created where host=%s", cn.Host)
	}

	done := make(chan error, 1)
	go func() {
		err := db.View(func(txn *badger.Txn) error {
			return nil
		})
		done <- err
	}()

	select {
	case err := <-done:
		if err != nil {
			return fmt.Errorf("test connection failed for badger where host=%s; %v", cn.Host, err)
		}
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (cn *AClientBadger) CloseConnection() error {
	cn.mu.Lock()
	defer cn.mu.Unlock()

	if cn.db != nil {
		if err := cn.db.Close(); err != nil {
			return fmt.Errorf("error in closing the badger db; %v", err)
		}
		cn.db = nil
		cn.UpdateHealth(aconns.HEALTHSTATUS_CLOSED)
	}

	return nil
}

func (cn *AClientBadger) DB() *badger.DB {
	cn.mu.RLock()
	if cn.IsHealthy() && !cn.GetHealth().IsStale(5*time.Minute) {
		defer cn.mu.RUnlock()
		return cn.db
	}
	cn.mu.RUnlock()

	// Upgrade to write lock for refresh
	cn.mu.Lock()
	defer cn.mu.Unlock()
	cn.test() // Refresh and test
	return cn.db
}

func (cn *AClientBadger) Refresh() error {
	cn.mu.Lock()
	defer cn.mu.Unlock()
	if cn.db != nil {
		cn.db.Close()
		cn.db = nil
	}
	return cn.openConnection()
}

func (cn *AClientBadger) GetBadgerStore(prefix string) *badgerstore.BadgerStore {
	cn.mu.Lock()
	defer cn.mu.Unlock()

	if cn.db == nil {
		return nil
	}
	prefix = strings.TrimSpace(prefix)
	if prefix == "" {
		return nil
	}
	if cn.bsMap == nil {
		cn.bsMap = NewBadgerStoreMap()
	}
	bsStore, exists := cn.bsMap.Get(prefix)
	if exists {
		return bsStore
	}
	newStore := badgerstore.NewWithPrefix(cn.db, prefix)
	cn.bsMap.Set(prefix, newStore)
	return newStore
}
