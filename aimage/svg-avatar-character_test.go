package aimage

import (
	"encoding/base64"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestGetDefaultImageDataCircleAvatar checks if the default avatar image data is generated correctly.
func TestGetDefaultImageDataCircleAvatar(t *testing.T) {
	defaultData := GetDefaultImageDataCircleAvatar()
	assert.NotEmpty(t, defaultData, "The default avatar image data should not be empty.")
}

// TestCreateImageDataCircleAvatar checks if the avatar image data is generated correctly.
func TestCreateImageDataCircleAvatar(t *testing.T) {
	character := "A"
	imageData := CreateImageDataCircleAvatar(character)
	assert.NotEmpty(t, imageData, "The avatar image data should not be empty.")

	// Decode the base64 data URI to get the SVG content.
	dataURI := strings.SplitN(imageData, ",", 2)
	if len(dataURI) != 2 {
		t.Fatal("The data URI is not in the correct format.")
	}
	svgData, err := base64.StdEncoding.DecodeString(dataURI[1])
	assert.NoError(t, err, "Decoding base64 data should not produce an error.")

	// Check if the SVG content contains the expected character (updated for simple text).
	assert.Contains(t, string(svgData), character, "The SVG content should contain the character.")
}

//// Using "github.com/gomig/avatar"
//// The avatar library uses presets. This is a different approach to what I had before but still using a SVG template.
//// Let's just compare output to the preset.
//var compareCharacter = `<path fill="#03071e" d="M70.6,93.64l-1.18-7.63h-10.77l-1.34,7.63h-16.51l14.23-59.29h19.58l12.58,59.29h-16.59Zm-3.07-20.21l-1.73-13.76c-.47-3.7-.79-7.47-1.18-11.17h-.16c-.47,3.7-.87,7.47-1.49,11.17l-2.2,13.76h6.76Z"/>`
//
//// TestAvatarLibraryIntegration checks if the avatar library generates the expected SVG content.
//func TestAvatarLibraryIntegration(t *testing.T) {
//	character := "A"
//	avCircle := avatar.NewTextAvatar(character)
//	svgContent := avCircle.InlineSVG()
//
//	assert.Contains(t, svgContent, compareCharacter, "The SVG content from the avatar library should contain the character.")
//}
/*
package aimage

import (
	"encoding/base64"
	"github.com/gomig/avatar"
	"regexp"
)

// IMAGE_NAME_AVATAR is the default name for avatar images.
const IMAGE_NAME_AVATAR = "avatar"

// imageDataCircleQuestion stores the default avatar image data.
var imageDataCircleQuestion string

// GetDefaultImageDataCircleAvatar returns the default avatar image data,
// creating it if it doesn't already exist.
func GetDefaultImageDataCircleAvatar() string {
	if imageDataCircleQuestion == "" {
		imageDataCircleQuestion = CreateImageDataCircleAvatar("?")
	}
	return imageDataCircleQuestion
}

// CreateImageDataCircleAvatar creates an SVG image data URI with a circle and a character in the center.
func CreateImageDataCircleAvatar(target string) string {
	image := CreateImageCircleAvatar(target)
	return image.ToImageData()
}

// CreateImageCircleAvatar creates an Image with a circle and a character in the center.
func CreateImageCircleAvatar(target string) *Image {
	// Use a default character if none is provided.
	if target == "" {
		target = "?"
	}
	// Create an Image struct with the SVG data.
	avCircle := avatar.NewTextAvatar(target)
	//return avCircle.Base64()

	svgBytes := []byte(RemoveSVGDescElements(avCircle.InlineSVG()))

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

// Precompile the regular expression for removing <desc> tags.
var descRegex = regexp.MustCompile(`<desc>.*?</desc>`)

// RemoveSVGDescElements removes the <desc></desc> element from an SVG string.
// Doing so reduces bloat.
func RemoveSVGDescElements(svgContent string) string {
	// Remove the <desc> tags and any content between them from the SVG content.
	return descRegex.ReplaceAllString(svgContent, "")
}
*/
