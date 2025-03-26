package acrypt

import (
	"fmt"
	"math/rand"
	"time"
)

// Adapted from:
// Copyright 2017 Seth Vargo <seth@sethvargo.com>
// MIT License
// ref: https://github.com/sethvargo/go-password

// MiniRandomCodes represents a slice of random codes.
type MiniRandomCodes []string

// ToStringArray converts MiniRandomCodes to a slice of strings.
func (codes MiniRandomCodes) ToStringArray() []string {
	return codes
}

// MatchAndRemove finds a target code and removes it from the slice.
func (codes *MiniRandomCodes) MatchAndRemove(target string) bool {
	if len(*codes) == 0 || target == "" {
		return false
	}

	for i, code := range *codes {
		if code == target {
			*codes = append((*codes)[:i], (*codes)[i+1:]...)
			return true
		}
	}
	return false
}

// Generate creates a specified number of random codes with a given length and optional divider.
func (codes *MiniRandomCodes) Generate(count int, length int, divider string) error {
	return generateRandomCodes(codes, count, length, divider, "")
}

// GenerateWithCharSet creates a specified number of random codes with a given length and optional divider, using the provided character set.
// The generated codes are stored in the MiniRandomCodes slice.
func (codes *MiniRandomCodes) GenerateWithCharSet(count int, length int, divider string, charset string) error {
	return generateRandomCodes(codes, count, length, divider, charset)
}

// generateRandomCodes is a helper function that generates random codes with an optional divider and charset.
func generateRandomCodes(codes *MiniRandomCodes, count int, length int, divider string, charset string) error {
	if count <= 0 {
		return fmt.Errorf("count must be greater than zero")
	}
	if length <= 0 {
		return fmt.Errorf("length must be greater than zero")
	}

	const defaultCharset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	if charset == "" {
		charset = defaultCharset
	}

	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	code := make([]byte, length)

	for i := 0; i < count; i++ {
		for j := range code {
			code[j] = charset[seededRand.Intn(len(charset))]
		}

		if divider != "" && length > 2 {
			mid := length / 2
			*codes = append(*codes, string(code[:mid])+divider+string(code[mid:]))
		} else {
			*codes = append(*codes, string(code))
		}
	}

	return nil
}

// GenerateMiniRandomCodes generates a specified number of random codes with a given length.
func GenerateMiniRandomCodes(count, length int) ([]string, error) {
	var codes MiniRandomCodes
	err := generateRandomCodes(&codes, count, length, "", "")
	if err != nil {
		return nil, err
	}
	return codes.ToStringArray(), nil
}
