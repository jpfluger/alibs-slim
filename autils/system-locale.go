package autils

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// GetSystemLanguageType retrieves the system language and locale in the format "language-locale".
// If only the language is available, it returns "language". If neither is available, it returns an empty string.
func GetSystemLanguageType() LanguageType {
	lang, locale := GetSystemLanguage()
	if lang == "" {
		return ""
	}
	if locale == "" {
		return LanguageType(lang)
	}
	return LanguageType(fmt.Sprintf("%s-%s", lang, locale))
}

// GetSystemLanguage detects the system's language and locale based on the operating system.
// It returns the default language "en" and locale "US" if it cannot detect the system settings.
func GetSystemLanguage() (lang string, locale string) {
	osHost := runtime.GOOS
	defaultLang := "en"
	defaultLoc := "US"

	switch osHost {
	case "windows":
		// Execute PowerShell command to get the system culture on Windows.
		cmd := exec.Command("powershell", "Get-Culture | select -exp Name")
		output, err := cmd.Output()
		if err == nil {
			langLocRaw := strings.TrimSpace(string(output))
			langLoc := strings.Split(langLocRaw, "-")
			if len(langLoc) == 2 {
				return langLoc[0], langLoc[1]
			}
		}
	case "darwin":
		// Execute osascript to get the user locale of the system on macOS.
		cmd := exec.Command("sh", "-c", "osascript -e 'user locale of (get system info)'")
		output, err := cmd.Output()
		if err == nil {
			langLocRaw := strings.TrimSpace(string(output))
			langLoc := strings.Split(langLocRaw, "_")
			if len(langLoc) == 2 {
				return langLoc[0], langLoc[1]
			}
		}
	case "linux":
		// Look up the LANG environment variable on Linux.
		envlang, ok := os.LookupEnv("LANG")
		if ok {
			langLocRaw := strings.Split(envlang, ".")[0]
			langLoc := strings.Split(langLocRaw, "_")
			if len(langLoc) == 2 {
				return langLoc[0], langLoc[1]
			}
		}
	}

	// Return default language and locale if detection fails.
	return defaultLang, defaultLoc
}
