package aclient_redis

import "fmt"

// ToAClientRedis attempts to cast any interface{} to *AClientRedis.
// It handles both *AClientRedis and AClientRedis input types.
func ToAClientRedis(v interface{}) (*AClientRedis, error) {
	switch t := v.(type) {
	case *AClientRedis:
		// Already a pointer, just return it
		return t, nil
	case AClientRedis:
		// Value, take its address
		return &t, nil
	default:
		return nil, fmt.Errorf("cannot cast %T to *AClientRedis", v)
	}
}
