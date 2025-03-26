package alog

import (
	"github.com/rs/zerolog"
	"os"
	"testing"
)

//type ChannelProvisioner struct {
//	ChannelProvisionerBase
//}
//
//func (cp *ChannelProvisioner) AddWith(logger zerolog.Logger) zerolog.Logger {
//	return logger.With().
//		Timestamp().
//		Str("app", "app-name").
//		Str("svr", "server-id").
//		Logger()
//}
//
//func TestLOGGER(t *testing.T) {
//
//	channels := Channels{
//		&Channel{
//			Name:        LOGGER_APP,
//			LogLevel:    zerolog.LevelErrorValue, // "" (empty) is NoLevel, "error" is correct whereas "err" is not found in zerolog.
//			WriterTypes: WriterTypes{WRITERTYPE_CONSOLE_STDOUT, WRITERTYPE_STDOUT, WRITERTYPE_CONSOLE_STDERR},
//		},
//	}
//	prov := &ChannelProvisioner{}
//
//	if err := SetGlobalLogger("", channels, prov); err != nil {
//		t.Error(err)
//		return
//	}
//
//	LOGGER(LOGGER_APP).Debug().Msg("test")
//	LOGGER(LOGGER_APP).Info().Msg("test")
//	LOGGER(LOGGER_APP).Warn().Msg("test")
//	LOGGER(LOGGER_APP).Err(fmt.Errorf("this is an error")).Msg("test")
//
//	LOGGER("unknown").Err(fmt.Errorf("unknown logger tripped")).Msg("unknown")
//}

// mockChannelProvisioner implements the IChannelProvisioner interface for testing purposes.
type mockChannelProvisioner struct {
	ChannelProvisionerBase
}

func (m *mockChannelProvisioner) GetFileLoggerOptions() *FileLoggerOptions {
	return &FileLoggerOptions{
		MaxSize:    1,
		MaxBackups: 1,
		MaxAge:     1,
		Compress:   false,
	}
}

func (m *mockChannelProvisioner) GetLogDir() string {
	return os.TempDir()
}

func (m *mockChannelProvisioner) AddWith(logger zerolog.Logger) zerolog.Logger {
	return logger.With().
		Timestamp().
		Str("app", "app-name").
		Str("svr", "server-id").
		Logger()
}

// TestGetGlobalLoggerConfig tests the retrieval of the global logger configuration.
func TestGetGlobalLoggerConfig(t *testing.T) {
	// Setup
	channels := Channels{
		&Channel{Name: LOGGER_APP, LogLevel: "info", WriterTypes: WriterTypes{WRITERTYPE_FILE}},
	}
	prov := &mockChannelProvisioner{}
	if err := SetGlobalLogger("", channels, prov); err != nil {
		t.Error(err)
		return
	}

	// Test
	config := GetGlobalLoggerConfig()
	if config == nil {
		t.Error("Expected non-nil config")
	}
	if len(config.Channels) != len(channels) {
		t.Errorf("Expected %d channels, got %d", len(channels), len(config.Channels))
	}
}

// TestLOGGER tests the LOGGER function for retrieving loggers.
func TestLOGGER(t *testing.T) {
	// Setup
	channels := Channels{
		&Channel{Name: LOGGER_APP, LogLevel: "info", WriterTypes: WriterTypes{WRITERTYPE_FILE}},
	}
	prov := &mockChannelProvisioner{}
	if err := SetGlobalLogger("", channels, prov); err != nil {
		t.Error(err)
		return
	}

	// Test
	logger := LOGGER(LOGGER_APP)
	if logger == nil {
		t.Error("Expected non-nil logger")
	}
}

// TestSetGlobalLogger tests the SetGlobalLogger function for initializing the global logger map.
func TestSetGlobalLogger(t *testing.T) {
	// Setup
	channels := Channels{
		&Channel{Name: LOGGER_APP, LogLevel: "info", WriterTypes: WriterTypes{WRITERTYPE_FILE}},
	}
	prov := &mockChannelProvisioner{}

	// Test
	err := SetGlobalLogger("", channels, prov)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify
	if globalLM == nil {
		t.Error("Expected globalLM to be non-nil")
	}
	if len(globalLM.Map) != len(channels) {
		t.Errorf("Expected %d loggers, got %d", len(channels), len(globalLM.Map))
	}
}

// TestGlobalLoggerMap_Get tests the Get method of globalLoggerMap.
func TestGlobalLoggerMap_Get(t *testing.T) {
	// Setup
	channels := Channels{
		&Channel{Name: LOGGER_APP, LogLevel: "info", WriterTypes: WriterTypes{WRITERTYPE_FILE}},
	}
	prov := &mockChannelProvisioner{}
	if err := SetGlobalLogger("", channels, prov); err != nil {
		t.Error(err)
		return
	}

	// Test
	logger := globalLM.Get(LOGGER_APP)
	if logger == nil {
		t.Error("Expected non-nil logger")
	}

	// Test unknown logger
	unknownLogger := globalLM.Get("unknown")
	if unknownLogger != globalLM.unknownLogger {
		t.Error("Expected the unknown logger")
	}
}
