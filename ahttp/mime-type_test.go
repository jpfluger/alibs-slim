package ahttp

import (
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestGetRequestContentType tests the GetRequestContentType function.
func TestGetRequestContentType(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	contentType := GetRequestContentType(c)
	if contentType != echo.MIMEApplicationJSON {
		t.Errorf("GetRequestContentType() = %v, want %v", contentType, echo.MIMEApplicationJSON)
	}
}

// TestIsRequestContentType tests the IsRequestContentType function.
func TestIsRequestContentType(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if !IsRequestContentType(c, CHECK_MIME_TYPE_JSON) {
		t.Errorf("IsRequestContentType() should return true for JSON content type")
	}

	if IsRequestContentType(c, CHECK_MIME_TYPE_XML) {
		t.Errorf("IsRequestContentType() should return false for non-XML content type")
	}
}

// TestRedirectCheckJSON tests the RedirectCheckJSON function.
func TestRedirectCheckJSON(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := RedirectCheckJSON(c, "http://example.com")
	if err != nil {
		t.Errorf("RedirectCheckJSON() returned an error: %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Errorf("RedirectCheckJSON() should respond with status OK for JSON requests")
	}

	// Test with non-JSON request
	req = httptest.NewRequest(http.MethodGet, "/", nil)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)

	err = RedirectCheckJSON(c, "http://example.com")
	if err != nil {
		t.Errorf("RedirectCheckJSON() returned an error: %v", err)
	}

	if rec.Code != http.StatusFound {
		t.Errorf("RedirectCheckJSON() should redirect for non-JSON requests")
	}
}

func TestDetectMimeTypeSendMessage(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := DetectMimeTypeSendMessage(c, 0, "test message", false)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "test message", rec.Body.String())
}

func TestDetectMimeTypeSendMessageWithCode(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	req.Header.Set(echo.HeaderContentType, "application/json")
	err := DetectMimeTypeSendMessageWithCode(c, http.StatusCreated, "test message", false)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)
	assert.Contains(t, rec.Body.String(), "test message")

	req.Header.Set(echo.HeaderContentType, "application/xml")
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	err = DetectMimeTypeSendMessageWithCode(c, http.StatusCreated, "test message", false)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)
	assert.Contains(t, rec.Body.String(), "test message")

	req.Header.Set(echo.HeaderContentType, "text/plain")
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	err = DetectMimeTypeSendMessageWithCode(c, http.StatusCreated, "test message", false)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)
	assert.Equal(t, "test message", rec.Body.String())
}

func TestDetectMimeTypeSendError(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := DetectMimeTypeSendError(c, 0, errors.New("test error"))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Equal(t, "test error", rec.Body.String())
}

func TestDetectMimeTypeSendErrorWithCode(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	req.Header.Set(echo.HeaderContentType, "application/json")
	err := DetectMimeTypeSendErrorWithCode(c, http.StatusBadRequest, errors.New("test error"))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "test error")

	req.Header.Set(echo.HeaderContentType, "application/xml")
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	err = DetectMimeTypeSendErrorWithCode(c, http.StatusBadRequest, errors.New("test error"))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "test error")

	req.Header.Set(echo.HeaderContentType, "text/plain")
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	err = DetectMimeTypeSendErrorWithCode(c, http.StatusBadRequest, errors.New("test error"))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Equal(t, "test error", rec.Body.String())
}
