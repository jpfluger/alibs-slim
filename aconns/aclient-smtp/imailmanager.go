package aclient_smtp

import "sync"

type IMailManager interface {
	Validate() error
	GetIsActive() bool
	FindTemplate(templateName MailTemplateName) *MailTemplate
	FromMAG(key MailAddressGroupKey, mergeMAG *MailAddressGroup) *MailAddressGroup
	SendWithRenderMAGKey(templateName MailTemplateName, magKey MailAddressGroupKey, subjectMerge []interface{}, dataBody interface{}) error
	SendWithRender(templateName MailTemplateName, mergeAddressGroup *MailAddressGroup, subjectMerge []interface{}, dataBody interface{}) error
	SendWithRenderAsync(templateName MailTemplateName, mergeAddressGroup *MailAddressGroup, subjectMerge []interface{}, dataBody interface{}, fnCallback FNMailSendCallback)
	SendWithRenderMAGKeyAsync(templateName MailTemplateName, magKey MailAddressGroupKey, subjectMerge []interface{}, dataBody interface{}, fnCallback FNMailSendCallback)
}

var (
	globalMailManager   IMailManager
	muGlobalMailManager = sync.RWMutex{}
)

func MAILMANAGER() IMailManager {
	muGlobalMailManager.RLock()
	defer muGlobalMailManager.RUnlock()
	return globalMailManager
}

func SetMailManager(mailManager IMailManager) {
	muGlobalMailManager.Lock()
	defer muGlobalMailManager.Unlock()
	globalMailManager = mailManager
}
