package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"

	"github.com/gix-coder/gix-coder/gateway/internal/audit"
	"github.com/gix-coder/gix-coder/gateway/internal/auth"
	"github.com/gix-coder/gix-coder/gateway/internal/authz"
	"github.com/gix-coder/gix-coder/gateway/internal/metrics"
	"github.com/gix-coder/gix-coder/gateway/internal/persistence"
)

const (
	IdempotencyKeyHeader = "Idempotency-Key"
)

// tokenBucket represents a token bucket for rate limiting.
type tokenBucket struct {
	mu         sync.Mutex
	tokens     float64
	lastRefill time.Time
	rate       float64 // tokens per second
	burst      float64 // maximum tokens
}

// RateLimitConfig holds rate limiting configuration.
type RateLimitConfig struct {
	Enabled           bool
	RequestsPerMinute int
	Burst             int
	RedisHost         string
	RedisPort         int
	RedisPassword     string
	RedisDB           int
	FailOpen          bool // if true, allow requests when Redis is unavailable
}

// RateLimiterInterface defines the interface for rate limiters.
type RateLimiterInterface interface {
	Middleware() gin.HandlerFunc
	Close() error
}

// RedisRateLimiter implements distributed token bucket rate limiting using Redis.
type RedisRateLimiter struct {
	config RateLimitConfig
	client *redis.Client
	local  *RateLimiter // fallback in-memory limiter
}

// RateLimiter implements in-memory token bucket rate limiting (fallback).
type RateLimiter struct {
	config  RateLimitConfig
	buckets sync.Map // key -> *tokenBucket
}

// NewRateLimiter creates a new rate limiter (in-memory only, for testing/fallback).
func NewRateLimiter(config RateLimitConfig) *RateLimiter {
	return &RateLimiter{
		config: config,
	}
}

// NewRedisRateLimiter creates a new Redis-backed rate limiter.
func NewRedisRateLimiter(config RateLimitConfig) (*RedisRateLimiter, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.RedisHost, config.RedisPort),
		Password: config.RedisPassword,
		DB:       config.RedisDB,
	})

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &RedisRateLimiter{
		config: config,
		client: client,
		local:  NewRateLimiter(config),
	}, nil
}

func (r *RedisRateLimiter) Close() error {
	return r.client.Close()
}

func (r *RedisRateLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !r.config.Enabled {
			c.Next()
			return
		}

		key := r.getRateLimitKey(c)
		allowed, retryAfter := r.checkRateLimit(c.Request.Context(), key)

		if !allowed {
			metrics.RateLimitRejectedTotal.WithLabelValues(key).Inc()
			if retryAfter > 0 {
				c.Header("Retry-After", strconv.Itoa(retryAfter))
			}
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": gin.H{
					"code":    "RATE_LIMITED",
					"message": "rate limit exceeded",
				},
			})
			return
		}

		c.Next()
	}
}

func (r *RedisRateLimiter) getRateLimitKey(c *gin.Context) string {
	// Use tenant ID + IP as key
	tenantID := c.GetString("tenant_id")
	if tenantID == "" {
		tenantID = "anonymous"
	}
	ip := c.ClientIP()
	return fmt.Sprintf("ratelimit:%s:%s", tenantID, ip)
}

func (r *RedisRateLimiter) checkRateLimit(ctx context.Context, key string) (bool, int) {
	// Try Redis first
	allowed, retryAfter, err := r.checkRateLimitRedis(ctx, key)
	if err != nil {
		// Redis error - fallback to in-memory if fail-open
		if r.config.FailOpen {
			return r.local.checkRateLimit(key), 0
		}
		// Fail-closed: reject request
		return false, 60
	}
	return allowed, retryAfter
}

func (r *RedisRateLimiter) checkRateLimitRedis(ctx context.Context, key string) (bool, int, error) {
	now := time.Now().Unix()
	windowStart := now - 60 // 1-minute sliding window
	rate := r.config.RequestsPerMinute
	burst := r.config.Burst

	// Lua script for atomic token bucket with sliding window
	script := redis.NewScript(`
		local key = KEYS[1]
		local now = tonumber(ARGV[1])
		local window_start = tonumber(ARGV[2])
		local rate = tonumber(ARGV[3])
		local burst = tonumber(ARGV[4])

		-- Remove expired entries
		redis.call('ZREMRANGEBYSCORE', key, '-inf', window_start)

		-- Count current requests in window
		local current = redis.call('ZCARD', key)

		if current >= burst then
			-- At burst limit, check if we can allow based on rate
			local oldest = redis.call('ZRANGE', key, 0, 0, 'WITHSCORES')
			if #oldest > 0 then
				local oldest_time = tonumber(oldest[2])
				local elapsed = now - oldest_time
				local expected = math.floor(elapsed * rate / 60)
				if current - expected >= burst then
					return {0, 60}
				end
			end
			return {0, 60}
		end

		-- Add current request
		redis.call('ZADD', key, now, now .. '-' .. math.random(1000000))
		redis.call('EXPIRE', key, 120)
		return {1, 0}
	`)

	result, err := script.Run(ctx, r.client, []string{key}, now, windowStart, rate, burst).Slice()
	if err != nil {
		return false, 0, err
	}

	// Redis returns int64 for integer values
	allowedVal, ok := result[0].(int64)
	if !ok {
		return false, 0, fmt.Errorf("unexpected type for allowed: %T", result[0])
	}
	allowed := allowedVal == 1

	retryAfterVal, ok := result[1].(int64)
	if !ok {
		return false, 0, fmt.Errorf("unexpected type for retryAfter: %T", result[1])
	}
	retryAfter := int(retryAfterVal)

	return allowed, retryAfter, nil
}

// In-memory rate limiter methods (for fallback/testing)
func (r *RateLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !r.config.Enabled {
			c.Next()
			return
		}

		key := r.getRateLimitKey(c)
		allowed := r.checkRateLimit(key)

		if !allowed {
			metrics.RateLimitRejectedTotal.WithLabelValues(key).Inc()
			c.Header("Retry-After", "60")
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": gin.H{
					"code":    "RATE_LIMITED",
					"message": "rate limit exceeded",
				},
			})
			return
		}

		c.Next()
	}
}

func (r *RateLimiter) getRateLimitKey(c *gin.Context) string {
	tenantID := c.GetString("tenant_id")
	if tenantID == "" {
		tenantID = "anonymous"
	}
	ip := c.ClientIP()
	return fmt.Sprintf("%s:%s", tenantID, ip)
}

func (r *RateLimiter) checkRateLimit(key string) bool {
	now := time.Now()
	rate := float64(r.config.RequestsPerMinute) / 60.0
	burst := float64(r.config.Burst)

	bucketInterface, loaded := r.buckets.LoadOrStore(key, &tokenBucket{
		tokens:     burst,
		lastRefill: now,
		rate:       rate,
		burst:      burst,
	})
	_ = loaded
	bucket := bucketInterface.(*tokenBucket)

	bucket.mu.Lock()
	defer bucket.mu.Unlock()

	elapsed := now.Sub(bucket.lastRefill).Seconds()
	bucket.tokens += elapsed * bucket.rate
	if bucket.tokens > bucket.burst {
		bucket.tokens = bucket.burst
	}
	bucket.lastRefill = now

	if bucket.tokens >= 1.0 {
		bucket.tokens -= 1.0
		return true
	}

	return false
}

// Close is a no-op for in-memory rate limiter (implements RateLimiterInterface).
func (r *RateLimiter) Close() error {
	return nil
}

// IdempotencyMiddleware handles idempotency keys.
type IdempotencyMiddleware struct {
	repo persistence.IdempotencyRepository
	ttl  time.Duration
}

func NewIdempotencyMiddleware(repo persistence.IdempotencyRepository, ttl time.Duration) *IdempotencyMiddleware {
	return &IdempotencyMiddleware{
		repo: repo,
		ttl:  ttl,
	}
}

func (m *IdempotencyMiddleware) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		idempotencyKey := c.GetHeader(IdempotencyKeyHeader)
		if idempotencyKey == "" {
			c.Next()
			return
		}

		// Check if we have a cached response
		response, found, err := m.repo.Get(c.Request.Context(), idempotencyKey)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": gin.H{
					"code":    "IDEMPOTENCY_ERROR",
					"message": err.Error(),
				},
			})
			return
		}

		if found {
			// Return cached response
			c.Header("X-Idempotency-Replay", "true")
			c.Data(response.StatusCode, response.ContentType, response.Body)
			c.Abort()
			return
		}

		// Try to acquire advisory lock for this idempotency key
		// This ensures only one request processes, others wait or get conflict
		locked, err := m.repo.TryLock(c.Request.Context(), idempotencyKey)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": gin.H{
					"code":    "IDEMPOTENCY_ERROR",
					"message": "failed to acquire lock: " + err.Error(),
				},
			})
			return
		}

		if !locked {
			// Another request is processing this key - wait and retry once
			// Small delay to let the first request complete
			time.Sleep(50 * time.Millisecond)

			// Check again after waiting
			response, found, err := m.repo.Get(c.Request.Context(), idempotencyKey)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"error": gin.H{
						"code":    "IDEMPOTENCY_ERROR",
						"message": err.Error(),
					},
				})
				return
			}

			if found {
				c.Header("X-Idempotency-Replay", "true")
				c.Data(response.StatusCode, response.ContentType, response.Body)
				c.Abort()
				return
			}

			// Still not found - this is a conflict (different payload for same key)
			metrics.IdempotencyConflictsTotal.WithLabelValues().Inc()
			c.AbortWithStatusJSON(http.StatusConflict, gin.H{
				"error": gin.H{
					"code":    "IDEMPOTENCY_CONFLICT",
					"message": "idempotency key in use with different payload",
				},
			})
			return
		}

		// We got the lock - proceed with request
		c.Set("idempotency_key", idempotencyKey)
		c.Next()

		// Lock is automatically released at transaction end (advisory xact lock)
	}
}

// AuditMiddleware handles audit logging.
type AuditMiddleware struct {
	client *audit.Client
}

func NewAuditMiddleware(client *audit.Client) *AuditMiddleware {
	return &AuditMiddleware{client: client}
}

func (m *AuditMiddleware) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Process request
		c.Next()

		// Log audit event for security-sensitive operations
		if c.Writer.Status() >= 400 {
			if err := m.client.Log(c.Request.Context(), audit.LogRequest{
				Event: audit.Event{
					EventType: "gateway.request.failed",
					Actor: audit.Actor{
						Type:     "user",
						ID:       c.GetString("user_id"),
						TenantID: c.GetString("tenant_id"),
					},
					Resource: audit.Resource{
						Type: "gateway.request",
						ID:   c.FullPath(),
					},
					Action:  c.Request.Method,
					Outcome: "failure",
					Details: map[string]string{
						"status_code": fmt.Sprintf("%d", c.Writer.Status()),
						"path":        c.Request.URL.Path,
						"method":      c.Request.Method,
					},
				},
			}); err != nil {
				// Log the error but don't fail the request
				// In production, this would be logged to a monitoring system
				_ = err
			}
		}
	}
}

// AuthMiddleware validates JWT tokens.
func Auth(authClient *auth.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := authClient.ExtractBearerToken(c.Request)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": gin.H{
					"code":    "AUTH_FAILED",
					"message": err.Error(),
				},
			})
			return
		}

		claims, err := authClient.ValidateToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": gin.H{
					"code":    "AUTH_FAILED",
					"message": err.Error(),
				},
			})
			return
		}

		// Store claims in context
		c.Set("claims", claims)
		c.Set("user_id", claims.Subject)
		c.Set("tenant_id", claims.TenantID)
		c.Set("scopes", claims.Scopes)

		c.Next()
	}
}

// AuthzMiddleware handles authorization.
type AuthzMiddleware struct {
	client *authz.Client
}

func NewAuthzMiddleware(client *authz.Client) *AuthzMiddleware {
	return &AuthzMiddleware{client: client}
}

func (m *AuthzMiddleware) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, ok := auth.GetClaims(c)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": gin.H{
					"code":    "AUTH_FAILED",
					"message": "no claims in context",
				},
			})
			return
		}

		// Extract resource and action from request
		resource := c.FullPath()
		action := c.Request.Method

		userID, ok := auth.GetUserID(c)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": gin.H{
					"code":    "AUTH_FAILED",
					"message": "user ID not found",
				},
			})
			return
		}

		tenantID, ok := auth.GetTenantID(c)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": gin.H{
					"code":    "AUTH_FAILED",
					"message": "tenant ID not found",
				},
			})
			return
		}

		// Evaluate policy
		resp, err := m.client.Evaluate(c.Request.Context(), &authz.EvaluateRequest{
			Subject:  userID,
			TenantID: tenantID,
			Resource: resource,
			Action:   action,
			Context:  map[string]string{},
		})
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": gin.H{
					"code":    "AUTHZ_ERROR",
					"message": err.Error(),
				},
			})
			return
		}

		if !resp.Allowed {
			metrics.AuthzDenialsTotal.WithLabelValues().Inc()
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": gin.H{
					"code":    "AUTHZ_DENIED",
					"message": "access denied",
					"details": resp.Reason,
				},
			})
			return
		}

		// Store grants in context
		c.Set("grants", resp.Grants)
		c.Next()
	}
}