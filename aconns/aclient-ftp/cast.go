package aclient_ftp

import "fmt"

// ToAClientFTP attempts to cast any interface{} to *AClientFTP.
// It handles both *AClientFTP and AClientFTP input types.
func ToAClientFTP(v interface{}) (*AClientFTP, error) {
	switch t := v.(type) {
	case *AClientFTP:
		// Already a pointer, just return it
		return t, nil
	case AClientFTP:
		// Value, take its address
		return &t, nil
	default:
		return nil, fmt.Errorf("cannot cast %T to *AClientFTP", v)
	}
}
