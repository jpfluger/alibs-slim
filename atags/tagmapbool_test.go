package atags

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTagMapBool(t *testing.T) {

	var tm TagMapBool

	if tm == nil {
		tm = TagMapBool{}
	}

	tm.Add("test1", false)
	tm.Add("test2", true)
	tm.Add("test3", false)

	// duplicate
	if err := tm.Add("test1", true); err != nil {
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

	if tm.Value("test1") != true {
		t.Fatal("value wrong for key 'test1'")
	}

	if tm.Value("test2") != true {
		t.Fatal("value wrong for key 'test2'")
	}

	if tm.Value("test3") != false {
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
