package ascripts

import (
	"fmt"
	"github.com/jpfluger/alibs-slim/ajson"
)

type Scriptables []*Scriptable

func (sis Scriptables) Find(key ajson.JsonKey) *Scriptable {
	if sis == nil || len(sis) == 0 || key.IsEmpty() {
		return nil
	}
	for _, si := range sis {
		if si.Key == key {
			return si
		}
	}
	return nil
}

func (sis Scriptables) HasItem(key ajson.JsonKey) bool {
	return sis.Find(key) != nil
}

func (sis Scriptables) ToMap() (ScriptablesMap, error) {
	sim := ScriptablesMap{}
	if sis == nil || len(sis) == 0 {
		return sim, nil
	}
	for _, si := range sis {
		if sim.HasItem(si.Key) {
			return nil, fmt.Errorf("duplicate ScriptKey detected '%s'", si.Key.String())
		}
		sim[si.Key] = si
	}
	return sim, nil
}
