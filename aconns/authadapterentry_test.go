package aconns

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func makeTestEntry(id ConnId, priority int) *AuthAdapterEntry {
	return &AuthAdapterEntry{
		ConnId:   id,
		Adapter:  &DummyAdapter{},
		Priority: priority,
	}
}

func TestAuthAdapterEntries_GetByConnId(t *testing.T) {
	id1 := NewConnId()
	id2 := NewConnId()

	entries := AuthAdapterEntries{
		makeTestEntry(id1, 1),
		makeTestEntry(id2, 2),
	}

	entry, found := entries.GetByConnId(id2)
	assert.True(t, found)
	assert.Equal(t, id2, entry.ConnId)

	_, notFound := entries.GetByConnId(NewConnId())
	assert.False(t, notFound)
}

func TestAuthAdapterEntries_GetAdapters(t *testing.T) {
	entries := AuthAdapterEntries{
		makeTestEntry(NewConnId(), 1),
		makeTestEntry(NewConnId(), 2),
		nil, // should be skipped
	}

	adapters := entries.GetAdapters()
	assert.Len(t, adapters, 2)
	for _, a := range adapters {
		assert.IsType(t, &DummyAdapter{}, a)
	}
}

func TestAuthAdapterEntries_GetConnIds(t *testing.T) {
	id1 := NewConnId()
	id2 := NewConnId()

	entries := AuthAdapterEntries{
		makeTestEntry(id1, 1),
		makeTestEntry(id2, 2),
	}

	ids := entries.GetConnIds()
	assert.Equal(t, []ConnId{id1, id2}, ids)
}

func TestAuthAdapterEntries_SortByPriority(t *testing.T) {
	id1 := NewConnId()
	id2 := NewConnId()
	id3 := NewConnId()

	entries := AuthAdapterEntries{
		makeTestEntry(id3, 3),
		makeTestEntry(id1, 1),
		makeTestEntry(id2, 2),
	}

	entries.SortByPriority()

	assert.Equal(t, id1, entries[0].ConnId)
	assert.Equal(t, id2, entries[1].ConnId)
	assert.Equal(t, id3, entries[2].ConnId)
}

func TestAuthAdapterEntries_Validate(t *testing.T) {
	valid := AuthAdapterEntries{
		makeTestEntry(NewConnId(), 0),
		makeTestEntry(NewConnId(), 1),
	}
	assert.NoError(t, valid.Validate())

	invalidNil := AuthAdapterEntries{nil}
	assert.Error(t, invalidNil.Validate())

	invalidAdapter := AuthAdapterEntries{
		{ConnId: NewConnId(), Adapter: nil, Priority: 0},
	}
	assert.Error(t, invalidAdapter.Validate())

	invalidId := AuthAdapterEntries{
		{ConnId: ConnId{}, Adapter: &DummyAdapter{}, Priority: 0},
	}
	assert.Error(t, invalidId.Validate())
}
