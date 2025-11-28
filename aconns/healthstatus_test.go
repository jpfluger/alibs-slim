package aconns

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// MockAdapter extends DummyAdapter for testing, with counters and controllable behavior.
type MockAdapter struct {
	DummyAdapter
	mu          sync.Mutex
	testCalls   int
	healthCalls int
	mockHealth  *HealthCheck // Assume HealthCheck is defined elsewhere; mock it here.
	shouldFail  bool         // Inherited, but override if needed.
}

func (m *MockAdapter) Test() (bool, TestStatus, error) {
	m.mu.Lock()
	m.testCalls++
	m.mu.Unlock()
	if m.shouldFail {
		return false, TESTSTATUS_FAILED, fmt.Errorf("mock test failed")
	}
	return true, TESTSTATUS_INITIALIZED_SUCCESSFUL, nil
}

func (m *MockAdapter) GetHealth() *HealthCheck {
	m.mu.Lock()
	m.healthCalls++
	m.mu.Unlock()
	return m.mockHealth
}

// Helper to get call counts (thread-safe).
func (m *MockAdapter) GetTestCalls() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.testCalls
}

func (m *MockAdapter) GetHealthCalls() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.healthCalls
}

func TestConn_StartMonitor(t *testing.T) {
	tests := []struct {
		name              string
		initialValidated  bool
		mockHealth        *HealthCheck
		shouldFailTest    bool
		interval          time.Duration
		waitDuration      time.Duration
		expectedTestCalls int // Minimum expected calls (e.g., >0)
		expectedValidated bool
	}{
		{
			name:              "Triggers Validate on Invalid State with Successful Test",
			initialValidated:  false,
			mockHealth:        NewHealthCheck(HEALTHSTATUS_HEALTHY), // Assume healthy, not stale/failed.
			shouldFailTest:    false,
			interval:          100 * time.Millisecond,
			waitDuration:      300 * time.Millisecond, // Enough for 2-3 ticks.
			expectedTestCalls: 1,                      // At least one call.
			expectedValidated: true,
		},
		{
			name:              "Handles Stale Health by Calling Test",
			initialValidated:  true,
			mockHealth:        &HealthCheck{LastCheck: time.Now().Add(-time.Hour)}, // Stale.
			shouldFailTest:    false,
			interval:          100 * time.Millisecond,
			waitDuration:      300 * time.Millisecond,
			expectedTestCalls: 1,
			expectedValidated: true, // Remains true.
		},
		{
			name:              "No Action on Healthy Validated State",
			initialValidated:  true,
			mockHealth:        NewHealthCheck(HEALTHSTATUS_HEALTHY), // Fresh and healthy.
			shouldFailTest:    false,
			interval:          100 * time.Millisecond,
			waitDuration:      300 * time.Millisecond,
			expectedTestCalls: 0, // No needsCheck.
			expectedValidated: true,
		},
		{
			name:              "Handles Failed Test Without Panic",
			initialValidated:  false,
			mockHealth:        NewHealthCheck(HEALTHSTATUS_HEALTHY),
			shouldFailTest:    true,
			interval:          100 * time.Millisecond,
			waitDuration:      300 * time.Millisecond,
			expectedTestCalls: 1,
			expectedValidated: true, // Validated during Test(), even if adapter Test() fails.
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &MockAdapter{
				mockHealth: tt.mockHealth,
				shouldFail: tt.shouldFailTest,
			}
			conn := &Conn{
				Adapter:     mock,
				IsValidated: tt.initialValidated,
			}

			// Start the monitor.
			conn.StartMonitor(tt.interval)

			// Wait for ticks to occur.
			time.Sleep(tt.waitDuration)

			// Assertions.
			assert.Equal(t, tt.expectedValidated, conn.GetIsValidated())
			assert.GreaterOrEqual(t, mock.GetTestCalls(), tt.expectedTestCalls)
			assert.Greater(t, mock.GetHealthCalls(), 0) // Always called in loop.
		})
	}
}
