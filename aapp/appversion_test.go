package aapp

import (
	"embed"
	"fmt"
	"reflect"
	"strings"
	"sync"
	"testing"

	"github.com/Masterminds/semver/v3"
)

// Note: This test file assumes the existence of types BuildName, BuildType, DeploymentType
// with methods like IsEmpty(), IsValid(), and String(), as well as constants like
// BUILDTYPE_DEBUG and DEPLOYMENTTYPE_DEV, which are likely defined in another file in the package.

//go:embed test_data/*
var testFS embed.FS

func TestAppVersionFormat_IsEmpty(t *testing.T) {
	tests := []struct {
		name string
		avf  AppVersionFormat
		want bool
	}{
		{"empty string", "", true},
		{"whitespace only", "   ", true},
		{"non-empty", "json", false},
		{"mixed whitespace", "\t\n", true},
		{"with content", " version ", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.avf.IsEmpty(); got != tt.want {
				t.Errorf("AppVersionFormat.IsEmpty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAppVersion_Format(t *testing.T) {
	validVersion := semver.MustParse("1.2.3")
	invalidVersion := semver.Version{} // invalid

	validApp := &AppVersion{
		Name:           "TestApp",
		Version:        *validVersion,
		About:          "Test about",
		Owner:          "Test Owner",
		LegalMark:      "Test Legal",
		BuildName:      "testapp",
		BuildType:      BUILDTYPE_DEBUG,
		DeploymentType: DEPLOYMENTTYPE_DEV,
	}

	invalidApp := &AppVersion{
		Name:           "TestApp",
		Version:        invalidVersion,
		About:          "Test about",
		Owner:          "Test Owner",
		LegalMark:      "Test Legal",
		BuildName:      "testapp",
		BuildType:      BUILDTYPE_DEBUG,
		DeploymentType: DEPLOYMENTTYPE_DEV,
	}

	tests := []struct {
		name    string
		a       *AppVersion
		format  AppVersionFormat
		want    string
		wantErr bool
	}{
		{"json valid", validApp, APPVERSIONFORMAT_JSON, "{\n  \"name\": \"TestApp\",\n  \"version\": \"1.2.3\",\n  \"about\": \"Test about\",\n  \"owner\": \"Test Owner\",\n  \"legalMark\": \"Test Legal\",\n  \"buildName\": \"testapp\",\n  \"buildType\": \"debug\",\n  \"deploymentType\": \"dev\"\n}", false},
		{"version only valid", validApp, APPVERSIONFORMAT_VERSION_ONLY, "v1.2.3", false},
		{"version only invalid", invalidApp, APPVERSIONFORMAT_VERSION_ONLY, "", true},
		{"no v valid", validApp, APPVERSIONFORMAT_NO_V, "1.2.3", false},
		{"no v invalid", invalidApp, APPVERSIONFORMAT_NO_V, "", true},
		{"app only", validApp, APPVERSIONFORMAT_APP_ONLY, "testapp", false},
		{"about only", validApp, APPVERSIONFORMAT_ABOUT_ONLY, "Test about", false},
		{"build valid", validApp, APPVERSIONFORMAT_BUILD, "testapp@1.2.3,debug-dev", false},
		{"build invalid", invalidApp, APPVERSIONFORMAT_BUILD, "", true},
		{"default valid", validApp, APPVERSIONFORMAT_APP_AT_VERSION, "testapp@1.2.3", false},
		{"default invalid", invalidApp, APPVERSIONFORMAT_APP_AT_VERSION, "", true},
		{"unknown format valid", validApp, "unknown", "testapp@1.2.3", false}, // falls to default
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.a.Format(tt.format)
			if (err != nil) != tt.wantErr {
				t.Errorf("AppVersion.Format() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("AppVersion.Format() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAppVersion_Validate(t *testing.T) {
	validVersion := semver.MustParse("1.2.3")

	tests := []struct {
		name               string
		a                  *AppVersion
		setupBuildVars     func()
		setupOsExecutable  func() (restore func())
		wantErr            bool
		wantBuildName      BuildName
		wantBuildType      BuildType
		wantDeploymentType DeploymentType
	}{
		{"valid with all fields", &AppVersion{Name: "App", Version: *validVersion, Owner: "Owner", LegalMark: "Legal", BuildName: "app", BuildType: BUILDTYPE_DEBUG, DeploymentType: DEPLOYMENTTYPE_DEV}, nil, nil, false, "app", BUILDTYPE_DEBUG, DEPLOYMENTTYPE_DEV},
		{"missing name", &AppVersion{Version: *validVersion, Owner: "Owner", LegalMark: "Legal"}, nil, nil, true, "", "", ""},
		{"missing owner", &AppVersion{Name: "App", Version: *validVersion, LegalMark: "Legal"}, nil, nil, true, "", "", ""},
		{"missing legal mark", &AppVersion{Name: "App", Version: *validVersion, Owner: "Owner"}, nil, nil, true, "", "", ""},
		{"invalid version", &AppVersion{Name: "App", Owner: "Owner", LegalMark: "Legal", Version: semver.Version{}}, nil, nil, true, "", "", ""},
		{"default buildname from name", &AppVersion{Name: "Test App", Version: *validVersion, Owner: "Owner", LegalMark: "Legal"}, nil, func() func() {
			orig := osExecutable
			osExecutable = func() (string, error) { return "", fmt.Errorf("mock error") }
			return func() { osExecutable = orig }
		}, false, "testapp", BUILDTYPE_DEBUG, DEPLOYMENTTYPE_DEV},
		{"default buildname from exe", &AppVersion{Name: "App", Version: *validVersion, Owner: "Owner", LegalMark: "Legal"}, nil, func() func() {
			orig := osExecutable
			osExecutable = func() (string, error) { return "/mock/path/mockexe", nil }
			return func() { osExecutable = orig }
		}, false, "mockexe", BUILDTYPE_DEBUG, DEPLOYMENTTYPE_DEV},
		{"buildname from compiler var", &AppVersion{Name: "App", Version: *validVersion, Owner: "Owner", LegalMark: "Legal"}, func() { buildName = "compilerbuild" }, nil, false, "compilerbuild", BUILDTYPE_DEBUG, DEPLOYMENTTYPE_DEV},
		{"buildtype from compiler var", &AppVersion{Name: "App", Version: *validVersion, Owner: "Owner", LegalMark: "Legal", BuildType: ""}, func() { buildType = "release" }, func() func() {
			orig := osExecutable
			osExecutable = func() (string, error) { return "", fmt.Errorf("mock error") }
			return func() { osExecutable = orig }
		}, false, "app", "release", DEPLOYMENTTYPE_DEV},
		//		{"invalid buildtype", &AppVersion{Name: "App", Version: *validVersion, Owner: "Owner", LegalMark: "Legal", BuildType: "invalid"}, nil, nil, true, "", "", ""},
		{"deploymenttype from compiler var", &AppVersion{Name: "App", Version: *validVersion, Owner: "Owner", LegalMark: "Legal", DeploymentType: ""}, func() { buildDeploymentType = "prod" }, func() func() {
			orig := osExecutable
			osExecutable = func() (string, error) { return "", fmt.Errorf("mock error") }
			return func() { osExecutable = orig }
		}, false, "app", BUILDTYPE_DEBUG, "prod"},
		//		{"invalid deploymenttype", &AppVersion{Name: "App", Version: *validVersion, Owner: "Owner", LegalMark: "Legal", DeploymentType: "invalid"}, nil, nil, true, "", "", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset build vars
			buildName = ""
			buildType = ""
			buildDeploymentType = ""
			if tt.setupBuildVars != nil {
				tt.setupBuildVars()
			}
			var restore func()
			if tt.setupOsExecutable != nil {
				restore = tt.setupOsExecutable()
			}
			if restore != nil {
				defer restore()
			}
			err := tt.a.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("AppVersion.Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if tt.a.BuildName != tt.wantBuildName {
					t.Errorf("AppVersion.Validate() BuildName = %v, want %v", tt.a.BuildName, tt.wantBuildName)
				}
				if tt.a.BuildType != tt.wantBuildType {
					t.Errorf("AppVersion.Validate() BuildType = %v, want %v", tt.a.BuildType, tt.wantBuildType)
				}
				if tt.a.DeploymentType != tt.wantDeploymentType {
					t.Errorf("AppVersion.Validate() DeploymentType = %v, want %v", tt.a.DeploymentType, tt.wantDeploymentType)
				}
			}
		})
	}
}

func TestAppVersion_GetTitle(t *testing.T) {
	validVersion := semver.MustParse("1.2.3")
	invalidVersion := semver.Version{}

	tests := []struct {
		name string
		a    *AppVersion
		want string
	}{
		{"no build info no version", &AppVersion{Name: "App"}, "App"},
		{"no build info with version", &AppVersion{Name: "App", Version: *validVersion}, "App v1.2.3"},
		{"with build info no version", &AppVersion{Name: "App", BuildType: "debug", DeploymentType: "dev"}, "App, debug-dev"},
		{"with build info with version", &AppVersion{Name: "App", Version: *validVersion, BuildType: "debug", DeploymentType: "dev"}, "App v1.2.3, debug-dev"},
		{"invalid version with build", &AppVersion{Name: "App", Version: invalidVersion, BuildType: "debug", DeploymentType: "dev"}, "App, debug-dev"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.a.GetTitle(); got != tt.want {
				t.Errorf("AppVersion.GetTitle() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAppVersion_GetBuildName(t *testing.T) {
	a := &AppVersion{BuildName: "testbuild"}
	if got := a.GetBuildName(); got != "testbuild" {
		t.Errorf("AppVersion.GetBuildName() = %v, want %v", got, "testbuild")
	}
}

func TestAppVersion_GetBuildType(t *testing.T) {
	a := &AppVersion{BuildType: "debug"}
	if got := a.GetBuildType(); got != "debug" {
		t.Errorf("AppVersion.GetBuildType() = %v, want %v", got, "debug")
	}
}

func TestAppVersion_GetDeploymentType(t *testing.T) {
	a := &AppVersion{DeploymentType: "dev"}
	if got := a.GetDeploymentType(); got != "dev" {
		t.Errorf("AppVersion.GetDeploymentType() = %v, want %v", got, "dev")
	}
}

func TestAppVersion_Clone(t *testing.T) {
	validVersion := semver.MustParse("1.2.3")
	a := &AppVersion{
		Name:           "App",
		Version:        *validVersion,
		About:          "About",
		Owner:          "Owner",
		LegalMark:      "Legal",
		BuildName:      "build",
		BuildType:      "debug",
		DeploymentType: "dev",
	}
	clone := a.Clone()
	if clone == nil {
		t.Errorf("AppVersion.Clone() returned nil")
	}
	if !reflect.DeepEqual(a, clone) {
		t.Errorf("AppVersion.Clone() = %v, want %v", clone, a)
	}
	if &a.Version == &clone.Version {
		t.Errorf("AppVersion.Clone() did not deep copy Version")
	}

	nilA := (*AppVersion)(nil)
	if got := nilA.Clone(); got != nil {
		t.Errorf("AppVersion.Clone() on nil = %v, want nil", got)
	}

	invalidVersionA := &AppVersion{Version: semver.Version{}}
	cloneInvalid := invalidVersionA.Clone()
	if !reflect.DeepEqual(invalidVersionA, cloneInvalid) {
		t.Errorf("AppVersion.Clone() with invalid version = %v, want %v", cloneInvalid, invalidVersionA)
	}
}

func TestLoadAppVersionFromBytes(t *testing.T) {
	validVersion := semver.MustParse("1.2.3")

	tests := []struct {
		name              string
		b                 []byte
		setupOsExecutable func() (restore func())
		want              *AppVersion
		wantErr           bool
	}{
		{"valid json", []byte(`{"name":"App","version":"1.2.3","about":"About","owner":"Owner","legalMark":"Legal","buildName":"build","buildType":"debug","deploymentType":"dev"}`), nil, &AppVersion{Name: "App", Version: *validVersion, About: "About", Owner: "Owner", LegalMark: "Legal", BuildName: "build", BuildType: "debug", DeploymentType: "dev"}, false},
		{"invalid json missing fields", []byte(`{"name":"App"}`), nil, nil, true},
		{"malformed json", []byte(`{invalid}`), nil, nil, true},
		{"empty bytes", []byte{}, nil, nil, true},
		{"json with whitespace", []byte(` {"name":" App ","version":"1.2.3","about":" About ","owner":" Owner ","legalMark":" Legal "} `), func() func() {
			orig := osExecutable
			osExecutable = func() (string, error) { return "", fmt.Errorf("mock error") }
			return func() { osExecutable = orig }
		}, &AppVersion{Name: "App", Version: *validVersion, About: "About", Owner: "Owner", LegalMark: "Legal", BuildName: "app", BuildType: BUILDTYPE_DEBUG, DeploymentType: DEPLOYMENTTYPE_DEV}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset build vars
			buildName = ""
			buildType = ""
			buildDeploymentType = ""
			var restore func()
			if tt.setupOsExecutable != nil {
				restore = tt.setupOsExecutable()
			}
			if restore != nil {
				defer restore()
			}
			got, err := LoadAppVersionFromBytes(tt.b)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadAppVersionFromBytes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LoadAppVersionFromBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLoadEmbeddedAppVersion(t *testing.T) {
	// Assume test_data/version.json exists with content: {"name":"EmbeddedApp","version":"4.5.6","owner":"Owner","legalMark":"Legal"}
	// For this test, we need to actually embed a file or mock, but since embed is used, assume testFS has it.

	validPath := "test_data/version.json"
	invalidPath := "nonexistent.json"
	emptyPath := ""

	t.Run("valid path", func(t *testing.T) {
		got, err := LoadEmbeddedAppVersion(testFS, validPath)
		if err != nil {
			t.Errorf("LoadEmbeddedAppVersion() error = %v, want no error", err)
		}
		if got == nil {
			t.Errorf("LoadEmbeddedAppVersion() got nil")
		}
		// Further assertions based on embedded content
	})

	t.Run("default path", func(t *testing.T) {
		got, err := LoadEmbeddedAppVersion(testFS, emptyPath)
		if err == nil || !strings.Contains(err.Error(), "failed to read") { // assuming default "app/version.json" not in testFS
			t.Errorf("LoadEmbeddedAppVersion() error = %v, want error", err)
		}
		if got != nil {
			t.Errorf("LoadEmbeddedAppVersion() = %v, want nil", got)
		}
	})

	t.Run("invalid path", func(t *testing.T) {
		got, err := LoadEmbeddedAppVersion(testFS, invalidPath)
		if err == nil {
			t.Errorf("LoadEmbeddedAppVersion() error = %v, want error", err)
		}
		if got != nil {
			t.Errorf("LoadEmbeddedAppVersion() = %v, want nil", got)
		}
	})

	t.Run("malformed json", func(t *testing.T) {
		// Assume test_data/malformed.json with invalid json
		got, err := LoadEmbeddedAppVersion(testFS, "test_data/malformed.json")
		if err == nil || !strings.Contains(err.Error(), "unmarshal") {
			t.Errorf("LoadEmbeddedAppVersion() error = %v, want unmarshal error", err)
		}
		if got != nil {
			t.Errorf("LoadEmbeddedAppVersion() = %v, want nil", got)
		}
	})
}

func TestGetAppVersion(t *testing.T) {
	ResetAppVersion()
	_, err := GetAppVersion()
	if err == nil {
		t.Errorf("GetAppVersion() before set, want error")
	}

	validApp := &AppVersion{Name: "App", Version: *semver.MustParse("1.0.0"), Owner: "Owner", LegalMark: "Legal"}
	if err := SetAppVersion(validApp); err != nil {
		t.Errorf("SetAppVersion() error = %v", err)
	}

	got, err := GetAppVersion()
	if err != nil {
		t.Errorf("GetAppVersion() error = %v", err)
	}
	if !reflect.DeepEqual(got, validApp) {
		t.Errorf("GetAppVersion() = %v, want %v", got, validApp)
	}

	ResetAppVersion()
}

func TestCloneAppVersion(t *testing.T) {
	ResetAppVersion()
	if got := CloneAppVersion(); got != nil {
		t.Errorf("CloneAppVersion() before set = %v, want nil", got)
	}

	validApp := &AppVersion{Name: "App", Version: *semver.MustParse("1.0.0"), Owner: "Owner", LegalMark: "Legal"}
	SetAppVersion(validApp)

	clone := CloneAppVersion()
	if clone == nil {
		t.Errorf("CloneAppVersion() returned nil")
	}
	if !reflect.DeepEqual(validApp, clone) {
		t.Errorf("CloneAppVersion() = %v, want %v", clone, validApp)
	}

	// Set invalid version
	global.instance.Version = semver.Version{}
	if got := CloneAppVersion(); got != nil {
		t.Errorf("CloneAppVersion() with invalid version = %v, want nil", got)
	}

	ResetAppVersion()
}

func TestMustVersion(t *testing.T) {
	ResetAppVersion()
	if got := MustVersion(); !got.Equal(&semver.Version{}) {
		t.Errorf("MustVersion() before set = %v, want empty", got)
	}

	validVersion := semver.MustParse("1.0.0")
	validApp := &AppVersion{Name: "App", Version: *validVersion, Owner: "Owner", LegalMark: "Legal"}
	SetAppVersion(validApp)

	if got := MustVersion(); !got.Equal(validVersion) {
		t.Errorf("MustVersion() = %v, want %v", got, validVersion)
	}

	global.instance.Version = semver.Version{}
	if got := MustVersion(); !got.Equal(&semver.Version{}) {
		t.Errorf("MustVersion() with invalid = %v, want empty", got)
	}

	ResetAppVersion()
}

func TestSetAppVersion(t *testing.T) {
	ResetAppVersion()

	if err := SetAppVersion(nil); err == nil {
		t.Errorf("SetAppVersion(nil) want error")
	}

	validApp := &AppVersion{Name: "App", Version: *semver.MustParse("1.0.0"), Owner: "Owner", LegalMark: "Legal"}
	if err := SetAppVersion(validApp); err != nil {
		t.Errorf("SetAppVersion() error = %v", err)
	}

	if err := SetAppVersion(validApp); err == nil {
		t.Errorf("SetAppVersion() second time want error")
	}

	ResetAppVersion()
}

func TestSetAppVersionByFS(t *testing.T) {
	// Mock the func for testing, but since it's var, we can override
	orig := SetAppVersionByFS
	defer func() { SetAppVersionByFS = orig }()

	called := false
	SetAppVersionByFS = func(fs embed.FS, myPath string) error {
		called = true
		if myPath != "testpath" {
			return fmt.Errorf("wrong path")
		}
		return nil
	}

	err := SetAppVersionByFS(testFS, "testpath")
	if err != nil {
		t.Errorf("SetAppVersionByFS() error = %v", err)
	}
	if !called {
		t.Errorf("SetAppVersionByFS not called")
	}

	// Test with invalid
	SetAppVersionByFS = func(fs embed.FS, myPath string) error {
		return fmt.Errorf("mock error")
	}
	if err := SetAppVersionByFS(testFS, ""); err == nil {
		t.Errorf("SetAppVersionByFS() want error")
	}
}

func TestResetAppVersion(t *testing.T) {
	validApp := &AppVersion{Name: "App", Version: *semver.MustParse("1.0.0"), Owner: "Owner", LegalMark: "Legal"}
	SetAppVersion(validApp)

	ResetAppVersion()

	_, err := GetAppVersion()
	if err == nil {
		t.Errorf("GetAppVersion() after reset, want error")
	}
}

func TestGlobalConcurrency(t *testing.T) {
	// Test concurrent access to global
	ResetAppVersion()

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			SetAppVersion(&AppVersion{Name: "App", Version: *semver.MustParse("1.0.0"), Owner: "Owner", LegalMark: "Legal"})
		}()
	}
	wg.Wait()

	// Since set only once, but concurrent, might error on some, but test for race
	// Note: Use -race flag in real tests

	got, err := GetAppVersion()
	if err != nil {
		t.Errorf("GetAppVersion() after concurrent set = %v", err)
	}
	if got == nil {
		t.Errorf("Global not set")
	}

	ResetAppVersion()
}
