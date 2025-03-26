package anotes

import (
	"github.com/jpfluger/alibs-slim/azb"
	"strings"
)

const NOTETYPE_FLAG NoteType = "flag"

// NoteFlag extends Note with additional flag-related information.
type NoteFlag struct {
	NoteText // Embedding Note struct to inherit its fields and methods.

	// FlagType is optional and a manual setting.
	// The Flag is used to bring attention to users about something more "permanent".
	// The ideal state is for the FlagType to be empty ("").
	// FlagType uses rules in alerts.types.*.
	FlagType azb.ZBType `json:"flagType,omitempty"`
}

// FixIntegrity for NoteFlag ensures the integrity of both Note and NoteFlag fields.
func (nf *NoteFlag) FixIntegrity() bool {
	if nf == nil {
		return false
	}
	if nf.NoteText.Type.IsEmpty() {
		nf.NoteText.Type = NOTETYPE_FLAG // Set the default NoteType to NOTETYPE_FLAG if empty.
	}
	if !nf.NoteText.FixIntegrity() {
		return false
	}
	if strings.TrimSpace(nf.Text) != "" {
		return true
	}
	return false
}

// NoteFlags is a slice of pointers to NoteFlag objects.
type NoteFlags []*NoteFlag

// Clean filters the NoteFlags slice to only include items with integrity.
func (nfs NoteFlags) Clean() NoteFlags {
	nfsNew := NoteFlags{}
	for _, nf := range nfs {
		if nf.FixIntegrity() {
			nfsNew = append(nfsNew, nf)
		}
	}
	return nfsNew
}
