package aconns

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAuthPipeline_GetAdaptersAndIds(t *testing.T) {
	id1 := NewConnId()
	id2 := NewConnId()

	entries := AuthAdapterEntries{
		makeTestEntry(id1, 1),
		makeTestEntry(id2, 2),
	}

	pipeline := AuthPipeline{
		AUTHMETHOD_PRIMARY: entries,
	}

	adapters := pipeline.GetAdapters(AUTHMETHOD_PRIMARY)
	ids := pipeline.GetConnIds(AUTHMETHOD_PRIMARY)

	assert.Len(t, adapters, 2)
	assert.Len(t, ids, 2)
	assert.Equal(t, id1, ids[0])
	assert.Equal(t, id2, ids[1])
}

func TestAuthPipeline_Validate(t *testing.T) {
	valid := AuthPipeline{
		AUTHMETHOD_PRIMARY: AuthAdapterEntries{
			makeTestEntry(NewConnId(), 0),
			makeTestEntry(NewConnId(), 1),
		},
		AUTHMETHOD_MFA: AuthAdapterEntries{
			makeTestEntry(NewConnId(), 2),
		},
	}
	assert.NoError(t, valid.Validate())

	invalid := AuthPipeline{
		AUTHMETHOD_SSPR: AuthAdapterEntries{
			{ConnId: ConnId{}, Adapter: &DummyAdapter{}},
		},
	}
	err := invalid.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "empty ConnId")
}

func TestAuthPipeline_Methods(t *testing.T) {
	pipeline := AuthPipeline{
		AUTHMETHOD_PRIMARY: AuthAdapterEntries{},
		AUTHMETHOD_MFA:     AuthAdapterEntries{},
	}

	methods := pipeline.Methods()
	assert.Contains(t, methods, AUTHMETHOD_PRIMARY)
	assert.Contains(t, methods, AUTHMETHOD_MFA)
	assert.Len(t, methods, 2)
}
