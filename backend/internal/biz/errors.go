package biz

import "errors"

// Returned by repositories.
var (
	ErrNotFound  = errors.New("not found")
	ErrDuplicate = errors.New("duplicate")
)

// Business errors, mapped to HTTP status codes by the handler layer.
var (
	ErrEmailTaken      = errors.New("email already registered")
	ErrInvalidCode     = errors.New("login code does not match")
	ErrTooManyAttempts = errors.New("too many failed login attempts")
	ErrUnauthenticated = errors.New("not authenticated")
)

// ValidationError's Message is safe to show the user.
type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string { return e.Message }

func invalid(msg string) error { return &ValidationError{Message: msg} }
