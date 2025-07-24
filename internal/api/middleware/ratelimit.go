package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type RateLimiter struct {
	limiters map[string]*rate.Limiter
	mu       sync.RWMutex
	rate     rate.Limit
	burst    int
}

func NewRateLimiter(r rate.Limit, burst int) *RateLimiter {
	return &RateLimiter{
		limiters: make(map[string]*rate.Limiter),
		rate:     r,
		burst:    burst,
	}
}

func (rl *RateLimiter) getLimiter(ip string) *rate.Limiter {
	rl.mu.RLock()
	limiter, exists := rl.limiters[ip]
	rl.mu.RUnlock()

	if !exists {
		rl.mu.Lock()
		limiter = rate.NewLimiter(rl.rate, rl.burst)
		rl.limiters[ip] = limiter
		rl.mu.Unlock()
	}

	return limiter
}

func (rl *RateLimiter) RateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		limiter := rl.getLimiter(ip)

		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":       "Rate limit exceeded",
				"retry_after": "1s",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

func NewGlobalSemaphore(maxConcurrent int) *GlobalSemaphore {
	return &GlobalSemaphore{
		sem: make(chan struct{}, maxConcurrent),
	}
}

type GlobalSemaphore struct {
	sem chan struct{}
}

func (gs *GlobalSemaphore) Acquire() bool {
	select {
	case gs.sem <- struct{}{}:
		return true
	default:
		return false
	}
}

func (gs *GlobalSemaphore) Release() {
	<-gs.sem
}

func (gs *GlobalSemaphore) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !gs.Acquire() {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"error":       "Server overloaded",
				"retry_after": "5s",
			})
			c.Abort()
			return
		}
		defer gs.Release()
		c.Next()
	}
}

// Cleanup old rate limiters periodically
func (rl *RateLimiter) Cleanup() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		rl.mu.Lock()
		// Simple cleanup: remove old limiters (this is a basic implementation)
		// In a production environment, you might want to track last usage time
		for ip := range rl.limiters {
			// Remove limiters randomly to prevent memory leak
			// This is a simple approach - in production you'd track usage time
			if len(rl.limiters) > 1000 { // arbitrary threshold
				delete(rl.limiters, ip)
				break
			}
		}
		rl.mu.Unlock()
	}
}
