package atags

import (
	"fmt"
)

type TagArr []*TagKeyValue

func (tm TagArr) Find(key TagKey) *TagKeyValue {
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

func (tm TagArr) HasKey(key TagKey) bool {
	return tm.Find(key) != nil
}

func (tm *TagArr) Add(tag *TagKeyValue) error {
	if tag == nil || tag.Key.IsEmpty() {
		return fmt.Errorf("tag not defined")
	}
	if tag.Key.IsEmpty() {
		return fmt.Errorf("tag.Key is empty")
	}
	if tm == nil {
		*tm = TagArr{}
	} else {
		if tm.HasKey(tag.Key) {
			return fmt.Errorf("tag.Key already added; use tag.Set to overwrite")
		}
	}
	*tm = append(*tm, tag)
	return nil
}

func (tm *TagArr) Set(tag *TagKeyValue) error {
	if tag == nil || tag.Key.IsEmpty() {
		return fmt.Errorf("tag not defined")
	}
	if tag.Key.IsEmpty() {
		return fmt.Errorf("tag.Key is empty")
	}
	if tm == nil {
		*tm = TagArr{}
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

func (tm *TagArr) Delete(key TagKey) {
	if tm == nil || len(*tm) == 0 || key.IsEmpty() {
		return
	}
	tmNew := TagArr{}
	for _, t := range *tm {
		if key == t.Key {
			continue
		}
		tmNew = append(tmNew, t)
	}
	*tm = tmNew
}

func (tm TagArr) Value(key TagKey) interface{} {
	t := tm.Find(key)
	if t == nil {
		return nil
	}
	return t.Value
}

func (tm TagArr) ToMap() TagMap {
	tmap := TagMap{}
	if tm == nil || len(tm) == 0 {
		return tmap
	}
	for _, t := range tm {
		tmap[t.Key] = t.Value
	}
	return tmap
}
