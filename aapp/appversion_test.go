package aapp

import (
	"testing"

	"github.com/Masterminds/semver/v3"
	"github.com/stretchr/testify/assert"
)

// TestAppVersion tests the loading and formatting of AppVersion from JSON.
func TestAppVersion(t *testing.T) {
	const testAppVersionJSON = `{
  "name": "NAME",
  "version": "0.1.0",
  "about": "ABOUT",
  "owner": "OWNER",
  "legalMark": "LEGALMARK",
  "buildName": "name",
  "buildType": "debug",
  "deploymentType": "dev"
}`

	appVersion, err := LoadAppVersionFromBytes([]byte(testAppVersionJSON))
	if err != nil {
		t.Error(err)
		return
	}

	s, err := appVersion.Format(APPVERSION_FORMAT_JSON)
	assert.Equal(t, testAppVersionJSON, s)
	s, err = appVersion.Format(APPVERSION_FORMAT_VERSION_ONLY)
	assert.Equal(t, "v0.1.0", s)
	s, err = appVersion.Format(APPVERSION_FORMAT_NO_V)
	assert.Equal(t, "0.1.0", s)
	s, err = appVersion.Format(APPVERSION_FORMAT_APP_ONLY)
	assert.Equal(t, "name", s)
	s, err = appVersion.Format(APPVERSION_FORMAT_ABOUT_ONLY)
	assert.Equal(t, "ABOUT", s)
	s, err = appVersion.Format(APPVERSION_FORMAT_APP_AT_VERSION)
	assert.Equal(t, "name@0.1.0", s)
}

// TestAppVersion_Validate tests the validation of AppVersion fields.
func TestAppVersion_Validate(t *testing.T) {
	appVersion := &AppVersion{
		Name:      "TestApp",
		Version:   semver.MustParse("1.0.0"),
		About:     "A test application",
		Owner:     "TestOwner",
		LegalMark: "TestLegalMark",
	}

	err := appVersion.Validate()
	assert.NoError(t, err)
}

// TestLoadAppVersionFromBytes tests the loading of AppVersion from a byte slice.
func TestLoadAppVersionFromBytes(t *testing.T) {
	const versionJSON = `{"name":"TestApp","version":"1.0.0","about":"A test application","owner":"TestOwner","legalMark":"TestLegalMark"}`
	appVersion, err := LoadAppVersionFromBytes([]byte(versionJSON))
	assert.NoError(t, err)
	assert.Equal(t, "TestApp", appVersion.Name)
	assert.Equal(t, "1.0.0", appVersion.Version.String())
	assert.Equal(t, "A test application", appVersion.About)
	assert.Equal(t, "TestOwner", appVersion.Owner)
	assert.Equal(t, "TestLegalMark", appVersion.LegalMark)
}
