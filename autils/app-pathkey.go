package autils

import "strings"

// AppPathKey constants represent different directory keys for an application.
const (
	DIR_ROOT                AppPathKey = "DIR_ROOT"
	DIR_ETC                 AppPathKey = "DIR_ETC"
	DIR_DATA                AppPathKey = "DIR_DATA"
	DIR_LOGS                AppPathKey = "DIR_LOGS"
	DIR_TMP_CACHE           AppPathKey = "DIR_TMP_CACHE"
	DIR_WEBDISTRO           AppPathKey = "DIR_WEBDISTRO"
	DIR_WEBDISTRO_PUBLIC    AppPathKey = "DIR_WEBDISTRO_PUBLIC"
	DIR_WEBDISTRO_PROTED    AppPathKey = "DIR_WEBDISTRO_PROTED" // short for PROTECTED
	DIR_WEBDISTRO_TEMPLATES AppPathKey = "DIR_WEBDISTRO_TEMPLATES"
)

// AppPathKey is a custom type for application path keys.
type AppPathKey string

// IsEmpty checks if the AppPathKey is empty after trimming whitespace.
func (key AppPathKey) IsEmpty() bool {
	return strings.TrimSpace(string(key)) == ""
}

// TrimSpace trims leading and trailing whitespace from the AppPathKey.
func (key AppPathKey) TrimSpace() AppPathKey {
	return AppPathKey(strings.TrimSpace(string(key)))
}

// String returns the string representation of the AppPathKey.
func (key AppPathKey) String() string {
	return string(key)
}

// ToStringTrimLower trims whitespace and converts the AppPathKey to lowercase.
func (key AppPathKey) ToStringTrimLower() string {
	return strings.ToLower(strings.TrimSpace(string(key)))
}

// Type extracts the type of path key, assuming it starts with 'dir' or 'file'.
func (key AppPathKey) Type() string {
	ss := strings.Split(key.ToStringTrimLower(), "_")
	if len(ss) > 0 {
		switch ss[0] {
		case "dir", "file":
			return ss[0]
		}
	}
	return ""
}

// IsDir checks if the AppPathKey represents a directory.
func (key AppPathKey) IsDir() bool {
	return key.Type() == "dir"
}

// IsFile checks if the AppPathKey represents a file.
func (key AppPathKey) IsFile() bool {
	return key.Type() == "file"
}

type AppPathKeys []AppPathKey
