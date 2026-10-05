package middleware

import (
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"sync"
	"synthori/ediary/m/handlers"
	"synthori/ediary/m/models"
	"time"
)

func RateLimitMiddleware(next http.Handler) http.Handler {
	usersRequests := make(map[string]*models.UserStats)
	var mu sync.Mutex
	go func() {
		ticker := time.NewTicker(60 * time.Second)
		for range ticker.C {
			mu.Lock()
			usersRequests = make(map[string]*models.UserStats)
			slog.Info("Global rate limit map flushed successfully")
			mu.Unlock()
		}
	}()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			handlers.WriteError(w, 401, errors.New("ip error"))
			return
		}
		mu.Lock()

		if _, exists := usersRequests[ip]; !exists {
			usersRequests[ip] = &models.UserStats{}
		}

		stats := usersRequests[ip]
		slog.Info("RateLimit middleware", "Request count", fmt.Sprintf("%v", stats.RequestCount))

		if stats.RequestCount >= 10 {
			handlers.WriteError(w, 429, errors.New("too many requests"))
			if stats.IsBanned {
				mu.Unlock()
				return
			}

			stats.IsBanned = true
			mu.Unlock()
			go func(ip string) {
				fmt.Println("start timer")
				timer := time.NewTimer(30 * time.Second)
				<-timer.C
				fmt.Println("end timer")
				mu.Lock()
				stats.RequestCount = 0
				mu.Unlock()
			}(ip)
			return
		}

		stats.RequestCount++
		mu.Unlock()

		next.ServeHTTP(w, r)
	})
}
