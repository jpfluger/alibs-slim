package aclient_smtp

import "fmt"

// ToAClientSMTP attempts to cast any interface{} to *AClientSMTP.
// It handles both *AClientSMTP and AClientSMTP input types.
func ToAClientSMTP(v interface{}) (*AClientSMTP, error) {
	switch t := v.(type) {
	case *AClientSMTP:
		// Already a pointer, just return it
		return t, nil
	case AClientSMTP:
		// Value, take its address
		return &t, nil
	default:
		return nil, fmt.Errorf("cannot cast %T to *AClientSMTP", v)
	}
}
