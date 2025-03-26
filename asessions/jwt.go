package asessions

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/jpfluger/alibs-slim/auser"
	"time"
)

// JWTClaim extends jwt.RegisteredClaims to include custom fields for a JWT.
type JWTClaim struct {
	jwt.RegisteredClaims                // Embedded type to provide standard claims like expiry time.
	Username             auser.Username `json:"username"` // The username associated with the JWT.
	RedirectUrl          string         `json:"redirect"` // URL to redirect after processing the JWT.
	Action               ActionKey      `json:"action"`   // The associated action with the JWT.
}

// JWTClaimMeta holds metadata for a JWT claim.
type JWTClaimMeta struct {
	Id          string         `json:"id,omitempty"`       // Unique identifier for the JWT claim.
	Username    auser.Username `json:"username,omitempty"` // Username associated with the JWT claim.
	RedirectUrl string         `json:"redirect,omitempty"` // Redirect URL associated with the JWT claim.
	Action      ActionKey      `json:"action,omitempty"`   // Action associated with the JWT claim.
}

// JWTAttemptId represents an individual attempt with a JWT.
type JWTAttemptId struct {
	Id           string `json:"jti"`           // JWT ID, a unique identifier for the attempt.
	ExpiresAt    int64  `json:"exp,omitempty"` // Expiration time of the JWT as Unix time.
	HasClicked   bool   `json:"hasClicked"`    // Indicates if the JWT has been clicked/used.
	IsSaved      bool   `json:"isSaved"`       // Indicates if the attempt has been saved after a successful post request.
	SaveAttempts int    `json:"saveAttempts"`  // Number of times a save has been attempted.
}

// JWTAttemptIds is a slice of pointers to JWTAttemptId.
type JWTAttemptIds []*JWTAttemptId

// JWTUserAttemptCounter tracks the attempts made by a user.
type JWTUserAttemptCounter struct {
	Username            auser.Username `json:"username"` // Username of the user making attempts.
	Attempts            JWTAttemptIds  `json:"attempts"` // List of attempts made by the user.
	hasUnclickedAttempt bool           // Indicates if there is an unclicked attempt.
}

// SyncAttempts updates the Attempts slice to remove expired attempts and count valid ones.
// SyncAttempts brings Attempts up-to-date with the permitted duration.
// It returns the number of attempts within the duration time-window.
// Example: if a user clicks the JWT thereby deactivating it and starts a new one,
// limit to 4 new attempts per a duration of 60 minutes. Retain the expired
// attempts for historical until the items are reached.
func (uac *JWTUserAttemptCounter) SyncAttempts(durationMinutes int) int {
	if uac.Attempts == nil {
		uac.Attempts = JWTAttemptIds{}
	}
	if len(uac.Attempts) == 0 {
		return 0
	}
	now := time.Now().Unix()
	future := time.Now().Add(time.Duration(durationMinutes) * time.Minute).Unix()
	newAttempts := JWTAttemptIds{}
	for _, attempt := range uac.Attempts {
		if now < attempt.ExpiresAt {
			newAttempts = append(newAttempts, attempt)
			if !attempt.HasClicked {
				uac.hasUnclickedAttempt = true
			}
		} else {
			newExpire := time.Unix(attempt.ExpiresAt, 0)
			if newExpire.Add(time.Duration(durationMinutes)*time.Minute).Unix() < future {
				newAttempts = append(newAttempts, attempt)
			}
		}
	}
	uac.Attempts = newAttempts
	return len(uac.Attempts)
}

// HasUnclickedAttempt checks if there is an unclicked attempt.
func (uac *JWTUserAttemptCounter) HasUnclickedAttempt() bool {
	return uac.hasUnclickedAttempt
}

// SetHasClicked updates the clicked status of an attempt by ID.
func (uac *JWTUserAttemptCounter) SetHasClicked(id string, hasClicked bool) {
	if uac.Attempts == nil {
		uac.Attempts = JWTAttemptIds{}
	}
	if len(uac.Attempts) == 0 {
		return
	}
	for _, attempt := range uac.Attempts {
		if attempt.Id == id {
			attempt.HasClicked = hasClicked
			return
		}
	}
}

// SetIsSaved updates the saved status of an attempt by ID.
func (uac *JWTUserAttemptCounter) SetIsSaved(id string, isSaved bool) {
	if uac.Attempts == nil {
		uac.Attempts = JWTAttemptIds{}
	}
	if len(uac.Attempts) == 0 {
		return
	}
	for _, attempt := range uac.Attempts {
		if attempt.Id == id {
			attempt.IsSaved = isSaved
			return
		}
	}
}

// IncrementSaveCount increases the save attempt count for an attempt by ID.
func (uac *JWTUserAttemptCounter) IncrementSaveCount(id string) {
	if uac.Attempts == nil {
		uac.Attempts = JWTAttemptIds{}
	}
	if len(uac.Attempts) == 0 {
		return
	}
	for _, attempt := range uac.Attempts {
		if attempt.Id == id {
			attempt.SaveAttempts++
			return
		}
	}
}

// FindAttempt locates an attempt by ID.
func (uac *JWTUserAttemptCounter) FindAttempt(id string) *JWTAttemptId {
	if uac.Attempts == nil || len(uac.Attempts) == 0 {
		return nil
	}
	for _, attempt := range uac.Attempts {
		if attempt.Id == id {
			return attempt
		}
	}
	return nil
}
