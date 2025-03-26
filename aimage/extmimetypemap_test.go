package aimage

import (
	"github.com/jpfluger/alibs-slim/autils"
	"mime"
	"strings"
	"testing"
)

func TestExtMimeType_GetCleanedExt(t *testing.T) {
	tests := []struct {
		name     string
		input    ExtMimeType
		expected string
	}{
		{"No CleanedExt", ExtMimeType{Ext: "jpg", Mime: "image/jpeg"}, "jpg"},
		{"With CleanedExt", ExtMimeType{Ext: "tif", Mime: "image/tiff", CleanedExt: "tiff"}, "tiff"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.input.GetCleanedExt()
			if result != tt.expected {
				t.Errorf("GetCleanedExt() = %q; expected %q", result, tt.expected)
			}
		})
	}
}

func TestExtMimeTypeMap_FindMime(t *testing.T) {
	originalExtMimeTypes := customExtMimeTypes
	defer func() {
		customExtMimeTypes = originalExtMimeTypes
	}()

	customExtMimeTypes = ExtMimeTypeMap{
		"jpg": {Ext: "jpg", Mime: "image/jpeg"},
		"tif": {Ext: "tif", Mime: "image/tiff", CleanedExt: "tiff"},
	}

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"Existing Extension", "jpg", "image/jpeg"},
		{"Non-Existing Extension", "png", ""},
		{"Extension With Cleaning", "tif", "image/tiff"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := customExtMimeTypes.FindMime(tt.input)
			if result != tt.expected {
				t.Errorf("FindMime(%q) = %q; expected %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestExtMimeTypeMap_GetCleanedExt(t *testing.T) {
	originalExtMimeTypes := customExtMimeTypes
	defer func() {
		customExtMimeTypes = originalExtMimeTypes
	}()

	customExtMimeTypes = ExtMimeTypeMap{
		"jpg": {Ext: "jpg", Mime: "image/jpeg"},
		"tif": {Ext: "tif", Mime: "image/tiff", CleanedExt: "tiff"},
	}

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"Existing Extension", "jpg", "jpg"},
		{"Extension With Cleaning", "tif", "tiff"},
		{"Non-Existing Extension", "png", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := customExtMimeTypes.GetCleanedExt(tt.input)
			if result != tt.expected {
				t.Errorf("GetCleanedExt(%q) = %q; expected %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestSetExtMimeType(t *testing.T) {
	originalExtMimeTypes := customExtMimeTypes
	defer func() {
		customExtMimeTypes = originalExtMimeTypes
	}()

	customExtMimeTypes = ExtMimeTypeMap{}

	tests := []struct {
		name       string
		ext        string
		mimeType   string
		cleanedExt string
		tags       []string
		expected   ExtMimeType
	}{
		{
			"Add New Extension",
			"jpg",
			"image/jpeg",
			"",
			nil,
			ExtMimeType{Ext: "jpg", Mime: "image/jpeg", CleanedExt: "jpg"},
		},
		{
			"Add With CleanedExt",
			"tif",
			"image/tiff",
			"tiff",
			nil,
			ExtMimeType{Ext: "tif", Mime: "image/tiff", CleanedExt: "tiff"},
		},
		{
			"Update Existing Extension",
			"jpg",
			"image/custom-jpeg",
			"jpeg",
			nil,
			ExtMimeType{Ext: "jpg", Mime: "image/custom-jpeg", CleanedExt: "jpeg"},
		},
		{
			"Add With Leading Dot",
			".png",
			"image/png",
			"",
			nil,
			ExtMimeType{Ext: "png", Mime: "image/png", CleanedExt: "png"},
		},
		{
			"Add With Space in Mime",
			"gif",
			" image/gif ",
			"",
			nil,
			ExtMimeType{Ext: "gif", Mime: "image/gif", CleanedExt: "gif"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetExtMimeType(tt.ext, tt.mimeType, nil, []string{tt.cleanedExt}...)
			result, exists := customExtMimeTypes[ImageType(autils.StripExtensionPrefix(tt.ext))]
			if !exists {
				t.Errorf("Extension %q not added to customExtMimeTypes", tt.ext)
			}
			if result.Mime != tt.expected.Mime && result.CleanedExt != tt.expected.CleanedExt && result.Ext != tt.expected.Ext {
				t.Errorf("SetExtMimeType(%q, %q, %v) = %+v; expected %+v", tt.ext, tt.mimeType, tt.cleanedExt, result, tt.expected)
			}
		})
	}
}

func TestCleanMimeType(t *testing.T) {
	originalExtMimeTypes := customExtMimeTypes
	defer func() {
		customExtMimeTypes = originalExtMimeTypes
	}()

	customExtMimeTypes = ExtMimeTypeMap{
		"jpg": {Ext: "jpg", Mime: "image/jpeg", CleanedExt: "jpeg"},
		"tif": {Ext: "tif", Mime: "image/tiff", CleanedExt: "tiff"},
	}

	tests := []struct {
		name           string
		mimeType       string
		fileNameOrExt  string
		expectedResult string
	}{
		{
			"Known Extension",
			"",
			"image.jpg",
			"image/jpeg",
		},
		{
			"Known Extension With CleanedExt",
			"",
			"file.tif",
			"image/tiff",
		},
		{
			"Unknown Extension",
			"",
			"file.unknown",
			"",
		},
		{
			"Fallback to Provided MimeType",
			"application/pdf",
			"file.pdf",
			"application/pdf",
		},
		{
			"Fallback to Go's Built-in Detection",
			"",
			"file.css",
			strings.TrimSpace(mime.TypeByExtension(".css")),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CleanMimeType(tt.mimeType, tt.fileNameOrExt)
			if result != tt.expectedResult {
				t.Errorf("CleanMimeType(%q, %q) = %q; expected %q", tt.mimeType, tt.fileNameOrExt, result, tt.expectedResult)
			}
		})
	}
}
