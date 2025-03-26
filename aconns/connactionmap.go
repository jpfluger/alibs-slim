package aconns

import (
	"fmt"
	"github.com/jpfluger/alibs-slim/autils"
	"strings"
	//	"github.com/Masterminds/semver/v3"
)

// ConnActionMap is a map of ConnActionType to ConnSystemItem.
type ConnActionMap map[ConnActionType]*ConnSystemItem

// NewConnActionMap creates a new ConnActionMap with default ConnSystemItems.
func NewConnActionMap() ConnActionMap {
	return ConnActionMap{
		CONNACTIONTYPE_CREATE:    &ConnSystemItem{},
		CONNACTIONTYPE_UPGRADE:   &ConnSystemItem{},
		CONNACTIONTYPE_DOWNGRADE: &ConnSystemItem{},
		CONNACTIONTYPE_DELETE:    &ConnSystemItem{},
	}
}

// Validate checks if the ConnActionMap is valid.
func (cvm ConnActionMap) Validate() error {
	if cvm == nil || len(cvm) == 0 {
		return fmt.Errorf("map is nil or empty")
	}
	for key, val := range cvm {
		if err := val.Validate(); err != nil {
			return fmt.Errorf("failed to validate item '%s'; %v", key.String(), err)
		}
	}
	if !cvm.HasSqlText(CONNACTIONTYPE_CREATE) {
		return fmt.Errorf("must have SQL for '%s'", CONNACTIONTYPE_CREATE.String())
	}
	return nil
}

// MustGetItem returns the ConnSystemItem for the given ConnActionType, creating a new one if it does not exist.
func (cvm ConnActionMap) MustGetItem(cvType ConnActionType) *ConnSystemItem {
	doSkip := cvType == CONNACTIONTYPE_CREATE
	if cvm == nil || len(cvm) == 0 {
		return &ConnSystemItem{DoSkip: doSkip}
	}
	val, ok := cvm[cvType]
	if !ok {
		return &ConnSystemItem{DoSkip: doSkip}
	}
	return val
}

// GetItem returns the ConnSystemItem for the given ConnActionType, or nil if it does not exist.
func (cvm ConnActionMap) GetItem(cvType ConnActionType) *ConnSystemItem {
	if cvm == nil || len(cvm) == 0 {
		return nil
	}
	val, ok := cvm[cvType]
	if !ok {
		return nil
	}
	return val
}

// HasSqlText checks if the ConnSystemItem for the given ConnActionType has non-empty text.
func (cvm ConnActionMap) HasSqlText(cvType ConnActionType) bool {
	item := cvm.GetItem(cvType)
	if item == nil {
		return false
	}
	return strings.TrimSpace(item.Text) != ""
}

// ConnSystemItem represents an item with versioning information.
type ConnSystemItem struct {
	DoSkip         bool   `json:"doSkip,omitempty"` // Flag to indicate if the item should be skipped.
	Text           string `json:"text,omitempty"`   // Text content of the item.
	importFilePath string // Path to the import file.
}

// Validate checks if the ConnSystemItem is valid.
func (cvi *ConnSystemItem) Validate() error {
	if cvi == nil {
		return fmt.Errorf("item is nil")
	}
	cvi.Text = strings.TrimSpace(cvi.Text)
	if cvi.DoSkip {
		return nil
	}
	if cvi.Text == "" {
		return fmt.Errorf("text is empty")
	}
	return nil
}

// UnpackImportWithMacroDirective processes the macro directive for importing content.
// If the text begins with "#igorg-do-import:", it reads the content from the specified file path.
func (cvi *ConnSystemItem) UnpackImportWithMacroDirective(dirOptions string) error {
	if cvi == nil {
		return fmt.Errorf("item is nil")
	}
	cvi.Text = strings.TrimSpace(cvi.Text)
	// If text begins with "#igorg-do-import:", then the text after ":" is the file path.
	if strings.HasPrefix(cvi.Text, "#igorg-do-import:") {
		myPath := strings.TrimSpace(strings.TrimPrefix(cvi.Text, "#igorg-do-import:"))
		if myPath == "" {
			return fmt.Errorf("import path is empty")
		}
		bytesContent, err := autils.ReadByFilePathWithDirOption(myPath, dirOptions)
		if err != nil {
			return fmt.Errorf("import path is invalid; %v", err)
		}
		cvi.Text = string(bytesContent)
	}
	return nil
}
