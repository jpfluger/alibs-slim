package aapp

import (
	"embed"
	"encoding/json"
	"fmt"
	"github.com/Masterminds/semver/v3"
	"github.com/jpfluger/alibs-slim/autils"
	"strings"
)

var (
	buildName      string
	buildType      string
	deploymentType string
)

func BUILDNAME() BuildName {
	return BuildName(buildName)
}

func BUILDTYPE() BuildType {
	return BuildType(buildType)
}

func DEPLOYMENTTYPE() DeploymentType {
	return DeploymentType(deploymentType)
}

// AppVersion represents the version information of an application.
type AppVersion struct {
	Name      string          `json:"name,omitempty"`      // The name of the application.
	Version   *semver.Version `json:"version,omitempty"`   // The semantic version of the application.
	Title     string          `json:"title,omitempty"`     // The title of the application.
	About     string          `json:"about,omitempty"`     // Additional information about the application.
	Owner     string          `json:"owner,omitempty"`     // The owner of the application.
	LegalMark string          `json:"legalMark,omitempty"` // Legal trademark or copyright notice.

	// BuildName should be set by the compiler.
	// In absence of BuildName a trimmed/lowercase Name is used.
	BuildName BuildName `json:"buildName,omitempty"`

	// BuildType should be set by the compiler.
	// If empty, then BUILDTYPE_DEBUG is used.
	BuildType BuildType `json:"buildType,omitempty"`

	// DeploymentType should be set by the compiler.
	// If empty, then DEPLOYMENTTYPE_DEV is used.
	DeploymentType DeploymentType `json:"deploymentType,omitempty"`
}

// AppVersionFormat defines the format in which the version information can be presented.
type AppVersionFormat string

// Constants for different version formats.
const (
	APPVERSION_FORMAT_JSON           AppVersionFormat = "json"
	APPVERSION_FORMAT_VERSION_ONLY   AppVersionFormat = "version"
	APPVERSION_FORMAT_NO_V           AppVersionFormat = "version-no-v"
	APPVERSION_FORMAT_APP_ONLY       AppVersionFormat = "name"
	APPVERSION_FORMAT_ABOUT_ONLY     AppVersionFormat = "about"
	APPVERSION_FORMAT_APP_AT_VERSION AppVersionFormat = "app@version" // default
	APPVERSION_FORMAT_BUILD          AppVersionFormat = "build"       // default
)

// IsEmpty checks if the AppVersionFormat is empty after trimming whitespace.
func (bt AppVersionFormat) IsEmpty() bool {
	return strings.TrimSpace(string(bt)) == ""
}

// Format returns the version information in the specified format.
func (a *AppVersion) Format(format AppVersionFormat) (string, error) {
	switch format {
	case APPVERSION_FORMAT_JSON:
		b, err := json.MarshalIndent(a, "", "  ")
		if err != nil {
			return "", fmt.Errorf("failed to marshal json for AppVersion: %v", err)
		}
		return string(b), nil
	case APPVERSION_FORMAT_VERSION_ONLY:
		return fmt.Sprintf("v%s", a.Version.String()), nil
	case APPVERSION_FORMAT_NO_V:
		return a.Version.String(), nil
	case APPVERSION_FORMAT_APP_ONLY:
		return a.BuildName.String(), nil
	case APPVERSION_FORMAT_ABOUT_ONLY:
		return a.About, nil
	case APPVERSION_FORMAT_BUILD:
		return fmt.Sprintf("%s@%s,%s-%s", a.BuildName.String(), a.Version.String(), a.BuildType.String(), a.DeploymentType.String()), nil
	default: // APPVERSION_FORMAT_APP_AT_VERSION
		return fmt.Sprintf("%s@%s", a.BuildName.String(), a.Version.String()), nil
	}
}

// Validate checks if the AppVersion fields are valid.
func (a *AppVersion) Validate() error {
	if a.Name == "" {
		return fmt.Errorf("version name cannot be empty")
	}
	if a.Owner == "" {
		return fmt.Errorf("version owner cannot be empty")
	}
	if a.LegalMark == "" {
		return fmt.Errorf("version legalMark cannot be empty")
	}
	if a.Version == nil || a.Version.String() == "" {
		return fmt.Errorf("version cannot be empty")
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
	}
	if a.DeploymentType.IsEmpty() {
		if DEPLOYMENTTYPE().IsEmpty() {
			a.DeploymentType = DEPLOYMENTTYPE_DEV
		} else {
			a.DeploymentType = DEPLOYMENTTYPE()
		}
	}
	return nil
}

func (a *AppVersion) GetTitle() string {
	return a.Title
}

func (a *AppVersion) GetLongTitle() string {
	if a.BuildType.IsEmpty() && a.DeploymentType.IsEmpty() {
		return fmt.Sprintf("%s v%s", a.Title, a.Version.String())
	}
	return fmt.Sprintf("%s v%s, %s-%s",
		a.Title,
		a.Version.String(),
		a.BuildType.String(),
		a.DeploymentType.String(),
	)
}

func (a *AppVersion) SetBuildName(newBuildName string) {
	buildName = newBuildName
	a.BuildName = BUILDNAME()
}

func (a *AppVersion) SetBuildType(newBuildType string) {
	buildType = newBuildType
	a.BuildType = BUILDTYPE()
}

func (a *AppVersion) SetDeploymentType(newDeploymentType string) {
	deploymentType = newDeploymentType
	a.DeploymentType = DEPLOYMENTTYPE()
}

func (a *AppVersion) GetBuildName() BuildName {
	return a.BuildName
}

func (a *AppVersion) GetBuildType() BuildType {
	return a.BuildType
}

func (a *AppVersion) GetDeploymentType() DeploymentType {
	return a.DeploymentType
}

// LoadAppVersionFromBytes loads the AppVersion from a byte slice.
func LoadAppVersionFromBytes(b []byte) (*AppVersion, error) {
	if len(b) == 0 {
		return nil, fmt.Errorf("no bytes to load")
	}
	appVersion := &AppVersion{}
	if err := json.Unmarshal(b, appVersion); err != nil {
		return nil, fmt.Errorf("failed to unmarshal AppVersion: %v", err)
	}
	appVersion.Name = strings.TrimSpace(appVersion.Name)
	appVersion.About = strings.TrimSpace(appVersion.About)
	appVersion.Owner = strings.TrimSpace(appVersion.Owner)
	appVersion.LegalMark = strings.TrimSpace(appVersion.LegalMark)
	if err := appVersion.Validate(); err != nil {
		return nil, fmt.Errorf("AppVersion failed validation: %v", err)
	}
	return appVersion, nil
}

func LoadEmbeddedAppVersion(fs embed.FS, myPath string) (*AppVersion, error) {
	myPath = strings.TrimSpace(myPath)
	if myPath == "" {
		myPath = "app/version.json"
	}
	if b, err := fs.ReadFile(myPath); err != nil {
		return nil, fmt.Errorf("failed to retrieve embedded file; %v", err)
	} else {
		appVersion := &AppVersion{}
		if err = json.Unmarshal(b, appVersion); err != nil {
			return nil, fmt.Errorf("failed to unmarshal AppVersion; %v", err)
		}
		appVersion.Name = strings.TrimSpace(appVersion.Name)
		appVersion.About = strings.TrimSpace(appVersion.About)
		appVersion.Owner = strings.TrimSpace(appVersion.Owner)
		appVersion.LegalMark = strings.TrimSpace(appVersion.LegalMark)
		if err = appVersion.Validate(); err != nil {
			return nil, fmt.Errorf("AppVersion failed validation; %v", err)
		}
		return appVersion, nil
	}
}

var appVersionInstance *AppVersion

func GetAppVersion() *AppVersion {
	if appVersionInstance == nil {
		panic("appVersionInstance is not initialized")
	}
	return appVersionInstance
}

func SetAppVersionByFS(fs embed.FS, myPath string) error {
	appVersion, err := LoadEmbeddedAppVersion(fs, myPath)
	if err != nil {
		return err
	}
	SetAppVersion(appVersion)
	return nil
}

func SetAppVersion(appVersion *AppVersion) {
	if appVersionInstance != nil {
		panic("appVersionInstance already initialized")
	}
	appVersionInstance = appVersion
}
