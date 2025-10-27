package auuids

import (
	"errors"
	"time"
)

// AuditInfo holds lifecycle tracking fields for entities.
type AuditInfo struct {
	CreatedAt time.Time `json:"createdAt,omitempty"`
	UpdatedAt time.Time `json:"updatedAt,omitempty"`
	CreatedBy UUID      `json:"createdBy,omitempty"` // Optional user ID
	UpdatedBy UUID      `json:"updatedBy,omitempty"` // Optional user ID
}

func NewAuditInfo(userID UUID) AuditInfo {
	au := AuditInfo{}
	au.CreatedBy = userID
	au.UpdatedBy = userID
	au.Normalize()
	return au
}

// Validate ensures AuditInfo is valid (e.g., timestamps logical, no future dates).
// Does not set defaults; only normalizes and checks.
func (ai *AuditInfo) Validate() error {
	now := time.Now().UTC()
	if !ai.CreatedAt.IsZero() && ai.CreatedAt.After(now) {
		return errors.New("createdAt cannot be in the future")
	}
	if !ai.UpdatedAt.IsZero() && ai.UpdatedAt.After(now) {
		return errors.New("updatedAt cannot be in the future")
	}
	if !ai.UpdatedAt.IsZero() && !ai.CreatedAt.IsZero() && ai.UpdatedAt.Before(ai.CreatedAt) {
		return errors.New("updatedAt cannot be before createdAt")
	}
	// Normalize to UTC if not already (optional; time.Time handles this)
	ai.CreatedAt = ai.CreatedAt.UTC()
	ai.UpdatedAt = ai.UpdatedAt.UTC()
	return nil
}

// SetCreatedNowIfZero sets CreatedAt to now (UTC) if zero, optionally with user.
func (ai *AuditInfo) SetCreatedNowIfZero(userID UUID) {
	if ai.CreatedAt.IsZero() {
		ai.CreatedAt = time.Now().UTC()
		ai.CreatedBy = userID
	}
}

// SetUpdatedNow sets UpdatedAt to now (UTC), optionally with user.
func (ai *AuditInfo) SetUpdatedNow(userID UUID) {
	ai.UpdatedAt = time.Now().UTC()
	ai.UpdatedBy = userID
}

// IsNew returns true if CreatedAt is zero (e.g., for creation flows).
func (ai *AuditInfo) IsNew() bool {
	return ai.CreatedAt.IsZero()
}

// Normalize normalizes the AuditInfo by setting timestamps to now (UTC) if zero.
// This is a soft operation with no errors; user IDs are not modified.
func (ai *AuditInfo) Normalize() {
	now := time.Now().UTC()
	if ai.CreatedAt.IsZero() {
		ai.CreatedAt = now
	}
	if ai.UpdatedAt.IsZero() {
		ai.UpdatedAt = now
	}
	// Normalize to UTC (already done via now, but explicit for safety)
	ai.CreatedAt = ai.CreatedAt.UTC()
	ai.UpdatedAt = ai.UpdatedAt.UTC()
}
