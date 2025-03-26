package ahttp

import (
	"github.com/alexedwards/scs/v2"
	"github.com/labstack/echo/v4"
	"github.com/jpfluger/alibs-slim/arob"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// TestForceRedirectRawQuery tests the ForceRedirectRawQuery function.
func TestForceRedirectRawQuery(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/?param=value", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := ForceRedirectRawQuery(c, "http://example.com")
	if err != nil {
		t.Errorf("ForceRedirectRawQuery() returned an error: %v", err)
	}
	expectedLocation := "http://example.com/?param=value"
	if rec.Header().Get("Location") != expectedLocation {
		t.Errorf("ForceRedirectRawQuery() should redirect to %v, got %v", expectedLocation, rec.Header().Get("Location"))
	}
}

// TestForceRedirectRequestToBase64 tests the ForceRedirectRequestToBase64 function.
func TestForceRedirectRequestToBase64(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := ForceRedirectRequestToBase64(c, "http://example.com")
	if err != nil {
		t.Errorf("ForceRedirectRequestToBase64() returned an error: %v", err)
	}
	// Check if the Location header contains the base64 query parameter.
	if !strings.Contains(rec.Header().Get("Location"), "burl=") {
		t.Errorf("ForceRedirectRequestToBase64() should contain a base64 query parameter")
	}
}

func TestResetUserSessionWithRedirect(t *testing.T) {
	e := echo.New()
	manager := scs.New()
	manager.Lifetime = 1 * time.Hour
	manager.Cookie.Name = "sessionid"

	// Create a middleware function that loads and saves the session data.
	sessionMiddleware := echo.WrapMiddleware(manager.LoadAndSave)

	// Wrap the handler with the session middleware.
	handler := func(c echo.Context) error {
		return ResetUserSessionWithRedirect(c, manager, "http://example.com")
	}
	wrappedHandler := sessionMiddleware(handler)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Call the wrapped handler, which includes the session middleware.
	if err := wrappedHandler(c); err != nil {
		t.Errorf("ResetUserSessionWithRedirect() returned an error: %v", err)
	}

	expectedLocation := "http://example.com"
	if rec.Header().Get("Location") != expectedLocation {
		t.Errorf("ResetUserSessionWithRedirect() should redirect to %v, got %v", expectedLocation, rec.Header().Get("Location"))
	}
}

func TestResetUserSessionWithRedirectJSONMessage(t *testing.T) {
	e := echo.New()
	manager := scs.New()
	manager.Lifetime = 1 * time.Hour
	manager.Cookie.Name = "sessionid"

	// Create a middleware function that loads and saves the session data.
	sessionMiddleware := echo.WrapMiddleware(manager.LoadAndSave)

	// Wrap the handler with the session middleware.
	handler := func(c echo.Context) error {
		message := arob.NewROBMessage("Test message")
		return ResetUserSessionWithRedirectJSONMessage(c, manager, "http://example.com", message)
	}
	wrappedHandler := sessionMiddleware(handler)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Call the wrapped handler, which includes the session middleware.
	if err := wrappedHandler(c); err != nil {
		t.Errorf("ResetUserSessionWithRedirectJSONMessage() returned an error: %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Errorf("ResetUserSessionWithRedirectJSONMessage() should respond with status OK for JSON requests")
	}
}
