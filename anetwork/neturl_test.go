package anetwork

import (
	"net/url"
	"testing"
)

// TestParseNetURL tests the ParseNetURL function for parsing URLs.
func TestParseNetURL(t *testing.T) {
	validURL := "http://example.com"
	_, err := ParseNetURL(validURL)
	if err != nil {
		t.Errorf("ParseNetURL() error = %v, wantErr %v", err, false)
	}

	invalidURL := "http//example.com"
	_, err = ParseNetURL(invalidURL)
	if err == nil {
		t.Errorf("ParseNetURL() error = %v, wantErr %v", err, true)
	}
}

// TestMustParseNetURL tests the MustParseNetURL function for parsing URLs without returning errors.
func TestMustParseNetURL(t *testing.T) {
	validURL := "http://example.com"
	result := MustParseNetURL(validURL)
	if result.String() != validURL {
		t.Errorf("MustParseNetURL() = %v, want %v", result.String(), validURL)
	}

	invalidURL := "http//example.com"
	result = MustParseNetURL(invalidURL)
	if result.String() != "" {
		t.Errorf("MustParseNetURL() = %v, want %v", result.String(), "")
	}
}

// TestParseNetURLNoError tests the ParseNetURLNoError function for parsing URLs and returning nil on error.
func TestParseNetURLNoError(t *testing.T) {
	validURL := "http://example.com"
	result := ParseNetURLNoError(validURL)
	if result.String() != validURL {
		t.Errorf("ParseNetURLNoError() = %v, want %v", result.String(), validURL)
	}

	invalidURL := "http//example.com"
	result = ParseNetURLNoError(invalidURL)
	if result != nil {
		t.Errorf("ParseNetURLNoError() = %v, want %v", result, nil)
	}
}

// TestGetUrlPathOrRoot tests the GetUrlPathOrRoot function for returning the URL path or root.
func TestGetUrlPathOrRoot(t *testing.T) {
	u, _ := url.Parse("http://example.com/path")
	result := GetUrlPathOrRoot(u)
	expected := "/path"
	if result != expected {
		t.Errorf("GetUrlPathOrRoot() = %v, want %v", result, expected)
	}

	u, _ = url.Parse("http://example.com")
	result = GetUrlPathOrRoot(u)
	expected = "/"
	if result != expected {
		t.Errorf("GetUrlPathOrRoot() = %v, want %v", result, expected)
	}
}

// TestNetURL_GetSchemeHost tests the GetSchemeHost method for returning the scheme and host of the URL.
func TestNetURL_GetSchemeHost(t *testing.T) {
	u := &NetURL{URL: &url.URL{Scheme: "http", Host: "example.com"}}
	result := u.GetSchemeHost()
	expected := "http://example.com"
	if result != expected {
		t.Errorf("GetSchemeHost() = %v, want %v", result, expected)
	}
}

// TestNetURL_IsHttps tests the IsHttps method for checking if the URL scheme is HTTPS.
func TestNetURL_IsHttps(t *testing.T) {
	u := &NetURL{URL: &url.URL{Scheme: "https"}}
	if !u.IsHttps() {
		t.Errorf("IsHttps() = %v, want %v", u.IsHttps(), true)
	}
}

// TestNetURL_IsUrl tests the IsUrl method for checking if the NetURL is a valid URL.
func TestNetURL_IsUrl(t *testing.T) {
	u := &NetURL{URL: &url.URL{Scheme: "http", Host: "example.com"}}
	if !u.IsUrl() {
		t.Errorf("IsUrl() = %v, want %v", u.IsUrl(), true)
	}

	urls := map[string]bool{
		"https":                        false,
		"https://":                     false,
		"":                             false,
		"http://www":                   true,
		"http://localhost":             true,
		"http://www.example.com":       true,
		"http://www.example.com:80":    true,
		"http://www.example.com:8080":  true,
		"https://www.example.com":      true,
		"https://www.example.com:443":  true,
		"https://www.example.com:8080": true,
		"/testing-path":                false,
		"testing-path":                 false,
		"alskjff#?asf//dfas":           false,
	}

	for key, val := range urls {
		nu, err := ParseNetURL(key)
		if err != nil {
			if val {
				t.Errorf("expected valid url for '%s'; %v", key, err)
				continue
			}
		}

		if val != nu.IsUrl() {
			if val {
				t.Errorf("expected valid url for '%s'", key)
			} else {
				t.Errorf("expected invalid url for '%s'", key)
			}
		}
	}
}

// TestNetURL_IsReachable tests the IsReachable method for checking if the URL is reachable.
func TestNetURL_IsReachable(t *testing.T) {
	u := &NetURL{URL: &url.URL{Scheme: "http", Host: "example.com"}}
	ok, err := u.IsReachable()
	if err != nil {
		t.Errorf("IsReachable() error = %v", err)
	}
	if !ok {
		t.Errorf("IsReachable() = %v, want %v", ok, true)
	}
}

// TestNetURL_GetListenerKey tests the GetListenerKey method for returning the listener key.
func TestNetURL_GetListenerKey(t *testing.T) {
	u := &NetURL{URL: &url.URL{Scheme: "http", Host: "example.com:8080"}}
	result := u.GetListenerKey()
	expected := "example.com:8080"
	if result != expected {
		t.Errorf("GetListenerKey() = %v, want %v", result, expected)
	}

	urls := map[string]string{
		"https":                        "",
		"https://":                     "",
		"":                             "",
		"http://www":                   "www:80",
		"http://localhost":             "localhost:80",
		"http://www.example.com":       "www.example.com:80",
		"http://www.example.com:80":    "www.example.com:80",
		"http://www.example.com:8080":  "www.example.com:8080",
		"https://www.example.com":      "www.example.com:443",
		"https://www.example.com:443":  "www.example.com:443",
		"https://www.example.com:8080": "www.example.com:8080",
		"/testing-path":                "",
		"testing-path":                 "",
		"alskjff#?asf//dfas":           "",
	}

	for key, val := range urls {
		nu, err := ParseNetURL(key)
		if err != nil {
			continue
		}

		if val != nu.GetListenerKey() {
			t.Errorf("expected listener key '%s' for '%s' but got '%s'", val, key, nu.GetListenerKey())
			continue
		}
	}
}

func TestNewUrlJoinPath(t *testing.T) {
	tests := []struct {
		baseURL     string
		paths       []string
		expectedURL string
		expectError bool
	}{
		{
			baseURL:     "https://example.com",
			paths:       []string{"new", "path"},
			expectedURL: "https://example.com/new/path",
			expectError: false,
		},
		{
			baseURL:     "https://example.com/base",
			paths:       []string{"new", "path"},
			expectedURL: "https://example.com/base/new/path",
			expectError: false,
		},
		{
			baseURL:     "https://example.com/base/",
			paths:       []string{"new", "path"},
			expectedURL: "https://example.com/base/new/path",
			expectError: false,
		},
		{
			baseURL:     "https://example.com",
			paths:       []string{},
			expectedURL: "https://example.com",
			expectError: false,
		},
		{
			baseURL:     "",
			paths:       []string{"new", "path"},
			expectedURL: "",
			expectError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.baseURL, func(t *testing.T) {
			parsedURL, err := url.Parse(test.baseURL)
			if err != nil && !test.expectError {
				t.Fatalf("unexpected error parsing base URL: %v", err)
			}

			nu := &NetURL{URL: parsedURL}
			result, err := nu.NewUrlJoinPath(test.paths...)

			if test.expectError {
				if err == nil {
					t.Fatalf("expected an error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if result != test.expectedURL {
				t.Errorf("expected %s, got %s", test.expectedURL, result)
			}
		})
	}
}
