package middleware

import (
	"net/http"

	"villainrsty-ecommerce-server/internal/adapters/http/lib/httpx"

	"github.com/casbin/casbin/v3"
)

func CasbinRBAC(enforcer *casbin.Enforcer) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user := GetUserFromContext(*r)
			if user == nil {
				httpx.Error(w, http.StatusUnauthorized, "unauthorized", "UNAUTHORIZED")
				return
			}

			role := user.Role
			if role == "" {
				role = "customer"
			}

			ok, err := enforcer.Enforce(role, r.URL.Path, r.Method)
			if err != nil {
				httpx.ErrorWithDetails(w, http.StatusInternalServerError, "unauthorized error", "AUTHZ_ERROR", err.Error())
				return
			}

			if !ok {
				httpx.Error(w, http.StatusForbidden, "access denied", "FORBIDDEN")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
