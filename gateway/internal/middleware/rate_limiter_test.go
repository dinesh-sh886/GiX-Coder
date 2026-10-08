package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestRateLimiter_WithinLimit(t *testing.T) {
	gin.SetMode(gin.TestMode)

	config := RateLimitConfig{
		Enabled:           true,
		RequestsPerMinute: 10,
		Burst:             5,
	}
	limiter := NewRateLimiter(config)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/test", nil)
	c.Set("tenant_id", "tenant-1")

	called := false
	handler := func(c *gin.Context) {
		called = true
	}

	// Make 5 requests (within burst limit)
	for i := 0; i < 5; i++ {
		w = httptest.NewRecorder()
		c, _ = gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/test", nil)
		c.Set("tenant_id", "tenant-1")
		limiter.Middleware()(c)
		handler(c)
		require.Equal(t, http.StatusOK, w.Code)
	}
	require.True(t, called)
}

func TestRateLimiter_ExceedsLimit(t *testing.T) {
	gin.SetMode(gin.TestMode)

	config := RateLimitConfig{
		Enabled:           true,
		RequestsPerMinute: 10,
		Burst:             3,
	}
	limiter := NewRateLimiter(config)

	// Make 3 requests (at burst limit)
	for i := 0; i < 3; i++ {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/test", nil)
		c.Set("tenant_id", "tenant-1")
		limiter.Middleware()(c)
		require.Equal(t, http.StatusOK, w.Code)
	}

	// 4th request should be rate limited
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/test", nil)
	c.Set("tenant_id", "tenant-1")
	limiter.Middleware()(c)
	require.Equal(t, http.StatusTooManyRequests, w.Code)
	require.Contains(t, w.Body.String(), "RATE_LIMITED")
	require.Equal(t, "60", w.Header().Get("Retry-After"))
}

func TestRateLimiter_Disabled(t *testing.T) {
	gin.SetMode(gin.TestMode)

	config := RateLimitConfig{
		Enabled:           false,
		RequestsPerMinute: 1,
		Burst:             1,
	}
	limiter := NewRateLimiter(config)

	// Make many requests - should all pass since disabled
	for i := 0; i < 100; i++ {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/test", nil)
		c.Set("tenant_id", "tenant-1")
		limiter.Middleware()(c)
		require.Equal(t, http.StatusOK, w.Code)
	}
}

func TestRateLimiter_TenantIsolation(t *testing.T) {
	gin.SetMode(gin.TestMode)

	config := RateLimitConfig{
		Enabled:           true,
		RequestsPerMinute: 10,
		Burst:             2,
	}
	limiter := NewRateLimiter(config)

	// Tenant 1 makes 2 requests (at burst limit)
	for i := 0; i < 2; i++ {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/test", nil)
		c.Set("tenant_id", "tenant-1")
		limiter.Middleware()(c)
		require.Equal(t, http.StatusOK, w.Code)
	}

	// Tenant 1 3rd request should be rate limited
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/test", nil)
	c.Set("tenant_id", "tenant-1")
	limiter.Middleware()(c)
	require.Equal(t, http.StatusTooManyRequests, w.Code)

	// Tenant 2 should still be able to make requests
	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/test", nil)
	c.Set("tenant_id", "tenant-2")
	limiter.Middleware()(c)
	require.Equal(t, http.StatusOK, w.Code)
}

func TestRateLimiter_AnonymousTenant(t *testing.T) {
	gin.SetMode(gin.TestMode)

	config := RateLimitConfig{
		Enabled:           true,
		RequestsPerMinute: 10,
		Burst:             2,
	}
	limiter := NewRateLimiter(config)

	// Requests without tenant_id should use "anonymous"
	for i := 0; i < 2; i++ {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/test", nil)
		// No tenant_id set
		limiter.Middleware()(c)
		require.Equal(t, http.StatusOK, w.Code)
	}

	// 3rd request should be rate limited
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/test", nil)
	limiter.Middleware()(c)
	require.Equal(t, http.StatusTooManyRequests, w.Code)
}

func TestRateLimiter_RefillOverTime(t *testing.T) {
	gin.SetMode(gin.TestMode)

	config := RateLimitConfig{
		Enabled:           true,
		RequestsPerMinute: 60, // 1 per second
		Burst:             1,
	}
	limiter := NewRateLimiter(config)

	// First request allowed
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/test", nil)
	c.Set("tenant_id", "tenant-1")
	limiter.Middleware()(c)
	require.Equal(t, http.StatusOK, w.Code)

	// Immediate second request should be rate limited (burst=1)
	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/test", nil)
	c.Set("tenant_id", "tenant-1")
	limiter.Middleware()(c)
	require.Equal(t, http.StatusTooManyRequests, w.Code)

	// Wait for refill (1 second)
	time.Sleep(1100 * time.Millisecond)

	// Should be allowed again
	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/test", nil)
	c.Set("tenant_id", "tenant-1")
	limiter.Middleware()(c)
	require.Equal(t, http.StatusOK, w.Code)
}