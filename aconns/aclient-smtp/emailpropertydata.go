package aclient_smtp

import (
	"sync"
	"time"
)

var (
	appEmailPropertyDataMu sync.RWMutex
	appEmailPropertyData   *EmailPropertyData
)

// SetEmailPropertyData sets the global EmailPropertyData.
// It is thread-safe and allows setting only once; subsequent calls are ignored.
// For reinitialization in tests or dynamic apps, consider a separate Reset function if needed.
func SetEmailPropertyData(e *EmailPropertyData) {
	appEmailPropertyDataMu.Lock()
	defer appEmailPropertyDataMu.Unlock()
	appEmailPropertyData = e
}

func EMAILPROPERTYDATA() *EmailPropertyData {
	return getAppEmailPropertyData()
}

// getAppEmailPropertyData returns a clone of the global EmailPropertyData.
// It is thread-safe for concurrent reads.
func getAppEmailPropertyData() *EmailPropertyData {
	appEmailPropertyDataMu.RLock()
	defer appEmailPropertyDataMu.RUnlock()
	if appEmailPropertyData == nil {
		return nil
	}
	return appEmailPropertyData.Clone()
}

// NewEmailPropertyData creates a new EmailPropertyData, merging in global app defaults.
// It uses the provided values and fills in missing ones from the global appEmailPropertyData.
// This ensures consistency across emails without hardcoding.
func NewEmailPropertyData(username string, clickLink string, expirationTime time.Time, orgContactName string) *EmailPropertyData {
	myData := &EmailPropertyData{
		Username:       username,
		ClickLink:      clickLink,
		ExpirationTime: expirationTime,
		OrgContactName: orgContactName,
	}
	globalData := getAppEmailPropertyData()
	if globalData != nil {
		globalData.MergeTo(myData)
	}
	return myData
}

// NewEmailPropertyDataFromData creates a new EmailPropertyData by merging global defaults into the provided overrideData.
// If overrideData is nil, it returns a clone of the global data (or an empty struct if no global is set).
// This is useful for scenarios where overrides take precedence but globals fill in gaps.
func NewEmailPropertyDataFromData(overrideData *EmailPropertyData) *EmailPropertyData {
	globalData := getAppEmailPropertyData()
	if overrideData == nil {
		if globalData == nil {
			return &EmailPropertyData{}
		}
		return globalData
	}
	if globalData != nil {
		globalData.MergeTo(overrideData)
	}
	return overrideData
}

// EmailPropertyData defines common placeholders used in email templates.
type EmailPropertyData struct {
	// Section 1: Used
	Username       string    `json:"username"`       // Recipient's username
	AppName        string    `json:"appName"`        // Application name
	PublicUrl      string    `json:"publicUrl"`      // The public url
	ClickLink      string    `json:"clickLink"`      // URL or token link for actions
	ExpirationTime time.Time `json:"expirationTime"` // Link expiration duration
	OrgName        string    `json:"orgName"`        // Organization name
	OrgContactName string    `json:"orgContactName"` // Support contact's name

	// Section 2: Optional
	//OrgPhone     string `json:"orgPhone"`     // Organization phone
	//OrgEmail     string `json:"orgEmail"`     // Organization email
	//UserDisplayName     string `json:"userDisplayName"`     // Recipient's display name
}

// MergeTo merges the receiver's data into `data` where `data`'s fields are empty or zero.
func (e *EmailPropertyData) MergeTo(data *EmailPropertyData) {
	if data.Username == "" {
		data.Username = e.Username
	}
	if data.AppName == "" {
		data.AppName = e.AppName
	}
	if data.PublicUrl == "" {
		data.PublicUrl = e.PublicUrl
	}
	if data.ClickLink == "" {
		data.ClickLink = e.ClickLink
	}
	if data.ExpirationTime.IsZero() {
		data.ExpirationTime = e.ExpirationTime
	}
	if data.OrgName == "" {
		data.OrgName = e.OrgName
	}
	if data.OrgContactName == "" {
		data.OrgContactName = e.OrgContactName
	}
}

// Clone returns a deep copy of the EmailPropertyData.
// Since all fields are value types (strings and time.Time), a shallow copy suffices.
func (e *EmailPropertyData) Clone() *EmailPropertyData {
	if e == nil {
		return nil
	}
	return &EmailPropertyData{
		Username:       e.Username,
		AppName:        e.AppName,
		PublicUrl:      e.PublicUrl,
		ClickLink:      e.ClickLink,
		ExpirationTime: e.ExpirationTime,
		OrgName:        e.OrgName,
		OrgContactName: e.OrgContactName,
	}
}

// EmailPropertyDataContent extends EmailPropertyData with additional ContentText and ContentHTML fields.
// It is useful for email bodies that require custom content alongside property data, supporting both text and HTML variants.
type EmailPropertyDataContent struct {
	EmailPropertyData
	ContentText string `json:"contentText"`
	ContentHTML string `json:"contentHTML"`
}

// NewEmailPropertyDataContent creates a new EmailPropertyDataContent.
// It merges global defaults into the overrideData (via NewEmailPropertyDataFromData) and sets the ContentText and ContentHTML.
func NewEmailPropertyDataContent(overrideData *EmailPropertyData, contentText string, contentHTML string) *EmailPropertyDataContent {
	return &EmailPropertyDataContent{
		EmailPropertyData: *NewEmailPropertyDataFromData(overrideData),
		ContentText:       contentText,
		ContentHTML:       contentHTML,
	}
}
