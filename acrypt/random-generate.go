package acrypt

// Copyright 2017 Seth Vargo <seth@sethvargo.com>
// MIT License
// ref: https://github.com/sethvargo/go-password

import (
	"fmt"
	"github.com/sethvargo/go-password/password"
	"math/rand"
	"time"
)

// RandomTextGenerator is a struct that defines parameters for generating random text.
type RandomTextGenerator struct {
	Length      int  // Total length of the random text
	NumDigits   int  // Number of digits in the random text
	NumSymbols  int  // Number of symbols in the random text
	NoUpper     bool // If true, no uppercase letters will be included
	AllowRepeat bool // If true, characters can be repeated
}

// Generate produces a random string based on the RandomTextGenerator's settings.
func (rg *RandomTextGenerator) Generate() (string, error) {
	return password.Generate(rg.Length, rg.NumDigits, rg.NumSymbols, rg.NoUpper, rg.AllowRepeat)
}

// RandGenerate4Digits generates a 4-digit random numbers.
func RandGenerate4Digits() (string, error) {
	return password.Generate(4, 4, 0, false, false)
}

//// randGenerate4Digits generates 4 random digits.
//func randGenerate4Digits() (string, error) {
//	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
//	code := make([]byte, 4)
//	for i := range code {
//		code[i] = '0' + byte(seededRand.Intn(10))
//	}
//	return string(code), nil
//}

// RandGenerate16 generates a 16-character random string.
func RandGenerate16() (string, error) {
	return password.Generate(16, 0, 0, false, false)
}

// RandGenerate20 generates a 20-character random string.
func RandGenerate20() (string, error) {
	return password.Generate(20, 0, 0, false, false)
}

// RandGenerate32 generates a 32-character random string.
func RandGenerate32() (string, error) {
	return password.Generate(32, 0, 0, false, false)
}

// RandGenerate64 generates a 64-character random string.
func RandGenerate64() (string, error) {
	return password.Generate(64, 0, 0, false, true)
}

// RandGenerate generates a random string of specified length.
func RandGenerate(length, numDigits, numSymbols int, noUpper, allowRepeat bool) (string, error) {
	return password.Generate(length, numDigits, numSymbols, noUpper, allowRepeat)
}

// HandleError is a helper function to handle errors in a user-defined way.
func HandleError(err error) {
	// Log the error, send it to a monitoring service, or handle it as needed.
	fmt.Println("Error generating random string:", err)
}

// TryRandGenerate4Digits tries to generate a 4-digit random string and handles errors.
func TryRandGenerate4Digits() string {
	pass, err := RandGenerate4Digits()
	if err != nil {
		HandleError(err)
		return ""
	}
	return pass
}

// MustRandGenerate4Digits generates a 4-digit random string and panics with a descriptive error message if there's an error.
func MustRandGenerate4Digits() string {
	pass, err := RandGenerate4Digits()
	if err != nil {
		panic(fmt.Sprintf("error generating 4-digit random string: %v", err))
	}
	return pass
}

// MustRandGenerate16 generates a 16-character random string and panics with a descriptive error message if there's an error.
func MustRandGenerate16() string {
	pass, err := RandGenerate16()
	if err != nil {
		panic(fmt.Sprintf("error generating 16-character random string: %v", err))
	}
	return pass
}

// MustRandGenerate20 generates a 20-character random string and panics with a descriptive error message if there's an error.
func MustRandGenerate20() string {
	pass, err := RandGenerate20()
	if err != nil {
		panic(fmt.Sprintf("error generating 20-character random string: %v", err))
	}
	return pass
}

// MustRandGenerate32 generates a 32-character random string and panics with a descriptive error message if there's an error.
func MustRandGenerate32() string {
	pass, err := RandGenerate32()
	if err != nil {
		panic(fmt.Sprintf("error generating 32-character random string: %v", err))
	}
	return pass
}

// MustRandGenerate64 generates a 64-character random string and panics with a descriptive error message if there's an error.
func MustRandGenerate64() string {
	pass, err := RandGenerate64()
	if err != nil {
		panic(fmt.Sprintf("error generating 64-character random string: %v", err))
	}
	return pass
}

// MustRandGenerate generates a random string of specified length and panics with a descriptive error message if there's an error.
func MustRandGenerate(length, numDigits, numSymbols int, noUpper, allowRepeat bool) string {
	pass, err := password.Generate(length, numDigits, numSymbols, noUpper, allowRepeat)
	if err != nil {
		panic(fmt.Sprintf("error generating random string: %v", err))
	}
	return pass
}

// GenerateRandomInt100KTo1B generates a random integer between 100,000 and 999,999,999.
func GenerateRandomInt100KTo1B() int {
	return GenerateRandomIntWithOptions(100000, 999999999)
}

// GenerateRandomInt100KTo1M generates a random integer between 100,000 and 999,999.
func GenerateRandomInt100KTo1M() int {
	return GenerateRandomIntWithOptions(100000, 999999)
}

// GenerateRandomIntWithOptions generates a random integer between min and max.
func GenerateRandomIntWithOptions(min int, max int) int {
	if min < 0 {
		min = 0
	}
	if max <= min {
		max = min + 999 // Adjust max to be within a **reasonable** range
	}

	// Use rand.New with a new source instead of rand.Seed
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Intn(max-min+1) + min
}
