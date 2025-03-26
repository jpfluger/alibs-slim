package aclient_http

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/jpfluger/alibs-slim/aconns"
	"github.com/jpfluger/alibs-slim/anetwork"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
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
func (a *AClientHTTP) GetWithOptions(hob *HOB) (body []byte, respContentType string, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("recovered from panic: %v", r)
		}
	}()

	if hob == nil {
		return nil, "", fmt.Errorf("HOB is nil")
	}

	a.mu.RLock()
	defer a.mu.RUnlock()

	myUrl, err := a.joinUrl(hob.Path)
	if err != nil {
		return nil, "", err
	}

	client := &http.Client{Timeout: time.Duration(hob.ConnectionTimeout) * time.Second}
	resp, err := client.Get(myUrl)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	respContentType = resp.Header.Get("Content-Type")
	if hob.ExpectedType != "" && !strings.HasPrefix(respContentType, hob.ExpectedType) {
		return nil, respContentType, fmt.Errorf("unexpected content type: %s", respContentType)
	}

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, respContentType, err
	}

	return body, respContentType, nil
}

// Post performs a POST request with the given payload and returns the response body and content type.
func (a *AClientHTTP) Post(hob *HOB) ([]byte, string, error) {
	return a.PostWithOptions(hob)
}

// PostWithOptions performs a POST request with options and returns the response body and content type.
func (a *AClientHTTP) PostWithOptions(hob *HOB) (body []byte, respContentType string, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("recovered from panic: %v", r)
		}
	}()

	if hob == nil {
		return nil, "", fmt.Errorf("HOB is nil")
	}

	a.mu.RLock()
	defer a.mu.RUnlock()

	myUrl, err := a.joinUrl(hob.Path)
	if err != nil {
		return nil, "", err
	}

	hob.ContentType = strings.TrimSpace(hob.ContentType)
	if hob.ContentType == "" {
		return nil, "", fmt.Errorf("content type not defined")
	}

	var data []byte
	if len(hob.Raw) > 0 {
		data = hob.Raw
	} else {
		switch hob.ContentType {
		case "application/json":
			data, err = json.Marshal(hob.Payload)
		case "application/xml":
			data, err = xml.Marshal(hob.Payload)
		default:
			return nil, "", fmt.Errorf("unsupported content type: %s", hob.ContentType)
		}
		if err != nil {
			return nil, "", err
		}
	}

	client := &http.Client{Timeout: time.Duration(hob.ConnectionTimeout) * time.Second}
	resp, err := client.Post(myUrl, hob.ContentType, bytes.NewBuffer(data))
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	respContentType = resp.Header.Get("Content-Type")
	if hob.ExpectedType != "" && !strings.HasPrefix(respContentType, hob.ExpectedType) {
		return nil, respContentType, fmt.Errorf("unexpected content type: %s", respContentType)
	}

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, respContentType, err
	}

	return body, respContentType, nil
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
		return fmt.Errorf("type is empty")
	}
	a.Name = a.Name.TrimSpace()
	if a.Name.IsEmpty() {
		return fmt.Errorf("name is empty")
	}
	if a.Url.URL == nil {
		a.Url.URL = &url.URL{}
	}
	if a.Url.Hostname() == "" {
		return aconns.ErrHostIsEmpty
	}
	if _, err := a.Url.GetPortInt(); err != nil {
		return err
	}
	return nil
}

// Test attempts to validate the AClientHTTP and returns the test status and error if any.
// It uses a write lock to ensure thread-safe access.
func (a *AClientHTTP) Test() (bool, aconns.TestStatus, error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	if err := a.validate(); err != nil {
		return false, aconns.TESTSTATUS_FAILED, err
	}
	return false, aconns.TESTSTATUS_FAILED, fmt.Errorf("not implemented")
}
