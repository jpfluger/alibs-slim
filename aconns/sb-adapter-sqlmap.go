package aconns

// ISBAdapterSqlMap is for aconns with sql capability and added mapping component
// to predefined sql.
type ISBAdapterSqlMap interface {
	ISBAdapterSql

	// RunCommand runs an exe command against the data source.
	RunCommand(text string) error
	RunMapByAction() error
	RunMapAction(action ConnActionType) error
	GetAdapterHelper() ISBAdapterHelper
}
