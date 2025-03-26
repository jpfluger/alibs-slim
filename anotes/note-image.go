package anotes

import (
	"github.com/jpfluger/alibs-slim/aimage"
	"strings"
)

const NOTETYPE_IMAGE = NoteType("image") // Default NoteType for NoteImage.

// NoteImage extends Note with additional image-related information.
type NoteImage struct {
	NoteText // Embedding Note struct to inherit its fields and methods.

	// Images are usually optional.
	// A script could be written to make it mandatory.
	Images aimage.Images `json:"images,omitempty"` // The images associated with the note.
}

// FixIntegrity for NoteImage ensures the integrity of both Note and NoteImage fields.
func (ni *NoteImage) FixIntegrity() bool {
	if ni == nil {
		return false
	}
	if ni.NoteText.Type.IsEmpty() {
		ni.NoteText.Type = NOTETYPE_IMAGE // Set the default NoteType to NOTETYPE_IMAGE if empty.
	}
	if !ni.NoteText.FixIntegrity() {
		return false
	}
	if ni.Images.Count() > 0 {
		ni.Images = ni.Images.Clean() // Clean the images.
	}
	// Return true if Text is not empty or if there are images present.
	return strings.TrimSpace(ni.Text) != "" || ni.Images.Count() > 0
}

// NoteImages is a slice of pointers to NoteImage objects.
type NoteImages []*NoteImage

// Clean filters the NoteImages slice to only include items with integrity.
func (nis NoteImages) Clean() NoteImages {
	nisNew := NoteImages{}
	for _, ni := range nis {
		if ni.FixIntegrity() {
			nisNew = append(nisNew, ni)
		}
	}
	return nisNew
}
