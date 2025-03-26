package atags

import (
	"fmt"
)

type TagArrBools []*TagKeyValueBool

func (tm TagArrBools) Find(key TagKey) *TagKeyValueBool {
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

func (tm TagArrBools) HasKey(key TagKey) bool {
	return tm.Find(key) != nil
}

func (tm *TagArrBools) Add(tag *TagKeyValueBool) error {
	if tag == nil || tag.Key.IsEmpty() {
		return fmt.Errorf("tag not defined")
	}
	if tag.Key.IsEmpty() {
		return fmt.Errorf("tag.Key is empty")
	}
	if tm == nil {
		*tm = TagArrBools{}
	} else {
		if tm.HasKey(tag.Key) {
			return fmt.Errorf("tag.Key already added; use tag.Set to overwrite")
		}
	}
	*tm = append(*tm, tag)
	return nil
}

func (tm *TagArrBools) Set(tag *TagKeyValueBool) error {
	if tag == nil || tag.Key.IsEmpty() {
		return fmt.Errorf("tag not defined")
	}
	if tag.Key.IsEmpty() {
		return fmt.Errorf("tag.Key is empty")
	}
	if tm == nil {
		*tm = TagArrBools{}
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

func (tm *TagArrBools) Delete(key TagKey) {
	if tm == nil || len(*tm) == 0 || key.IsEmpty() {
		return
	}
	tmNew := TagArrBools{}
	for _, t := range *tm {
		if key == t.Key {
			continue
		}
		tmNew = append(tmNew, t)
	}
	*tm = tmNew
}

func (tm TagArrBools) Value(key TagKey) bool {
	t := tm.Find(key)
	if t == nil {
		return false
	}
	return t.Value
}

func (tm TagArrBools) ToMap() TagMapBool {
	tmap := TagMapBool{}
	if tm == nil || len(tm) == 0 {
		return tmap
	}
	for _, t := range tm {
		tmap[t.Key] = t.Value
	}
	return tmap
}
