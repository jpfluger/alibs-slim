package aconns

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// MockSBAdapter is a mock implementation of ISBAdapterSql for testing purposes.
type MockSBAdapter struct {
	adapterType AdapterType
	adapterName AdapterName
	host        string
	port        int
	database    string
	username    string
	password    string
}

func (m *MockSBAdapter) GetType() AdapterType {
	return m.adapterType
}

func (m *MockSBAdapter) GetName() AdapterName {
	return m.adapterName
}

func (m *MockSBAdapter) GetHost() string {
	return m.host
}

func (m *MockSBAdapter) GetPort() int {
	return m.port
}

func (m *MockSBAdapter) GetDatabase() string {
	return m.database
}

func (m *MockSBAdapter) GetUsername() string {
	return m.username
}

func (m *MockSBAdapter) GetPassword() string {
	return m.password
}

func (m *MockSBAdapter) Validate() error {
	if m.database == "" || m.username == "" || m.password == "" {
		return fmt.Errorf("validation failed: missing database, username, or password")
	}
	return nil
}

func (m *MockSBAdapter) Test() (bool, TestStatus, error) {
	if err := m.Validate(); err != nil {
		return false, TESTSTATUS_FAILED, err
	}
	return true, TESTSTATUS_INITIALIZED_SUCCESSFUL, nil
}

func TestNewSBAdapterSql(t *testing.T) {
	db := &sql.DB{}
	adapter := &MockSBAdapter{
		adapterType: "mockType",
		adapterName: "mockName",
		host:        "localhost",
	}
	sbAdapter := NewSBAdapterSql(adapter, db)
	assert.NotNil(t, sbAdapter)
	assert.Equal(t, "mockType", sbAdapter.GetType().String())
	assert.Equal(t, "mockName", sbAdapter.GetName().String())
	assert.Equal(t, "localhost", sbAdapter.GetHost())
}

func TestSBAdapterSql_SupportsModels(t *testing.T) {
	db := &sql.DB{}
	adapter := &MockSBAdapter{}
	sbAdapter := NewSBAdapterSql(adapter, db)
	assert.False(t, sbAdapter.SupportsModels())
}

func TestSBAdapterSql_Query(t *testing.T) {
	db := &sql.DB{}
	adapter := &MockSBAdapter{}
	sbAdapter := NewSBAdapterSql(adapter, db)
	_, err := sbAdapter.Query("SELECT * FROM test")
	assert.Error(t, err)
}

func TestSBAdapterSql_QueryArgs(t *testing.T) {
	db := &sql.DB{}
	adapter := &MockSBAdapter{}
	sbAdapter := NewSBAdapterSql(adapter, db)
	_, err := sbAdapter.QueryArgs("SELECT * FROM test WHERE id = ?", 1)
	assert.Error(t, err)
}

func TestSBAdapterSql_QueryModel(t *testing.T) {
	db := &sql.DB{}
	adapter := &MockSBAdapter{}
	sbAdapter := NewSBAdapterSql(adapter, db)
	err := sbAdapter.QueryModel("SELECT * FROM test", nil)
	assert.Error(t, err)
}

func TestSBAdapterSql_QueryModelArgs(t *testing.T) {
	db := &sql.DB{}
	adapter := &MockSBAdapter{}
	sbAdapter := NewSBAdapterSql(adapter, db)
	err := sbAdapter.QueryModelArgs("SELECT * FROM test WHERE id = ?", nil, 1)
	assert.Error(t, err)
}

func TestSBAdapterSql_QueryArgs_PanicRecovery(t *testing.T) {
	db := &sql.DB{}
	adapter := &MockSBAdapter{}
	sbAdapter := NewSBAdapterSql(adapter, db)

	// Simulate a panic scenario
	defer func() {
		if r := recover(); r != nil {
			assert.Fail(t, "Panic occurred during QueryArgs")
		}
	}()
	_, err := sbAdapter.QueryArgs("SELECT * FROM test WHERE id = ?", 1)
	assert.Error(t, err)
}
