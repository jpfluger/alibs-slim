package atags

import (
	"fmt"
	"strings"
)

type TagArrStrings []*TagKeyValueString

func (tm TagArrStrings) Find(key TagKey) *TagKeyValueString {
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

func (tm TagArrStrings) HasKey(key TagKey) bool {
	return tm.Find(key) != nil
}

func (tm *TagArrStrings) Add(tag *TagKeyValueString) error {
	if tag == nil || tag.Key.IsEmpty() {
		return fmt.Errorf("tag not defined")
	}
	if tag.Key.IsEmpty() {
		return fmt.Errorf("tag.Key is empty")
	}
	if tm == nil {
		*tm = TagArrStrings{}
	} else {
		if tm.HasKey(tag.Key) {
			return fmt.Errorf("tag.Key already added; use tag.Set to overwrite")
		}
	}
	*tm = append(*tm, tag)
	return nil
}

func (tm *TagArrStrings) SetByKeyValue(key TagKey, value string) error {
	return tm.Set(&TagKeyValueString{Key: key, Value: value})
}

func (tm *TagArrStrings) Set(tag *TagKeyValueString) error {
	if tag == nil || tag.Key.IsEmpty() {
		return fmt.Errorf("tag is not defined")
	}

	// Initialize if tm is nil
	if tm == nil || *tm == nil {
		*tm = TagArrStrings{}
	}

	// Update existing or append new
	for ii, t := range *tm {
		if t.Key == tag.Key {
			(*tm)[ii] = tag
			return nil // Early return
		}
	}

	// Key not found, append as new
	*tm = append(*tm, tag)
	return nil
}

func (tm *TagArrStrings) Delete(key TagKey) {
	if tm == nil || len(*tm) == 0 || key.IsEmpty() {
		return
	}
	tmNew := TagArrStrings{}
	for _, t := range *tm {
		if key == t.Key {
			continue
		}
		tmNew = append(tmNew, t)
	}
	*tm = tmNew
}

func (tm TagArrStrings) Value(key TagKey) string {
	t := tm.Find(key)
	if t == nil {
		return ""
	}
	return t.Value
}

func (tm TagArrStrings) ToMap() TagMapString {
	tmap := TagMapString{}
	if tm == nil || len(tm) == 0 {
		return tmap
	}
	for _, t := range tm {
		tmap[t.Key] = t.Value
	}
	return tmap
}

func (tm TagArrStrings) Clean() TagArrStrings {
	arr := TagArrStrings{}
	if tm == nil || len(tm) == 0 {
		return arr
	}
	for _, t := range tm {
		if t.Key.IsEmpty() {
			continue
		}
		val := strings.TrimSpace(t.Value)
		if val == "" {
			continue
		}
		arr = append(arr, t)
	}
	return arr
}
