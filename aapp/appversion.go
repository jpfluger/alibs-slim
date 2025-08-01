package aapp

import (
	"embed"
	"encoding/json"
	"fmt"
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
	Name           string          `json:"name,omitempty"`           // The name of the application.
	Version        *semver.Version `json:"version,omitempty"`        // The semantic version of the application.
	Title          string          `json:"title,omitempty"`          // The title of the application.
	About          string          `json:"about,omitempty"`          // Additional information about the application.
	Owner          string          `json:"owner,omitempty"`          // The owner of the application.
	LegalMark      string          `json:"legalMark,omitempty"`      // Legal trademark or copyright notice.
	BuildName      BuildName       `json:"buildName,omitempty"`      // The name of the built binary, defaults to trimmed/lowercase Name if empty.
	BuildType      BuildType       `json:"buildType,omitempty"`      // The build type, defaults to BUILDTYPE_DEBUG if empty.
	DeploymentType DeploymentType  `json:"deploymentType,omitempty"` // The deployment type, defaults to DEPLOYMENTTYPE_DEV if empty.
}

// globalAppVersion holds the global AppVersion instance and related state.
type globalAppVersion struct {
	mu             sync.RWMutex
	instance       *AppVersion
	buildName      BuildName
	buildType      BuildType
	deploymentType DeploymentType
}

var global = &globalAppVersion{}

// BUILDNAME returns the global BuildName value, typically set by the compiler.
func BUILDNAME() BuildName {
	global.mu.RLock()
	defer global.mu.RUnlock()
	return global.buildName
}

// BUILDTYPE returns the global BuildType value, typically set by the compiler.
func BUILDTYPE() BuildType {
	global.mu.RLock()
	defer global.mu.RUnlock()
	return global.buildType
}

// DEPLOYMENTTYPE returns the global DeploymentType value, typically set by the compiler.
func DEPLOYMENTTYPE() DeploymentType {
	global.mu.RLock()
	defer global.mu.RUnlock()
	return global.deploymentType
}

// Format returns the version information in the specified format.
// Returns an error if marshaling fails for JSON format or if version is nil for certain formats.
func (a *AppVersion) Format(format AppVersionFormat) (string, error) {
	switch format {
	case APPVERSIONFORMAT_JSON:
		data, err := json.MarshalIndent(a, "", "  ")
		if err != nil {
			return "", fmt.Errorf("failed to marshal JSON for AppVersion: %w", err)
		}
		return string(data), nil
	case APPVERSIONFORMAT_VERSION_ONLY:
		if a.Version == nil {
			return "", fmt.Errorf("version cannot be nil")
		}
		return fmt.Sprintf("v%s", a.Version.String()), nil
	case APPVERSIONFORMAT_NO_V:
		if a.Version == nil {
			return "", fmt.Errorf("version cannot be nil")
		}
		return a.Version.String(), nil
	case APPVERSIONFORMAT_APP_ONLY:
		return a.BuildName.String(), nil
	case APPVERSIONFORMAT_ABOUT_ONLY:
		return a.About, nil
	case APPVERSIONFORMAT_BUILD:
		if a.Version == nil {
			return "", fmt.Errorf("version cannot be nil")
		}
		return fmt.Sprintf("%s@%s,%s-%s", a.BuildName.String(), a.Version.String(), a.BuildType.String(), a.DeploymentType.String()), nil
	default: // APPVERSIONFORMAT_APP_AT_VERSION
		if a.Version == nil {
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
	if a.Version == nil || a.Version.String() == "" {
		return fmt.Errorf("version cannot be nil or empty")
	}
	if a.BuildName.IsEmpty() {
		if BUILDNAME().IsEmpty() {
			a.BuildName = BuildName(autils.ToStringTrimLower(a.Name))
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
		if DEPLOYMENTTYPE().IsEmpty() {
			a.DeploymentType = DEPLOYMENTTYPE_DEV
		} else {
			a.DeploymentType = DEPLOYMENTTYPE()
		}
	} else if !a.DeploymentType.IsValid() {
		return fmt.Errorf("invalid deployment type %q", a.DeploymentType)
	}
	return nil
}

// GetTitle returns the title of the application.
func (a *AppVersion) GetTitle() string {
	return a.Title
}

// GetLongTitle returns a detailed title including version, build type, and deployment type.
func (a *AppVersion) GetLongTitle() string {
	if a.BuildType.IsEmpty() && a.DeploymentType.IsEmpty() {
		if a.Version == nil {
			return a.Title
		}
		return fmt.Sprintf("%s v%s", a.Title, a.Version.String())
	}
	if a.Version == nil {
		return fmt.Sprintf("%s, %s-%s", a.Title, a.BuildType.String(), a.DeploymentType.String())
	}
	return fmt.Sprintf("%s v%s, %s-%s", a.Title, a.Version.String(), a.BuildType.String(), a.DeploymentType.String())
}

// SetBuildName sets the global build name and updates the AppVersion instance.
func (a *AppVersion) SetBuildName(newBuildName BuildName) {
	global.mu.Lock()
	defer global.mu.Unlock()
	global.buildName = newBuildName
	a.BuildName = newBuildName
}

// SetBuildType sets the global build type and updates the AppVersion instance.
// Returns an error if the build type is invalid.
func (a *AppVersion) SetBuildType(newBuildType BuildType) error {
	if !newBuildType.IsValid() {
		return fmt.Errorf("invalid build type %q", newBuildType)
	}
	global.mu.Lock()
	defer global.mu.Unlock()
	global.buildType = newBuildType
	a.BuildType = newBuildType
	return nil
}

// SetDeploymentType sets the global deployment type and updates the AppVersion instance.
// Returns an error if the deployment type is invalid.
func (a *AppVersion) SetDeploymentType(newDeploymentType DeploymentType) error {
	if !newDeploymentType.IsValid() {
		return fmt.Errorf("invalid deployment type %q", newDeploymentType)
	}
	global.mu.Lock()
	defer global.mu.Unlock()
	global.deploymentType = newDeploymentType
	a.DeploymentType = newDeploymentType
	return nil
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
func GetAppVersion() (*AppVersion, error) {
	global.mu.RLock()
	defer global.mu.RUnlock()
	if global.instance == nil {
		return nil, fmt.Errorf("AppVersion not initialized")
	}
	return global.instance, nil
}

// MustAppVersion returns the global AppVersion instance.
func MustAppVersion() *semver.Version {
	global.mu.RLock()
	defer global.mu.RUnlock()
	if global.instance == nil || global.instance.Version == nil {
		return &semver.Version{}
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
func SetAppVersionByFS(fs embed.FS, myPath string) error {
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
	global.buildName = ""
	global.buildType = ""
	global.deploymentType = ""
}
