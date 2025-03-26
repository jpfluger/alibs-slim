package alog

import (
	"github.com/rs/zerolog"
	"io"
)

// MockWriter captures logs into an array of strings for testing
type MockWriter struct {
	Logs []string
}

// Write appends log entries to the internal Logs array
func (mw *MockWriter) Write(p []byte) (n int, err error) {
	mw.Logs = append(mw.Logs, string(p))
	return len(p), nil
}

// Reset clears the captured logs
func (mw *MockWriter) Reset() {
	mw.Logs = []string{}
}

// MockLogChannelProvisioner is a mock logger provisioner for unit testing
type MockLogChannelProvisioner struct {
	ChannelProvisionerBase
	App    string
	Writer *MockWriter
}

// AddWith adds metadata to the logger
func (cp *MockLogChannelProvisioner) AddWith(logger zerolog.Logger) zerolog.Logger {
	return logger.With().
		Timestamp().
		Str("app", cp.App).
		Logger()
}

// GetWriters returns the mock writer as the logger output
func (cp *MockLogChannelProvisioner) GetWriters(ch *Channel, prov IChannelProvisioner) ([]io.Writer, error) {
	if cp.Writer == nil {
		cp.Writer = &MockWriter{}
	}
	return []io.Writer{cp.Writer}, nil
}

// NewMockLogChannelProvisioner creates a new MockLogChannelProvisioner
func NewMockLogChannelProvisioner(app string) *MockLogChannelProvisioner {
	return &MockLogChannelProvisioner{
		ChannelProvisionerBase: ChannelProvisionerBase{
			DirLog:            "",
			FileLoggerOptions: nil,
		},
		App:    app,
		Writer: &MockWriter{},
	}
}

// SetupMockLogger sets up a mock logger for testing
func SetupMockLogger(channelName ChannelLabel, logLevel zerolog.Level) (*MockLogChannelProvisioner, error) {
	channels := Channels{
		&Channel{
			Name:        channelName,
			LogLevel:    logLevel.String(),
			WriterTypes: WriterTypes{"custom"},
		},
	}

	prov := NewMockLogChannelProvisioner("tester")

	// Set global logger
	if err := setGlobalLogger("", channels, prov); err != nil {
		return nil, err
	}

	return prov, nil
}
