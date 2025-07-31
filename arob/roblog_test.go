package arob

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestROBLog_GetType_NormalizesOnce(t *testing.T) {
	log := &ROBLog{Type: "crit", Message: "critical failure"}
	normalized := log.GetType()
	assert.Equal(t, ROBTYPE_CRITICAL, normalized)

	log2 := ROBLog{Type: "crit", Message: "critical failure"}
	normalized = log2.GetType()
	assert.Equal(t, ROBTYPE_CRITICAL, normalized)
}

func TestROBLog_String(t *testing.T) {
	log := ROBLog{Type: "warn", Message: "watch this"}
	assert.Equal(t, "[warning] watch this", log.String())
}

func TestROBLogs_HasLogType(t *testing.T) {
	logs := ROBLogs{
		{Type: "debug", Message: "init"},
		{Type: "err", Message: "failure"},
	}

	assert.True(t, logs.HasLogType(ROBTYPE_ERROR))
	assert.False(t, logs.HasLogType(ROBTYPE_CRITICAL))
}

func TestROBLogs_FilterByType(t *testing.T) {
	logs := ROBLogs{
		{Type: "info", Message: "ok"},
		{Type: "warn", Message: "warning issued"},
		{Type: "debug", Message: "trace"},
		{Type: "error", Message: "failure"},
	}

	filtered := logs.FilterByType(ROBTYPE_WARNING, ROBTYPE_DEBUG)
	assert.Len(t, filtered, 2)
	assert.Equal(t, "[warning] warning issued", filtered[0].String())
	assert.Equal(t, "[debug] trace", filtered[1].String())
}

func TestROBLogs_ToStringArray(t *testing.T) {
	logs := ROBLogs{
		{Type: "notice", Message: "processing"},
		{Type: "err", Message: "bad state"},
	}
	strs := logs.ToStringArray()
	expected := []string{"[notice] processing", "[error] bad state"}
	assert.Equal(t, expected, strs)
}
