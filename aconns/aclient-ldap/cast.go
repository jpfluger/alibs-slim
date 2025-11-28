package aclient_ldap

import "fmt"

// ToAClientLDAP attempts to cast any interface{} to *AClientLDAP.
// It handles both *AClientLDAP and AClientLDAP input types.
func ToAClientLDAP(v interface{}) (*AClientLDAP, error) {
	switch t := v.(type) {
	case *AClientLDAP:
		// Already a pointer, just return it
		return t, nil
	case AClientLDAP:
		// Value, take its address
		return &t, nil
	default:
		return nil, fmt.Errorf("cannot cast %T to *AClientLDAP", v)
	}
}
