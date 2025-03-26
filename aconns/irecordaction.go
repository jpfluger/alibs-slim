package aconns

type IRecordAction interface {
	// Import is required for imports to work.
	// Notice the function returns "ImportRecordResults". This must
	// be true for a single object or for an array of objects.
	// Why? When the object is created from model type in the import
	// function, the same "Import" function is run regardless if the
	// object is a single object or an array of objects. The array of
	// objects would normally and naturally return "ImportRecordResults"
	// but a single object would naturually return a single "ImportRecordResult".
	// Although the implementation is counter-intuitive, it makes sense
	// from a single interface implementation since both a single object
	// and array of objects are treated the same.
	Import(ri *RI, irw IImportRecordWrapper) (ImportRecordResults, error)
}
