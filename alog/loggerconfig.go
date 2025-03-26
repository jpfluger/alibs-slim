package alog

// LoggerConfig represents the configuration for loggers,
// encapsulating the Channels in a structure suitable for JSON serialization.
type LoggerConfig struct {
	Channels Channels `json:"channels,omitempty"` // Channels is a slice of Channel configurations.
}

// HasChannels checks if the LoggerConfig has any Channels configured.
func (lc *LoggerConfig) HasChannels() bool {
	return lc != nil && len(lc.Channels) > 0 // Simplified check for non-nil and non-empty Channels slice.
}
