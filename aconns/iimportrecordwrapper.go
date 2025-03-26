package aconns

import (
	"github.com/jpfluger/alibs-slim/atags"
	"time"
)

type RecImportOnExistsType int

// These values work in conjunction with those on individual records.
// The logic is as follows:
//  1. Examine the record for a ModifiedTime and RecImportOnExistsType value.
//     a. If the ModifiedTime is valid, then use this one for instructions.
//     b. If the ModifiedTime is invalid, then proceed to 2.
//  2. Use the RecImportOnExistsType value from RecordImport
//     and the ModifiedDate from RecordImportsWrapper to make the import decision.
//     --> a sophisticated program might employ a conflict manager to allow or disallow updates.
const (
	REC_IMPORT_ON_EXISTS_IGNORE RecImportOnExistsType = iota
	REC_IMPORT_ON_EXISTS_UPDATE
	REC_IMPORT_ON_EXISTS_UPDATE_IF_NEWER
)

type IImportRecordWrapper interface {
	GetDirImports() string
	GetRecImportOnExistsType() RecImportOnExistsType
	GetModifiedDate() time.Time
	GetTag(key atags.TagKey) *atags.TagKeyValueString
}
