package ahttp

import (
	"github.com/jpfluger/alibs-slim/aapp"
	"github.com/jpfluger/alibs-slim/anetwork"
	"github.com/jpfluger/alibs-slim/asessions"
)

type IVersionProvider interface {
	GetAppVersion() *aapp.AppVersion
}

type IUrlProvider interface {
	RouteExists(httpRouteId HttpRouteId) bool
	MustUrl(httpRouteId HttpRouteId) string
	IfActiveUrlThenValue(activeUrl string, targetUrlKey HttpRouteId, value string) string
}

type ISiteConfig interface {
	GetMinExtension() string
	GetIsPrivateSite() bool
	GetPublicUrl() string
	GetConst(target string) string
	GetPublicNetUrl() *anetwork.NetURL
}

type IPageSessionController interface {
	IVersionProvider
	IUrlProvider
	ISiteConfig
	WRC() IWebRouteController
	// Permissions methods
	HasPerm(us asessions.ILoginSessionPerm, target asessions.Perm) bool
	HasPermS(us asessions.ILoginSessionPerm, keyPermValue string) bool
	HasPermSV(us asessions.ILoginSessionPerm, key string, value string) bool
	HasPermB(us asessions.ILoginSessionPerm, keyBits string) bool
	HasPermBV(us asessions.ILoginSessionPerm, key string, bit int) bool
	HasPermSet(us asessions.ILoginSessionPerm, target asessions.PermSet) bool
	HasPermKeyValueConst(us asessions.ILoginSessionPerm, key string, value string) bool
}

type PageSessionController struct {
	AppVersion         *aapp.AppVersion
	WebRouteController IWebRouteController
	MinExtension       string
	IsPrivateSite      bool
	PublicUrl          string
	Constants          map[string]string
	publicNetUrl       *anetwork.NetURL
}

func NewPageSessionController(appVersion *aapp.AppVersion, webRouteController IWebRouteController, minExt string, isPrivate bool, publicUrl string, constants map[string]string) *PageSessionController {
	if constants == nil {
		constants = make(map[string]string)
	}
	publicNetUrl, err := anetwork.ParseNetURL(publicUrl)
	if err != nil {
		panic(err)
	}
	return &PageSessionController{
		AppVersion:         appVersion,
		WebRouteController: webRouteController,
		MinExtension:       minExt,
		IsPrivateSite:      isPrivate,
		PublicUrl:          publicUrl,
		Constants:          constants,
		publicNetUrl:       publicNetUrl,
	}
}

// GetAppVersion returns the application version
func (psc *PageSessionController) GetAppVersion() *aapp.AppVersion {
	return psc.AppVersion
}

// RouteExists checks if a route exists for the given HttpRouteId
func (psc *PageSessionController) RouteExists(httpRouteId HttpRouteId) bool {
	return psc.WebRouteController.RouteExists(httpRouteId)
}

// MustUrl returns the URL for the specified HttpRouteId or panics if not found
func (psc *PageSessionController) MustUrl(httpRouteId HttpRouteId) string {
	return psc.WebRouteController.MustUrl(httpRouteId)
}

// GetMinExtension returns the minimum extension required
func (psc *PageSessionController) GetMinExtension() string {
	return psc.MinExtension
}

// GetIsPrivateSite indicates whether the site is private
func (psc *PageSessionController) GetIsPrivateSite() bool {
	return psc.IsPrivateSite
}

// GetPublicUrl returns the public URL of the site
func (psc *PageSessionController) GetPublicUrl() string {
	return psc.PublicUrl
}

// GetPublicNetUrl returns the public NetURL of the site
func (psc *PageSessionController) GetPublicNetUrl() *anetwork.NetURL {
	return psc.publicNetUrl
}

func (psc *PageSessionController) WRC() IWebRouteController {
	return psc.WebRouteController
}

// Define a global instance for IPageSessionController
var pscInstance IPageSessionController

// InitializePSC initializes the global instance of IPageSessionController.
// This function should be called once at program startup.
func InitializePSC(controller IPageSessionController) {
	if pscInstance != nil {
		panic("pscInstance already initialized")
	}
	pscInstance = controller
}

// PSC returns the global instance of IPageSessionController.
func PSC() IPageSessionController {
	if pscInstance == nil {
		panic("pscInstance is not initialized")
	}
	return pscInstance
}

func (ps *PageSessionController) IfActiveUrlThenValue(activeUrl string, targetUrlKey HttpRouteId, value string) string {
	if activeUrl == ps.MustUrl(targetUrlKey) {
		return value
	}
	return ""
}

func (ps *PageSessionController) GetConst(target string) string {
	return ps.Constants[target]
}

// HasPermS checks if the user session has a specific permission as a key-perm-value string.
func (ps *PageSessionController) HasPermS(us asessions.ILoginSessionPerm, keyPermValue string) bool {
	if us == nil || !us.IsLoggedIn() {
		return false
	}
	return us.HasPermS(keyPermValue)
}

// HasPermSV checks if the user session has a specific permission value for a given key.
func (ps *PageSessionController) HasPermSV(us asessions.ILoginSessionPerm, key string, value string) bool {
	if us == nil || !us.IsLoggedIn() {
		return false
	}
	if key == "" || value == "" {
		return false
	}
	return us.HasPermSV(key, value)
}

// HasPermB checks if the user session has a specific permission represented as a key-bit string.
func (ps *PageSessionController) HasPermB(us asessions.ILoginSessionPerm, keyBits string) bool {
	if us == nil || !us.IsLoggedIn() {
		return false
	}
	return us.HasPermB(keyBits)
}

// HasPermBV checks if the user session has a specific permission value for a given key using bit representation.
func (ps *PageSessionController) HasPermBV(us asessions.ILoginSessionPerm, key string, bit int) bool {
	if us == nil || !us.IsLoggedIn() {
		return false
	}
	if key == "" || bit <= 0 {
		return false
	}
	return us.HasPermBV(key, bit)
}

// HasPermSet checks if the user session has any matching permission with the target PermSet.
func (ps *PageSessionController) HasPermSet(us asessions.ILoginSessionPerm, target asessions.PermSet) bool {
	if us == nil || !us.IsLoggedIn() {
		return false
	}
	return us.HasPermSet(target)
}

// HasPermKeyValueConst checks if the user session has a specific permission value for a given key using constants.
func (ps *PageSessionController) HasPermKeyValueConst(us asessions.ILoginSessionPerm, key string, value string) bool {
	if us == nil || !us.IsLoggedIn() {
		return false
	}
	key = ps.GetConst(key)
	value = ps.GetConst(value)
	if key == "" || value == "" {
		return false
	}
	return us.HasPermSV(key, value)
}

// HasPerm checks if the user session has a specific permission object.
func (ps *PageSessionController) HasPerm(us asessions.ILoginSessionPerm, target asessions.Perm) bool {
	if us == nil || !us.IsLoggedIn() {
		return false
	}
	return us.HasPerm(target)
}

//func (ps *PageSessionController) HasPerm(us asessions.ILoginSessionPerm, keyValue string) bool {
//	if us == nil || !us.IsLoggedIn() {
//		return false
//	}
//	return us.HasPerm(keyValue)
//}
//
//func (ps *PageSessionController) HasPermKeyValue(us asessions.ILoginSessionPerm, key string, value string) bool {
//	if us == nil || !us.IsLoggedIn() {
//		return false
//	}
//	if key == "" || value == "" {
//		return false
//	}
//	return us.HasPermValue(key, value)
//}
//
//func (ps *PageSessionController) HasPermKeyValueConst(us asessions.ILoginSessionPerm, key string, value string) bool {
//	if us == nil || !us.IsLoggedIn() {
//		return false
//	}
//	key = ps.GetConst(key)
//	value = ps.GetConst(value)
//	if key == "" || value == "" {
//		return false
//	}
//	return us.HasPermValue(key, value)
//}

func GetPSCConstantsMapDefault() map[string]string {
	return map[string]string{
		"ACTIONKEY_2FACTOR":              asessions.ACTIONKEY_2FACTOR.String(),
		"ACTIONKEY_NEW_PASSWORD":         asessions.ACTIONKEY_NEW_PASSWORD.String(),
		"ACTIONKEY_ACCEPT_TERMS":         asessions.ACTIONKEY_ACCEPT_TERMS.String(),
		"ACTIONKEY_AFFIRM_PASSWORD":      asessions.ACTIONKEY_AFFIRM_PASSWORD.String(),
		"ACTIONKEY_CLICK_RESET_PASSWORD": asessions.ACTIONKEY_CLICK_RESET_PASSWORD.String(),
		"LOGINTYPE_SIMPLEAUTH":           asessions.LOGINTYPE_SIMPLEAUTH.String(),
		"LOGINTYPE_STEP_USERNAME":        asessions.LOGINTYPE_STEP_USERNAME.String(),
		"LOGINTYPE_STEP_PASSWORD":        asessions.LOGINTYPE_STEP_PASSWORD.String(),
		"LOGINTYPE_SIGNUP":               asessions.LOGINTYPE_SIGNUP.String(),
		"LOGINTYPE_FORGOT_LOGIN":         asessions.LOGINTYPE_FORGOT_LOGIN.String(),
		"PERMVALUE_X":                    asessions.PERMS_X,
		"PERMVALUE_L":                    asessions.PERMS_L,
		"PERMVALUE_C":                    asessions.PERMS_C,
		"PERMVALUE_R":                    asessions.PERMS_R,
		"PERMVALUE_U":                    asessions.PERMS_U,
		"PERMVALUE_D":                    asessions.PERMS_D,
		"PUI_DIRECTION_LEFT":             PUI_DIRECTION_LEFT.String(),
		"PUI_DIRECTION_RIGHT":            PUI_DIRECTION_RIGHT.String(),
		"PUI_SIDEBAR_STATUS_OPEN":        PUI_SIDEBAR_STATUS_OPEN.String(),
		"PUI_SIDEBAR_STATUS_CLOSE":       PUI_SIDEBAR_STATUS_CLOSE.String(),
		"PUI_SIDEBAR_STATUS_DISMISS":     PUI_SIDEBAR_STATUS_DISMISS.String(),
	}
}
