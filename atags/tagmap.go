package atags

import (
	"errors"
	"github.com/gofrs/uuid/v5"
	"github.com/jpfluger/alibs-slim/atypeconvert"
	"github.com/jpfluger/alibs-slim/autils"
)

// Predefined errors for common issues.
var ErrKeyIsEmpty = errors.New("key is empty")
var ErrValueNotFound = errors.New("value not found")

// TagMap is a map that associates TagKeys with values.
type TagMap map[TagKey]interface{}

// Find retrieves a value associated with the key in the TagMap.
// Returns nil if the key is empty or not found.
func (tmap TagMap) Find(key TagKey) interface{} {
	if key.IsEmpty() {
		return nil
	}
	return tmap[key]
}

// HasValue checks if there is a value associated with the key in the TagMap.
// Returns false if the key is empty or the value is nil.
func (tmap TagMap) HasValue(key TagKey) bool {
	return !key.IsEmpty() && tmap[key] != nil
}

// Add inserts or updates a value associated with the key in the TagMap.
// If the data is nil, it removes the key from the map.
func (tmap *TagMap) Add(key TagKey, data interface{}) {
	if key.IsEmpty() {
		return
	}
	if data == nil {
		delete(*tmap, key)
	} else {
		(*tmap)[key] = data
	}
}

// Remove deletes the value associated with the key in the TagMap.
// Does nothing if the key is empty.
func (tmap *TagMap) Remove(key TagKey) {
	if !key.IsEmpty() {
		delete(*tmap, key)
	}
}

// GetValue retrieves a value associated with the key in the TagMap.
// Returns an error if the key is empty or the value is not found.
func (tmap TagMap) GetValue(key TagKey) (interface{}, error) {
	if key.IsEmpty() {
		return nil, ErrKeyIsEmpty
	}
	data := tmap.Find(key)
	if data == nil {
		return nil, ErrValueNotFound
	}
	return data, nil
}

// GetValueAsString attempts to retrieve a value as a string.
// Returns an error if the value is not found or cannot be converted.
func (tmap TagMap) GetValueAsString(key TagKey) (string, error) {
	val, err := tmap.GetValue(key)
	if err != nil {
		return "", err
	}
	return atypeconvert.ConvertToStringFrom(val)
}

// GetValueAsInt attempts to retrieve a value as an int.
// Returns an error if the value is not found or cannot be converted.
func (tmap TagMap) GetValueAsInt(key TagKey) (int, error) {
	val, err := tmap.GetValue(key)
	if err != nil {
		return 0, err
	}
	return atypeconvert.ConvertToIntFrom(val)
}

// GetValueAsFloat attempts to retrieve a value as a float64.
// Returns an error if the value is not found or cannot be converted.
func (tmap TagMap) GetValueAsFloat(key TagKey) (float64, error) {
	val, err := tmap.GetValue(key)
	if err != nil {
		return 0, err
	}
	return atypeconvert.ConvertToFloatFrom(val)
}

// GetValueAsBool attempts to retrieve a value as a bool.
// Returns an error if the value is not found or cannot be converted.
func (tmap TagMap) GetValueAsBool(key TagKey) (bool, error) {
	val, err := tmap.GetValue(key)
	if err != nil {
		return false, err
	}
	return atypeconvert.ConvertToBoolFrom(val)
}

// GetValueAsUUID attempts to retrieve a value as a UUID.
// Returns an error if the value is not found or cannot be converted.
func (tmap TagMap) GetValueAsUUID(key TagKey) (uuid.UUID, error) {
	val, err := tmap.GetValueAsString(key)
	if err != nil {
		return uuid.Nil, err
	}
	return autils.ParseUUID(val), nil
}

// ToArray converts the TagMap to an array of TagKeyValue pointers.
// Returns an empty array if the TagMap is nil or empty.
func (tmap TagMap) ToArray() TagArr {
	var tarr TagArr
	for key, val := range tmap {
		tarr = append(tarr, &TagKeyValue{Key: key, Value: val})
	}
	return tarr
}

// MergeFrom combines two TagMaps into one.
// Values from the target TagMap override those from the source TagMap.
func (tmap TagMap) MergeFrom(target TagMap) TagMap {
	newMap := make(TagMap)
	for key, val := range tmap {
		newMap[key] = val
	}
	for key, val := range target {
		newMap[key] = val
	}
	return newMap
}
