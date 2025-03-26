package ahttp

import (
	"fmt"
	"github.com/jpfluger/alibs-slim/atemplates"
	"github.com/jpfluger/alibs-slim/azb"
	"strings"
)

const (
	PUI_DIRECTION_LEFT         azb.ZBType = "left"
	PUI_DIRECTION_RIGHT        azb.ZBType = "right"
	PUI_SIDEBAR_STATUS_OPEN    azb.ZBType = "open"
	PUI_SIDEBAR_STATUS_CLOSE   azb.ZBType = "close"
	PUI_SIDEBAR_STATUS_DISMISS azb.ZBType = "dismiss"
)

// UIElemOption represents an option within a UI element, such as a dropdown.
type UIElemOption struct {
	Type      azb.ZBType `json:"type,omitempty"`      // The internal type of the option.
	Title     string     `json:"title,omitempty"`     // The display title of the option.
	IsDefault bool       `json:"isDefault,omitempty"` // Indicates if this option is the default selection.
}

// GetType returns the internal type of the UI option.
func (ui *UIElemOption) GetType() azb.ZBType {
	return ui.Type
}

// GetTitle returns the display title of the UI option, defaulting to the capitalized type if not set.
func (ui *UIElemOption) GetTitle() string {
	if ui.Title == "" {
		return atemplates.ToUpperFirst(ui.Type.String())
	}
	return ui.Title
}

// GetValue returns the string representation of the option's type.
func (ui *UIElemOption) GetValue() string {
	return ui.Type.String()
}

// GetIsDefault returns whether the UI option is the default selection.
func (ui *UIElemOption) GetIsDefault() bool {
	return ui.IsDefault
}

// UIElemOptions is a slice of UIElemOption pointers.
type UIElemOptions []*UIElemOption

// FindDefault searches for and returns the default UI option, if any.
func (uis UIElemOptions) FindDefault() *UIElemOption {
	for _, ui := range uis {
		if ui != nil && ui.IsDefault {
			return ui
		}
	}
	return nil
}

// FindByType searches for and returns the UI option with the matching type.
func (uis UIElemOptions) FindByType(target azb.ZBType) *UIElemOption {
	for _, ui := range uis {
		if ui != nil && ui.Type == target {
			return ui
		}
	}
	return nil
}

// FindByValue is an alias for FindByType, as they perform the same function.
func (uis UIElemOptions) FindByValue(target azb.ZBType) *UIElemOption {
	return uis.FindByType(target)
}

// HasByType checks if there is a UI option with the matching type.
func (uis UIElemOptions) HasByType(target azb.ZBType) bool {
	return uis.FindByType(target) != nil
}

// ToHTML generates an HTML string representing the options as HTML 'option' elements.
func (uis UIElemOptions) ToHTML(currValue string) string {
	var sb strings.Builder
	for _, opt := range uis {
		if opt != nil {
			selectedAttr := ""
			if opt.Type.String() == currValue {
				selectedAttr = " selected"
			}
			sb.WriteString(fmt.Sprintf(`<option value="%s"%s>%s</option>`, opt.GetValue(), selectedAttr, opt.GetTitle()))
		}
	}
	return sb.String()
}
