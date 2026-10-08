package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/gix-coder/gix-coder/gateway/internal/metrics"
	"github.com/gix-coder/gix-coder/shared/logging"
	"github.com/gix-coder/gix-coder/shared/tracing"
	"go.opentelemetry.io/otel/propagation"
)

const (
	RequestIDHeader     = "X-Request-ID"
	CorrelationIDHeader = "X-Correlation-ID"
	TraceParentHeader   = "traceparent"
)

// RequestID generates a unique request ID for each request.
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader(RequestIDHeader)
		if requestID == "" {
			requestID = uuid.New().String()
		}
		c.Set("request_id", requestID)
		c.Header(RequestIDHeader, requestID)
		c.Next()
	}
}

// Logging middleware for structured request/response logging.
func Logging(logger *logging.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()

		log := logger.Zerolog().Info().
			Str("method", method).
			Str("path", path).
			Int("status", status).
			Dur("latency", latency).
			Str("request_id", c.GetString("request_id"))

		if status >= 400 {
			log = logger.Zerolog().Error()
		}

		log.Msg("HTTP request")
	}
}

// Tracing middleware for OpenTelemetry tracing.
func Tracing() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		// Extract trace context from headers
		ctx = tracing.Extract(ctx, propagation.HeaderCarrier(c.Request.Header))

		// Create span
		ctx, span := tracing.StartSpan(ctx, "HTTP "+c.Request.Method+" "+c.FullPath())
		defer span.End()

		// Add request attributes
		span.SetAttributes(
			tracing.StringAttribute("http.method", c.Request.Method),
			tracing.StringAttribute("http.path", c.Request.URL.Path),
			tracing.StringAttribute("http.route", c.FullPath()),
		)

		c.Request = c.Request.WithContext(ctx)
		c.Set("trace_context", ctx)

		c.Next()

		// Add response attributes
		span.SetAttributes(
			tracing.IntAttribute("http.status_code", c.Writer.Status()),
		)

		if len(c.Errors) > 0 {
			span.RecordError(c.Errors.Last())
		}
	}
}

// Metrics middleware for Prometheus metrics.
func Metrics() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		duration := time.Since(start).Milliseconds()
		status := c.Writer.Status()

		// Record metrics
		metrics.HTTPRequestsTotal.WithLabelValues(
			c.Request.Method,
			c.FullPath(),
			http.StatusText(status),
		).Inc()

		metrics.HTTPRequestDurationMs.WithLabelValues(
			c.Request.Method,
			c.FullPath(),
		).Observe(float64(duration))

		if c.Writer.Status() >= 400 {
			metrics.HTTPErrorsTotal.WithLabelValues(
				c.Request.Method,
				c.FullPath(),
				c.Errors.Last().Error(),
			).Inc()
		}
	}
}

// CORS middleware for Cross-Origin Resource Sharing.
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization, X-Request-ID, X-Correlation-ID, Idempotency-Key")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Max-Age", "86400")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// RequestIDPropagation extracts and propagates request/correlation IDs.
func RequestIDPropagation() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader(RequestIDHeader)
		if requestID == "" {
			requestID = uuid.New().String()
		}
		c.Set("request_id", requestID)
		c.Header(RequestIDHeader, requestID)

		correlationID := c.GetHeader(CorrelationIDHeader)
		if correlationID == "" {
			correlationID = uuid.New().String()
		}
		c.Set("correlation_id", correlationID)
		c.Header(CorrelationIDHeader, correlationID)

		// Extract traceparent header for W3C trace context
		traceParent := c.GetHeader("traceparent")
		if traceParent != "" {
			c.Set("traceparent", traceParent)
		}

		c.Next()
	}
}

// Recovery middleware for panic recovery.
func Recovery() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		// Log the panic
		// log.Error().Interface("panic", recovered).Msg("Panic recovered")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Internal server error",
			},
		})
	})
}

// RequestSizeLimit limits the request body size.
func RequestSizeLimit(maxSize int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.ContentLength > maxSize {
			c.AbortWithStatusJSON(http.StatusRequestEntityTooLarge, gin.H{
				"error": gin.H{
					"code":    "REQUEST_TOO_LARGE",
					"message": "Request body too large",
				},
			})
			return
		}
		c.Next()
	}
}

// SecurityHeaders adds security-related headers.
func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Next()
	}
}
