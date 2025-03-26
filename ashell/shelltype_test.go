package ashell

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShellType_IsEmpty(t *testing.T) {
	assert.True(t, ShellType("").IsEmpty())
	assert.True(t, ShellType("   ").IsEmpty())
	assert.False(t, ShellType("bash").IsEmpty())
}

func TestShellType_TrimSpace(t *testing.T) {
	assert.Equal(t, ShellType("bash"), ShellType("  bash  ").TrimSpace())
	assert.Equal(t, ShellType("sh"), ShellType("sh").TrimSpace())
}

func TestShellType_String(t *testing.T) {
	assert.Equal(t, "bash", ShellType("bash").String())
	assert.Equal(t, "sh", ShellType("sh").String())
}

func TestShellType_ToStringTrimLower(t *testing.T) {
	assert.Equal(t, "bash", ShellType("  BASH  ").ToStringTrimLower())
	assert.Equal(t, "sh", ShellType("SH").ToStringTrimLower())
}

func TestShellTypes_Contains(t *testing.T) {
	shellTypes := ShellTypes{SHELLTYPE_BASH, SHELLTYPE_SH}
	assert.True(t, shellTypes.Contains(SHELLTYPE_BASH))
	assert.False(t, shellTypes.Contains(ShellType("zsh")))
}

func TestShellTypes_Add(t *testing.T) {
	shellTypes := ShellTypes{}
	shellTypes.Add(SHELLTYPE_BASH)
	assert.True(t, shellTypes.Contains(SHELLTYPE_BASH))

	shellTypes.Add(SHELLTYPE_BASH) // Adding the same type again should not duplicate it
	assert.Len(t, shellTypes, 1)
}

func TestShellTypes_Remove(t *testing.T) {
	shellTypes := ShellTypes{SHELLTYPE_BASH, SHELLTYPE_SH}
	shellTypes.Remove(SHELLTYPE_BASH)
	assert.False(t, shellTypes.Contains(SHELLTYPE_BASH))
	assert.True(t, shellTypes.Contains(SHELLTYPE_SH))
}
