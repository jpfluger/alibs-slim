package aclient_smtp

import "time"

var appEmailPropertyData *EmailPropertyData

func SetEmailPropertyData(e *EmailPropertyData) {
	if appEmailPropertyData != nil {
		panic("appEmailPropertyData already initialized")
	}
	appEmailPropertyData = e
}

func NewEmailPropertyData(username string, clickLink string, expirationTime time.Time, orgContactName string) *EmailPropertyData {
	myData := &EmailPropertyData{
		Username:       username,
		ClickLink:      clickLink,
		ExpirationTime: expirationTime,
		OrgContactName: orgContactName,
	}
	myData.AppName = appEmailPropertyData.AppName
	myData.PublicUrl = appEmailPropertyData.PublicUrl
	myData.OrgName = appEmailPropertyData.OrgName
	return myData
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
}
