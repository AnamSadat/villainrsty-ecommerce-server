package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"strings"

	"villainrsty-ecommerce-server/internal/adapters/http/lib/httpx"
	"villainrsty-ecommerce-server/internal/core/auth/ports"
	"villainrsty-ecommerce-server/internal/core/shared/models"
)

func AuthJWT(jwtService ports.JWTService, logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := extractToken(r)
			if token == "" {
				if logger != nil {
					logger.Warn("missing authorization token", "path", r.URL.Path, "method", r.Method)
				}
				httpx.Error(w, http.StatusUnauthorized, "missing authorization token", "MISSING_AUTHORIZATION_TOKEN")
				return
			}

			user, err := jwtService.ValidateToken(token)
			if err != nil {
				if logger != nil {
					logger.Warn("invalid token", "error", err.Error(), "path", r.URL.Path, "method", r.Method)
				}
				httpx.ErrorWithDetails(w, http.StatusUnauthorized, "invalid token", "INVALID_TOKEN", err.Error())
				return
			}

			if logger != nil {
				logger.Info("user authenticated", "user_id", user.ID.String(), "path", r.URL.Path, "method", r.Method)
			}

			ctx := context.WithValue(r.Context(), "user", user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func extractToken(r *http.Request) string {
	auth := r.Header.Get("Authorization")
	if auth == "" {
		return ""
	}

	parts := strings.Split(auth, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return ""
	}

	return parts[1]
}

func GetUserFromContext(r http.Request) *models.User {
	user, ok := r.Context().Value("user").(*models.User)
	if !ok {
		return nil
	}

	return user
}
