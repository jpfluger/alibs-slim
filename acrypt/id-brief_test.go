package acrypt

import (
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
	"time"
)

// TestIsEmpty checks the IsEmpty method.
func TestIsEmpty(t *testing.T) {
	emptyID := IdBrief(" ")
	if !emptyID.IsEmpty() {
		t.Error("IsEmpty should return true for whitespace-only IdBrief")
	}

	nonEmptyID := IdBrief("id")
	if nonEmptyID.IsEmpty() {
		t.Error("IsEmpty should return false for non-empty IdBrief")
	}
}

// TestTrimSpace checks the TrimSpace method.
func TestTrimSpace(t *testing.T) {
	idWithSpace := IdBrief(" id ")
	if idWithSpace.TrimSpace() != "id" {
		t.Error("TrimSpace should remove leading and trailing whitespace")
	}
}

// TestHasMatch checks the HasMatch method.
func TestHasMatch(t *testing.T) {
	id := IdBrief("id")
	matchingID := IdBrief("id")
	nonMatchingID := IdBrief("id2")

	if !id.HasMatch(matchingID) {
		t.Error("HasMatch should return true for matching IdBriefs")
	}

	if id.HasMatch(nonMatchingID) {
		t.Error("HasMatch should return false for non-matching IdBriefs")
	}
}

// TestString checks the String method.
func TestString(t *testing.T) {
	id := IdBrief("id")
	if id.String() != "id" {
		t.Error("String should return the underlying string of IdBrief")
	}
}

// TestToUpper checks the ToUpper method.
func TestToUpper(t *testing.T) {
	id := IdBrief("id")
	if id.ToUpper() != "ID" {
		t.Error("ToUpper should return an uppercase version of IdBrief")
	}
}

// TestNewIdBrief4Digits checks the NewIdBrief4Digits function.
func TestNewIdBrief4Digits(t *testing.T) {
	id, err := NewIdBrief4Digits()
	if err != nil {
		t.Errorf("NewIdBrief4Digits should not return an error: %v", err)
	}
	if len(id) != 4 {
		t.Errorf("NewIdBrief4Digits should return 4 digits, got: %s", id)
	}
	if !regexp.MustCompile(`^\d{4}$`).MatchString(string(id)) {
		t.Errorf("NewIdBrief4Digits should return only digits, got: %s", id)
	}
}

// TestNewIdBriefToDay checks the NewIdBriefToDay function.
func TestNewIdBriefToDay(t *testing.T) {
	id, err := NewIdBriefToDay()
	if err != nil {
		t.Errorf("NewIdBriefToDay should not return an error: %v", err)
	}
	if !regexp.MustCompile(`^\d{6}-[A-Z0-9]{6}$`).MatchString(string(id)) {
		t.Errorf("NewIdBriefToDay should match the pattern 'YYMMDD-RANDOM', got: %s", id)
	}
}

// TestNewIdBriefToHour checks the NewIdBriefToHour function.
func TestNewIdBriefToHour(t *testing.T) {
	id, err := NewIdBriefToHour()
	if err != nil {
		t.Errorf("NewIdBriefToHour should not return an error: %v", err)
	}
	if !regexp.MustCompile(`^\d{8}-[A-Z0-9]{6}$`).MatchString(string(id)) {
		t.Errorf("NewIdBriefToHour should match the pattern 'YYMMDDHH-RANDOM', got: %s", id)
	}
}

// TestNewIdBriefToMinute checks the NewIdBriefToMinute function.
func TestNewIdBriefToMinute(t *testing.T) {
	id, err := NewIdBriefToMinute()
	if err != nil {
		t.Errorf("NewIdBriefToMinute should not return an error: %v", err)
	}
	if !regexp.MustCompile(`^\d{10}-[A-Z0-9]{6}$`).MatchString(string(id)) {
		t.Errorf("NewIdBriefToMinute should match the pattern 'YYMMDDHHMM-RANDOM', got: %s", id)
	}
}

// TestNewIdBriefWithOptions checks the NewIdBriefWithOptions function.
func TestNewIdBriefWithOptions(t *testing.T) {
	timeFormat := "060102"
	length := 6
	id, err := newIdBriefWithOptions(time.Time{}, timeFormat, length)
	if err != nil {
		t.Errorf("newIdBriefWithOptions should not return an error: %v", err)
	}
	if !regexp.MustCompile(`^\d{6}-[A-Z0-9]{6}$`).MatchString(string(id)) {
		t.Errorf("newIdBriefWithOptions should match the pattern 'YYMMDD-RANDOM', got: %s", id)
	}
}

func TestIdBrief_GenerateOne(t *testing.T) {
	// 4Digits
	idBrief, err := NewIdBrief4Digits()
	assert.NoError(t, err)
	assert.False(t, idBrief.IsEmpty())
	assert.Len(t, idBrief.String(), 4)
	assert.Regexp(t, regexp.MustCompile(`^\d{4}$`), idBrief.String())

	// DAY
	idBrief, err = NewIdBriefToDay()
	assert.NoError(t, err)
	assert.False(t, idBrief.IsEmpty())
	assert.Len(t, idBrief.String(), 13)
	assert.Regexp(t, regexp.MustCompile(`^\d{6}-[A-Z0-9]{6}$`), idBrief.String())

	// HOUR
	idBrief, err = NewIdBriefToHour()
	assert.NoError(t, err)
	assert.False(t, idBrief.IsEmpty())
	assert.Len(t, idBrief.String(), 15)
	assert.Regexp(t, regexp.MustCompile(`^\d{8}-[A-Z0-9]{6}$`), idBrief.String())

	// MINUTE
	idBrief, err = NewIdBriefToMinute()
	assert.NoError(t, err)
	assert.False(t, idBrief.IsEmpty())
	assert.Len(t, idBrief.String(), 17)
	assert.Regexp(t, regexp.MustCompile(`^\d{10}-[A-Z0-9]{6}$`), idBrief.String())

	// CUSTOM: ToMonth w/ 2 random = 0601-###
	idBrief, err = newIdBriefWithOptions(time.Time{}, "0601", 3)
	assert.NoError(t, err)
	assert.False(t, idBrief.IsEmpty())
	assert.Len(t, idBrief.String(), 8)
	assert.Regexp(t, regexp.MustCompile(`^\d{4}-[A-Z0-9]{3}$`), idBrief.String())
}
