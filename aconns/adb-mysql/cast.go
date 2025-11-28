package adb_mysql

import "fmt"

// ToADBMysql attempts to cast any interface{} to *ADBMysql.
// It handles both *ADBMysql and ADBMysql input types.
func ToADBMysql(v interface{}) (*ADBMysql, error) {
	switch t := v.(type) {
	case *ADBMysql:
		// Already a pointer, just return it
		return t, nil
	case ADBMysql:
		// Value, take its address
		return &t, nil
	default:
		return nil, fmt.Errorf("cannot cast %T to *ADBMysql", v)
	}
}
