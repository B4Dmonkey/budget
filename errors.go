package main

// ValidationError is an error type that represents an error that occurred when checking inputs
type ValidationError struct{ Message string }

// Implement the Error interface for the ValidationError type
func (e *ValidationError) Error() string { return e.Message }

// newValidationError creates a new ValidationError with the given message
func newValidationError(message string) error { return &ValidationError{Message: message} }

// ParseError is an error type that represents an error that occurred when parsing data
type ParseError struct{ Message string }

// Implement the Error interface for the ParseError type
func (e *ParseError) Error() string { return e.Message }

// newParseError creates a new ParseError with the given message
func newParseError(message string) error { return &ParseError{Message: message} }
