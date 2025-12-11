package asessions

import (
	"time"

	"github.com/jpfluger/alibs-slim/auser"

	"github.com/jpfluger/alibs-slim/atags"
	"github.com/jpfluger/alibs-slim/azb"
)

// Constants for different types of login actions.
const (
	// Login types with a status of LOGIN_SESSION_STATUS_NONE.
	LOGINTYPE_SIMPLEAUTH    = azb.ZBType("simple-auth")
	LOGINTYPE_STEP_USERNAME = azb.ZBType("step-username")
	LOGINTYPE_STEP_PASSWORD = azb.ZBType("step-password")
	LOGINTYPE_SIGNUP        = azb.ZBType("signup")
	LOGINTYPE_FORGOT_LOGIN  = azb.ZBType("forgot-login")

	// Action keys with a status of LOGIN_SESSION_STATUS_ACTIONS.
	ACTIONKEY_2FACTOR               = ActionKey("2factor")
	ACTIONKEY_NEW_PASSWORD          = ActionKey("new-password")
	ACTIONKEY_ACCEPT_TERMS          = ActionKey("accept-terms")
	ACTIONKEY_AFFIRM_PASSWORD       = ActionKey("affirm-password")
	ACTIONKEY_CLICK_RESET_PASSWORD  = ActionKey("click-reset-password")
	ACTIONKEY_CLICK_EMAILVERIFY     = ActionKey("click-email-verify")
	ACTIONKEY_JWT_USER_FORGOT_LOGIN = ActionKey("jwt-user-forgot-login")

	// Action keys enabling policy-based authentication solutions.
	ACTIONKEY_IAM_PRIMARY = ActionKey("iam-primary")
	ACTIONKEY_IAM_MFA     = ActionKey("iam-mfa")
	ACTIONKEY_IAM_SSPR    = ActionKey("iam-sspr")

	// Tag key for JWT metadata.
	TAGKEY_JWT_META = atags.TagKey("jwt-meta")
)

// IUILoginUser defines an interface for UI login user data.
type IUILoginUser interface {
	GetLoginType() azb.ZBType
	GetUrlPost() string
	GetUsername() auser.Username
	GetFormDate() time.Time
	SetFormDate(formDate time.Time)
	GetSecret() string
	GetIsOnRememberMe() bool
	GetTitle() string
	GetLead() string
	GetFooter() string
	GetActiveTypes() azb.ZBTypes
	GetTags() atags.TagMapString
}

// UILoginUser holds the data for a user login interface.
type UILoginUser struct {
	LoginType azb.ZBType `json:"loginType"` // Type of login action.

	UrlPost string `json:"urlPost,omitempty"` // URL for POST requests.

	FormDate time.Time `json:"formDate,omitempty"` // Date of the form submission.

	Username       auser.Username `json:"username,omitempty"`       // Username of the user.
	Secret         string         `json:"secret,omitempty"`         // Secret for authentication.
	IsOnRememberMe bool           `json:"isOnRememberMe,omitempty"` // Flag for 'remember me' feature.

	Title  string `json:"title"`  // Title for the login form.
	Lead   string `json:"lead"`   // Lead text for the login form.
	Footer string `json:"footer"` // Footer text for the login form.

	ActiveTypes azb.ZBTypes `json:"activeTypes"` // Active types for the login form.

	Tags atags.TagMapString `json:"tags"` // Tags associated with the login form.
}

// NewUILoginUser creates a new UILoginUser with default options.
func NewUILoginUser(loginType azb.ZBType, username auser.Username, urlPost string) *UILoginUser {
	return NewUILoginUserWithOptions(loginType, username, urlPost, "", "", "", time.Time{}, nil, false, nil)
}

// NewUILoginUserWithOptions creates a new UILoginUser with specified options.
func NewUILoginUserWithOptions(loginType azb.ZBType, username auser.Username, urlPost string, title string, leader string, footer string, formDate time.Time, activeTypes azb.ZBTypes, isOnRememberMe bool, tags atags.TagMapString) *UILoginUser {
	// Initialize empty types and tags if nil.
	if activeTypes == nil {
		activeTypes = azb.ZBTypes{}
	}
	if tags == nil {
		tags = atags.TagMapString{}
	}
	// Set the form date to current UTC time if not provided.
	if formDate.IsZero() {
		formDate = time.Now().UTC()
	}
	return &UILoginUser{
		LoginType:      loginType,
		UrlPost:        urlPost,
		Username:       username,
		FormDate:       formDate,
		Title:          title,
		Lead:           leader,
		Footer:         footer,
		ActiveTypes:    activeTypes,
		IsOnRememberMe: isOnRememberMe,
		Tags:           tags,
	}
}

// Getters and setters for UILoginUser fields.
func (u *UILoginUser) GetLoginType() azb.ZBType {
	return u.LoginType
}

func (u *UILoginUser) GetUrlPost() string {
	return u.UrlPost
}

func (u *UILoginUser) GetFormDate() time.Time {
	return u.FormDate
}

func (u *UILoginUser) SetFormDate(formDate time.Time) {
	u.FormDate = formDate
}

func (u *UILoginUser) GetUsername() auser.Username {
	return u.Username
}

func (u *UILoginUser) GetSecret() string {
	return u.Secret
}

func (u *UILoginUser) GetIsOnRememberMe() bool {
	return u.IsOnRememberMe
}

func (u *UILoginUser) GetTitle() string {
	return u.Title
}

func (u *UILoginUser) GetLead() string {
	return u.Lead
}

func (u *UILoginUser) GetFooter() string {
	return u.Footer
}

func (u *UILoginUser) GetActiveTypes() azb.ZBTypes {
	return u.ActiveTypes
}

func (u *UILoginUser) GetTags() atags.TagMapString {
	return u.Tags
}
