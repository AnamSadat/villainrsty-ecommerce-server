package httpx

import (
	"net/http"

	"villainrsty-ecommerce-server/internal/core/shared/errors"
)

type Logger interface {
	Error(msg string, args ...interface{})
	Warn(msg string, args ...interface{})
}

func HandleError(w http.ResponseWriter, err error, log Logger) {
	appErr, ok := errors.AsAppError(err)
	if !ok {
		if log != nil {
			log.Error("internal server error", "error", err.Error())
		}
		ErrorWithDetails(
			w,
			http.StatusInternalServerError,
			"Internal Server Error",
			"INTERNAL_ERROR",
			err.Error(),
		)
		return
	}

	switch appErr.Kind {

	case errors.ErrNotFound:
		if log != nil {
			log.Warn("resource not found", "error", appErr.Error())
		}
		Error(w, http.StatusNotFound, "Resource not found", "NOT_FOUND")

	case errors.ErrBadRequest:
		if log != nil {
			log.Warn("bad request", "error", appErr.Error())
		}

		msg := "Bad request"
		if appErr.Message != "" {
			msg = appErr.Message
		}

		Error(w, http.StatusBadRequest, msg, "BAD_REQUEST")

	case errors.ErrValidation:
		if log != nil {
			log.Warn("validation error", "error", appErr.Error())
		}

		// Mapping field errors
		fieldErrors := make([]FieldError, 0)
		for field, msg := range appErr.Fields {
			fieldErrors = append(fieldErrors, FieldError{
				Field:   field,
				Message: msg,
			})
		}
		ValidationError(w, fieldErrors)

	case errors.ErrUnauthorized:
		if log != nil {
			log.Warn("unauthorized access", "error", appErr.Error())
		}
		Error(w, http.StatusUnauthorized, "Invalid credentials", "UNAUTHORIZED")

	case errors.ErrForbidden:
		if log != nil {
			log.Warn("forbidden access", "error", appErr.Error())
		}
		msg := "Access denied"
		if appErr.Message != "" {
			msg = appErr.Message
		}
		Error(w, http.StatusForbidden, msg, "FORBIDDEN")

	case errors.ErrConflict:
		if log != nil {
			log.Warn("resource conflict", "error", appErr.Error())
		}
		ErrorWithDetails(w, http.StatusConflict, "Resource already exists", "CONFLICT", err.Error())

	default:
		// Default ke 500 jika Kind tidak dikenali
		if log != nil {
			log.Error("server error response", "kind", appErr.Kind, "message", appErr.Error())
		}
		ErrorWithDetails(
			w,
			http.StatusInternalServerError,
			"Internal server error",
			"INTERNAL_ERROR",
			appErr.Message,
		)
	}
}
