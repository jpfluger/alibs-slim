package aapp

import (
	"embed"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/Masterminds/semver/v3"
	"github.com/jpfluger/alibs-slim/autils"
)

// AppVersionFormat defines the format in which the version information can be presented.
type AppVersionFormat string

// Constants for different version formats.
const (
	APPVERSIONFORMAT_JSON           AppVersionFormat = "json"
	APPVERSIONFORMAT_VERSION_ONLY   AppVersionFormat = "version"
	APPVERSIONFORMAT_NO_V           AppVersionFormat = "version-no-v"
	APPVERSIONFORMAT_APP_ONLY       AppVersionFormat = "name"
	APPVERSIONFORMAT_ABOUT_ONLY     AppVersionFormat = "about"
	APPVERSIONFORMAT_APP_AT_VERSION AppVersionFormat = "app@version" // default
	APPVERSIONFORMAT_BUILD          AppVersionFormat = "build"
)

// IsEmpty checks if the AppVersionFormat is empty after trimming whitespace.
// Returns true if the string is empty or contains only whitespace.
func (avf AppVersionFormat) IsEmpty() bool {
	return strings.TrimSpace(string(avf)) == ""
}

// AppVersion represents the version information of an application.
type AppVersion struct {
	Name           string         `json:"name,omitempty"`           // The name of the application.
	Version        semver.Version `json:"version,omitempty"`        // The semantic version of the application.
	About          string         `json:"about,omitempty"`          // Additional information about the application.
	Owner          string         `json:"owner,omitempty"`          // The owner of the application.
	LegalMark      string         `json:"legalMark,omitempty"`      // Legal trademark or copyright notice.
	BuildName      BuildName      `json:"buildName,omitempty"`      // The name of the built binary, defaults to trimmed/lowercase Name if empty.
	BuildType      BuildType      `json:"buildType,omitempty"`      // The build type, defaults to BUILDTYPE_DEBUG if empty.
	DeploymentType DeploymentType `json:"deploymentType,omitempty"` // The deployment type, defaults to DEPLOYMENTTYPE_DEV if empty.
}

// globalAppVersion holds the global AppVersion instance and related state.
type globalAppVersion struct {
	mu       sync.RWMutex
	instance *AppVersion
}

var global = &globalAppVersion{}

var (
	// Plain string vars for build-time injection via -X
	buildNameStr           string = ""
	buildTypeStr           string = ""
	buildDeploymentTypeStr string = ""

	// Your typed vars (keep these for type safety in the code)
	buildName           BuildName      = ""
	buildType           BuildType      = ""
	buildDeploymentType DeploymentType = ""
)

func init() {
	// Assign the injected strings to the typed vars at runtime
	buildName = BuildName(buildNameStr)
	buildType = BuildType(buildTypeStr)
	buildDeploymentType = DeploymentType(buildDeploymentTypeStr)
}

// osExecutable is a variable for os.Executable to allow mocking in tests.
var osExecutable = os.Executable

// BUILDNAME returns the buildName BuildName value, typically set by the compiler.
func BUILDNAME() BuildName {
	return buildName
}

// BUILDTYPE returns the buildType BuildType value, typically set by the compiler.
func BUILDTYPE() BuildType {
	return buildType
}

// BUILDDEPLOYMENTTYPE returns the buildDeploymentType DeploymentType value, typically set by the compiler.
func BUILDDEPLOYMENTTYPE() DeploymentType {
	return buildDeploymentType
}

// Format returns the version information in the specified format.
// Returns an error if marshaling fails for JSON format or if version is nil for certain formats.
func (a *AppVersion) Format(format AppVersionFormat) (string, error) {
	isVersionValid := autils.IsSemverValid(&a.Version)
	switch format {
	case APPVERSIONFORMAT_JSON:
		data, err := json.MarshalIndent(a, "", "  ")
		if err != nil {
			return "", fmt.Errorf("failed to marshal JSON for AppVersion: %w", err)
		}
		return string(data), nil
	case APPVERSIONFORMAT_VERSION_ONLY:
		if !isVersionValid {
			return "", fmt.Errorf("version cannot be nil")
		}
		return fmt.Sprintf("v%s", a.Version.String()), nil
	case APPVERSIONFORMAT_NO_V:
		if !isVersionValid {
			return "", fmt.Errorf("version cannot be nil")
		}
		return a.Version.String(), nil
	case APPVERSIONFORMAT_APP_ONLY:
		return a.BuildName.String(), nil
	case APPVERSIONFORMAT_ABOUT_ONLY:
		return a.About, nil
	case APPVERSIONFORMAT_BUILD:
		if !isVersionValid {
			return "", fmt.Errorf("version cannot be nil")
		}
		return fmt.Sprintf("%s@%s,%s-%s", a.BuildName.String(), a.Version.String(), a.BuildType.String(), a.DeploymentType.String()), nil
	default: // APPVERSIONFORMAT_APP_AT_VERSION
		if !isVersionValid {
			return "", fmt.Errorf("version cannot be nil")
		}
		return fmt.Sprintf("%s@%s", a.BuildName.String(), a.Version.String()), nil
	}
}

// Validate checks if the AppVersion fields are valid and sets defaults if necessary.
// Returns an error if required fields are invalid.
func (a *AppVersion) Validate() error {
	if a.Name == "" {
		return fmt.Errorf("name cannot be empty")
	}
	if a.Owner == "" {
		return fmt.Errorf("owner cannot be empty")
	}
	if a.LegalMark == "" {
		return fmt.Errorf("legal mark cannot be empty")
	}
	if !autils.IsSemverValid(&a.Version) {
		return fmt.Errorf("version cannot be nil or empty")
	}
	if a.BuildName.IsEmpty() {
		if BUILDNAME().IsEmpty() {
			if exePath, err := osExecutable(); err == nil {
				a.BuildName = BuildName(filepath.Base(exePath))
			}
			if a.BuildName.IsEmpty() {
				trimmedLower := autils.ToStringTrimLower(a.Name)
				noSpaces := strings.ReplaceAll(trimmedLower, " ", "")
				a.BuildName = BuildName(noSpaces)
			}
		} else {
			a.BuildName = BUILDNAME()
		}
	}
	if a.BuildType.IsEmpty() {
		if BUILDTYPE().IsEmpty() {
			a.BuildType = BUILDTYPE_DEBUG
		} else {
			a.BuildType = BUILDTYPE()
		}
	} else if !a.BuildType.IsValid() {
		return fmt.Errorf("invalid build type %q", a.BuildType)
	}
	if a.DeploymentType.IsEmpty() {
		if BUILDDEPLOYMENTTYPE().IsEmpty() {
			a.DeploymentType = DEPLOYMENTTYPE_DEV
		} else {
			a.DeploymentType = BUILDDEPLOYMENTTYPE()
		}
	} else if !a.DeploymentType.IsValid() {
		return fmt.Errorf("invalid deployment type %q", a.DeploymentType)
	}
	return nil
}

// GetTitle returns the title of the application.
func (a *AppVersion) GetTitle() string {
	if a.BuildType.IsEmpty() && a.DeploymentType.IsEmpty() {
		if !autils.IsSemverValid(&a.Version) {
			return a.Name
		}
		return fmt.Sprintf("%s v%s", a.Name, a.Version.String())
	}
	if !autils.IsSemverValid(&a.Version) {
		return fmt.Sprintf("%s, %s-%s", a.Name, a.BuildType.String(), a.DeploymentType.String())
	}
	return fmt.Sprintf("%s v%s, %s-%s", a.Name, a.Version.String(), a.BuildType.String(), a.DeploymentType.String())
}

// GetBuildName returns the build name of the application.
func (a *AppVersion) GetBuildName() BuildName {
	return a.BuildName
}

// GetBuildType returns the build type of the application.
func (a *AppVersion) GetBuildType() BuildType {
	return a.BuildType
}

// GetDeploymentType returns the deployment type of the application.
func (a *AppVersion) GetDeploymentType() DeploymentType {
	return a.DeploymentType
}

// Clone returns a deep copy of the AppVersion instance.
// If the original Version is nil, the clone's Version will also be nil.
func (a *AppVersion) Clone() *AppVersion {
	if a == nil {
		return nil
	}
	clone := &AppVersion{
		Name:           a.Name,
		Version:        a.Version,
		About:          a.About,
		Owner:          a.Owner,
		LegalMark:      a.LegalMark,
		BuildName:      a.BuildName,
		BuildType:      a.BuildType,
		DeploymentType: a.DeploymentType,
	}
	return clone
}

// LoadAppVersionFromBytes loads an AppVersion from a JSON byte slice.
// Returns an error if the input is empty or invalid.
func LoadAppVersionFromBytes(b []byte) (*AppVersion, error) {
	if len(b) == 0 {
		return nil, fmt.Errorf("input bytes cannot be empty")
	}
	var appVersion AppVersion
	if err := json.Unmarshal(b, &appVersion); err != nil {
		return nil, fmt.Errorf("failed to unmarshal AppVersion: %w", err)
	}
	// Trim whitespace from string fields
	appVersion.Name = strings.TrimSpace(appVersion.Name)
	appVersion.About = strings.TrimSpace(appVersion.About)
	appVersion.Owner = strings.TrimSpace(appVersion.Owner)
	appVersion.LegalMark = strings.TrimSpace(appVersion.LegalMark)
	if err := appVersion.Validate(); err != nil {
		return nil, fmt.Errorf("AppVersion validation failed: %w", err)
	}
	return &appVersion, nil
}

// LoadEmbeddedAppVersion loads an AppVersion from an embedded JSON file.
// Defaults to "app/version.json" if path is empty.
// Returns an error if the file cannot be read or is invalid.
func LoadEmbeddedAppVersion(fs embed.FS, myPath string) (*AppVersion, error) {
	myPath = strings.TrimSpace(myPath)
	if myPath == "" {
		myPath = "app/version.json"
	}
	b, err := fs.ReadFile(myPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read embedded file %q: %w", myPath, err)
	}
	var appVersion AppVersion
	if err := json.Unmarshal(b, &appVersion); err != nil {
		return nil, fmt.Errorf("failed to unmarshal AppVersion: %w", err)
	}
	// Trim whitespace from string fields
	appVersion.Name = strings.TrimSpace(appVersion.Name)
	appVersion.About = strings.TrimSpace(appVersion.About)
	appVersion.Owner = strings.TrimSpace(appVersion.Owner)
	appVersion.LegalMark = strings.TrimSpace(appVersion.LegalMark)
	if err := appVersion.Validate(); err != nil {
		return nil, fmt.Errorf("AppVersion validation failed: %w", err)
	}
	return &appVersion, nil
}

// GetAppVersion returns the global AppVersion instance.
// Returns an error if not initialized.
var GetAppVersion = func() (*AppVersion, error) {
	global.mu.RLock()
	defer global.mu.RUnlock()
	if global.instance == nil {
		return nil, fmt.Errorf("AppVersion not initialized")
	}
	return global.instance, nil
}

// APPVERSION returns the global AppVersion instance.
// Panics if not initialized.
func APPVERSION() *AppVersion {
	return MustAppVersion()
}

// MustAppVersion returns the global AppVersion instance.
// Panics if not initialized.
func MustAppVersion() *AppVersion {
	appVersion, err := GetAppVersion()
	if err != nil {
		panic(err)
	}
	return appVersion
}

// MustAppAtVersion returns the formatted string "app@version" of the global AppVersion instance.
// Panics if not initialized.
func MustAppAtVersion() string {
	appVersion, err := GetAppVersion()
	if err != nil {
		panic(err)
	}
	format, err := appVersion.Format(APPVERSIONFORMAT_APP_AT_VERSION)
	if err != nil {
		panic(err)
	}
	return format
}

// CloneAppVersion returns a clone of the global AppVersion instance.
// Returns nil if the global instance is not initialized or its Version is nil.
func CloneAppVersion() *AppVersion {
	global.mu.RLock()
	defer global.mu.RUnlock()
	if global.instance == nil || !autils.IsSemverValid(&global.instance.Version) {
		return nil
	}
	return global.instance.Clone()
}

// MustVersion returns the semantic version from the global AppVersion instance.
// Returns an empty semver.Version if not initialized or invalid.
func MustVersion() semver.Version {
	global.mu.RLock()
	defer global.mu.RUnlock()
	if global.instance == nil || !autils.IsSemverValid(&global.instance.Version) {
		return semver.Version{}
	}
	return global.instance.Version
}

// SetAppVersion sets the global AppVersion instance.
// Returns an error if already initialized.
func SetAppVersion(appVersion *AppVersion) error {
	global.mu.Lock()
	defer global.mu.Unlock()
	if global.instance != nil {
		return fmt.Errorf("AppVersion already initialized")
	}
	if appVersion == nil {
		return fmt.Errorf("AppVersion cannot be nil")
	}
	global.instance = appVersion
	return nil
}

// SetAppVersionByFS sets the global AppVersion instance from an embedded JSON file.
// Returns an error if the file cannot be read or is invalid.
var SetAppVersionByFS = func(fs embed.FS, myPath string) error {
	appVersion, err := LoadEmbeddedAppVersion(fs, myPath)
	if err != nil {
		return err
	}
	return SetAppVersion(appVersion)
}

// ResetAppVersion clears the global AppVersion instance and related state.
// Allows reinitialization for testing or dynamic reconfiguration.
func ResetAppVersion() {
	global.mu.Lock()
	defer global.mu.Unlock()
	global.instance = nil
}
