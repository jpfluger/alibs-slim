package ahttp

import (
	"errors"
	"github.com/labstack/echo/v4"
	"net/http"
	"net/http/httptest"
	"testing"
)

// mockEchoContext is a helper function to create an Echo context with a recorder.
func mockEchoContext(httpMethod, target string) (echo.Context, *httptest.ResponseRecorder) {
	e := echo.New()
	req := httptest.NewRequest(httpMethod, target, nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	return c, rec
}

// TestNewHTTPErrorHandlerBase tests the creation of a new HTTPErrorHandlerBase.
func TestNewHTTPErrorHandlerBase(t *testing.T) {
	c, _ := mockEchoContext(http.MethodGet, "/")
	logger := echo.New().Logger
	isOnDebug := true
	err := errors.New("test error")

	handler := NewHTTPErrorHandlerBase(err, c, logger, isOnDebug)

	if handler.GetErr() != err {
		t.Errorf("Expected error to be %v, got %v", err, handler.GetErr())
	}
	if handler.GetContext() != c {
		t.Errorf("Expected context to be %v, got %v", c, handler.GetContext())
	}
	if handler.GetLogger() != logger {
		t.Errorf("Expected logger to be %v, got %v", logger, handler.GetLogger())
	}
	if handler.GetIsOnDebug() != isOnDebug {
		t.Errorf("Expected isOnDebug to be %v, got %v", isOnDebug, handler.GetIsOnDebug())
	}
}

// MockHTTPErrorHandler implements the IHTTPErrorHandler interface for testing.
type MockHTTPErrorHandler struct {
	*HTTPErrorHandlerBase
}

// NewMockHTTPErrorHandler creates a new instance of MockHTTPErrorHandler.
func NewMockHTTPErrorHandler(err error, c echo.Context, logger echo.Logger, isOnDebug bool) *MockHTTPErrorHandler {
	return &MockHTTPErrorHandler{
		HTTPErrorHandlerBase: NewHTTPErrorHandlerBase(err, c, logger, isOnDebug),
	}
}

// HandleResponse is a mock method to satisfy the IHTTPErrorHandler interface.
func (m *MockHTTPErrorHandler) HandleResponse() error {
	// Mock response handling logic here.
	// For example, return nil to simulate a successful handling of the response.
	//return nil
	return m.c.String(m.httpCode, m.httpMessage)
}

// TestDefaultHTTPErrorHandler tests the DefaultHTTPErrorHandler function.
func TestDefaultHTTPErrorHandler(t *testing.T) {
	c, rec := mockEchoContext(http.MethodGet, "/")
	logger := echo.New().Logger
	isOnDebug := false
	err := echo.NewHTTPError(http.StatusNotFound, "resource not found")

	mockHandler := NewMockHTTPErrorHandler(err, c, logger, isOnDebug)

	DefaultHTTPErrorHandler(mockHandler)

	if rec.Code != http.StatusNotFound {
		t.Errorf("Expected status code to be %v, got %v", http.StatusNotFound, rec.Code)
	}
}
