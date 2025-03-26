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

// cron-ological
//

// notebook-esri.git
// * scripts-linux
// * scripts-windows
//   > milsoft
//     > job.json -> point to a totally different file/directory
