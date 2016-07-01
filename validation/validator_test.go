package validation_test

import (
	"testing"

	. "github.com/pivotal-sg/ochoku/validation"
)

func TestValidStructInputToNonBlankString(t *testing.T) {
	v := struct {
		name string
	}{"hello world"}

	if err := ValidateStringNotBlank(v, "name"); err != nil {
		t.Errorf("Error should be nil, was '%v'", err)
	}

}

func TestValidStructInputToNonEmptyByteSlice(t *testing.T) {
	v := struct {
		data []byte
	}{[]byte("hello world")}

	if err := ValidateByteSliceNotEmpty(v, "data"); err != nil {
		t.Errorf("Error should be nil, was '%v'", err)
	}

}

func TestNoPanicOnBadFieldNameBlankStringValidator(t *testing.T) {
	v := struct {
		field string
	}{"test"}

	err := ValidateStringNotBlank(v, "wrongName")

	expectedMessage := "No field of name 'wrongName' or nil pointer"

	if err.Error() != expectedMessage {
		t.Errorf("Error string should '%s', was '%s'", expectedMessage, err.Error())
	}

}

func TestNoPanicOnBadFieldNameEmptyByteSliceValidator(t *testing.T) {
	v := struct {
		field string
	}{"test"}

	err := ValidateByteSliceNotEmpty(v, "wrongName")

	expectedMessage := "No field of name 'wrongName' or nil pointer"

	if err.Error() != expectedMessage {
		t.Errorf("Error string should '%s', was '%s'", expectedMessage, err.Error())
	}

}

func TestNoPanicOnPointerValueBlankStringValidator(t *testing.T) {
	var value string = "test"
	v := struct {
		field *string
	}{&value}

	err := ValidateStringNotBlank(v, "field")

	if err != nil {
		t.Errorf("Error should be '%v', was '%v'", nil, err)
	}

}

func TestNoPanicOnPointerByteSliceValidator(t *testing.T) {
	var value []byte = []byte{1, 2, 3}
	v := struct {
		field *[]byte
	}{&value}

	err := ValidateByteSliceNotEmpty(v, "field")

	if err != nil {
		t.Errorf("Error should be '%v', was '%v'", nil, err)
	}

}

func TestNoPanicOnNilValueBlankStringValidator(t *testing.T) {
	v := struct {
		field *string
	}{nil}

	err := ValidateStringNotBlank(v, "field")

	expectedMessage := "No field of name 'field' or nil pointer"

	if err.Error() != expectedMessage {
		t.Errorf("Error string should '%s', was '%s'", expectedMessage, err.Error())
	}

}

func TestNoPanicOnNilByteSliceValidator(t *testing.T) {
	v := struct {
		field *[]byte
	}{nil}

	err := ValidateByteSliceNotEmpty(v, "field")

	expectedMessage := "No field of name 'field' or nil pointer"

	if err.Error() != expectedMessage {
		t.Errorf("Error string should '%s', was '%s'", expectedMessage, err.Error())
	}

}
