package aclient_http

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/jpfluger/alibs-slim/aconns"
	"github.com/jpfluger/alibs-slim/anetwork"
)

const (
	ADAPTERTYPE_HTTP        aconns.AdapterType = "http"
	HTTP_CONNECTION_TIMEOUT                    = 10
)

// AClientHTTP satisfies IAdapter.
// Use it to compose connector structs.
type AClientHTTP struct {
	Type aconns.AdapterType `json:"type,omitempty"`
	Name aconns.AdapterName `json:"name,omitempty"`
	Url  anetwork.NetURL    `json:"url,omitempty"`

	health aconns.HealthCheck

	mu sync.RWMutex
}

func NewAClientHTTP(myUrl string) (*AClientHTTP, error) {
	u, err := anetwork.ParseNetURL(myUrl)
	if err != nil {
		return nil, err
	}
	return &AClientHTTP{
		Type: ADAPTERTYPE_HTTP,
		Name: "http-client",
		Url:  *u,
	}, nil
}

func (a *AClientHTTP) GetType() aconns.AdapterType {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.Type
}

func (a *AClientHTTP) GetName() aconns.AdapterName {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.Name
}

func (a *AClientHTTP) GetHost() string {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.Url.Hostname()
}

func (a *AClientHTTP) GetPort() int {
	a.mu.RLock()
	defer a.mu.RUnlock()
	port, err := a.Url.GetPortInt()
	if err != nil {
		return -1
	}
	return port
}

// joinUrl joins the given path to the base URL.
func (a *AClientHTTP) joinUrl(withPath string) (string, error) {
	myUrl, err := a.Url.NewUrlJoinPath(withPath)
	if err != nil {
		return "", err
	}
	return myUrl, nil
}

// Get performs a GET request and returns the response body and content type.
func (a *AClientHTTP) Get(hob *HOB) ([]byte, string, error) {
	return a.GetWithOptions(hob)
}

// GetWithOptions performs a GET request with options and returns the response body and content type.
func (a *AClientHTTP) GetWithOptions(hob *HOB) ([]byte, string, error) {
	a.mu.RLock()
	if a.health.IsHealthy && !a.health.IsStale(5*time.Minute) {
		defer a.mu.RUnlock()
	} else {
		a.mu.RUnlock()
		a.mu.Lock()
		defer a.mu.Unlock()
		if _, _, err := a.test(); err != nil {
			return nil, "", err
		}
	}

	fullURL, err := a.joinUrl(hob.Path)
	if err != nil {
		return nil, "", err
	}

	return DoHTTPGet(fullURL, hob)
}

// Post performs a POST request with the given payload and returns the response body and content type.
func (a *AClientHTTP) Post(hob *HOB) ([]byte, string, error) {
	return a.PostWithOptions(hob)
}

// PostWithOptions performs a POST request with options and returns the response body and content type.
func (a *AClientHTTP) PostWithOptions(hob *HOB) ([]byte, string, error) {
	a.mu.RLock()
	if a.health.IsHealthy && !a.health.IsStale(5*time.Minute) {
		defer a.mu.RUnlock()
	} else {
		a.mu.RUnlock()
		a.mu.Lock()
		defer a.mu.Unlock()
		if _, _, err := a.test(); err != nil {
			return nil, "", err
		}
	}

	fullURL, err := a.joinUrl(hob.Path)
	if err != nil {
		return nil, "", err
	}

	return DoHTTPPost(fullURL, hob)
}

// GetJSON performs a GET request and parses the response as JSON into the provided interface.
func (a *AClientHTTP) GetJSON(hob *HOB, v interface{}) error {
	body, _, err := a.GetWithOptions(hob)
	if err != nil {
		return err
	}
	return json.Unmarshal(body, v)
}

// GetXML performs a GET request and parses the response as XML into the provided interface.
func (a *AClientHTTP) GetXML(hob *HOB, v interface{}) error {
	body, _, err := a.GetWithOptions(hob)
	if err != nil {
		return err
	}
	return xml.Unmarshal(body, v)
}

// PostJSON performs a POST request with a JSON payload and parses the response as JSON into the provided interface.
func (a *AClientHTTP) PostJSON(hob *HOB, v interface{}) error {
	body, _, err := a.PostWithOptions(hob)
	if err != nil {
		return err
	}
	return json.Unmarshal(body, v)
}

// PostXML performs a POST request with an XML payload and parses the response as XML into the provided interface.
func (a *AClientHTTP) PostXML(hob *HOB, v interface{}) error {
	body, _, err := a.PostWithOptions(hob)
	if err != nil {
		return err
	}
	return xml.Unmarshal(body, v)
}

// Validate checks if the AClientHTTP is valid.
func (a *AClientHTTP) Validate() error {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.validate()
}

// validate checks if the AClientHTTP is valid.
func (a *AClientHTTP) validate() error {
	a.Type = a.Type.TrimSpace()
	if a.Type.IsEmpty() {
		a.health.Update(aconns.HEALTHSTATUS_VALIDATE_FAILED)
		return fmt.Errorf("type is empty")
	}
	a.Name = a.Name.TrimSpace()
	if a.Name.IsEmpty() {
		a.health.Update(aconns.HEALTHSTATUS_VALIDATE_FAILED)
		return fmt.Errorf("name is empty")
	}
	if a.Url.URL == nil {
		a.Url.URL = &url.URL{}
	}
	if a.Url.Hostname() == "" {
		a.health.Update(aconns.HEALTHSTATUS_VALIDATE_FAILED)
		return aconns.ErrHostIsEmpty
	}
	if _, err := a.Url.GetPortInt(); err != nil {
		a.health.Update(aconns.HEALTHSTATUS_VALIDATE_FAILED)
		return err
	}
	a.health.Update(aconns.HEALTHSTATUS_HEALTHY)
	return nil
}

// Test attempts to validate the AClientHTTP and test connectivity with a HEAD request.
func (a *AClientHTTP) Test() (bool, aconns.TestStatus, error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.test()
}

func (a *AClientHTTP) test() (bool, aconns.TestStatus, error) {
	if err := a.validate(); err != nil {
		return false, aconns.TESTSTATUS_FAILED, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) // Test timeout
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "HEAD", a.Url.String(), nil)
	if err != nil {
		a.health.Update(aconns.HEALTHSTATUS_OPEN_FAILED)
		return false, aconns.TESTSTATUS_FAILED, err
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		status := aconns.HEALTHSTATUS_PING_FAILED
		if context.DeadlineExceeded == err {
			status = aconns.HEALTHSTATUS_TIMEOUT
		} else if strings.Contains(err.Error(), "network") || strings.Contains(err.Error(), "connection refused") {
			status = aconns.HEALTHSTATUS_NETWORK_ERROR
		}
		a.health.Update(status)
		return false, aconns.TESTSTATUS_FAILED, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		a.health.Update(aconns.HEALTHSTATUS_PING_FAILED)
		return false, aconns.TESTSTATUS_FAILED, fmt.Errorf("HEAD request failed with status: %d", resp.StatusCode)
	}

	a.health.Update(aconns.HEALTHSTATUS_HEALTHY)
	return true, aconns.TESTSTATUS_INITIALIZED_SUCCESSFUL, nil
}

// Refresh is a no-op for stateless HTTP clients.
func (a *AClientHTTP) Refresh() error {
	// No state to refresh; no-op
	return nil
}

//package aclient_http
//
//import (
//	"encoding/json"
//	"encoding/xml"
//	"fmt"
//	"github.com/jpfluger/alibs-slim/aconns"
//	"github.com/jpfluger/alibs-slim/anetwork"
//	"net/url"
//	"sync"
//)
//
//const (
//	ADAPTERTYPE_HTTP        aconns.AdapterType = "http"
//	HTTP_CONNECTION_TIMEOUT                    = 10
//)
//
//// AClientHTTP satisfies IAdapter.
//// Use it to compose connector structs.
//type AClientHTTP struct {
//	Type aconns.AdapterType `json:"type,omitempty"`
//	Name aconns.AdapterName `json:"name,omitempty"`
//	Url  anetwork.NetURL    `json:"url,omitempty"`
//
//	mu sync.RWMutex
//}
//
//func NewAClientHTTP(myUrl string) (*AClientHTTP, error) {
//	u, err := anetwork.ParseNetURL(myUrl)
//	if err != nil {
//		return nil, err
//	}
//	return &AClientHTTP{
//		Type: ADAPTERTYPE_HTTP,
//		Name: "http-client",
//		Url:  *u,
//	}, nil
//}
//
//func (a *AClientHTTP) GetType() aconns.AdapterType {
//	a.mu.RLock()
//	defer a.mu.RUnlock()
//	return a.Type
//}
//
//func (a *AClientHTTP) GetName() aconns.AdapterName {
//	a.mu.RLock()
//	defer a.mu.RUnlock()
//	return a.Name
//}
//
//func (a *AClientHTTP) GetHost() string {
//	a.mu.RLock()
//	defer a.mu.RUnlock()
//	return a.Url.Hostname()
//}
//
//func (a *AClientHTTP) GetPort() int {
//	a.mu.RLock()
//	defer a.mu.RUnlock()
//	port, err := a.Url.GetPortInt()
//	if err != nil {
//		return -1
//	}
//	return port
//}
//
//// joinUrl joins the given path to the base URL.
//func (a *AClientHTTP) joinUrl(withPath string) (string, error) {
//	myUrl, err := a.Url.NewUrlJoinPath(withPath)
//	if err != nil {
//		return "", err
//	}
//	return myUrl, nil
//}
//
//// Get performs a GET request and returns the response body and content type.
//func (a *AClientHTTP) Get(hob *HOB) ([]byte, string, error) {
//	return a.GetWithOptions(hob)
//}
//
//// GetWithOptions performs a GET request with options and returns the response body and content type.
//func (a *AClientHTTP) GetWithOptions(hob *HOB) ([]byte, string, error) {
//	a.mu.RLock()
//	defer a.mu.RUnlock()
//
//	fullURL, err := a.joinUrl(hob.Path)
//	if err != nil {
//		return nil, "", err
//	}
//
//	return DoHTTPGet(fullURL, hob)
//}
//
//// Post performs a POST request with the given payload and returns the response body and content type.
//func (a *AClientHTTP) Post(hob *HOB) ([]byte, string, error) {
//	return a.PostWithOptions(hob)
//}
//
//// PostWithOptions performs a POST request with options and returns the response body and content type.
//func (a *AClientHTTP) PostWithOptions(hob *HOB) ([]byte, string, error) {
//	a.mu.RLock()
//	defer a.mu.RUnlock()
//
//	fullURL, err := a.joinUrl(hob.Path)
//	if err != nil {
//		return nil, "", err
//	}
//
//	return DoHTTPPost(fullURL, hob)
//}
//
//// GetJSON performs a GET request and parses the response as JSON into the provided interface.
//func (a *AClientHTTP) GetJSON(hob *HOB, v interface{}) error {
//	body, _, err := a.GetWithOptions(hob)
//	if err != nil {
//		return err
//	}
//	return json.Unmarshal(body, v)
//}
//
//// GetXML performs a GET request and parses the response as XML into the provided interface.
//func (a *AClientHTTP) GetXML(hob *HOB, v interface{}) error {
//	body, _, err := a.GetWithOptions(hob)
//	if err != nil {
//		return err
//	}
//	return xml.Unmarshal(body, v)
//}
//
//// PostJSON performs a POST request with a JSON payload and parses the response as JSON into the provided interface.
//func (a *AClientHTTP) PostJSON(hob *HOB, v interface{}) error {
//	body, _, err := a.PostWithOptions(hob)
//	if err != nil {
//		return err
//	}
//	return json.Unmarshal(body, v)
//}
//
//// PostXML performs a POST request with an XML payload and parses the response as XML into the provided interface.
//func (a *AClientHTTP) PostXML(hob *HOB, v interface{}) error {
//	body, _, err := a.PostWithOptions(hob)
//	if err != nil {
//		return err
//	}
//	return xml.Unmarshal(body, v)
//}
//
//// Validate checks if the AClientHTTP is valid.
//func (a *AClientHTTP) Validate() error {
//	a.mu.Lock()
//	defer a.mu.Unlock()
//	return a.validate()
//}
//
//// validate checks if the AClientHTTP is valid.
//func (a *AClientHTTP) validate() error {
//	a.Type = a.Type.TrimSpace()
//	if a.Type.IsEmpty() {
//		return fmt.Errorf("type is empty")
//	}
//	a.Name = a.Name.TrimSpace()
//	if a.Name.IsEmpty() {
//		return fmt.Errorf("name is empty")
//	}
//	if a.Url.URL == nil {
//		a.Url.URL = &url.URL{}
//	}
//	if a.Url.Hostname() == "" {
//		return aconns.ErrHostIsEmpty
//	}
//	if _, err := a.Url.GetPortInt(); err != nil {
//		return err
//	}
//	return nil
//}
//
//// Test attempts to validate the AClientHTTP and returns the test status and error if any.
//// It uses a write lock to ensure thread-safe access.
//func (a *AClientHTTP) Test() (bool, aconns.TestStatus, error) {
//	a.mu.Lock()
//	defer a.mu.Unlock()
//	if err := a.validate(); err != nil {
//		return false, aconns.TESTSTATUS_FAILED, err
//	}
//	return false, aconns.TESTSTATUS_FAILED, fmt.Errorf("not implemented")
//}
