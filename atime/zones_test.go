package atime

import (
	"runtime"
	"testing"
	"time"
)

func TestTimeIn(t *testing.T) {
	_, err := TimeIn(time.Now(), "America/New_York")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	_, err = TimeIn(time.Now(), "Invalid/Timezone")
	if err == nil {
		t.Errorf("Expected an error for invalid timezone, got nil")
	}
}

func TestTimeInNoError(t *testing.T) {
	_ = TimeInNoError(time.Now(), "America/New_York")
	_ = TimeInNoError(time.Now(), "Invalid/Timezone")
}

func TestTimeInPointer(t *testing.T) {
	_, err := TimeInPointer(time.Now(), "America/New_York")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	_, err = TimeInPointer(time.Now(), "Invalid/Timezone")
	if err == nil {
		t.Errorf("Expected an error for invalid timezone, got nil")
	}
}

func TestTimeInPointerNoError(t *testing.T) {
	_ = TimeInPointerNoError(time.Now(), "America/New_York")
	_ = TimeInPointerNoError(time.Now(), "Invalid/Timezone")
}

func TestGetLocation(t *testing.T) {
	_, err := GetLocation("America/New_York")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	_, err = GetLocation("Invalid/Timezone")
	if err == nil {
		t.Errorf("Expected an error for invalid timezone, got nil")
	}
}

func TestGetCurrentTimeInZone(t *testing.T) {
	_, err := GetCurrentTimeInZone("America/New_York")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	_, err = GetCurrentTimeInZone("Invalid/Timezone")
	if err == nil {
		t.Errorf("Expected an error for invalid timezone, got nil")
	}
}

func TestConvertToTimeZone(t *testing.T) {
	now := time.Now()
	converted := ConvertToTimeZone(now, "America/New_York")
	if converted.Location().String() != "America/New_York" {
		t.Errorf("Expected location America/New_York, got %v", converted.Location())
	}

	converted = ConvertToTimeZone(now, "Invalid/Timezone")
	if converted.Location().String() != "UTC" {
		t.Errorf("Expected location UTC for invalid timezone, got %v", converted.Location())
	}
}

func TestTimeZoneOffset(t *testing.T) {
	offset, err := TimeZoneOffset("America/New_York")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if offset != -5 && offset != -4 { // Considering daylight saving time
		t.Errorf("Expected offset -5 or -4, got %d", offset)
	}

	_, err = TimeZoneOffset("Invalid/Timezone")
	if err == nil {
		t.Errorf("Expected an error for invalid timezone, got nil")
	}
}

func TestGetSystemTimeZone(t *testing.T) {
	tz := GetSystemTimeZone()
	if tz == "" {
		t.Errorf("Expected a timezone, got empty string")
	}
}

func TestGetOSTimeZones(t *testing.T) {
	osHost := runtime.GOOS
	switch osHost {
	case "windows":
		t.Skip("Skipping test on Windows system")
	case "darwin", "linux":
		tzList := GetOSTimeZones()
		if len(tzList) == 0 {
			t.Errorf("Expected a list of timezones, got empty list")
		}
	default:
		t.Skip("Skipping test on unsupported system")
	}
}
