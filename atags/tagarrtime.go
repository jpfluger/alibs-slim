package atags

import (
	"fmt"
	"time"
)

type TagArrTimes []*TagKeyValueTime

func (tm TagArrTimes) Find(key TagKey) *TagKeyValueTime {
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

func (tm TagArrTimes) HasKey(key TagKey) bool {
	return tm.Find(key) != nil
}

func (tm *TagArrTimes) Add(tag *TagKeyValueTime) error {
	if tag == nil || tag.Key.IsEmpty() {
		return fmt.Errorf("tag not defined")
	}
	if tag.Key.IsEmpty() {
		return fmt.Errorf("tag.Key is empty")
	}
	if tm == nil {
		*tm = TagArrTimes{}
	} else {
		if tm.HasKey(tag.Key) {
			return fmt.Errorf("tag.Key already added; use tag.Set to overwrite")
		}
	}
	*tm = append(*tm, tag)
	return nil
}

func (tm *TagArrTimes) Set(tag *TagKeyValueTime) error {
	if tag == nil || tag.Key.IsEmpty() {
		return fmt.Errorf("tag not defined")
	}
	if tag.Key.IsEmpty() {
		return fmt.Errorf("tag.Key is empty")
	}
	if tm == nil {
		*tm = TagArrTimes{}
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

func (tm *TagArrTimes) Delete(key TagKey) {
	if tm == nil || len(*tm) == 0 || key.IsEmpty() {
		return
	}
	tmNew := TagArrTimes{}
	for _, t := range *tm {
		if key == t.Key {
			continue
		}
		tmNew = append(tmNew, t)
	}
	*tm = tmNew
}

func (tm TagArrTimes) Value(key TagKey) time.Time {
	t := tm.Find(key)
	if t == nil {
		return time.Time{}
	}
	return t.Value
}

func (tm TagArrTimes) ToMap() TagMapTime {
	tmap := TagMapTime{}
	if tm == nil || len(tm) == 0 {
		return tmap
	}
	for _, t := range tm {
		tmap[t.Key] = t.Value
	}
	return tmap
}
