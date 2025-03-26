package aerr

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNew(t *testing.T) {
	err := New("test error")
	assert.NotNil(t, err)
	assert.Equal(t, "test error", err.Error())
}

func TestNewf(t *testing.T) {
	err := Newf("test error: %d", 123)
	assert.NotNil(t, err)
	assert.Equal(t, "test error: 123", err.Error())
}

func TestNewError(t *testing.T) {
	baseErr := fmt.Errorf("base error")
	err := NewError(baseErr)
	assert.NotNil(t, err)
	assert.Equal(t, "base error", err.Error())

	nilErr := NewError(nil)
	assert.Nil(t, nilErr)
}

func TestNewErrorFromString(t *testing.T) {
	errMsg := "test error"
	aerr := NewErrorFromString(errMsg)
	if aerr.Error() != errMsg {
		t.Errorf("NewErrorFromString did not create the error correctly")
	}
}

func TestIsNil(t *testing.T) {
	var aerr *Error
	if !aerr.IsNil() {
		t.Errorf("IsNil should return true for nil *Error")
	}

	aerr = NewError(nil)
	if !aerr.IsNil() {
		t.Errorf("IsNil should return true for nil embedded error")
	}
}

func TestMarshalJSON(t *testing.T) {
	aerr := NewErrorFromString("test error")
	data, err := json.Marshal(aerr)
	if err != nil {
		t.Fatal("Failed to marshal JSON:", err)
	}
	if string(data) != "\"test error\"" {
		t.Errorf("MarshalJSON did not return the correct JSON representation")
	}
}

func TestUnmarshalJSON(t *testing.T) {
	var aerr Error
	err := json.Unmarshal([]byte("\"test error\""), &aerr)
	if err != nil {
		t.Fatal("Failed to unmarshal JSON:", err)
	}
	if aerr.Error() != "test error" {
		t.Errorf("UnmarshalJSON did not set the correct error message")
	}
}

func TestToError(t *testing.T) {
	err := errors.New("test error")
	aerr := NewError(err)
	if aerr.ToError() != err {
		t.Errorf("ToError did not return the original error")
	}
}

func TestString(t *testing.T) {
	errMsg := "test error"
	aerr := NewErrorFromString(errMsg)
	if aerr.String() != errMsg {
		t.Errorf("String did not return the correct error message")
	}
}

func TestUnwrap(t *testing.T) {
	err := errors.New("test error")
	aerr := NewError(err)
	if aerr.Unwrap() != err {
		t.Errorf("Unwrap did not return the original error")
	}
}

func TestIsEqual(t *testing.T) {
	err := errors.New("test error")
	aerr := NewError(err)
	otherErr := fmt.Errorf("test error")

	if !aerr.IsEqual(err) {
		t.Errorf("IsEqual should return true for the same error")
	}

	if !aerr.IsEqual(otherErr) {
		t.Errorf("IsEqual should return true for different error with same message")
	}

	if aerr.IsEqual(errors.New("other error")) {
		t.Errorf("IsEqual should return false for different errors")
	}
}

// TestError is a struct for testing JSON marshaling of Error types.
type TestError struct {
	Err         Error  `json:"err"`
	ErrEmpty    *Error `json:"errEmpty"`
	ErrPtr      *Error `json:"errPtr"`
	ErrNil      *Error `json:"errNil"`
	ErrPtrEmpty *Error `json:"errPtrEmpty"`
}

// TestError_MarshalJSON tests the marshaling of Error types to JSON.
func TestError_MarshalJSON(t *testing.T) {
	jsonRoot := `{"err":"this is an error","errEmpty":"","errPtr":"this is a ptr error","errNil": null, "errPtrEmpty":""}`

	testErr := &TestError{}
	err := json.Unmarshal([]byte(jsonRoot), testErr)
	if err != nil {
		t.Fatal(err)
	}

	if testErr.Err.Error() != "this is an error" {
		t.Fatal("Err field did not unmarshal correctly")
	}

	if testErr.ErrEmpty != nil && testErr.ErrEmpty.Error() != "" {
		t.Fatal("ErrEmpty field should be nil or empty")
	}

	if testErr.ErrPtr == nil || testErr.ErrPtr.Error() != "this is a ptr error" {
		t.Fatal("ErrPtr field did not unmarshal correctly")
	}

	if testErr.ErrNil != nil {
		t.Fatal("ErrNil field should be nil")
	}

	if testErr.ErrPtrEmpty != nil && testErr.ErrPtrEmpty.Error() != "" {
		t.Fatal("ErrPtrEmpty field should be nil or empty")
	}
}

// getErrorNil returns a nil *Error.
func getErrorNil() *Error {
	return NewError(nil)
}

// getErrorNilDefault returns a nil error.
func getErrorNilDefault() error {
	return NewError(nil)
}

// TestError_Assignment tests the assignment of errors and nil values.
func TestError_Assignment(t *testing.T) {
	err1 := NewError(fmt.Errorf("this is a test"))
	if err1 == nil || err1.Error() != "this is a test" {
		t.Fatalf("err1 should equal 'this is a test', got: %v", err1)
	}

	err2 := NewError(nil)
	if err2 != nil {
		t.Fatalf("err2 should be nil, got: %v", err2)
	}

	err3 := getErrorNil()
	if err3 != nil {
		t.Fatalf("err3 should be nil, got: %v", err3)
	}

	err4 := getErrorNilDefault()
	if err4 != nil && err4.Error() != "" {
		t.Fatalf("err4 should be empty, got: %v", err4)
	}
}
