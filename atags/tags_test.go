package atags

import (
	"testing"
)

func TestTags_Add(t *testing.T) {

	var tags Tags

	if tags == nil {
		tags = Tags{}
	}

	tags.Add("test1")
	tags.Add("test2")
	tags.Add("test3")
	// duplicate
	tags.Add("test3")

	if len(tags) != 3 {
		t.Fatalf("expected tag count of 3 but instead have %d", len(tags))
	}

	tags.Remove("test1")
	tags.Remove("test2")
	tags.Remove("test3")

	if len(tags) != 0 {
		t.Fatalf("expected tag count of 0 but instead have %d", len(tags))
	}
}
