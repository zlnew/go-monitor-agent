package middleware

import (
	"context"
	"net/http"
	"strings"

	"horizonx-server/internal/domain"
)

type key int

const ServerIDKey key = 0

func AgentAuth(serverRepo domain.ServerRepository) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Missing Authorization Header", http.StatusUnauthorized)
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				http.Error(w, "Invalid Token Format", http.StatusUnauthorized)
				return
			}
			token := parts[1]

			server, err := serverRepo.GetByToken(r.Context(), token)
			if err != nil {
				http.Error(w, "Invalid or Revoked Token", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), ServerIDKey, server.ID)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
