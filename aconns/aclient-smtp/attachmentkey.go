package aclient_smtp

import (
	"errors"
	"strings"
)

// AttachmentKey represents an attachment key type
type AttachmentKey string

// Predefined AttachmentKeys
const (
	ATTACHMENTKEY_NONE AttachmentKey = ""
	ATTACHMENTKEY_FILE AttachmentKey = "file"
	ATTACHMENTKEY_ID   AttachmentKey = "id"
)

// IsEmpty checks if the AttachmentKey is empty
func (ak AttachmentKey) IsEmpty() bool {
	return ak == ""
}

// TrimSpace returns the trimmed string representation of the AttachmentKey
func (ak AttachmentKey) TrimSpace() AttachmentKey {
	return AttachmentKey(strings.TrimSpace(string(ak)))
}

// String returns the string representation of the AttachmentKey
func (ak AttachmentKey) String() string {
	return string(ak)
}

// Matches checks if the AttachmentKey matches the given string
func (ak AttachmentKey) Matches(s string) bool {
	return string(ak) == s
}

// ToStringTrimLower returns the trimmed and lowercased string representation of the AttachmentKey
func (ak AttachmentKey) ToStringTrimLower() string {
	return strings.ToLower(strings.TrimSpace(string(ak)))
}

// Validate checks if the AttachmentKey is valid
func (ak AttachmentKey) Validate() error {
	if ak.IsEmpty() {
		return nil
	}

	parts := strings.Split(string(ak), ":")
	if len(parts) != 2 {
		return errors.New("AttachmentKey must contain exactly one ':'")
	}

	for _, part := range parts {
		if strings.TrimSpace(part) == "" {
			return errors.New("AttachmentKey parts must not be empty")
		}
	}

	return nil
}

// IsFile checks if the AttachmentKey has the "file" prefix
func (ak AttachmentKey) IsFile() bool {
	return strings.HasPrefix(string(ak), string(ATTACHMENTKEY_FILE)+":")
}

// IsId checks if the AttachmentKey has the "id" prefix
func (ak AttachmentKey) IsId() bool {
	return strings.HasPrefix(string(ak), string(ATTACHMENTKEY_ID)+":")
}

// GetParts splits the AttachmentKey into its key and target parts
func (ak AttachmentKey) GetParts() (key AttachmentKey, target string, err error) {
	if ak.IsEmpty() {
		return ATTACHMENTKEY_NONE, "", nil
	}
	if err = ak.Validate(); err != nil {
		return "", "", err
	}
	parts := strings.SplitN(string(ak), ":", 2)
	return AttachmentKey(parts[0]), parts[1], nil
}

// AttachmentKeys represents a slice of AttachmentKey
type AttachmentKeys []AttachmentKey

// IsEmpty checks if the AttachmentKeys slice is empty
func (aks AttachmentKeys) IsEmpty() bool {
	return len(aks) == 0
}

// String returns the string representation of the AttachmentKeys slice
func (aks AttachmentKeys) String() string {
	return strings.Join(aks.ToStringArray(), ", ")
}

// ToStringArray returns an array of AttachmentKeys as strings
func (aks AttachmentKeys) ToStringArray() []string {
	strArray := make([]string, len(aks))
	for i, ak := range aks {
		strArray[i] = ak.String()
	}
	return strArray
}

// Find returns the AttachmentKey if found, otherwise an empty AttachmentKey
func (aks AttachmentKeys) Find(ak AttachmentKey) AttachmentKey {
	for _, v := range aks {
		if v == ak {
			return v
		}
	}
	return ""
}

// HasKey checks if the AttachmentKeys slice contains the given AttachmentKey
func (aks AttachmentKeys) HasKey(s AttachmentKey) bool {
	return aks.Find(s) != ""
}

// Matches checks if any AttachmentKey in the AttachmentKeys slice matches the given string
func (aks AttachmentKeys) Matches(s string) bool {
	for _, v := range aks {
		if v.Matches(s) {
			return true
		}
	}
	return false
}
