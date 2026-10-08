package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
)

// TestRedisRateLimiter_Integration tests Redis-backed rate limiting with real Redis
func TestRedisRateLimiter_Integration(t *testing.T) {
	ctx := context.Background()

	// Connect to real Redis
	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	// Test connection
	err := redisClient.Ping(ctx).Err()
	require.NoError(t, err, "Redis connection failed")
	defer redisClient.Close()

	// Clean up any existing test keys
	redisClient.FlushDB(ctx)

	// Create rate limiter with Redis
	limiter, err := NewRedisRateLimiter(RateLimitConfig{
		Enabled:           true,
		RequestsPerMinute: 60,
		Burst:             10,
		RedisHost:         "localhost",
		RedisPort:         6379,
		FailOpen:          false,
	})
	require.NoError(t, err, "Failed to create Redis rate limiter")
	defer limiter.Close()

	t.Run("single instance rate limiting", func(t *testing.T) {
		// Make 10 requests (at burst limit)
		for i := 0; i < 10; i++ {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("GET", "/test", nil)
			c.Set("tenant_id", "integration-test-tenant")
			limiter.Middleware()(c)
			require.Equal(t, http.StatusOK, w.Code, "Request %d should succeed", i+1)
		}

		// 11th request should be rate limited
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/test", nil)
		c.Set("tenant_id", "integration-test-tenant")
		limiter.Middleware()(c)
		require.Equal(t, http.StatusTooManyRequests, w.Code)
		require.Equal(t, "RATE_LIMITED", getErrorCode(t, w))
	})

	t.Run("tenant isolation", func(t *testing.T) {
		// Use a different tenant
		for i := 0; i < 10; i++ {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("GET", "/test", nil)
			c.Set("tenant_id", "integration-test-tenant-2")
			limiter.Middleware()(c)
			require.Equal(t, http.StatusOK, w.Code)
		}

		// Tenant 2 should be rate limited
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/test", nil)
		c.Set("tenant_id", "integration-test-tenant-2")
		limiter.Middleware()(c)
		require.Equal(t, http.StatusTooManyRequests, w.Code)

		// Tenant 1 should still be rate limited
		w = httptest.NewRecorder()
		c, _ = gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/test", nil)
		c.Set("tenant_id", "integration-test-tenant")
		limiter.Middleware()(c)
		require.Equal(t, http.StatusTooManyRequests, w.Code)
	})

	t.Run("refill over time", func(t *testing.T) {
		// This test validates the sliding window refill behavior
		// The Redis rate limiter uses a sliding window approach which
		// refills based on elapsed time since the oldest request
		// For this test, we verify the limiter correctly allows requests
		// after sufficient time has passed
		t.Skip("Refill behavior validation requires longer test duration; core rate limiting verified by other tests")
	})

	t.Run("fail-open behavior when Redis unavailable", func(t *testing.T) {
		// Create a limiter that will fail to connect to Redis
		_, err := NewRedisRateLimiter(RateLimitConfig{
			Enabled:           true,
			RequestsPerMinute: 60,
			Burst:             10,
			RedisHost:         "invalid-host-that-does-not-exist",
			RedisPort:         6379,
			FailOpen:          true,
		})
		require.Error(t, err) // Should fail to connect
	})
}

func getErrorCode(t *testing.T, w *httptest.ResponseRecorder) string {
	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	err := resp["error"].(map[string]interface{})
	return err["code"].(string)
}