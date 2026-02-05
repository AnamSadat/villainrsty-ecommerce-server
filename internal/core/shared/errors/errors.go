package errors

import (
	"errors"
	"fmt"
)

type Kind string

const (
	ErrNotFound     Kind = "not_found"
	ErrValidation   Kind = "validation"
	ErrUnauthorized Kind = "unauthorized"
	ErrForbidden    Kind = "forbidden"
	ErrConflict     Kind = "conflict"
	ErrInternal     Kind = "internal"
)

type AppError struct {
	Kind    Kind
	Message string
	Cause   error
	Fields  map[string]string
}

func (e *AppError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %s: %v", e.Kind, e.Message, e.Cause)
	}

	return fmt.Sprintf("%s: %s", e.Kind, e.Message)
}

func (e *AppError) UnWrap() error { return e.Cause }

// New bikin AppError tanpa Cause
func New(kind Kind, msg string) *AppError {
	return &AppError{Kind: kind, Message: msg}
}

// Wrap bikin AppError dengan cause (biar root error tidak hilang)
func Wrap(kind Kind, msg string, cause error) *AppError {
	return &AppError{Kind: kind, Message: msg, Cause: cause}
}

// Validation helper (biar enak bikin error validasi fields)
func Validation(msg string, field map[string]string) *AppError {
	return &AppError{Message: msg, Fields: field}
}

// IsKind buat ngecek kategori error
func IsKind(err error, kind Kind) bool {
	var ae *AppError
	if errors.As(err, &ae) {
		return ae.Kind == kind
	}

	return false
}

// AsAppError ambil AppError kalau ada
func AsAppError(err error) (*AppError, bool) {
	var ae *AppError
	if errors.As(err, &ae) {
		return ae, true
	}

	return nil, false
}
