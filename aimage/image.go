package aimage

import (
	"encoding/base64"
	"fmt"
	"github.com/jpfluger/alibs-slim/anetwork"
	"github.com/jpfluger/alibs-slim/autils"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

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
	if err := ji.ImageInfo.Validate(); err != nil {
		return err
	}
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
	return fmt.Sprintf("data:%s;base64,%s", GetMimeType(ji.Type.String()), ji.Data)
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

	im, err := CreateFromBytes(buf, path.Base(filePath), nil)
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
	return base64.StdEncoding.DecodeString(ji.Data)
}

// MustToBytes decodes the base64 image data to bytes and panics on failure.
func (ji *Image) MustToBytes() []byte {
	b, err := ji.ToBytes()
	if err != nil {
		panic(err)
	}
	return b
}

// CreateFromFile creates an Image from a file on disk.
func CreateFromFile(filePath string, filterOptions *ImageFilterOption) (*Image, error) {
	buf, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("cannot read file '%s': %w", filePath, err)
	}
	return CreateFromBytes(buf, filepath.Base(filePath), filterOptions)
}

// CreateFromBytes creates an Image from a byte slice.
func CreateFromBytes(buf []byte, name string, filterOptions *ImageFilterOption) (*Image, error) {
	if len(buf) == 0 {
		return nil, fmt.Errorf("buffer is empty")
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

	return &Image{
		ImageInfo: ImageInfo{
			Type:         ImageType(extType),
			Name:         strings.ToLower(autils.SanitizeName(strings.TrimSuffix(name, filepath.Ext(name)))),
			OriginalName: strings.TrimSpace(name),
			Size:         len(buf),
			Url:          nil,
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

	// Fetch the image data
	resp, err := http.Get(parsedUrl.String())
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
		},
		Data: data64,
	}, nil
}
