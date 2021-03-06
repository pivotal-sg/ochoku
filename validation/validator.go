package validation

import (
	"errors"
	"fmt"
	"reflect"
)

// Validator is a function that will return a a ValidationError if the fieldName
// fails whatever validation tests it has.  Potentially other errors (like "DB Not Found"
// or something).  It returns nil if there are no errors
type Validator func(i interface{}, fieldName string) error

// Validation holds a field level valdiation error.  It also implements the error
// interface.
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// Error stringifies the ValidationError, implements the error interface
func (err ValidationError) Error() string {
	return fmt.Sprintf("%s is %s", err.Field, err.Message)
}

type Validation struct {
	V         Validator
	FieldName string
}

type Validations []Validation

// getConcreteValue will get the reflect.Value associated with the field
// or return an error if it doesn't exist.  It will dereference any pointers
// as well
func getConcreteValue(i interface{}, fieldName string) (reflect.Value, error) {
	v := reflect.ValueOf(i)
	v = reflect.Indirect(v)
	fieldValue := v.FieldByName(fieldName)
	fieldValue = reflect.Indirect(fieldValue)
	if !fieldValue.IsValid() {
		return fieldValue, fmt.Errorf("No field of name '%s' or nil pointer", fieldName)
	}

	return fieldValue, nil
}

// Ensure that the field of name `fieldName` is both a string, and
// not blank ("")
func ValidateStringNotBlank(i interface{}, fieldName string) error {
	fieldValue, err := getConcreteValue(i, fieldName)
	if err != nil {
		return err
	}

	if fieldValue.Type() != reflect.TypeOf("") {
		return errors.New("Not a String")
	}
	if fieldValue.String() == "" {
		return ValidationError{Field: fieldName, Message: "missing"}
	}
	return nil
}

// Ensure that the field of name `fieldName` is both a []byte, and
// not empty
func ValidateByteSliceNotEmpty(i interface{}, fieldName string) error {
	fieldValue, err := getConcreteValue(i, fieldName)
	if err != nil {
		return err
	}
	if fieldValue.Type() != reflect.TypeOf([]byte{}) {
		return errors.New("Not a []byte")
	}
	if len(fieldValue.Bytes()) == 0 {
		return ValidationError{Field: fieldName, Message: "missing"}
	}
	return nil
}

// Validate that the passed in interface is good.  It must be a concrete type;
// not a pointer or reference type
func (v Validations) Validate(i interface{}) []error {
	errorSlice := make([]error, 0, 0)
	// Grab the underlying value, and dereference it if that is needed

	for _, validation := range v {
		if err := validation.V(i, validation.FieldName); err != nil {
			errorSlice = append(errorSlice, err)
		}
	}

	if len(errorSlice) != 0 {
		return errorSlice
	}
	return nil
}
