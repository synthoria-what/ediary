package middleware

import (
	"log"
	"net/http"
	"time"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now().UTC()
		log.Printf("start logging: %s", r.Method)
		next.ServeHTTP(w, r)
		log.Printf("end logging: %s, %s", r.Method, time.Since(start))
	})
}
