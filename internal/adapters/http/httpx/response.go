package httpx

import (
	"encoding/json"
	"net/http"
)

type (
	ErrorInfo struct {
		Code    string `json:"code"`
		Details any    `json:"details,omitempty"`
	}

	Meta struct {
		Page      int `json:"page"`
		Limit     int `json:"limit"`
		Total     int `json:"total"`
		TotalPage int `json:"total_page"`
	}

	FieldError struct {
		Field   string `json:"field"`
		Message string `json:"message"`
	}

	BaseResponse[T any] struct {
		Success bool       `json:"success"`
		Message string     `json:"message"`
		Data    T          `json:"data,omitempty"`
		Error   *ErrorInfo `json:"error,omitempty"`
		Meta    *Meta      `json:"meta,omitempty"`
	}
)

func JSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func Success[T any](w http.ResponseWriter, status int, msg string, data T) {
	resp := BaseResponse[T]{
		Success: true,
		Message: msg,
		Data:    data,
	}

	JSON(w, status, resp)
}

func SuccessWithMeta[T any](w http.ResponseWriter, status int, msg string, data T, meta *Meta) {
	resp := BaseResponse[T]{
		Success: true,
		Message: msg,
		Data:    data,
		Meta:    meta,
	}

	JSON(w, status, resp)
}

func Error(w http.ResponseWriter, status int, msg string, code string) {
	resp := BaseResponse[any]{
		Success: false,
		Message: msg,
		Data:    nil,
		Error: &ErrorInfo{
			Code: code,
		},
	}

	JSON(w, status, resp)
}

func ErrorWithDetails(w http.ResponseWriter, status int, msg string, code string, details any) {
	resp := BaseResponse[any]{
		Success: false,
		Message: msg,
		Data:    nil,
		Error: &ErrorInfo{
			Code:    code,
			Details: details,
		},
	}

	JSON(w, status, resp)
}

func ValidationError(w http.ResponseWriter, filedErrors []FieldError) {
	resp := BaseResponse[any]{
		Success: false,
		Message: "validation failed",
		Data:    nil,
		Error: &ErrorInfo{
			Code:    "VALIDATION_ERROR",
			Details: filedErrors,
		},
	}

	JSON(w, http.StatusBadRequest, resp)
}

func DecodeJSON[T any](w http.ResponseWriter, r *http.Request, dst *T) bool {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(dst); err != nil {
		ErrorWithDetails(w, http.StatusBadRequest, "Invalid JSON", "INVALID_JSON", err.Error())

		return false
	}

	return true
}
