package atags

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTagMapFloat64(t *testing.T) {

	var tm TagMapFloat64

	if tm == nil {
		tm = TagMapFloat64{}
	}

	tm.Add("test1", 1.1)
	tm.Add("test2", 11.1)
	tm.Add("test3", 12.1)

	// duplicate
	if err := tm.Add("test1", 2.1); err != nil {
		t.Fatal(err)
	}
	if err := tm.AddInt("test1", 10); err != nil {
		t.Fatal(err)
	}
	if err := tm.Add("test1", 10.1); err != nil {
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

	if tm.Value("test1") != 10.1 {
		t.Fatal("value wrong for key 'test1'")
	}

	if tm.Value("test2") != 11.1 {
		t.Fatal("value wrong for key 'test2'")
	}

	if tm.Value("test3") != 12.1 {
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
