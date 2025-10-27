package anetwork

import (
	"net"
	"net/url"
	"strconv"
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
	if result.IsEmpty() || result.String() != validURL {
		t.Errorf("MustParseNetURL() = empty or %q, want %q", result.String(), validURL)
	}

	invalidURL := "http//example.com"
	result = MustParseNetURL(invalidURL)
	if !result.IsEmpty() {
		t.Errorf("MustParseNetURL() = %q, want empty", result.String())
	}

	// Test file URL
	fileURL := "file:///path/to/file"
	result = MustParseNetURL(fileURL)
	if result.IsEmpty() || result.String() != fileURL {
		t.Errorf("MustParseNetURL() for file = empty or %q, want %q", result.String(), fileURL)
	}
}

// TestParseNetURLNoError tests the ParseNetURLNoError function for parsing URLs and returning nil on error.
func TestParseNetURLNoError(t *testing.T) {
	validURL := "http://example.com"
	result := ParseNetURLNoError(validURL)
	if result == nil || result.String() != validURL {
		t.Errorf("ParseNetURLNoError() = nil or %q, want %q", result.String(), validURL)
	}

	invalidURL := "http//example.com"
	result = ParseNetURLNoError(invalidURL)
	if result != nil {
		t.Errorf("ParseNetURLNoError() = %q, want nil", result.String())
	}

	// Test file URL
	fileURL := "file:///path/to/file"
	result = ParseNetURLNoError(fileURL)
	if result == nil || result.String() != fileURL {
		t.Errorf("ParseNetURLNoError() for file = nil or %q, want %q", result.String(), fileURL)
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

func TestNetURL_IsHttps(t *testing.T) {
	tests := []struct {
		name     string
		url      *NetURL
		expected bool
	}{
		{
			name:     "nil URL",
			url:      &NetURL{},
			expected: false,
		},
		{
			name:     "https URL with host",
			url:      MustParseNetURL("https://example.com"),
			expected: true,
		},
		{
			name:     "uppercase HTTPS URL",
			url:      &NetURL{URL: &url.URL{Scheme: "HTTPS", Host: "example.com"}},
			expected: true,
		},
		{
			name:     "https URL no host",
			url:      MustParseNetURL("https://"),
			expected: false,
		},
		{
			name:     "http URL",
			url:      MustParseNetURL("http://example.com"),
			expected: false,
		},
		{
			name:     "file URL",
			url:      MustParseNetURL("file:///path"),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.url.IsHttps(); got != tt.expected {
				t.Errorf("IsHttps() = %v, want %v", got, tt.expected)
			}
		})
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
		"http://www":                   true, // single-label host allowed here
		"oci://registry/image":         true, // authority scheme with host
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

// Mock?
//// TestNetURL_IsReachable tests the IsReachable method for checking if the URL is reachable.
//func TestNetURL_IsReachable(t *testing.T) {
//	u := &NetURL{URL: &url.URL{Scheme: "http", Host: "example.com"}}
//	ok, err := u.IsReachable()
//	if err != nil {
//		t.Errorf("IsReachable() error = %v", err)
//	}
//	if !ok {
//		t.Errorf("IsReachable() = %v, want %v", ok, true)
//	}
//}

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

// helper that ensures the port returned by the probe is still available
func assertPortIsFree(t *testing.T, host string, port int) {
	t.Helper()
	l, err := net.Listen("tcp4", net.JoinHostPort(host, strconv.Itoa(port)))
	if err != nil {
		t.Fatalf("port %d is not free on %s: %v", port, host, err)
	}
	_ = l.Close()
}

func TestFindNextOpenIPv4Port(t *testing.T) {
	t.Parallel()

	t.Run("defaults to localhost & first ephemeral port", func(t *testing.T) {
		p, err := FindNextOpenPort(0)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if p < 49152 || p > 65535 {
			t.Fatalf("returned port %d outside expected ephemeral range", p)
		}
		assertPortIsFree(t, "127.0.0.1", p)
	})

	t.Run("skips a port that is already bound", func(t *testing.T) {
		ln, err := net.Listen("tcp4", "127.0.0.1:0") // let the OS choose a free one
		if err != nil {
			t.Fatalf("failed to acquire test listener: %v", err)
		}
		defer ln.Close()

		occupied := ln.Addr().(*net.TCPAddr).Port
		free, err := FindNextOpenPort(occupied)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if free == occupied {
			t.Fatalf("expected a different port than the occupied %d", occupied)
		}
		assertPortIsFree(t, "127.0.0.1", free)
	})
}

func TestWithNextOpenPort(t *testing.T) {
	t.Parallel()

	base := MustParseNetURL("http://localhost:8080")
	updated, newPort, err := base.WithNextOpenPort(0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if updated == nil || !updated.IsUrl() {
		t.Fatalf("returned NetURL is nil or invalid")
	}
	if newPortStr := updated.Port(); newPortStr == "" {
		t.Fatalf("updated URL has no port: %s", updated.String())
	} else if p, _ := strconv.Atoi(newPortStr); p != newPort {
		t.Fatalf("reported port %d does not match URL port %d", newPort, p)
	}
	assertPortIsFree(t, "127.0.0.1", newPort)
}

func TestNetURL_IsEmpty(t *testing.T) {
	tests := []struct {
		name     string
		url      *NetURL
		expected bool
	}{
		{
			name:     "nil URL",
			url:      &NetURL{},
			expected: true,
		},
		{
			name:     "empty string URL",
			url:      &NetURL{URL: &url.URL{}},
			expected: true,
		},
		{
			name:     "whitespace URL",
			url:      MustParseNetURL(" "),
			expected: true,
		},
		{
			name:     "non-empty URL",
			url:      MustParseNetURL("file://localhost/path"),
			expected: false,
		},
		{
			name:     "http URL",
			url:      MustParseNetURL("https://example.com"),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.url.IsEmpty(); got != tt.expected {
				t.Errorf("IsEmpty() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestNetURL_IsFile(t *testing.T) {
	tests := []struct {
		name     string
		url      *NetURL
		expected bool
	}{
		{
			name:     "nil URL",
			url:      &NetURL{},
			expected: false,
		},
		{
			name:     "file URL with path",
			url:      MustParseNetURL("file://localhost/path"),
			expected: true,
		},
		{
			name:     "uppercase file URL",
			url:      MustParseNetURL("FILE://localhost/path"),
			expected: true,
		},
		{
			name:     "file URL no path",
			url:      MustParseNetURL("file:///"),
			expected: false,
		},
		{
			name:     "http URL",
			url:      MustParseNetURL("https://example.com"),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.url.IsFile(); got != tt.expected {
				t.Errorf("IsFile() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestFileNormalization(t *testing.T) {
	type okCase struct {
		in   string
		want string
	}
	ok := []okCase{
		{"file:/var/log/app.log", "file:///var/log/app.log"},
		{"file:///var/log/app.log", "file:///var/log/app.log"},
		{"file://server/share/log", "file://server/share/log"},
		{"file:///C:/Windows/notepad", "file:///C:/Windows/notepad"},
		// If you choose to accept root, uncomment the next line and remove it from "bad":
		// {"file:/", "file:///"},
	}

	for _, tc := range ok {
		t.Run("OK_"+tc.in, func(t *testing.T) {
			got, err := ParseNetURL(tc.in)
			if err != nil {
				t.Fatalf("ParseNetURL(%q) unexpected error: %v", tc.in, err)
			}
			if gotStr := got.String(); gotStr != tc.want {
				t.Fatalf("ParseNetURL(%q) => %q; want %q\n  scheme=%q host=%q path=%q",
					tc.in, gotStr, tc.want, got.Scheme, got.Host, got.Path)
			}
		})
	}

	bad := []string{
		"file://",     // no path
		"file:/",      // root only (reject if you apply Option B)
		"file:foo",    // relative/opaque
		"file://host", // no path
	}

	for _, in := range bad {
		t.Run("BAD_"+in, func(t *testing.T) {
			got, err := ParseNetURL(in)
			if err == nil {
				t.Fatalf("expected error for %q, got success: %q\n  scheme=%q host=%q path=%q",
					in, got.String(), got.Scheme, got.Host, got.Path)
			}
		})
	}
}

func TestNormalizeURL(t *testing.T) {
	tests := []struct {
		name     string
		input    NetURL
		expected string
	}{
		{"Empty", NetURL{}, ""},

		// Already web
		{"Valid HTTPS", *MustParseNetURL("https://example.com"), "https://example.com"},
		{"Valid HTTP", *MustParseNetURL("http://example.com"), "http://example.com"},

		// Schemeless -> assume https
		{"Scheme-less Domain", NetURL{URL: &url.URL{Path: "example.com"}}, "https://example.com"},
		// Protocol-relative (//host) -> https://host
		{"Scheme-less with //", NetURL{URL: &url.URL{Host: "example.com"}}, "https://example.com"},
		{"Scheme-less with Path", NetURL{URL: &url.URL{Path: "example.com/path"}}, "https://example.com/path"},
		{"Scheme-less with Host and Path", NetURL{URL: &url.URL{Host: "example.com", Path: "/path"}}, "https://example.com/path"},

		// file: left as-is by NormalizeURL; Go prints file:/// when Host is empty and Path is absolute
		{"File Scheme Unchanged", NetURL{URL: &url.URL{Scheme: "file", Path: "/path/to/file"}}, "file:///path/to/file"},

		// Open-ended non-web scheme (authority form) stays unchanged
		{"OCI Scheme Unchanged", *MustParseNetURL("oci://registry/image"), "oci://registry/image"},

		// Inputs that should fail normalization to a valid web URL
		// Host that starts with "//" ends up as "invalid" after trimming; since it's single-label, reject.
		{"Invalid After Fix", NetURL{URL: &url.URL{Host: "//invalid"}}, ""},
		// Leading/trailing spaces in schemeless input are trimmed
		{"Spaced Input", NetURL{URL: &url.URL{Path: " example.com "}}, "https://example.com"},
		// https with single-label host is invalid per policy (must be IP, localhost, or contain a dot)
		{"Invalid HTTPS", NetURL{URL: &url.URL{Scheme: "https", Host: "invalid"}}, ""},

		// Allowed special hosts
		{"Valid Localhost", *MustParseNetURL("http://localhost"), "http://localhost"},
		{"Valid IP", *MustParseNetURL("https://192.168.0.1"), "https://192.168.0.1"},

		// Schemeless invalid single-label host
		{"Scheme-less Invalid", NetURL{URL: &url.URL{Path: "invalid"}}, ""},

		// Provide query/fragment via RawQuery/Fragment (don’t embed in Path)
		{"With Query", NetURL{URL: &url.URL{Path: "example.com", RawQuery: "query=1"}}, "https://example.com?query=1"},
		{"With Fragment", NetURL{URL: &url.URL{Path: "example.com", Fragment: "frag"}}, "https://example.com#frag"},

		// Bonus: protocol-relative with path, query, fragment
		{"Protocol-relative full", NetURL{URL: &url.URL{
			Host:     "example.com",
			Path:     "/x",
			RawQuery: "a=1",
			Fragment: "b",
		}}, "https://example.com/x?a=1#b"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NormalizeURL(tt.input)
			if got.String() != tt.expected {
				t.Errorf("NormalizeURL() = %q, want %q", got.String(), tt.expected)
			}
		})
	}
}
