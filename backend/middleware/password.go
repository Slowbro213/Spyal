package middleware

import (
	"net/http"

	"spyal/pkg/utils/ip"

	"go.uber.org/zap"
)

func UsernamePassword(username, password string, handler http.Handler, myLogger zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u, p, ok := r.BasicAuth()
		if !ok || u != username || p != password {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			http.Error(w, "Unauthorized.", http.StatusUnauthorized)
			clientIP := ip.ReadUserIP(r)
			myLogger.Warn("Unauthorized attemp by IP: " + clientIP)
			return
		}
		handler.ServeHTTP(w, r)
	}
}
