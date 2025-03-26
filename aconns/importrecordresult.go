package aconns

import "time"

type ImportRecordResult struct {
	Id           string    `json:"id,omitempty"`
	ImportedDate time.Time `json:"importedDate,omitempty"`
	Error        error     `json:"error,omitempty"`
}

type ImportRecordResults []*ImportRecordResult

func NewImportRecordResults(id string, err error) ImportRecordResults {
	return ImportRecordResults{&ImportRecordResult{
		Id:           id,
		ImportedDate: time.Now().UTC(),
		Error:        err,
	}}
}
