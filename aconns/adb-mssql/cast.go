package adb_mssql

import "fmt"

// ToADBMSSql attempts to cast any interface{} to *ADBMSSql.
// It handles both *ADBMSSql and ADBMSSql input types.
func ToADBMSSql(v interface{}) (*ADBMSSql, error) {
	switch t := v.(type) {
	case *ADBMSSql:
		// Already a pointer, just return it
		return t, nil
	case ADBMSSql:
		// Value, take its address
		return &t, nil
	default:
		return nil, fmt.Errorf("cannot cast %T to *ADBMSSql", v)
	}
}
