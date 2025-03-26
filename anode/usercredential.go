package anode

import (
	"encoding/base64"
	"fmt"
	"github.com/jpfluger/alibs-slim/acrypt"
	"github.com/jpfluger/alibs-slim/anetwork"
	"github.com/jpfluger/alibs-slim/auser"
	"strings"
	"sync"
)

// UserCredential represents a basic node user credential with tokens.
type UserCredential struct {
	JWTUserToken
	Username auser.Username `json:"username,omitempty"`
	Password string         `json:"password,omitempty"` // Password is a hashed password using Argon2id.

	// AuthName is optionally used as the name of the
	// connection adapter, link or id used for authentication.
	AuthName string `json:"authName"`

	WhitelistIPs []string `json:"whitelistIPs,omitempty"`
	parsedIPs    anetwork.NetIPs
	mu           sync.RWMutex // Added mutex for thread-safety
}

// GetHasCredential returns if one credential form is valid.
func (uc *UserCredential) GetHasCredential() (isJWT, isPassword, ok bool) {
	isJWT = uc.GetHasJWT()
	isPassword = uc.GetHasPassword()
	ok = isJWT || isPassword
	return
}

// GetHasJWT returns if this credential uses JWT.
func (uc *UserCredential) GetHasJWT() bool {
	return uc.GetIsTokenValid()
}

// GetHasPassword returns if this credential uses passwords.
func (uc *UserCredential) GetHasPassword() bool {
	uc.mu.RLock()
	defer uc.mu.RUnlock()
	return uc.Password != ""
}

// GetUsername returns the username.
func (uc *UserCredential) GetUsername() auser.Username {
	uc.mu.RLock()
	defer uc.mu.RUnlock()
	return uc.Username
}

// GetWhitelistIPs returns the access white-listed IPs.
func (uc *UserCredential) GetWhitelistIPs() []string {
	uc.mu.RLock()
	defer uc.mu.RUnlock()
	return uc.WhitelistIPs
}

// MatchUsernamePassword checks if the provided username and password match the stored ones.
// By default case sensitivity for username matches is false.
func (uc *UserCredential) MatchUsernamePassword(username auser.Username, password string) bool {
	return uc.MatchUsernamePasswordWithCaseSensitive(username, password, false)
}

// MatchUsernamePasswordWithCaseSensitive checks if the provided username and password match the stored ones with an optional case-sensitivity modifier for the username.
func (uc *UserCredential) MatchUsernamePasswordWithCaseSensitive(username auser.Username, password string, isUserCaseSensitive bool) bool {
	uc.mu.RLock()
	defer uc.mu.RUnlock()
	if username == "" || password == "" {
		return false
	}
	if isUserCaseSensitive {
		if uc.Username != username {
			return false
		}
	} else {
		if uc.Username.ToStringTrimLower() != username.ToStringTrimLower() {
			return false
		}
	}
	isValid, err := acrypt.IsEqualArgon2id(uc.Password, password)
	if !isValid || err != nil {
		return false
	}
	return true
}

// CheckAuthorizationHeaderWithSecretKey processes the Authorization header for Basic or Bearer token authentication.
func (uc *UserCredential) CheckAuthorizationHeaderWithSecretKey(authHeader string, secretKey []byte) (bool, error) {
	if strings.HasPrefix(authHeader, "Basic ") {
		payload, _ := base64.StdEncoding.DecodeString(strings.TrimPrefix(authHeader, "Basic "))
		pair := strings.SplitN(string(payload), ":", 2)
		if len(pair) != 2 || !uc.MatchUsernamePassword(auser.Username(pair[0]), pair[1]) {
			return false, fmt.Errorf("unauthorized")
		}
	} else if strings.HasPrefix(authHeader, "Bearer ") {
		token := strings.TrimPrefix(authHeader, "Bearer ")
		isValid, _, err := uc.VerifyTokenTimeWithSecretKey(token, secretKey)
		if err != nil || !isValid {
			return false, fmt.Errorf("unauthorized")
		}
	} else {
		return false, fmt.Errorf("unauthorized")
	}
	return true, nil
}

// CheckClientIP checks if the client's IP is in the whitelist.
func (uc *UserCredential) CheckClientIP(clientIP string) bool {
	hasParsedIPs, containsIP := uc.WhitelistIPContains(clientIP)
	return !hasParsedIPs || containsIP
}

// WhitelistIPContains checks if the ipTarget is in the list of parsed IPs or subnets.
func (uc *UserCredential) WhitelistIPContains(ipTarget string) (hasParsedIPs, containsIP bool) {
	uc.mu.RLock()
	defer uc.mu.RUnlock()
	if uc.parsedIPs == nil || len(uc.parsedIPs) == 0 {
		return false, false
	}
	return true, uc.parsedIPs.Contains(ipTarget)
}

func (uc *UserCredential) GetParsedWhitelistIPs() anetwork.NetIPs {
	uc.mu.RLock()
	defer uc.mu.RUnlock()
	return uc.parsedIPs
}

// EncryptPassword hashes the secret using Argon2id and stores it in the Password field.
func (uc *UserCredential) EncryptPassword(secret string) error {
	uc.mu.Lock()
	defer uc.mu.Unlock()
	if acrypt.IsArgon2idHash(secret) {
		return fmt.Errorf("secret already hashed")
	}
	hashedPassword, err := acrypt.EncryptArgon2id(secret, nil)
	if err != nil {
		return fmt.Errorf("failed to encrypt secret: %v", err)
	}
	uc.Password = hashedPassword
	uc.ClearToken()
	return nil
}

// CheckPassword compares a plaintext secret with the hashed secret stored in the Password field.
func (uc *UserCredential) CheckPassword(secret string) (bool, error) {
	uc.mu.RLock()
	defer uc.mu.RUnlock()

	isEqual, err := acrypt.IsEqualArgon2id(uc.Password, secret)
	if err != nil {
		return false, fmt.Errorf("failed to check secret: %v", err)
	}

	return isEqual, nil
}

// ValidateWithSecretKey validates the UserCredential.
func (uc *UserCredential) ValidateWithSecretKey(secretKey []byte) error {
	uc.mu.Lock()
	defer uc.mu.Unlock()

	uc.AuthName = strings.TrimSpace(uc.AuthName)

	uc.Username = uc.Username.TrimSpace()
	uc.Password = strings.TrimSpace(uc.Password)

	if uc.Username.IsEmpty() {
		return fmt.Errorf("username is empty")
	}
	//if uc.Password == "" {
	//	return fmt.Errorf("password is required")
	//}
	if uc.Password != "" && !acrypt.IsArgon2idHash(uc.Password) {
		return fmt.Errorf("invalid password; not hashed")
	}

	// Validate IPs
	if len(uc.WhitelistIPs) > 0 {
		netIPs, err := anetwork.ToNetIP(uc.WhitelistIPs)
		if err != nil {
			return fmt.Errorf("failed validating user whitelistIPs: %v", err)
		}
		uc.parsedIPs = netIPs
	}

	// Validate Tokens
	if uc.Password != "" {
		uc.ClearToken()
	} else {
		if err := uc.JWTUserToken.ValidateTokenWithSecretKey(secretKey); err != nil {
			return fmt.Errorf("failed validating JWTUserToken: %v", err)
		}
	}

	return nil
}

// GenerateTokensWithSecretKey generates an access token and a refresh token using the Credential's Username and passed-in secretKey.
func (uc *UserCredential) GenerateTokensWithSecretKey(accessExpiresInHours, refreshExpiresInDays int, secretKey []byte) (string, string, error) {
	accessToken, refreshToken, err := uc.GenerateTokensWithUsername(uc.GetUsername(), accessExpiresInHours, refreshExpiresInDays, secretKey)
	if err != nil {
		return "", "", err
	}
	uc.mu.Lock()
	defer uc.mu.Unlock()
	uc.Password = ""
	return accessToken, refreshToken, nil
}
