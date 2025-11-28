package adb_pg

import "fmt"

// ToADBPG attempts to cast any interface{} to *ADBPG.
// It handles both *ADBPG and ADBPG input types.
func ToADBPG(v interface{}) (*ADBPG, error) {
	switch t := v.(type) {
	case *ADBPG:
		// Already a pointer, just return it
		return t, nil
	case ADBPG:
		// Value, take its address
		return &t, nil
	default:
		return nil, fmt.Errorf("cannot cast %T to *ADBPG", v)
	}
}
