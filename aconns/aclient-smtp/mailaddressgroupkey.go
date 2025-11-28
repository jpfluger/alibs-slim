package aclient_smtp

import "strings"

// MailAddressGroupKey is a string type representing keys for MailAddressGroupMap.
// It provides methods for common string operations like checking emptiness and trimming.
type MailAddressGroupKey string

const (
	MAG_KEY_SYSTEM MailAddressGroupKey = "system"
	MAG_KEY_USERS  MailAddressGroupKey = "users"
)

// IsEmpty checks if the MailAddressGroupKey is empty after trimming whitespace.
func (mkey MailAddressGroupKey) IsEmpty() bool {
	return strings.TrimSpace(string(mkey)) == ""
}

// TrimSpace returns a new MailAddressGroupKey with leading and trailing whitespace removed.
func (mkey MailAddressGroupKey) TrimSpace() MailAddressGroupKey {
	return MailAddressGroupKey(strings.TrimSpace(string(mkey)))
}

// String returns the string representation of the MailAddressGroupKey.
func (mkey MailAddressGroupKey) String() string {
	return string(mkey)
}
