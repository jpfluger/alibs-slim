package aimage

import (
	"encoding/base64"
	"strings"
	"testing"

	"github.com/gomig/avatar"
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

	// Check if the SVG content contains the expected character.
	assert.Contains(t, string(svgData), compareCharacter, "The SVG content should contain the character.")
}

// The avatar library uses presets. This is a different approach to what I had before but still using a SVG template.
// Let's just compare output to the preset.
var compareCharacter = `<path fill="#03071e" d="M70.6,93.64l-1.18-7.63h-10.77l-1.34,7.63h-16.51l14.23-59.29h19.58l12.58,59.29h-16.59Zm-3.07-20.21l-1.73-13.76c-.47-3.7-.79-7.47-1.18-11.17h-.16c-.47,3.7-.87,7.47-1.49,11.17l-2.2,13.76h6.76Z"/>`

// TestAvatarLibraryIntegration checks if the avatar library generates the expected SVG content.
func TestAvatarLibraryIntegration(t *testing.T) {
	character := "A"
	avCircle := avatar.NewTextAvatar(character)
	svgContent := avCircle.InlineSVG()

	assert.Contains(t, svgContent, compareCharacter, "The SVG content from the avatar library should contain the character.")
}
