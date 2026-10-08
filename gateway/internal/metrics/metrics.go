package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// HTTPRequestsTotal is a counter for total HTTP requests.
	HTTPRequestsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "gateway",
		Subsystem: "http",
		Name:      "requests_total",
		Help:      "Total number of HTTP requests",
	}, []string{"method", "path", "status"})

	// HTTPRequestDurationMs is a histogram for HTTP request duration in milliseconds.
	HTTPRequestDurationMs = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "gateway",
		Subsystem: "http",
		Name:      "request_duration_ms",
		Help:      "HTTP request duration in milliseconds",
		Buckets:   []float64{5, 10, 25, 50, 100, 250, 500, 1000, 2500, 5000, 10000},
	}, []string{"method", "path"})

	// HTTPErrorsTotal is a counter for HTTP request errors.
	HTTPErrorsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "gateway",
		Subsystem: "http",
		Name:      "errors_total",
		Help:      "Total number of HTTP request errors",
	}, []string{"method", "path", "error"})

	// AuthSuccessTotal is a counter for successful authentications.
	AuthSuccessTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "gateway",
		Subsystem: "auth",
		Name:      "success_total",
		Help:      "Total successful authentications",
	}, []string{})

	// AuthFailureTotal is a counter for failed authentications.
	AuthFailureTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "gateway",
		Subsystem: "auth",
		Name:      "failure_total",
		Help:      "Total failed authentications",
	}, []string{"reason"})

	// RateLimitRejectedTotal is a counter for rate limited requests.
	RateLimitRejectedTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "gateway",
		Subsystem: "rate_limit",
		Name:      "rejected_total",
		Help:      "Total rate limited requests",
	}, []string{"key"})

	// IdempotencyConflictsTotal is a counter for idempotency conflicts.
	IdempotencyConflictsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "gateway",
		Subsystem: "idempotency",
		Name:      "conflicts_total",
		Help:      "Total idempotency conflicts",
	}, []string{})

	// AuthzDenialsTotal is a counter for authorization denials.
	AuthzDenialsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "gateway",
		Subsystem: "authz",
		Name:      "denials_total",
		Help:      "Total authorization denials",
	}, []string{})
)
