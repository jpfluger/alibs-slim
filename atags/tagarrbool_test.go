package atags

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTagArrBool(t *testing.T) {
	tarr := TagArrBools{}

	_ = tarr.Add(&TagKeyValueBool{Key: "key1", Value: true})
	assert.Equal(t, 1, len(tarr))

	_ = tarr.Add(&TagKeyValueBool{Key: "key2", Value: true})
	assert.Equal(t, 2, len(tarr))

	assert.NotNil(t, tarr.Find("key1"))
	assert.Equal(t, true, tarr.Find("key2").Value)

	if err := tarr.Add(&TagKeyValueBool{Key: "key1", Value: true}); err == nil {
		t.Errorf("should have erred because of duplicate key")
	}

	assert.Equal(t, 2, len(tarr.ToMap()))

	tarr.Delete("key1")
	assert.Equal(t, 1, len(tarr))

	tarr.Delete("key2")
	assert.Equal(t, 0, len(tarr))
}
