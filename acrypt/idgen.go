package acrypt

import (
	"crypto/rand"
	"fmt"
	"strings"
	"time"
)

// GenerateIDBase36Caps generates a random Base36 string of the specified length
// that are all capitals and no "0" and "O" characters.
func GenerateIDBase36Caps(length int) string {
	//const base36Chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	// remove "0" and "O". Less confusing
	const base36Chars = "123456789ABCDEFGHIJKLMNPQRSTUVWXYZ"
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		panic("failed to generate random bytes")
	}

	var result strings.Builder
	for _, b := range bytes {
		result.WriteByte(base36Chars[int(b)%len(base36Chars)])
	}
	return result.String()
}

// NewIdGenReadableWithOptions generates human-readable numbers, such
// as for invoice, purchase orders, quotes or other identifiers.
// If running 100,000 of these within a second, there can be a few
// collisions. When using in your own app, double-check uniqueness
// prior to saving, such as to a database.
func NewIdGenReadableWithOptions(format string, prefix string, date time.Time, length int) string {
	dateString := date.Format("20060102")
	randomPart := GenerateIDBase36Caps(length)
	return fmt.Sprintf(format, prefix, dateString, randomPart)
}

func NewIdGenReadableShort(prefix string, date time.Time) string {
	return fmt.Sprintf("%s%s-%s", prefix, date.Format("20060102"), GenerateIDBase36Caps(7))
}

func NewIdGenReadableLong(prefix string, date time.Time) string {
	return fmt.Sprintf("%s-%s-%s", prefix, date.Format("20060102"), GenerateIDBase36Caps(7))
}
