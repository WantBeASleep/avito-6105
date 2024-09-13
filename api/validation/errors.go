package validation

import (
	"errors"
	"fmt"
)

var (
	ErrParsed = errors.New("wrong data format")
)

func NewValidateError(msg string) error {
	return &ValidateError{msg: msg}
}

type ValidateError struct {
	msg string
}

func (e *ValidateError) Error() string {
	return fmt.Sprintf("validation failed: %s", e.msg)
}
