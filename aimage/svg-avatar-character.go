package aimage

import (
	"encoding/base64"
	"fmt"
	"strings"
	"sync" // Added for RWMutex
)

// IMAGE_NAME_AVATAR is the default name for avatar images.
const IMAGE_NAME_AVATAR = "avatar"

// defaultAvatarData stores the default avatar image data.
var defaultAvatarData string
var defaultAvatarMu sync.RWMutex // Added for thread-safe lazy init

// GetDefaultImageDataCircleAvatar returns the default avatar image data,
// creating it if it doesn't already exist (thread-safe).
func GetDefaultImageDataCircleAvatar() string {
	defaultAvatarMu.RLock()
	if defaultAvatarData != "" {
		data := defaultAvatarData
		defaultAvatarMu.RUnlock()
		return data
	}
	defaultAvatarMu.RUnlock()

	defaultAvatarMu.Lock()
	defer defaultAvatarMu.Unlock()
	if defaultAvatarData == "" {
		defaultAvatarData = CreateImageDataCircleAvatar("?")
	}
	return defaultAvatarData
}

// CreateImageDataCircleAvatar creates an SVG image data URI with a circle and a character in the center.
func CreateImageDataCircleAvatar(target string) string {
	image := CreateImageCircleAvatar(target)
	return image.ToImageData()
}

// CreateImageCircleAvatar creates an Image with a circle and a character in the center.
// Simplified to remove external dependency; uses basic text for letter (can extend with paths if needed).
func CreateImageCircleAvatar(target string) *Image {
	if target == "" {
		target = "?"
	}
	target = strings.ToUpper(string([]rune(target)[0])) // Take first character, uppercase

	// Simple SVG template (replaced library for reliability; add letter paths if fancy shapes needed)
	svg := fmt.Sprintf(`
<svg width="128" height="128" viewBox="0 0 128 128" xmlns="http://www.w3.org/2000/svg">
  <circle cx="64" cy="64" r="64" fill="#007bff"/>
  <text x="64" y="64" font-size="64" text-anchor="middle" dy=".35em" fill="white" font-family="sans-serif">%s</text>
</svg>`, target)

	svgBytes := []byte(svg) // No need for desc removal

	return &Image{
		ImageInfo: ImageInfo{
			Type:         "svg",
			Name:         IMAGE_NAME_AVATAR,
			OriginalName: "",
			Size:         len(svgBytes),
			Url:          nil,
		},
		Data: base64.StdEncoding.EncodeToString(svgBytes),
	}
}
