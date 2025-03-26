package acrypt

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/anthonynsimon/bild/blur"
	"github.com/pquerna/otp/totp"
	"image"
	"image/png"
	"os"
	"strings"
)

// TOTP struct holds the secret and image for a TOTP (Time-based One-Time Password).
type TOTP struct {
	Secret string `json:"secret,omitempty"` // The TOTP secret key
	Image  []byte `json:"image,omitempty"`  // PNG image bytes of the QR code
}

// HasSecret checks if the TOTP struct has a secret set.
func (tot *TOTP) HasSecret() bool {
	return tot != nil && tot.Secret != ""
}

// GetImageBase64 returns the base64-encoded string of the TOTP image.
func (tot *TOTP) GetImageBase64() string {
	return base64.StdEncoding.EncodeToString(tot.Image)
}

// GetImageAsSrcAttrValue returns the TOTP image as a data URI suitable for HTML src attribute.
func (tot *TOTP) GetImageAsSrcAttrValue() string {
	return fmt.Sprintf(`data:image/png;base64,%s`, tot.GetImageBase64())
}

// BlurImage applies a blur effect to the TOTP image and returns the new image bytes.
func (tot *TOTP) BlurImage() ([]byte, error) {
	img, _, err := image.Decode(bytes.NewReader(tot.Image))
	if err != nil {
		return nil, fmt.Errorf("failed to decode raw image; %v", err)
	}

	imgRGBA := blur.Box(img, 9.0)

	buf := new(bytes.Buffer)
	if err := png.Encode(buf, imgRGBA); err != nil {
		return nil, fmt.Errorf("failed to encode blurred image; %v", err)
	}

	return buf.Bytes(), nil
}

// SaveImage saves the TOTP image to the specified file path.
func (tot *TOTP) SaveImage(filepath string) error {
	if strings.TrimSpace(filepath) == "" {
		return fmt.Errorf("filepath is empty")
	}

	img, _, err := image.Decode(bytes.NewReader(tot.Image))
	if err != nil {
		return fmt.Errorf("failed to decode raw image; %v", err)
	}

	f, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("failed to create file; %v", err)
	}
	defer f.Close()

	if err := png.Encode(f, img); err != nil {
		return fmt.Errorf("failed to write file; %v", err)
	}

	return nil
}

// SaveImageBlur saves a blurred version of the TOTP image to the specified file path.
func (tot *TOTP) SaveImageBlur(filepath string) error {
	if strings.TrimSpace(filepath) == "" {
		return fmt.Errorf("filepath is empty")
	}

	blurredImage, err := tot.BlurImage()
	if err != nil {
		return fmt.Errorf("failed to blur image; %v", err)
	}

	img, _, err := image.Decode(bytes.NewReader(blurredImage))
	if err != nil {
		return fmt.Errorf("failed to decode blurred image; %v", err)
	}

	f, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("failed to create file; %v", err)
	}
	defer f.Close()

	if err := png.Encode(f, img); err != nil {
		return fmt.Errorf("failed to write file; %v", err)
	}

	return nil
}

// TOTPGenerate generates a new TOTP object including the secret and QR code image.
func TOTPGenerate(issuer string, account string, imageDimension int) (*TOTP, error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      issuer,
		AccountName: account,
	})
	if err != nil {
		return nil, err
	}

	ud := &TOTP{Secret: key.Secret()}

	if imageDimension < 100 {
		imageDimension = 400
	}

	var buf bytes.Buffer
	img, err := key.Image(imageDimension, imageDimension)
	if err != nil {
		return nil, fmt.Errorf("failed to generate image; %v", err)
	}
	if err := png.Encode(&buf, img); err != nil {
		return nil, fmt.Errorf("failed to encode png; %v", err)
	}

	ud.Image = buf.Bytes()

	return ud, nil
}

// TOTPValidate validates a submitted TOTP against the system's secret.
func TOTPValidate(submittedSecret string, systemSecret string) bool {
	if submittedSecret == "" || systemSecret == "" {
		return false
	}
	return totp.Validate(submittedSecret, systemSecret)
}
