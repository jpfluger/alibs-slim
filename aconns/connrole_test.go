package aconns

import (
	"testing"
)

func TestConnRole_IsEmpty(t *testing.T) {
	tests := []struct {
		input    ConnRole
		expected bool
	}{
		{"", true},
		{"   ", true},
		{"master", false},
		{" auth ", false},
	}

	for _, test := range tests {
		if got := test.input.IsEmpty(); got != test.expected {
			t.Errorf("ConnRole(%q).IsEmpty() = %v; want %v", test.input, got, test.expected)
		}
	}
}

func TestConnRole_TrimSpace(t *testing.T) {
	input := ConnRole("  auth ")
	expected := ConnRole("auth")
	if result := input.TrimSpace(); result != expected {
		t.Errorf("TrimSpace() = %q; want %q", result, expected)
	}
}

func TestConnRole_String(t *testing.T) {
	role := ConnRole("  tenant  ")
	if got := role.String(); got != "tenant" {
		t.Errorf("String() = %q; want %q", got, "tenant")
	}
}

func TestConnRole_IsValid(t *testing.T) {
	tests := []struct {
		role     ConnRole
		expected bool
	}{
		{CONNROLE_MASTER, true},
		{CONNROLE_AUTH, true},
		{CONNROLE_TENANT, true},
		{"invalid", false},
		{"", false},
	}

	for _, test := range tests {
		if got := test.role.IsValid(); got != test.expected {
			t.Errorf("ConnRole(%q).IsValid() = %v; want %v", test.role, got, test.expected)
		}
	}
}

func TestConnRoles_HasRole(t *testing.T) {
	roles := ConnRoles{CONNROLE_MASTER, "  auth ", "TENANT"}
	tests := []struct {
		needle   ConnRole
		expected bool
	}{
		{"master", true},
		{"auth", true},
		{"tenant", false}, // case-sensitive check
		{"  auth  ", true},
		{"unknown", false},
	}

	for _, test := range tests {
		if got := roles.HasRole(test.needle); got != test.expected {
			t.Errorf("HasRole(%q) = %v; want %v", test.needle, got, test.expected)
		}
	}
}

func TestConnRoles_Validate(t *testing.T) {
	valid := ConnRoles{CONNROLE_MASTER, "auth", "tenant"}
	if err := valid.Validate(); err != nil {
		t.Errorf("Validate() returned unexpected error: %v", err)
	}

	invalid := ConnRoles{"master", " ", "unsupported"}
	if err := invalid.Validate(); err == nil {
		t.Error("Validate() expected error for invalid roles, got nil")
	}
}
