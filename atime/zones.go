package atime

import (
	"github.com/mileusna/timezones"
	"os"
	"runtime"
	"time"
)

// TimeIn returns the time in a specified timezone.
// If the name is empty or "UTC", it returns the time in UTC.
// If the name is "Local", it returns the local time.
// Otherwise, it assumes the name is a location name in the IANA Time Zone database.
func TimeIn(t time.Time, timeZoneId string) (time.Time, error) {
	if t.IsZero() {
		t = time.Now().UTC()
	}
	loc, err := time.LoadLocation(timeZoneId)
	if err != nil {
		return time.Time{}, err
	}
	return t.In(loc), nil
}

// TimeInNoError is similar to TimeIn but does not return an error.
func TimeInNoError(t time.Time, timeZoneId string) time.Time {
	tt, _ := TimeIn(t, timeZoneId)
	return tt
}

// TimeInPointer returns a pointer to the time in a specified timezone.
func TimeInPointer(t time.Time, timeZoneId string) (*time.Time, error) {
	tt, err := TimeIn(t, timeZoneId)
	return &tt, err
}

// TimeInPointerNoError is similar to TimeInPointer but does not return an error.
func TimeInPointerNoError(t time.Time, timeZoneId string) *time.Time {
	tt := TimeInNoError(t, timeZoneId)
	return &tt
}

// GetLocation returns the time.Location for a given timezone ID.
func GetLocation(timeZoneID string) (*time.Location, error) {
	return time.LoadLocation(timeZoneID)
}

// GetCurrentTimeInZone returns the current time in the specified timezone.
func GetCurrentTimeInZone(timeZoneID string) (time.Time, error) {
	loc, err := GetLocation(timeZoneID)
	if err != nil {
		return time.Time{}, err
	}
	return time.Now().In(loc), nil
}

// ConvertToTimeZone converts the given time to the specified time zone.
func ConvertToTimeZone(t interface{}, timeZoneId string) time.Time {
	dt := EnsureDateTime(t)
	if dt.IsZero() {
		return dt
	}
	loc, err := time.LoadLocation(timeZoneId)
	if err != nil {
		return dt.In(time.UTC)
	}
	return dt.In(loc)
}

// TimeZoneOffset returns the offset in hours for the specified timezone.
func TimeZoneOffset(timeZoneID string) (int, error) {
	loc, err := GetLocation(timeZoneID)
	if err != nil {
		return 0, err
	}

	_, offset := time.Now().In(loc).Zone()
	return offset / 3600, nil // Convert offset from seconds to hours
}

// GetSystemTimeZone attempts to determine the system's timezone.
func GetSystemTimeZone() string {
	osHost := runtime.GOOS
	switch osHost {
	case "windows":
		// Windows timezone names do not align with IANA timezone names.
		return "UTC"
	case "darwin", "linux":
		// Darwin (macOS) and Linux timezones can be determined from the TZ environment variable.
		if zone, ok := os.LookupEnv("TZ"); ok {
			return zone
		}
		// Fallback to local time zone if TZ is not set.
		loc, err := time.LoadLocation("Local")
		if err == nil {
			return loc.String()
		}
	}
	return "UTC"
}

// GetOSTimeZones retrieves a list of valid timezones for the operating system.
func GetOSTimeZones() []string {
	osHost := runtime.GOOS
	switch osHost {
	case "windows":
		// Windows timezone names do not align with IANA timezone names.
		return nil
	case "darwin", "linux":
		// Darwin (macOS) and Linux timezones can be determined using timezones package.
		return timezones.List()
	}
	return nil
}
