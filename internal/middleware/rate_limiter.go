package middleware

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// rateLimiter struct holds a map of rate limiters per IP
type ipRateLimiter struct {
	ips map[string]*rate.Limiter
	mu  *sync.RWMutex
	r   rate.Limit
	b   int
}

// newIPRateLimiter creates a new rate limiter that limits per IP
func newIPRateLimiter(r rate.Limit, b int) *ipRateLimiter {
	i := &ipRateLimiter{
		ips: make(map[string]*rate.Limiter),
		mu:  &sync.RWMutex{},
		r:   r,
		b:   b,
	}

	// Optional: a cleanup routine to remove old IPs could be added here for production
	return i
}

// getLimiter returns the rate limiter for the provided IP
func (i *ipRateLimiter) getLimiter(ip string) *rate.Limiter {
	i.mu.Lock()
	defer i.mu.Unlock()

	limiter, exists := i.ips[ip]
	if !exists {
		limiter = rate.NewLimiter(i.r, i.b)
		i.ips[ip] = limiter
	}

	return limiter
}

// RateLimitMiddleware applies an IP-based rate limit of r requests per second with burst b.
func RateLimitMiddleware() gin.HandlerFunc {
	// Allow 5 requests per second with a burst of 10
	limiter := newIPRateLimiter(5, 10)

	return func(c *gin.Context) {
		ip := c.ClientIP()
		if !limiter.getLimiter(ip).Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"success": false,
				"message": "Rate limit exceeded",
			})
			return
		}
		c.Next()
	}
}
