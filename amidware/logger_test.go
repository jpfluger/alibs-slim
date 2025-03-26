package amidware

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/jpfluger/alibs-slim/alog"
)

// Mock logger for testing
type mockLogger struct{}

func (m mockLogger) Info() *zerolog.Event {
	return &zerolog.Event{}
}

func (m mockLogger) Err(err error) *zerolog.Event {
	return &zerolog.Event{}
}

// TestLoggerMiddleware verifies the middleware logging behavior
func TestLoggerMiddleware(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mw := LoggerWithConfig(DefaultLoggerConfig)
	handler := mw(func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	err := handler(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

// TestLoggerMiddlewareWithError verifies logging when handler returns an error
func TestLoggerMiddlewareWithError(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/error", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mw := LoggerWithConfig(DefaultLoggerConfig)
	handler := mw(func(c echo.Context) error {
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal Server Error")
	})

	err := handler(c)
	assert.Error(t, err)
	assert.Equal(t, http.StatusInternalServerError, err.(*echo.HTTPError).Code)
}

// TestLoggerMiddlewareWithCustomTag verifies custom logging functionality
func TestLoggerMiddlewareWithCustomTag(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/custom", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	config := LoggerConfig{
		LogChannel:   alog.LOGGER_HTTP,
		LogSeparator: &sync.Map{},
		CustomTagFunc: func(c echo.Context, event *zerolog.Event) {
			event.Str("custom_tag", "custom_value")
		},
	}

	mw := LoggerWithConfig(config)
	handler := mw(func(c echo.Context) error {
		return c.String(http.StatusCreated, "Created")
	})

	err := handler(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)
}

// TestLoggerMiddlewareWithLogSeparator verifies dynamic log separation
func TestLoggerMiddlewareWithLogSeparator(t *testing.T) {
	config := LoggerConfig{
		LogChannel:   alog.LOGGER_HTTP,
		LogSeparator: &sync.Map{},
	}

	// Add log separator for a specific IP
	AddLogSeparator(&config, "192.168.1.1", alog.LOGGER_APP)
	logChannel := getLogChannelForIP("192.168.1.1", config.LogSeparator, alog.LOGGER_HTTP)
	assert.Equal(t, alog.LOGGER_APP, logChannel)

	// Remove log separator and check it falls back to default
	RemoveLogSeparator(&config, "192.168.1.1")
	logChannel = getLogChannelForIP("192.168.1.1", config.LogSeparator, alog.LOGGER_HTTP)
	assert.Equal(t, alog.LOGGER_HTTP, logChannel)
}

// TestLoggerMiddlewareLatency verifies latency calculation in logs
func TestLoggerMiddlewareLatency(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/latency", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mw := LoggerWithConfig(DefaultLoggerConfig)
	handler := mw(func(c echo.Context) error {
		time.Sleep(10 * time.Millisecond)
		return c.String(http.StatusOK, "OK")
	})

	start := time.Now()
	err := handler(c)
	stop := time.Now()

	latency := stop.Sub(start)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.GreaterOrEqual(t, latency.Milliseconds(), int64(10))
}

// TestLoggerMiddlewareContentLength verifies logging of request content length
func TestLoggerMiddlewareContentLength(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/contentlength", nil)
	req.Header.Set(echo.HeaderContentLength, strconv.Itoa(256))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mw := LoggerWithConfig(DefaultLoggerConfig)
	handler := mw(func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	err := handler(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

// TestInitLogSeparator verifies that InitLogSeparator correctly initializes LogSeparator with multiple IP mappings.
func TestInitLogSeparator(t *testing.T) {
	// Create a LoggerConfig instance
	config := LoggerConfig{
		LogChannel:   alog.LOGGER_HTTP, // Default log channel
		LogSeparator: &sync.Map{},
	}

	// Define an IPLogChannelMap with multiple mappings
	ipLogMap := IPLogChannelMap{
		"192.168.1.1": alog.LOGGER_APP,
		"10.0.0.2":    alog.LOGGER_SQL,
		"172.16.0.3":  alog.LOGGER_HTTP,
	}

	// Call InitLogSeparator to populate LogSeparator
	InitLogSeparator(&config, ipLogMap)

	// Verify that the IPs have the correct log channels assigned
	for ip, expectedLogChannel := range ipLogMap {
		actualLogChannel := getLogChannelForIP(ip, config.LogSeparator, alog.LOGGER_HTTP)
		assert.Equal(t, expectedLogChannel, actualLogChannel, "Log channel mismatch for IP: %s", ip)
	}

	// Verify that an unknown IP falls back to the default log channel
	unknownIP := "8.8.8.8"
	defaultLogChannel := getLogChannelForIP(unknownIP, config.LogSeparator, alog.LOGGER_HTTP)
	assert.Equal(t, alog.LOGGER_HTTP, defaultLogChannel, "Unknown IP should return the default log channel")
}
