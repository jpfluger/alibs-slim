package aimage

import (
	"encoding/base64"
	"fmt"
	"github.com/jpfluger/alibs-slim/autils"
	"net/http"
	"os"
	"strings"
)

// ToBytes decodes a base64 encoded string into a byte slice.
func ToBytes(base64Str string) ([]byte, error) {
	if base64Str == "" {
		return nil, fmt.Errorf("base64 string is empty")
	}
	data, err := base64.StdEncoding.DecodeString(base64Str)
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64: %v", err)
	}
	return data, nil
}

// ToBase64 encodes a byte slice into a base64 string.
func ToBase64(target []byte) string {
	if len(target) == 0 {
		return ""
	}
	return base64.StdEncoding.EncodeToString(target)
}

// ToBase64ImageData encodes a byte slice into a base64 string with a data URI scheme including the MIME type.
func ToBase64ImageData(target []byte, altMimeType string) string {
	if len(target) == 0 {
		return ""
	}
	mimeType := CleanMimeType(http.DetectContentType(target), altMimeType)
	return ToImageData(mimeType, ToBase64(target))
}

// ToBase64ImageDataWithMimeType creates a base64 data URI string with the specified MIME type.
func ToBase64ImageDataWithMimeType(mimeType string, target []byte) string {
	return ToImageData(mimeType, ToBase64(target))
}

// ToImageData creates a data URI string with the specified MIME type and base64 encoded data.
func ToImageData(mimeType, base64Str string) string {
	if mimeType == "" || base64Str == "" {
		return ""
	}
	return fmt.Sprintf("data:%s;base64,%s", mimeType, base64Str)
}

// ImageLoadFromFile loads an image file from the specified path and returns its MIME type and data.
func ImageLoadFromFile(filePath string) (mimeType string, data []byte, err error) {
	filePath = strings.TrimSpace(filePath)
	if filePath == "" {
		err = fmt.Errorf("file path is empty")
		return
	}

	data, err = os.ReadFile(filePath)
	if err != nil {
		return
	}

	mimeType = CleanMimeType(http.DetectContentType(data), autils.GetFileNamePartExt(filePath))
	return
}

// Base64LoadFromFile loads an image file and returns its MIME type and base64 encoded data.
func Base64LoadFromFile(filePath string) (mimeType, base64Str string, err error) {
	mimeType, data, err := ImageLoadFromFile(filePath)
	if err != nil {
		return
	}
	base64Str = ToBase64(data)
	return
}

// ImageDataLoadFromFile loads an image file and returns its MIME type and base64 data URI string.
func ImageDataLoadFromFile(filePath string) (mimeType, imageData string, err error) {
	mimeType, data, err := ImageLoadFromFile(filePath)
	if err != nil {
		return
	}
	imageData = ToBase64ImageDataWithMimeType(mimeType, data)
	return
}

// Base64SaveToFileAsBytes saves base64 encoded data to a file after decoding it.
func Base64SaveToFileAsBytes(filePath, base64Str string) error {
	if strings.TrimSpace(base64Str) == "" {
		return fmt.Errorf("base64 data is empty")
	}
	b, err := ToBytes(base64Str)
	if err != nil {
		return fmt.Errorf("failed to convert base64 data to bytes: %v", err)
	}
	err = os.WriteFile(filePath, b, autils.PATH_CHMOD_FILE)
	if err != nil {
		return fmt.Errorf("failed to write decoded base64 data to file: %v", err)
	}
	return nil
}
