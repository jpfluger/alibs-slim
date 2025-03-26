package atags

import (
	"fmt"
)

type TagMapBool map[TagKey]bool

func (tm *TagMapBool) Add(key TagKey, value bool) error {
	if tm == nil {
		*tm = map[TagKey]bool{}
	}
	if key.IsEmpty() {
		return fmt.Errorf("tag key not found")
	}
	(*tm)[key] = value
	return nil
}

func (tm TagMapBool) HasTag(key TagKey) bool {
	if _, ok := tm[key]; ok {
		return true
	}
	return false
}

func (tm TagMapBool) Value(key TagKey) bool {
	if val, ok := tm[key]; ok {
		return val
	}
	return false
}

func (tm *TagMapBool) Remove(key TagKey) {
	if tm == nil {
		return
	}
	if tm.HasTag(key) {
		delete(*tm, key)
	}
}

func (tm TagMapBool) FilterByType(tag TagKey) Tags {
	newTags := Tags{}
	for t, _ := range tm {
		if t.GetType() == tag.GetType() {
			newTags = append(newTags, t)
		}
	}
	return newTags
}

func (tmap TagMapBool) ToArray() TagArrBools {
	tarr := TagArrBools{}
	if tmap == nil || len(tmap) == 0 {
		return tarr
	}
	for key, val := range tmap {
		tarr = append(tarr, &TagKeyValueBool{Key: key, Value: val})
	}
	return tarr
}
