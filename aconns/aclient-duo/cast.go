package aclient_duo

import "fmt"

// ToAClientDuo attempts to cast any interface{} to *AClientDuo.
// It handles both *AClientDuo and AClientDuo input types.
func ToAClientDuo(v interface{}) (*AClientDuo, error) {
	switch t := v.(type) {
	case *AClientDuo:
		// Already a pointer, just return it
		return t, nil
	case AClientDuo:
		// Value, take its address
		return &t, nil
	default:
		return nil, fmt.Errorf("cannot cast %T to *AClientDuo", v)
	}
}
