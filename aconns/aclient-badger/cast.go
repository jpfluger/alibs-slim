package aclient_badger

import "fmt"

// ToAClientBadger attempts to cast any interface{} to *AClientBadger.
// It handles both *AClientBadger and AClientBadger input types.
func ToAClientBadger(v interface{}) (*AClientBadger, error) {
	switch t := v.(type) {
	case *AClientBadger:
		// Already a pointer, just return it
		return t, nil
	case AClientBadger:
		// Value, take its address
		return &t, nil
	default:
		return nil, fmt.Errorf("cannot cast %T to *AClientBadger", v)
	}
}
