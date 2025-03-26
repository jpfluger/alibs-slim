package anode

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jpfluger/alibs-slim/auser"
	"strings"
	"sync"
	"time"
)

// NewClaimsInMinutes creates a new JWT claims map with the specified username and expiration time given in minutes from now.
func NewClaimsInMinutes(username auser.Username, expirationTimeMinutes int) jwt.MapClaims {
	if expirationTimeMinutes <= 0 {
		expirationTimeMinutes = 5
	}
	return NewClaims(username, time.Now().Add(time.Duration(expirationTimeMinutes)*time.Minute))
}

// NewClaims creates a new JWT claims map with the specified username and expiration time.
func NewClaims(username auser.Username, expirationTime time.Time) jwt.MapClaims {
	return jwt.MapClaims{
		"username": username.String(),
		"exp":      expirationTime.Unix(),
	}
}

func JWTMapClaimSignedString(claims jwt.MapClaims, secretKey []byte) (tokenString string, err error) {
	if claims == nil {
		return "", errors.New("claims is nil")
	}
	if secretKey == nil || len(secretKey) == 0 {
		return "", errors.New("secretKey is nil")
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

// JWTMapClaimValidate validates a JWT token and matches claims against the provided mapMatch.
// If the mapMatch values match the claims, it returns the username from the claims.
func JWTMapClaimValidate(tokenString string, secretKey []byte, mapMatch map[string]string) (username auser.Username, isExpired bool, err error) {
	// Parse and validate JWT token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Method)
		}
		return secretKey, nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			isExpired = true
			err = fmt.Errorf("token has expired: %w", err)
			return
		}
		err = fmt.Errorf("failed to parse token: %w", err)
		return
	}

	// Extract claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		err = fmt.Errorf("invalid or expired token")
		return
	}

	// Check if claims match the mapMatch criteria (optional)
	if mapMatch != nil && len(mapMatch) > 0 {
		for key, value := range mapMatch {
			claimValue, ok := claims[key]
			if !ok || claimValue != value {
				err = fmt.Errorf("claim mismatch for key '%s': expected '%s', got '%v'", key, value, claimValue)
				return
			}
		}
	}

	// Retrieve the username from claims
	usernameClaim, ok := claims["username"]
	if !ok {
		return "", false, fmt.Errorf("username claim is missing")
	}

	// Ensure the username is a string
	user, ok := usernameClaim.(string)
	if !ok {
		return "", false, fmt.Errorf("username claim is not a valid string")
	}

	username = auser.Username(user)
	if username.IsEmpty() {
		return "", false, fmt.Errorf("username is empty")
	}

	return username, false, nil
}

// JWTUserToken handles JWT-related functionality for user tokens.
type JWTUserToken struct {
	AccessToken    string `json:"accessToken,omitempty"`  // Access tokens are expected to be JWT.
	RefreshToken   string `json:"refreshToken,omitempty"` // Refresh tokens are expected to be JWT.
	isTokenValid   bool
	tokenExpires   time.Time
	refreshExpires time.Time
	mu             sync.RWMutex // Added mutex for thread-safety
}

// GetAccessToken returns the access token.
func (uc *UserCredential) GetAccessToken() string {
	uc.mu.RLock()
	defer uc.mu.RUnlock()
	return uc.AccessToken
}

// GetRefreshToken returns the refresh token.
func (uc *UserCredential) GetRefreshToken() string {
	uc.mu.RLock()
	defer uc.mu.RUnlock()
	return uc.RefreshToken
}

// ValidateTokenWithSecretKey validates the access and refresh tokens using the provided secret key.
func (jt *JWTUserToken) ValidateTokenWithSecretKey(secretKey []byte) error {
	jt.mu.Lock()
	defer jt.mu.Unlock()
	jt.AccessToken = strings.TrimSpace(jt.AccessToken)
	jt.RefreshToken = strings.TrimSpace(jt.RefreshToken)

	if jt.AccessToken != "" {
		isValid, expTime, _ := jt.verifyTokenTimeWithSecretKey(jt.AccessToken, secretKey)
		jt.isTokenValid = isValid
		if expTime == nil {
			jt.tokenExpires = time.Time{}
		} else {
			jt.tokenExpires = *expTime
		}
	}
	if jt.RefreshToken != "" {
		isValid, expTime, _ := jt.verifyTokenTimeWithSecretKey(jt.RefreshToken, secretKey)
		if !isValid || expTime == nil || expTime.IsZero() {
			jt.refreshExpires = time.Time{}
		} else {
			jt.refreshExpires = *expTime
		}
	}
	return nil
}

// ClearToken clears the access and refresh tokens.
func (jt *JWTUserToken) ClearToken() {
	jt.mu.Lock()
	defer jt.mu.Unlock()
	jt.AccessToken = ""
	jt.RefreshToken = ""
	jt.isTokenValid = false
	jt.tokenExpires = time.Time{}
	jt.refreshExpires = time.Time{}
}

// GenerateTokensWithUsername generates an access token and a refresh token for the user.
func (jt *JWTUserToken) GenerateTokensWithUsername(username auser.Username, accessExpiresInHours, refreshExpiresInDays int, secretKey []byte) (string, string, error) {
	jt.mu.Lock()
	defer jt.mu.Unlock()

	accessExpirationTime := time.Now().Add(time.Duration(accessExpiresInHours) * time.Hour)
	refreshExpirationTime := time.Now().Add(time.Duration(refreshExpiresInDays) * 24 * time.Hour)

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, NewClaims(username, accessExpirationTime))
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, NewClaims(username, refreshExpirationTime))

	accessTokenString, err := accessToken.SignedString(secretKey)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate access token: %v", err)
	}

	refreshTokenString, err := refreshToken.SignedString(secretKey)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate refresh token: %v", err)
	}

	jt.AccessToken = accessTokenString
	jt.RefreshToken = refreshTokenString
	jt.tokenExpires = accessExpirationTime
	jt.refreshExpires = refreshExpirationTime
	jt.isTokenValid = true
	return accessTokenString, refreshTokenString, nil
}

// VerifyTokenWithSecretKey verifies the JWT token stored in the AccessToken field.
func (jt *JWTUserToken) VerifyTokenWithSecretKey(secretKey []byte) (bool, error) {
	jt.mu.RLock()
	defer jt.mu.RUnlock()
	isValid, _, err := jt.verifyTokenTimeWithSecretKey(jt.AccessToken, secretKey)
	return isValid, err
}

// VerifyTokenTimeWithSecretKey verifies the JWT token and returns the expiration time.
func (jt *JWTUserToken) VerifyTokenTimeWithSecretKey(tokenString string, secretKey []byte) (bool, *time.Time, error) {
	jt.mu.RLock()
	defer jt.mu.RUnlock()
	return jt.verifyTokenTimeWithSecretKey(tokenString, secretKey)
}

// verifyTokenTimeWithSecretKey verifies the JWT token and returns the expiration time.
func (jt *JWTUserToken) verifyTokenTimeWithSecretKey(tokenString string, secretKey []byte) (bool, *time.Time, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return false, nil, fmt.Errorf("token has expired")
		}
		return false, nil, fmt.Errorf("failed to parse token: %v", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if exp, ok := claims["exp"].(float64); ok {
			expTime := time.Unix(int64(exp), 0)
			if expTime.Before(time.Now()) {
				return false, &expTime, fmt.Errorf("token has expired")
			}
			return true, &expTime, nil
		}
		return true, nil, fmt.Errorf("expiration time not found in token")
	}

	return false, nil, fmt.Errorf("invalid token")
}

// RefreshAccessTokenWithSecretKey generates a new access token using the refresh token.
func (jt *JWTUserToken) RefreshAccessTokenWithSecretKey(username auser.Username, secretKey []byte) (string, error) {
	jt.mu.Lock()
	defer jt.mu.Unlock()

	isValid, expTime, err := jt.verifyTokenTimeWithSecretKey(jt.RefreshToken, secretKey)
	if err != nil || !isValid || expTime.Before(time.Now()) {
		return "", fmt.Errorf("invalid or expired refresh token")
	}

	accessExpirationTime := time.Now().Add(1 * time.Hour) // Example: 1 hour expiration for access token
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, NewClaims(username, accessExpirationTime))

	accessTokenString, err := accessToken.SignedString(secretKey)
	if err != nil {
		return "", fmt.Errorf("failed to generate access token: %v", err)
	}

	jt.AccessToken = accessTokenString
	jt.tokenExpires = accessExpirationTime
	jt.isTokenValid = true
	return accessTokenString, nil
}

// GetTokenExpiration returns the token expiration time.
func (jt *JWTUserToken) GetTokenExpiration() time.Time {
	jt.mu.RLock()
	defer jt.mu.RUnlock()
	return jt.tokenExpires
}

// SetTokenExpiration sets the token expiration time.
func (jt *JWTUserToken) SetTokenExpiration(tokenExpires time.Time) {
	jt.mu.Lock()
	defer jt.mu.Unlock()
	jt.tokenExpires = tokenExpires
}

// GetIsTokenValid returns the token validity status.
func (jt *JWTUserToken) GetIsTokenValid() bool {
	jt.mu.RLock()
	defer jt.mu.RUnlock()
	return jt.isTokenValid
}

// GetRefreshTokenExpiration returns the refresh token expiration time.
func (jt *JWTUserToken) GetRefreshTokenExpiration() time.Time {
	jt.mu.RLock()
	defer jt.mu.RUnlock()
	return jt.refreshExpires
}
