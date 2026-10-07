package config

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// Config represents the base configuration structure.
type Config struct {
	Service  ServiceConfig
	Server   ServerConfig
	Logging  LoggingConfig
	Tracing  TracingConfig
	Metrics  MetricsConfig
	Health   HealthConfig
	Database DatabaseConfig
	Redis    RedisConfig
	Vault    VaultConfig
	Features FeatureFlags
	Security SecurityConfig
}

// ServiceConfig holds service identification config.
type ServiceConfig struct {
	Name        string
	Version     string
	Environment string
}

// ServerConfig holds HTTP/gRPC server config.
type ServerConfig struct {
	HTTP HTTPServerConfig
	GRPC GRPCServerConfig
}

type HTTPServerConfig struct {
	Port         int
	Host         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

type GRPCServerConfig struct {
	Port int
	Host string
}

// LoggingConfig holds logging config.
type LoggingConfig struct {
	Level  string
	Format string
}

// TracingConfig holds tracing config.
type TracingConfig struct {
	Enabled  bool
	Sampler  string
	Ratio    float64
	Exporter string
	Endpoint string
}

// MetricsConfig holds metrics config.
type MetricsConfig struct {
	Enabled bool
	Path    string
	Port    int
}

// HealthConfig holds health check config.
type HealthConfig struct {
	Liveness  string
	Readiness string
}

// DatabaseConfig holds PostgreSQL config.
type DatabaseConfig struct {
	Host     string
	Port     int
	Database string
	User     string
	Password string
	SSLMode  string
	Schema   string
}

// RedisConfig holds Redis config.
type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
}

// VaultConfig holds Vault config.
type VaultConfig struct {
	Address string
	Path    string
	Token   string
}

// FeatureFlags holds feature flag config.
type FeatureFlags struct {
	NewRouting   bool
	EnhancedAuth bool
}

// SecurityConfig holds security config.
type SecurityConfig struct {
	TLS  TLSConfig
	CORS CORSConfig
}

type TLSConfig struct {
	Enabled  bool
	CertFile string
	KeyFile  string
}

type CORSConfig struct {
	AllowedOrigins []string
	AllowedMethods []string
	AllowedHeaders []string
}

// Load loads configuration from files, environment variables, and secrets.
func Load(serviceName string) (*Config, error) {
	v := viper.New()

	// Set defaults
	setDefaults(v, serviceName)

	// Load config files (priority: config file > env vars > secrets)
	if err := loadConfigFiles(v, serviceName); err != nil {
		return nil, fmt.Errorf("load config files: %w", err)
	}

	// Bind environment variables
	if err := bindEnvVars(v); err != nil {
		return nil, fmt.Errorf("bind env vars: %w", err)
	}

	// Load secrets from Vault (if configured)
	if err := loadSecrets(v); err != nil {
		return nil, fmt.Errorf("load secrets: %w", err)
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}

	// Validate required fields
	if err := validate(&cfg); err != nil {
		return nil, fmt.Errorf("validate config: %w", err)
	}

	return &cfg, nil
}

// MustLoad loads config or panics.
func MustLoad(serviceName string) *Config {
	cfg, err := Load(serviceName)
	if err != nil {
		panic(fmt.Sprintf("failed to load config for %s: %v", serviceName, err))
	}
	return cfg
}

func setDefaults(v *viper.Viper, serviceName string) {
	// Service defaults
	v.SetDefault("service.name", serviceName)
	v.SetDefault("service.version", "dev")
	v.SetDefault("service.environment", "development")

	// Server defaults
	v.SetDefault("server.http.port", 8080)
	v.SetDefault("server.http.host", "0.0.0.0")
	v.SetDefault("server.http.read_timeout", "30s")
	v.SetDefault("server.http.write_timeout", "30s")
	v.SetDefault("server.http.idle_timeout", "120s")
	v.SetDefault("server.grpc.port", 9090)
	v.SetDefault("server.grpc.host", "0.0.0.0")

	// Logging defaults
	v.SetDefault("logging.level", "info")
	v.SetDefault("logging.format", "json")

	// Tracing defaults
	v.SetDefault("tracing.enabled", true)
	v.SetDefault("tracing.sampler", "parentbased_traceidratio")
	v.SetDefault("tracing.ratio", 0.1)
	v.SetDefault("tracing.exporter", "otlp")
	v.SetDefault("tracing.endpoint", "")

	// Metrics defaults
	v.SetDefault("metrics.enabled", true)
	v.SetDefault("metrics.path", "/metrics")
	v.SetDefault("metrics.port", 9091)

	// Health defaults
	v.SetDefault("health.liveness", "/health/live")
	v.SetDefault("health.readiness", "/health/ready")

	// Database defaults
	v.SetDefault("database.host", "localhost")
	v.SetDefault("database.port", 5432)
	v.SetDefault("database.sslmode", "disable")

	// Redis defaults
	v.SetDefault("redis.host", "localhost")
	v.SetDefault("redis.port", 6379)
	v.SetDefault("redis.db", 0)

	// Vault defaults
	v.SetDefault("vault.address", "http://localhost:8200")
	v.SetDefault("vault.path", "secret/"+serviceName)

	// Feature flags defaults
	v.SetDefault("features.new_routing", false)
	v.SetDefault("features.enhanced_auth", true)

	// Security defaults
	v.SetDefault("security.tls.enabled", false)
	v.SetDefault("security.cors.allowed_origins", []string{"*"})
	v.SetDefault("security.cors.allowed_methods", []string{"GET", "POST", "PUT", "DELETE"})
	v.SetDefault("security.cors.allowed_headers", []string{"*"})
}

func loadConfigFiles(v *viper.Viper, serviceName string) error {
	// Config file paths (in order of priority)
	configPaths := []string{
		".",
		"./configs",
		"/etc/gix-coder",
		"$HOME/.config/gix-coder",
	}

	configNames := []string{
		"config",
		serviceName,
		fmt.Sprintf("%s.%s", serviceName, getEnv("ENVIRONMENT", "development")),
	}

	for _, path := range configPaths {
		for _, name := range configNames {
			v.AddConfigPath(path)
			v.SetConfigName(name)
			v.SetConfigType("yaml")

			if err := v.MergeInConfig(); err != nil {
				var cfgErr viper.ConfigFileNotFoundError
				if !errors.As(err, &cfgErr) {
					return err
				}
			}
		}
	}
	return nil
}

func bindEnvVars(v *viper.Viper) error {
	v.SetEnvPrefix("GIX")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	// Explicit bindings for nested config
	if err := v.BindEnv("service.name", "GIX_SERVICE_NAME"); err != nil {
		return err
	}
	if err := v.BindEnv("service.version", "GIX_SERVICE_VERSION"); err != nil {
		return err
	}
	if err := v.BindEnv("service.environment", "GIX_ENVIRONMENT"); err != nil {
		return err
	}
	if err := v.BindEnv("server.http.port", "GIX_SERVER_HTTP_PORT"); err != nil {
		return err
	}
	if err := v.BindEnv("server.http.host", "GIX_SERVER_HTTP_HOST"); err != nil {
		return err
	}
	if err := v.BindEnv("server.grpc.port", "GIX_SERVER_GRPC_PORT"); err != nil {
		return err
	}
	if err := v.BindEnv("logging.level", "GIX_LOG_LEVEL"); err != nil {
		return err
	}
	if err := v.BindEnv("logging.format", "GIX_LOG_FORMAT"); err != nil {
		return err
	}
	if err := v.BindEnv("tracing.enabled", "GIX_TRACING_ENABLED"); err != nil {
		return err
	}
	if err := v.BindEnv("tracing.endpoint", "GIX_TRACING_ENDPOINT"); err != nil {
		return err
	}
	if err := v.BindEnv("database.host", "GIX_POSTGRES_HOST"); err != nil {
		return err
	}
	if err := v.BindEnv("database.port", "GIX_POSTGRES_PORT"); err != nil {
		return err
	}
	if err := v.BindEnv("database.database", "GIX_POSTGRES_DATABASE"); err != nil {
		return err
	}
	if err := v.BindEnv("database.user", "GIX_POSTGRES_USER"); err != nil {
		return err
	}
	if err := v.BindEnv("redis.host", "GIX_REDIS_HOST"); err != nil {
		return err
	}
	if err := v.BindEnv("redis.port", "GIX_REDIS_PORT"); err != nil {
		return err
	}
	if err := v.BindEnv("vault.address", "GIX_VAULT_ADDRESS"); err != nil {
		return err
	}
	if err := v.BindEnv("vault.token", "GIX_VAULT_TOKEN"); err != nil {
		return err
	}
	return nil
}

func loadSecrets(v *viper.Viper) error {
	// Secrets are loaded from Vault at runtime via the Vault agent
	// This is a placeholder for the Vault integration
	vaultAddr := v.GetString("vault.address")
	if vaultAddr == "" {
		return nil
	}
	// Vault integration would go here
	// For now, we rely on environment variables for secrets
	return nil
}

func validate(cfg *Config) error {
	if cfg.Service.Name == "" {
		return fmt.Errorf("service.name is required")
	}
	if cfg.Server.HTTP.Port <= 0 || cfg.Server.HTTP.Port > 65535 {
		return fmt.Errorf("invalid HTTP port: %d", cfg.Server.HTTP.Port)
	}
	if cfg.Server.GRPC.Port <= 0 || cfg.Server.GRPC.Port > 65535 {
		return fmt.Errorf("invalid gRPC port: %d", cfg.Server.GRPC.Port)
	}
	if cfg.Database.Host == "" {
		return fmt.Errorf("database.host is required")
	}
	return nil
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// GetDSN returns the PostgreSQL connection string.
func (c *Config) GetDSN() string {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		c.Database.User,
		c.Database.Password,
		c.Database.Host,
		c.Database.Port,
		c.Database.Database,
		c.Database.SSLMode,
	)
	if c.Database.Schema != "" {
		dsn += "&search_path=" + c.Database.Schema
	}
	return dsn
}

// GetRedisAddr returns the Redis address.
func (c *Config) GetRedisAddr() string {
	return fmt.Sprintf("%s:%d", c.Redis.Host, c.Redis.Port)
}
