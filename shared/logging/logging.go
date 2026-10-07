package logging

import (
	"context"
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Level represents log levels.
type Level = zerolog.Level

const (
	DebugLevel = zerolog.DebugLevel
	InfoLevel  = zerolog.InfoLevel
	WarnLevel  = zerolog.WarnLevel
	ErrorLevel = zerolog.ErrorLevel
	FatalLevel = zerolog.FatalLevel
	PanicLevel = zerolog.PanicLevel
)

// Logger wraps zerolog.Logger with structured fields.
type Logger struct {
	logger zerolog.Logger
}

// Config holds logging configuration.
type Config struct {
	Level      Level
	Format     string // "json" or "console"
	Output     io.Writer
	Sampling   *SamplingConfig
	TimeFormat string
}

// SamplingConfig configures log sampling.
type SamplingConfig struct {
	Initial    uint32
	Thereafter uint32
}

// DefaultConfig returns a default logging configuration.
func DefaultConfig() Config {
	return Config{
		Level:      InfoLevel,
		Format:     "json",
		Output:     os.Stdout,
		TimeFormat: time.RFC3339Nano,
	}
}

// New creates a new Logger with the given config.
func New(cfg Config) *Logger {
	output := cfg.Output
	if cfg.Format == "console" {
		output = zerolog.ConsoleWriter{Out: cfg.Output, TimeFormat: cfg.TimeFormat}
	}

	zlog := zerolog.New(output).Level(cfg.Level).With().Timestamp().Logger()

	if cfg.Sampling != nil {
		zlog = zlog.Sample(&zerolog.BasicSampler{N: cfg.Sampling.Thereafter})
	}

	return &Logger{logger: zlog}
}

// InitGlobal initializes the global logger.
func InitGlobal(cfg Config) {
	logger := New(cfg)
	log.Logger = logger.logger
}

// FromContext extracts a logger from context, or returns the global logger.
func FromContext(ctx context.Context) *Logger {
	if ctx == nil {
		return &Logger{logger: log.Logger}
	}
	if logger, ok := ctx.Value(loggerKey{}).(*Logger); ok {
		return logger
	}
	return &Logger{logger: log.Logger}
}

// WithContext returns a new context with the logger.
func (l *Logger) WithContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, loggerKey{}, l)
}

// loggerKey is the context key for logger.
type loggerKey struct{}

// With adds fields to the logger.
func (l *Logger) With() *Logger {
	return &Logger{logger: l.logger.With().Logger()}
}

// WithField adds a single field.
func (l *Logger) WithField(key string, value any) *Logger {
	return &Logger{logger: l.logger.With().Interface(key, value).Logger()}
}

// WithFields adds multiple fields.
func (l *Logger) WithFields(fields map[string]any) *Logger {
	ctx := l.logger.With()
	for k, v := range fields {
		ctx = ctx.Interface(k, v)
	}
	return &Logger{logger: ctx.Logger()}
}

// WithError adds an error field.
func (l *Logger) WithError(err error) *Logger {
	if err == nil {
		return l
	}
	return &Logger{logger: l.logger.With().Err(err).Logger()}
}

// WithCorrelationID adds correlation ID.
func (l *Logger) WithCorrelationID(correlationID string) *Logger {
	return l.WithField("correlation_id", correlationID)
}

// WithTraceID adds trace ID.
func (l *Logger) WithTraceID(traceID string) *Logger {
	return l.WithField("trace_id", traceID)
}

// WithSpanID adds span ID.
func (l *Logger) WithSpanID(spanID string) *Logger {
	return l.WithField("span_id", spanID)
}

// WithExecutionID adds execution ID.
func (l *Logger) WithExecutionID(executionID string) *Logger {
	return l.WithField("execution_id", executionID)
}

// WithTenantID adds tenant ID.
func (l *Logger) WithTenantID(tenantID string) *Logger {
	return l.WithField("tenant_id", tenantID)
}

// WithRequestID adds request ID.
func (l *Logger) WithRequestID(requestID string) *Logger {
	return l.WithField("request_id", requestID)
}

// Debug logs at debug level.
func (l *Logger) Debug(msg string, args ...any) {
	l.logger.Debug().Msgf(msg, args...)
}

// Info logs at info level.
func (l *Logger) Info(msg string, args ...any) {
	l.logger.Info().Msgf(msg, args...)
}

// Warn logs at warn level.
func (l *Logger) Warn(msg string, args ...any) {
	l.logger.Warn().Msgf(msg, args...)
}

// Error logs at error level.
func (l *Logger) Error(msg string, args ...any) {
	l.logger.Error().Msgf(msg, args...)
}

// Fatal logs at fatal level and exits.
func (l *Logger) Fatal(msg string, args ...any) {
	l.logger.Fatal().Msgf(msg, args...)
}

// Debugf logs at debug level with formatting.
func (l *Logger) Debugf(format string, args ...any) {
	l.logger.Debug().Msgf(format, args...)
}

// Infof logs at info level with formatting.
func (l *Logger) Infof(format string, args ...any) {
	l.logger.Info().Msgf(format, args...)
}

// Warnf logs at warn level with formatting.
func (l *Logger) Warnf(format string, args ...any) {
	l.logger.Warn().Msgf(format, args...)
}

// Errorf logs at error level with formatting.
func (l *Logger) Errorf(format string, args ...any) {
	l.logger.Error().Msgf(format, args...)
}

// Fatalf logs at fatal level with formatting and exits.
func (l *Logger) Fatalf(format string, args ...any) {
	l.logger.Fatal().Msgf(format, args...)
}

// StdLogger returns a standard library compatible logger.
func (l *Logger) StdLogger() *StdLogger {
	return &StdLogger{logger: l}
}

// StdLogger wraps Logger to implement io.Writer for standard library logging.
type StdLogger struct {
	logger *Logger
}

func (s *StdLogger) Write(p []byte) (n int, err error) {
	s.logger.logger.Info().Msg(string(p))
	return len(p), nil
}

// NewContextLogger creates a logger with context values.
func NewContextLogger(ctx context.Context, base *Logger) *Logger {
	if base == nil {
		base = &Logger{logger: log.Logger}
	}
	l := base
	if correlationID, ok := ctx.Value("correlation_id").(string); ok && correlationID != "" {
		l = l.WithCorrelationID(correlationID)
	}
	if traceID, ok := ctx.Value("trace_id").(string); ok && traceID != "" {
		l = l.WithTraceID(traceID)
	}
	if spanID, ok := ctx.Value("span_id").(string); ok && spanID != "" {
		l = l.WithSpanID(spanID)
	}
	if executionID, ok := ctx.Value("execution_id").(string); ok && executionID != "" {
		l = l.WithExecutionID(executionID)
	}
	if tenantID, ok := ctx.Value("tenant_id").(string); ok && tenantID != "" {
		l = l.WithTenantID(tenantID)
	}
	if requestID, ok := ctx.Value("request_id").(string); ok && requestID != "" {
		l = l.WithRequestID(requestID)
	}
	return l
}

// WithComponent adds a component field.
func (l *Logger) WithComponent(component string) *Logger {
	return l.WithField("component", component)
}

// WithOperation adds an operation field.
func (l *Logger) WithOperation(operation string) *Logger {
	return l.WithField("operation", operation)
}

// WithDuration adds a duration field in milliseconds.
func (l *Logger) WithDuration(d time.Duration) *Logger {
	return l.WithField("duration_ms", d.Milliseconds())
}
