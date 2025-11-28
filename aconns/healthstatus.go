package aconns

import (
	"strings"
	"time"
)

// HealthStatus represents the current health state of a connection or adapter.
// It is used to categorize the outcome of validation, opening, pinging, or other operations.
type HealthStatus string

const (
	// HEALTHSTATUS_HEALTHY indicates the connection is fully operational and healthy.
	HEALTHSTATUS_HEALTHY HealthStatus = "healthy"

	// HEALTHSTATUS_UNKNOWN indicates the health status has not been determined yet.
	HEALTHSTATUS_UNKNOWN HealthStatus = "unknown"

	// HEALTHSTATUS_VALIDATE_FAILED indicates a failure during configuration validation.
	HEALTHSTATUS_VALIDATE_FAILED HealthStatus = "validate failed"

	// HEALTHSTATUS_OPEN_FAILED indicates a failure when attempting to open the connection.
	HEALTHSTATUS_OPEN_FAILED HealthStatus = "open failed"

	// HEALTHSTATUS_PING_FAILED indicates a failure during a ping or test query.
	HEALTHSTATUS_PING_FAILED HealthStatus = "ping failed"

	// HEALTHSTATUS_CLOSED indicates the connection has been explicitly closed.
	HEALTHSTATUS_CLOSED HealthStatus = "closed"

	// HEALTHSTATUS_TIMEOUT indicates a timeout occurred during an operation.
	HEALTHSTATUS_TIMEOUT HealthStatus = "timeout"

	// HEALTHSTATUS_NETWORK_ERROR indicates a network-related error, such as connectivity issues.
	HEALTHSTATUS_NETWORK_ERROR HealthStatus = "network error"

	// HEALTHSTATUS_AUTH_FAILED indicates authentication or credential validation failed.
	HEALTHSTATUS_AUTH_FAILED HealthStatus = "auth failed"
)

// IsEmpty returns true if the HealthStatus is empty or contains only whitespace.
func (hs HealthStatus) IsEmpty() bool {
	return strings.TrimSpace(string(hs)) == ""
}

// IsOK returns true if the HealthStatus is HEALTHSTATUS_HEALTHY.
func (hs HealthStatus) IsOK() bool {
	return hs == HEALTHSTATUS_HEALTHY
}

// IsFailed returns true if the HealthStatus indicates any failure state.
func (hs HealthStatus) IsFailed() bool {
	switch hs {
	case HEALTHSTATUS_VALIDATE_FAILED, HEALTHSTATUS_OPEN_FAILED, HEALTHSTATUS_PING_FAILED,
		HEALTHSTATUS_TIMEOUT, HEALTHSTATUS_NETWORK_ERROR, HEALTHSTATUS_AUTH_FAILED:
		return true
	default:
		return false
	}
}

// String returns the string representation of the HealthStatus.
func (hs HealthStatus) String() string {
	return string(hs)
}

// HealthCheck encapsulates the health information of a connection or adapter,
// including whether it is currently healthy, the timestamp of the last check,
// and the status from the last operation.
type HealthCheck struct {
	IsHealthy  bool         `json:"isHealthy,omitempty"`
	LastCheck  time.Time    `json:"lastCheck,omitempty"`
	LastStatus HealthStatus `json:"lastStatus,omitempty"`
}

// NewHealthCheck creates a new HealthCheck instance with the given initial status.
// If the provided status is empty, it defaults to HEALTHSTATUS_UNKNOWN.
// The IsHealthy field is set to true only if the status is HEALTHSTATUS_HEALTHY.
// The LastCheck timestamp is initialized to the zero time value.
func NewHealthCheck(lastStatus HealthStatus) *HealthCheck {
	if lastStatus.IsEmpty() {
		lastStatus = HEALTHSTATUS_UNKNOWN
	}
	return &HealthCheck{
		IsHealthy:  lastStatus == HEALTHSTATUS_HEALTHY,
		LastCheck:  time.Now(),
		LastStatus: lastStatus,
	}
}

// Update refreshes the HealthCheck with a new status and sets the current timestamp as LastCheck.
// It updates IsHealthy based on whether the new status is HEALTHSTATUS_HEALTHY.
func (hc *HealthCheck) Update(newStatus HealthStatus) {
	if hc == nil {
		return
	}
	hc.LastStatus = newStatus
	hc.IsHealthy = newStatus == HEALTHSTATUS_HEALTHY
	hc.LastCheck = time.Now()
}

// IsStale returns true if the last health check is older than the given duration.
// If LastCheck is zero (never checked), it returns true.
func (hc *HealthCheck) IsStale(maxAge time.Duration) bool {
	if hc == nil || hc.LastCheck.IsZero() {
		return true
	}
	return time.Since(hc.LastCheck) > maxAge
}

// IsFailed returns true if the last check failed.
func (hc *HealthCheck) IsFailed() bool {
	if hc == nil {
		return true
	}
	return hc.LastStatus.IsFailed()
}
