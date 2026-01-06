package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type rateLimiter struct {
	limiters map[string]*rate.Limiter
	mu       sync.RWMutex
	rate     rate.Limit
	burst    int
}

var (
	loginLimiter *rateLimiter
	once         sync.Once
)

func initLoginRateLimiter(requestsPerMinute int, burst int) {
	once.Do(func() {
		loginLimiter = &rateLimiter{
			limiters: make(map[string]*rate.Limiter),
			rate:     rate.Limit(float64(requestsPerMinute) / 60.0),
			burst:    burst,
		}

		go func() {
			ticker := time.NewTicker(5 * time.Minute)
			defer ticker.Stop()
			for range ticker.C {
				loginLimiter.cleanup()
			}
		}()
	})
}

func (rl *rateLimiter) getLimiter(key string) *rate.Limiter {
	rl.mu.RLock()
	limiter, exists := rl.limiters[key]
	rl.mu.RUnlock()

	if !exists {
		rl.mu.Lock()
		limiter = rate.NewLimiter(rl.rate, rl.burst)
		rl.limiters[key] = limiter
		rl.mu.Unlock()
	}

	return limiter
}

func (rl *rateLimiter) cleanup() {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	for key := range rl.limiters {
		if rl.limiters[key].Tokens() >= float64(rl.burst) {
			delete(rl.limiters, key)
		}
	}
}

func getClientIP(c *gin.Context) string {
	ip := c.GetHeader("X-Forwarded-For")
	if ip != "" {
		return ip
	}
	ip = c.GetHeader("X-Real-IP")
	if ip != "" {
		return ip
	}
	return c.ClientIP()
}

func RateLimitMiddleware(requestsPerMinute int, burst int) gin.HandlerFunc {
	initLoginRateLimiter(requestsPerMinute, burst)

	return func(c *gin.Context) {
		ip := getClientIP(c)
		limiter := loginLimiter.getLimiter(ip)

		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Too many login attempts. Please try again later.",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
