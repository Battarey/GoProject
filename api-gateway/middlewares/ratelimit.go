package middlewares

import (
	"net/http"
	"strings"
	"sync"
	"time"
)

var rateLimit = 10 // requests
var rateWindow = time.Minute
var clients = make(map[string][]time.Time)
var mu sync.Mutex

// RateLimitMiddleware ограничивает частоту запросов с одного IP
func RateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr
		if i := strings.LastIndex(ip, ":"); i != -1 {
			ip = ip[:i]
		}
		mu.Lock()
		times := clients[ip]
		now := time.Now()
		var filtered []time.Time
		for _, t := range times {
			if now.Sub(t) < rateWindow {
				filtered = append(filtered, t)
			}
		}
		if len(filtered) >= rateLimit {
			mu.Unlock()
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}
		filtered = append(filtered, now)
		clients[ip] = filtered
		mu.Unlock()
		next.ServeHTTP(w, r)
	})
}
