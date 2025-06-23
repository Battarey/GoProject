package handler

import (
	"sync"
	"time"
)

type RateLimiter struct {
	mu      sync.Mutex
	buckets map[string][]int64 // ключ: email или IP, значения: unix timestamps
	limit   int
	window  int64 // в секундах
}

func NewRateLimiter(limit int, windowSec int64) *RateLimiter {
	return &RateLimiter{
		buckets: make(map[string][]int64),
		limit:   limit,
		window:  windowSec,
	}
}

func (r *RateLimiter) Allow(key string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	now := time.Now().Unix()
	windowStart := now - r.window
	times := r.buckets[key]
	// Оставляем только те, что в окне
	var filtered []int64
	for _, t := range times {
		if t > windowStart {
			filtered = append(filtered, t)
		}
	}
	if len(filtered) >= r.limit {
		return false
	}
	filtered = append(filtered, now)
	r.buckets[key] = filtered
	return true
}

var RegLimiter = NewRateLimiter(5, 60)    // 5 регистраций в минуту на email
var LoginLimiter = NewRateLimiter(10, 60) // 10 логинов в минуту на email
