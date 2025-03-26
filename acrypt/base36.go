package acrypt

import "strconv"

const base36 = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"

// EncodeBase36 converts an integer to a base36 string
func EncodeBase36(n int) string {
	if n == 0 {
		return "0"
	}
	var result string
	for n > 0 {
		result = string(base36[n%36]) + result
		n /= 36
	}
	return result
}

// DecodeBase36 converts a base36 string to an integer
func DecodeBase36(s string) (int, error) {
	ii, err := strconv.ParseInt(s, 36, 64)
	return int(ii), err
}
