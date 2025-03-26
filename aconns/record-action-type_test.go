package aconns

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRecordActionType_IsEmpty(t *testing.T) {
	assert.True(t, RecordActionType(" ").IsEmpty())
	assert.False(t, RecordActionType("insert").IsEmpty())
}

func TestRecordActionType_TrimSpace(t *testing.T) {
	assert.Equal(t, RecordActionType("insert"), RecordActionType(" insert ").TrimSpace())
}

func TestRecordActionType_String(t *testing.T) {
	assert.Equal(t, "insert", RecordActionType(" insert ").String())
}

func TestRecordActionType_Constants(t *testing.T) {
	assert.Equal(t, RecordActionType("insert"), REC_ACTION_INSERT)
	assert.Equal(t, RecordActionType("update"), REC_ACTION_UPDATE)
	assert.Equal(t, RecordActionType("delete"), REC_ACTION_DELETE)
	assert.Equal(t, RecordActionType("upsert"), REC_ACTION_UPSERT)
	assert.Equal(t, RecordActionType("import"), REC_EVENT_IMPORT)
	assert.Equal(t, RecordActionType("adminfix"), REC_EVENT_ADMINFIX)
}
