package atags

import (
	"fmt"
	"strings"
)

type TagMapString map[TagKey]string

func (tm *TagMapString) Set(key TagKey, value string) error {
	if tm == nil {
		*tm = make(TagMapString)
	}
	if key.IsEmpty() {
		return fmt.Errorf("tag key not found")
	}
	(*tm)[key] = value
	return nil
}

func (tm TagMapString) HasTag(key TagKey) bool {
	if _, ok := tm[key]; ok {
		return true
	}
	return false
}

func (tm TagMapString) Value(key TagKey) string {
	if val, ok := tm[key]; ok {
		return val
	}
	return ""
}

func (tm *TagMapString) Remove(key TagKey) {
	if tm == nil {
		return
	}
	if tm.HasTag(key) {
		delete(*tm, key)
	}
}

func (tm TagMapString) Clean() TagMapString {
	return tm.CleanWithOptions(false)
}

func (tm TagMapString) CleanWithOptions(removeEmptyValues bool) TagMapString {
	tmap := TagMapString{}
	if tm == nil || len(tm) == 0 {
		return tmap
	}
	for k, v := range tm {
		if k.IsEmpty() {
			continue
		}
		if removeEmptyValues {
			if strings.TrimSpace(v) == "" {
				continue
			}
		}
		tmap[k] = v
	}
	return tmap
}

func (tm TagMapString) FilterByType(tag TagKey) Tags {
	newTags := Tags{}
	for t, _ := range tm {
		if t.GetType() == tag.GetType() {
			newTags = append(newTags, t)
		}
	}
	return newTags
}

func (tmap TagMapString) ToArray() TagArrStrings {
	tarr := TagArrStrings{}
	if tmap == nil || len(tmap) == 0 {
		return tarr
	}
	for key, val := range tmap {
		tarr = append(tarr, &TagKeyValueString{Key: key, Value: val})
	}
	return tarr
}

func (tmap TagMapString) IsEmpty() bool {
	return tmap == nil || len(tmap) == 0
}

// Keys returns the list of tag keys in the map.
func (tmap TagMapString) Keys() TagKeys {
	keys := make(TagKeys, 0, len(tmap))
	for k := range tmap {
		keys = append(keys, k)
	}
	return keys
}

// ToStringArr returns the tag keys as a slice of strings.
func (tmap TagMapString) ToStringArr() []string {
	strs := make([]string, 0, len(tmap))
	for k := range tmap {
		strs = append(strs, string(k))
	}
	return strs
}
