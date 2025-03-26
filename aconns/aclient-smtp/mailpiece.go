package aclient_smtp

import (
	"fmt"
	"github.com/jhillyerd/enmime"
	"net/mail"
	"strings"
)

// MailPiece represents an email with various fields
type MailPiece struct {
	From mail.Address `json:"from,omitempty"`

	To  []mail.Address `json:"to,omitempty"`
	CC  []mail.Address `json:"cc,omitempty"`
	BCC []mail.Address `json:"bcc,omitempty"`

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
	if _, err = mail.ParseAddress(mp.From.Address); err != nil {
		return bdr, fmt.Errorf("invalid from address; %v", err)
	}
	bdr = bdr.From(mp.From.Name, mp.From.Address)

	// Validate TO, CC, BCC addresses
	var canSend bool

	if mp.To != nil && len(mp.To) > 0 {
		for _, addr := range mp.To {
			if _, err = mail.ParseAddress(addr.Address); err != nil {
				return bdr, fmt.Errorf("invalid TO address '%s'; %v", addr.Address, err)
			}
		}
		bdr = bdr.ToAddrs(mp.To)
		canSend = true
	}

	if mp.CC != nil && len(mp.CC) > 0 {
		for _, addr := range mp.CC {
			if _, err := mail.ParseAddress(addr.Address); err != nil {
				return bdr, fmt.Errorf("invalid CC address '%s'; %v", addr.Address, err)
			}
		}
		bdr = bdr.CCAddrs(mp.CC)
		canSend = true
	}

	if mp.BCC != nil && len(mp.BCC) > 0 {
		for _, addr := range mp.BCC {
			if _, err := mail.ParseAddress(addr.Address); err != nil {
				return bdr, fmt.Errorf("invalid BCC address '%s'; %v", addr.Address, err)
			}
		}
		bdr = bdr.BCCAddrs(mp.BCC)
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
	// Create a new MailPiece and copy over simple fields
	clone := &MailPiece{
		From:              mp.From,
		AllowEmptySubject: mp.AllowEmptySubject,
		Subject:           mp.Subject,
		AllowNoTextHMTL:   mp.AllowNoTextHMTL,
		Text:              mp.Text,
		HTML:              mp.HTML,
		AllowAttachments:  mp.AllowAttachments,
	}

	// Deep copy To, CC, and BCC slices only if they're not nil
	if mp.To != nil && len(mp.To) > 0 {
		clone.To = make([]mail.Address, len(mp.To))
		copy(clone.To, mp.To)
	}

	if mp.CC != nil && len(mp.CC) > 0 {
		clone.CC = make([]mail.Address, len(mp.CC))
		copy(clone.CC, mp.CC)
	}

	if mp.BCC != nil && len(mp.BCC) > 0 {
		clone.BCC = make([]mail.Address, len(mp.BCC))
		copy(clone.BCC, mp.BCC)
	}

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
