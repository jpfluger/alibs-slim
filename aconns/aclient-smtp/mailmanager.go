package aclient_smtp

import (
	"fmt"
	"github.com/jpfluger/alibs-slim/aconns"
)

const MAIL_MANAGER_SMTP_DEFAULT aconns.AdapterName = "default"

// MailManager is responsible for managing SMTP connections and mail sender groups.
type MailManager struct {
	SMTPs           AClientSMTPs     `json:"smtps,omitempty"`
	Templates       MailTemplates    `json:"templates,omitempty"`
	MAGDefaults     MailAddressGroup `json:"magDefaults,omitempty"`
	mailTemplateMap MailTemplateMap
}

// Validate checks the integrity of the MailManager instance.
// It ensures that the MailManager itself, SMTP connections, and mail sender groups are valid.
// It associates the smtp connection with the template.
// Returns an error if any validation step fails.
func (mm *MailManager) Validate() error {
	// Check if the MailManager instance is nil.
	if mm == nil {
		return fmt.Errorf("mail manager is nil")
	}

	// Ensure SMTP connections and sender groups are not empty.
	if mm.SMTPs == nil || len(mm.SMTPs) == 0 {
		return fmt.Errorf("mail manager has no SMTP connections")
	}

	if len(mm.Templates) == 0 {
		return fmt.Errorf("mail manager has no templates")
	}

	var defaultAdapter *AClientSMTP

	// Validate each SMTP connection in the Conns collection.
	for ii, smtp := range mm.SMTPs {
		if smtp == nil {
			return fmt.Errorf("SMTP connection at index %d is nil", ii)
		}
		if err := smtp.Validate(); err != nil {
			return fmt.Errorf("failed to validate SMTP connection at index %d: %v", ii, err)
		}
		if smtp.GetName() == MAIL_MANAGER_SMTP_DEFAULT {
			if defaultAdapter != nil {
				return fmt.Errorf("SMTP 'default' connection at index %d: duplicate initialization", ii)
			}
			defaultAdapter = smtp
		}
	}

	if defaultAdapter == nil {
		return fmt.Errorf("smtp connection does not have an SMTP adapter named 'default'")
	}

	// Perform pre-validation on the mail sender groups.
	for ii, template := range mm.Templates {
		if err := template.PreValidate(); err != nil {
			return fmt.Errorf("mail template #%d with pre-validate error: %v", ii, err)
		}
		if template.SmtpOverride.IsEmpty() {
			template.smtpAuth = defaultAdapter
		} else {
			smtp := mm.SMTPs.FindByName(template.SmtpOverride)
			if smtp == nil {
				return fmt.Errorf("failed to locate SMTP Override '%s' at template index %d", template.SmtpOverride.String(), ii)
			}
			template.smtpAuth = smtp
		}
	}

	mm.mailTemplateMap = mm.Templates.ToMap()

	return nil
}

func (mm *MailManager) FindTemplate(templateName MailTemplateName) *MailTemplate {
	if template, ok := mm.mailTemplateMap[templateName]; ok {
		return template
	}
	return nil
}

func (mm *MailManager) SendWithRender(templateName MailTemplateName, mergeAddressGroup *MailAddressGroup, subjectMerge []interface{}, dataBody interface{}) error {
	return mm.SendWithRenderOptions(templateName, nil, mergeAddressGroup, subjectMerge, dataBody)
}

func (mm *MailManager) SendWithRenderOptions(templateName MailTemplateName, smtpAuth ISMTPAuth, mergeAddressGroup *MailAddressGroup, subjectMerge []interface{}, dataBody interface{}) error {
	template := mm.FindTemplate(templateName)
	if template == nil {
		return fmt.Errorf("mail template '%s' not found", templateName)
	}
	// Ensure mergeAddressGroup is initialized
	if mergeAddressGroup == nil {
		mergeAddressGroup = &MailAddressGroup{}
	}
	// Merge MAGDefaults into mergeAddressGroup
	mm.MAGDefaults.MergeInto(mergeAddressGroup)
	// Render!
	return template.SendWithRenderOptions(smtpAuth, mergeAddressGroup, subjectMerge, dataBody)
}
