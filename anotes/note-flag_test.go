package anotes

import (
	"testing"
	"time"

	"github.com/jpfluger/alibs-slim/azb"
)

// TestNoteFlag_FixIntegrity tests the FixIntegrity method of NoteFlag.
func TestNoteFlag_FixIntegrity(t *testing.T) {
	tests := []struct {
		name     string
		noteFlag *NoteFlag
		want     bool
	}{
		{
			name:     "NilNoteFlag",
			noteFlag: nil,
			want:     false,
		},
		{
			name: "EmptyNoteFlag",
			noteFlag: &NoteFlag{
				NoteText: NoteText{
					Type:   "",
					Date:   time.Time{},
					Text:   "",
					UserId: "",
				},
			},
			want: false,
		},
		{
			name: "ValidNoteFlag",
			noteFlag: &NoteFlag{
				NoteText: NoteText{
					Type:   NOTETYPE_FLAG,
					Date:   time.Now(),
					Text:   "Attention needed",
					UserId: "user123",
				},
				FlagType: azb.ZBType("permanent"),
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.noteFlag.FixIntegrity(); got != tt.want {
				t.Errorf("NoteFlag.FixIntegrity() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestNoteFlags_Clean tests the Clean method of NoteFlags.
func TestNoteFlags_Clean(t *testing.T) {
	noteFlags := NoteFlags{
		&NoteFlag{
			NoteText: NoteText{
				Type:   NOTETYPE_FLAG,
				Date:   time.Now(),
				Text:   "Attention needed",
				UserId: "user123",
			},
			FlagType: azb.ZBType("permanent"),
		},
		&NoteFlag{
			NoteText: NoteText{
				Type:   NOTETYPE_FLAG,
				Date:   time.Time{}, // Invalid date
				Text:   "No date",
				UserId: "user123",
			},
		},
	}

	cleanedNoteFlags := noteFlags.Clean()
	if len(cleanedNoteFlags) != 1 {
		t.Errorf("NoteFlags.Clean() did not clean the slice correctly, got %d, want 1", len(cleanedNoteFlags))
	}
}
