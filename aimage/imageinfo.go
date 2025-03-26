package aimage

import (
	"fmt"
	"github.com/jpfluger/alibs-slim/anetwork"
	"github.com/jpfluger/alibs-slim/autils"
	"strings"
)

type ImageInfo struct {
	Type         ImageType        `json:"type,omitempty"`         // Type is the extension type with no preceeding period (".").
	Name         string           `json:"name,omitempty"`         // Human-readable name
	OriginalName string           `json:"originalName,omitempty"` // Original filename
	Size         int              `json:"size,omitempty"`         // Size of the raw data in bytes
	Url          *anetwork.NetURL `json:"url,omitempty"`          // Optional reference for external data
}

// Validate checks if the image is valid.
func (ii *ImageInfo) Validate() error {
	if ii == nil {
		return fmt.Errorf("image info is nil")
	}
	ii.OriginalName = strings.TrimSpace(ii.OriginalName)
	ii.Type = ii.Type.TrimSpace()
	if ii.Type.IsEmpty() {
		if ii.OriginalName != "" {
			ext := autils.StripExtensionPrefix(ii.OriginalName)
			if ext != "" {
				extType := GetCleanedExt(ext)
				if extType != "" {
					ii.Type = ImageType(extType)
				}
			}
		}
		if ii.Type.IsEmpty() {
			return fmt.Errorf("image type is empty")
		}
	}
	ii.Name = strings.TrimSpace(ii.Name)
	if ii.Name == "" {
		return fmt.Errorf("image name is empty")
	}
	if ii.Size < 0 {
		return fmt.Errorf("image size is negative")
	}
	if ii.Url != nil && !ii.Url.IsUrl() {
		return fmt.Errorf("image url is invalid")
	}
	return nil
}

// HasType checks if the image has a type set.
func (ii *ImageInfo) HasType() bool {
	return ii != nil && !ii.Type.IsEmpty()
}

// ToImageMimeType returns the MIME type of the image.
func (ii *ImageInfo) ToImageMimeType() string {
	if ii == nil {
		return ""
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
