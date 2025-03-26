package aconns

import (
	"strings"
)

// AdapterName represents an adapter name
type AdapterName string

// IsEmpty checks if the AdapterName is empty
func (an AdapterName) IsEmpty() bool {
	return an == ""
}

// TrimSpace returns the trimmed string representation of the AdapterName
func (an AdapterName) TrimSpace() AdapterName {
	return AdapterName(strings.TrimSpace(string(an)))
}

// String returns the string representation of the AdapterName
func (an AdapterName) String() string {
	return string(an)
}

// Matches checks if the AdapterName matches the given string
func (an AdapterName) Matches(s string) bool {
	return string(an) == s
}

// AdapterNames represents a slice of AdapterName
type AdapterNames []AdapterName

// IsEmpty checks if the AdapterNames slice is empty
func (ans AdapterNames) IsEmpty() bool {
	return len(ans) == 0
}

// String returns the string representation of the AdapterNames slice
func (ans AdapterNames) String() string {
	return strings.Join(ans.ToStringArray(), ", ")
}

// ToStringArray returns an array of AdapterNames as strings
func (ans AdapterNames) ToStringArray() []string {
	strArray := make([]string, len(ans))
	for i, an := range ans {
		strArray[i] = an.String()
	}
	return strArray
}

// Find returns the AdapterName if found, otherwise an empty AdapterName
func (ans AdapterNames) Find(an AdapterName) AdapterName {
	for _, v := range ans {
		if v == an {
			return v
		}
	}
	return ""
}

// HasKey checks if the AdapterNames slice contains the given AdapterName
func (ans AdapterNames) HasKey(s AdapterName) bool {
	return ans.Find(s) != ""
}

// Matches checks if any AdapterName in the AdapterNames slice matches the given string
func (ans AdapterNames) Matches(s string) bool {
	for _, v := range ans {
		if v.Matches(s) {
			return true
		}
	}
	return false
}
