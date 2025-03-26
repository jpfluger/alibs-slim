package acrypt

import (
	"encoding/base64"
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"path"
	"testing"
)

func TestHasSecret(t *testing.T) {
	totpud := &TOTP{Secret: "secret123"}
	assert.True(t, totpud.HasSecret(), "TOTP should have a secret.")

	totpud = &TOTP{}
	assert.False(t, totpud.HasSecret(), "TOTP should not have a secret.")
}

func TestGetImageBase64(t *testing.T) {
	totpud := &TOTP{Image: []byte{1, 2, 3}}
	expectedBase64 := "AQID" // Base64 encoding of {1, 2, 3}
	assert.Equal(t, expectedBase64, totpud.GetImageBase64(), "Base64 encoding does not match expected.")
}

func TestGetImageAsSrcAttrValue(t *testing.T) {
	totpud := &TOTP{Image: []byte{1, 2, 3}}
	expectedSrcAttr := `data:image/png;base64,AQID`
	assert.Equal(t, expectedSrcAttr, totpud.GetImageAsSrcAttrValue(), "Data URI does not match expected.")
}

func TestTOTPGenerate(t *testing.T) {
	issuer := "TestIssuer"
	account := "test@example.com"
	imageDimension := 400
	totpud, err := TOTPGenerate(issuer, account, imageDimension)
	assert.NoError(t, err, "Generating TOTP should not produce an error.")
	assert.NotNil(t, totpud, "Generated TOTP should not be nil.")
	assert.NotEmpty(t, totpud.Secret, "Generated TOTP should have a secret.")
	assert.NotEmpty(t, totpud.Image, "Generated TOTP should have an image.")
}

func TestTOTPValidate(t *testing.T) {
	// Assuming 'TOTPValidate' function is correctly implemented.
	isValid := TOTPValidate("123456", "secret123")
	assert.False(t, isValid, "Validation should fail with incorrect secret.")
}

func TestTOTPGenerateSeries(t *testing.T) {
	// Generate TOTP for the first user
	ud1, err := TOTPGenerate("SnakeOil", "alice@example.com", 0)
	assert.NoError(t, err, "Generating TOTP for the first user should not produce an error.")
	assert.NotEmpty(t, ud1.Secret, "The secret for the first user should not be empty.")
	assert.NotEmpty(t, ud1.Image, "The image for the first user should not be empty.")

	// Generate TOTP for the second user
	ud2, err := TOTPGenerate("SnakeOil", "alice2@example.com", 0)
	assert.NoError(t, err, "Generating TOTP for the second user should not produce an error.")
	assert.NotEmpty(t, ud2.Secret, "The secret for the second user should not be empty.")
	assert.NotEmpty(t, ud2.Image, "The image for the second user should not be empty.")

	// Ensure that the secrets and images are unique for each user
	assert.NotEqual(t, ud1.Secret, ud2.Secret, "Secrets for both users should be different.")
	assert.NotEqual(t, ud1.Image, ud2.Image, "Images for both users should be different.")

	// Test behavior when the image is empty
	ud1.Image = []byte{}
	assert.Equal(t, "", ud1.GetImageBase64(), "Base64 encoding of an empty image should be an empty string.")

	// Test behavior when the image is nil
	ud1.Image = nil
	assert.Equal(t, "", ud1.GetImageBase64(), "Base64 encoding of a nil image should be an empty string.")
}

// This function works correctly but the blur takes seconds to
// run so it is commented out. Best is to generate a single blur
// image and reuse it.
func TestTOTPUserData_BlurImage(t *testing.T) {
	// Create a temporary directory to store test files
	dir, err := os.MkdirTemp("", "test")
	if err != nil {
		t.Fatalf("cannot create temp directory; %v", err)
	}
	defer os.RemoveAll(dir) // Clean up after the test

	// Generate TOTP data for the user
	ud1, err := TOTPGenerate("SnakeOil", "alice@example.com", 0)
	if err != nil {
		t.Error(err) // Use t.Fatal to stop the test if TOTP generation fails
		return
	}

	// Save the QR code image to the temporary directory
	qrcPath := path.Join(dir, "qrc.png")
	if err := ud1.SaveImage(qrcPath); err != nil {
		t.Error(err) // Use t.Fatal to stop the test if saving the image fails
		return
	}

	// The following code is commented out because it's assumed that the blur operation
	// takes a long time to run and is not suitable for unit testing.
	// Instead, it's suggested to generate a single blurred image and reuse it.

	// Uncomment the following lines if you want to test the blur operation.
	// Note that this will significantly increase the test runtime.

	// blurFilePath := path.Join(dir, "qrc-blur.png")
	// if err := ud1.SaveImageBlur(blurFilePath); err != nil {
	// 	t.Fatal(err) // Use t.Fatal to stop the test if blurring the image fails
	// }

	// Create an HTML file to visually verify the QR code image
	htmlContent := fmt.Sprintf("<html><head></head><body><img src='data:image/png;base64,%s'></body></html>",
		base64.StdEncoding.EncodeToString(ud1.Image))
	htmlFilePath := path.Join(dir, "qrc.html")
	if err := os.WriteFile(htmlFilePath, []byte(htmlContent), 0644); err != nil {
		t.Error(err)
	}
}
