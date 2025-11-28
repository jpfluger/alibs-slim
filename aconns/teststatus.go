package aconns

import (
	"fmt"
	"time"
)

// TestStatus represents the status of the connection test.
type TestStatus string

const (
	TESTSTATUS_INITIALIZED            TestStatus = "initialized"
	TESTSTATUS_INITIALIZED_SUCCESSFUL TestStatus = "initialized+test-successful"
	TESTSTATUS_FAILED                 TestStatus = "test-failed"
)

// ConnMonitorConfig holds configuration for connection health testing and monitoring.
// It includes options for retries during tests and periodic health checks.
type ConnMonitorConfig struct {
	Ignore                 bool `json:"ignore,omitempty"`                 // If true, disables all testing retries and monitoring. Default: false.
	Retries                int  `json:"retries,omitempty"`                // Number of retry attempts for connection tests. Default: 3.
	RetryDelaySeconds      int  `json:"retryDelaySeconds,omitempty"`      // Delay between retry attempts in seconds. Default: 5.
	MonitorIntervalSeconds int  `json:"monitorIntervalSeconds,omitempty"` // Interval for health monitoring in seconds. Default: 300 (5 minutes).
}

// NewConnMonitorConfig returns a new ConnMonitorConfig instance with default values.
func NewConnMonitorConfig() *ConnMonitorConfig {
	return &ConnMonitorConfig{
		Ignore:                 false,
		Retries:                3,
		RetryDelaySeconds:      5,
		MonitorIntervalSeconds: 5 * 60,
	}
}

// IsEnabled returns true if monitoring and retries are active (i.e., not ignored).
func (c *ConnMonitorConfig) IsEnabled() bool {
	return !c.Ignore
}

// GetRetries returns the number of retry attempts.
// Returns 0 if disabled (ignored); defaults to 3 if zero.
func (c *ConnMonitorConfig) GetRetries() int {
	if !c.IsEnabled() {
		return 0
	}
	if c.Retries <= 0 {
		return 3 // Default fallback.
	}
	return c.Retries
}

// GetDelay returns the retry delay as a time.Duration.
// Returns 0 if disabled (ignored); defaults to 5s if zero.
func (c *ConnMonitorConfig) GetDelay() time.Duration {
	if !c.IsEnabled() {
		return 0
	}
	if c.RetryDelaySeconds <= 0 {
		return 5 * time.Second // Default fallback.
	}
	return time.Duration(c.RetryDelaySeconds) * time.Second
}

// GetInterval returns the monitoring interval as a time.Duration.
// Returns 0 if disabled (ignored); defaults to 5m if zero.
func (c *ConnMonitorConfig) GetInterval() time.Duration {
	if !c.IsEnabled() {
		return 0
	}
	if c.MonitorIntervalSeconds <= 0 {
		return 5 * time.Minute // Default fallback.
	}
	return time.Duration(c.MonitorIntervalSeconds) * time.Second
}

// Validate checks if the ConnMonitorConfig values are within acceptable ranges.
// If Ignore is true, skips validation and returns nil.
// Otherwise, ensures Retries >=1 and seconds fields >0.
// Returns an error if any value is invalid; otherwise, nil.
func (c *ConnMonitorConfig) Validate() error {
	if c.Ignore {
		return nil // Disabled; no further checks.
	}
	if c.Retries < 1 {
		return fmt.Errorf("retries must be at least 1 when not ignored")
	}
	if c.RetryDelaySeconds <= 0 {
		return fmt.Errorf("retry delay must be positive when not ignored")
	}
	if c.MonitorIntervalSeconds <= 0 {
		return fmt.Errorf("monitor interval must be positive when not ignored")
	}
	return nil
}
