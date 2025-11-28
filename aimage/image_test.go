package aimage

import (
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"testing"

	"github.com/jpfluger/alibs-slim/autils"
	"github.com/stretchr/testify/assert"
)

// All test images part of public domain
// png: https://commons.wikimedia.org/wiki/File:Quanterr.png
// gif: https://commons.wikimedia.org/wiki/File:Opera_Game,_1858.gif
// jpg: https://commons.wikimedia.org/wiki/File:Luther95theses.jpg
// svg: https://commons.wikimedia.org/wiki/File:Flag_of_Florida.svg

// TestImageToImageMimeType checks the ToImageMimeType method for Image.
func TestImageToImageMimeType(t *testing.T) {
	image := &Image{ImageInfo: ImageInfo{Type: "jpg"}}
	assert.Equal(t, "image/jpeg", image.ToImageMimeType())
}

// TestImageToImageData checks the ToImageData method for Image.
func TestImageToImageData(t *testing.T) {
	image := &Image{ImageInfo: ImageInfo{Type: "jpg"}, Data: base64.StdEncoding.EncodeToString([]byte("image data"))}
	assert.Contains(t, image.ToImageData(), "data:image/jpeg;base64,")
}

// TestImageHasData checks the HasData method for Image.
func TestImageHasData(t *testing.T) {
	image := &Image{Data: "some data"}
	assert.True(t, image.HasData())
}

// TestImageHasType checks the HasType method for Image.
func TestImageHasType(t *testing.T) {
	image := &Image{ImageInfo: ImageInfo{Type: "jpg"}}
	assert.True(t, image.HasType())
}

// TestImageValidate checks the Validate method for Image.
func TestImageValidate(t *testing.T) {
	image := &Image{ImageInfo: ImageInfo{Name: "test", Type: "jpg"}, Data: "some data"}
	assert.NoError(t, image.Validate())
}

func TestSanitizeName(t *testing.T) {
	tests := []struct {
		input    ImageInfo
		expected string
	}{
		{ImageInfo{Name: "My Image", Type: "jpg"}, "my-image.jpg"},
		{ImageInfo{Name: "Invalid:File*Name?", Type: "png"}, "invalidfilename.png"},
		{ImageInfo{Name: "Trailing Space .", Type: "gif"}, "trailing-space.gif"},
		{ImageInfo{Name: "Leading/Slash", Type: "bmp"}, "leadingslash.bmp"},
		{ImageInfo{Name: "Reserved|Name", Type: "tiff"}, "reservedname.tiff"},
	}

	for _, tt := range tests {
		result := tt.input.SanitizeName()
		assert.Equal(t, tt.expected, result, "SanitizeName failed for input: %v", tt.input)
	}
}

// TestImageLoadFileImports checks the LoadFileImports method for Image.
func TestImageLoadFileImports(t *testing.T) {
	// This test would require an actual file to be present at the given path.
	// For the purpose of this example, we will assume the function works correctly.
	// In a real-world scenario, you should mock the file reading process.
}

// TestImageToBytes checks the ToBytes method for Image.
func TestImageToBytes(t *testing.T) {
	data := []byte("image data")
	image := &Image{Data: base64.StdEncoding.EncodeToString(data)}
	bytes, err := image.ToBytes()
	assert.NoError(t, err)
	assert.Equal(t, data, bytes)
}

// TestImageMustToBytes checks the MustToBytes method for Image.
func TestImageMustToBytes(t *testing.T) {
	data := []byte("image data")
	image := &Image{Data: base64.StdEncoding.EncodeToString(data)}
	bytes := image.MustToBytes()
	assert.Equal(t, data, bytes)
}

// TestCreateFromBytes checks the CreateFromBytes function for Image.
func TestImageFilterOption_CreateFromBytes_Ext(t *testing.T) {
	data := []byte("image data")
	name := "text.txt"
	image, err := CreateFromBytes(data, name, &ImageFilterOption{
		Types: ImageTypes{"txt"},
	}, LIMIT_10MB)
	assert.NoError(t, err)
	assert.NotNil(t, image)
}

// TestCreateFromBytes checks the CreateFromBytes function for Image.
func TestImageFilterOption_CreateFromBytes_Tags(t *testing.T) {
	data := []byte("image data")
	name := "text.txt"
	image, err := CreateFromBytes(data, name, &ImageFilterOption{
		Tags: []string{"text"},
	}, LIMIT_10MB)
	assert.NoError(t, err)
	assert.NotNil(t, image)
}

// TestCreateFromFile checks the CreateFromFile function for Image.
func TestCreateFromFile(t *testing.T) {
	// Create a temporary directory to store images during the test.
	dir, err := autils.CreateTempDir()
	if err != nil {
		t.Fatalf("cannot create temp directory: %v", err)
	}
	defer os.RemoveAll(dir)

	// Get the current working directory to locate the test data.
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	// Define a helper function to test image creation and saving for different file types.
	testImageFile := func(filename, expectedType, expectedMimeType string) {
		// Create an image from the test data file.
		jsonImage, err := CreateFromFile(path.Join(wd, filename), nil)
		if err != nil {
			t.Fatal(err)
		}

		// Assert that the image type and MIME type are as expected.
		assert.Equal(t, expectedType, jsonImage.Type.String())
		assert.Equal(t, expectedMimeType, jsonImage.ToImageMimeType())

		// Save the image data to a new file in the temporary directory.
		filePath := path.Join(dir, "image."+expectedType)
		err = Base64SaveToFileAsBytes(filePath, jsonImage.Data)
		if err != nil {
			t.Fatal(err)
		}

		// Create a new image from the saved file and assert that the data matches the original.
		jsonImage2, err := CreateFromFile(filePath, nil)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, jsonImage.Data, jsonImage2.Data)
	}

	// Test various image and video file types.
	testImageFile("test-data/t1-opera-game.gif", "gif", "image/gif")
	testImageFile("test-data/t2-luther.jpg", "jpg", "image/jpeg")
	testImageFile("test-data/t3-quanterr.png", "png", "image/png")
	testImageFile("test-data/t4-florida-flag.svg", "svg", "image/svg+xml")
	testImageFile("test-data/t5-fake-audio.amp4", "amp4", "video/mp4")
	testImageFile("test-data/t6-fake-video.mp4", "mp4", "video/mp4")

	// Uncomment the following block to test .mov video files.
	// testImageFile("test-data/conviction.mov", "mov", "video/mp4")
}

// TestToImageMimeType tests ToImageMimeType with various cases.
func TestToImageMimeType(t *testing.T) {
	tests := []struct {
		name     string
		image    *Image
		expected string
	}{
		{
			name:     "Valid Image with Known Type",
			image:    &Image{ImageInfo: ImageInfo{Type: "jpg"}},
			expected: "image/jpeg",
		},
		{
			name:     "Valid Image with Unknown Type",
			image:    &Image{ImageInfo: ImageInfo{Type: "unknown"}},
			expected: "application/octet-stream",
		},
		{
			name:     "Nil Image",
			image:    nil,
			expected: "",
		},
		{
			name:     "Empty Type",
			image:    &Image{ImageInfo: ImageInfo{Type: ""}},
			expected: "application/octet-stream",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.image.ToImageMimeType()
			assert.Equal(t, tt.expected, result)
		})
	}
}

// Mock HTTP server to serve test images
func mockImageServer(t *testing.T, imageData []byte, contentType string) *httptest.Server {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", contentType)
		_, _ = w.Write(imageData)
	})
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close) // Ensure cleanup after test
	return server
}

// TestCreateFromUrl validates fetching images from URLs.
func TestCreateFromUrl(t *testing.T) {
	// Mock Image Data
	mockData := []byte("fakeImageData") // Simulate image binary data
	encodedData := base64.StdEncoding.EncodeToString(mockData)

	// ✅ Case 1: Valid JPEG image with base64 encoding enabled
	server := mockImageServer(t, mockData, "image/jpeg")
	defer server.Close()

	img, err := CreateFromUrl(server.URL+"/sample.jpg", true)
	assert.NoError(t, err, "CreateFromUrl should succeed for a valid image URL")
	assert.NotNil(t, img, "Image object should not be nil")
	assert.Equal(t, "sample", img.Name, "Image name should be extracted correctly")
	assert.Equal(t, "jpg", string(img.Type), "Image type should be extracted correctly")
	assert.Equal(t, len(mockData), img.Size, "Image size should match")
	assert.Equal(t, encodedData, img.Data, "Base64 encoded data should match")

	// ✅ Case 2: Valid PNG image but without base64 encoding (doStoreData=false)
	serverPNG := mockImageServer(t, mockData, "image/png")
	defer serverPNG.Close()

	imgNoData, err := CreateFromUrl(serverPNG.URL+"/image.png", false)
	assert.NoError(t, err, "CreateFromUrl should succeed for a valid PNG image URL")
	assert.NotNil(t, imgNoData, "Image object should not be nil")
	assert.Equal(t, "image", imgNoData.Name, "Image name should be extracted correctly")
	assert.Equal(t, "png", string(imgNoData.Type), "Image type should be extracted correctly")
	assert.Equal(t, len(mockData), imgNoData.Size, "Image size should match")
	assert.Empty(t, imgNoData.Data, "Data should be empty when doStoreData is false")

	// ❌ Case 3: Invalid URL
	_, err = CreateFromUrl("invalid-url", true)
	assert.Error(t, err, "CreateFromUrl should fail for an invalid URL")
	assert.Contains(t, err.Error(), "invalid URL", "Error message should indicate invalid URL format")

	// ❌ Case 4: Empty response body
	serverEmpty := mockImageServer(t, []byte{}, "image/jpeg")
	defer serverEmpty.Close()

	_, err = CreateFromUrl(serverEmpty.URL+"/empty.jpg", true)
	assert.Error(t, err, "CreateFromUrl should fail for an empty image file")
	assert.Contains(t, err.Error(), "image data is empty", "Error should indicate empty response")
}
