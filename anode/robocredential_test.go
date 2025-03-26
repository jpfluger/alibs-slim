package anode

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/jpfluger/alibs-slim/acrypt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRoboCredential_GenerateKeyPair(t *testing.T) {
	rc := &RoboCredential{}
	masterPassword := "secure-password"

	err := rc.GenerateKeyPair(masterPassword, 0)
	assert.NoError(t, err)
	assert.NotEmpty(t, rc.PublicKey, "Public key should not be empty")

	privKey, err := rc.GetDecodedPrivateKey(masterPassword)
	assert.NoError(t, err)
	assert.NotEmpty(t, privKey, "Private key should not be empty")
}

func TestRoboCredential_RotateKeys(t *testing.T) {
	rc := &RoboCredential{}
	masterPassword := "secure-password"

	err := rc.GenerateKeyPair(masterPassword, 0)
	assert.NoError(t, err)

	oldPublicKey := rc.PublicKey
	oldPrivKey, _ := rc.GetDecodedPrivateKey(masterPassword)

	err = rc.RotateKeys(masterPassword)
	assert.NoError(t, err)

	newPublicKey := rc.PublicKey
	newPrivKey, _ := rc.GetDecodedPrivateKey(masterPassword)

	assert.NotEqual(t, oldPublicKey, newPublicKey, "Public key should change after rotation")
	assert.NotEqual(t, acrypt.EncodeToBase64(oldPrivKey), acrypt.EncodeToBase64(newPrivKey), "Private key should change after rotation")
	assert.Empty(t, rc.AccessToken, "Access token should be invalidated after rotation")
	assert.Empty(t, rc.RefreshToken, "Refresh token should be invalidated after rotation")
}

func TestRoboCredential_ValidateToken(t *testing.T) {
	rc := &RoboCredential{}
	rc.TokenExpiresAt = time.Now().Add(1 * time.Hour)
	rc.IsTokenValid = true

	assert.True(t, rc.ValidateToken(), "Token should be valid")

	rc.TokenExpiresAt = time.Now().Add(-1 * time.Hour)
	assert.False(t, rc.ValidateToken(), "Token should be invalid after expiration")
}

func TestRoboCredential_InvalidateTokens(t *testing.T) {
	rc := &RoboCredential{}
	rc.AccessToken = "some-token"
	rc.RefreshToken = "some-refresh-token"
	rc.TokenExpiresAt = time.Now().Add(1 * time.Hour)
	rc.RefreshExpiresAt = time.Now().Add(24 * time.Hour)
	rc.IsTokenValid = true

	rc.InvalidateTokens()

	assert.Empty(t, rc.AccessToken, "Access token should be invalidated")
	assert.Empty(t, rc.RefreshToken, "Refresh token should be invalidated")
	assert.False(t, rc.IsTokenValid, "Token validity should be set to false")
	assert.True(t, rc.TokenExpiresAt.IsZero(), "Token expiration time should be reset")
	assert.True(t, rc.RefreshExpiresAt.IsZero(), "Refresh token expiration time should be reset")
}

func TestRoboCredential_RefreshAccessToken(t *testing.T) {
	rc := &RoboCredential{}
	rc.RefreshToken = "valid-refresh-token"
	rc.RefreshExpiresAt = time.Now().Add(1 * time.Hour)

	refreshFunc := func(refreshToken string) (string, time.Time, error) {
		if refreshToken != "valid-refresh-token" {
			return "", time.Time{}, fmt.Errorf("invalid refresh token")
		}
		return "new-access-token", time.Now().Add(1 * time.Hour), nil
	}

	err := rc.RefreshAccessToken(refreshFunc)
	assert.NoError(t, err)
	assert.Equal(t, "new-access-token", rc.AccessToken, "Access token should be updated")
	assert.True(t, rc.TokenExpiresAt.After(time.Now()), "Token expiration time should be in the future")
	assert.True(t, rc.IsTokenValid, "Token validity should be set to true")
}

// Peer represents a single node in the peer-to-peer network.
type Peer struct {
	Name          string
	Credential    *RoboCredential
	SharedSecrets map[string]string // Map to store shared secrets (e.g., messages)
	NonceStore    *acrypt.NonceStore
	Mutex         sync.RWMutex // Mutex for thread-safe access
}

func NewPeer(name string, masterPassword string) (*Peer, error) {
	cred := &RoboCredential{}
	if err := cred.GenerateKeyPair(masterPassword, 0); err != nil {
		return nil, fmt.Errorf("failed to generate key pair for peer %s: %v", name, err)
	}

	return &Peer{
		Name:          name,
		Credential:    cred,
		SharedSecrets: make(map[string]string),
		NonceStore:    acrypt.NewNonceStore(),
	}, nil
}

func (p *Peer) AuthenticatePeer(payload string, signedPayload acrypt.SignedPayload, base64PublicKey string) error {
	// Decode the base64 public key
	pubKeyBytes, err := base64.StdEncoding.DecodeString(base64PublicKey)
	if err != nil {
		return fmt.Errorf("failed to decode public key: %v", err)
	}

	// Verify the signature
	isValid, err := signedPayload.Verify(payload, ed25519.PublicKey(pubKeyBytes))
	if err != nil || !isValid {
		return fmt.Errorf("authentication failed: %v", err)
	}

	// Check nonce for replay protection
	if !p.NonceStore.Add(signedPayload.Nonce) {
		return fmt.Errorf("replay attack detected: nonce %s already used", signedPayload.Nonce)
	}

	return nil
}

func (p *Peer) SendMessage(receiver *Peer, masterPassword, message string) error {
	// Sign the message using the sender's private key
	privKey, err := p.Credential.GetDecodedPrivateKey(masterPassword)
	if err != nil {
		return fmt.Errorf("failed to retrieve private key: %v", err)
	}

	signedPayload := &acrypt.SignedPayload{}
	signedMsg, err := signedPayload.Sign(message, ed25519.PrivateKey(privKey))
	if err != nil {
		return fmt.Errorf("failed to sign message: %v", err)
	}

	// Send the base64-encoded public key
	err = receiver.AuthenticatePeer(message, signedMsg, string(p.Credential.PublicKey))
	if err != nil {
		return fmt.Errorf("receiver failed to authenticate message: %v", err)
	}

	// Store the message in the receiver's shared secrets
	receiver.Mutex.Lock()
	defer receiver.Mutex.Unlock()
	receiver.SharedSecrets[p.Name] = message

	return nil
}

func (p *Peer) EncryptMessage(message, sharedKey string) (string, error) {
	// Derive AES key
	key := sha256.Sum256([]byte(sharedKey))
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return "", fmt.Errorf("failed to create AES cipher: %v", err)
	}

	nonce := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("failed to generate nonce: %v", err)
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM cipher: %v", err)
	}

	ciphertext := aesGCM.Seal(nonce, nonce, []byte(message), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func (p *Peer) DecryptMessage(encryptedMessage, sharedKey string) (string, error) {
	// Derive AES key
	key := sha256.Sum256([]byte(sharedKey))
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return "", fmt.Errorf("failed to create AES cipher: %v", err)
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM cipher: %v", err)
	}

	decodedMessage, err := base64.StdEncoding.DecodeString(encryptedMessage)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64 message: %v", err)
	}

	nonce, ciphertext := decodedMessage[:12], decodedMessage[12:]
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt message: %v", err)
	}

	return string(plaintext), nil
}

func (p *Peer) SendMessageWithUpdatedKey(receiver *Peer, masterPassword, message string, updatedPublicKey acrypt.CryptKeyBase64) error {
	// Update the receiver's public key
	receiver.Credential.PublicKey = updatedPublicKey

	// Reuse the original SendMessage logic
	return p.SendMessage(receiver, masterPassword, message)
}

func TestRoboCredential_PeerToPeerMessaging(t *testing.T) {
	masterPassword := "secure-pass"

	// Create two peers
	peer1, err := NewPeer("Peer1", masterPassword)
	assert.NoError(t, err, "Peer1 should be created successfully")

	peer2, err := NewPeer("Peer2", masterPassword)
	assert.NoError(t, err, "Peer2 should be created successfully")

	// Debugging logs
	t.Logf("Peer1 Public Key (Base64): %s", peer1.Credential.PublicKey.Encoded())
	t.Logf("Peer2 Public Key (Base64): %s", peer2.Credential.PublicKey.Encoded())

	// Peer1 sends a message to Peer2
	message := "Hello from Peer1"
	err = peer1.SendMessage(peer2, masterPassword, message)
	assert.NoError(t, err, "Message should be sent successfully")

	// Verify the message was received
	peer2.Mutex.RLock()
	receivedMessage, exists := peer2.SharedSecrets[peer1.Name]
	peer2.Mutex.RUnlock()

	assert.True(t, exists, "Message should be received by Peer2")
	assert.Equal(t, message, receivedMessage, "Received message should match the sent message")
}

func TestRoboCredential_UnifiedAuthenticatePeer(t *testing.T) {
	// Setup
	peer1, err := NewPeer("Peer1", "secure-pass-1")
	assert.NoError(t, err)

	peer2, err := NewPeer("Peer2", "secure-pass-2")
	assert.NoError(t, err)

	// Simulate signed payload
	payload := "Hello, Peer2!"
	privKey, _ := peer1.Credential.GetDecodedPrivateKey("secure-pass-1")
	signedPayload := acrypt.SignedPayload{}
	signedMsg, _ := signedPayload.Sign(payload, ed25519.PrivateKey(privKey))

	// Authenticate
	err = peer2.AuthenticatePeer(payload, signedMsg, string(peer1.Credential.PublicKey))
	assert.NoError(t, err, "Authentication should succeed")
}

func TestRoboCredential_TokenExpiration(t *testing.T) {
	peer, err := NewPeer("TestPeer", "secure-pass")
	assert.NoError(t, err)

	// Simulate an expired token
	peer.Credential.TokenExpiresAt = time.Now().Add(-1 * time.Minute)
	assert.False(t, peer.Credential.ValidateToken(), "Token should be invalid after expiration")
}

func TestRoboCredential_KeyRotationDuringCommunication(t *testing.T) {
	masterPassword := "secure-pass"

	// Create peers
	peer1, err := NewPeer("Peer1", masterPassword)
	assert.NoError(t, err)
	assert.NotNil(t, peer1, "Peer1 should be initialized")
	assert.NotNil(t, peer1.Credential, "Peer1.Credential should be initialized")

	peer2, err := NewPeer("Peer2", masterPassword)
	assert.NoError(t, err)
	assert.NotNil(t, peer2, "Peer2 should be initialized")
	assert.NotNil(t, peer2.Credential, "Peer2.Credential should be initialized")

	// Debugging logs
	t.Logf("Peer1 Public Key (Base64): %s", peer1.Credential.PublicKey)
	t.Logf("Peer2 Public Key (Base64): %s", peer2.Credential.PublicKey)

	// Send a message before rotation
	messageBefore := "Message before key rotation"
	err = peer1.SendMessage(peer2, masterPassword, messageBefore)
	assert.NoError(t, err, "Message should be sent successfully before rotation")

	// Rotate keys on Peer2
	err = peer2.Credential.RotateKeys(masterPassword)
	assert.NoError(t, err, "Key rotation should succeed")
	t.Logf("Peer2 New Public Key (Base64): %s", peer2.Credential.PublicKey)

	// Simulate key sharing: update Peer1's reference to Peer2's public key
	peer1SharedPublicKey := peer2.Credential.PublicKey

	// Try to send a message after rotation
	messageAfter := "Message after key rotation"
	err = peer1.SendMessageWithUpdatedKey(peer2, masterPassword, messageAfter, peer1SharedPublicKey)

	// Expect success because the key has been updated
	assert.NoError(t, err, "Message should succeed after key sharing")
}

type PeerServer struct {
	Name          string
	Credential    *RoboCredential
	SharedKeys    map[string]acrypt.CryptKeyBase64
	NonceStore    *acrypt.NonceStore
	Server        *echo.Echo
	JWTSigningKey []byte
	Mutex         sync.RWMutex
}

func NewPeerServer(name, masterPassword string) (*PeerServer, error) {
	cred := &RoboCredential{}
	if err := cred.GenerateKeyPair(masterPassword, 0); err != nil {
		return nil, fmt.Errorf("failed to generate key pair for peer %s: %v", name, err)
	}

	server := echo.New()
	peer := &PeerServer{
		Name:          name,
		Credential:    cred,
		SharedKeys:    make(map[string]acrypt.CryptKeyBase64),
		NonceStore:    acrypt.NewNonceStore(),
		JWTSigningKey: []byte("super-secret-key"),
		Server:        server,
	}

	// Register API handlers: Signed key verify every transaction
	server.POST("/exchange-public-key", peer.ExchangePublicKeyHandler)
	server.POST("/send-message", peer.ReceiveMessageHandler)

	// Register API handlers: Signed key verify for initial auth, then JWT thereafter
	server.POST("/authenticate", peer.AuthenticateHandler)
	server.POST("/secure-endpoint", peer.SecureEndpointHandler)

	return peer, nil
}

func (p *PeerServer) ExchangePublicKeyHandler(c echo.Context) error {
	var request struct {
		PeerName  string                `json:"peerName"`
		PublicKey acrypt.CryptKeyBase64 `json:"publicKey"`
	}

	// Parse the incoming request
	if err := c.Bind(&request); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request payload")
	}

	// Decode the base64-encoded public key
	senderPublicKey, err := request.PublicKey.Decoded()
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid public key: base64 decoding failed")
	}

	// Validate the public key size
	if len(senderPublicKey) != ed25519.PublicKeySize {
		return echo.NewHTTPError(http.StatusBadRequest,
			fmt.Sprintf("Invalid public key size: expected %d bytes, got %d bytes", ed25519.PublicKeySize, len(senderPublicKey)))
	}

	// Store the valid base64-encoded public key
	p.Mutex.Lock()
	p.SharedKeys[request.PeerName] = request.PublicKey
	p.Mutex.Unlock()

	return c.JSON(http.StatusOK, map[string]string{"message": "Public key exchanged successfully"})
}

func (p *PeerServer) ReceiveMessageHandler(c echo.Context) error {
	var request struct {
		SenderName    string               `json:"senderName"`
		Message       string               `json:"message"`
		SignedPayload acrypt.SignedPayload `json:"signedPayload"`
	}
	if err := c.Bind(&request); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request payload")
	}

	// Retrieve the sender's public key
	p.Mutex.RLock()
	senderPublicKeyBase64, ok := p.SharedKeys[request.SenderName]
	p.Mutex.RUnlock()
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "Sender's public key not found")
	}

	senderPublicKey, err := senderPublicKeyBase64.Decoded()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to decode sender's public key")
	}

	// Authenticate the sender
	err = p.AuthenticatePeer(request.Message, request.SignedPayload, senderPublicKey)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, fmt.Sprintf("Failed to authenticate sender: %v", err))
	}

	// Store the message (for simplicity, just print it here)
	c.Logger().Infof("Message received from %s: %s", request.SenderName, request.Message)
	return c.JSON(http.StatusOK, map[string]string{"message": "Message received successfully"})
}

func (p *PeerServer) AuthenticatePeer(payload string, signedPayload acrypt.SignedPayload, pubKey []byte) error {
	// Verify the signature
	isValid, err := signedPayload.Verify(payload, ed25519.PublicKey(pubKey))
	if err != nil || !isValid {
		return fmt.Errorf("authentication failed: %v", err)
	}

	// Check nonce for replay protection
	if !p.NonceStore.Add(signedPayload.Nonce) {
		return fmt.Errorf("replay attack detected: nonce %s already used", signedPayload.Nonce)
	}

	return nil
}

func (p *PeerServer) AuthenticateHandler(c echo.Context) error {
	var request struct {
		SenderName    string               `json:"senderName"`
		Message       string               `json:"message"`
		SignedPayload acrypt.SignedPayload `json:"signedPayload"`
	}
	if err := c.Bind(&request); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request payload")
	}

	// Validate public key
	p.Mutex.RLock()
	senderPublicKeyBase64, ok := p.SharedKeys[request.SenderName]
	p.Mutex.RUnlock()
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "Sender's public key not found")
	}

	senderPublicKey, err := senderPublicKeyBase64.Decoded()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to decode sender's public key")
	}

	// Authenticate the sender
	isValid, err := request.SignedPayload.Verify(request.Message, ed25519.PublicKey(senderPublicKey))
	if err != nil || !isValid {
		return echo.NewHTTPError(http.StatusUnauthorized, "Failed to authenticate sender")
	}

	// Issue JWT
	expiration := time.Now().Add(10 * time.Minute) // Short-lived token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"senderName": request.SenderName,
		"exp":        expiration.Unix(),
	})

	tokenString, err := token.SignedString(p.JWTSigningKey)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to issue token")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"accessToken": tokenString,
		"expiresAt":   expiration,
	})
}

func (p *PeerServer) SecureEndpointHandler(c echo.Context) error {
	tokenString := c.Request().Header.Get("Authorization")
	if tokenString == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, "Missing Authorization header")
	}

	tokenString = strings.TrimPrefix(tokenString, "Bearer ")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Method)
		}
		return p.JWTSigningKey, nil
	})

	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid or expired token: %v", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Process the request (e.g., validate claims, perform actions)
		return c.JSON(http.StatusOK, map[string]string{"message": "Access granted", "sender": claims["senderName"].(string)})
	}

	return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token")
}

func TestRoboCredential_ExchangePublicKeyHandlerValidation(t *testing.T) {
	peer, err := NewPeerServer("Peer1", "secure-pass")
	assert.NoError(t, err)

	server := httptest.NewServer(peer.Server)
	defer server.Close()

	tests := []struct {
		name         string
		publicKey    string
		expectedCode int
		expectedMsg  string
	}{
		{
			name:         "Valid public key",
			publicKey:    base64.StdEncoding.EncodeToString(ed25519.NewKeyFromSeed(make([]byte, 32)).Public().(ed25519.PublicKey)),
			expectedCode: http.StatusOK,
			expectedMsg:  "Public key exchanged successfully",
		},
		{
			name:         "Valid generated public key",
			publicKey:    peer.Credential.PublicKey.Encoded(),
			expectedCode: http.StatusOK,
			expectedMsg:  "Public key exchanged successfully",
		},
		{
			name:         "Invalid base64",
			publicKey:    "invalid-base64@@@",
			expectedCode: http.StatusBadRequest,
			expectedMsg:  "Invalid public key: base64 decoding failed",
		},
		{
			name:         "Invalid key size",
			publicKey:    base64.StdEncoding.EncodeToString([]byte("short-key")),
			expectedCode: http.StatusBadRequest,
			expectedMsg:  fmt.Sprintf("Invalid public key size: expected %d bytes, got %d bytes", ed25519.PublicKeySize, len([]byte("short-key"))),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payload := fmt.Sprintf(`{"peerName":"Peer2","publicKey":"%s"}`, tt.publicKey)
			resp, err := http.Post(server.URL+"/exchange-public-key", "application/json", strings.NewReader(payload))
			assert.NoError(t, err)

			body, _ := io.ReadAll(resp.Body)
			defer resp.Body.Close()

			assert.Equal(t, tt.expectedCode, resp.StatusCode)
			if tt.expectedCode == http.StatusOK {
				assert.Contains(t, string(body), tt.expectedMsg)
			} else {
				assert.Contains(t, string(body), tt.expectedMsg)
			}
		})
	}
}

func TestRoboCredential_AuthenticatePeerInvalidKeySize(t *testing.T) {
	peer, err := NewPeerServer("Peer1", "secure-pass")
	assert.NoError(t, err)

	payload := "test-message"
	signedPayload := acrypt.SignedPayload{}

	// Use an invalid base64-encoded public key
	invalidKey := base64.StdEncoding.EncodeToString([]byte("invalid-key-size"))

	err = peer.AuthenticatePeer(payload, signedPayload, []byte(invalidKey))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid public key size", "Should fail due to invalid key size")
}

func TestRoboCredential_PeerToPeerAPIWithRoboCredential(t *testing.T) {
	masterPassword1 := "secure-pass-1"
	masterPassword2 := "secure-pass-2"

	// Create peers and their servers
	peer1, err := NewPeerServer("Peer1", masterPassword1)
	assert.NoError(t, err)
	ts1 := httptest.NewServer(peer1.Server)
	defer ts1.Close()

	peer2, err := NewPeerServer("Peer2", masterPassword2)
	assert.NoError(t, err)
	ts2 := httptest.NewServer(peer2.Server)
	defer ts2.Close()

	// Exchange public keys
	resp, err := http.Post(ts2.URL+"/exchange-public-key", "application/json", strings.NewReader(`{
		"peerName": "Peer1",
		"publicKey": "`+peer1.Credential.PublicKey.Encoded()+`"
	}`))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	resp, err = http.Post(ts1.URL+"/exchange-public-key", "application/json", strings.NewReader(`{
		"peerName": "Peer2",
		"publicKey": "`+peer2.Credential.PublicKey.Encoded()+`"
	}`))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Peer1 sends a message to Peer2
	message := "Hello from Peer1"
	privKey, _ := peer1.Credential.GetDecodedPrivateKey(masterPassword1)
	signedPayload := &acrypt.SignedPayload{}
	signedMsg, _ := signedPayload.Sign(message, ed25519.PrivateKey(privKey))

	payload, _ := json.Marshal(map[string]interface{}{
		"senderName":    "Peer1",
		"message":       message,
		"signedPayload": signedMsg,
	})

	resp, err = http.Post(ts2.URL+"/send-message", "application/json", strings.NewReader(string(payload)))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestRoboCredential_HybridAuthentication(t *testing.T) {
	masterPassword := "secure-pass"

	// Setup Peer1 (client)
	peer1, err := NewPeer("Peer1", masterPassword)
	assert.NoError(t, err)

	// Setup Peer2 (server)
	peer2, err := NewPeerServer("Peer2", masterPassword)
	assert.NoError(t, err)

	ts := httptest.NewServer(peer2.Server)
	defer ts.Close()

	// Exchange public keys
	resp, err := http.Post(ts.URL+"/exchange-public-key", "application/json", strings.NewReader(`{
		"peerName": "Peer1",
		"publicKey": "`+peer1.Credential.PublicKey.Encoded()+`"
	}`))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	resp, err = http.Post(ts.URL+"/exchange-public-key", "application/json", strings.NewReader(`{
		"peerName": "Peer2",
		"publicKey": "`+peer2.Credential.PublicKey.Encoded()+`"
	}`))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Peer1 sends initial signed payload
	initialPayload := "First request to authenticate"
	privKey, _ := peer1.Credential.GetDecodedPrivateKey(masterPassword)
	signedPayload := &acrypt.SignedPayload{}
	signedMsg, _ := signedPayload.Sign(initialPayload, ed25519.PrivateKey(privKey))

	requestBody := map[string]interface{}{
		"senderName":    "Peer1",
		"message":       initialPayload,
		"signedPayload": signedMsg,
	}

	body, _ := json.Marshal(requestBody)
	resp, err = http.Post(ts.URL+"/authenticate", "application/json", strings.NewReader(string(body)))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Extract JWT from the response
	var response map[string]interface{}
	_ = json.NewDecoder(resp.Body).Decode(&response)
	jwtToken := response["accessToken"].(string)
	assert.NotEmpty(t, jwtToken, "JWT token should be issued")

	// Peer1 uses JWT for subsequent requests
	req, _ := http.NewRequest(http.MethodPost, ts.URL+"/secure-endpoint", strings.NewReader(`{"data":"secure request"}`))
	req.Header.Set("Authorization", "Bearer "+jwtToken)
	resp, err = http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Verify the server's response
	var secureResponse map[string]string
	_ = json.NewDecoder(resp.Body).Decode(&secureResponse)
	assert.Equal(t, "Access granted", secureResponse["message"])
	assert.Equal(t, "Peer1", secureResponse["sender"])
}
