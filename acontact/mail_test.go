package acontact

import (
	"testing"
)

// TestMail_IsDefault tests if the IsDefault flag is set correctly.
func TestMail_IsDefault(t *testing.T) {
	mail := &Mail{IsDefault: true}
	if !mail.IsDefault {
		t.Errorf("Expected IsDefault to be true, got false")
	}
}

// TestMails_FindByType tests the FindByType method of the Mails type.
func TestMails_FindByType(t *testing.T) {
	mails := Mails{
		&Mail{Type: MAILTYPE_HOME, IsDefault: false},
		&Mail{Type: MAILTYPE_WORK, IsDefault: true},
	}

	if got := mails.FindByType(MAILTYPE_WORK); !got.IsDefault {
		t.Errorf("FindByType(MAILTYPE_WORK) should return the default work mail")
	}
	if got := mails.FindByType(MAILTYPE_HOME); got.IsDefault {
		t.Errorf("FindByType(MAILTYPE_HOME) should not return the default home mail")
	}
}

// TestMails_FindByTypeOrDefault tests the FindByTypeOrDefault method of the Mails type.
func TestMails_FindByTypeOrDefault(t *testing.T) {
	mails := Mails{
		&Mail{Type: MAILTYPE_HOME, IsDefault: true},
		&Mail{Type: MAILTYPE_WORK, IsDefault: false},
	}

	if got := mails.FindByTypeOrDefault(MAILTYPE_ALTERNATE); got.Type != MAILTYPE_HOME {
		t.Errorf("FindByTypeOrDefault(MAILTYPE_ALTERNATE) should return the default home mail")
	}
}

// TestMails_HasType tests the HasType method of the Mails type.
func TestMails_HasType(t *testing.T) {
	mails := Mails{
		&Mail{Type: MAILTYPE_HOME, IsDefault: true},
	}

	if !mails.HasType(MAILTYPE_HOME) {
		t.Errorf("HasType(MAILTYPE_HOME) should return true")
	}
	if mails.HasType(MAILTYPE_WORK) {
		t.Errorf("HasType(MAILTYPE_WORK) should return false")
	}
}

// TestMails_Clone tests the Clone method of the Mails type.
func TestMails_Clone(t *testing.T) {
	mails := Mails{
		&Mail{Type: MAILTYPE_HOME, IsDefault: true},
	}
	clonedMails := mails.Clone()
	if len(clonedMails) != 1 || clonedMails[0].Type != MAILTYPE_HOME {
		t.Errorf("Clone() did not properly clone the mails slice")
	}
}

// Additional tests should be written for MergeFrom, Set, and Remove methods.
