package ahttp

import (
	"github.com/jpfluger/alibs-slim/azb"
	"testing"
)

// TestUIElemOption_GetType tests the GetType method of UIElemOption.
func TestUIElemOption_GetType(t *testing.T) {
	ui := UIElemOption{Type: "button"}
	if ui.GetType() != "button" {
		t.Errorf("GetType() = %v, want %v", ui.GetType(), "button")
	}
}

// TestUIElemOption_GetTitle tests the GetTitle method of UIElemOption.
func TestUIElemOption_GetTitle(t *testing.T) {
	ui := UIElemOption{Type: "button", Title: "Submit"}
	if ui.GetTitle() != "Submit" {
		t.Errorf("GetTitle() = %v, want %v", ui.GetTitle(), "Submit")
	}

	// Test with empty title, should return capitalized type.
	ui = UIElemOption{Type: "button"}
	if ui.GetTitle() != "Button" {
		t.Errorf("GetTitle() = %v, want %v", ui.GetTitle(), "Button")
	}
}

// TestUIElemOption_GetValue tests the GetValue method of UIElemOption.
func TestUIElemOption_GetValue(t *testing.T) {
	ui := UIElemOption{Type: "button"}
	if ui.GetValue() != "button" {
		t.Errorf("GetValue() = %v, want %v", ui.GetValue(), "button")
	}
}

// TestUIElemOption_GetIsDefault tests the GetIsDefault method of UIElemOption.
func TestUIElemOption_GetIsDefault(t *testing.T) {
	ui := UIElemOption{IsDefault: true}
	if !ui.GetIsDefault() {
		t.Errorf("GetIsDefault() should return true")
	}
}

// TestUIElemOptions_FindDefault tests the FindDefault method of UIElemOptions.
func TestUIElemOptions_FindDefault(t *testing.T) {
	uis := UIElemOptions{
		&UIElemOption{Type: "button", IsDefault: false},
		&UIElemOption{Type: "submit", IsDefault: true},
	}
	defaultOpt := uis.FindDefault()
	if defaultOpt == nil || defaultOpt.Type != "submit" {
		t.Errorf("FindDefault() should return the 'submit' option")
	}
}

// TestUIElemOptions_FindByType tests the FindByType method of UIElemOptions.
func TestUIElemOptions_FindByType(t *testing.T) {
	uis := UIElemOptions{
		&UIElemOption{Type: "button"},
		&UIElemOption{Type: "submit"},
	}
	opt := uis.FindByType("submit")
	if opt == nil || opt.Type != "submit" {
		t.Errorf("FindByType() should return the 'submit' option")
	}
}

// TestUIElemOptions_HasByType tests the HasByType method of UIElemOptions.
func TestUIElemOptions_HasByType(t *testing.T) {
	uis := UIElemOptions{
		&UIElemOption{Type: "button"},
		&UIElemOption{Type: "submit"},
	}
	if !uis.HasByType("button") {
		t.Errorf("HasByType() should return true for 'button'")
	}
}

// TestUIElemOptions_ToHTML tests the ToHTML method of UIElemOptions.
func TestUIElemOptions_ToHTML(t *testing.T) {
	uis := UIElemOptions{
		{Type: azb.ZBType("button"), Title: "Button", IsDefault: false},
		{Type: azb.ZBType("submit"), Title: "Submit", IsDefault: true},
	}
	currValue := "submit" // Current value should match the Type of the second option.

	html := uis.ToHTML(currValue)
	expectedHTML := `<option value="button">Button</option><option value="submit" selected>Submit</option>`

	if html != expectedHTML {
		t.Errorf("ToHTML() = %v, want %v", html, expectedHTML)
	}
}
