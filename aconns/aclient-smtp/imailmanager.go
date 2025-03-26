package aclient_smtp

import "sync"

type IMailManager interface {
	FindTemplate(templateName MailTemplateName) *MailTemplate
	SendWithRender(templateName MailTemplateName, addressGroup *MailAddressGroup, subjectMerge []interface{}, dataBody interface{}) error
	SendWithRenderOptions(templateName MailTemplateName, smtpAuth ISMTPAuth, addressGroup *MailAddressGroup, subjectMerge []interface{}, dataBody interface{}) error
}

var (
	globalMailManager   IMailManager
	muGlobalMailManager = sync.RWMutex{}
)

func SetMailManager(mailManager IMailManager) {
	muGlobalMailManager.Lock()
	defer muGlobalMailManager.Unlock()
	globalMailManager = mailManager
}

func MAILMANAGER() IMailManager {
	muGlobalMailManager.RLock()
	defer muGlobalMailManager.RUnlock()
	return globalMailManager
}
