package middleware

import (
	"golang.org/x/time/rate"
	"net/http"
)

const (
	reqPerSec = 5
	burst = 10
)

func RateLimitMiddleware(next http.Handler) http.Handler {
	limiter := rate.NewLimiter(reqPerSec, burst)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !limiter.Allow() {
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}
