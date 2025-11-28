package aclient_sftp

import "fmt"

// ToAClientSFTP attempts to cast any interface{} to *AClientSFTP.
// It handles both *AClientSFTP and AClientSFTP input types.
func ToAClientSFTP(v interface{}) (*AClientSFTP, error) {
	switch t := v.(type) {
	case *AClientSFTP:
		// Already a pointer, just return it
		return t, nil
	case AClientSFTP:
		// Value, take its address
		return &t, nil
	default:
		return nil, fmt.Errorf("cannot cast %T to *AClientSFTP", v)
	}
}
