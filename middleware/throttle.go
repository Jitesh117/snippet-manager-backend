package middleware

import (
	"net/http"

	"golang.org/x/time/rate"
)

var rateLimiter = rate.NewLimiter(1, 5)

func RateLimiter(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !rateLimiter.Allow() {
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}
