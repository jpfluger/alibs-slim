package alog

// FileLoggerOptions defines the configuration options for a file-based logger.
type FileLoggerOptions struct {
	MaxSize    int  `json:"maxSize,omitempty"`    // MaxSize is the maximum size in megabytes of the log file before it gets rotated.
	MaxBackups int  `json:"maxBackups,omitempty"` // MaxBackups is the maximum number of old log files to retain.
	MaxAge     int  `json:"maxAge,omitempty"`     // MaxAge is the maximum number of days to retain old log files.
	Compress   bool `json:"compress,omitempty"`   // Compress determines if the rotated log files should be compressed using gzip.
}
