package atags

import (
	"fmt"
)

type TagMapInt64 map[TagKey]int64

func (tm *TagMapInt64) Add(key TagKey, value int64) error {
	if tm == nil {
		*tm = map[TagKey]int64{}
	}
	if key.IsEmpty() {
		return fmt.Errorf("tag key not found")
	}
	(*tm)[key] = value
	return nil
}

func (tm *TagMapInt64) AddInt(key TagKey, value int) error {
	return tm.Add(key, int64(value))
}

func (tm TagMapInt64) HasTag(key TagKey) bool {
	if _, ok := tm[key]; ok {
		return true
	}
	return false
}

func (tm TagMapInt64) Value(key TagKey) int64 {
	if val, ok := tm[key]; ok {
		return val
	}
	return 0
}

func (tm *TagMapInt64) Remove(key TagKey) {
	if tm == nil {
		return
	}
	if tm.HasTag(key) {
		delete(*tm, key)
	}
}

func (tm TagMapInt64) FilterByType(tag TagKey) Tags {
	newTags := Tags{}
	for t, _ := range tm {
		if t.GetType() == tag.GetType() {
			newTags = append(newTags, t)
		}
	}
	return newTags
}

func (tmap TagMapInt64) ToArray() TagArrInt64 {
	tarr := TagArrInt64{}
	if tmap == nil || len(tmap) == 0 {
		return tarr
	}
	for key, val := range tmap {
		tarr = append(tarr, &TagKeyValueInt64{Key: key, Value: val})
	}
	return tarr
}
