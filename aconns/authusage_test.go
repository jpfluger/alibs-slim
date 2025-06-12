package aconns

import "testing"

func TestAuthUsage_IsValid(t *testing.T) {
	valid := []AuthUsage{
		AUTHUSAGE_PRIMARY,
		AUTHUSAGE_MFA,
		AUTHUSAGE_SSPR,
	}

	for _, ut := range valid {
		if !ut.IsValid() {
			t.Errorf("expected %q to be valid", ut)
		}
	}

	if AuthUsage("bad").IsValid() {
		t.Error("expected invalid AuthUsage")
	}
}

func TestAuthUsage_IsEmpty(t *testing.T) {
	if !AuthUsage("").IsEmpty() {
		t.Error("expected empty AuthUsage")
	}
	if AUTHUSAGE_MFA.IsEmpty() {
		t.Error("expected non-empty AuthUsage")
	}
}

func TestAuthUsage_TrimSpace(t *testing.T) {
	input := AuthUsage(" sspr ")
	if input.TrimSpace() != AUTHUSAGE_SSPR {
		t.Error("expected trimmed to match AUTHUSAGE_SSPR")
	}
}

func TestAuthUsage_String(t *testing.T) {
	ut := AuthUsage("  primary ")
	if ut.String() != "primary" {
		t.Errorf("unexpected string: %q", ut.String())
	}
}

func TestAuthUsages_Has(t *testing.T) {
	usages := AuthUsages{AUTHUSAGE_MFA}
	if !usages.Has(AUTHUSAGE_MFA) {
		t.Error("expected to find AUTHUSAGE_MFA")
	}
	if usages.Has(AUTHUSAGE_PRIMARY) {
		t.Error("did not expect to find AUTHUSAGE_PRIMARY")
	}
}

func TestAuthUsages_Validate(t *testing.T) {
	valid := AuthUsages{AUTHUSAGE_MFA}
	if err := valid.Validate(); err != nil {
		t.Errorf("expected valid, got error: %v", err)
	}

	invalid := AuthUsages{AuthUsage("bogus")}
	if err := invalid.Validate(); err == nil {
		t.Error("expected validation error for invalid usage type")
	}
}
