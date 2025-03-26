package amidware

import (
	"github.com/jpfluger/alibs-slim/alog"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
)

// LoggerConfig defines the config for Logger middleware.
type LoggerConfig struct {
	// Skipper defines a function to skip middleware.
	Skipper middleware.Skipper

	// Zerolog logger instance to be used for structured logging.
	LogChannel alog.ChannelLabel

	// Thread-safe dynamic log separation.
	LogSeparator *sync.Map

	// CustomTagFunc is a function to handle `${custom}` tag.
	CustomTagFunc func(c echo.Context, event *zerolog.Event)
}

// DefaultLoggerConfig is the default Logger middleware config.
var DefaultLoggerConfig = LoggerConfig{
	Skipper:      middleware.DefaultSkipper,
	LogChannel:   alog.LOGGER_HTTP,
	LogSeparator: &sync.Map{},
}

// Logger returns a middleware that logs HTTP requests using zerolog.
func Logger() echo.MiddlewareFunc {
	return LoggerWithConfig(DefaultLoggerConfig)
}

// LoggerWithIPLogChannelMap returns a middleware that logs HTTP requests using zerolog.
func LoggerWithIPLogChannelMap(ipLogMap IPLogChannelMap) echo.MiddlewareFunc {
	config := LoggerConfig{
		Skipper:      middleware.DefaultSkipper,
		LogChannel:   alog.LOGGER_HTTP,
		LogSeparator: &sync.Map{},
	}
	InitLogSeparator(&config, ipLogMap)
	return LoggerWithConfig(config)
}

// LoggerWithConfig returns a Logger middleware with config.
func LoggerWithConfig(config LoggerConfig) echo.MiddlewareFunc {
	if config.Skipper == nil {
		config.Skipper = DefaultLoggerConfig.Skipper
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			if config.Skipper(c) {
				return next(c)
			}

			req := c.Request()
			res := c.Response()
			start := time.Now()

			err = next(c)
			if err != nil {
				c.Error(err)
			}

			stop := time.Now()
			latency := stop.Sub(start)

			// Get Real IP and determine the log channel dynamically
			realIP := c.RealIP()
			logChannel := getLogChannelForIP(realIP, config.LogSeparator, config.LogChannel)

			// Start building the zerolog event
			var event *zerolog.Event
			if err == nil {
				event = alog.LOGGER(logChannel).Info()
			} else {
				event = alog.LOGGER(logChannel).Err(err)
			}

			event = event.
				//Str("time_rfc3339", time.Now().Format(time.RFC3339)).
				Str("id", req.Header.Get(echo.HeaderXRequestID)).
				Str("remote_ip", realIP).
				Str("host", req.Host).
				Str("method", req.Method).
				Str("uri", req.RequestURI).
				Str("user_agent", req.UserAgent()).
				Int("status", res.Status).
				Dur("latency", latency).
				Str("latency_human", latency.String()).
				Int64("bytes_out", res.Size)

			// Handle bytes_in
			cl := req.Header.Get(echo.HeaderContentLength)
			if cl == "" {
				cl = "0"
			}
			if contentLength, parseErr := strconv.Atoi(cl); parseErr == nil {
				event.Int("bytes_in", contentLength)
			}

			// Dynamic keyword detection
			config.parseCustomTags(event, c, req, res, start, stop)

			// Log errors, if any
			if err != nil {
				event.Err(err)
			}

			event.Msg("request")
			return err
		}
	}
}

// parseCustomTags dynamically handles keywords for additional log fields.
func (config *LoggerConfig) parseCustomTags(event *zerolog.Event, c echo.Context, req *http.Request, res *echo.Response, start, stop time.Time) {
	tags := []string{
		"protocol",
		"referer",
		"path",
		"route",
	}

	for _, tag := range tags {
		switch tag {
		case "protocol":
			event.Str(tag, req.Proto)
		case "referer":
			event.Str(tag, req.Referer())
		case "path":
			p := req.URL.Path
			if p == "" {
				p = "/"
			}
			event.Str(tag, p)
		case "route":
			event.Str(tag, c.Path())
		}
	}

	// Custom tag handler
	if config.CustomTagFunc != nil {
		config.CustomTagFunc(c, event)
	}
}

// getLogChannelForIP checks if the RealIP has a specific log channel; otherwise, it returns the default.
func getLogChannelForIP(realIP string, logSeparator *sync.Map, defaultChannel alog.ChannelLabel) alog.ChannelLabel {
	if logSeparator != nil {
		if value, ok := logSeparator.Load(realIP); ok {
			if logChannel, valid := value.(alog.ChannelLabel); valid {
				return logChannel
			}
		}
	}
	return defaultChannel
}

// AddLogSeparator dynamically updates the LogSeparator.
func AddLogSeparator(config *LoggerConfig, ip string, logChannel alog.ChannelLabel) {
	if config.LogSeparator == nil {
		config.LogSeparator = &sync.Map{}
	}
	config.LogSeparator.Store(ip, logChannel)
}

// RemoveLogSeparator dynamically removes an IP from LogSeparator.
func RemoveLogSeparator(config *LoggerConfig, ip string) {
	if config.LogSeparator != nil {
		config.LogSeparator.Delete(ip)
	}
}

type IPLogChannelMap map[string]alog.ChannelLabel

// InitLogSeparator initializes LogSeparator with a given IPLogChannelMap.
func InitLogSeparator(config *LoggerConfig, ipLogMap IPLogChannelMap) {
	if config.LogSeparator == nil {
		config.LogSeparator = &sync.Map{}
	}

	for ip, logChannel := range ipLogMap {
		config.LogSeparator.Store(ip, logChannel)
	}
}
