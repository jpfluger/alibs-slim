package aimage

import (
	"fmt"
	"strings"

	"github.com/jpfluger/alibs-slim/anetwork"
	"github.com/jpfluger/alibs-slim/autils"
)

type ImageInfo struct {
	Type         ImageType        `json:"type,omitempty"`         // Type is the extension type with no preceeding period (".").
	Name         string           `json:"name,omitempty"`         // Human-readable name
	OriginalName string           `json:"originalName,omitempty"` // Original filename
	Size         int              `json:"size,omitempty"`         // Size of the raw data in bytes
	Url          *anetwork.NetURL `json:"url,omitempty"`          // Optional reference for external data
	MimeType     string           `json:"mimeType,omitempty"`     // Optional explicit MIME type override (e.g., "text/css")
}

// Validate checks if the image is valid and returns a validated copy without mutating the original.
func (ii ImageInfo) Validate() (ImageInfo, error) {
	validated := ii // Copy
	validated.OriginalName = strings.TrimSpace(validated.OriginalName)
	validated.Type = validated.Type.TrimSpace()
	if validated.Type.IsEmpty() {
		if validated.OriginalName != "" {
			ext := autils.StripExtensionPrefix(validated.OriginalName)
			if ext != "" {
				extType := GetCleanedExt(ext)
				if extType != "" {
					validated.Type = ImageType(extType)
				}
			}
		}
		if validated.Type.IsEmpty() && validated.MimeType == "" { // Allow MIME-only if no type
			return validated, fmt.Errorf("image type and mime type are both empty")
		}
	}
	validated.Name = strings.TrimSpace(validated.Name)
	if validated.Name == "" {
		return validated, fmt.Errorf("image name is empty")
	}
	if validated.Size < 0 {
		return validated, fmt.Errorf("image size is negative")
	}
	if validated.Url != nil && !validated.Url.IsUrl() {
		return validated, fmt.Errorf("image url is invalid")
	}
	validated.MimeType = strings.TrimSpace(validated.MimeType) // Trim new field
	return validated, nil
}

// HasType checks if the image has a type set.
func (ii *ImageInfo) HasType() bool {
	return ii != nil && !ii.Type.IsEmpty()
}

// ToImageMimeType returns the MIME type of the image, preferring explicit MimeType if set.
func (ii *ImageInfo) ToImageMimeType() string {
	if ii == nil {
		return ""
	}
	if ii.MimeType != "" {
		return ii.MimeType // Prefer override
	}
	mimeType := GetMimeType(ii.Type.String())
	if mimeType == "" {
		return "application/octet-stream" // Default MIME type for unknown binary data
	}
	return mimeType
}

// SanitizeName returns a file-system-safe and URL-friendly file name.
func (ii ImageInfo) SanitizeName() string {
	return fmt.Sprintf("%s.%s", strings.ToLower(autils.SanitizeName(ii.Name)), ii.Type.String())
}
