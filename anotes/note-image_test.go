package anotes

import (
	"fmt"
	"testing"
	"time"

	"github.com/jpfluger/alibs-slim/aimage"
)

// mockImages is a helper function to create a mock aimage.Images.
func mockImages(count int) aimage.Images {
	images := make(aimage.Images, count)
	for i := 0; i < count; i++ {
		images[i] = &aimage.Image{Data: "http://example.com/image" + fmt.Sprintf("%d", i+1)}
	}
	return images
}

// TestNoteImage_FixIntegrity tests the FixIntegrity method of NoteImage.
func TestNoteImage_FixIntegrity(t *testing.T) {
	tests := []struct {
		name      string
		noteImage *NoteImage
		want      bool
	}{
		{
			name:      "NilNoteImage",
			noteImage: nil,
			want:      false,
		},
		{
			name: "EmptyNoteImage",
			noteImage: &NoteImage{
				NoteText: NoteText{
					Type:   "",
					Date:   time.Time{},
					Text:   "",
					UserId: "",
				},
				Images: nil,
			},
			want: false,
		},
		{
			name: "ValidNoteImageWithText",
			noteImage: &NoteImage{
				NoteText: NoteText{
					Type:   NOTETYPE_IMAGE,
					Date:   time.Now(),
					Text:   "Sample image note",
					UserId: "user123",
				},
				Images: mockImages(1),
			},
			want: true,
		},
		{
			name: "ValidNoteImageWithImages",
			noteImage: &NoteImage{
				NoteText: NoteText{
					Type:   NOTETYPE_IMAGE,
					Date:   time.Now(),
					Text:   "",
					UserId: "user123",
				},
				Images: mockImages(2),
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.noteImage.FixIntegrity(); got != tt.want {
				t.Errorf("NoteImage.FixIntegrity() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestNoteImages_Clean tests the Clean method of NoteImages.
func TestNoteImages_Clean(t *testing.T) {
	noteImages := NoteImages{
		&NoteImage{
			NoteText: NoteText{
				Type:   NOTETYPE_IMAGE,
				Date:   time.Now(),
				Text:   "Sample image note",
				UserId: "user123",
			},
			Images: mockImages(1),
		},
		&NoteImage{
			NoteText: NoteText{
				Type:   NOTETYPE_IMAGE,
				Date:   time.Time{}, // Invalid date
				Text:   "No date",
				UserId: "user123",
			},
			Images: nil,
		},
	}

	cleanedNoteImages := noteImages.Clean()
	if len(cleanedNoteImages) != 1 {
		t.Errorf("NoteImages.Clean() did not clean the slice correctly, got %d, want 1", len(cleanedNoteImages))
	}
}
