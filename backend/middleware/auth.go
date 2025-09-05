package middleware

import (
	"net/http"
	"strings"

	"spyal/auth"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "authorization required", http.StatusUnauthorized)
			return
		}

		tokenLength := 2
		parts := strings.SplitN(authHeader, " ", tokenLength)
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "invalid authorization header", http.StatusUnauthorized)
			return
		}

		username, ok := auth.VerifyToken(parts[1])
		if !ok {
			http.Error(w, "invalid or expired token", http.StatusUnauthorized)
			return
		}

		ctx := r.Context()
		ctx = auth.WithUsername(ctx, username)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
