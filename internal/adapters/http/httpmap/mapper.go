package httpmap

import (
	"net/http"

	"villainrsty-ecommerce-server/internal/adapters/http/httpx"
	"villainrsty-ecommerce-server/internal/core/shared/errors"
)

type ErrorResponse struct {
	Error   string            `json:"error"`
	Message string            `json:"message,omitempty"`
	Fields  map[string]string `json:"field,omiempty"`
}

func WriteError(w http.ResponseWriter, err error) {
	// default kalo bukan error AppError
	status := http.StatusInternalServerError
	resp := ErrorResponse{Error: "internal", Message: "internal server error"}

	if ae, ok := errors.AsAppError(err); ok {
		status = toStatus(ae.Kind)
		resp = ErrorResponse{
			Error:   string(ae.Kind),
			Message: ae.Message,
			Fields:  ae.Fields,
		}
	}

	httpx.JSON(w, status, resp)
}

func toStatus(kind errors.Kind) int {
	switch kind {
	case errors.ErrValidation:
		return http.StatusBadRequest
	case errors.ErrUnauthorized:
		return http.StatusUnauthorized
	case errors.ErrForbidden:
		return http.StatusForbidden
	case errors.ErrNotFound:
		return http.StatusNotFound
	case errors.ErrConflict:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}
