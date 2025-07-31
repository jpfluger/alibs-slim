package rruleplus

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/rickar/cal/v2"
	cal_us "github.com/rickar/cal/v2/us"
)

// ICalendar defines the expected calendar interface for holiday support.
type ICalendar interface {
	AddHoliday(holiday ...*cal.Holiday)
	IsHoliday(date time.Time) (actual, observed bool, h *cal.Holiday)
}

// calendarRegistry stores ISO-to-calendar mappings with thread safety.
var (
	calendarRegistry = make(map[string]ICalendar)
	registryMutex    sync.RWMutex
)

// NewCalendar initializes a calendar by ISO code.
// It currently supports "us" (United States). Returns an error for unsupported codes.
func NewCalendar(iso string) (ICalendar, error) {
	iso = CleanISO(iso)

	if iso == "" {
		return nil, fmt.Errorf("invalid or empty ISO code")
	}

	bc := cal.NewBusinessCalendar()

	switch iso {
	case "us":
		bc.AddHoliday(cal_us.Holidays...)
	default:
		return nil, fmt.Errorf("iso code not supported: %s", iso)
	}

	return bc, nil
}

// GetCalendar retrieves a calendar from the global registry by ISO code.
// Returns nil if no calendar is registered under the given code.
func GetCalendar(iso string) (ICalendar, error) {
	iso = CleanISO(iso)

	registryMutex.RLock()
	defer registryMutex.RUnlock()

	calendar, exists := calendarRegistry[iso]
	if !exists {
		return nil, fmt.Errorf("calendar not found for ISO code: %s", iso)
	}

	return calendar, nil
}

// SetCalendar registers a calendar under a normalized ISO code.
func SetCalendar(iso string, c ICalendar) {
	iso = CleanISO(iso)

	registryMutex.Lock()
	defer registryMutex.Unlock()

	calendarRegistry[iso] = c
}

// CleanISO normalizes ISO codes to lowercase and trims whitespace.
func CleanISO(code string) string {
	return strings.TrimSpace(strings.ToLower(code))
}
