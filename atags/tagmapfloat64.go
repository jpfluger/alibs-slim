package atags

import (
	"fmt"
)

type TagMapFloat64 map[TagKey]float64

func (tm *TagMapFloat64) Add(key TagKey, value float64) error {
	if tm == nil {
		*tm = map[TagKey]float64{}
	}
	if key.IsEmpty() {
		return fmt.Errorf("tag key not found")
	}
	(*tm)[key] = value
	return nil
}

func (tm *TagMapFloat64) AddInt(key TagKey, value int) error {
	return tm.Add(key, float64(value))
}

// Precision issue
//func (tm *TagMapFloat64)  AddFloat(key string, value float32) error {
//	v2, err := typeconvert.ConvertToFloatFrom(value)
//	if err != nil {
//		return err
//	}
//	return tm.Add(key, v2)
//}

func (tm TagMapFloat64) HasTag(key TagKey) bool {
	if _, ok := tm[key]; ok {
		return true
	}
	return false
}

func (tm TagMapFloat64) Value(key TagKey) float64 {
	if val, ok := tm[key]; ok {
		return val
	}
	return 0
}

func (tm *TagMapFloat64) Remove(key TagKey) {
	if tm == nil {
		return
	}
	if tm.HasTag(key) {
		delete(*tm, key)
	}
}

func (tm TagMapFloat64) FilterByType(tag TagKey) Tags {
	newTags := Tags{}
	for t, _ := range tm {
		if t.GetType() == tag.GetType() {
			newTags = append(newTags, t)
		}
	}
	return newTags
}

func (tmap TagMapFloat64) ToArray() TagArrFloat64 {
	tarr := TagArrFloat64{}
	if tmap == nil || len(tmap) == 0 {
		return tarr
	}
	for key, val := range tmap {
		tarr = append(tarr, &TagKeyValueFloat64{Key: key, Value: val})
	}
	return tarr
}
