package aclient_duo

import (
	"encoding/json"
	"fmt"
	"net/url"
	"sync"
	"time"

	duoapi "github.com/duosecurity/duo_api_golang"
	"github.com/jpfluger/alibs-slim/aconns"
	"github.com/jpfluger/alibs-slim/anetwork"
	"github.com/jpfluger/alibs-slim/atags"
)

const (
	ADAPTERTYPE_DUO aconns.AdapterType = "duo"
)

// AClientDuo represents a Duo client adapter.
type AClientDuo struct {
	Type       aconns.AdapterType `json:"type,omitempty"`
	Name       aconns.AdapterName `json:"name,omitempty"`
	Url        anetwork.NetURL    `json:"url,omitempty"`
	IKey       string             `json:"ikey,omitempty"`
	SKey       string             `json:"skey,omitempty"`
	ApiHost    string             `json:"apiHost,omitempty"`
	Parameters atags.TagMapString `json:"parameters,omitempty"`

	duoClient *duoapi.DuoApi

	health aconns.HealthCheck

	mu sync.RWMutex
}

func (a *AClientDuo) GetType() aconns.AdapterType {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.Type
}

func (a *AClientDuo) GetName() aconns.AdapterName {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.Name
}

func (a *AClientDuo) GetHost() string {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.Url.Hostname()
}

func (a *AClientDuo) GetPort() int {
	a.mu.RLock()
	defer a.mu.RUnlock()
	port, err := a.Url.GetPortInt()
	if err != nil {
		return -1
	}
	return port
}

// Validate checks if the AClientDuo is valid.
func (a *AClientDuo) Validate() error {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.validate()
}

// validate checks if the AClientDuo is valid.
func (a *AClientDuo) validate() error {
	if a.Type.TrimSpace().IsEmpty() {
		a.health.Update(aconns.HEALTHSTATUS_VALIDATE_FAILED)
		return fmt.Errorf("type is empty")
	}
	if a.Name.TrimSpace().IsEmpty() {
		a.health.Update(aconns.HEALTHSTATUS_VALIDATE_FAILED)
		return fmt.Errorf("name is empty")
	}
	if a.IKey == "" || a.SKey == "" || a.ApiHost == "" {
		a.health.Update(aconns.HEALTHSTATUS_VALIDATE_FAILED)
		return fmt.Errorf("duo credentials are incomplete")
	}
	if !a.Url.IsUrl() {
		a.health.Update(aconns.HEALTHSTATUS_VALIDATE_FAILED)
		return fmt.Errorf("invalid URL")
	}
	a.health.Update(aconns.HEALTHSTATUS_HEALTHY)
	return nil
}

// Test attempts to validate the AClientDuo, initialize if necessary, and test the connection.
func (a *AClientDuo) Test() (bool, aconns.TestStatus, error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.test()
}

func (a *AClientDuo) test() (bool, aconns.TestStatus, error) {
	if err := a.validate(); err != nil {
		return false, aconns.TESTSTATUS_FAILED, err
	}

	if a.duoClient == nil {
		if err := a.init(); err != nil {
			a.health.Update(aconns.HEALTHSTATUS_OPEN_FAILED)
			return false, aconns.TESTSTATUS_FAILED, err
		}
	}

	if err := a.testConnection(); err != nil {
		a.health.Update(aconns.HEALTHSTATUS_PING_FAILED)
		return false, aconns.TESTSTATUS_FAILED, err
	}

	a.health.Update(aconns.HEALTHSTATUS_HEALTHY)
	return true, aconns.TESTSTATUS_INITIALIZED_SUCCESSFUL, nil
}

// Refresh re-initializes the Duo client.
func (a *AClientDuo) Refresh() error {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.duoClient = nil
	return a.init()
}

// Init initializes the Duo client.
func (a *AClientDuo) Init() error {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.init()
}

// init is the internal, lock-free version of Init.
func (a *AClientDuo) init() error {
	if a.duoClient != nil {
		return nil
	}

	if err := a.validate(); err != nil {
		return err
	}

	a.duoClient = duoapi.NewDuoApi(a.IKey, a.SKey, a.ApiHost, "AClientDuo")
	return nil
}

// testConnection tests the Duo connection by calling the /ping endpoint.
func (a *AClientDuo) testConnection() error {
	if a.duoClient == nil {
		return fmt.Errorf("duo client has not been initialized")
	}

	_, body, err := a.duoClient.Call("GET", "/auth/v2/ping", nil)
	if err != nil {
		return fmt.Errorf("failed to ping Duo API: %w", err)
	}

	var parsed map[string]interface{}
	if err := json.Unmarshal(body, &parsed); err != nil {
		return fmt.Errorf("failed to parse Duo ping response: %w", err)
	}

	if stat, ok := parsed["stat"].(string); !ok || stat != "OK" {
		return fmt.Errorf("Duo ping failed: %v", parsed["message"])
	}

	return nil
}

// DuoClient returns the Duo client.
func (a *AClientDuo) DuoClient() *duoapi.DuoApi {
	a.mu.RLock()
	if a.health.IsHealthy && !a.health.IsStale(5*time.Minute) {
		defer a.mu.RUnlock()
		return a.duoClient
	}
	a.mu.RUnlock()

	// Upgrade to write lock for refresh
	a.mu.Lock()
	defer a.mu.Unlock()
	if _, _, err := a.test(); err != nil {
		return nil
	}
	return a.duoClient
}

type DuoAuthResponse struct {
	duoapi.StatResult
	Response DuoAuthData `json:"response"`
}

type DuoAuthData struct {
	Result    string `json:"result"`
	Status    string `json:"status"`
	StatusMsg string `json:"status_msg"`
	TxID      string `json:"txid"`
}

type DuoResult struct {
	Allowed   bool
	Status    string
	StatusMsg string
	TxID      string
}

func (a *AClientDuo) Push2FA(username string, overrides atags.TagMapString) (*DuoResult, error) {
	params := url.Values{}
	params.Set("username", username)

	// Apply defaults from a.Parameters
	for k, v := range a.Parameters {
		params.Set(k.String(), v)
	}
	for k, v := range overrides {
		params.Set(k.String(), v)
	}

	if params.Get("factor") == "" {
		params.Set("factor", "push")
	}
	if params.Get("device") == "" {
		params.Set("device", "auto")
	}

	duoClient := a.DuoClient()
	if duoClient == nil {
		return nil, fmt.Errorf("duo client has not been initialized")
	}

	_, body, err := duoClient.Call("POST", "/auth/v2/auth", params)
	if err != nil {
		return nil, err
	}

	var parsed DuoAuthResponse
	if err := json.Unmarshal(body, &parsed); err != nil {
		return nil, fmt.Errorf("failed to parse Duo response: %w", err)
	}

	if parsed.Stat != "OK" {
		msg := "<unknown>"
		if parsed.Message != nil {
			msg = *parsed.Message
		}
		return nil, fmt.Errorf("duo error: %s", msg)
	}

	allowed := parsed.Response.Result == "allow"
	return &DuoResult{
		Allowed:   allowed,
		Status:    parsed.Response.Status,
		StatusMsg: parsed.Response.StatusMsg,
		TxID:      parsed.Response.TxID,
	}, nil
}
