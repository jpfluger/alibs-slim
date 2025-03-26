package aconns

import "strings"

// DbColumnName is a type that represents a database column name.
type DbColumnName string

// IsEmpty checks if the DbColumnName is empty after trimming whitespace.
func (cn DbColumnName) IsEmpty() bool {
	return strings.TrimSpace(string(cn)) == ""
}

// TrimSpace returns a new DbColumnName with leading and trailing whitespace removed.
func (cn DbColumnName) TrimSpace() DbColumnName {
	return DbColumnName(strings.TrimSpace(string(cn)))
}

// String converts the DbColumnName to a regular string.
func (cn DbColumnName) String() string {
	return string(cn)
}

// Bytes converts the DbColumnName to a slice of bytes.
func (cn DbColumnName) Bytes() []byte {
	return []byte(cn)
}
