package atags

import (
	"time"
)

type TagKeyValue struct {
	Key   TagKey      `json:"key,omitempty"`
	Value interface{} `json:"value,omitempty"`
}

type TagKeyValueString struct {
	Key   TagKey `json:"key,omitempty"`
	Value string `json:"value,omitempty"`
}

type TagKeyValueInt64 struct {
	Key   TagKey `json:"key,omitempty"`
	Value int64  `json:"value,omitempty"`
}

type TagKeyValueFloat64 struct {
	Key   TagKey  `json:"key,omitempty"`
	Value float64 `json:"value,omitempty"`
}

type TagKeyValueBool struct {
	Key   TagKey `json:"key,omitempty"`
	Value bool   `json:"value,omitempty"`
}

type TagKeyValueTime struct {
	Key   TagKey    `json:"key,omitempty"`
	Value time.Time `json:"value,omitempty"`
}
