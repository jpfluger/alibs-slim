package acontact

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jpfluger/alibs-slim/aimage"
	"github.com/jpfluger/alibs-slim/atags"
)

// IContactCore defines the interface for ContactCore behavior.
type IContactCore interface {
	GetTypes() ContactTypes
	IsValid() bool // ContactCore Validity is based off the Name.
	GetName() string
	GetNameOptions() *Name
	IsPerson() bool
	IsEntity() bool
	IsBothPersonAndEntity() bool
	GetEmails() Emails
	GetPhones() Phones
	GetMails() Mails
	GetUrls() Urls
	GetDates() atags.TagArrTimes
}

// ContactCore represents the core structure for a contact.
type ContactCore struct {
	Name    Name                `json:"name,omitempty"`
	Emails  Emails              `json:"emails,omitempty"` // Collection of email addresses
	Phones  Phones              `json:"phones,omitempty"` // Collection of phone numbers
	Mails   Mails               `json:"mails,omitempty"`  // Collection of mailing addresses
	Urls    Urls                `json:"urls,omitempty"`   // Collection of URLs
	Dates   atags.TagArrTimes   `json:"dates,omitempty"`
	Images  aimage.Images       `json:"images,omitempty"`
	Tags    atags.TagArrStrings `json:"tags,omitempty"`
	Details json.RawMessage     `json:"details,omitempty"`
}

// Validate ensures that the ContactCore is valid.
func (c *ContactCore) Validate() error {
	if c == nil {
		return errors.New("contact core is nil")
	}
	if err := c.Name.Validate(); err != nil {
		return fmt.Errorf("name validation failed: %w", err)
	}
	return nil
}

func (c *ContactCore) GetTypes() ContactTypes {
	if c.IsBothPersonAndEntity() {
		return ContactTypes{CONTACTTYPE_ENTITY, CONTACTTYPE_PERSON}
	} else if c.IsPerson() {
		return ContactTypes{CONTACTTYPE_PERSON}
	} else if c.IsEntity() {
		return ContactTypes{CONTACTTYPE_ENTITY}
	}
	return ContactTypes{}
}

// IsValid checks to ensure at least the name is set.
func (c *ContactCore) IsValid() bool {
	return c.Name.IsPerson() || c.Name.IsEntity()
}

// GetName returns the common name.
func (c *ContactCore) GetName() string {
	return c.Name.GetName()
}

// GetNameOptions returns more naming parameters
func (c *ContactCore) GetNameOptions() *Name {
	return &c.Name
}

func (c *ContactCore) IsPerson() bool {
	return c.Name.IsPerson()
}

func (c *ContactCore) IsEntity() bool {
	return c.Name.IsEntity()
}

func (c *ContactCore) IsBothPersonAndEntity() bool {
	return c.Name.IsBothPersonAndEntity()
}

// GetEmails ensures the Emails field is initialized and returns it.
func (c *ContactCore) GetEmails() Emails {
	if c.Emails == nil {
		c.Emails = Emails{}
	}
	return c.Emails
}

// GetPhones ensures the Phones field is initialized and returns it.
func (c *ContactCore) GetPhones() Phones {
	if c.Phones == nil {
		c.Phones = Phones{}
	}
	return c.Phones
}

// GetMails ensures the Mails field is initialized and returns it.
func (c *ContactCore) GetMails() Mails {
	if c.Mails == nil {
		c.Mails = Mails{}
	}
	return c.Mails
}

// GetUrls ensures the Urls field is initialized and returns it.
func (c *ContactCore) GetUrls() Urls {
	if c.Urls == nil {
		c.Urls = Urls{}
	}
	return c.Urls
}

// GetDates ensures the Dates field is initialized and returns it.
func (c *ContactCore) GetDates() atags.TagArrTimes {
	if c.Dates == nil {
		c.Dates = atags.TagArrTimes{}
	}
	return c.Dates
}

// Clone creates a deep copy of the ContactCore.
func (c *ContactCore) Clone() *ContactCore {
	b, err := json.Marshal(c)
	if err != nil {
		return nil // Return nil if marshaling fails
	}
	target := &ContactCore{}
	if err := json.Unmarshal(b, target); err != nil {
		return nil // Return nil if unmarshaling fails
	}
	return target // Return the deep copy
}

// LoadDetails unmarshals the Details field into the provided target interface.
func (c *ContactCore) LoadDetails(target interface{}) error {
	if len(c.Details) == 0 {
		return nil // No details to load
	}
	if target == nil {
		return errors.New("target cannot be nil")
	}
	return json.Unmarshal(c.Details, target)
}

// SaveDetails marshals the provided source interface into the Details field.
func (c *ContactCore) SaveDetails(source interface{}) error {
	if source == nil {
		return errors.New("source cannot be nil")
	}
	data, err := json.Marshal(source)
	if err != nil {
		return fmt.Errorf("failed to marshal details: %w", err)
	}
	c.Details = data
	return nil
}

type ContactValidator struct {
	Name     string
	Required bool
	Validate func() bool
}

func RunContactValidators(validators []ContactValidator) error {
	for _, v := range validators {
		if v.Required && !v.Validate() {
			return fmt.Errorf("no %s", v.Name)
		}
	}
	return nil
}

type ContactValidators []ContactValidator

func (c *ContactCore) GetValidatorEmail(mustEmail EmailType, allowDefault bool) ContactValidator {
	return ContactValidator{
		Name:     fmt.Sprintf("email of type '%s'", mustEmail.String()),
		Required: !mustEmail.IsEmpty(),
		Validate: func() bool {
			return c.Emails != nil && c.Emails.HasTypeWithDefault(mustEmail, allowDefault)
		},
	}
}

func (c *ContactCore) RunValidatorEmail(mustEmail EmailType, allowDefault bool) error {
	return RunContactValidators(ContactValidators{c.GetValidatorEmail(mustEmail, allowDefault)})
}

func (c *ContactCore) GetValidatorPhone(mustPhone PhoneType, allowDefault bool) ContactValidator {
	return ContactValidator{
		Name:     fmt.Sprintf("phone of type '%s'", mustPhone.String()),
		Required: !mustPhone.IsEmpty(),
		Validate: func() bool {
			return c.Phones != nil && c.Phones.HasTypeWithDefault(mustPhone, allowDefault)
		},
	}
}

func (c *ContactCore) RunValidatorPhone(mustPhone PhoneType, allowDefault bool) error {
	return RunContactValidators(ContactValidators{c.GetValidatorPhone(mustPhone, allowDefault)})
}

func (c *ContactCore) GetValidatorMail(mustMail MailType, allowDefault bool) ContactValidator {
	return ContactValidator{
		Name:     fmt.Sprintf("mail of type '%s'", mustMail.String()),
		Required: !mustMail.IsEmpty(),
		Validate: func() bool {
			return c.Mails != nil && c.Mails.HasTypeWithDefault(mustMail, allowDefault)
		},
	}
}

func (c *ContactCore) RunValidatorMail(mustMail MailType, allowDefault bool) error {
	return RunContactValidators(ContactValidators{c.GetValidatorMail(mustMail, allowDefault)})
}

func (c *ContactCore) GetValidatorUrl(mustUrl UrlType, allowDefault bool) ContactValidator {
	return ContactValidator{
		Name:     fmt.Sprintf("url of type '%s'", mustUrl.String()),
		Required: !mustUrl.IsEmpty(),
		Validate: func() bool {
			return c.Urls != nil && c.Urls.HasTypeWithDefault(mustUrl, allowDefault)
		},
	}
}

func (c *ContactCore) RunValidatorUrl(mustUrl UrlType, allowDefault bool) error {
	return RunContactValidators(ContactValidators{c.GetValidatorUrl(mustUrl, allowDefault)})
}

func (c *ContactCore) ValidateByEPMTypes(mustEmail EmailType, mustPhone PhoneType, mustMail MailType, allowDefault bool) error {
	validators := ContactValidators{
		c.GetValidatorEmail(mustEmail, allowDefault),
		c.GetValidatorPhone(mustPhone, allowDefault),
		c.GetValidatorMail(mustMail, allowDefault),
	}
	return RunContactValidators(validators)
}

type ContactCores []*ContactCore

// Find filters ContactCores by the provided ContactTypes.
// Returns a slice of ContactCores that match any of the given ContactTypes.
func (ccs ContactCores) Find(contactTypes ...ContactType) ContactCores {
	if len(contactTypes) == 0 {
		return nil // No types specified, return empty result
	}

	// Convert contactTypes to a map for faster lookup
	typesMap := make(map[ContactType]struct{})
	for _, ct := range contactTypes {
		typesMap[ct] = struct{}{}
	}

	// Filter ContactCores
	var filtered ContactCores
	for _, contact := range ccs {
		if contact == nil {
			continue
		}

		for _, ct := range contact.GetTypes() {
			if _, found := typesMap[ct]; found {
				filtered = append(filtered, contact)
				break
			}
		}
	}

	return filtered
}
