package errors

import (
	"errors"

	stackTraceError "github.com/pkg/errors"
)

var EnableStackTrace = true

var New = func(text string) error {
	if EnableStackTrace {
		return stackTraceError.New(text)
	}

	return errors.New(text)
}
