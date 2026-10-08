package tracing

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	"go.opentelemetry.io/otel/trace"
)

// Exporter types for tracing.
const (
	ExporterOTLPgRPC = "otlpgrpc"
	ExporterOTLPHTTP = "otlphttp"
	ExporterStdout   = "stdout"
)

// Config holds tracing configuration.
type Config struct {
	Enabled     bool
	ServiceName string
	Environment string
	Endpoint    string
	Sampler     string
	Ratio       float64
	Exporter    string
	Insecure    bool
}

// DefaultConfig returns default tracing config.
func DefaultConfig(serviceName, environment string) Config {
	return Config{
		Enabled:     true,
		ServiceName: serviceName,
		Environment: environment,
		Endpoint:    "",
		Sampler:     "parentbased_traceidratio",
		Ratio:       0.1,
		Exporter:    ExporterOTLPgRPC,
		Insecure:    true,
	}
}

// Tracer wraps the OpenTelemetry tracer.
type Tracer struct {
	tracer trace.Tracer
}

// Init initializes OpenTelemetry tracing.
func Init(cfg Config) (func(context.Context) error, error) {
	if !cfg.Enabled {
		return func(context.Context) error { return nil }, nil
	}

	ctx := context.Background()

	// Create resource
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName(cfg.ServiceName),
			semconv.DeploymentEnvironment(cfg.Environment),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("create resource: %w", err)
	}

	// Create trace exporter
	var traceExporter sdktrace.SpanExporter

	switch cfg.Exporter {
	case ExporterOTLPgRPC:
		traceExporter, err = otlptracegrpc.New(ctx,
			otlptracegrpc.WithEndpoint(cfg.Endpoint),
			otlptracegrpc.WithInsecure(),
		)
	case ExporterOTLPHTTP:
		traceExporter, err = otlptracehttp.New(ctx,
			otlptracehttp.WithEndpoint(cfg.Endpoint),
			otlptracehttp.WithInsecure(),
		)
	default:
		return nil, fmt.Errorf("unsupported trace exporter: %s", cfg.Exporter)
	}

	if err != nil {
		return nil, fmt.Errorf("create trace exporter: %w", err)
	}

	// Create sampler
	sampler := createSampler(cfg)

	// Create trace provider
	traceProvider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(traceExporter),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sampler),
	)

	otel.SetTracerProvider(traceProvider)

	// Create metric exporter
	var metricExporter metric.Exporter
	switch cfg.Exporter {
	case ExporterOTLPgRPC:
		metricExporter, err = otlpmetricgrpc.New(ctx,
			otlpmetricgrpc.WithEndpoint(cfg.Endpoint),
			otlpmetricgrpc.WithInsecure(),
		)
	case ExporterOTLPHTTP:
		metricExporter, err = otlpmetrichttp.New(ctx,
			otlpmetrichttp.WithEndpoint(cfg.Endpoint),
			otlpmetrichttp.WithInsecure(),
		)
	default:
		return nil, fmt.Errorf("unsupported metric exporter: %s", cfg.Exporter)
	}

	if err != nil {
		return nil, fmt.Errorf("create metric exporter: %w", err)
	}

	// Create metric provider
	metricProvider := metric.NewMeterProvider(
		metric.WithReader(metric.NewPeriodicReader(metricExporter)),
		metric.WithResource(res),
	)

	otel.SetMeterProvider(metricProvider)

	// Set propagators
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	return func(ctx context.Context) error {
		var errs []error
		if err := traceProvider.Shutdown(ctx); err != nil {
			errs = append(errs, fmt.Errorf("shutdown trace provider: %w", err))
		}
		if err := metricProvider.Shutdown(ctx); err != nil {
			errs = append(errs, fmt.Errorf("shutdown metric provider: %w", err))
		}
		if len(errs) > 0 {
			return fmt.Errorf("%v", errs)
		}
		return nil
	}, nil
}

func createSampler(cfg Config) sdktrace.Sampler {
	switch cfg.Sampler {
	case "always_on":
		return sdktrace.AlwaysSample()
	case "always_off":
		return sdktrace.NeverSample()
	case "traceidratio":
		return sdktrace.TraceIDRatioBased(cfg.Ratio)
	case "parentbased_traceidratio":
		return sdktrace.ParentBased(sdktrace.TraceIDRatioBased(cfg.Ratio))
	case "parentbased_always_on":
		return sdktrace.ParentBased(sdktrace.AlwaysSample())
	case "parentbased_always_off":
		return sdktrace.ParentBased(sdktrace.NeverSample())
	default:
		return sdktrace.ParentBased(sdktrace.TraceIDRatioBased(cfg.Ratio))
	}
}

// NewTracer creates a new tracer for the given service.
func NewTracer(serviceName string) *Tracer {
	return &Tracer{
		tracer: otel.Tracer(serviceName),
	}
}

// Start starts a new span.
func (t *Tracer) Start(ctx context.Context, name string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	return t.tracer.Start(ctx, name, opts...)
}

// StartWithAttributes starts a new span with attributes.
func (t *Tracer) StartWithAttributes(ctx context.Context, name string, attrs ...attribute.KeyValue) (context.Context, trace.Span) {
	return t.tracer.Start(ctx, name, trace.WithAttributes(attrs...))
}

// Propagation keys
const (
	TraceParentHeader = "traceparent"
	TraceStateHeader  = "tracestate"
	BaggageHeader     = "baggage"
)

// Inject injects trace context into carrier.
func Inject(ctx context.Context, carrier propagation.TextMapCarrier) {
	otel.GetTextMapPropagator().Inject(ctx, carrier)
}

// Extract extracts trace context from carrier.
func Extract(ctx context.Context, carrier propagation.TextMapCarrier) context.Context {
	return otel.GetTextMapPropagator().Extract(ctx, carrier)
}

// StartSpan starts a new span with the given name.
func StartSpan(ctx context.Context, name string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	return otel.Tracer("gateway").Start(ctx, name, opts...)
}

// SpanFromContext returns the span from context.
func SpanFromContext(ctx context.Context) trace.Span {
	return trace.SpanFromContext(ctx)
}

// TraceIDFromContext returns the trace ID from context.
func TraceIDFromContext(ctx context.Context) string {
	span := trace.SpanFromContext(ctx)
	if span.SpanContext().HasTraceID() {
		return span.SpanContext().TraceID().String()
	}
	return ""
}

// SpanIDFromContext returns the span ID from context.
func SpanIDFromContext(ctx context.Context) string {
	span := trace.SpanFromContext(ctx)
	if span.SpanContext().HasSpanID() {
		return span.SpanContext().SpanID().String()
	}
	return ""
}

// AddEvent adds an event to the current span.
func AddEvent(ctx context.Context, name string, attrs ...attribute.KeyValue) {
	span := trace.SpanFromContext(ctx)
	if span.IsRecording() {
		span.AddEvent(name, trace.WithAttributes(attrs...))
	}
}

// RecordError records an error on the current span.
func RecordError(ctx context.Context, err error, attrs ...attribute.KeyValue) {
	span := trace.SpanFromContext(ctx)
	if span.IsRecording() {
		span.RecordError(err, trace.WithAttributes(attrs...))
	}
}

// SetAttributes sets attributes on the current span.
func SetAttributes(ctx context.Context, attrs ...attribute.KeyValue) {
	span := trace.SpanFromContext(ctx)
	if span.IsRecording() {
		span.SetAttributes(attrs...)
	}
}

// StringAttribute creates a string attribute for tracing.
func StringAttribute(key, value string) attribute.KeyValue {
	return attribute.String(key, value)
}

// IntAttribute creates an int attribute for tracing.
func IntAttribute(key string, value int) attribute.KeyValue {
	return attribute.Int(key, value)
}

// WithSpan runs a function with a new span.
func WithSpan(ctx context.Context, tracer *Tracer, name string, fn func(context.Context) error, attrs ...attribute.KeyValue) error {
	ctx, span := tracer.Start(ctx, name, trace.WithAttributes(attrs...))
	defer span.End()

	if err := fn(ctx); err != nil {
		RecordError(ctx, err)
		span.SetAttributes(attribute.String("error", "true"))
		return err
	}
	return nil
}

// HTTPClientTrace holds trace state for HTTP client.
type HTTPClientTrace struct {
	ctx context.Context
}

// NewHTTPClientTrace creates an HTTP client trace for tracing.
func NewHTTPClientTrace(ctx context.Context) *HTTPClientTrace {
	return &HTTPClientTrace{ctx: ctx}
}

// GetContext returns the context.
func (h *HTTPClientTrace) GetContext() context.Context {
	return h.ctx
}

// WithSpanContext returns a new HTTPClientTrace with span context.
func (h *HTTPClientTrace) WithSpanContext(ctx context.Context) *HTTPClientTrace {
	return &HTTPClientTrace{ctx: ctx}
}
