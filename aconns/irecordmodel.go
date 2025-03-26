package aconns

// IRecordModel is used to reduce "mostly" repeatable code within object models-tenants.
// An example can be found when inheriting pg/pg-model-base.go and running OnImportRunAction
type IRecordModel interface {
	IRecordAction

	Insert(ri *RI) error
	Update(ri *RI) error
	Upsert(ri *RI) error
	Select(ri *RI) error

	SelectIntoNewObject(ri *RI) (IRecordModel, error)
	GetRecordSecurity() RecordSecurity
	GetPrimaryKeyAsString() string
}
