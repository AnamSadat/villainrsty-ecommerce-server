package handler

import (
	"log"
	"net/http"

	"villainrsty-ecommerce-server/internal/adapters/http/auth/models"
	"villainrsty-ecommerce-server/internal/adapters/http/httpx"
	"villainrsty-ecommerce-server/internal/core/auth/ports"
	"villainrsty-ecommerce-server/internal/core/shared/errors"

	sharedModel "villainrsty-ecommerce-server/internal/core/shared/models"
)

type AuthHandler struct {
	authService ports.AuthService
}

func NewAuthHandler(service ports.AuthService) *AuthHandler {
	return &AuthHandler{authService: service}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if !httpx.DecodeJSON(w, r, &req) {
		return
	}

	if err := req.Validate(); err != nil {
		return
	}

	user, token, err := h.authService.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		hanlderError(w, err)
		return
	}

	res := models.LoginResponse{
		User:  mapUserToDTO(user),
		Token: token,
	}

	httpx.JSON(w, http.StatusOK, res)
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req models.RegisterRequest

	if !httpx.DecodeJSON(w, r, &req) {
		return
	}

	log.Printf("DEBUG: Register request - Email: %s, Name: %s, Password: %s", req.Email, req.Name, req.Password)

	if err := req.Validate(); err != nil {
		log.Printf("DEBUG: Validation error: %v", err)
		return
	}

	user, err := h.authService.Register(r.Context(), req.Email, req.Password, req.Name)
	if err != nil {
		log.Printf("DEBUG: Service error: %v", err)
		hanlderError(w, err)
		return
	}

	res := models.RegisterResponse{
		User: mapUserToDTO(user),
	}

	httpx.JSON(w, http.StatusOK, res)
}

func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var req models.RefreshTokenRequest
	if !httpx.DecodeJSON(w, r, &req) {
		return
	}

	if err := req.Validate(); err != nil {
		return
	}

	token, err := h.authService.RefreshToken(r.Context(), req.RefreshToken)
	if err != nil {
		hanlderError(w, err)
		return
	}

	res := models.RefreshTokenResponse{
		Token: token,
	}

	httpx.JSON(w, http.StatusOK, res)
}

func hanlderError(w http.ResponseWriter, err error) {
	appErr, ok := errors.AsAppError(err)
	if !ok {
		httpx.Error(w, http.StatusInternalServerError, "internal server error")
	}

	switch appErr.Kind {
	case errors.ErrNotFound:
		httpx.Error(w, http.StatusNotFound, appErr.Message)
	case errors.ErrValidation:
		httpx.JSON(w, http.StatusBadRequest, map[string]any{
			"error": appErr.Message,
			"field": appErr.Fields,
		})
	case errors.ErrUnauthorized:
		httpx.Error(w, http.StatusUnauthorized, appErr.Message)
	case errors.ErrForbidden:
		httpx.Error(w, http.StatusForbidden, appErr.Message)
	case errors.ErrConflict:
		httpx.Error(w, http.StatusConflict, appErr.Message)
	default:
		httpx.Error(w, http.StatusInternalServerError, "internal server error")

	}
}

func mapUserToDTO(user *sharedModel.User) models.UserDTO {
	return models.UserDTO{
		ID:    user.ID.String(),
		Email: user.Email,
		Name:  user.Name,
	}
}
