package model

// ValidationResult is non-nil in case of errors
// This could contain lists of Warnings and Errors
type ValidationResult error

type Validatable interface {
	// Validate run checks on the struct
	// Strict mode means that values that are technically valid but not sane are
	// treated as errors
	Validate(strict bool) ValidationResult
}
