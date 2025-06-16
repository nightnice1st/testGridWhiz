package ratelimit

import (
	"sync"
	"time"
)

type RateLimiter struct {
	attempts    map[string]*attemptInfo
	mu          sync.RWMutex
	maxAttempts int
	window      time.Duration
}

type attemptInfo struct {
	count     int
	firstTime time.Time
}

func NewRateLimiter(maxAttempts int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		attempts:    make(map[string]*attemptInfo),
		maxAttempts: maxAttempts,
		window:      window,
	}

	go rl.cleanup()
	return rl
}

func (rl *RateLimiter) Allow(key string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	info, exists := rl.attempts[key]

	if !exists || now.Sub(info.firstTime) > rl.window {
		rl.attempts[key] = &attemptInfo{
			count:     1,
			firstTime: now,
		}
		return true
	}

	if info.count >= rl.maxAttempts {
		return false
	}

	info.count++
	return true
}

func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(rl.window)
	defer ticker.Stop()

	for range ticker.C {
		rl.mu.Lock()
		now := time.Now()
		for key, info := range rl.attempts {
			if now.Sub(info.firstTime) > rl.window {
				delete(rl.attempts, key)
			}
		}
		rl.mu.Unlock()
	}
}
