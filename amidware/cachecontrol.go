package amidware

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"strconv"
)

// CacheControlConfig defines configurable cache headers.
type CacheControlConfig struct {
	NoStore        bool
	NoCache        bool
	MustRevalidate bool
	PragmaNoCache  bool
	Expires        string
	MaxAge         int  // in seconds
	ETagSupport    bool // If true, sets a static weak ETag header

	headers map[string]string // precomputed headers
}

// SetHeaders precomputes the headers and optionally merges additional ones.
// If extra has keys that match existing ones, extra overrides them.
func (cfg *CacheControlConfig) SetHeaders(extra ...map[string]string) {
	headers := make(map[string]string)
	var cacheControl []string

	if cfg.NoStore {
		cacheControl = append(cacheControl, "no-store")
	}
	if cfg.NoCache {
		cacheControl = append(cacheControl, "no-cache")
	}
	if cfg.MustRevalidate {
		cacheControl = append(cacheControl, "must-revalidate")
	}
	if cfg.MaxAge > 0 {
		cacheControl = append(cacheControl, fmt.Sprintf("max-age=%d", cfg.MaxAge))
	}

	if len(cacheControl) > 0 {
		headers["Cache-Control"] = joinDirectives(cacheControl)
	}
	if cfg.PragmaNoCache {
		headers["Pragma"] = "no-cache"
	}
	if cfg.Expires != "" {
		headers["Expires"] = cfg.Expires
	}
	if cfg.ETagSupport {
		headers["ETag"] = fmt.Sprintf(`W/"static-%s"`, strconv.Itoa(cfg.MaxAge))
	}

	// Merge any provided additional headers
	if len(extra) > 0 {
		for k, v := range extra[0] {
			headers[k] = v // extra overrides
		}
	}

	cfg.headers = headers
}

// Helper to join directives
func joinDirectives(directives []string) string {
	result := directives[0]
	for _, d := range directives[1:] {
		result += ", " + d
	}
	return result
}

// GetHeaders returns precomputed headers.
func (cfg CacheControlConfig) GetHeaders() map[string]string {
	return cfg.headers
}

// DefaultCacheControlConfig returns strict no-cache config.
func DefaultCacheControlConfig() CacheControlConfig {
	cfg := CacheControlConfig{
		NoStore:        true,
		NoCache:        true,
		MustRevalidate: true,
		PragmaNoCache:  true,
		Expires:        "0",
		MaxAge:         0,
		ETagSupport:    false,
	}
	cfg.SetHeaders()
	return cfg
}

// CacheControlMiddleware returns Echo middleware with precomputed headers.
// The default applies global middleware to disable client-side caching for all routes.
// This sets headers like Cache-Control, Pragma, and Expires to prevent the browser from caching pages.
//
// ✅ Use cases:
// - Security-sensitive apps (e.g. banking, finance) where cached pages must not be revisited
// - After logout to prevent the user from going back to a sensitive page via the back button
// - Dynamic UIs (e.g. dashboards) where stale data must never be shown
// - Multi-user environments where data should not persist across sessions
func CacheControlMiddleware(cfg CacheControlConfig) echo.MiddlewareFunc {
	headers := cfg.GetHeaders()

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			for k, v := range headers {
				c.Response().Header().Set(k, v)
			}
			return next(c)
		}
	}
}
