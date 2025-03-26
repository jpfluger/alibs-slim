package atags

import (
	"fmt"
	"time"
)

type TagMapTime map[TagKey]time.Time

func (tm *TagMapTime) Add(key TagKey, value time.Time) error {
	if tm == nil {
		*tm = TagMapTime{}
	}
	if key.IsEmpty() {
		return fmt.Errorf("tag key not found")
	}
	(*tm)[key] = value
	return nil
}

func (tm TagMapTime) HasTag(key TagKey) bool {
	if _, ok := tm[key]; ok {
		return true
	}
	return false
}

func (tm TagMapTime) Value(key TagKey) time.Time {
	if val, ok := tm[key]; ok {
		return val
	}
	return time.Time{}
}

func (tm *TagMapTime) Remove(key TagKey) {
	if tm == nil {
		return
	}
	if tm.HasTag(key) {
		delete(*tm, key)
	}
}

func (tm TagMapTime) FilterByType(tag TagKey) Tags {
	newTags := Tags{}
	for t, _ := range tm {
		if t.GetType() == tag.GetType() {
			newTags = append(newTags, t)
		}
	}
	return newTags
}

func (tmap TagMapTime) ToArray() TagArrTimes {
	tarr := TagArrTimes{}
	if tmap == nil || len(tmap) == 0 {
		return tarr
	}
	for key, val := range tmap {
		tarr = append(tarr, &TagKeyValueTime{Key: key, Value: val})
	}
	return tarr
}
