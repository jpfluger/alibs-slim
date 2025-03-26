package acrypt

import (
	"fmt"
	"strings"
	"time"
)

// IdBrief represents a brief identifier, typically a string.
type IdBrief string

// DefaultTimeFormat is the default format for time-based IdBriefs.
const DefaultTimeFormat = "060102"

// IsEmpty checks if the IdBrief is empty after trimming whitespace.
func (id IdBrief) IsEmpty() bool {
	return strings.TrimSpace(string(id)) == ""
}

// TrimSpace trims whitespace from the IdBrief and returns a new IdBrief.
func (id IdBrief) TrimSpace() IdBrief {
	return IdBrief(strings.TrimSpace(string(id)))
}

// HasMatch checks if the target IdBrief matches the current IdBrief.
func (id IdBrief) HasMatch(target IdBrief) bool {
	return target == id
}

// String converts the IdBrief to a regular string.
func (id IdBrief) String() string {
	return string(id)
}

// ToUpper converts the IdBrief to uppercase.
func (id IdBrief) ToUpper() string {
	return strings.ToUpper(id.String())
}

// MustNewIdBrief is a generic wrapper that panics on error.
func MustNewIdBrief(generatorFunc func() (IdBrief, error)) IdBrief {
	id, err := generatorFunc()
	if err != nil {
		panic(err)
	}
	return id
}

// NewIdBrief4Digits creates a new IdBrief with 4 random digits.
func NewIdBrief4Digits() (IdBrief, error) {
	id, err := RandGenerate4Digits()
	if err != nil {
		return "", err
	}
	return IdBrief(id), nil
}

// NewIdBriefToDay creates a new IdBrief with the current day.
func NewIdBriefToDay() (IdBrief, error) {
	return newIdBriefWithOptions(time.Time{}, DefaultTimeFormat, 6)
}

// NewIdBriefToHour creates a new IdBrief with the current hour.
func NewIdBriefToHour() (IdBrief, error) {
	return newIdBriefWithOptions(time.Time{}, DefaultTimeFormat+"15", 6)
}

// NewIdBriefToMinute creates a new IdBrief with the current minute.
func NewIdBriefToMinute() (IdBrief, error) {
	return newIdBriefWithOptions(time.Time{}, DefaultTimeFormat+"1504", 6)
}

// newIdBriefWithOptions creates a new IdBrief with the given options.
func newIdBriefWithOptions(targetTime time.Time, timeFormat string, length int) (IdBrief, error) {
	codes, err := GenerateMiniRandomCodes(1, length)
	if err != nil {
		return "", fmt.Errorf("failed to generate IdBrief; %v", err)
	}
	if targetTime.IsZero() {
		targetTime = time.Now().UTC()
	}
	if strings.TrimSpace(timeFormat) == "" {
		timeFormat = DefaultTimeFormat
	}
	idBrief := fmt.Sprintf("%s-%s", targetTime.Format(timeFormat), codes[0])
	return IdBrief(strings.ToUpper(idBrief)), nil
}
