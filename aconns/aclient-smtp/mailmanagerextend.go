package aclient_smtp

import (
	"fmt"
	"strings"
	"time"
)

const (
	MAILTEMPLATE_TEST          MailTemplateName = "test"          // Template name for test emails.
	MAILTEMPLATE_SYSTEM_NOTIFY MailTemplateName = "system-notify" // Template name for system notifications.
)

// MailManagerExtend extends MailManager with additional methods to simplify bootstrapping and common email operations in web services.
type MailManagerExtend struct {
	MailManager // Embeds MailManager for core functionality.
}

// SetupMailManager initializes the MailManager with app-specific properties and sets it as the global instance.
// It validates the configuration and sets global email property data for consistent use in templates.
func (mme *MailManagerExtend) SetupMailManager(appName string, publicURL string, orgName string) error {
	if !mme.MailManager.GetIsActive() {
		return nil
	}

	if err := mme.MailManager.Validate(); err != nil {
		return fmt.Errorf("failed to initialize mail manager: %v", err)
	}
	SetMailManager(&mme.MailManager)
	SetEmailPropertyData(&EmailPropertyData{
		AppName:   appName,
		PublicUrl: publicURL,
		OrgName:   orgName,
	})

	return nil
}

// SendTestEmail sends a test email using the "test-email" template.
// It merges the provided MAG with the specified key's default (defaults to system if empty),
// constructs a subject with the app name and timestamp, and uses default email properties with the provided click URL and expiration.
func (mme *MailManagerExtend) SendTestEmail(magKey MailAddressGroupKey, mag *MailAddressGroup, urlClick string, urlExpires time.Time) error {
	if !mme.MailManager.GetIsActive() {
		return fmt.Errorf("inactive mail manager")
	}
	if magKey.IsEmpty() {
		magKey = MAG_KEY_SYSTEM
	}
	mergedMAG := mme.MailManager.FromMAG(magKey, mag)
	if mergedMAG == nil {
		return fmt.Errorf("failed to get merged MAG for key %q", magKey)
	}
	subjectMerge := []interface{}{fmt.Sprintf("%s, %s", getAppEmailPropertyData().AppName, time.Now().Format("2006-01-02 15:04:05"))}
	testData := NewEmailPropertyData("Test User", urlClick, urlExpires, "")
	if err := mme.MailManager.SendWithRender(MAILTEMPLATE_TEST, mergedMAG, subjectMerge, testData); err != nil {
		return err
	}
	return nil
}

// SendEmailSystemNotify sends a system notification email using the "system-notify" template.
// It merges the provided MAG with the system default if applicable, constructs the subject with a timestamp,
// and uses EmailPropertyDataContent for the body with custom text and HTML content.
func (mme *MailManagerExtend) SendEmailSystemNotify(mag *MailAddressGroup, notifySubject string, contentText string, contentHTML string) error {
	if !mme.MailManager.GetIsActive() {
		return fmt.Errorf("inactive mail manager")
	}
	var mergedMAG *MailAddressGroup
	if mag == nil {
		mergedMAG = mme.MailManager.MAGS.Get(MAG_KEY_SYSTEM)
	} else {
		// mag is optional; merge with system default
		mergedMAG = mme.MailManager.FromMAG(MAG_KEY_SYSTEM, mag)
	}
	if mergedMAG == nil {
		return fmt.Errorf("failed to get system MAG")
	}
	notifySubject = strings.TrimSpace(notifySubject)
	if notifySubject == "" {
		notifySubject = "General"
	}
	subjectMerge := []interface{}{fmt.Sprintf("%s, %s", notifySubject, time.Now().Format("2006-01-02 15:04:05"))}
	bodyData := NewEmailPropertyDataContent(nil, contentText, contentHTML)
	if err := mme.MailManager.SendWithRender(MAILTEMPLATE_SYSTEM_NOTIFY, mergedMAG, subjectMerge, bodyData); err != nil {
		return err
	}
	return nil
}

// SendEmailSystem sends a general system email using the specified template.
// It merges the provided MAG with the system default if applicable.
func (mme *MailManagerExtend) SendEmailSystem(templateName MailTemplateName, mag *MailAddressGroup, subjectMerge []interface{}, dataBody interface{}) error {
	if !mme.MailManager.GetIsActive() {
		return fmt.Errorf("inactive mail manager")
	}
	if templateName.IsEmpty() {
		return fmt.Errorf("template name is empty")
	}
	var mergedMAG *MailAddressGroup
	if mag == nil {
		mergedMAG = mme.MailManager.MAGS.Get(MAG_KEY_SYSTEM)
	} else {
		// mag is optional; merge with system default
		mergedMAG = mme.MailManager.FromMAG(MAG_KEY_SYSTEM, mag)
	}
	if mergedMAG == nil {
		return fmt.Errorf("failed to get system MAG")
	}
	if err := mme.MailManager.SendWithRender(templateName, mergedMAG, subjectMerge, dataBody); err != nil {
		return err
	}
	return nil
}

// SendEmailUser sends a general user email using the specified template.
// It merges the provided MAG with the users default if applicable.
func (mme *MailManagerExtend) SendEmailUser(templateName MailTemplateName, mag *MailAddressGroup, subjectMerge []interface{}, dataBody interface{}) error {
	if !mme.MailManager.GetIsActive() {
		return fmt.Errorf("inactive mail manager")
	}
	if templateName.IsEmpty() {
		return fmt.Errorf("template name is empty")
	}
	var mergedMAG *MailAddressGroup
	if mag == nil {
		mergedMAG = mme.MailManager.MAGS.Get(MAG_KEY_USERS)
	} else {
		// mag is optional; merge with users default
		mergedMAG = mme.MailManager.FromMAG(MAG_KEY_USERS, mag)
	}
	if mergedMAG == nil {
		return fmt.Errorf("failed to get users MAG")
	}
	if err := mme.MailManager.SendWithRender(templateName, mergedMAG, subjectMerge, dataBody); err != nil {
		return err
	}
	return nil
}

// SendEmailSystemAsync sends a system email asynchronously using the specified template.
// It merges the provided MAG with the system default if applicable.
func (mme *MailManagerExtend) SendEmailSystemAsync(templateName MailTemplateName, mag *MailAddressGroup, subjectMerge []interface{}, dataBody interface{}, fnCallback FNMailSendCallback) {
	var mergedMAG *MailAddressGroup
	if mag == nil {
		mergedMAG = mme.MailManager.MAGS.Get(MAG_KEY_SYSTEM)
	} else {
		mergedMAG = mme.MailManager.FromMAG(MAG_KEY_SYSTEM, mag)
	}
	mme.MailManager.SendWithRenderAsync(templateName, mergedMAG, subjectMerge, dataBody, fnCallback)
}

// SendEmailUserAsync sends a user email asynchronously using the specified template.
// It merges the provided MAG with the users default if applicable.
func (mme *MailManagerExtend) SendEmailUserAsync(templateName MailTemplateName, mag *MailAddressGroup, subjectMerge []interface{}, dataBody interface{}, fnCallback FNMailSendCallback) {
	var mergedMAG *MailAddressGroup
	if mag == nil {
		mergedMAG = mme.MailManager.MAGS.Get(MAG_KEY_USERS)
	} else {
		mergedMAG = mme.MailManager.FromMAG(MAG_KEY_USERS, mag)
	}
	mme.MailManager.SendWithRenderAsync(templateName, mergedMAG, subjectMerge, dataBody, fnCallback)
}

//package aclient_smtp
//
//import (
//	"fmt"
//	"strings"
//	"time"
//)
//
//const (
//	MAILTEMPLATE_TEST_EMAIL    MailTemplateName = "test-email"    // Template name for test emails.
//	MAILTEMPLATE_SYSTEM_NOTIFY MailTemplateName = "system-notify" // Template name for system notifications.
//)
//
//// MailManagerExtend extends MailManager with additional methods to simplify bootstrapping and common email operations in web services.
//type MailManagerExtend struct {
//	MailManager // Embeds MailManager for core functionality.
//}
//
//// SetupMailManager initializes the MailManager with app-specific properties and sets it as the global instance.
//// It validates the configuration and sets global email property data for consistent use in templates.
//func (mme *MailManagerExtend) SetupMailManager(appName string, publicURL string, orgName string) error {
//	if !mme.MailManager.GetIsActive() {
//		return nil
//	}
//
//	if err := mme.MailManager.Validate(); err != nil {
//		return fmt.Errorf("failed to initialize mail manager: %v", err)
//	}
//	SetMailManager(&mme.MailManager)
//	SetEmailPropertyData(&EmailPropertyData{
//		AppName:   appName,
//		PublicUrl: publicURL,
//		OrgName:   orgName,
//	})
//
//	return nil
//}
//
//// SendTestEmail sends a test email using the "test-email" template.
//// It merges the provided MAG with the specified key's default (defaults to system if empty),
//// constructs a subject with the app name and timestamp, and uses default email properties with the provided click URL and expiration.
//func (mme *MailManagerExtend) SendTestEmail(magKey MailAddressGroupKey, mag *MailAddressGroup, urlClick string, urlExpires time.Time) error {
//	if !mme.MailManager.GetIsActive() {
//		return fmt.Errorf("inactive mail manager")
//	}
//	if magKey.IsEmpty() {
//		magKey = MAG_KEY_SYSTEM
//	}
//	mergedMAG := mme.MailManager.FromMAG(magKey, mag)
//	if mergedMAG == nil {
//		return fmt.Errorf("failed to get merged MAG for key %q", magKey)
//	}
//	subjectMerge := []interface{}{fmt.Sprintf("%s, %s", getAppEmailPropertyData().AppName, time.Now().Format("2006-01-02 15:04:05"))}
//	testData := NewEmailPropertyData("Test User", urlClick, urlExpires, "")
//	if err := mme.MailManager.SendWithRender(MAILTEMPLATE_TEST_EMAIL, mergedMAG, subjectMerge, testData); err != nil {
//		return err
//	}
//	return nil
//}
//
//// SendEmailSystemNotify sends a system notification email using the "system-notify" template.
//// It merges the provided MAG with the system default if applicable, constructs the subject with a timestamp,
//// and uses EmailPropertyDataContent for the body with custom text and HTML content.
//func (mme *MailManagerExtend) SendEmailSystemNotify(mag *MailAddressGroup, notifySubject string, contentText string, contentHTML string) error {
//	if !mme.MailManager.GetIsActive() {
//		return fmt.Errorf("inactive mail manager")
//	}
//	var mergedMAG *MailAddressGroup
//	if mag == nil {
//		mergedMAG = mme.MailManager.MAGS.Get(MAG_KEY_SYSTEM)
//	} else {
//		// mag is optional; merge with system default
//		mergedMAG = mme.MailManager.FromMAG(MAG_KEY_SYSTEM, mag)
//	}
//	if mergedMAG == nil {
//		return fmt.Errorf("failed to get system MAG")
//	}
//	notifySubject = strings.TrimSpace(notifySubject)
//	if notifySubject == "" {
//		notifySubject = "General"
//	}
//	subjectMerge := []interface{}{fmt.Sprintf("%s, %s", notifySubject, time.Now().Format("2006-01-02 15:04:05"))}
//	bodyData := NewEmailPropertyDataContent(nil, contentText, contentHTML)
//	if err := mme.MailManager.SendWithRender(MAILTEMPLATE_SYSTEM_NOTIFY, mergedMAG, subjectMerge, bodyData); err != nil {
//		return err
//	}
//	return nil
//}
//
//// SendEmailSystem sends a general system email using the specified template.
//// It merges the provided MAG with the system default if applicable.
//func (mme *MailManagerExtend) SendEmailSystem(templateName MailTemplateName, mag *MailAddressGroup, subjectMerge []interface{}, dataBody interface{}) error {
//	if !mme.MailManager.GetIsActive() {
//		return fmt.Errorf("inactive mail manager")
//	}
//	if templateName.IsEmpty() {
//		return fmt.Errorf("template name is empty")
//	}
//	var mergedMAG *MailAddressGroup
//	if mag == nil {
//		mergedMAG = mme.MailManager.MAGS.Get(MAG_KEY_SYSTEM)
//	} else {
//		// mag is optional; merge with system default
//		mergedMAG = mme.MailManager.FromMAG(MAG_KEY_SYSTEM, mag)
//	}
//	if mergedMAG == nil {
//		return fmt.Errorf("failed to get system MAG")
//	}
//	if err := mme.MailManager.SendWithRender(templateName, mergedMAG, subjectMerge, dataBody); err != nil {
//		return err
//	}
//	return nil
//}
//
//// SendEmailUser sends a general user email using the specified template.
//// It merges the provided MAG with the users default if applicable.
//func (mme *MailManagerExtend) SendEmailUser(templateName MailTemplateName, mag *MailAddressGroup, subjectMerge []interface{}, dataBody interface{}) error {
//	if !mme.MailManager.GetIsActive() {
//		return fmt.Errorf("inactive mail manager")
//	}
//	if templateName.IsEmpty() {
//		return fmt.Errorf("template name is empty")
//	}
//	var mergedMAG *MailAddressGroup
//	if mag == nil {
//		mergedMAG = mme.MailManager.MAGS.Get(MAG_KEY_USERS)
//	} else {
//		// mag is optional; merge with users default
//		mergedMAG = mme.MailManager.FromMAG(MAG_KEY_USERS, mag)
//	}
//	if mergedMAG == nil {
//		return fmt.Errorf("failed to get users MAG")
//	}
//	if err := mme.MailManager.SendWithRender(templateName, mergedMAG, subjectMerge, dataBody); err != nil {
//		return err
//	}
//	return nil
//}
//
//// SendEmailSystemAsync sends a system email asynchronously using the specified template.
//// It merges the provided MAG with the system default if applicable.
//func (mme *MailManagerExtend) SendEmailSystemAsync(templateName MailTemplateName, mag *MailAddressGroup, subjectMerge []interface{}, dataBody interface{}, fnCallback FNMailSendCallback) {
//	var mergedMAG *MailAddressGroup
//	if mag == nil {
//		mergedMAG = mme.MailManager.MAGS.Get(MAG_KEY_SYSTEM)
//	} else {
//		mergedMAG = mme.MailManager.FromMAG(MAG_KEY_SYSTEM, mag)
//	}
//	mme.MailManager.SendWithRenderAsync(templateName, mergedMAG, subjectMerge, dataBody, fnCallback)
//}
//
//// SendEmailUserAsync sends a user email asynchronously using the specified template.
//// It merges the provided MAG with the users default if applicable.
//func (mme *MailManagerExtend) SendEmailUserAsync(templateName MailTemplateName, mag *MailAddressGroup, subjectMerge []interface{}, dataBody interface{}, fnCallback FNMailSendCallback) {
//	var mergedMAG *MailAddressGroup
//	if mag == nil {
//		mergedMAG = mme.MailManager.MAGS.Get(MAG_KEY_USERS)
//	} else {
//		mergedMAG = mme.MailManager.FromMAG(MAG_KEY_USERS, mag)
//	}
//	mme.MailManager.SendWithRenderAsync(templateName, mergedMAG, subjectMerge, dataBody, fnCallback)
//}
