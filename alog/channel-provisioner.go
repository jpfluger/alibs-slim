package alog

import (
	"errors"
	"fmt"
	"github.com/rs/zerolog"
	"github.com/jpfluger/alibs-slim/autils"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
	"time"
)

// IChannelProvisioner defines the interface for a channel provisioner,
// which provides methods for configuring and managing logging channels.
type IChannelProvisioner interface {
	// GetLogDir returns the directory path where log files will be stored.
	GetLogDir() string

	// GetFileLoggerOptions returns the configuration options for file-based logging,
	// providing settings such as file rotation, max size, and retention policies.
	GetFileLoggerOptions() *FileLoggerOptions

	// AddWith allows additional structured logging fields to be added to the logger.
	// This is intended to enrich the log output with consistent metadata.
	AddWith(logger zerolog.Logger) zerolog.Logger

	// GetWriters
	GetWriters(ch *Channel, prov IChannelProvisioner) ([]io.Writer, error)
}

// ChannelProvisionerBase provides a base implementation of the IChannelProvisioner interface,
// offering core configurations for logging channels such as the log directory and file options.
type ChannelProvisionerBase struct {
	DirLog            string             // Directory path for storing log files.
	FileLoggerOptions *FileLoggerOptions // Options struct containing file logger settings.
}

// GetLogDir returns the directory path where log files will be stored.
func (cp *ChannelProvisionerBase) GetLogDir() string {
	return cp.DirLog
}

// GetFileLoggerOptions returns the configuration options for the file logger,
// enabling structured and file-based log management.
func (cp *ChannelProvisionerBase) GetFileLoggerOptions() *FileLoggerOptions {
	return cp.FileLoggerOptions
}

func (cp *ChannelProvisionerBase) GetWriters(ch *Channel, prov IChannelProvisioner) ([]io.Writer, error) {
	if ch == nil {
		return nil, errors.New("channel is nil")
	}
	if prov == nil {
		return nil, errors.New("prov is nil")
	}
	if len(ch.WriterTypes) == 0 {
		return nil, errors.New("writer types is empty")
	}

	writers := []io.Writer{}

	for _, wt := range ch.WriterTypes {
		switch wt {
		case WRITERTYPE_CONSOLE_STDOUT:
			writers = append(writers, zerolog.ConsoleWriter{
				Out:        os.Stdout,
				TimeFormat: time.RFC3339,
			})
		case WRITERTYPE_CONSOLE_STDERR:
			writers = append(writers, zerolog.ConsoleWriter{
				Out:        os.Stderr,
				TimeFormat: time.RFC3339,
			})
		case WRITERTYPE_STDOUT:
			writers = append(writers, os.Stdout)
		case WRITERTYPE_STDERR:
			writers = append(writers, os.Stderr)
		case WRITERTYPE_FILE:
			// Ensure the log directory is resolved.
			if _, err := autils.ResolveDirectory(prov.GetLogDir()); err != nil {
				return nil, fmt.Errorf("log dir not defined; %v", err)
			}
			// Setup the file logger with the specified options.
			writers = append(writers, &lumberjack.Logger{
				Filename:   fmt.Sprintf("%s/%s.log", prov.GetLogDir(), ch.Name.String()),
				MaxSize:    ch.FileLoggerOptions.MaxSize,
				MaxBackups: ch.FileLoggerOptions.MaxBackups,
				MaxAge:     ch.FileLoggerOptions.MaxAge,
				Compress:   ch.FileLoggerOptions.Compress,
			})
		}
	}

	return writers, nil
}

// ChannelProvisioner extends ChannelProvisionerBase to provide additional
// context-specific logging configurations, such as application and server identifiers.
// Use ChannelProvisioner as an example of how to extend ChannelProvisionerBase in your
// own applications.
type ChannelProvisioner struct {
	ChannelProvisionerBase        // Embeds ChannelProvisionerBase for base functionality.
	App                    string // Identifier for the application, used as a log field.
	Svr                    string // Identifier for the server, used as a log field.
}

// AddWith enriches the provided zerolog.Logger instance by adding structured fields.
// These fields include the application identifier (App) and server identifier (Svr),
// along with a timestamp to track when each log entry was created.
// This method returns the enriched zerolog.Logger instance for further logging use.
func (cp *ChannelProvisioner) AddWith(logger zerolog.Logger) zerolog.Logger {
	return logger.With().
		Timestamp().        // Adds a timestamp to each log entry.
		Str("app", cp.App). // Adds the "app" field with the application identifier.
		Str("svr", cp.Svr). // Adds the "svr" field with the server identifier.
		Logger()            // Returns the updated logger with the added context.
}
