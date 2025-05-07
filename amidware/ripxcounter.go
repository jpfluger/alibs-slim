package amidware

import (
	"net/http"
	"sync"
	"time"

	"github.com/jpfluger/alibs-slim/alog"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type TimeRange struct {
	Start time.Time
	End   time.Time
}

type IPExclusion struct {
	IP            string
	AlwaysExclude bool
	TimeRanges    []TimeRange

	// Internal runtime state (unexported, not serialized)
	counter       int
	windowStart   time.Time
	hasLoggedOnce bool
}

func (ex *IPExclusion) ResetWindow() {
	ex.windowStart = time.Time{}
	ex.counter = 0
	ex.hasLoggedOnce = false
}

func (ex *IPExclusion) IncrementCounter() {
	ex.counter++
}

func (ex *IPExclusion) GetCounter() int {
	return ex.counter
}

type RIPXCounterConfig struct {
	Skipper           middleware.Skipper
	IsOnExclusions    bool
	Exclusions        map[string]*IPExclusion
	IsOnGeneralCounts bool
	GeneralCounts     map[string]int
	Mutex             sync.Mutex
	LogChannel        alog.ChannelLabel
	FlushInterval     int // in seconds
}

func (cfg *RIPXCounterConfig) FlushDuration() time.Duration {
	return time.Duration(cfg.FlushInterval) * time.Second
}

const LOGGER_RIPXC alog.ChannelLabel = "ripxc"

var DefaultRIPXCounterConfig = RIPXCounterConfig{
	Skipper:           middleware.DefaultSkipper,
	LogChannel:        LOGGER_RIPXC,
	IsOnExclusions:    true,
	IsOnGeneralCounts: true,
}

func RIPXCounterMiddleware(config *RIPXCounterConfig) echo.MiddlewareFunc {
	if config.Skipper == nil {
		config.Skipper = middleware.DefaultSkipper
	}
	if config.Exclusions == nil {
		config.Exclusions = map[string]*IPExclusion{}
	}
	if config.GeneralCounts == nil {
		config.GeneralCounts = map[string]int{}
	}
	if config.FlushInterval < 1 {
		config.FlushInterval = 60 // fallback to 1 second
	}
	hasLogger := !config.LogChannel.IsEmpty()

	go startFlushScheduler(config)

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if config.Skipper(c) {
				return next(c)
			}

			ip := c.RealIP()
			now := time.Now()

			config.Mutex.Lock()

			// General count tracking
			if config.IsOnGeneralCounts {
				config.GeneralCounts[ip]++
			}

			// Exclusion tracking
			if config.IsOnExclusions {
				if exclusion, ok := config.Exclusions[ip]; ok {
					doExclusion := false
					if exclusion.AlwaysExclude {
						exclusion.counter++
						if hasLogger && !exclusion.hasLoggedOnce {
							alog.LOGGER(config.LogChannel).
								Info().
								Str("ip", ip).
								Str("mode", "always").
								Int("req-counter", exclusion.counter).
								Msg("ip always excluded")
							exclusion.hasLoggedOnce = true
						}
						config.Mutex.Unlock()
						doExclusion = true
					}

					if !doExclusion {
						for _, tr := range exclusion.TimeRanges {
							start := time.Date(now.Year(), now.Month(), now.Day(), tr.Start.Hour(), tr.Start.Minute(), 0, 0, time.Local)
							end := time.Date(now.Year(), now.Month(), now.Day(), tr.End.Hour(), tr.End.Minute(), 0, 0, time.Local)

							if now.After(start) && now.Before(end) {
								if exclusion.windowStart != start {
									exclusion.windowStart = start
									exclusion.counter = 0
									exclusion.hasLoggedOnce = false
								}
								exclusion.counter++
								if hasLogger && !exclusion.hasLoggedOnce {
									alog.LOGGER(config.LogChannel).
										Info().
										Str("ip", ip).
										Str("start", tr.Start.Format("15:04")).
										Str("end", tr.End.Format("15:04")).
										Int("req-counter", 0).
										Msg("ip entered exclusion window")
									exclusion.hasLoggedOnce = true
								}
								config.Mutex.Unlock()
							}
							doExclusion = true
						}
					}
					if doExclusion {
						return c.JSON(http.StatusForbidden, map[string]string{"error": "restricted"})
					}
				}
			}

			config.Mutex.Unlock()
			return next(c)
		}
	}
}

func startFlushScheduler(config *RIPXCounterConfig) {
	for {
		next := time.Now().Add(config.FlushDuration())
		time.Sleep(time.Until(next))
		FlushIPStats(config)
	}
}

func FlushIPStats(config *RIPXCounterConfig) {
	config.Mutex.Lock()
	defer config.Mutex.Unlock()

	hasLogger := !config.LogChannel.IsEmpty()
	if !hasLogger {
		return
	}

	// Track IPs already logged (to avoid duplicates)
	logged := map[string]bool{}

	// Step 1: Handle exclusions first
	if config.IsOnExclusions {
		for ip, ex := range config.Exclusions {
			if ex.hasLoggedOnce && ex.counter > 0 {
				startTime := ""
				stopTime := ""
				if !ex.windowStart.IsZero() {
					startTime = ex.windowStart.Format("15:04")
					stopTime = ex.windowStart.Add(config.FlushDuration()).Format("15:04")
				}

				generalCount := config.GeneralCounts[ip]

				alog.LOGGER(config.LogChannel).
					Info().
					Str("ip", ip).
					Int("req-counter", generalCount).
					Int("exclusion-counter", ex.counter).
					Str("start", startTime).
					Str("stop", stopTime).
					Msg("flush combined ip req-counters")

				// Reset exclusion state
				ex.ResetWindow()

				// Mark as logged
				logged[ip] = true
				delete(config.GeneralCounts, ip)
			}
		}
	}

	// Step 2: Log remaining general counts (IPs not in exclusion)
	if config.IsOnGeneralCounts {
		for ip, count := range config.GeneralCounts {
			if logged[ip] {
				continue
			}
			alog.LOGGER(config.LogChannel).
				Info().
				Str("ip", ip).
				Int("req-counter", count).
				Msg("flush general ip req-counter")
		}
		// Reset general counts
		config.GeneralCounts = map[string]int{}
	}
}

// RIPXCounterOpts defines user-facing configuration for the RIPX middleware.
// This is intended to be filled from JSON/YAML or environment variables.
type RIPXCounterOpts struct {
	Enabled             bool           `json:"enabled"`                // Enable/disable RIPX middleware
	FlushInterval       int            `json:"flushInterval"`          // Interval (in seconds) to flush counters
	EnableExclusions    bool           `json:"enableExclusions"`       // Enable time-based exclusions
	EnableGeneralCounts bool           `json:"enableGeneralCounts"`    // Enable general per-IP counting
	IPExclusions        []*IPExclusion `json:"ipExclusions,omitempty"` // Full IP exclusion definition including time windows
}

func (opts *RIPXCounterOpts) ToConfig(logChannel alog.ChannelLabel) *RIPXCounterConfig {
	if opts == nil || !opts.Enabled || logChannel.IsEmpty() {
		return nil
	}

	return &RIPXCounterConfig{
		Skipper:           middleware.DefaultSkipper,
		IsOnExclusions:    opts.EnableExclusions,
		IsOnGeneralCounts: opts.EnableGeneralCounts,
		FlushInterval:     opts.FlushInterval,
		LogChannel:        logChannel,
		GeneralCounts:     make(map[string]int),
		Exclusions:        buildExclusionsFromOpts(opts),
	}
}

func buildExclusionsFromOpts(opts *RIPXCounterOpts) map[string]*IPExclusion {
	exclusions := make(map[string]*IPExclusion)

	for _, ipEx := range opts.IPExclusions {
		if ipEx == nil || ipEx.IP == "" {
			continue
		}

		// Sanitize and copy valid time ranges
		var validRanges []TimeRange
		for _, tr := range ipEx.TimeRanges {
			if !tr.Start.IsZero() && !tr.End.IsZero() && tr.End.After(tr.Start) {
				validRanges = append(validRanges, tr)
			}
		}

		if !ipEx.AlwaysExclude && len(validRanges) == 0 {
			continue
		}

		// Deep copy into new struct (so it’s isolated from config)
		exclusions[ipEx.IP] = &IPExclusion{
			IP:            ipEx.IP,
			AlwaysExclude: ipEx.AlwaysExclude,
			TimeRanges:    validRanges,
			counter:       0,
			windowStart:   time.Time{},
			hasLoggedOnce: false,
		}
	}

	return exclusions
}
