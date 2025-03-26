package anotes

import (
	"github.com/jpfluger/alibs-slim/azb"
)

const NOTETYPE_IMAGE_FLAG = NoteType("image-flag") // Default NoteType for

// NoteImageFlag extends NoteImage with additional flag-related information.
type NoteImageFlag struct {
	NoteImage // Embedding NoteImage struct to inherit its fields and methods.

	// FlagType is optional and a manual setting.
	// The Flag is used to bring attention to users about something more "permanent".
	// The ideal state is for the FlagType to be empty ("").
	// FlagType uses rules in alerts.types.*.
	FlagType azb.ZBType `json:"flagType,omitempty"` // The flag type associated with the note.
}

// FixIntegrity for NoteImageFlag ensures the integrity of NoteImage and the FlagType.
func (nif *NoteImageFlag) FixIntegrity() bool {
	if nif == nil {
		return false
	}
	if nif.NoteImage.Type.IsEmpty() {
		nif.NoteImage.Type = NOTETYPE_IMAGE_FLAG // Set the default NoteType to NOTETYPE_IMAGE_FLAG if empty.
	}
	if !nif.NoteImage.FixIntegrity() {
		return false
	}
	// Additional integrity checks can be added here if needed.
	return true // Assuming the integrity is valid if it passes the NoteImage checks.
}

// NoteImageFlags is a slice of pointers to NoteImageFlag objects.
type NoteImageFlags []*NoteImageFlag

// Clean filters the NoteImageFlags slice to only include items with integrity.
func (nifs NoteImageFlags) Clean() NoteImageFlags {
	nifsNew := NoteImageFlags{}
	for _, nif := range nifs {
		if nif.FixIntegrity() {
			nifsNew = append(nifsNew, nif)
		}
	}
	return nifsNew
}
