package aclient_smtp

import "net/mail"

// MailAddressGroup represents a group of email addresses, including From, To, CC, and BCC fields.
type MailAddressGroup struct {
	From mail.Address   `json:"from,omitempty"`
	To   []mail.Address `json:"to,omitempty"`
	CC   []mail.Address `json:"cc,omitempty"`
	BCC  []mail.Address `json:"bcc,omitempty"`
}

// MergeInto merges the non-empty properties of the source (mag) into the target (mergeTo)
// only if the target properties are empty.
func (mag *MailAddressGroup) MergeInto(mergeTo *MailAddressGroup) {
	// Merge From
	if mag.From.Address != "" && mergeTo.From.Address == "" {
		mergeTo.From = mag.From
	}

	// Merge To
	if len(mergeTo.To) == 0 && len(mag.To) > 0 {
		mergeTo.To = mag.To
	}

	// Merge CC
	if len(mergeTo.CC) == 0 && len(mag.CC) > 0 {
		mergeTo.CC = mag.CC
	}

	// Merge BCC
	if len(mergeTo.BCC) == 0 && len(mag.BCC) > 0 {
		mergeTo.BCC = mag.BCC
	}
}

// MergeIntoTemplate merges non-empty values from a MailAddressGroup into a MailPiece.
// This ensures that only fields with data in the MailAddressGroup overwrite corresponding fields in MailPiece.
func (mag *MailAddressGroup) MergeIntoTemplate(template *MailPiece) {
	if mag == nil {
		return
	}
	if mag.From.Address != "" {
		template.From = mag.From
	}

	var mustMerge bool

	if mag.To != nil && len(mag.To) > 0 {
		mustMerge = true
	} else if mag.CC != nil && len(mag.CC) > 0 {
		mustMerge = true
	} else if mag.BCC != nil && len(mag.BCC) > 0 {
		mustMerge = true
	}

	if mustMerge {
		template.To = mag.To
		template.CC = mag.CC
		template.BCC = mag.BCC
	}
}

// NewMAG creates a new MailAddressGroup with all fields provided.
func NewMAG(from mail.Address, to []mail.Address, cc []mail.Address, bcc []mail.Address) *MailAddressGroup {
	return &MailAddressGroup{
		From: from,
		To:   to,
		CC:   cc,
		BCC:  bcc,
	}
}

// NewMAGWithTo creates a new MailAddressGroup with a single "To" address.
func NewMAGWithTo(from mail.Address, to mail.Address) *MailAddressGroup {
	return &MailAddressGroup{
		From: from,
		To:   []mail.Address{to},
	}
}

// NewMAGWithCC creates a new MailAddressGroup with "From" and "CC" addresses.
func NewMAGWithCC(from mail.Address, cc []mail.Address) *MailAddressGroup {
	return &MailAddressGroup{
		From: from,
		CC:   cc,
	}
}

// NewMAGWithBCC creates a new MailAddressGroup with "From" and "BCC" addresses.
func NewMAGWithBCC(from mail.Address, bcc []mail.Address) *MailAddressGroup {
	return &MailAddressGroup{
		From: from,
		BCC:  bcc,
	}
}

// NewMAGWithToAndCC creates a new MailAddressGroup with a single "To" address and "CC" addresses.
func NewMAGWithToAndCC(from mail.Address, to mail.Address, cc []mail.Address) *MailAddressGroup {
	return &MailAddressGroup{
		From: from,
		To:   []mail.Address{to},
		CC:   cc,
	}
}

// NewMAGWithToAndBCC creates a new MailAddressGroup with a single "To" address and "BCC" addresses.
func NewMAGWithToAndBCC(from mail.Address, to mail.Address, bcc []mail.Address) *MailAddressGroup {
	return &MailAddressGroup{
		From: from,
		To:   []mail.Address{to},
		BCC:  bcc,
	}
}

// NewMAGWithSingleToAndCCBCC creates a new MailAddressGroup with a single "To", "CC", and "BCC" addresses.
func NewMAGWithSingleToAndCCBCC(from mail.Address, to mail.Address, cc []mail.Address, bcc []mail.Address) *MailAddressGroup {
	return &MailAddressGroup{
		From: from,
		To:   []mail.Address{to},
		CC:   cc,
		BCC:  bcc,
	}
}

// NewMailAddress creates a new mail.Address with only the Address string.
// The Name field remains empty.
func NewMailAddress(email string) mail.Address {
	return mail.Address{
		Name:    "",
		Address: email,
	}
}

// NewMAGWithToString creates a MailAddressGroup with a single "To" email address string.
func NewMAGWithToString(from string, to string) *MailAddressGroup {
	return &MailAddressGroup{
		From: NewMailAddress(from),
		To:   []mail.Address{NewMailAddress(to)},
	}
}

// NewMAGWithCCString creates a MailAddressGroup with "From" and "CC" email address strings.
func NewMAGWithCCString(from string, cc []string) *MailAddressGroup {
	return &MailAddressGroup{
		From: NewMailAddress(from),
		CC:   convertStringsToAddresses(cc),
	}
}

// NewMAGWithBCCString creates a MailAddressGroup with "From" and "BCC" email address strings.
func NewMAGWithBCCString(from string, bcc []string) *MailAddressGroup {
	return &MailAddressGroup{
		From: NewMailAddress(from),
		BCC:  convertStringsToAddresses(bcc),
	}
}

// NewMAGWithToCCBCCString creates a MailAddressGroup with "To", "CC", and "BCC" email address strings.
func NewMAGWithToCCBCCString(from string, to string, cc []string, bcc []string) *MailAddressGroup {
	return &MailAddressGroup{
		From: NewMailAddress(from),
		To:   []mail.Address{NewMailAddress(to)},
		CC:   convertStringsToAddresses(cc),
		BCC:  convertStringsToAddresses(bcc),
	}
}

// Helper function to convert a slice of email strings into a slice of mail.Address.
func convertStringsToAddresses(emails []string) []mail.Address {
	addresses := make([]mail.Address, len(emails))
	for i, email := range emails {
		addresses[i] = NewMailAddress(email)
	}
	return addresses
}
