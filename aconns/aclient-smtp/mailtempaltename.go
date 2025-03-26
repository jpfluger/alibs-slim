package aclient_smtp

import "strings"

// MailTemplateName represents a mail template name
type MailTemplateName string

// IsEmpty checks if the MailTemplateName is empty
func (mtn MailTemplateName) IsEmpty() bool {
	return mtn == ""
}

// TrimSpace returns the trimmed string representation of the MailTemplateName
func (mtn MailTemplateName) TrimSpace() MailTemplateName {
	return MailTemplateName(strings.TrimSpace(string(mtn)))
}

// String returns the string representation of the MailTemplateName
func (mtn MailTemplateName) String() string {
	return string(mtn)
}

// Matches checks if the MailTemplateName matches the given string
func (mtn MailTemplateName) Matches(s string) bool {
	return string(mtn) == s
}

// MailTemplateNames represents a slice of MailTemplateName
type MailTemplateNames []MailTemplateName

// IsEmpty checks if the MailTemplateNames slice is empty
func (mtns MailTemplateNames) IsEmpty() bool {
	return len(mtns) == 0
}

// String returns the string representation of the MailTemplateNames slice
func (mtns MailTemplateNames) String() string {
	return strings.Join(mtns.ToStringArray(), ", ")
}

// ToStringArray returns an array of MailTemplateNames as strings
func (mtns MailTemplateNames) ToStringArray() []string {
	strArray := make([]string, len(mtns))
	for i, mtn := range mtns {
		strArray[i] = mtn.String()
	}
	return strArray
}

// Find returns the MailTemplateName if found, otherwise an empty MailTemplateName
func (mtns MailTemplateNames) Find(mtn MailTemplateName) MailTemplateName {
	for _, v := range mtns {
		if v == mtn {
			return v
		}
	}
	return ""
}

// HasKey checks if the MailTemplateNames slice contains the given MailTemplateName
func (mtns MailTemplateNames) HasKey(s MailTemplateName) bool {
	return mtns.Find(s) != ""
}

// Matches checks if any MailTemplateName in the MailTemplateNames slice matches the given string
func (mtns MailTemplateNames) Matches(s string) bool {
	for _, v := range mtns {
		if v.Matches(s) {
			return true
		}
	}
	return false
}
