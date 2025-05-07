package alog

import (
	"fmt"
	"github.com/rs/zerolog"
	"io"
	"strings"
)

// Predefined WriterType constants.
const (
	WRITERTYPE_CONSOLE_STDOUT WriterType = "console-stdout"
	WRITERTYPE_CONSOLE_STDERR WriterType = "console-stderr"
	WRITERTYPE_STDOUT         WriterType = "stdout"
	WRITERTYPE_STDERR         WriterType = "stderr"
	WRITERTYPE_FILE           WriterType = "file"
)

// Channel represents a logging channel with specific configurations.
type Channel struct {
	Name              ChannelLabel       `json:"name,omitempty"`
	LogLevel          string             `json:"logLevel,omitempty"`
	WriterTypes       WriterTypes        `json:"writerTypes,omitempty"`
	FileLoggerOptions *FileLoggerOptions `json:"fileLoggerOptions,omitempty"`

	level  zerolog.Level
	logger zerolog.Logger
}

// Channels is a slice of pointers to Channel.
type Channels []*Channel

// Initialize sets up the logging channel with the provided configurations.
func (ch *Channel) Initialize(prov IChannelProvisioner) error {
	if ch == nil {
		return fmt.Errorf("channel is nil")
	}
	if prov == nil {
		return fmt.Errorf("channel provisioner is nil")
	}
	if ch.Name.IsEmpty() {
		return fmt.Errorf("channel name is empty")
	}

	// Parse the log level from the configuration.
	lvl, err := zerolog.ParseLevel(ch.LogLevel)
	if err != nil {
		ch.level = zerolog.ErrorLevel
	} else {
		ch.level = lvl
	}

	// Setup file logger options if the WriterType includes file logging.
	if ch.WriterTypes.HasMatch(WRITERTYPE_FILE) {
		if ch.FileLoggerOptions == nil {
			ch.FileLoggerOptions = prov.GetFileLoggerOptions()
			if ch.FileLoggerOptions == nil {
				ch.FileLoggerOptions = &FileLoggerOptions{
					MaxSize:    25,
					MaxBackups: 10,
					MaxAge:     14,
					Compress:   true,
				}
			}
		}
	}

	writers, err := prov.GetWriters(ch, prov)
	if err != nil {
		return fmt.Errorf("get channel writers failed: %s", err)
	}

	if len(writers) == 0 {
		return fmt.Errorf("no writer types found")
	}

	// Create a new logger with the configured writers and log level.
	ch.logger = prov.AddWith(zerolog.New(io.MultiWriter(writers...)).Level(ch.level))

	return nil
}

func (ch *Channel) Validate() error {
	if ch == nil {
		return fmt.Errorf("channel is nil")
	}
	if ch.Name.IsEmpty() {
		return fmt.Errorf("channel name is empty")
	}
	if strings.TrimSpace(ch.LogLevel) == "" {
		return fmt.Errorf("channel log level is empty")
	}
	if ch.WriterTypes == nil || len(ch.WriterTypes) == 0 {
		return fmt.Errorf("channel writer types is empty")
	}
	return nil
}

// ApplyOverrides applies the overrides from the LogChannelConfigMap to the Channels array.
// Returns a new array of Channels with overrides applied, a boolean indicating if any changes were made, or an error.
func (cns Channels) ApplyOverrides(overrideMap LogChannelConfigMap) (Channels, bool, error) {
	if len(cns) == 0 {
		return nil, false, fmt.Errorf("no channels to apply overrides to")
	}

	// Create a copy of the channels to apply changes without mutating the original
	var result Channels
	isChanged := false

	for _, channel := range cns {
		if channel == nil {
			return nil, false, fmt.Errorf("encountered a nil channel in Channels array")
		}

		// Make a copy of the current channel to modify
		newChannel := *channel

		// Check and apply overrides
		if overrideMap.HasOverride(channel.Name) {
			isChanged = true
			overrideMap.ApplyOverrides(&newChannel)
		}

		if err := newChannel.Validate(); err != nil {
			return nil, false, fmt.Errorf("channel validation failed for name '%s': %s", newChannel.Name.String(), err)
		}

		// Append the modified or unmodified channel to the result
		result = append(result, &newChannel)
	}

	return result, isChanged, nil
}

// ToMap converts a Channels array into a LogChannelConfigMap.
// Each Channel is converted to a corresponding LogChannelConfig entry.
func (cns Channels) ToMap() LogChannelConfigMap {
	result := make(LogChannelConfigMap)

	for _, channel := range cns {
		if channel == nil {
			continue // Skip nil channels
		}

		result[channel.Name] = LogChannelConfig{
			LogLevel:    channel.LogLevel,
			WriterTypes: channel.WriterTypes,
		}
	}

	return result
}
