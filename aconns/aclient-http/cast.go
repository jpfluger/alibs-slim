package aclient_http

import "fmt"

// ToAClientHTTP attempts to cast any interface{} to *AClientHTTP.
// It handles both *AClientHTTP and AClientHTTP input types.
func ToAClientHTTP(v interface{}) (*AClientHTTP, error) {
	switch t := v.(type) {
	case *AClientHTTP:
		// Already a pointer, just return it
		return t, nil
	case AClientHTTP:
		// Value, take its address
		return &t, nil
	default:
		return nil, fmt.Errorf("cannot cast %T to *AClientHTTP", v)
	}
}
