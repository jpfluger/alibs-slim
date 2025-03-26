package atags

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTagMapString(t *testing.T) {

	var tm TagMapString

	if tm == nil {
		tm = TagMapString{}
	}

	tm.Set("test1", "test1")
	tm.Set("test2", "value2")
	tm.Set("test3", "value3")

	// duplicate
	if err := tm.Set("test1", "value1"); err != nil {
		t.Fatal(err)
	}

	if len(tm) != 3 {
		t.Fatalf("expected tagmap count of 3 but instead have %d", len(tm))
	}

	if !tm.HasTag("test1") {
		t.Fatal("does not have key 'test1'")
	}

	if !tm.HasTag("test2") {
		t.Fatal("does not have key 'test2'")
	}

	if !tm.HasTag("test3") {
		t.Fatal("does not have key 'test3'")
	}

	if tm.Value("test1") != "value1" {
		t.Fatal("value wrong for key 'test1'")
	}

	if tm.Value("test2") != "value2" {
		t.Fatal("value wrong for key 'test2'")
	}

	if tm.Value("test3") != "value3" {
		t.Fatal("value wrong for key 'test3'")
	}

	assert.Equal(t, 3, len(tm.ToArray()))

	tm.Remove("test1")
	tm.Remove("test2")
	tm.Remove("test3")

	if len(tm) != 0 {
		t.Fatalf("expected tagmap count of 0 but instead have %d", len(tm))
	}
}
