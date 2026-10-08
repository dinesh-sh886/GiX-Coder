package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestValidationMiddleware_ValidRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	middleware := ValidationMiddleware()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/test", nil)
	c.Request.Header.Set("Content-Type", "application/json")

	called := false
	handler := func(c *gin.Context) {
		called = true
		c.JSON(http.StatusOK, gin.H{"result": "success"})
	}

	middleware(c)
	handler(c)

	require.True(t, called)
	require.Equal(t, http.StatusOK, w.Code)
}

func TestValidateProtoMessage(t *testing.T) {
	// Test the ValidateProtoMessage function with a nil message
	err := ValidateProtoMessage(nil)
	require.NoError(t, err)
}