package aconns

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConnActionType_IsEmpty(t *testing.T) {
	assert.True(t, ConnActionType("").IsEmpty())
	assert.False(t, ConnActionType("CREATE").IsEmpty())
}

func TestConnActionType_TrimSpace(t *testing.T) {
	assert.Equal(t, ConnActionType("CREATE"), ConnActionType(" CREATE ").TrimSpace())
}

func TestConnActionType_String(t *testing.T) {
	assert.Equal(t, "CREATE", ConnActionType("CREATE").String())
}

func TestConnActionType_Matches(t *testing.T) {
	assert.True(t, ConnActionType("CREATE").Matches("CREATE"))
	assert.False(t, ConnActionType("CREATE").Matches("UPGRADE"))
}

func TestConnActionType_ToStringTrimLower(t *testing.T) {
	assert.Equal(t, "create", ConnActionType(" CREATE ").ToStringTrimLower())
}

func TestConnActionType_Validate(t *testing.T) {
	assert.NoError(t, ConnActionType("CREATE").Validate())
	assert.Error(t, ConnActionType("").Validate())
	assert.Error(t, ConnActionType("CREATE!").Validate())
}

func TestConnActionTypes_IsEmpty(t *testing.T) {
	assert.True(t, ConnActionTypes{}.IsEmpty())
	assert.False(t, ConnActionTypes{"CREATE"}.IsEmpty())
}

func TestConnActionTypes_String(t *testing.T) {
	assert.Equal(t, "CREATE, UPGRADE", ConnActionTypes{"CREATE", "UPGRADE"}.String())
}

func TestConnActionTypes_ToStringArray(t *testing.T) {
	assert.Equal(t, []string{"CREATE", "UPGRADE"}, ConnActionTypes{"CREATE", "UPGRADE"}.ToStringArray())
}

func TestConnActionTypes_Find(t *testing.T) {
	assert.Equal(t, ConnActionType("CREATE"), ConnActionTypes{"CREATE", "UPGRADE"}.Find("CREATE"))
	assert.Equal(t, ConnActionType(""), ConnActionTypes{"CREATE", "UPGRADE"}.Find("DELETE"))
}

func TestConnActionTypes_HasKey(t *testing.T) {
	assert.True(t, ConnActionTypes{"CREATE", "UPGRADE"}.HasKey("CREATE"))
	assert.False(t, ConnActionTypes{"CREATE", "UPGRADE"}.HasKey("DELETE"))
}

func TestConnActionTypes_Matches(t *testing.T) {
	assert.True(t, ConnActionTypes{"CREATE", "UPGRADE"}.Matches("CREATE"))
	assert.False(t, ConnActionTypes{"CREATE", "UPGRADE"}.Matches("DELETE"))
}
