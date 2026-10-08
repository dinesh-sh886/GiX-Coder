package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"google.golang.org/grpc/reflection"

	"github.com/gix-coder/gix-coder/api/proto/gix/gateway/v1"
	"github.com/gix-coder/gix-coder/gateway/internal/api"
	"github.com/gix-coder/gix-coder/gateway/internal/audit"
	"github.com/gix-coder/gix-coder/gateway/internal/auth"
	"github.com/gix-coder/gix-coder/gateway/internal/authz"
	"github.com/gix-coder/gix-coder/gateway/internal/config"
	"github.com/gix-coder/gix-coder/gateway/internal/middleware"
	"github.com/gix-coder/gix-coder/gateway/internal/persistence"
	"github.com/gix-coder/gix-coder/gateway/internal/server"
	"github.com/gix-coder/gix-coder/shared/logging"
	"github.com/gix-coder/gix-coder/shared/metrics"
	"github.com/gix-coder/gix-coder/shared/tracing"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	// Initialize logging
	logCfg := logging.DefaultConfig()
	logCfg.Level = parseLogLevel(cfg.Logging.Level)
	logCfg.Format = cfg.Logging.Format
	logging.InitGlobal(logCfg)

	logger := logging.FromContext(context.Background())
	logger.Zerolog().Info().Msg("Starting Gateway service")

	// Initialize tracing
	traceCfg := tracing.DefaultConfig("gateway", cfg.Service.Environment)
	traceCfg.Enabled = cfg.Tracing.Enabled
	traceCfg.Endpoint = cfg.Tracing.Endpoint
	traceCfg.Sampler = cfg.Tracing.Sampler
	traceCfg.Ratio = cfg.Tracing.Ratio
	traceCfg.Exporter = cfg.Tracing.Exporter

	shutdownTracing, err := tracing.Init(traceCfg)
	if err != nil {
		logger.Zerolog().Fatal().Err(err).Msg("Failed to initialize tracing")
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := shutdownTracing(ctx); err != nil {
			logger.Zerolog().Error().Err(err).Msg("Failed to shutdown tracing")
		}
	}()

	// Set global propagator
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	// Initialize metrics
	metricsBuilder := metrics.NewMetricBuilder("gateway", "http")
	setupMetrics(metricsBuilder)

	// Initialize audit client
	auditClient := audit.NewClient()
	auditMiddleware := middleware.NewAuditMiddleware(auditClient)

	// Initialize persistence
	pgxPool, err := persistence.NewPostgresPool(cfg)
	if err != nil {
		logger.Zerolog().Fatal().Err(err).Msg("Failed to connect to database")
	}
	defer pgxPool.Close()

	// Run migrations
	if err := persistence.RunMigrations(pgxPool.Pool); err != nil {
		logger.Zerolog().Fatal().Err(err).Msg("Failed to run migrations")
	}

	// Initialize repositories
	executionRepo := persistence.NewExecutionRepository(pgxPool.Pool)
	idempotencyRepo := persistence.NewIdempotencyRepository(pgxPool.Pool)

	// Initialize auth
	authClient, err := auth.NewClient(auth.Config{
		JWKSURL:   cfg.Auth.JWKSURL,
		Issuer:    cfg.Auth.Issuer,
		Audience:  cfg.Auth.Audience,
		Algorithm: cfg.Auth.Algorithm,
	})
	if err != nil {
		logger.Zerolog().Fatal().Err(err).Msg("Failed to initialize auth client")
	}

	// Initialize authorization
	authzClient := authz.NewClient()

	// Initialize rate limiter (Redis-backed if Redis config provided, otherwise in-memory)
	var rateLimiter middleware.RateLimiterInterface
	if cfg.RateLimit.RedisHost != "" {
		redisLimiter, err := middleware.NewRedisRateLimiter(middleware.RateLimitConfig{
			Enabled:           cfg.RateLimit.Enabled,
			RequestsPerMinute: cfg.RateLimit.RequestsPerMinute,
			Burst:             cfg.RateLimit.Burst,
			RedisHost:         cfg.RateLimit.RedisHost,
			RedisPort:         cfg.RateLimit.RedisPort,
			RedisPassword:     cfg.RateLimit.RedisPassword,
			RedisDB:           cfg.RateLimit.RedisDB,
			FailOpen:          true, // Fail-open for development
		})
		if err != nil {
			logger.Zerolog().Warn().Err(err).Msg("Failed to initialize Redis rate limiter, falling back to in-memory")
			rateLimiter = middleware.NewRateLimiter(middleware.RateLimitConfig{
				Enabled:           cfg.RateLimit.Enabled,
				RequestsPerMinute: cfg.RateLimit.RequestsPerMinute,
				Burst:             cfg.RateLimit.Burst,
			})
		} else {
			rateLimiter = redisLimiter
			defer func() {
				if err := redisLimiter.Close(); err != nil {
					logger.Zerolog().Error().Err(err).Msg("Failed to close Redis rate limiter")
				}
			}()
		}
	} else {
		rateLimiter = middleware.NewRateLimiter(middleware.RateLimitConfig{
			Enabled:           cfg.RateLimit.Enabled,
			RequestsPerMinute: cfg.RateLimit.RequestsPerMinute,
			Burst:             cfg.RateLimit.Burst,
		})
	}

	// Initialize idempotency
	idempotencyMiddleware := middleware.NewIdempotencyMiddleware(idempotencyRepo, cfg.Idempotency.TTL)

	// Initialize authorization middleware
	authzMiddleware := middleware.NewAuthzMiddleware(authzClient)

	// Create HTTP server
	httpServer := server.NewHTTPServer(server.HTTPConfig{
		Host:         cfg.Server.HTTP.Host,
		Port:         cfg.Server.HTTP.Port,
		ReadTimeout:  cfg.Server.HTTP.ReadTimeout,
		WriteTimeout: cfg.Server.HTTP.WriteTimeout,
		IdleTimeout:  cfg.Server.HTTP.IdleTimeout,
	})

	// Setup HTTP routes
	router := setupHTTPRouter(gin.Default(), cfg, authClient, authzMiddleware, rateLimiter, idempotencyMiddleware, auditMiddleware, logger, executionRepo)

	// Setup HTTP server with metrics endpoint
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))
	httpServer.Engine().GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Create gRPC server
	grpcServer := server.NewGRPCServer(server.GRPCConfig{
		Host: cfg.Server.GRPC.Host,
		Port: cfg.Server.GRPC.Port,
	})

	// Register gRPC services
	grpcHandler := api.NewGRPCHandler(executionRepo)
	gatewayv1.RegisterGatewayServiceServer(grpcServer.Server(), grpcHandler)
	reflection.Register(grpcServer.Server())

	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start HTTP server
	go func() {
		logger.Zerolog().Info().Int("port", cfg.Server.HTTP.Port).Msg("Starting HTTP server")
		if err := httpServer.Start(); err != nil && err != http.ErrServerClosed {
			logger.Zerolog().Fatal().Err(err).Msg("HTTP server failed")
		}
	}()

	// Start gRPC server
	go func() {
		logger.Zerolog().Info().Int("port", cfg.Server.GRPC.Port).Msg("Starting gRPC server")
		lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", cfg.Server.GRPC.Host, cfg.Server.GRPC.Port))
		if err != nil {
			logger.Zerolog().Fatal().Err(err).Msg("Failed to listen on gRPC port")
		}
		if err := grpcServer.Serve(lis); err != nil {
			logger.Zerolog().Fatal().Err(err).Msg("gRPC server failed")
		}
	}()

	// Wait for shutdown signal
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	logger.Zerolog().Info().Msg("Shutdown signal received, gracefully shutting down...")

	cancel()

	// Graceful shutdown
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		logger.Zerolog().Error().Err(err).Msg("HTTP server shutdown error")
	}
	grpcServer.GracefulStop()

	logger.Zerolog().Info().Msg("Gateway service stopped")
}

func setupMetrics(builder *metrics.MetricBuilder) {
	// HTTP metrics
	builder.CounterVec("requests_total", "Total number of HTTP requests", []string{"method", "path", "status"})
	builder.HistogramVecMs("request_duration_ms", "HTTP request duration in milliseconds", []string{"method", "path"})
	builder.CounterVec("requests_errors_total", "Total number of HTTP request errors", []string{"method", "path", "error"})

	// Auth metrics
	builder.CounterVec("auth_success_total", "Total successful authentications", []string{})
	builder.CounterVec("auth_failure_total", "Total failed authentications", []string{"reason"})

	// Rate limit metrics
	builder.CounterVec("rate_limit_rejected_total", "Total rate limited requests", []string{"key"})

	// Idempotency metrics
	builder.CounterVec("idempotency_conflicts_total", "Total idempotency conflicts", []string{})
}

func setupHTTPRouter(
	engine *gin.Engine,
	cfg *config.Config,
	authClient *auth.Client,
	authzMiddleware *middleware.AuthzMiddleware,
	rateLimiter middleware.RateLimiterInterface,
	idempotencyMiddleware *middleware.IdempotencyMiddleware,
	auditMiddleware *middleware.AuditMiddleware,
	logger *logging.Logger,
	executionRepo *persistence.ExecutionRepository,
) *gin.Engine {
	// Recovery middleware
	engine.Use(gin.Recovery())

	// Request ID middleware
	engine.Use(middleware.RequestID())

	// Logging middleware
	engine.Use(middleware.Logging(logger))

	// Tracing middleware
	engine.Use(middleware.Tracing())

	// Metrics middleware
	engine.Use(middleware.Metrics())

	// CORS
	engine.Use(middleware.CORS())

	// Request ID propagation
	engine.Use(middleware.RequestIDPropagation())

	// Add execution repository to context
	engine.Use(func(c *gin.Context) {
		c.Set("execution_repo", executionRepo)
		c.Next()
	})

	// Public health endpoints (no auth)
	public := engine.Group("/")
	{
		public.GET("/health/live", api.HealthLive)
		public.GET("/health/ready", api.HealthReady)
	}

	// Protected routes (require authentication)
	protected := engine.Group("/")
	protected.Use(middleware.Auth(authClient))
	protected.Use(authzMiddleware.Middleware())
	protected.Use(rateLimiter.Middleware())
	protected.Use(idempotencyMiddleware.Middleware())
	protected.Use(auditMiddleware.Middleware())
	{
		protected.POST("/workflows/execute", api.ExecuteWorkflow)
		protected.GET("/workflows/:execution_id", api.GetExecution)
		protected.POST("/workflows/:execution_id:cancel", api.CancelExecution)
	}

	return engine
}

func parseLogLevel(level string) logging.Level {
	switch level {
	case "debug":
		return logging.DebugLevel
	case "info":
		return logging.InfoLevel
	case "warn":
		return logging.WarnLevel
	case "error":
		return logging.ErrorLevel
	case "fatal":
		return logging.FatalLevel
	default:
		return logging.InfoLevel
	}
}
