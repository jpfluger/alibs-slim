package atags

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTagArrString(t *testing.T) {
	tarr := TagArrStrings{}

	_ = tarr.Add(&TagKeyValueString{Key: "key1", Value: "value1"})
	assert.Equal(t, 1, len(tarr))

	_ = tarr.Add(&TagKeyValueString{Key: "key2", Value: "value2"})
	assert.Equal(t, 2, len(tarr))

	assert.NotNil(t, tarr.Find("key1"))
	assert.Equal(t, "value2", tarr.Find("key2").Value)

	if err := tarr.Add(&TagKeyValueString{Key: "key1", Value: "value1"}); err == nil {
		t.Errorf("should have erred because of duplicate key")
	}

	assert.Equal(t, 2, len(tarr.ToMap()))

	tarr.Delete("key1")
	assert.Equal(t, 1, len(tarr))

	tarr.Delete("key2")
	assert.Equal(t, 0, len(tarr))
}
