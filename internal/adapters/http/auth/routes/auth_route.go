package routes

import (
	"villainrsty-ecommerce-server/internal/adapters/http/auth/handler"

	"github.com/go-chi/chi/v5"
)

func RegisterRoute(r chi.Router, handler *handler.AuthHandler) {
	r.Route("/auth", func(r chi.Router) {
		r.Post("/login", handler.Login)
		r.Post("/login-2fa", handler.Login2FA)
		r.Post("/verify-login-2fa", handler.VerifyLogin2FA)
		r.Post("/register", handler.Register)
		r.Post("/refresh", handler.RefreshToken)
		r.Post("/logout", handler.Logout)
		r.Post("/forgot-password", handler.ForgotPassword)
		r.Post("/reset-password", handler.ResetPassword)
	})
}
