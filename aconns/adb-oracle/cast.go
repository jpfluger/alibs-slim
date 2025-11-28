package adb_oracle

import "fmt"

// ToADBOracle attempts to cast any interface{} to *ADBOracle.
// It handles both *ADBOracle and ADBOracle input types.
func ToADBOracle(v interface{}) (*ADBOracle, error) {
	switch t := v.(type) {
	case *ADBOracle:
		// Already a pointer, just return it
		return t, nil
	case ADBOracle:
		// Value, take its address
		return &t, nil
	default:
		return nil, fmt.Errorf("cannot cast %T to *ADBOracle", v)
	}
}
