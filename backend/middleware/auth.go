package middleware

import (
	"context"
	"net/http"

	"spyal/auth"
)

//nolint
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rawToken := ""

		if c, err := r.Cookie("auth"); err == nil {
			rawToken = c.Value
		}
		if rawToken == "" {
			http.Error(w, "Duhet te logohesh", http.StatusUnauthorized)
			return
		}

		id, username, ok := auth.VerifyToken(rawToken)
		if !ok {
			http.Error(w, "Duhet te logohesh", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "username", username)
		ctx = context.WithValue(ctx, "id", int64(id))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
