package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

const (
	defaultMaxEntries = 10_000
	defaultIdleTTL    = 30 * time.Minute
)

type ipLimiterEntry struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

// IPRateLimiter is a per-IP token-bucket limiter.
type IPRateLimiter struct {
	mu         sync.Mutex
	limiters   map[string]*ipLimiterEntry
	limit      rate.Limit
	burst      int
	maxEntries int
	idleTTL    time.Duration
	now        func() time.Time
}

func NewIPRateLimiter(perMinute float64, burst int) *IPRateLimiter {
	if perMinute <= 0 {
		perMinute = 60
	}
	if burst <= 0 {
		burst = int(perMinute)
		if burst < 1 {
			burst = 1
		}
	}
	return &IPRateLimiter{
		limiters:   make(map[string]*ipLimiterEntry),
		limit:      rate.Limit(perMinute / 60.0), // tokens per second
		burst:      burst,
		maxEntries: defaultMaxEntries,
		idleTTL:    defaultIdleTTL,
		now:        time.Now,
	}
}

func (l *IPRateLimiter) get(ip string) *rate.Limiter {
	l.mu.Lock()
	defer l.mu.Unlock()
	now := l.now()
	if entry, ok := l.limiters[ip]; ok {
		entry.lastSeen = now
		return entry.limiter
	}
	// Cap growth: when map is full, drop idle entries before inserting.
	if len(l.limiters) >= l.maxEntries {
		l.cleanupIdleLocked(now)
	}
	lim := rate.NewLimiter(l.limit, l.burst)
	l.limiters[ip] = &ipLimiterEntry{limiter: lim, lastSeen: now}
	return lim
}

// Middleware rejects requests that exceed the per-IP rate.
func (l *IPRateLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		if !l.get(ip).Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"detail": "Rate limit exceeded. Please try again later.",
			})
			return
		}
		c.Next()
	}
}

// cleanupIdleLocked drops entries idle longer than idleTTL. Caller must hold l.mu.
func (l *IPRateLimiter) cleanupIdleLocked(now time.Time) int {
	if l.idleTTL <= 0 {
		return 0
	}
	cutoff := now.Add(-l.idleTTL)
	removed := 0
	for ip, entry := range l.limiters {
		if entry.lastSeen.Before(cutoff) {
			delete(l.limiters, ip)
			removed++
		}
	}
	return removed
}

// Len reports how many IP entries are tracked (for tests/metrics).
func (l *IPRateLimiter) Len() int {
	l.mu.Lock()
	defer l.mu.Unlock()
	return len(l.limiters)
}
