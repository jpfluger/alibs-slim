package alog

import (
	"github.com/rs/zerolog"
)

// LoggerMap is a map that associates ChannelLabel with a pointer to a zerolog.Logger.
type LoggerMap map[ChannelLabel]*zerolog.Logger

// Get retrieves the logger associated with the given ChannelLabel.
// If there is no logger associated with the label, it returns nil.
func (lm LoggerMap) Get(name ChannelLabel) *zerolog.Logger {
	return lm[name]
}
