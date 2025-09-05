package middleware

import (
	"net/http"

	"spyal/auth"
)

func GuestMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader != "" {
			if username, ok := auth.VerifyToken(authHeader[len("Bearer "):]); ok && username != "" {
				http.Error(w, "guests only endpoint", http.StatusForbidden)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}
