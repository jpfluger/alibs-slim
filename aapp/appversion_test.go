package aapp

import (
	"embed"
	"github.com/jpfluger/alibs-slim/autils"
	"testing"

	"github.com/Masterminds/semver/v3"
	"github.com/stretchr/testify/assert"
)

//go:embed test_data
var testFS embed.FS

func TestAppVersion_Format(t *testing.T) {
	tests := []struct {
		name     string
		app      *AppVersion
		format   AppVersionFormat
		expected string
		wantErr  bool
		errMsg   string
	}{
		{
			name: "json_format",
			app: &AppVersion{
				Name:           "TestApp",
				Version:        semver.MustParse("1.0.0"),
				Title:          "Test Application",
				About:          "A test application",
				Owner:          "TestOwner",
				LegalMark:      "TestLegalMark",
				BuildName:      "testapp",
				BuildType:      BUILDTYPE_DEBUG,
				DeploymentType: DEPLOYMENTTYPE_DEV,
			},
			format: APPVERSIONFORMAT_JSON,
			expected: `{
  "name": "TestApp",
  "version": "1.0.0",
  "title": "Test Application",
  "about": "A test application",
  "owner": "TestOwner",
  "legalMark": "TestLegalMark",
  "buildName": "testapp",
  "buildType": "debug",
  "deploymentType": "dev"
}`,
		},
		{
			name: "version_only",
			app: &AppVersion{
				Version:   semver.MustParse("1.0.0"),
				BuildName: "testapp",
			},
			format:   APPVERSIONFORMAT_VERSION_ONLY,
			expected: "v1.0.0",
		},
		{
			name: "no_v",
			app: &AppVersion{
				Version:   semver.MustParse("1.0.0"),
				BuildName: "testapp",
			},
			format:   APPVERSIONFORMAT_NO_V,
			expected: "1.0.0",
		},
		{
			name: "app_only",
			app: &AppVersion{
				BuildName: "testapp",
			},
			format:   APPVERSIONFORMAT_APP_ONLY,
			expected: "testapp",
		},
		{
			name: "about_only",
			app: &AppVersion{
				About: "A test application",
			},
			format:   APPVERSIONFORMAT_ABOUT_ONLY,
			expected: "A test application",
		},
		{
			name: "app_at_version",
			app: &AppVersion{
				Version:   semver.MustParse("1.0.0"),
				BuildName: "testapp",
			},
			format:   APPVERSIONFORMAT_APP_AT_VERSION,
			expected: "testapp@1.0.0",
		},
		{
			name: "build",
			app: &AppVersion{
				Version:        semver.MustParse("1.0.0"),
				BuildName:      "testapp",
				BuildType:      BUILDTYPE_DEBUG,
				DeploymentType: DEPLOYMENTTYPE_DEV,
			},
			format:   APPVERSIONFORMAT_BUILD,
			expected: "testapp@1.0.0,debug-dev",
		},
		{
			name: "nil_version",
			app: &AppVersion{
				BuildName: "testapp",
			},
			format:  APPVERSIONFORMAT_VERSION_ONLY,
			wantErr: true,
			errMsg:  "version cannot be nil",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := tc.app.Format(tc.format)
			if tc.wantErr {
				assert.ErrorContains(t, err, tc.errMsg, "Expected error containing %q", tc.errMsg)
				return
			}
			assert.NoError(t, err, "Unexpected error")
			assert.Equal(t, tc.expected, result, "Format result mismatch")
		})
	}
}

func TestAppVersion_Validate(t *testing.T) {
	ResetAppVersion() // Ensure clean global state
	defer ResetAppVersion()

	tests := []struct {
		name    string
		app     *AppVersion
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid",
			app: &AppVersion{
				Name:           "TestApp",
				Version:        semver.MustParse("1.0.0"),
				Owner:          "TestOwner",
				LegalMark:      "TestLegalMark",
				BuildName:      "testapp",
				BuildType:      BUILDTYPE_DEBUG,
				DeploymentType: DEPLOYMENTTYPE_DEV,
			},
		},
		{
			name: "empty_name",
			app: &AppVersion{
				Version:        autils.MustNewVersionPtr("1.0.0"),
				Owner:          "TestOwner",
				LegalMark:      "TestLegalMark",
				BuildName:      "testapp",
				BuildType:      BUILDTYPE_DEBUG,
				DeploymentType: DEPLOYMENTTYPE_DEV,
			},
			wantErr: true,
			errMsg:  "name cannot be empty",
		},
		{
			name: "empty_owner",
			app: &AppVersion{
				Name:           "TestApp",
				Version:        semver.MustParse("1.0.0"),
				LegalMark:      "TestLegalMark",
				BuildName:      "testapp",
				BuildType:      BUILDTYPE_DEBUG,
				DeploymentType: DEPLOYMENTTYPE_DEV,
			},
			wantErr: true,
			errMsg:  "owner cannot be empty",
		},
		{
			name: "empty_legal_mark",
			app: &AppVersion{
				Name:           "TestApp",
				Version:        semver.MustParse("1.0.0"),
				Owner:          "TestOwner",
				BuildName:      "testapp",
				BuildType:      BUILDTYPE_DEBUG,
				DeploymentType: DEPLOYMENTTYPE_DEV,
			},
			wantErr: true,
			errMsg:  "legal mark cannot be empty",
		},
		{
			name: "nil_version",
			app: &AppVersion{
				Name:           "TestApp",
				Owner:          "TestOwner",
				LegalMark:      "TestLegalMark",
				BuildName:      "testapp",
				BuildType:      BUILDTYPE_DEBUG,
				DeploymentType: DEPLOYMENTTYPE_DEV,
			},
			wantErr: true,
			errMsg:  "version cannot be nil",
		},
		{
			name: "invalid_build_type",
			app: &AppVersion{
				Name:           "TestApp",
				Version:        semver.MustParse("1.0.0"),
				Owner:          "TestOwner",
				LegalMark:      "TestLegalMark",
				BuildName:      "testapp",
				BuildType:      "invalid@type",
				DeploymentType: DEPLOYMENTTYPE_DEV,
			},
			wantErr: true,
			errMsg:  "invalid build type",
		},
		{
			name: "invalid_deployment_type",
			app: &AppVersion{
				Name:           "TestApp",
				Version:        semver.MustParse("1.0.0"),
				Owner:          "TestOwner",
				LegalMark:      "TestLegalMark",
				BuildName:      "testapp",
				BuildType:      BUILDTYPE_DEBUG,
				DeploymentType: "invalid@type",
			},
			wantErr: true,
			errMsg:  "invalid deployment type",
		},
		{
			name: "empty_build_name_with_global",
			app: &AppVersion{
				Name:           "TestApp",
				Version:        semver.MustParse("1.0.0"),
				Owner:          "TestOwner",
				LegalMark:      "TestLegalMark",
				BuildType:      BUILDTYPE_DEBUG,
				DeploymentType: DEPLOYMENTTYPE_DEV,
			},
			wantErr: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ResetAppVersion() // Reset global state for each test
			err := tc.app.Validate()
			if tc.wantErr {
				assert.ErrorContains(t, err, tc.errMsg, "Expected error containing %q", tc.errMsg)
				return
			}
			assert.NoError(t, err, "Unexpected error")
			if tc.app.BuildName.IsEmpty() {
				assert.Equal(t, BuildName(autils.ToStringTrimLower(tc.app.Name)), tc.app.BuildName, "BuildName default mismatch")
			}
		})
	}
}

func TestAppVersion_GetTitle(t *testing.T) {
	tests := []struct {
		name     string
		app      *AppVersion
		expected string
	}{
		{
			name: "valid_title",
			app: &AppVersion{
				Title: "Test Application",
			},
			expected: "Test Application",
		},
		{
			name:     "empty_title",
			app:      &AppVersion{},
			expected: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, tc.app.GetTitle(), "GetTitle result mismatch")
		})
	}
}

func TestAppVersion_GetLongTitle(t *testing.T) {
	tests := []struct {
		name     string
		app      *AppVersion
		expected string
	}{
		{
			name: "full_title",
			app: &AppVersion{
				Title:          "Test Application",
				Version:        semver.MustParse("1.0.0"),
				BuildType:      BUILDTYPE_DEBUG,
				DeploymentType: DEPLOYMENTTYPE_DEV,
			},
			expected: "Test Application v1.0.0, debug-dev",
		},
		{

			name: "no_build_deployment",
			app: &AppVersion{
				Title:   "Test Application",
				Version: semver.MustParse("1.0.0"),
			},
			expected: "Test Application v1.0.0",
		},
		{
			name: "no_version",
			app: &AppVersion{
				Title:          "Test Application",
				BuildType:      BUILDTYPE_DEBUG,
				DeploymentType: DEPLOYMENTTYPE_DEV,
			},
			expected: "Test Application, debug-dev",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, tc.app.GetLongTitle(), "GetLongTitle result mismatch")
		})
	}
}

func TestAppVersion_SetGetBuildName(t *testing.T) {
	tests := []struct {
		name     string
		setName  BuildName
		expected BuildName
	}{
		{
			name:     "set_valid",
			setName:  "testapp",
			expected: "testapp",
		},
		{
			name:     "set_empty",
			setName:  "",
			expected: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ResetAppVersion() // Reset global state
			app := &AppVersion{}
			app.SetBuildName(tc.setName)
			assert.Equal(t, tc.expected, app.GetBuildName(), "GetBuildName result mismatch")
			assert.Equal(t, tc.expected, BUILDNAME(), "Global BuildName mismatch")
		})
	}
}

func TestAppVersion_SetGetBuildType(t *testing.T) {
	tests := []struct {
		name     string
		setType  BuildType
		expected BuildType
		wantErr  bool
		errMsg   string
	}{
		{
			name:     "set_valid",
			setType:  BUILDTYPE_DEBUG,
			expected: BUILDTYPE_DEBUG,
		},
		{
			name:    "set_invalid",
			setType: "invalid@type",
			wantErr: true,
			errMsg:  "invalid build type",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ResetAppVersion() // Reset global state
			app := &AppVersion{}
			err := app.SetBuildType(tc.setType)
			if tc.wantErr {
				assert.ErrorContains(t, err, tc.errMsg, "Expected error containing %q", tc.errMsg)
				return
			}
			assert.NoError(t, err, "Unexpected error")
			assert.Equal(t, tc.expected, app.GetBuildType(), "GetBuildType result mismatch")
			assert.Equal(t, tc.expected, BUILDTYPE(), "Global BuildType mismatch")
		})
	}
}

func TestAppVersion_SetGetDeploymentType(t *testing.T) {
	tests := []struct {
		name     string
		setType  DeploymentType
		expected DeploymentType
		wantErr  bool
		errMsg   string
	}{
		{
			name:     "set_valid",
			setType:  DEPLOYMENTTYPE_DEV,
			expected: DEPLOYMENTTYPE_DEV,
		},
		{
			name:    "set_invalid",
			setType: "invalid@type",
			wantErr: true,
			errMsg:  "invalid deployment type",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ResetAppVersion() // Reset global state
			app := &AppVersion{}
			err := app.SetDeploymentType(tc.setType)
			if tc.wantErr {
				assert.ErrorContains(t, err, tc.errMsg, "Expected error containing %q", tc.errMsg)
				return
			}
			assert.NoError(t, err, "Unexpected error")
			assert.Equal(t, tc.expected, app.GetDeploymentType(), "GetDeploymentType result mismatch")
			assert.Equal(t, tc.expected, DEPLOYMENTTYPE(), "Global DeploymentType mismatch")
		})
	}
}

func TestLoadAppVersionFromBytes(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected *AppVersion
		wantErr  bool
		errMsg   string
	}{
		{
			name: "valid_json",
			input: `{
				"name": "TestApp",
				"version": "1.0.0",
				"about": "A test application",
				"owner": "TestOwner",
				"legalMark": "TestLegalMark",
				"buildName": "testapp",
				"buildType": "debug",
				"deploymentType": "dev"
			}`,
			expected: &AppVersion{
				Name:           "TestApp",
				Version:        semver.MustParse("1.0.0"),
				About:          "A test application",
				Owner:          "TestOwner",
				LegalMark:      "TestLegalMark",
				BuildName:      "testapp",
				BuildType:      BUILDTYPE_DEBUG,
				DeploymentType: DEPLOYMENTTYPE_DEV,
			},
		},
		{
			name:    "empty_input",
			input:   "",
			wantErr: true,
			errMsg:  "input bytes cannot be empty",
		},
		{
			name:    "invalid_json",
			input:   `{"name":"TestApp"`,
			wantErr: true,
			errMsg:  "failed to unmarshal AppVersion",
		},
		{
			name: "invalid_build_type",
			input: `{
				"name": "TestApp",
				"version": "1.0.0",
				"about": "A test application",
				"owner": "TestOwner",
				"legalMark": "TestLegalMark",
				"buildType": "invalid@type"
			}`,
			wantErr: true,
			errMsg:  "invalid build type",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ResetAppVersion() // Reset global state
			result, err := LoadAppVersionFromBytes([]byte(tc.input))
			if tc.wantErr {
				assert.ErrorContains(t, err, tc.errMsg, "Expected error containing %q", tc.errMsg)
				return
			}
			assert.NoError(t, err, "Unexpected error")
			assert.Equal(t, tc.expected, result, "Loaded AppVersion mismatch")
		})
	}
}

func TestLoadEmbeddedAppVersion(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected *AppVersion
		wantErr  bool
		errMsg   string
	}{
		{
			name: "valid_file",
			path: "test_data/app_version.json",
			expected: &AppVersion{
				Name:           "TestApp",
				Version:        semver.MustParse("1.0.0"),
				About:          "A test application",
				Owner:          "TestOwner",
				LegalMark:      "TestLegalMark",
				BuildName:      "testapp",
				BuildType:      BUILDTYPE_DEBUG,
				DeploymentType: DEPLOYMENTTYPE_DEV,
			},
		},
		{
			name:    "non_existent_file",
			path:    "test_data/nonexistent.json",
			wantErr: true,
			errMsg:  "failed to read embedded file",
		},
		{
			name:    "invalid_json",
			path:    "test_data/malformed.json",
			wantErr: true,
			errMsg:  "failed to unmarshal AppVersion",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ResetAppVersion() // Reset global state
			result, err := LoadEmbeddedAppVersion(testFS, tc.path)
			if tc.wantErr {
				assert.ErrorContains(t, err, tc.errMsg, "Expected error containing %q", tc.errMsg)
				return
			}
			assert.NoError(t, err, "Unexpected error")
			assert.Equal(t, tc.expected, result, "Loaded AppVersion mismatch")
		})
	}
}

func TestGetSetAppVersion(t *testing.T) {
	tests := []struct {
		name    string
		setApp  *AppVersion
		wantErr bool
		errMsg  string
	}{
		{
			name: "set_valid",
			setApp: &AppVersion{
				Name:           "TestApp",
				Version:        semver.MustParse("1.0.0"),
				Owner:          "TestOwner",
				LegalMark:      "TestLegalMark",
				BuildName:      "testapp",
				BuildType:      BUILDTYPE_DEBUG,
				DeploymentType: DEPLOYMENTTYPE_DEV,
			},
		},
		{
			name:    "set_nil",
			setApp:  nil,
			wantErr: true,
			errMsg:  "AppVersion cannot be nil",
		},
		{
			name: "set_twice",
			setApp: &AppVersion{
				Name:      "TestApp",
				Version:   semver.MustParse("1.0.0"),
				Owner:     "TestOwner",
				LegalMark: "TestLegalMark",
			},
			wantErr: true,
			errMsg:  "AppVersion already initialized",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ResetAppVersion() // Reset global state
			if tc.name == "set_twice" {
				_ = SetAppVersion(&AppVersion{
					Name:      "TestApp",
					Version:   semver.MustParse("1.0.0"),
					Owner:     "TestOwner",
					LegalMark: "TestLegalMark",
				})
			}

			err := SetAppVersion(tc.setApp)
			if tc.wantErr {
				assert.ErrorContains(t, err, tc.errMsg, "Expected error containing %q", tc.errMsg)
				return
			}
			assert.NoError(t, err, "Unexpected error")
			result, err := GetAppVersion()
			assert.NoError(t, err, "Unexpected error in GetAppVersion")
			assert.Equal(t, tc.setApp, result, "GetAppVersion result mismatch")
		})
	}
}

func TestGetAppVersion_Uninitialized(t *testing.T) {
	ResetAppVersion() // Ensure clean state
	_, err := GetAppVersion()
	assert.ErrorContains(t, err, "AppVersion not initialized", "Expected error for uninitialized AppVersion")
}
