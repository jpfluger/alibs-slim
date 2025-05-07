package amidware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestDefaultCacheControlConfig(t *testing.T) {
	cfg := DefaultCacheControlConfig()
	expected := map[string]string{
		"Cache-Control": "no-store, no-cache, must-revalidate",
		"Pragma":        "no-cache",
		"Expires":       "0",
	}

	for k, v := range expected {
		assert.Equal(t, v, cfg.GetHeaders()[k], "Header %s should match", k)
	}
}

func TestCacheControlWithMaxAgeAndETag(t *testing.T) {
	cfg := CacheControlConfig{
		MaxAge:      3600,
		ETagSupport: true,
	}
	cfg.SetHeaders()
	headers := cfg.GetHeaders()

	assert.Equal(t, "max-age=3600", headers["Cache-Control"])
	assert.Equal(t, `W/"static-3600"`, headers["ETag"])
}

func TestCacheControlHeaderOverride(t *testing.T) {
	cfg := CacheControlConfig{
		MaxAge:      600,
		NoStore:     true,
		ETagSupport: true,
	}
	extra := map[string]string{
		"Cache-Control": "public, max-age=900",
		"X-Custom":      "yes",
	}
	cfg.SetHeaders(extra)
	headers := cfg.GetHeaders()

	assert.Equal(t, "public, max-age=900", headers["Cache-Control"])
	assert.Equal(t, "yes", headers["X-Custom"])
}

func TestCacheControlMiddlewareIntegration(t *testing.T) {
	e := echo.New()
	cfg := DefaultCacheControlConfig()
	mw := CacheControlMiddleware(cfg)

	e.GET("/test", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello")
	}, mw)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	e.Router().Find(http.MethodGet, "/test", c)
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "no-store, no-cache, must-revalidate", rec.Header().Get("Cache-Control"))
	assert.Equal(t, "no-cache", rec.Header().Get("Pragma"))
	assert.Equal(t, "0", rec.Header().Get("Expires"))
}
