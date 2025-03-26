package auser

import (
	"encoding/json"
	"fmt"
	"github.com/jpfluger/alibs-slim/aemail"
	"github.com/jpfluger/alibs-slim/atags"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestNewRecordUserIdentity(t *testing.T) {
	uid := NewUID()
	ids := atags.TagMapString{
		"email": "user@example.com",
		"altId": "12345",
	}

	rui := NewRecordUserIdentity(uid, ids)

	assert.Equal(t, uid, rui.UID)
	assert.Equal(t, "user@example.com", rui.IDs.Value("email"))
	assert.Equal(t, "12345", rui.IDs.Value("altId"))
}

func TestRecordUserIdentity_IsEmpty(t *testing.T) {
	tests := []struct {
		name string
		rui  RecordUserIdentity
		want bool
	}{
		{"Empty RecordUserIdentity", RecordUserIdentity{}, true},
		{"Non-empty RecordUserIdentity", NewRecordUserIdentityByUID(NewUID()), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.rui.IsEmpty())
		})
	}
}

func TestNewRecordUserIdentityByEmail(t *testing.T) {
	email := aemail.EmailAddress("user@example.com")
	rui := NewRecordUserIdentityByEmail(email)

	assert.Equal(t, "user@example.com", rui.IDs.Value("email"))
	assert.True(t, rui.UID.IsNil())
}

func TestNewRecordUserIdentityByUID(t *testing.T) {
	uid := NewUID()
	rui := NewRecordUserIdentityByUID(uid)

	assert.Equal(t, uid, rui.UID)
	assert.True(t, rui.IDs.IsEmpty())
}

func TestNewRecordUserIdentityById(t *testing.T) {
	label := atags.TagKey("altId")
	value := "12345"
	rui := NewRecordUserIdentityById(label, value)

	assert.Equal(t, value, rui.IDs.Value(label))
	assert.True(t, rui.UID.IsNil())
}

func TestRecordUserIdentity_String(t *testing.T) {
	uid := NewUID()
	ids := atags.TagMapString{
		"email": "user@example.com",
		"altId": "12345",
	}
	rui := NewRecordUserIdentity(uid, ids)

	// Convert to a map for order-agnostic comparison
	expectedMap := map[string]string{
		"uid":   uid.String(),
		"email": "user@example.com",
		"altId": "12345",
	}

	actualMap := make(map[string]string)
	parts := strings.Split(rui.String(), ",")
	for _, part := range parts {
		kv := strings.SplitN(part, "=", 2)
		if len(kv) == 2 {
			actualMap[kv[0]] = kv[1]
		}
	}

	assert.Equal(t, expectedMap, actualMap)
}

func TestRecordUserIdentity_MarshalJSON(t *testing.T) {
	uid := NewUID()
	ids := atags.TagMapString{
		"altId": "12345",
		"email": "user@example.com",
	}
	rui := NewRecordUserIdentity(uid, ids)

	// Marshal the struct
	data, err := json.Marshal(rui)
	assert.NoError(t, err)

	// Expected JSON as a parsed structure
	expected := map[string]interface{}{
		"uid": uid.String(),
		"ids": map[string]interface{}{
			"altId": "12345",
			"email": "user@example.com",
		},
	}

	// Unmarshal both expected and actual JSON for comparison
	var actual map[string]interface{}
	err = json.Unmarshal(data, &actual)
	assert.NoError(t, err)

	assert.Equal(t, expected, actual)
}

func TestRecordUserIdentity_UnmarshalJSON(t *testing.T) {
	uid := NewUID()
	ids := atags.TagMapString{
		"email": "user@example.com",
		"altId": "12345",
	}
	rui := NewRecordUserIdentity(uid, ids)

	data := []byte(fmt.Sprintf(`"%s"`, rui.String()))
	var unmarshaledRui RecordUserIdentity

	err := json.Unmarshal(data, &unmarshaledRui)
	assert.NoError(t, err)
	assert.Equal(t, rui, unmarshaledRui)
}

func TestRecordUserIdentity_GetEmail(t *testing.T) {
	email := "user@example.com"
	rui := NewRecordUserIdentityByEmail(aemail.EmailAddress(email))

	assert.Equal(t, email, rui.GetEmail().String())
}

func TestRecordUserIdentity_FindLabel(t *testing.T) {
	label := atags.TagKey("altId")
	value := "12345"
	rui := NewRecordUserIdentityById(label, value)

	assert.Equal(t, value, rui.FindLabel(label))
	assert.Empty(t, rui.FindLabel("nonexistent"))
}

func TestRecordUserIdentity_HasMatch(t *testing.T) {
	uid := NewUID()
	email := aemail.EmailAddress("user@example.com")
	altId := "12345"

	tests := []struct {
		name     string
		rui      RecordUserIdentity
		user     RecordUserIdentity
		expected bool
	}{
		{
			name:     "Match by UID",
			rui:      NewRecordUserIdentityByUID(uid),
			user:     NewRecordUserIdentityByUID(uid),
			expected: true,
		},
		{
			name:     "Match by Email",
			rui:      NewRecordUserIdentityByEmail(email),
			user:     NewRecordUserIdentityById("email", email.String()),
			expected: true,
		},
		{
			name:     "Match by AltId",
			rui:      NewRecordUserIdentityById("altId", altId),
			user:     NewRecordUserIdentityById("altId", altId),
			expected: true,
		},
		{
			name:     "No Match - Different UID",
			rui:      NewRecordUserIdentityByUID(NewUID()),
			user:     NewRecordUserIdentityByUID(NewUID()),
			expected: false,
		},
		{
			name:     "No Match - Empty RecordUserIdentity",
			rui:      RecordUserIdentity{},
			user:     RecordUserIdentity{},
			expected: false,
		},
		{
			name:     "No Match - Different Email",
			rui:      NewRecordUserIdentityByEmail(aemail.EmailAddress("user1@example.com")),
			user:     NewRecordUserIdentityByEmail(aemail.EmailAddress("user2@example.com")),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.rui.HasMatch(tt.user)
			assert.Equal(t, tt.expected, result)
		})
	}
}
