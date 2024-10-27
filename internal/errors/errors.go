package errors

import (
	"strings"
)

type Verifier struct {
	errors []VerificationError
}

func (v *Verifier) That(condition bool, message string) {
	if !condition {
		v.errors = append(v.errors, VerificationError{Message: message})
	}
}

func (v *Verifier) Flush() error {
	if len(v.errors) == 0 {
		return nil
	}

	errMessages := make([]string, len(v.errors))

	for i, err := range v.errors {
		errMessages[i] = "\n\t" + err.Message
	}
	
	message := "Verification failed:" + strings.Join(errMessages, ";")
	
	return &VerificationError{Message: message}
}

// VerificationError is an error type that represents an error that occurred when checking inputs
type VerificationError struct{ Message string }

// Implement the Error interface for the ValidationError type
func (e VerificationError) Error() string { return e.Message }

// ParseError is an error type that represents an error that occurred when parsing data
type ParseError struct{ Message string }

// Implement the Error interface for the ParseError type
func (e *ParseError) Error() string { return e.Message }

// newParseError creates a new ParseError with the given message
func NewParseError(message string) error { return &ParseError{Message: message} }
