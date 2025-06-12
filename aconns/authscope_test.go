package aconns

import "testing"

func TestAuthScope_IsValid(t *testing.T) {
	valid := []AuthScope{
		AUTHSCOPE_MASTER,
		AUTHSCOPE_DOMAIN,
		AUTHSCOPE_MODULE,
		AUTHSCOPE_ADMIN,
	}

	for _, tier := range valid {
		if !tier.IsValid() {
			t.Errorf("expected %q to be valid", tier)
		}
	}

	if AuthScope("bad").IsValid() {
		t.Error("expected invalid AuthScope")
	}
}

func TestAuthScope_IsEmpty(t *testing.T) {
	if !AuthScope("").IsEmpty() {
		t.Error("expected empty AuthScope")
	}
	if AUTHSCOPE_DOMAIN.IsEmpty() {
		t.Error("expected non-empty AuthScope")
	}
}

func TestAuthScope_TrimSpace(t *testing.T) {
	input := AuthScope(" module ")
	if input.TrimSpace() != AUTHSCOPE_MODULE {
		t.Error("expected trimmed to match AUTHSCOPE_MODULE")
	}
}

func TestAuthScope_String(t *testing.T) {
	tier := AuthScope("  master ")
	if tier.String() != "master" {
		t.Errorf("unexpected string: %q", tier.String())
	}
}

func TestAuthScopes_Has(t *testing.T) {
	tiers := AuthScopes{AUTHSCOPE_MASTER}
	if !tiers.Has(AUTHSCOPE_MASTER) {
		t.Error("expected to find AUTHSCOPE_MASTER")
	}
	if tiers.Has(AUTHSCOPE_DOMAIN) {
		t.Error("did not expect to find AUTHSCOPE_DOMAIN")
	}
}

func TestAuthScopes_Validate(t *testing.T) {
	valid := AuthScopes{AUTHSCOPE_MASTER}
	if err := valid.Validate(); err != nil {
		t.Errorf("expected valid, got error: %v", err)
	}

	invalid := AuthScopes{AuthScope("invalid")}
	if err := invalid.Validate(); err == nil {
		t.Error("expected validation error for invalid tier")
	}
}
