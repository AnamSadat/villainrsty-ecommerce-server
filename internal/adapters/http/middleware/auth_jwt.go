package middleware

import (
	"context"
	"net/http"
	"strings"

	"villainrsty-ecommerce-server/internal/adapters/http/httpx"
	"villainrsty-ecommerce-server/internal/core/auth/models"
	"villainrsty-ecommerce-server/internal/core/auth/ports"
)

func AuthJWT(jwtService ports.JWTService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := extractToken(r)
			if token == "" {
				httpx.Error(w, http.StatusUnauthorized, "missing authorization token", "MISSING_AUTHORIZATION_TOKEN")
				return
			}

			user, err := jwtService.ValidateToken(token)
			if err != nil {
				httpx.ErrorWithDetails(w, http.StatusUnauthorized, "invalid token", "INVALID_TOKEN", err.Error())
				return
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
