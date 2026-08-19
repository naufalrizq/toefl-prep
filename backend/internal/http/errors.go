package http

import (
	"errors"
	"net/http"
)

var (
	ErrUnauthorized = errors.New("auth_required")
	ErrForbidden    = errors.New("forbidden")
	ErrNotFound     = errors.New("not_found")
	ErrConflict     = errors.New("conflict")
	ErrValidation   = errors.New("validation_failed")
	ErrRateLimited  = errors.New("rate_limited")
)

type codedError struct {
	status  int
	code    string
	message string
}

func (e *codedError) Error() string { return e.message }

// NewError builds an error with a stable code and HTTP status.
func NewError(status int, code, message string) error {
	return &codedError{status: status, code: code, message: message}
}

// Map turns any error into (httpStatus, errorCode, message).
func Map(err error) (int, string, string) {
	var ce *codedError
	if errors.As(err, &ce) {
		return ce.status, ce.code, ce.message
	}
	switch {
	case errors.Is(err, ErrUnauthorized):
		return http.StatusUnauthorized, "auth_required", "authentication required"
	case errors.Is(err, ErrForbidden):
		return http.StatusForbidden, "forbidden", "you do not have permission"
	case errors.Is(err, ErrNotFound):
		return http.StatusNotFound, "not_found", "resource not found"
	case errors.Is(err, ErrConflict):
		return http.StatusConflict, "conflict", "conflicts with existing state"
	case errors.Is(err, ErrValidation):
		return http.StatusUnprocessableEntity, "validation_failed", err.Error()
	case errors.Is(err, ErrRateLimited):
		return http.StatusTooManyRequests, "rate_limited", "too many requests"
	default:
		return http.StatusInternalServerError, "internal", "internal server error"
	}
}