package ahttp

import (
	"encoding/base64"
	"encoding/json"
	"github.com/labstack/echo/v4"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestRedirectFlashMessage tests the RedirectFlashMessage function.
func TestRedirectFlashMessage(t *testing.T) {
	// Create a new Echo instance for testing.
	e := echo.New()

	// Create a new request and recorder.
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Define the FlashMessageData to be used in the test.
	fmData := &FlashMessageData{
		Path:     "/testpath",
		Action:   "testaction",
		Username: "testuser",
	}

	// Define the cookie name and target URL for the test.
	cookieName := COOKIE_FLASH_FORGOT_LOGIN
	targetUrl := "/target"

	// Call the RedirectFlashMessage function.
	err := RedirectFlashMessage(c, cookieName, fmData, targetUrl)
	if err != nil {
		t.Errorf("RedirectFlashMessage returned an error: %v", err)
	}

	// Retrieve the 'Set-Cookie' header from the response.
	setCookieHeader := rec.Header().Get("Set-Cookie")
	if setCookieHeader == "" {
		t.Errorf("Set-Cookie header not found")
		return
	}

	// Parse the 'Set-Cookie' header to get the cookie value.
	cookieValue := ""
	parts := strings.Split(setCookieHeader, ";")
	for _, part := range parts {
		if strings.HasPrefix(part, cookieName.String()+"=") {
			cookieValue = strings.TrimPrefix(part, cookieName.String()+"=")
			break
		}
	}
	if cookieValue == "" {
		t.Errorf("Cookie %s value not found", cookieName.String())
		return
	}

	// Decode the cookie value.
	decoded, err := base64.URLEncoding.DecodeString(cookieValue)
	if err != nil {
		t.Errorf("Failed to decode cookie value: %v", err)
		return
	}

	// Unmarshal the JSON data from the cookie.
	var decodedData FlashMessageData
	err = json.Unmarshal(decoded, &decodedData)
	if err != nil {
		t.Errorf("Failed to unmarshal cookie value: %v", err)
		return
	}

	// Assert that the decoded data matches the original data.
	if decodedData.Path != fmData.Path {
		t.Errorf("Expected Path to be %v, got %v", fmData.Path, decodedData.Path)
	}
	if decodedData.Action != fmData.Action {
		t.Errorf("Expected Action to be %v, got %v", fmData.Action, decodedData.Action)
	}
	if decodedData.Username != fmData.Username {
		t.Errorf("Expected Username to be %v, got %v", fmData.Username, decodedData.Username)
	}

	// Assert that the cookie has the correct attributes.
	if !strings.Contains(setCookieHeader, "Path=/") {
		t.Errorf("Expected cookie Path to be '/', got %v", setCookieHeader)
	}
	if !strings.Contains(setCookieHeader, "HttpOnly") {
		t.Errorf("Expected cookie HttpOnly to be set")
	}

	// Assert that the redirect location is correct.
	location := rec.Header().Get("Location")
	if location != targetUrl {
		t.Errorf("Expected redirect location to be %v, got %v", targetUrl, location)
	}

	// Assert that the status code is correct.
	if status := rec.Code; status != http.StatusFound {
		t.Errorf("Expected status code to be %v, got %v", http.StatusFound, status)
	}
}
