package atags

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestTagMapTime(t *testing.T) {

	var tm TagMapTime

	if tm == nil {
		tm = TagMapTime{}
	}

	t1 := time.Now()
	t2 := time.Now()
	tm.Add("test1", t1)
	tm.Add("test2", t2)
	tm.Add("test3", time.Time{})

	// duplicate
	if err := tm.Add("test1", t1); err != nil {
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

	if tm.Value("test1") != t1 {
		t.Fatal("value wrong for key 'test1'")
	}

	if tm.Value("test2") != t2 {
		t.Fatal("value wrong for key 'test2'")
	}

	if !tm.Value("test3").IsZero() {
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
