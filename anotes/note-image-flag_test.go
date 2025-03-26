package anotes

import (
	"testing"
	"time"

	"github.com/jpfluger/alibs-slim/azb"
)

// TestNoteImageFlag_FixIntegrity tests the FixIntegrity method of NoteImageFlag.
func TestNoteImageFlag_FixIntegrity(t *testing.T) {
	tests := []struct {
		name          string
		noteImageFlag *NoteImageFlag
		want          bool
	}{
		{
			name:          "NilNoteImageFlag",
			noteImageFlag: nil,
			want:          false,
		},
		{
			name: "EmptyNoteImageFlag",
			noteImageFlag: &NoteImageFlag{
				NoteImage: NoteImage{
					NoteText: NoteText{
						Type:   "",
						Date:   time.Time{},
						Text:   "",
						UserId: "",
					},
				},
			},
			want: false,
		},
		{
			name: "ValidNoteImageFlag",
			noteImageFlag: &NoteImageFlag{
				NoteImage: NoteImage{
					NoteText: NoteText{
						Type:   NOTETYPE_IMAGE_FLAG,
						Date:   time.Now(),
						Text:   "Attention needed",
						UserId: "user123",
					},
				},
				FlagType: azb.ZBType("permanent"),
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.noteImageFlag.FixIntegrity(); got != tt.want {
				t.Errorf("NoteImageFlag.FixIntegrity() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestNoteImageFlags_Clean tests the Clean method of NoteImageFlags.
func TestNoteImageFlags_Clean(t *testing.T) {
	noteImageFlags := NoteImageFlags{
		&NoteImageFlag{
			NoteImage: NoteImage{
				NoteText: NoteText{
					Type:   NOTETYPE_IMAGE_FLAG,
					Date:   time.Now(),
					Text:   "Attention needed",
					UserId: "user123",
				},
			},
			FlagType: azb.ZBType("permanent"),
		},
		&NoteImageFlag{
			NoteImage: NoteImage{
				NoteText: NoteText{
					Type:   NOTETYPE_IMAGE_FLAG,
					Date:   time.Time{}, // Invalid date
					Text:   "No date",
					UserId: "user123",
				},
			},
		},
	}

	cleanedNoteImageFlags := noteImageFlags.Clean()
	if len(cleanedNoteImageFlags) != 1 {
		t.Errorf("NoteImageFlags.Clean() did not clean the slice correctly, got %d, want 1", len(cleanedNoteImageFlags))
	}
}
