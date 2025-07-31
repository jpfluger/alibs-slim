package rruleplus

import "time"

type RROccurrence struct {
	Time     time.Time `json:"time"`
	IsDeny   bool      `json:"isDeny,omitempty"`
	Priority int       `json:"priority,omitempty"`
	Name     string    `json:"name,omitempty"`
}

type RROccurrences []*RROccurrence

type RROccurrenceMap map[int]RROccurrences
