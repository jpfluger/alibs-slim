package ahttp

import (
	"github.com/mileusna/useragent"
	"testing"
)

// TestParseUserAgentString tests the ParseUserAgentString function.
func TestParseUserAgentString(t *testing.T) {
	uaString := "Mozilla/5.0 (iPhone; CPU iPhone OS 10_3 like Mac OS X) AppleWebKit/602.1.50 (KHTML, like Gecko) CriOS/56.0.2924.75 Mobile/14E5239e Safari/602.1"
	ua, deviceName := ParseUserAgentString(uaString)

	if ua.Name != "Chrome" || !ua.Mobile {
		t.Errorf("ParseUserAgentString() failed to parse user agent string correctly")
	}

	expectedDeviceName := "Chrome v56.0.2924.75; iOS v10.3"
	if deviceName != expectedDeviceName {
		t.Errorf("ParseUserAgentString() deviceName = %v, want %v", deviceName, expectedDeviceName)
	}
}

// TestIsMobileDevice tests the IsMobileDevice function.
func TestIsMobileDevice(t *testing.T) {
	ua := useragent.UserAgent{Mobile: true}
	if !IsMobileDevice(&ua) {
		t.Errorf("IsMobileDevice() should return true for mobile user agents")
	}
}

// TestIsSpecificBrowser tests the IsSpecificBrowser function.
func TestIsSpecificBrowser(t *testing.T) {
	ua := useragent.UserAgent{Name: "Firefox"}
	if !IsSpecificBrowser(&ua, "Firefox") {
		t.Errorf("IsSpecificBrowser() should return true for Firefox user agents")
	}
}

// TestIsBot tests the IsBot function.
func TestIsBot(t *testing.T) {
	ua := useragent.UserAgent{Bot: true}
	if !IsBot(&ua) {
		t.Errorf("IsBot() should return true for bot user agents")
	}
}

// TestPrettyPrint tests the PrettyPrint function.
func TestPrettyPrint(t *testing.T) {
	ua := useragent.UserAgent{
		Name:      "Chrome",
		Version:   "56.0.2924.75",
		OS:        "iOS",
		OSVersion: "10.3",
		Device:    "iPhone",
		Mobile:    true,
	}
	expected := "Browser: Chrome v56.0.2924.75\nOS: iOS v10.3\nDevice: iPhone (Mobile)"
	result := PrettyPrint(&ua)
	if result != expected {
		t.Errorf("PrettyPrint() = %v, want %v", result, expected)
	}
}
