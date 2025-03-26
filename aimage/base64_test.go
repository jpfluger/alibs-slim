package aimage

import (
	"github.com/stretchr/testify/assert"
	"os"
	"path"
	"testing"
)

// All test images part of public domain
// png: https://commons.wikimedia.org/wiki/File:Quanterr.png
// gif: https://commons.wikimedia.org/wiki/File:Opera_Game,_1858.gif
// jpg: https://commons.wikimedia.org/wiki/File:Luther95theses.jpg
// svg: https://commons.wikimedia.org/wiki/File:Flag_of_Florida.svg

// TestBase64Files verifies that the Base64LoadFromFile function correctly identifies the MIME type of various image files.
func TestBase64Files(t *testing.T) {
	// Get the current working directory to locate the test data.
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}

	// Define a helper function to test loading and MIME type identification for image files.
	testLoadAndCheckMIME := func(filename, expectedMIME string) {
		mimeType, _, err := Base64LoadFromFile(path.Join(wd, filename))
		if err != nil {
			t.Errorf("Failed to load file '%s': %v", filename, err)
		}
		assert.Equal(t, expectedMIME, mimeType, "MIME type should match for file '%s'", filename)
	}

	// Test various image file types.
	testLoadAndCheckMIME("test-data/t1-opera-game.gif", "image/gif")
	testLoadAndCheckMIME("test-data/t2-luther.jpg", "image/jpeg")
	testLoadAndCheckMIME("test-data/t3-quanterr.png", "image/png")
	testLoadAndCheckMIME("test-data/t4-florida-flag.svg", "image/svg+xml")
}

// TestToBytes checks if the ToBytes function correctly decodes a base64 string.
func TestToBytes(t *testing.T) {
	encoded := "aGVsbG8=" // "hello" in base64
	expected := []byte("hello")

	decoded, err := ToBytes(encoded)
	assert.NoError(t, err)
	assert.Equal(t, expected, decoded)
}

// TestToBase64 checks if the ToBase64 function correctly encodes bytes to a base64 string.
func TestToBase64(t *testing.T) {
	data := []byte("hello")
	expected := "aGVsbG8="

	encoded := ToBase64(data)
	assert.Equal(t, expected, encoded)
}

// TestToBase64ImageData checks if the ToBase64ImageData function correctly creates a base64 image data URI.
func TestToBase64ImageData(t *testing.T) {
	data := []byte("hello")
	altMimeType := "text/plain"

	imageData := ToBase64ImageData(data, altMimeType)
	assert.Contains(t, imageData, "data:text/plain; charset=utf-8;base64,aGVsbG8=")
}

// TestToImageData checks if the ToImageData function correctly creates a data URI string.
func TestToImageData(t *testing.T) {
	mimeType := "text/plain"
	base64Str := "aGVsbG8="

	imageData := ToImageData(mimeType, base64Str)
	assert.Equal(t, "data:text/plain;base64,aGVsbG8=", imageData)
}
