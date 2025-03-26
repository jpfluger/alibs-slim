package aconns

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewDbColumnNames(t *testing.T) {
	names := NewDbColumnNames("id", "name", "age")
	assert.Equal(t, DbColumnNames{"id", "name", "age"}, names)
}

func TestDbColumnNameIsEmpty(t *testing.T) {
	var name DbColumnName = "  "
	assert.True(t, name.IsEmpty())

	name = "id"
	assert.False(t, name.IsEmpty())
}

func TestDbColumnNameTrimSpace(t *testing.T) {
	name := DbColumnName("  id  ").TrimSpace()
	assert.Equal(t, DbColumnName("id"), name)
}

func TestDbColumnNamesFind(t *testing.T) {
	names := DbColumnNames{"id", "name", "age"}
	assert.Equal(t, DbColumnName("name"), names.Find("name"))
	assert.Equal(t, DbColumnName(""), names.Find("not_exist"))
}

func TestDbColumnNamesHas(t *testing.T) {
	names := DbColumnNames{"id", "name", "age"}
	assert.True(t, names.Has("name"))
	assert.False(t, names.Has("not_exist"))
}

func TestDbColumnNamesLength(t *testing.T) {
	names := DbColumnNames{"id", "name", "age"}
	assert.Equal(t, 3, names.Length())

	var emptyNames DbColumnNames
	assert.Equal(t, 0, emptyNames.Length())
}

func TestDbColumnNamesAdd(t *testing.T) {
	names := DbColumnNames{"id", "name"}.Add("age", "name")
	assert.Equal(t, DbColumnNames{"id", "name", "age"}, names)
}

func TestDbColumnNamesRemove(t *testing.T) {
	names := DbColumnNames{"id", "name", "age"}.Remove("name")
	assert.Equal(t, DbColumnNames{"id", "age"}, names)
}

func TestDbColumnNamesEnsure(t *testing.T) {
	names := DbColumnNames{"id", "name"}.Ensure("age", "name")
	assert.Equal(t, DbColumnNames{"id", "name", "age"}, names)
}
