package ahttp

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/jpfluger/alibs-slim/azb"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
)

// TestRHBind tests the RHBind function for POST requests.
func TestRHBind(t *testing.T) {
	e := echo.New()
	data := `{"key":"value"}`
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(data))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	var din map[string]interface{}

	// Handler that uses RHBind
	handler := func(c echo.Context) error {
		if err := RHBind(c, &din); err != nil {
			return err
		}
		val, exists := din["key"]
		if !exists || val != "value" {
			t.Errorf("Expected data to be bound with 'key': 'value', got: %v", din)
		}
		return c.NoContent(http.StatusOK)
	}

	// Run the handler and test the results
	if assert.NoError(t, handler(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		// The following line is incorrect because it's trying to get the value from the context directly
		// assert.Equal(t, "value", c.Get("key"))
		// Instead, you should assert the value from the 'din' map
		assert.Equal(t, "value", din["key"])
	}
}

// TestRHBindNopCloser tests the RHBindNopCloser function for POST requests.
func TestRHBindNopCloser(t *testing.T) {
	e := echo.New()
	data := `{"key":"value"}`
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(data))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Handler that uses RHBindNopCloser
	handler := func(c echo.Context) error {
		var din map[string]interface{}
		if err := RHBindNopCloser(c, &din); err != nil {
			return err
		}
		val, exists := din["key"]
		if !exists || val != "value" {
			t.Errorf("Expected data to be bound with 'key': 'value', got: %v", din)
		}

		// Attempt to read the body again to ensure it was restored
		bodyBytes, _ := io.ReadAll(c.Request().Body)
		if string(bodyBytes) != data {
			t.Errorf("Expected request body to be restored, got: %s", string(bodyBytes))
		}

		return c.NoContent(http.StatusOK)
	}

	// Run the handler and test the results
	if assert.NoError(t, handler(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}

// TestRHBindListPostOnly tests the RHBindListPostOnly function for POST requests.
func TestRHBindListPostOnly(t *testing.T) {
	e := echo.New()
	data := `{"page":1,"pageSize":10}`
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(data))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Handler that uses RHBindListPostOnly
	handler := func(c echo.Context) error {
		var din azb.DIN
		if err := RHBindListPostOnly(c, &din); err != nil {
			return err
		}
		// Add assertions to check if the data is correctly bound
		// Example:
		// assert.Equal(t, 1, din.Page)
		// assert.Equal(t, 10, din.PageSize)
		return c.NoContent(http.StatusOK)
	}

	// Run the handler and test the results for a POST request
	if assert.NoError(t, handler(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}

	// Test for a non-POST request (e.g., GET)
	req = httptest.NewRequest(http.MethodGet, "/", nil)
	c = e.NewContext(req, rec)
	if err := handler(c); err != nil {
		t.Errorf("Expected no error for non-POST request, got: %v", err)
	}
}

// Mock implementation of IDINPaginate for testing purposes
type MockIDINPaginate struct {
	azb.DIN
	Page     int `query:"page"`
	PageSize int `query:"pageSize"`
}

// Validate mock implementation (replace with actual validation logic if necessary)
func (m *MockIDINPaginate) Validate() error {
	if m.Page < 1 {
		return fmt.Errorf("page must be greater than 0")
	}
	if m.PageSize < 1 || m.PageSize > 100 {
		return fmt.Errorf("pageSize must be between 1 and 100")
	}
	return nil
}

// TestRHBindListQuery tests the RHBindListQuery function for binding and validating query parameters.
func TestRHBindListQuery(t *testing.T) {
	e := echo.New()
	q := make(url.Values)
	q.Set("page", "1")
	q.Set("pageSize", "10")
	req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Handler that uses RHBindListQuery
	handler := func(c echo.Context) error {
		var din MockIDINPaginate // Replace with the actual type
		if err := RHBindListQuery(c, &din); err != nil {
			return err
		}
		// Add assertions to check if the data is correctly bound and validated
		assert.Equal(t, 1, din.Page)
		assert.Equal(t, 10, din.PageSize)
		return c.NoContent(http.StatusOK)
	}

	// Run the handler and test the results
	if assert.NoError(t, handler(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}
