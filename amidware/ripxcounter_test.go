package amidware

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func newEchoWithRIPX(config *RIPXCounterConfig) *echo.Echo {
	e := echo.New()
	e.Use(RIPXCounterMiddleware(config))
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})
	return e
}

func makeRequest(e *echo.Echo, ip string) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.RemoteAddr = ip + ":12345"
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
}

func TestGeneralRequestCounting(t *testing.T) {
	config := &RIPXCounterConfig{
		Exclusions:        make(map[string]*IPExclusion),
		GeneralCounts:     make(map[string]int),
		LogChannel:        LOGGER_RIPXC,
		FlushInterval:     3600, // 1 hour in seconds
		IsOnGeneralCounts: true,
	}
	e := newEchoWithRIPX(config)

	ip := "192.0.2.1"
	for i := 0; i < 5; i++ {
		makeRequest(e, ip)
	}

	config.Mutex.Lock()
	defer config.Mutex.Unlock()
	assert.Equal(t, 5, config.GeneralCounts[ip])
}

func TestExclusionCountingWithinWindow(t *testing.T) {
	now := time.Now()
	start := now.Add(-1 * time.Minute)
	end := now.Add(1 * time.Minute)

	ip := "192.0.2.2"
	config := &RIPXCounterConfig{
		Exclusions: map[string]*IPExclusion{
			ip: {
				IP: ip,
				TimeRanges: []TimeRange{
					{Start: start, End: end},
				},
			},
		},
		GeneralCounts:     make(map[string]int),
		LogChannel:        LOGGER_RIPXC,
		FlushInterval:     3600, // 1 hour in seconds
		IsOnExclusions:    true,
		IsOnGeneralCounts: true,
	}
	e := newEchoWithRIPX(config)

	for i := 0; i < 3; i++ {
		makeRequest(e, ip)
	}

	config.Mutex.Lock()
	defer config.Mutex.Unlock()
	ex := config.Exclusions[ip]
	assert.Equal(t, 3, ex.counter)
	assert.True(t, ex.hasLoggedOnce)
}

func TestFlushResetsCounters(t *testing.T) {
	ip1 := "192.0.2.10"
	ip2 := "192.0.2.20"

	config := &RIPXCounterConfig{
		Exclusions: map[string]*IPExclusion{
			ip2: {
				IP: ip2,
				TimeRanges: []TimeRange{
					{
						Start: time.Now().Add(-1 * time.Minute),
						End:   time.Now().Add(1 * time.Minute),
					},
				},
			},
		},
		GeneralCounts:     map[string]int{ip1: 4},
		LogChannel:        LOGGER_RIPXC,
		FlushInterval:     3600, // 1 hour in seconds
		IsOnExclusions:    true,
		IsOnGeneralCounts: true,
	}

	config.Exclusions[ip2].counter = 7
	config.Exclusions[ip2].hasLoggedOnce = true
	config.Exclusions[ip2].windowStart = time.Now()

	FlushIPStats(config)

	config.Mutex.Lock()
	defer config.Mutex.Unlock()

	assert.Equal(t, 0, config.GeneralCounts[ip1])
	assert.Equal(t, 0, config.Exclusions[ip2].counter)
	assert.False(t, config.Exclusions[ip2].hasLoggedOnce)
	assert.True(t, config.Exclusions[ip2].windowStart.IsZero())
}

func TestConcurrentRequests(t *testing.T) {
	config := &RIPXCounterConfig{
		Exclusions:        make(map[string]*IPExclusion),
		GeneralCounts:     make(map[string]int),
		LogChannel:        LOGGER_RIPXC,
		FlushInterval:     3600, // 1 hour in seconds
		IsOnGeneralCounts: true,
	}
	e := newEchoWithRIPX(config)

	ip := "192.0.2.3"
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			makeRequest(e, ip)
			wg.Done()
		}()
	}
	wg.Wait()

	config.Mutex.Lock()
	defer config.Mutex.Unlock()
	assert.Equal(t, 100, config.GeneralCounts[ip])
}
