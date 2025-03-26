package alog

import (
	"encoding/json"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMockLogChannelProvisioner(t *testing.T) {
	globalLM = nil
	// Define a test channel
	testChannel := ChannelLabel("testchannel")

	// Set up the mock logger
	prov, err := SetupMockLogger(testChannel, zerolog.InfoLevel)
	assert.NoError(t, err)
	assert.NotNil(t, prov)
	assert.NotNil(t, prov.Writer)

	// Log an example message
	logger := LOGGER(testChannel)
	logger.Info().
		Str("key1", "value1").
		Int("key2", 123).
		Msg("Test message")

	// Validate that the log was captured
	assert.Greater(t, len(prov.Writer.Logs), 0, "Expected at least one log entry")

	// Parse the captured log
	var logOutput map[string]interface{}
	err = json.Unmarshal([]byte(prov.Writer.Logs[0]), &logOutput)
	assert.NoError(t, err, "Failed to parse the captured log")

	// Validate log fields
	assert.Equal(t, "info", logOutput["level"], "Expected log level to be 'info'")
	assert.Equal(t, "Test message", logOutput["message"], "Expected log message to match")
	assert.Equal(t, "value1", logOutput["key1"], "Expected 'key1' to be 'value1'")
	assert.EqualValues(t, 123, logOutput["key2"], "Expected 'key2' to be 123")
	assert.NotNil(t, logOutput["time"], "Expected a time field in the log")
}
