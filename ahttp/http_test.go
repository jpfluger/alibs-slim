package ahttp

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestGetIsOnUnitTests tests the GetIsOnUnitTests function.
func TestGetIsOnUnitTests(t *testing.T) {
	ISON_UNITTESTS_WAIT_USER_SHUTDOWN = 1
	if !GetIsOnUnitTests() {
		t.Errorf("GetIsOnUnitTests() should return true when ISON_UNITTESTS_WAIT_USER_SHUTDOWN is not 0")
	}

	ISON_UNITTESTS_WAIT_USER_SHUTDOWN = 0
	if GetIsOnUnitTests() {
		t.Errorf("GetIsOnUnitTests() should return false when ISON_UNITTESTS_WAIT_USER_SHUTDOWN is 0")
	}
}

// TestGetIsOnUnitTestsHasSecret tests the GetIsOnUnitTestsHasSecret function.
func TestGetIsOnUnitTestsHasSecret(t *testing.T) {
	ISON_UNITTESTS_UPDOWN_SECRET = "secret"
	if !GetIsOnUnitTestsHasSecret() {
		t.Errorf("GetIsOnUnitTestsHasSecret() should return true when ISON_UNITTESTS_UPDOWN_SECRET is not empty")
	}

	ISON_UNITTESTS_UPDOWN_SECRET = ""
	if GetIsOnUnitTestsHasSecret() {
		t.Errorf("GetIsOnUnitTestsHasSecret() should return false when ISON_UNITTESTS_UPDOWN_SECRET is empty")
	}
}

// TestHasContentType tests the HasContentType function.
func TestHasContentType(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Content-Type", "application/json")

	if !HasContentType(req, "application/json") {
		t.Errorf("HasContentType() should return true when Content-Type header matches")
	}

	if HasContentType(req, "text/html") {
		t.Errorf("HasContentType() should return false when Content-Type header does not match")
	}
}

// TestExtractHttpRouteId tests the ExtractHttpRouteId function.
func TestExtractHttpRouteId(t *testing.T) {
	routeId, ok := ExtractHttpRouteId("hrt:RPAGE_LOGIN")
	if !ok || routeId != "RPAGE_LOGIN" {
		t.Errorf("ExtractHttpRouteId() failed to extract valid HttpRouteId")
	}

	_, ok = ExtractHttpRouteId("invalid")
	if ok {
		t.Errorf("ExtractHttpRouteId() should return false for invalid input")
	}
}

// TestJoinUrl tests the JoinUrl function.
func TestJoinUrl(t *testing.T) {
	urlRoot := "http://example.com"
	urlPath := "/path"

	joinedUrl := JoinUrl(urlRoot, urlPath)
	expectedUrl := "http://example.com/path"
	if joinedUrl != expectedUrl {
		t.Errorf("JoinUrl() = %v, want %v", joinedUrl, expectedUrl)
	}

	// Test with trailing slash in urlRoot and leading slash in urlPath
	urlRoot = "http://example.com/"
	urlPath = "/path"
	joinedUrl = JoinUrl(urlRoot, urlPath)
	if joinedUrl != expectedUrl {
		t.Errorf("JoinUrl() with slashes = %v, want %v", joinedUrl, expectedUrl)
	}
}
