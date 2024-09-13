package validation

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func compositeErrors(err error) error {
	if err == nil {
		return nil
	}

	var errorMessages []string
	for _, err := range err.(validator.ValidationErrors) {
		errorMessages = append(errorMessages,
			fmt.Sprintf("field %s failed validation: %s", err.Field(), err.Tag()),
		)
	}

	return &ValidateError{msg: strings.Join(errorMessages, "; ")}
}

func ValidateStruct(s any) error {
	if err := validate.Struct(s); err != nil {
		return compositeErrors(err)
	}

	return nil
}

func ValidateOneOf[T ~string](filter []T, value string, fieldName string) error {
	for _, v := range filter {
		if T(value) == v {
			return nil
		}
	}

	return NewValidateError(fmt.Sprintf("invalid %s format, must be in: %v", fieldName, filter))
}

func ValidateOneOfSlice[T ~string](filter []T, values []string, fieldName string) error {
	for _, v := range values {
		if err := ValidateOneOf(filter, v, fieldName); err != nil {
			return err
		}
	}

	return nil
}

func ValidateHTTPMethod(r *http.Request, expectedMethod string) error {
	if r.Method != expectedMethod {
		return fmt.Errorf("invalid HTTP method: expected %s, got %s", expectedMethod, r.Method)
	}
	return nil
}
