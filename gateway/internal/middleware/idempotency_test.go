package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gix-coder/gix-coder/gateway/internal/persistence"
	"github.com/stretchr/testify/require"
)

// MockIdempotencyRepository is a test implementation
type MockIdempotencyRepository struct {
	store map[string]*persistence.IdempotencyResponse
}

func NewMockIdempotencyRepository() *MockIdempotencyRepository {
	return &MockIdempotencyRepository{
		store: make(map[string]*persistence.IdempotencyResponse),
	}
}

func (m *MockIdempotencyRepository) Get(ctx context.Context, key string) (*persistence.IdempotencyResponse, bool, error) {
	resp, ok := m.store[key]
	return resp, ok, nil
}

func (m *MockIdempotencyRepository) Set(ctx context.Context, key, requestHash string, resp *persistence.IdempotencyResponse, ttl time.Duration) error {
	m.store[key] = resp
	return nil
}

func (m *MockIdempotencyRepository) TryLock(ctx context.Context, key string) (bool, error) {
	return true, nil
}

func (m *MockIdempotencyRepository) Unlock(ctx context.Context, key string) error {
	return nil
}

func (m *MockIdempotencyRepository) Cleanup(ctx context.Context) error {
	return nil
}

func TestIdempotencyMiddleware_FirstRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := NewMockIdempotencyRepository()
	middleware := NewIdempotencyMiddleware(repo, 24*time.Hour)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/test", nil)
	c.Request.Header.Set("Idempotency-Key", "test-key-1")

	called := false
	handler := func(c *gin.Context) {
		called = true
		c.JSON(http.StatusOK, gin.H{"result": "success"})
	}

	middleware.Middleware()(c)
	handler(c)

	require.True(t, called)
	require.Equal(t, http.StatusOK, w.Code)
	require.Equal(t, `{"result":"success"}`, w.Body.String())
}

func TestIdempotencyMiddleware_ReplayRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := NewMockIdempotencyRepository()
	// Pre-populate with a cached response
	repo.store["test-key-1"] = &persistence.IdempotencyResponse{
		StatusCode:  http.StatusOK,
		ContentType: "application/json",
		Body:        []byte(`{"cached":true}`),
	}

	middleware := NewIdempotencyMiddleware(repo, 24*time.Hour)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/test", nil)
	c.Request.Header.Set("Idempotency-Key", "test-key-1")

	// Middleware should write the cached response and abort
	middleware.Middleware()(c)

	// Verify the middleware handled the request and returned the cached response
	require.Equal(t, http.StatusOK, w.Code)
	require.Equal(t, `{"cached":true}`, w.Body.String())
	require.Equal(t, "true", w.Header().Get("X-Idempotency-Replay"))
}

func TestIdempotencyMiddleware_NoKey(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := NewMockIdempotencyRepository()
	middleware := NewIdempotencyMiddleware(repo, 24*time.Hour)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/test", nil)
	// No Idempotency-Key header

	called := false
	handler := func(c *gin.Context) {
		called = true
		c.JSON(http.StatusOK, gin.H{"result": "success"})
	}

	middleware.Middleware()(c)
	handler(c)

	require.True(t, called)
	require.Equal(t, http.StatusOK, w.Code)
}

func TestIdempotencyMiddleware_ConcurrentDuplicate(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := NewMockIdempotencyRepository()
	middleware := NewIdempotencyMiddleware(repo, 24*time.Hour)

	// First request - processes and caches response
	w1 := httptest.NewRecorder()
	c1, _ := gin.CreateTestContext(w1)
	c1.Request = httptest.NewRequest("POST", "/test", nil)
	c1.Request.Header.Set("Idempotency-Key", "concurrent-key")

	// Middleware allows through (no cached response)
	middleware.Middleware()(c1)
	// Simulate handler saving response to repo
	repo.store["concurrent-key"] = &persistence.IdempotencyResponse{
		StatusCode:  http.StatusOK,
		ContentType: "application/json",
		Body:        []byte(`{"result":"success"}`),
	}
	// Simulate handler response
	c1.JSON(http.StatusOK, gin.H{"result": "success"})

	require.Equal(t, http.StatusOK, w1.Code)

	// Second request with same key - should replay from cache
	w2 := httptest.NewRecorder()
	c2, _ := gin.CreateTestContext(w2)
	c2.Request = httptest.NewRequest("POST", "/test", nil)
	c2.Request.Header.Set("Idempotency-Key", "concurrent-key")

	middleware.Middleware()(c2)

	require.Equal(t, http.StatusOK, w2.Code)
	require.Equal(t, `{"result":"success"}`, w2.Body.String())
	require.Equal(t, "true", w2.Header().Get("X-Idempotency-Replay"))
}