package handler

import (
	"sync"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type rateLimiter struct {
	mu      sync.Mutex
	clients map[string]time.Time
	limit   time.Duration
}

func NewRateLimiter(limit time.Duration) *rateLimiter {
	return &rateLimiter{
		clients: make(map[string]time.Time),
		limit:   limit,
	}
}

func (r *rateLimiter) Allow(clientID string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	last, ok := r.clients[clientID]
	now := time.Now()
	if !ok || now.Sub(last) > r.limit {
		r.clients[clientID] = now
		return true
	}
	return false
}

// Пример middleware для gRPC
func (s *TaskServer) RateLimit(clientID string) error {
	if s.RateLimiter != nil && !s.RateLimiter.Allow(clientID) {
		return status.Error(codes.ResourceExhausted, "rate limit exceeded")
	}
	return nil
}
