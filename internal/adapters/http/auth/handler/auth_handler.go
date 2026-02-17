package handler

import (
	"fmt"
	"log/slog"
	"net/http"

	"villainrsty-ecommerce-server/internal/adapters/http/auth/models"
	"villainrsty-ecommerce-server/internal/adapters/http/httpx"
	"villainrsty-ecommerce-server/internal/core/auth/ports"
	"villainrsty-ecommerce-server/internal/core/shared/errors"

	sharedModel "villainrsty-ecommerce-server/internal/core/shared/models"
)

type AuthHandler struct {
	authService ports.AuthService
	logger      *slog.Logger
}

func NewAuthHandler(service ports.AuthService, logger *slog.Logger) *AuthHandler {
	return &AuthHandler{authService: service, logger: logger}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if !httpx.DecodeJSON(w, r, &req) {
		h.logger.Warn("failed to decode login json body")
		return
	}

	if err := req.Validate(); err != nil {
		h.hanlderError(w, err)
		return
	}

	user, accessToken, refreshToken, err := h.authService.Login(r.Context(), req.Email, req.Password, req.RememberMe)
	if err != nil {
		h.logger.Warn("login failed", "email", req.Email, "error", err.Error())
		h.hanlderError(w, err)
		return
	}

	resp := models.LoginResponse{
		User:         mapUserToDTO(user),
		Token:        accessToken,
		RefreshToken: refreshToken,
	}

	httpx.Success(w, http.StatusOK, "Login susccessfully", resp)
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req models.RegisterRequest

	if !httpx.DecodeJSON(w, r, &req) {
		return
	}

	if err := req.Validate(); err != nil {
		h.hanlderError(w, err)
		return
	}

	user, err := h.authService.Register(r.Context(), req.Email, req.Password, req.Name)
	if err != nil {
		h.hanlderError(w, err)
		return
	}

	resp := models.RegisterResponse{
		User: mapUserToDTO(user),
	}

	httpx.Success(w, http.StatusOK, "User registered successfully", resp)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	var req models.LogoutRequest

	if !httpx.DecodeJSON(w, r, &req) {
		return
	}

	if err := req.Validate(); err != nil {
		h.hanlderError(w, err)
		return
	}

	err := h.authService.Logout(r.Context(), req.RefreshToken)
	if err != nil {
		h.hanlderError(w, err)
		return
	}

	httpx.Success(w, http.StatusOK, "Successfully logout", "")
}

func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var req models.RefreshTokenRequest
	// h.logger.Info("req di handler", "validate diluar", req.Validate())

	if !httpx.DecodeJSON(w, r, &req) {
		return
	}
	h.logger.Info("req di handler", "refresh token", req.RefreshToken)

	if err := req.Validate(); err != nil {
		h.logger.Info("req di handler", "validate di dalam", req.Validate())
		h.hanlderError(w, err)
		return
	}

	accessToken, refreshToken, err := h.authService.RefreshToken(r.Context(), req.RefreshToken)
	if err != nil {
		fmt.Println("masuk ke sini")
		h.hanlderError(w, err)
		return
	}

	resp := models.RefreshTokenResponse{
		Token:        accessToken,
		RefreshToken: refreshToken,
	}

	httpx.Success(w, http.StatusOK, "Token refreshed successfully", resp)
}

func (h *AuthHandler) hanlderError(w http.ResponseWriter, err error) {
	appErr, ok := errors.AsAppError(err)
	if !ok {
		h.logger.Error("internal server error", "error", err.Error())
		httpx.ErrorWithDetails(
			w,
			http.StatusInternalServerError,
			"Internal Server Error",
			"INTERNAL_ERROR",
			err.Error(),
		)
		return
	}

	switch appErr.Kind {
	case errors.ErrNotFound, errors.ErrValidation, errors.ErrUnauthorized, errors.ErrForbidden, errors.ErrConflict:
		h.logger.Warn("client error response",
			"kind", appErr.Kind,
			"message", appErr.Error(),
		)
	default:
		h.logger.Error("server error response",
			"kind", appErr.Kind,
			"message", appErr.Error(),
		)
	}

	switch appErr.Kind {
	case errors.ErrNotFound:
		httpx.Error(w, http.StatusNotFound, "Resource not found", "NOT_FOUND")
	case errors.ErrValidation:
		fieldErrors := make([]httpx.FieldError, 0)
		for field, msg := range appErr.Fields {
			fieldErrors = append(fieldErrors, httpx.FieldError{
				Field:   field,
				Message: msg,
			})
		}
		httpx.ValidationError(w, fieldErrors)
	case errors.ErrUnauthorized:
		httpx.Error(w, http.StatusUnauthorized, "invalid credentials", "UNAUTHORIZED")
	case errors.ErrForbidden:
		httpx.Error(w, http.StatusForbidden, "access denied", "FORBIDDEN")
	case errors.ErrConflict:
		httpx.ErrorWithDetails(w, http.StatusConflict, "resource already exists", "CONFLICT", err.Error())
	default:
		httpx.ErrorWithDetails(
			w,
			http.StatusInternalServerError,
			"internal server error",
			"INTERNAL_ERROR",
			appErr.Message,
		)

	}
}

func mapUserToDTO(user *sharedModel.User) models.UserDTO {
	return models.UserDTO{
		ID:    user.ID.String(),
		Email: user.Email,
		Name:  user.Name,
	}
}

func maskToken(t string) string {
	if len(t) <= 16 {
		return t
	}
	return t[:8] + "..." + t[len(t)-8:]
}
