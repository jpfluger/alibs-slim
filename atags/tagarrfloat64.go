package atags

import (
	"fmt"
)

type TagArrFloat64 []*TagKeyValueFloat64

func (tm TagArrFloat64) Find(key TagKey) *TagKeyValueFloat64 {
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

func (tm TagArrFloat64) HasKey(key TagKey) bool {
	return tm.Find(key) != nil
}

func (tm *TagArrFloat64) Add(tag *TagKeyValueFloat64) error {
	if tag == nil || tag.Key.IsEmpty() {
		return fmt.Errorf("tag not defined")
	}
	if tag.Key.IsEmpty() {
		return fmt.Errorf("tag.Key is empty")
	}
	if tm == nil {
		*tm = TagArrFloat64{}
	} else {
		if tm.HasKey(tag.Key) {
			return fmt.Errorf("tag.Key already added; use tag.Set to overwrite")
		}
	}
	*tm = append(*tm, tag)
	return nil
}

func (tm *TagArrFloat64) Set(tag *TagKeyValueFloat64) error {
	if tag == nil || tag.Key.IsEmpty() {
		return fmt.Errorf("tag not defined")
	}
	if tag.Key.IsEmpty() {
		return fmt.Errorf("tag.Key is empty")
	}
	if tm == nil {
		*tm = TagArrFloat64{}
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

func (tm *TagArrFloat64) Delete(key TagKey) {
	if tm == nil || len(*tm) == 0 || key.IsEmpty() {
		return
	}
	tmNew := TagArrFloat64{}
	for _, t := range *tm {
		if key == t.Key {
			continue
		}
		tmNew = append(tmNew, t)
	}
	*tm = tmNew
}

func (tm TagArrFloat64) Value(key TagKey) float64 {
	t := tm.Find(key)
	if t == nil {
		return 0
	}
	return t.Value
}

func (tm TagArrFloat64) ToMap() TagMapFloat64 {
	tmap := TagMapFloat64{}
	if tm == nil || len(tm) == 0 {
		return tmap
	}
	for _, t := range tm {
		tmap[t.Key] = t.Value
	}
	return tmap
}
