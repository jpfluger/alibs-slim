package aclient_smtp

import (
	"fmt"
	"strings"

	"github.com/jhillyerd/enmime/v2"
	"github.com/jpfluger/alibs-slim/aemail"
)

// MailPiece represents an email with various fields
type MailPiece struct {
	From aemail.Address `json:"from,omitempty"`

	To  aemail.Addresses `json:"to,omitempty"`
	CC  aemail.Addresses `json:"cc,omitempty"`
	BCC aemail.Addresses `json:"bcc,omitempty"`

	AllowEmptySubject bool   `json:"allowEmptySubject,omitempty"`
	Subject           string `json:"subject,omitempty"`

	AllowNoTextHMTL bool   `json:"allowNoTextHMTL,omitempty"`
	Text            string `json:"text,omitempty"`
	HTML            string `json:"html,omitempty"`

	AllowAttachments bool `json:"allowAttachments,omitempty"`
	// Uses Attachments instead of []*enmime.Part, to allow more
	// flexibility, efficiency and long-term storage options at the app level.
	Attachments Attachments `json:"attachments,omitempty"`
	Inlines     Attachments `json:"inlines,omitempty"`
}

// Validate checks if the MailPiece object is valid
func (mp *MailPiece) Validate() error {
	_, err := mp.validateWithBuilder()
	if err != nil {
		return err
	}
	return nil
}

// validateWithBuilder validates the MailPiece and returns an enmime.MailBuilder
func (mp *MailPiece) validateWithBuilder() (bdr enmime.MailBuilder, err error) {
	bdr = enmime.Builder()

	if mp == nil {
		return bdr, fmt.Errorf("mail piece is nil")
	}

	// Validate FROM address
	if err = mp.From.Validate(); err != nil {
		return bdr, fmt.Errorf("invalid from address; %v", err)
	}
	bdr = bdr.From(mp.From.Name, mp.From.Address.String())

	// Validate TO, CC, BCC addresses
	var canSend bool

	if mp.To != nil && len(mp.To) > 0 {
		if err = mp.To.Validate(); err != nil {
			return bdr, fmt.Errorf("invalid TO addresses; %v", err)
		}
		bdr = bdr.ToAddrs(mp.To.ToMailAddresses())
		canSend = true
	}

	if mp.CC != nil && len(mp.CC) > 0 {
		if err = mp.CC.Validate(); err != nil {
			return bdr, fmt.Errorf("invalid CC addresses; %v", err)
		}
		bdr = bdr.CCAddrs(mp.CC.ToMailAddresses())
		canSend = true
	}

	if mp.BCC != nil && len(mp.BCC) > 0 {
		if err = mp.BCC.Validate(); err != nil {
			return bdr, fmt.Errorf("invalid BCC addresses; %v", err)
		}
		bdr = bdr.BCCAddrs(mp.BCC.ToMailAddresses())
		canSend = true
	}

	if !canSend {
		return bdr, fmt.Errorf("no TO, CC, or BCC addresses found")
	}

	// Validate SUBJECT
	mp.Subject = strings.TrimSpace(mp.Subject)
	if mp.Subject == "" {
		if !mp.AllowEmptySubject {
			return bdr, fmt.Errorf("subject is empty")
		}
	} else {
		bdr = bdr.Subject(mp.Subject)
	}

	// Validate TEXT or HTML content
	hasTextOrHTML := false

	mp.Text = strings.TrimSpace(mp.Text)
	if mp.Text != "" {
		bdr = bdr.Text([]byte(mp.Text))
		hasTextOrHTML = true
	}

	mp.HTML = strings.TrimSpace(mp.HTML)
	if mp.HTML != "" {
		bdr = bdr.HTML([]byte(mp.HTML))
		hasTextOrHTML = true
	}

	if !hasTextOrHTML {
		if !mp.AllowNoTextHMTL {
			return bdr, fmt.Errorf("no text or html content found")
		}
	}

	if mp.AllowAttachments {
		// Add attachments
		if mp.Attachments != nil && len(mp.Attachments) > 0 {
			for ii, att := range mp.Attachments {
				if att.HasEnmimePart() {
					epart := att.GetEnmimePart()
					bdr = bdr.AddAttachment(epart.Content, epart.ContentType, epart.FileName)
				} else {
					if err = att.LoadContent(); err != nil {
						return bdr, fmt.Errorf("failed to load attachment at index %d: %v", ii, err)
					}
					var content []byte
					content, err = att.GetContent()
					if err != nil {
						return bdr, fmt.Errorf("failed to get attachment content at index %d: %v", ii, err)
					}
					bdr = bdr.AddAttachment(content, att.ContentType, att.Name)
				}
			}
		}

		// Add inlines
		if mp.Inlines != nil && len(mp.Inlines) > 0 {
			for ii, inl := range mp.Inlines {
				if inl.HasEnmimePart() {
					epart := inl.GetEnmimePart()
					bdr = bdr.AddInline(epart.Content, epart.ContentType, epart.FileName, epart.ContentID)
				} else {
					if err = inl.LoadContent(); err != nil {
						return bdr, fmt.Errorf("failed to load inline attachment at index %d: %v", ii, err)
					}
					var content []byte
					content, err = inl.GetContent()
					if err != nil {
						return bdr, fmt.Errorf("failed to get inline attachment content at index %d: %v", ii, err)
					}
					bdr = bdr.AddInline(content, inl.ContentType, inl.Name, "")
				}
			}
		}
	}

	return bdr, nil
}

// Send sends the email using the provided SMTP authentication
func (mp *MailPiece) Send(smtpAuth ISMTPAuth) error {
	bdr, err := mp.validateWithBuilder()
	if err != nil {
		return err
	}
	if smtpAuth == nil {
		return fmt.Errorf("smtp auth is nil")
	}
	sender, err := smtpAuth.GetSender()
	if err != nil {
		return fmt.Errorf("failed to get smtp sender; %v", err)
	}
	if err = bdr.Send(sender); err != nil {
		return fmt.Errorf("failed to send email; %v", err)
	}
	return nil
}

// Clone creates a deep copy of the MailPiece object, including all nested fields.
func (mp *MailPiece) Clone() *MailPiece {
	if mp == nil {
		return nil
	}
	// Create a new MailPiece and copy over simple fields
	clone := &MailPiece{
		From:              mp.From, // Value type, shallow copy is fine
		AllowEmptySubject: mp.AllowEmptySubject,
		Subject:           mp.Subject,
		AllowNoTextHMTL:   mp.AllowNoTextHMTL,
		Text:              mp.Text,
		HTML:              mp.HTML,
		AllowAttachments:  mp.AllowAttachments,
	}

	// Deep copy To, CC, and BCC slices
	clone.To = mp.To.Clone()
	clone.CC = mp.CC.Clone()
	clone.BCC = mp.BCC.Clone()

	// Deep copy Attachments and Inlines only if they're not nil
	if mp.Attachments != nil {
		clone.Attachments = mp.Attachments.Clone()
		if clone.Attachments == nil {
			clone.Attachments = Attachments{}
		}
	}

	if mp.Inlines != nil {
		clone.Inlines = mp.Inlines.Clone()
		if clone.Inlines == nil {
			clone.Inlines = Attachments{}
		}
	}

	return clone
}
