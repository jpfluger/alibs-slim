package alog

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/rs/zerolog"
)

// Predefined channel labels for different loggers.
const (
	LOGGER_APP  ChannelLabel = "app"
	LOGGER_AUTH ChannelLabel = "auth"
	LOGGER_SQL  ChannelLabel = "sql"
	LOGGER_HTTP ChannelLabel = "http"
)

// globalLM holds the global logger map instance.
var globalLM *globalLoggerMap

// once ensures that the global logger map is only initialized once.
var once sync.Once

// globalLoggerMap maintains a map of loggers and associated channels.
type globalLoggerMap struct {
	Map           LoggerMap
	Channels      Channels
	unknownLogger *zerolog.Logger
}

// Get retrieves a logger by its channel label. If not found, returns the unknown logger.
func (glm *globalLoggerMap) Get(name ChannelLabel) *zerolog.Logger {
	if lg, ok := glm.Map[name]; ok {
		return lg
	}
	return glm.unknownLogger
}

// LOGGER provides access to the global logger map.
func LOGGER(name ChannelLabel) *zerolog.Logger {
	once.Do(func() {
		if globalLM != nil {
			return
		}
		// Initialize globalLM and any other necessary components here.
		globalLM = &globalLoggerMap{
			Map: make(LoggerMap),
			Channels: Channels{
				&Channel{Name: LOGGER_APP, LogLevel: "err", WriterTypes: WriterTypes{WRITERTYPE_CONSOLE_STDOUT, WRITERTYPE_CONSOLE_STDERR}},
			},
			// Initialize the unknownLogger with a default zerolog.Logger instance.
			unknownLogger: &zerolog.Logger{},
		}
	})
	return globalLM.Get(name)
}

// GetGlobalLoggerConfig returns the current global logger configuration.
func GetGlobalLoggerConfig() *LoggerConfig {
	if globalLM == nil {
		return nil // Return nil or an appropriate default if globalLM is not set.
	}
	return &LoggerConfig{
		Channels: globalLM.Channels,
	}
}

// SetGlobalLogger initializes the global logger map with the provided channels and provisioner.
func SetGlobalLogger(defaultTimeFormat string, channels Channels, prov IChannelProvisioner) error {
	var err error
	once.Do(func() {
		if err = setGlobalLogger(defaultTimeFormat, channels, prov); err != nil {
			return
		}
	})
	return err
}

// setGlobalLogger initializes the global logger map with the provided channels and provisioner.
func setGlobalLogger(defaultTimeFormat string, channels Channels, prov IChannelProvisioner) (err error) {
	if len(channels) == 0 {
		err = fmt.Errorf("channels are empty")
		return
	}
	if prov == nil {
		err = fmt.Errorf("provisioner is nil")
		return
	}

	// Set the default time format if not provided.
	if defaultTimeFormat == "" {
		defaultTimeFormat = time.RFC3339Nano
	}
	zerolog.TimeFieldFormat = defaultTimeFormat
	zerolog.TimestampFieldName = "time"

	// Initialize each channel.
	for _, ch := range channels {
		if initErr := ch.Initialize(prov); initErr != nil {
			err = fmt.Errorf("failed to initialize log channel '%s': %v", ch.Name.String(), initErr)
			return
		}
	}

	// Create a new logger map and populate it.
	mp := make(LoggerMap)
	for _, ch := range channels {
		mp[ch.Name] = &ch.logger
	}

	// Create an unknown logger for undefined channels.
	ul := prov.AddWith(zerolog.New(os.Stderr).Level(zerolog.ErrorLevel))

	// Assign the newly created map and unknown logger to the global logger map.
	globalLM = &globalLoggerMap{
		Map:           mp,
		Channels:      channels,
		unknownLogger: &ul,
	}

	return err
}
