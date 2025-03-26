package atags

import (
	"fmt"
)

type TagArrInt64 []*TagKeyValueInt64

func (tm TagArrInt64) Find(key TagKey) *TagKeyValueInt64 {
	if tm == nil || len(tm) == 0 || key.IsEmpty() {
		return nil
	}
	for _, tag := range tm {
		if tag.Key == key {
			return tag
		}
	}
	return nil
}

func (tm TagArrInt64) HasKey(key TagKey) bool {
	return tm.Find(key) != nil
}

func (tm *TagArrInt64) Add(tag *TagKeyValueInt64) error {
	if tag == nil || tag.Key.IsEmpty() {
		return fmt.Errorf("tag not defined")
	}
	if tag.Key.IsEmpty() {
		return fmt.Errorf("tag.Key is empty")
	}
	if tm == nil {
		*tm = TagArrInt64{}
	} else {
		if tm.HasKey(tag.Key) {
			return fmt.Errorf("tag.Key already added; use tag.Set to overwrite")
		}
	}
	*tm = append(*tm, tag)
	return nil
}

func (tm *TagArrInt64) Set(tag *TagKeyValueInt64) error {
	if tag == nil || tag.Key.IsEmpty() {
		return fmt.Errorf("tag not defined")
	}
	if tag.Key.IsEmpty() {
		return fmt.Errorf("tag.Key is empty")
	}
	if tm == nil {
		*tm = TagArrInt64{}
	}
	isFound := false
	for ii, t := range *tm {
		if tag.Key == t.Key {
			isFound = true
			(*tm)[ii] = tag
		}
	}
	if !isFound {
		*tm = append(*tm, tag)
	}
	return nil
}

func (tm *TagArrInt64) Delete(key TagKey) {
	if tm == nil || len(*tm) == 0 || key.IsEmpty() {
		return
	}
	tmNew := TagArrInt64{}
	for _, t := range *tm {
		if key == t.Key {
			continue
		}
		tmNew = append(tmNew, t)
	}
	*tm = tmNew
}

func (tm TagArrInt64) Value(key TagKey) int64 {
	t := tm.Find(key)
	if t == nil {
		return 0
	}
	return t.Value
}

func (tm TagArrInt64) ToMap() TagMapInt64 {
	tmap := TagMapInt64{}
	if tm == nil || len(tm) == 0 {
		return tmap
	}
	for _, t := range tm {
		tmap[t.Key] = t.Value
	}
	return tmap
}
