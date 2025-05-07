package aclient_duo

import (
	"encoding/json"
	"fmt"
	duoapi "github.com/duosecurity/duo_api_golang"
	"github.com/jpfluger/alibs-slim/aconns"
	"github.com/jpfluger/alibs-slim/anetwork"
	"github.com/jpfluger/alibs-slim/atags"
	"net/url"
	"sync"
)

const (
	ADAPTERTYPE_DUO aconns.AdapterType = "duo"
)

type AClientDuo struct {
	Type       aconns.AdapterType `json:"type,omitempty"`
	Name       aconns.AdapterName `json:"name,omitempty"`
	Url        anetwork.NetURL    `json:"url,omitempty"`
	IKey       string             `json:"ikey,omitempty"`
	SKey       string             `json:"skey,omitempty"`
	ApiHost    string             `json:"apiHost,omitempty"`
	Parameters atags.TagMapString `json:"parameters,omitempty"`

	duoClient *duoapi.DuoApi
	mu        sync.RWMutex
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

func (a *AClientDuo) Validate() error {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.validate()
}

func (a *AClientDuo) validate() error {
	if a.Type.TrimSpace().IsEmpty() {
		return fmt.Errorf("type is empty")
	}
	if a.Name.TrimSpace().IsEmpty() {
		return fmt.Errorf("name is empty")
	}
	if a.IKey == "" || a.SKey == "" || a.ApiHost == "" {
		return fmt.Errorf("duo credentials are incomplete")
	}
	if !a.Url.IsUrl() {
		return fmt.Errorf("duo URL is not set")
	}
	if a.Url.Hostname() == "" {
		return aconns.ErrHostIsEmpty
	}
	if _, err := a.Url.GetPortInt(); err != nil {
		return err
	}
	return nil
}

func (a *AClientDuo) OpenConnection() error {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.openConnection()
}

func (a *AClientDuo) openConnection() error {
	if a.duoClient != nil {
		// Already initialized
		return nil
	}

	if err := a.validate(); err != nil {
		return err
	}

	// Create Duo client instance
	client := duoapi.NewDuoApi(a.IKey, a.SKey, a.ApiHost, a.Name.String())
	a.duoClient = client

	return nil
}

func (a *AClientDuo) CloseConnection() error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.duoClient != nil {
		// While Duo client has no Close() method, clearing it allows re-init
		a.duoClient = nil
	}
	return nil
}

func (a *AClientDuo) Test() (bool, aconns.TestStatus, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.duoClient == nil {
		if err := a.openConnection(); err != nil {
			return false, aconns.TESTSTATUS_FAILED, err
		}
	}

	_, body, err := a.duoClient.Call("GET", "/auth/v2/ping", nil)
	if err != nil {
		return false, aconns.TESTSTATUS_FAILED, err
	}

	// Define inline struct to unmarshal into
	var pingResp struct {
		duoapi.StatResult
		Response string `json:"response"`
	}

	if err := json.Unmarshal(body, &pingResp); err != nil {
		return false, aconns.TESTSTATUS_FAILED, fmt.Errorf("failed to parse Duo ping response: %w", err)
	}

	if pingResp.Stat != "OK" {
		msg := "<unknown>"
		if pingResp.Message != nil {
			msg = *pingResp.Message
		}
		return false, aconns.TESTSTATUS_FAILED, fmt.Errorf("Duo ping failed: %s", msg)
	}

	if pingResp.Response == "pong" {
		return true, aconns.TESTSTATUS_INITIALIZED_SUCCESSFUL, nil
	}

	return false, aconns.TESTSTATUS_FAILED, fmt.Errorf("unexpected Duo ping response: %v", pingResp.Response)
}

func (a *AClientDuo) GetDuoClient() *duoapi.DuoApi {
	a.mu.RLock()
	defer a.mu.RUnlock()
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

	a.mu.RLock()
	duoClient := a.duoClient
	a.mu.RUnlock()
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
