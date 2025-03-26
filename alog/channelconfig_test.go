package alog

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestApplyOverrides_NoOverrides(t *testing.T) {
	channels := Channels{
		&Channel{
			Name:        LOGGER_APP,
			LogLevel:    "error",
			WriterTypes: WriterTypes{"console-stderr", "file"},
		},
		&Channel{
			Name:        LOGGER_AUTH,
			LogLevel:    "warn",
			WriterTypes: WriterTypes{"file"},
		},
		&Channel{
			Name:        LOGGER_HTTP,
			LogLevel:    "info",
			WriterTypes: WriterTypes{"file"},
		},
		&Channel{
			Name:        "LOGGER_PROXY",
			LogLevel:    "info",
			WriterTypes: WriterTypes{"console-stderr", "file"},
		},
	}

	overrides := LogChannelConfigMap{}

	modified, isChanged, err := channels.ApplyOverrides(overrides)
	assert.NoError(t, err, "Unexpected error occurred")
	assert.False(t, isChanged, "isChanged should be false when no overrides are applied")
	assert.Equal(t, channels, modified, "Channels should remain unchanged when no overrides are applied")
}

func TestApplyOverrides_WithOverrides(t *testing.T) {
	channels := Channels{
		&Channel{
			Name:        LOGGER_APP,
			LogLevel:    "error",
			WriterTypes: WriterTypes{"console-stderr", "file"},
		},
		&Channel{
			Name:        LOGGER_AUTH,
			LogLevel:    "warn",
			WriterTypes: WriterTypes{"file"},
		},
	}

	overrides := LogChannelConfigMap{
		LOGGER_APP: {
			LogLevel:    "info",
			WriterTypes: WriterTypes{"file"},
		},
		LOGGER_AUTH: {
			LogLevel: "error",
		},
	}

	expected := Channels{
		&Channel{
			Name:        LOGGER_APP,
			LogLevel:    "info",
			WriterTypes: WriterTypes{"file"},
		},
		&Channel{
			Name:        LOGGER_AUTH,
			LogLevel:    "error",
			WriterTypes: WriterTypes{"file"}, // WriterTypes not overridden
		},
	}

	modified, isChanged, err := channels.ApplyOverrides(overrides)
	assert.NoError(t, err, "Unexpected error occurred")
	assert.True(t, isChanged, "isChanged should be true when overrides are applied")
	assert.Equal(t, expected, modified, "Channels should be updated with overrides")
}

func TestApplyOverrides_PartialOverrides(t *testing.T) {
	channels := Channels{
		&Channel{
			Name:        LOGGER_HTTP,
			LogLevel:    "info",
			WriterTypes: WriterTypes{"file"},
		},
		&Channel{
			Name:        "LOGGER_PROXY",
			LogLevel:    "info",
			WriterTypes: WriterTypes{"console-stderr", "file"},
		},
	}

	overrides := LogChannelConfigMap{
		LOGGER_HTTP: {
			LogLevel: "debug",
		},
	}

	expected := Channels{
		&Channel{
			Name:        LOGGER_HTTP,
			LogLevel:    "debug",             // LogLevel overridden
			WriterTypes: WriterTypes{"file"}, // WriterTypes not overridden
		},
		&Channel{
			Name:        "LOGGER_PROXY",
			LogLevel:    "info",                                // No override
			WriterTypes: WriterTypes{"console-stderr", "file"}, // No override
		},
	}

	modified, isChanged, err := channels.ApplyOverrides(overrides)
	assert.NoError(t, err, "Unexpected error occurred")
	assert.True(t, isChanged, "isChanged should be true when partial overrides are applied")
	assert.Equal(t, expected, modified, "Channels should reflect partial overrides")
}

func TestApplyOverrides_NilChannel(t *testing.T) {
	channels := Channels{
		nil,
		&Channel{
			Name:        LOGGER_HTTP,
			LogLevel:    "info",
			WriterTypes: WriterTypes{"file"},
		},
	}

	overrides := LogChannelConfigMap{
		LOGGER_HTTP: {
			LogLevel: "debug",
		},
	}

	modified, isChanged, err := channels.ApplyOverrides(overrides)
	assert.Error(t, err, "Expected error due to nil channel")
	assert.Nil(t, modified, "Resulting channels array should be nil when error occurs")
	assert.False(t, isChanged, "isChanged should be false when error occurs")
}

func TestApplyOverrides_EmptyChannels(t *testing.T) {
	channels := Channels{}

	overrides := LogChannelConfigMap{
		LOGGER_APP: {
			LogLevel: "debug",
		},
	}

	modified, isChanged, err := channels.ApplyOverrides(overrides)
	assert.Error(t, err, "Expected error due to empty channels")
	assert.Nil(t, modified, "Resulting channels array should be nil when error occurs")
	assert.False(t, isChanged, "isChanged should be false when error occurs")
}

func TestChannels_ToMap(t *testing.T) {
	channels := Channels{
		&Channel{
			Name:        "LOGGER_APP",
			LogLevel:    "error",
			WriterTypes: WriterTypes{"console-stderr", "file"},
		},
		&Channel{
			Name:        "LOGGER_AUTH",
			LogLevel:    "warn",
			WriterTypes: WriterTypes{"file"},
		},
		&Channel{
			Name:        "LOGGER_HTTP",
			LogLevel:    "info",
			WriterTypes: WriterTypes{"file"},
		},
		nil, // Include a nil channel to ensure it is skipped
	}

	expected := LogChannelConfigMap{
		"LOGGER_APP": {
			LogLevel:    "error",
			WriterTypes: WriterTypes{"console-stderr", "file"},
		},
		"LOGGER_AUTH": {
			LogLevel:    "warn",
			WriterTypes: WriterTypes{"file"},
		},
		"LOGGER_HTTP": {
			LogLevel:    "info",
			WriterTypes: WriterTypes{"file"},
		},
	}

	result := channels.ToMap()

	assert.Equal(t, expected, result, "ToMap should convert Channels to LogChannelConfigMap correctly")
}
