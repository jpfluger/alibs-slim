package aclient_duo

import (
	"crypto/tls"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	duoapi "github.com/duosecurity/duo_api_golang"
	"github.com/jpfluger/alibs-slim/aconns"
	"github.com/jpfluger/alibs-slim/anetwork"
	"github.com/jpfluger/alibs-slim/atags"
	"github.com/stretchr/testify/assert"
)

// mockDuoServer creates a mock server for Duo API endpoints.
func mockDuoServer(t *testing.T, stat string, response interface{}) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" && r.URL.Path == "/auth/v2/ping" {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]interface{}{"stat": stat})
			return
		}
		if r.Method == "POST" && r.URL.Path == "/auth/v2/auth" {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"stat":     stat,
				"response": response,
			})
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
}

// TestAClientDuo_Validate_Success tests successful validation.
func TestAClientDuo_Validate_Success(t *testing.T) {
	client := &AClientDuo{
		Type:    ADAPTERTYPE_DUO,
		Name:    "test_duo",
		IKey:    "ikey",
		SKey:    "skey",
		ApiHost: "api.example.com",
		Url:     anetwork.NetURL{URL: &url.URL{Scheme: "https", Host: "api.example.com"}},
	}

	err := client.Validate()
	assert.NoError(t, err)
}

// TestAClientDuo_Validate_Failure tests validation failures.
func TestAClientDuo_Validate_Failure(t *testing.T) {
	tests := []struct {
		name  string
		setup func(*AClientDuo)
	}{
		{"Empty Type", func(c *AClientDuo) { c.Type = "" }},
		{"Empty Name", func(c *AClientDuo) { c.Name = "" }},
		{"Missing Credentials", func(c *AClientDuo) { c.IKey = ""; c.SKey = ""; c.ApiHost = "" }},
		{"Invalid URL", func(c *AClientDuo) { c.Url = anetwork.NetURL{} }},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			client := &AClientDuo{
				Type:    ADAPTERTYPE_DUO,
				Name:    "test_duo",
				IKey:    "ikey",
				SKey:    "skey",
				ApiHost: "api.example.com",
				Url:     anetwork.NetURL{URL: &url.URL{Scheme: "https", Host: "api.example.com"}},
			}
			tc.setup(client)
			err := client.Validate()
			assert.Error(t, err)
		})
	}
}

// TestAClientDuo_Test_Success tests successful Test with mock ping.
func TestAClientDuo_Test_Success(t *testing.T) {
	mockServer := mockDuoServer(t, "OK", nil)
	defer mockServer.Close()

	host := strings.TrimPrefix(mockServer.URL, "http://")

	client := &AClientDuo{
		Type:    ADAPTERTYPE_DUO,
		Name:    "test_duo",
		IKey:    "ikey",
		SKey:    "skey",
		ApiHost: host,
		Url:     anetwork.NetURL{URL: &url.URL{Scheme: "http", Host: host}},
	}

	// Custom HTTP client to use HTTP for mock
	duo := duoapi.NewDuoApi(client.IKey, client.SKey, client.ApiHost, "AClientDuo")
	duo.SetCustomHTTPClient(&http.Client{
		Transport: &httpTransport{},
	})
	client.duoClient = duo

	ok, status, err := client.Test()
	assert.True(t, ok)
	assert.Equal(t, aconns.TESTSTATUS_INITIALIZED_SUCCESSFUL, status)
	assert.NoError(t, err)
}

// TestAClientDuo_Test_Failure tests Test failure scenarios.
func TestAClientDuo_Test_Failure(t *testing.T) {
	mockServer := mockDuoServer(t, "FAIL", map[string]string{"message": "error"})
	defer mockServer.Close()

	host := strings.TrimPrefix(mockServer.URL, "http://")

	client := &AClientDuo{
		Type:    ADAPTERTYPE_DUO,
		Name:    "test_duo",
		IKey:    "ikey",
		SKey:    "skey",
		ApiHost: host,
		Url:     anetwork.NetURL{URL: &url.URL{Scheme: "http", Host: host}},
	}

	// Custom HTTP client to use HTTP for mock
	duo := duoapi.NewDuoApi(client.IKey, client.SKey, client.ApiHost, "AClientDuo")
	duo.SetCustomHTTPClient(&http.Client{
		Transport: &httpTransport{},
	})
	client.duoClient = duo

	ok, _, err := client.Test()
	assert.False(t, ok)
	assert.Error(t, err)
}

// TestAClientDuo_Refresh tests refreshing the Duo client.
func TestAClientDuo_Refresh(t *testing.T) {
	mockServer := mockDuoServer(t, "OK", nil)
	defer mockServer.Close()

	host := strings.TrimPrefix(mockServer.URL, "https://")

	client := &AClientDuo{
		Type:    ADAPTERTYPE_DUO,
		Name:    "test_duo",
		IKey:    "ikey",
		SKey:    "skey",
		ApiHost: host,
		Url:     anetwork.NetURL{URL: &url.URL{Scheme: "https", Host: host}},
	}

	// Custom HTTP client with InsecureSkipVerify for TLS mock
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	duo := duoapi.NewDuoApi(client.IKey, client.SKey, client.ApiHost, "AClientDuo")
	duo.SetCustomHTTPClient(&http.Client{Transport: tr})
	client.duoClient = duo

	// Refresh
	err := client.Refresh()
	assert.NoError(t, err)
	assert.NotEqual(t, client.duoClient, duo) // Should be reset
}

// TestAClientDuo_Push2FA_Success tests successful Push2FA with mock response.
func TestAClientDuo_Push2FA_Success(t *testing.T) {
	mockServer := mockDuoServer(t, "OK", map[string]string{
		"result":     "allow",
		"status":     "success",
		"status_msg": "OK",
		"txid":       "tx123",
	})
	defer mockServer.Close()

	host := strings.TrimPrefix(mockServer.URL, "http://")

	client := &AClientDuo{
		Type:    ADAPTERTYPE_DUO,
		Name:    "test_duo",
		IKey:    "ikey",
		SKey:    "skey",
		ApiHost: host,
		Url:     anetwork.NetURL{URL: &url.URL{Scheme: "http", Host: host}},
	}

	// Custom HTTP client to use HTTP for mock
	duo := duoapi.NewDuoApi(client.IKey, client.SKey, client.ApiHost, "AClientDuo")
	duo.SetCustomHTTPClient(&http.Client{
		Transport: &httpTransport{},
	})
	client.duoClient = duo

	result, err := client.Push2FA("user", atags.TagMapString{})
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.True(t, result.Allowed)
	assert.Equal(t, "success", result.Status)
}

// TestAClientDuo_Push2FA_Failure tests Push2FA failure.
func TestAClientDuo_Push2FA_Failure(t *testing.T) {
	mockServer := mockDuoServer(t, "FAIL", map[string]string{"message": "error"})
	defer mockServer.Close()

	host := strings.TrimPrefix(mockServer.URL, "http://")

	client := &AClientDuo{
		Type:    ADAPTERTYPE_DUO,
		Name:    "test_duo",
		IKey:    "ikey",
		SKey:    "skey",
		ApiHost: host,
		Url:     anetwork.NetURL{URL: &url.URL{Scheme: "http", Host: host}},
	}

	// Custom HTTP client to use HTTP for mock
	duo := duoapi.NewDuoApi(client.IKey, client.SKey, client.ApiHost, "AClientDuo")
	duo.SetCustomHTTPClient(&http.Client{
		Transport: &httpTransport{},
	})
	client.duoClient = duo

	_, err := client.Push2FA("user", atags.TagMapString{})
	assert.Error(t, err)
}

// httpTransport forces HTTP scheme for mock requests.
type httpTransport struct{}

func (httpTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.URL.Scheme = "http"
	return http.DefaultTransport.RoundTrip(req)
}
