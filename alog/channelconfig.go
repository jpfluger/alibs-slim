package alog

// LogChannelConfig holds the configuration overrides for a logging channel.
type LogChannelConfig struct {
	LogLevel    string      // Override log level as a string
	WriterTypes WriterTypes // Override writer types
}

// LogChannelConfigMap is a custom type for handling channel configuration overrides.
// Use this in a higher-order application to customize logging levels and writer types.
type LogChannelConfigMap map[ChannelLabel]LogChannelConfig

// HasOverride checks if any override exists for the given channel label.
func (lc LogChannelConfigMap) HasOverride(channel ChannelLabel) bool {
	config, exists := lc[channel]
	if !exists {
		return false
	}
	return config.LogLevel != "" || len(config.WriterTypes) > 0
}

// ApplyOverrides applies the log level and writer types overrides to the given channel, if they exist.
func (lc LogChannelConfigMap) ApplyOverrides(channel *Channel) {
	if config, exists := lc[channel.Name]; exists {
		// Apply log level override if it exists.
		if config.LogLevel != "" {
			channel.LogLevel = config.LogLevel
		}
		// Apply writer types override if it exists.
		if len(config.WriterTypes) > 0 {
			channel.WriterTypes = config.WriterTypes
		}
	}
}

// HasChannel checks if the given channel label exists in the log channel config map.
func (lc LogChannelConfigMap) HasChannel(channel ChannelLabel) bool {
	_, exists := lc[channel]
	return exists
}
