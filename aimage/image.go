package aimage

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/jpfluger/alibs-slim/anetwork"
	"github.com/jpfluger/alibs-slim/autils"
)

// Shared HTTP client for high-traffic efficiency
var defaultHTTPClient = &http.Client{
	Timeout: 10 * time.Second,
}

// Image represents an image with metadata and encoded data.
type Image struct {
	ImageInfo
	Data            string `json:"data,omitempty"`
	ImageImportPath string `json:"imageImportPath,omitempty"`
}

// Validate checks if the image is valid.
func (ji *Image) Validate() error {
	if ji == nil {
		return fmt.Errorf("image info is nil")
	}
	validated, err := ji.ImageInfo.Validate()
	if err != nil {
		return err
	}
	ji.ImageInfo = validated
	if !ji.HasData() {
		return fmt.Errorf("image info has no data")
	}
	return nil
}

// ToImageMimeType returns the MIME type of the image.
func (ji *Image) ToImageMimeType() string {
	if ji == nil {
		return ""
	}
	return ji.ImageInfo.ToImageMimeType()
}

// ToImageData returns the image data as a base64 encoded string with MIME type prefix.
func (ji *Image) ToImageData() string {
	if ji == nil {
		return ""
	}
	return fmt.Sprintf("data:%s;base64,%s", ji.ToImageMimeType(), ji.Data)
}

// HasData checks if the image has data.
func (ji *Image) HasData() bool {
	return ji != nil && strings.TrimSpace(ji.Data) != ""
}

// LoadFileImports loads image data from a file if the image data is not already set.
func (ji *Image) LoadFileImports(dirOptions string) error {
	if ji.HasData() || ji.ImageImportPath == "" {
		return nil
	}
	filePath := filepath.Join(dirOptions, ji.ImageImportPath)
	buf, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("cannot read file '%s': %w", filePath, err)
	}

	im, err := CreateFromBytes(buf, path.Base(filePath), nil, 0)
	if err != nil {
		return err
	}
	ji.Name = im.Name
	ji.Type = im.Type
	ji.OriginalName = im.OriginalName
	ji.Size = im.Size
	ji.Data = im.Data
	ji.ImageImportPath = ""
	return nil
}

// ToBytes decodes the base64 image data to bytes.
func (ji *Image) ToBytes() ([]byte, error) {
	if ji == nil || ji.Data == "" {
		return nil, fmt.Errorf("image data is empty")
	}
	b, err := base64.StdEncoding.DecodeString(ji.Data)
	if err != nil {
		return nil, fmt.Errorf("data is not valid base64: %w", err) // Error instead of fallback
	}
	return b, nil
}

// MustToBytes decodes the base64 image data to bytes and panics on failure.
func (ji *Image) MustToBytes() []byte {
	b, err := ji.ToBytes()
	if err != nil {
		panic(err)
	}
	return b
}

// SetFromBytes sets the image data from a byte slice by encoding it to base64.
func (ji *Image) SetFromBytes(data []byte) {
	ji.Data = base64.StdEncoding.EncodeToString(data)
	ji.Size = len(data)
}

// CreateFromFile creates an Image from a file on disk (no size limit).
func CreateFromFile(filePath string, filterOptions *ImageFilterOption) (*Image, error) {
	return CreateFromFileWithLimit(filePath, filterOptions, 0)
}

// CreateFromFileWithLimit creates an Image from a file on disk with an optional upload size limit (in bytes; 0 for no limit).
func CreateFromFileWithLimit(filePath string, filterOptions *ImageFilterOption, uploadLimit int) (*Image, error) {
	buf, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("cannot read file '%s': %w", filePath, err)
	}
	return CreateFromBytes(buf, filepath.Base(filePath), filterOptions, uploadLimit)
}

const (
	LIMIT_5MB   = 5 * 1024 * 1024
	LIMIT_10MB  = 10 * 1024 * 1024
	LIMIT_100MB = 100 * 1024 * 1024
	LIMIT_500MB = 500 * 1024 * 1024
	LIMIT_1G    = 1 * 1024 * 1024 * 1024
	LIMIT_5G    = 5 * 1024 * 1024 * 1024
	LIMIT_10G   = 10 * 1024 * 1024 * 1024
)

// CreateFromBytes creates an Image from a byte slice with an optional upload size limit (in bytes; 0 for no limit).
func CreateFromBytes(buf []byte, name string, filterOptions *ImageFilterOption, uploadLimit int) (*Image, error) {
	if len(buf) == 0 {
		return nil, fmt.Errorf("buffer is empty")
	}
	if uploadLimit > 0 && len(buf) > uploadLimit {
		return nil, fmt.Errorf("file too large; max %d bytes", uploadLimit)
	}
	if name = strings.TrimSpace(name); name == "" {
		return nil, fmt.Errorf("name is empty")
	}
	ext := autils.StripExtensionPrefix(name)
	if ext == "" {
		return nil, fmt.Errorf("file name '%s' does not contain an extension", name)
	}
	if !IsAllowedExtFileType(ext, filterOptions) {
		return nil, fmt.Errorf("file type '%s' is not allowed by filter options", ext)
	}
	extType := GetCleanedExt(ext)
	if extType == "" {
		return nil, fmt.Errorf("file type '%s' does not contain a known extension", ext)
	}

	// Infer MIME type from content and fallback to extension
	detectedMime := http.DetectContentType(buf)
	mimeType := CleanMimeType(detectedMime, name)

	return &Image{
		ImageInfo: ImageInfo{
			Type:         ImageType(extType),
			Name:         strings.ToLower(autils.SanitizeName(strings.TrimSuffix(name, filepath.Ext(name)))),
			OriginalName: strings.TrimSpace(name),
			Size:         len(buf),
			Url:          nil,
			MimeType:     mimeType, // Set inferred MIME
		},
		Data: base64.StdEncoding.EncodeToString(buf),
	}, nil
}

// IsAllowedExtFileType checks if the file type or tag is in the allowed filter options.
func IsAllowedExtFileType(ext string, filterOptions *ImageFilterOption) bool {
	emt := customExtMimeTypes.FindExtMime(ext)
	if emt == nil {
		return false // Unknown extension
	}
	// If no options, then all are available.
	if !filterOptions.HasOptions() {
		return true
	}
	// Check if the extension is allowed
	if filterOptions.Types != nil && len(filterOptions.Types) > 0 {
		for _, allowedExt := range filterOptions.Types {
			if allowedExt == ImageType(ext) {
				return true
			}
		}
	}
	// Check if the tags are allowed
	if filterOptions.Tags != nil && len(filterOptions.Tags) > 0 {
		tags := emt.GetTags()
		for _, tag := range tags {
			for _, allowedTag := range filterOptions.Tags {
				if strings.ToLower(tag) == strings.ToLower(allowedTag) {
					return true
				}
			}
		}
	}
	return false
}

// CreateFromUrl fetches an image from a given URL and returns an *Image.
func CreateFromUrl(imageUrl string, doStoreData bool) (*Image, error) {
	parsedUrl, err := anetwork.ParseNetURL(imageUrl)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %s", err)
	}

	// Use shared client with timeout
	resp, err := defaultHTTPClient.Get(parsedUrl.String())
	if err != nil {
		return nil, fmt.Errorf("failed to fetch image from URL: %s", err)
	}
	defer resp.Body.Close()

	// Read the image content
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read image data: %s", err)
	}

	if len(data) == 0 {
		return nil, fmt.Errorf("image data is empty")
	}

	// Extract file name from URL path
	name := filepath.Base(parsedUrl.Path)
	if name == "" || !strings.Contains(name, ".") {
		return nil, fmt.Errorf("invalid file name extracted from URL: %s", parsedUrl.Path)
	}

	// Get file extension
	ext := strings.TrimPrefix(filepath.Ext(name), ".")
	if ext == "" {
		return nil, fmt.Errorf("unable to determine file extension from URL")
	}

	// Infer MIME type from response header and fallback to content/ext
	headerMime := resp.Header.Get("Content-Type")
	mimeType := CleanMimeType(headerMime, name)
	if mimeType == "" {
		mimeType = http.DetectContentType(data)
	}

	var data64 string
	if doStoreData {
		data64 = base64.StdEncoding.EncodeToString(data)
	}

	// Create the Image object
	return &Image{
		ImageInfo: ImageInfo{
			Type:         ImageType(ext),
			Name:         strings.TrimSuffix(name, filepath.Ext(name)),
			OriginalName: name,
			Size:         len(data),
			Url:          parsedUrl,
			MimeType:     mimeType, // Set inferred MIME
		},
		Data: data64,
	}, nil
}
