package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Standard histogram buckets for latency metrics.
var (
	// LatencyBuckets are standard latency buckets in seconds.
	LatencyBuckets = []float64{
		0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10,
	}

	// LatencyBucketsMs are standard latency buckets in milliseconds.
	LatencyBucketsMs = []float64{
		5, 10, 25, 50, 100, 250, 500, 1000, 2500, 5000, 10000,
	}

	// SizeBuckets are standard size buckets in bytes.
	SizeBuckets = []float64{
		100, 1000, 10000, 100000, 1000000, 10000000, 100000000,
	}
)

// MetricBuilder helps create standardized metrics.
type MetricBuilder struct {
	namespace string
	subsystem string
}

func NewMetricBuilder(namespace, subsystem string) *MetricBuilder {
	return &MetricBuilder{
		namespace: namespace,
		subsystem: subsystem,
	}
}

// CounterVec creates a new CounterVec with standard labels.
func (m *MetricBuilder) CounterVec(name, help string, labels []string) *prometheus.CounterVec {
	return promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: m.namespace,
		Subsystem: m.subsystem,
		Name:      name,
		Help:      help,
	}, labels)
}

// GaugeVec creates a new GaugeVec with standard labels.
func (m *MetricBuilder) GaugeVec(name, help string, labels []string) *prometheus.GaugeVec {
	return promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: m.namespace,
		Subsystem: m.subsystem,
		Name:      name,
		Help:      help,
	}, labels)
}

// HistogramVec creates a new HistogramVec with standard latency buckets.
func (m *MetricBuilder) HistogramVec(name, help string, labels []string) *prometheus.HistogramVec {
	return promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: m.namespace,
		Subsystem: m.subsystem,
		Name:      name,
		Help:      help,
		Buckets:   LatencyBuckets,
	}, labels)
}

// HistogramVecMs creates a new HistogramVec with millisecond buckets.
func (m *MetricBuilder) HistogramVecMs(name, help string, labels []string) *prometheus.HistogramVec {
	return promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: m.namespace,
		Subsystem: m.subsystem,
		Name:      name,
		Help:      help,
		Buckets:   LatencyBucketsMs,
	}, labels)
}

// SummaryVec creates a new SummaryVec.
func (m *MetricBuilder) SummaryVec(name, help string, labels []string, objectives map[float64]float64) *prometheus.SummaryVec {
	return promauto.NewSummaryVec(prometheus.SummaryOpts{
		Namespace:  m.namespace,
		Subsystem:  m.subsystem,
		Name:       name,
		Help:       help,
		Objectives: objectives,
	}, labels)
}

// StandardLabels returns the standard label names for service metrics.
func StandardLabels() []string {
	return []string{"service", "environment"}
}

// ServiceLabels returns labels including service name.
func ServiceLabels(service string) []string {
	return []string{"service", "environment"}
}

// TenantLabels returns labels including tenant.
func TenantLabels() []string {
	return []string{"service", "environment", "tenant"}
}

// MethodLabels returns labels for HTTP/gRPC methods.
func MethodLabels() []string {
	return []string{"service", "environment", "method"}
}

// EndpointLabels returns labels for endpoints.
func EndpointLabels() []string {
	return []string{"service", "environment", "endpoint", "method"}
}

// StatusLabels returns labels including status code.
func StatusLabels() []string {
	return []string{"service", "environment", "status"}
}

// ErrorLabels returns labels including error type.
func ErrorLabels() []string {
	return []string{"service", "environment", "error_type"}
}

// ResultLabels returns labels including result (success/failure).
func ResultLabels() []string {
	return []string{"service", "environment", "result"}
}

// REDMetricBuilder provides RED (Rate, Errors, Duration) metrics.
type REDMetricBuilder struct {
	*MetricBuilder
}

func NewREDMetricBuilder(namespace, subsystem string) *REDMetricBuilder {
	return &REDMetricBuilder{MetricBuilder: NewMetricBuilder(namespace, subsystem)}
}

// RequestCounter creates a request counter.
func (m *REDMetricBuilder) RequestCounter() *prometheus.CounterVec {
	return m.CounterVec("requests_total", "Total number of requests", EndpointLabels())
}

// RequestErrors creates a request error counter.
func (m *REDMetricBuilder) RequestErrors() *prometheus.CounterVec {
	return m.CounterVec("request_errors_total", "Total number of request errors", ErrorLabels())
}

// RequestDuration creates a request duration histogram.
func (m *REDMetricBuilder) RequestDuration() *prometheus.HistogramVec {
	return m.HistogramVec("request_duration_seconds", "Request duration in seconds", EndpointLabels())
}

// RequestsInFlight creates an in-flight requests gauge.
func (m *REDMetricBuilder) RequestsInFlight() *prometheus.GaugeVec {
	return m.GaugeVec("requests_in_flight", "Number of requests currently in flight", MethodLabels())
}

// GRPCREDMetricBuilder provides RED metrics for gRPC.
type GRPCREDMetricBuilder struct {
	*MetricBuilder
}

func NewGRPCREDMetricBuilder(namespace, subsystem string) *GRPCREDMetricBuilder {
	return &GRPCREDMetricBuilder{MetricBuilder: NewMetricBuilder(namespace, subsystem)}
}

// GRPCRequestCounter creates a gRPC request counter.
func (m *GRPCREDMetricBuilder) GRPCRequestCounter() *prometheus.CounterVec {
	return m.CounterVec("grpc_requests_total", "Total number of gRPC requests", []string{"service", "environment", "method", "status"})
}

// GRPCRequestDuration creates a gRPC request duration histogram.
func (m *GRPCREDMetricBuilder) GRPCRequestDuration() *prometheus.HistogramVec {
	return m.HistogramVec("grpc_request_duration_seconds", "gRPC request duration in seconds", []string{"service", "environment", "method"})
}

// GRPCRequestsInFlight creates a gRPC in-flight requests gauge.
func (m *GRPCREDMetricBuilder) GRPCRequestsInFlight() *prometheus.GaugeVec {
	return m.GaugeVec("grpc_requests_in_flight", "Number of gRPC requests currently in flight", []string{"service", "environment", "method"})
}

// BusinessMetricBuilder provides business-level metrics.
type BusinessMetricBuilder struct {
	*MetricBuilder
}

func NewBusinessMetricBuilder(namespace, subsystem string) *BusinessMetricBuilder {
	return &BusinessMetricBuilder{MetricBuilder: NewMetricBuilder(namespace, subsystem)}
}

// ExecutionsCounter creates a workflow executions counter.
func (m *BusinessMetricBuilder) ExecutionsCounter() *prometheus.CounterVec {
	return m.CounterVec("executions_total", "Total number of workflow executions", []string{"service", "tenant", "status"})
}

// ExecutionDuration creates an execution duration histogram.
func (m *BusinessMetricBuilder) ExecutionDuration() *prometheus.HistogramVec {
	return m.HistogramVec("execution_duration_seconds", "Workflow execution duration in seconds", []string{"service", "tenant", "workflow_type", "status"})
}

// TokensConsumed creates a token consumption counter.
func (m *BusinessMetricBuilder) TokensConsumed() *prometheus.CounterVec {
	return m.CounterVec("tokens_consumed_total", "Total tokens consumed", []string{"service", "tenant", "model"})
}

// SandboxExecutions creates a sandbox executions counter.
func (m *BusinessMetricBuilder) SandboxExecutions() *prometheus.CounterVec {
	return m.CounterVec("sandbox_executions_total", "Total sandbox executions", []string{"service", "tenant", "status"})
}

// ActiveTenants creates an active tenants gauge.
func (m *BusinessMetricBuilder) ActiveTenants() *prometheus.GaugeVec {
	return m.GaugeVec("active_tenants", "Current active tenants", []string{"service"})
}

// SystemMetricBuilder provides system-level metrics.
type SystemMetricBuilder struct {
	*MetricBuilder
}

func NewSystemMetricBuilder(namespace, subsystem string) *SystemMetricBuilder {
	return &SystemMetricBuilder{MetricBuilder: NewMetricBuilder(namespace, subsystem)}
}

// CPUSeconds creates a CPU usage counter.
func (m *SystemMetricBuilder) CPUSeconds() *prometheus.CounterVec {
	return m.CounterVec("process_cpu_seconds_total", "Total CPU time consumed", []string{"service", "instance"})
}

// ResidentMemory creates a resident memory gauge.
func (m *SystemMetricBuilder) ResidentMemory() *prometheus.GaugeVec {
	return m.GaugeVec("process_resident_memory_bytes", "Resident memory in bytes", []string{"service", "instance"})
}

// OpenFDs creates an open file descriptors gauge.
func (m *SystemMetricBuilder) OpenFDs() *prometheus.GaugeVec {
	return m.GaugeVec("process_open_fds", "Number of open file descriptors", []string{"service", "instance"})
}

// Goroutines creates a goroutines gauge (Go specific).
func (m *SystemMetricBuilder) Goroutines() *prometheus.GaugeVec {
	return m.GaugeVec("go_goroutines", "Number of goroutines", []string{"service", "instance"})
}

// GoGC creates a GC duration histogram.
func (m *SystemMetricBuilder) GoGC() *prometheus.HistogramVec {
	return m.HistogramVec("go_gc_duration_seconds", "Go GC duration in seconds", []string{"service", "instance"})
}
