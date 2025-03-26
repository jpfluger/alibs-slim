package atags

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTagArrInt64(t *testing.T) {
	tarr := TagArrInt64{}

	_ = tarr.Add(&TagKeyValueInt64{Key: "key1", Value: 1})
	assert.Equal(t, 1, len(tarr))

	_ = tarr.Add(&TagKeyValueInt64{Key: "key2", Value: 2})
	assert.Equal(t, 2, len(tarr))

	assert.NotNil(t, tarr.Find("key1"))
	assert.Equal(t, int64(2), tarr.Find("key2").Value)

	if err := tarr.Add(&TagKeyValueInt64{Key: "key1", Value: 1}); err == nil {
		t.Errorf("should have erred because of duplicate key")
	}

	assert.Equal(t, 2, len(tarr.ToMap()))

	tarr.Delete("key1")
	assert.Equal(t, 1, len(tarr))

	tarr.Delete("key2")
	assert.Equal(t, 0, len(tarr))
}
