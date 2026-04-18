package middleware

import (
	"strings"
	"sync"
	"time"

	"go-server/internal/config"
	"go-server/internal/errcode"

	"github.com/gin-gonic/gin"
	"go.uber.org/ratelimit"
)

type clientRateLimiter struct {
	mu       sync.Mutex
	limiters map[string]ratelimit.Limiter
	rps      int
}

func NewRateLimit(cfg config.RateLimitConfig) gin.HandlerFunc {
	if !cfg.Enabled || cfg.RequestsPerSecond <= 0 {
		return func(c *gin.Context) { c.Next() }
	}

	store := &clientRateLimiter{
		limiters: make(map[string]ratelimit.Limiter),
		rps:      cfg.RequestsPerSecond,
	}

	return func(c *gin.Context) {
		if !shouldLimit(c.Request.URL.Path, cfg.ProtectedPrefixes) {
			c.Next()
			return
		}

		start := time.Now()
		store.get(c.ClientIP() + ":" + c.FullPath()).Take()
		if cfg.MaxDelay >= 0 && time.Since(start) > cfg.MaxDelay {
			_ = c.Error(errcode.ErrTooManyRequests.AsError())
			c.Abort()
			return
		}

		c.Next()
	}
}

func (s *clientRateLimiter) get(key string) ratelimit.Limiter {
	s.mu.Lock()
	defer s.mu.Unlock()

	if limiter, ok := s.limiters[key]; ok {
		return limiter
	}

	limiter := ratelimit.New(s.rps, ratelimit.WithoutSlack)
	s.limiters[key] = limiter
	return limiter
}

func shouldLimit(path string, protectedPrefixes []string) bool {
	if len(protectedPrefixes) == 0 {
		return true
	}
	for _, prefix := range protectedPrefixes {
		if strings.HasPrefix(path, prefix) {
			return true
		}
	}
	return false
}
