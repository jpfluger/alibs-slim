package acontact

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestContactCore_Validate(t *testing.T) {
	t.Run("Valid ContactCore", func(t *testing.T) {
		c := &ContactCore{
			Name: Name{First: "John", Last: "Doe"},
		}
		assert.NoError(t, c.Validate())
	})

	t.Run("Invalid ContactCore with nil Name", func(t *testing.T) {
		c := &ContactCore{}
		assert.Error(t, c.Validate())
	})
}

func TestContactCore_GetMethods(t *testing.T) {
	c := &ContactCore{}

	t.Run("GetEmails initializes if nil", func(t *testing.T) {
		emails := c.GetEmails()
		assert.NotNil(t, emails)
	})

	t.Run("GetPhones initializes if nil", func(t *testing.T) {
		phones := c.GetPhones()
		assert.NotNil(t, phones)
	})

	t.Run("GetMails initializes if nil", func(t *testing.T) {
		mails := c.GetMails()
		assert.NotNil(t, mails)
	})

	t.Run("GetUrls initializes if nil", func(t *testing.T) {
		urls := c.GetUrls()
		assert.NotNil(t, urls)
	})

	t.Run("GetDates initializes if nil", func(t *testing.T) {
		dates := c.GetDates()
		assert.NotNil(t, dates)
	})
}

// TestContactCore_GetName tests the GetName method of the ContactCore struct.
func TestContactCore_GetName(t *testing.T) {
	contact := &ContactCore{
		Name: Name{
			Legal: "John P Smith",
			Short: "John",
		},
	}
	if got := contact.Name.MustGetShort(); got != "John" {
		t.Errorf("ContactCore.GetName() = %v, want %v", got, "John")
	}

	contact.Name.Short = ""
	if got := contact.Name.MustGetFull(); got != "John P Smith" {
		t.Errorf("ContactCore.GetName() = %v, want %v", got, "John P Smith")
	}
}

// TestContactCore_GetNameLegal tests the GetLegalName method of the ContactCore struct.
func TestContactCore_GetNameLegal(t *testing.T) {
	contact := &ContactCore{
		Name: Name{Legal: "Acme, Inc."},
	}
	if got := contact.Name.MustGetLegal(); got != "Acme, Inc." {
		t.Errorf("ContactCore.GetNameLegal() = %v, want %v", got, "Acme, Inc.")
	}
}

// TestContactCore_GetNameShort tests the GetShortName method of the ContactCore struct.
func TestContactCore_GetNameShort(t *testing.T) {
	contact := &ContactCore{
		Name: Name{Short: "Acme"},
	}
	if got := contact.Name.MustGetShort(); got != "Acme" {
		t.Errorf("ContactCore.GetNameShort() = %v, want %v", got, "Acme")
	}
}

// TestContactCore_Clone tests the Clone method of the ContactCore struct.
func TestContactCore_Clone(t *testing.T) {
	contact := &ContactCore{
		Name: Name{
			Legal: "John P Smith",
			Short: "John",
		},
		Emails: Emails{&Email{Address: "john@example.com"}},
	}
	clonedContact := contact.Clone()
	if clonedContact.Name.Legal != contact.Name.Legal {
		t.Errorf("ContactCore.Clone() failed to clone NameLegal")
	}
	if len(clonedContact.Emails) == 0 || clonedContact.Emails[0].Address.String() != contact.Emails[0].Address.String() {
		t.Errorf("ContactCore.Clone() failed to clone Emails")
	}
}

// ContactDetails represents additional details for a contact.
type ContactDetails struct {
	SocialSecurity string `json:"socialSecurity,omitempty"`
	Passport       string `json:"passport,omitempty"`
	NationalID     string `json:"nationalId,omitempty"`
}

// TestContactCore_LoadDetails tests the LoadDetails method of the ContactCore struct.
func TestContactCore_LoadDetails(t *testing.T) {
	contact := &ContactCore{
		Details: json.RawMessage(`{"socialSecurity": "123-45-6789", "passport": "X1234567"}`),
	}

	var details ContactDetails
	err := contact.LoadDetails(&details)

	assert.NoError(t, err, "LoadDetails should not return an error")
	assert.Equal(t, "123-45-6789", details.SocialSecurity, "SocialSecurity should match the input")
	assert.Equal(t, "X1234567", details.Passport, "Passport should match the input")
	assert.Empty(t, details.NationalID, "NationalID should be empty")
}

// TestContactCore_LoadDetails_NoDetails tests LoadDetails when the Details field is empty.
func TestContactCore_LoadDetails_NoDetails(t *testing.T) {
	contact := &ContactCore{}

	var details ContactDetails
	err := contact.LoadDetails(&details)

	assert.NoError(t, err, "LoadDetails should not return an error when Details is empty")
	assert.Empty(t, details, "Details should remain empty when no Details are provided")
}

// TestContactCore_LoadDetails_NilTarget tests LoadDetails with a nil target.
func TestContactCore_LoadDetails_NilTarget(t *testing.T) {
	contact := &ContactCore{
		Details: json.RawMessage(`{"socialSecurity": "123-45-6789"}`),
	}

	err := contact.LoadDetails(nil)

	assert.Error(t, err, "LoadDetails should return an error when target is nil")
	assert.Equal(t, "target cannot be nil", err.Error())
}

// TestContactCore_SaveDetails tests the SaveDetails method of the ContactCore struct.
func TestContactCore_SaveDetails(t *testing.T) {
	contact := &ContactCore{}
	details := ContactDetails{
		SocialSecurity: "987-65-4321",
		Passport:       "Y7654321",
		NationalID:     "A123456",
	}

	err := contact.SaveDetails(details)

	assert.NoError(t, err, "SaveDetails should not return an error")
	assert.JSONEq(t, `{"socialSecurity": "987-65-4321", "passport": "Y7654321", "nationalId": "A123456"}`, string(contact.Details))
}

// TestContactCore_SaveDetails_NilSource tests SaveDetails with a nil source.
func TestContactCore_SaveDetails_NilSource(t *testing.T) {
	contact := &ContactCore{}

	err := contact.SaveDetails(nil)

	assert.Error(t, err, "SaveDetails should return an error when source is nil")
	assert.Equal(t, "source cannot be nil", err.Error())
}

// Test for GetValidatorEmail
func TestContactCore_GetValidatorEmail(t *testing.T) {
	contact := &ContactCore{
		Emails: Emails{{Type: "work"}, {Type: "personal", IsDefault: true}},
	}

	validator := contact.GetValidatorEmail("work", false)
	assert.True(t, validator.Validate(), "Validator should pass for email type 'work'")

	validator = contact.GetValidatorEmail("personal", true)
	assert.True(t, validator.Validate(), "Validator should pass for default email type 'personal'")

	validator = contact.GetValidatorEmail("nonexistent", false)
	assert.False(t, validator.Validate(), "Validator should fail for non-existent email type 'nonexistent'")
}

// Test for RunValidatorEmail
func TestContactCore_RunValidatorEmail(t *testing.T) {
	contact := &ContactCore{
		Emails: Emails{{Type: "work"}, {Type: "personal", IsDefault: true}},
	}

	err := contact.RunValidatorEmail("work", false)
	assert.NoError(t, err, "RunValidatorEmail should not return an error for email type 'work'")

	err = contact.RunValidatorEmail("nonexistent", false)
	assert.Error(t, err, "RunValidatorEmail should return an error for non-existent email type 'nonexistent'")
	assert.Equal(t, "no email of type 'nonexistent'", err.Error())
}

// Test for GetValidatorPhone
func TestContactCore_GetValidatorPhone(t *testing.T) {
	contact := &ContactCore{
		Phones: Phones{{Type: "mobile"}, {Type: "office", IsDefault: true}},
	}

	validator := contact.GetValidatorPhone("mobile", false)
	assert.True(t, validator.Validate(), "Validator should pass for phone type 'mobile'")

	validator = contact.GetValidatorPhone("office", true)
	assert.True(t, validator.Validate(), "Validator should pass for default phone type 'office'")

	validator = contact.GetValidatorPhone("home", false)
	assert.False(t, validator.Validate(), "Validator should fail for non-existent phone type 'home'")
}

// Test for RunValidatorPhone
func TestContactCore_RunValidatorPhone(t *testing.T) {
	contact := &ContactCore{
		Phones: Phones{{Type: "mobile"}, {Type: "office", IsDefault: true}},
	}

	err := contact.RunValidatorPhone("mobile", false)
	assert.NoError(t, err, "RunValidatorPhone should not return an error for phone type 'mobile'")

	err = contact.RunValidatorPhone("home", false)
	assert.Error(t, err, "RunValidatorPhone should return an error for non-existent phone type 'home'")
	assert.Equal(t, "no phone of type 'home'", err.Error())
}

// Test for GetValidatorMail
func TestContactCore_GetValidatorMail(t *testing.T) {
	contact := &ContactCore{
		Mails: Mails{{Type: "home"}, {Type: "business", IsDefault: true}},
	}

	validator := contact.GetValidatorMail("home", false)
	assert.True(t, validator.Validate(), "Validator should pass for mail type 'home'")

	validator = contact.GetValidatorMail("business", true)
	assert.True(t, validator.Validate(), "Validator should pass for default mail type 'business'")

	validator = contact.GetValidatorMail("other", false)
	assert.False(t, validator.Validate(), "Validator should fail for non-existent mail type 'other'")
}

// Test for RunValidatorMail
func TestContactCore_RunValidatorMail(t *testing.T) {
	contact := &ContactCore{
		Mails: Mails{{Type: "home"}, {Type: "business", IsDefault: true}},
	}

	err := contact.RunValidatorMail("home", false)
	assert.NoError(t, err, "RunValidatorMail should not return an error for mail type 'home'")

	err = contact.RunValidatorMail("other", false)
	assert.Error(t, err, "RunValidatorMail should return an error for non-existent mail type 'other'")
	assert.Equal(t, "no mail of type 'other'", err.Error())
}

// Test for ValidateByEPMTypes
func TestContactCore_ValidateByEPMTypes(t *testing.T) {
	contact := &ContactCore{
		Emails: Emails{{Type: "work"}},
		Phones: Phones{{Type: "mobile"}},
		Mails:  Mails{{Type: "home"}},
	}

	err := contact.ValidateByEPMTypes("work", "mobile", "home", false)
	assert.NoError(t, err, "ValidateByEPMTypes should pass for valid email, phone, and mail types")

	err = contact.ValidateByEPMTypes("work", "office", "home", false)
	assert.Error(t, err, "ValidateByEPMTypes should fail for missing phone type 'office'")
	assert.Equal(t, "no phone of type 'office'", err.Error())
}

func TestContactCores_Find(t *testing.T) {
	contactCores := ContactCores{
		&ContactCore{
			Name: Name{First: "John", Last: "Doe"},
		},
		&ContactCore{
			Name: Name{Company: "TechCorp"},
		},
		&ContactCore{
			Name: Name{First: "Jane", Last: "Smith", Company: "DesignStudio"},
		},
	}

	t.Run("Find Entity Contacts", func(t *testing.T) {
		results := contactCores.Find(CONTACTTYPE_ENTITY)
		assert.Equal(t, 2, len(results)) // Two contacts are entities
	})

	t.Run("Find Person Contacts", func(t *testing.T) {
		results := contactCores.Find(CONTACTTYPE_PERSON)
		assert.Equal(t, 2, len(results)) // Two contacts are persons
	})

	t.Run("Find Both Entity and Person Contacts", func(t *testing.T) {
		results := contactCores.Find(CONTACTTYPE_ENTITY, CONTACTTYPE_PERSON)
		assert.Equal(t, 3, len(results)) // All three contacts qualify
	})

	t.Run("Find with No Types", func(t *testing.T) {
		results := contactCores.Find()
		assert.Empty(t, results) // No types specified, should return empty
	})
}
