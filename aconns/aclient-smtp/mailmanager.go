package aclient_smtp

import (
	"fmt"

	"github.com/jpfluger/alibs-slim/aconns"
)

const MAIL_MANAGER_SMTP_DEFAULT aconns.AdapterName = "default"

// MailManager is responsible for managing SMTP connections and mail sender groups.
type MailManager struct {
	IsActive        bool                `json:"isActive"`
	SMTPs           AClientSMTPs        `json:"smtps,omitempty"`
	Templates       MailTemplates       `json:"templates,omitempty"`
	MAGS            MailAddressGroupMap `json:"mags,omitempty"`
	mailTemplateMap MailTemplateMap
}

// Validate checks the integrity of the MailManager instance.
// It ensures that the MailManager itself, SMTP connections, and mail sender groups are valid.
// It associates the smtp connection with the template.
// Returns an error if any validation step fails.
// This is intended to be called once during initialization in a single-threaded context.
func (mm *MailManager) Validate() error {
	// Check if the MailManager instance is nil.
	if mm == nil {
		return fmt.Errorf("mail manager is nil")
	}

	if !mm.IsActive {
		return nil
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

	// Validate mail address groups.
	if err := mm.MAGS.Validate(); err != nil {
		return fmt.Errorf("failed to validate mail address groups: %w", err)
	}
	if mag := mm.MAGS.Get(MAG_KEY_SYSTEM); mag == nil {
		return fmt.Errorf("mail address groups must include a '%s' key", MAG_KEY_SYSTEM)
	} else if !mag.HasAnyRecipients() {
		return fmt.Errorf("mag '%s' must include at least one recipient (To, CC, or BCC)", MAG_KEY_SYSTEM)
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

// EnsureTemplatesReady checks if all templates' snippets are loaded.
// This should be called after template loading (e.g., in activation phase).
func (mm *MailManager) EnsureTemplatesReady() error {
	for _, template := range mm.Templates {
		if err := template.CheckSnippetsLoaded(); err != nil {
			return fmt.Errorf("template %s failed snippet load check: %v", template.Name, err)
		}
	}
	return nil
}

func (mm *MailManager) GetIsActive() bool {
	if mm == nil {
		return false
	}
	return mm.IsActive
}

// FindTemplate returns a pointer to the MailTemplate with the given name.
// Returns nil if not found.
func (mm *MailManager) FindTemplate(name MailTemplateName) *MailTemplate {
	if mm == nil {
		return nil
	}
	return mm.mailTemplateMap[name]
}

// FromMAG creates a new MailAddressGroup by merging the default MAG for the given key into mergeMAG.
// If the key is empty, it defaults to MAG_KEY_SYSTEM.
// This avoids mutating the input and prevents shared slice issues via deep copying.
func (mm *MailManager) FromMAG(key MailAddressGroupKey, mergeMAG *MailAddressGroup) *MailAddressGroup {
	if key.IsEmpty() {
		key = MAG_KEY_SYSTEM
	}
	defaultMAG := mm.MAGS.Get(key)
	return FromMAG(defaultMAG, mergeMAG)
}

// SendWithRenderMAGKey sends using a default MAG as indicated by its key.
// It defaults to using the MAG_KEY_SYSTEM group if no specific key is provided.
func (mm *MailManager) SendWithRenderMAGKey(templateName MailTemplateName, magKey MailAddressGroupKey, subjectMerge []interface{}, dataBody interface{}) error {
	return mm.SendWithRenderMAGKeyOptions(templateName, magKey, nil, subjectMerge, dataBody)
}

// SendWithRenderMAGKeyOptions sends using a default MAG as indicated by its key.
// It defaults to using the MAG_KEY_SYSTEM group if no specific key is provided.
func (mm *MailManager) SendWithRenderMAGKeyOptions(templateName MailTemplateName, magKey MailAddressGroupKey, smtpAuth ISMTPAuth, subjectMerge []interface{}, dataBody interface{}) error {
	if magKey.IsEmpty() {
		magKey = MAG_KEY_SYSTEM
	}
	return mm.sendWithRenderOptions(templateName, smtpAuth, mm.FromMAG(magKey, nil), subjectMerge, dataBody)
}

// SendWithRender sends using a merged address group without mutating inputs.
// It defaults to using the MAG_KEY_SYSTEM group if no specific key is provided.
// Callers should pass mm.FromMAG(theirKey, theirMAG) as mergeAddressGroup.
func (mm *MailManager) SendWithRender(templateName MailTemplateName, mergeAddressGroup *MailAddressGroup, subjectMerge []interface{}, dataBody interface{}) error {
	return mm.sendWithRenderOptions(templateName, nil, mergeAddressGroup, subjectMerge, dataBody)
}

// sendWithRenderOptions is the private implementation for sending without locks.
// It assumes the MailManager is immutable post-validation and safe for concurrent reads.
func (mm *MailManager) sendWithRenderOptions(templateName MailTemplateName, smtpAuth ISMTPAuth, mergeAddressGroup *MailAddressGroup, subjectMerge []interface{}, dataBody interface{}) error {
	if !mm.IsActive {
		return fmt.Errorf("mail manager is inactive")
	}
	template := mm.FindTemplate(templateName)
	if template == nil {
		return fmt.Errorf("mail template '%s' not found", templateName)
	}
	if mergeAddressGroup == nil {
		return fmt.Errorf("mergeAddressGroup must not be nil; use FromMAG to create one with defaults")
	}
	// No merge here—assume caller has already applied defaults via FromMAG if desired.
	return template.SendWithRenderOptions(smtpAuth, mergeAddressGroup, subjectMerge, dataBody)
}

// FNMailSendCallback is a function type for handling the result of an asynchronous send.
// It receives the error (nil on success) after the send attempt completes.
type FNMailSendCallback func(err error)

// SendWithRenderAsync sends the email asynchronously without waiting for completion.
// If callback is provided, it is called with the result (err or nil); otherwise, errors are logged via fmt.Printf.
// For production, replace fmt with a proper logging system.
func (mm *MailManager) SendWithRenderAsync(templateName MailTemplateName, mergeAddressGroup *MailAddressGroup, subjectMerge []interface{}, dataBody interface{}, fnCallback FNMailSendCallback) {
	go func() {
		err := mm.SendWithRender(templateName, mergeAddressGroup, subjectMerge, dataBody)
		if fnCallback != nil {
			fnCallback(err)
		} else if err != nil {
			// Handle error, e.g., log it
			// For production, integrate with a logging system
			fmt.Printf("Async send error: %v\n", err)
		}
	}()
}

// SendWithRenderMAGKeyAsync sends asynchronously using a default MAG as indicated by its key.
// If callback is provided, it is called with the result (err or nil); otherwise, errors are logged via fmt.Printf.
// For production, replace fmt with a proper logging system.
func (mm *MailManager) SendWithRenderMAGKeyAsync(templateName MailTemplateName, magKey MailAddressGroupKey, subjectMerge []interface{}, dataBody interface{}, fnCallback FNMailSendCallback) {
	go func() {
		err := mm.SendWithRenderMAGKey(templateName, magKey, subjectMerge, dataBody)
		if fnCallback != nil {
			fnCallback(err)
		} else if err != nil {
			// Handle error, e.g., log it
			fmt.Printf("Async send error: %v\n", err)
		}
	}()
}

//package aclient_smtp
//
//import (
//	"fmt"
//
//	"github.com/jpfluger/alibs-slim/aconns"
//)
//
//const MAIL_MANAGER_SMTP_DEFAULT aconns.AdapterName = "default"
//
//// MailManager is responsible for managing SMTP connections and mail sender groups.
//type MailManager struct {
//	IsActive        bool                `json:"isActive"`
//	SMTPs           AClientSMTPs        `json:"smtps,omitempty"`
//	Templates       MailTemplates       `json:"templates,omitempty"`
//	MAGS            MailAddressGroupMap `json:"mags,omitempty"`
//	mailTemplateMap MailTemplateMap
//}
//
//// Validate checks the integrity of the MailManager instance.
//// It ensures that the MailManager itself, SMTP connections, and mail sender groups are valid.
//// It associates the smtp connection with the template.
//// Returns an error if any validation step fails.
//// This is intended to be called once during initialization in a single-threaded context.
//func (mm *MailManager) Validate() error {
//	// Check if the MailManager instance is nil.
//	if mm == nil {
//		return fmt.Errorf("mail manager is nil")
//	}
//
//	if !mm.IsActive {
//		return nil
//	}
//
//	// Ensure SMTP connections and sender groups are not empty.
//	if mm.SMTPs == nil || len(mm.SMTPs) == 0 {
//		return fmt.Errorf("mail manager has no SMTP connections")
//	}
//
//	if len(mm.Templates) == 0 {
//		return fmt.Errorf("mail manager has no templates")
//	}
//
//	var defaultAdapter *AClientSMTP
//
//	// Validate each SMTP connection in the Conns collection.
//	for ii, smtp := range mm.SMTPs {
//		if smtp == nil {
//			return fmt.Errorf("SMTP connection at index %d is nil", ii)
//		}
//		if err := smtp.Validate(); err != nil {
//			return fmt.Errorf("failed to validate SMTP connection at index %d: %v", ii, err)
//		}
//		if smtp.GetName() == MAIL_MANAGER_SMTP_DEFAULT {
//			if defaultAdapter != nil {
//				return fmt.Errorf("SMTP 'default' connection at index %d: duplicate initialization", ii)
//			}
//			defaultAdapter = smtp
//		}
//	}
//
//	if defaultAdapter == nil {
//		return fmt.Errorf("smtp connection does not have an SMTP adapter named 'default'")
//	}
//
//	// Validate mail address groups.
//	if err := mm.MAGS.Validate(); err != nil {
//		return fmt.Errorf("failed to validate mail address groups: %w", err)
//	}
//	if mag := mm.MAGS.Get(MAG_KEY_SYSTEM); mag == nil {
//		return fmt.Errorf("mail address groups must include a '%s' key", MAG_KEY_SYSTEM)
//	} else if !mag.HasAnyRecipients() {
//		return fmt.Errorf("mag '%s' must include at least one recipient (To, CC, or BCC)", MAG_KEY_SYSTEM)
//	}
//
//	// Perform pre-validation on the mail sender groups.
//	for ii, template := range mm.Templates {
//		if err := template.PreValidate(); err != nil {
//			return fmt.Errorf("mail template #%d with pre-validate error: %v", ii, err)
//		}
//		if template.SmtpOverride.IsEmpty() {
//			template.smtpAuth = defaultAdapter
//		} else {
//			smtp := mm.SMTPs.FindByName(template.SmtpOverride)
//			if smtp == nil {
//				return fmt.Errorf("failed to locate SMTP Override '%s' at template index %d", template.SmtpOverride.String(), ii)
//			}
//			template.smtpAuth = smtp
//		}
//	}
//
//	mm.mailTemplateMap = mm.Templates.ToMap()
//
//	return nil
//}
//
//func (mm *MailManager) GetIsActive() bool {
//	if mm == nil {
//		return false
//	}
//	return mm.IsActive
//}
//
//func (mm *MailManager) FindTemplate(templateName MailTemplateName) *MailTemplate {
//	if template, ok := mm.mailTemplateMap[templateName]; ok {
//		return template
//	}
//	return nil
//}
//
//// FromMAG creates a new MailAddressGroup by merging the default MAG for the given key into mergeMAG.
//// If the key is empty, it defaults to MAG_KEY_SYSTEM.
//// This avoids mutating the input and prevents shared slice issues via deep copying.
//func (mm *MailManager) FromMAG(key MailAddressGroupKey, mergeMAG *MailAddressGroup) *MailAddressGroup {
//	if key.IsEmpty() {
//		key = MAG_KEY_SYSTEM
//	}
//	defaultMAG := mm.MAGS.Get(key)
//	return FromMAG(defaultMAG, mergeMAG)
//}
//
//// SendWithRenderMAGKey sends using a default MAG as indicated by its key.
//// It directly uses the MAG without merging; for merging, use SendWithRender with FromMAG.
//func (mm *MailManager) SendWithRenderMAGKey(templateName MailTemplateName, magKey MailAddressGroupKey, subjectMerge []interface{}, dataBody interface{}) error {
//	if magKey.IsEmpty() {
//		magKey = MAG_KEY_SYSTEM
//	}
//	mag := mm.MAGS.Get(magKey)
//	if mag == nil {
//		return fmt.Errorf("MAG not found by key '%s'", magKey)
//	}
//	return mm.sendWithRenderOptions(templateName, nil, mag, subjectMerge, dataBody)
//}
//
//// SendWithRender sends using a merged address group without mutating inputs.
//// It defaults to using the MAG_KEY_SYSTEM group if no specific key is provided.
//// Callers should pass mm.FromMAG(theirKey, theirMAG) as mergeAddressGroup.
//func (mm *MailManager) SendWithRender(templateName MailTemplateName, mergeAddressGroup *MailAddressGroup, subjectMerge []interface{}, dataBody interface{}) error {
//	return mm.sendWithRenderOptions(templateName, nil, mergeAddressGroup, subjectMerge, dataBody)
//}
//
//// sendWithRenderOptions is the private implementation for sending without locks.
//// It assumes the MailManager is immutable post-validation and safe for concurrent reads.
//func (mm *MailManager) sendWithRenderOptions(templateName MailTemplateName, smtpAuth ISMTPAuth, mergeAddressGroup *MailAddressGroup, subjectMerge []interface{}, dataBody interface{}) error {
//	if !mm.IsActive {
//		return fmt.Errorf("mail manager is inactive")
//	}
//	template := mm.FindTemplate(templateName)
//	if template == nil {
//		return fmt.Errorf("mail template '%s' not found", templateName)
//	}
//	if mergeAddressGroup == nil {
//		return fmt.Errorf("mergeAddressGroup must not be nil; use FromMAG to create one with defaults")
//	}
//	// No merge here—assume caller has already applied defaults via FromMAG if desired.
//	return template.SendWithRenderOptions(smtpAuth, mergeAddressGroup, subjectMerge, dataBody)
//}
//
//// FNMailSendCallback is a function type for handling the result of an asynchronous send.
//// It receives the error (nil on success) after the send attempt completes.
//type FNMailSendCallback func(err error)
//
//// SendWithRenderAsync sends the email asynchronously without waiting for completion.
//// If callback is provided, it is called with the result (err or nil); otherwise, errors are logged via fmt.Printf.
//// For production, replace fmt with a proper logging system.
//func (mm *MailManager) SendWithRenderAsync(templateName MailTemplateName, mergeAddressGroup *MailAddressGroup, subjectMerge []interface{}, dataBody interface{}, fnCallback FNMailSendCallback) {
//	go func() {
//		err := mm.SendWithRender(templateName, mergeAddressGroup, subjectMerge, dataBody)
//		if fnCallback != nil {
//			fnCallback(err)
//		} else if err != nil {
//			// Handle error, e.g., log it
//			// For production, integrate with a logging system
//			fmt.Printf("Async send error: %v\n", err)
//		}
//	}()
//}
//
//// SendWithRenderMAGKeyAsync sends asynchronously using a default MAG as indicated by its key.
//// If callback is provided, it is called with the result (err or nil); otherwise, errors are logged via fmt.Printf.
//// For production, replace fmt with a proper logging system.
//func (mm *MailManager) SendWithRenderMAGKeyAsync(templateName MailTemplateName, magKey MailAddressGroupKey, subjectMerge []interface{}, dataBody interface{}, fnCallback FNMailSendCallback) {
//	go func() {
//		err := mm.SendWithRenderMAGKey(templateName, magKey, subjectMerge, dataBody)
//		if fnCallback != nil {
//			fnCallback(err)
//		} else if err != nil {
//			// Handle error, e.g., log it
//			fmt.Printf("Async send error: %v\n", err)
//		}
//	}()
//}
