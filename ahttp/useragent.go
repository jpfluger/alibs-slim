package ahttp

import (
	"github.com/mileusna/useragent"
	"strings"
)

// ParseUserAgentString parses the user agent string and returns the useragent.UserAgent object and a descriptive device name.
func ParseUserAgentString(uaString string) (ua *useragent.UserAgent, deviceName string) {
	// Parse the user agent string using the useragent package.
	myUA := useragent.Parse(uaString)

	// Build a descriptive device name based on the parsed user agent details.
	var sbDeviceName strings.Builder
	if myUA.Name != "" {
		sbDeviceName.WriteString(myUA.Name)
		if myUA.Version != "" {
			sbDeviceName.WriteString(" v")
			sbDeviceName.WriteString(myUA.Version)
		}
	}
	if myUA.OS != "" {
		if sbDeviceName.Len() > 0 {
			sbDeviceName.WriteString("; ")
		}
		sbDeviceName.WriteString(myUA.OS)
		if myUA.OSVersion != "" {
			sbDeviceName.WriteString(" v")
			sbDeviceName.WriteString(myUA.OSVersion)
		}
	}

	return &myUA, sbDeviceName.String()
}

// IsMobileDevice checks if the user agent represents a mobile device.
func IsMobileDevice(ua *useragent.UserAgent) bool {
	return ua.Mobile
}

// IsSpecificBrowser checks if the user agent represents a specific browser.
func IsSpecificBrowser(ua *useragent.UserAgent, browserName string) bool {
	return strings.EqualFold(ua.Name, browserName)
}

// IsBot checks if the user agent represents a search engine crawler or bot.
func IsBot(ua *useragent.UserAgent) bool {
	return ua.Bot
}

// PrettyPrint returns a formatted string with the user agent details.
func PrettyPrint(ua *useragent.UserAgent) string {
	var sb strings.Builder
	sb.WriteString("Browser: ")
	sb.WriteString(ua.Name)
	sb.WriteString(" v")
	sb.WriteString(ua.Version)
	sb.WriteString("\nOS: ")
	sb.WriteString(ua.OS)
	sb.WriteString(" v")
	sb.WriteString(ua.OSVersion)
	sb.WriteString("\nDevice: ")
	sb.WriteString(ua.Device)
	if ua.Mobile {
		sb.WriteString(" (Mobile)")
	}
	if ua.Bot {
		sb.WriteString(" (Bot)")
	}
	return sb.String()
}
