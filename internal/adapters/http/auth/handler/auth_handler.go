package handler

import (
	"fmt"
	"log/slog"
	"net/http"

	"villainrsty-ecommerce-server/internal/adapters/http/auth/models"
	"villainrsty-ecommerce-server/internal/adapters/http/lib/httpx"
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
		httpx.HandleError(w, err, h.logger)
		return
	}

	user, accessToken, refreshToken, err := h.authService.Login(r.Context(), req.Email, req.Password, req.RememberMe)
	if err != nil {
		h.logger.Warn("login failed", "email", req.Email, "error", err.Error())
		httpx.HandleError(w, err, h.logger)
		return
	}

	resp := models.LoginResponse{
		User:         mapUserToDTO(user),
		Token:        accessToken,
		RefreshToken: refreshToken,
	}

	httpx.Success(w, http.StatusOK, "Login succcessfully", resp)
}

func (h *AuthHandler) Login2FA(w http.ResponseWriter, r *http.Request) {
	var req models.Login2FARequest
	if !httpx.DecodeJSON(w, r, &req) {
		return
	}

	if err := req.Validate(); err != nil {
		httpx.HandleError(w, err, h.logger)
		return
	}

	challengeID, err := h.authService.LoginWith2FA(r.Context(), req.Email, req.Password, req.RememberMe)
	if err != nil {
		httpx.HandleError(w, err, h.logger)
		return
	}

	resp := models.Login2FAResponse{ChallengeID: challengeID}

	httpx.Success(w, http.StatusOK, "OTP send to email", resp)
}

func (h *AuthHandler) VerifyLogin2FA(w http.ResponseWriter, r *http.Request) {
	var req models.VerifyLogin2FARequest
	if !httpx.DecodeJSON(w, r, &req) {
		return
	}
	h.logger.Info("isinya",
		"challenge_id", req.ChallengeID,
		"otp_code", req.OTPCode,
		"remember", req.RememberMe,
	)

	if err := req.Validate(); err != nil {
		h.logger.Info("error validate", "error", err)
		httpx.HandleError(w, err, h.logger)
		return
	}

	user, accessToken, refreshToken, err := h.authService.VerifyLogin2FA(r.Context(), req.ChallengeID, req.OTPCode, req.RememberMe)
	if err != nil {
		h.logger.Info("ke error verify", "error", err)
		httpx.HandleError(w, err, h.logger)
		return
	}

	resp := models.LoginResponse{
		User:         mapUserToDTO(user),
		Token:        accessToken,
		RefreshToken: refreshToken,
	}

	httpx.Success(w, http.StatusOK, "2FA verified successfully", resp)
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req models.RegisterRequest

	if !httpx.DecodeJSON(w, r, &req) {
		return
	}

	if err := req.Validate(); err != nil {
		httpx.HandleError(w, err, h.logger)
		return
	}

	user, err := h.authService.Register(r.Context(), req.Email, req.Password, req.Name)
	if err != nil {
		httpx.HandleError(w, err, h.logger)
		return
	}

	resp := models.RegisterResponse{User: mapUserToDTO(user)}

	httpx.Success(w, http.StatusOK, "User registered successfully", resp)
}

func (h *AuthHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user").(*sharedModel.User)
	if !ok || user == nil {
		httpx.HandleError(w, errors.New(errors.ErrUnauthorized, "user not found in context"), h.logger)
		return
	}

	userDTO := models.UserDTO{
		ID:    user.ID.String(),
		Email: user.Email,
		Name:  user.Name,
		Role:  user.Role,
	}

	resp := models.GetProfileResponse{User: userDTO}
	httpx.Success(w, http.StatusOK, "Get profile success", resp)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	var req models.LogoutRequest

	if !httpx.DecodeJSON(w, r, &req) {
		return
	}

	if err := req.Validate(); err != nil {
		httpx.HandleError(w, err, h.logger)
		return
	}

	if err := h.authService.Logout(r.Context(), req.RefreshToken); err != nil {
		httpx.HandleError(w, err, h.logger)
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
		httpx.HandleError(w, err, h.logger)
		return
	}

	accessToken, refreshToken, err := h.authService.RefreshToken(r.Context(), req.RefreshToken)
	if err != nil {
		fmt.Println("masuk ke sini")
		httpx.HandleError(w, err, h.logger)
		return
	}

	resp := models.RefreshTokenResponse{
		Token:        accessToken,
		RefreshToken: refreshToken,
	}

	httpx.Success(w, http.StatusOK, "Token refreshed successfully", resp)
}

func (h *AuthHandler) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	var req models.ForgotPasswordRequest

	if !httpx.DecodeJSON(w, r, &req) {
		return
	}

	if err := req.Validate(); err != nil {
		httpx.HandleError(w, err, h.logger)
		return
	}

	if err := h.authService.RequestPasswordReset(r.Context(), req.Email); err != nil {
		httpx.HandleError(w, err, h.logger)
	}

	httpx.Success(w, http.StatusOK, "Request berhasil, link akan segera dikirim", "")
}

func (h *AuthHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var req models.ResetPasswordRequest

	if !httpx.DecodeJSON(w, r, &req) {
		return
	}

	if err := req.Validate(); err != nil {
		httpx.HandleError(w, err, h.logger)
		return
	}

	if err := h.authService.ConfirmPasswordReset(r.Context(), req.Token, req.NewPassword); err != nil {
		httpx.HandleError(w, err, h.logger)
		return
	}

	httpx.Success(w, http.StatusOK, "Password berhasil direset", "")
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
