package aclient_smtp

import (
	"fmt"
	"github.com/jpfluger/alibs-slim/aconns"
	"github.com/jpfluger/alibs-slim/atemplates"
	"strings"
)

// MailTemplate represents an individual email template.
type MailTemplate struct {
	Name MailTemplateName `json:"name,omitempty"` // The name of the mail template.

	SmtpOverride aconns.AdapterName `json:"smtpOverride,omitempty"` // Optional override of the SMTP adapter

	MailPiece

	SubjectMerge    string `json:"subjectMerge,omitempty"`    // Format string for subject merging.
	SnippetTextName string `json:"snippetTextName,omitempty"` // Name of the text snippet.
	SnippetHTMLName string `json:"snippetHTMLName,omitempty"` // Name of the HTML snippet.

	isPreValidate   bool // Flag indicating whether the template has been pre-validated.
	hasSubjectMerge bool // Flag indicating if SubjectMerge is set.
	hasSnippetText  bool // Flag indicating if SnippetTextName is set.
	hasSnippetHTML  bool // Flag indicating if SnippetHTMLName is set.

	smtpAuth ISMTPAuth
}

// PreValidate checks the validity of the MailTemplate.
// It trims unnecessary spaces and ensures required snippets are loaded.
// Returns an error if validation fails.
func (mt *MailTemplate) PreValidate() error {
	if mt == nil {
		return fmt.Errorf("mail template is nil")
	}

	mt.SubjectMerge = strings.TrimSpace(mt.SubjectMerge)
	if mt.SubjectMerge != "" {
		mt.hasSubjectMerge = true
	}

	mt.SnippetTextName = strings.TrimSpace(mt.SnippetTextName)
	mt.SnippetHTMLName = strings.TrimSpace(mt.SnippetHTMLName)

	if mt.SnippetTextName != "" {
		if err := atemplates.TEMPLATES().SnippetsText.IsLoaded(mt.SnippetTextName); err != nil {
			return err
		}
		mt.hasSnippetText = true
	}

	if mt.SnippetHTMLName != "" {
		if err := atemplates.TEMPLATES().SnippetsHTML.IsLoaded(mt.SnippetHTMLName); err != nil {
			return err
		}
		mt.hasSnippetHTML = true
	}

	mt.isPreValidate = true

	return nil
}

// SendWithRender prepares and sends an email using the MailTemplate.
func (mt *MailTemplate) SendWithRender(addressGroup *MailAddressGroup, subjectMerge []interface{}, dataBody interface{}) error {
	return mt.SendWithRenderOptions(nil, addressGroup, subjectMerge, dataBody)
}

// SendWithRenderOptions prepares and sends an email using the MailTemplate.
// It renders the subject and body with provided data and sends the email via SMTP authentication.
func (mt *MailTemplate) SendWithRenderOptions(smtpAuth ISMTPAuth, addressGroup *MailAddressGroup, subjectMerge []interface{}, dataBody interface{}) error {
	if mt == nil {
		return fmt.Errorf("mail template is nil")
	}

	if !mt.isPreValidate {
		if err := mt.PreValidate(); err != nil {
			return fmt.Errorf("failed pre-validate mail template: %w", err)
		}
	}

	if smtpAuth == nil {
		if mt.smtpAuth != nil {
			smtpAuth = mt.smtpAuth
		} else {
			return fmt.Errorf("smtpAuth is nil")
		}
	}

	// Clone the template to avoid modifying the original instance.
	clone := mt.Clone()

	// Merge address group if provided.
	if addressGroup != nil {
		addressGroup.MergeIntoTemplate(clone)
	}

	// Render the subject if SubjectMerge is set.
	if mt.hasSubjectMerge {
		// fmt.Sprintf will safely handle a nil or empty subjectMerge slice.
		clone.Subject = fmt.Sprintf(mt.SubjectMerge, subjectMerge...)
	}

	// Render the text snippet if specified.
	if mt.hasSnippetText {
		output, err := atemplates.TEMPLATES().SnippetsText.RenderSnippet(mt.SnippetTextName, dataBody)
		if err != nil {
			return fmt.Errorf("failed render snippet text: %w", err)
		}
		clone.Text = output
	}

	// Render the HTML snippet if specified.
	if mt.hasSnippetHTML {
		output, err := atemplates.TEMPLATES().SnippetsHTML.RenderSnippet(mt.SnippetHTMLName, dataBody)
		if err != nil {
			return fmt.Errorf("failed render snippet html: %w", err)
		}
		clone.HTML = output
	}

	// Send the email.
	return clone.Send(smtpAuth)
}

// MailTemplates represents a collection of MailTemplate pointers.
type MailTemplates []*MailTemplate

// FindByName searches for a MailTemplate by its name.
// Returns a pointer to the template if found, or nil otherwise.
func (mts MailTemplates) FindByName(name MailTemplateName) *MailTemplate {
	if name.IsEmpty() {
		return nil
	}
	for _, mt := range mts {
		if mt.Name == name {
			return mt
		}
	}
	return nil
}

// ToMap converts the array to a map of MailTemplates using MailTemplateName as the key.
func (mts MailTemplates) ToMap() MailTemplateMap {
	mmap := make(MailTemplateMap, len(mts))
	for _, mt := range mts {
		mmap[mt.Name] = mt
	}
	return mmap
}

type MailTemplateMap map[MailTemplateName]*MailTemplate
