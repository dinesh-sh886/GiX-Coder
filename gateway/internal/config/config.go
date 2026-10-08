package config

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"

	"github.com/gix-coder/gix-coder/shared/config"
)

// Config extends the shared config with gateway-specific settings.
type Config struct {
	*config.Config

	Auth        AuthConfig
	RateLimit   RateLimitConfig
	Idempotency IdempotencyConfig
}

// GetDatabaseHost returns the database host.
func (c *Config) GetDatabaseHost() string {
	return c.Database.Host
}

// GetDatabasePort returns the database port.
func (c *Config) GetDatabasePort() int {
	return c.Database.Port
}

// GetDatabaseDatabase returns the database name.
func (c *Config) GetDatabaseDatabase() string {
	return c.Database.Database
}

// GetDatabaseUser returns the database user.
func (c *Config) GetDatabaseUser() string {
	return c.Database.User
}

// GetDatabasePassword returns the database password.
func (c *Config) GetDatabasePassword() string {
	return c.Database.Password
}

// GetDatabaseSSLMode returns the database SSL mode.
func (c *Config) GetDatabaseSSLMode() string {
	return c.Database.SSLMode
}

// GetDatabaseSchema returns the database schema.
func (c *Config) GetDatabaseSchema() string {
	return c.Database.Schema
}

// AuthConfig holds authentication configuration.
type AuthConfig struct {
	JWKSURL   string
	Issuer    string
	Audience  string
	Algorithm string
}

// RateLimitConfig holds rate limiting configuration.
type RateLimitConfig struct {
	Enabled           bool
	RequestsPerMinute int
	Burst             int
	RedisHost         string
	RedisPort         int
	RedisPassword     string
	RedisDB           int
}

// IdempotencyConfig holds idempotency configuration.
type IdempotencyConfig struct {
	Enabled bool
	TTL     time.Duration
}

// Load loads the gateway configuration.
func Load() (*Config, error) {
	baseCfg, err := config.Load("gateway")
	if err != nil {
		return nil, fmt.Errorf("load base config: %w", err)
	}

	v := viper.New()

	// Set gateway-specific defaults
	setGatewayDefaults(v)

	// Load config files
	if err := loadConfigFiles(v, "gateway"); err != nil {
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

	var gatewayCfg Config
	if err := v.Unmarshal(&gatewayCfg); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}

	gatewayCfg.Config = baseCfg

	// Validate gateway-specific config
	if err := validateGatewayConfig(&gatewayCfg); err != nil {
		return nil, fmt.Errorf("validate gateway config: %w", err)
	}

	return &gatewayCfg, nil
}

func setGatewayDefaults(v *viper.Viper) {
	// Auth defaults
	v.SetDefault("auth.jwks_url", "")
	v.SetDefault("auth.issuer", "")
	v.SetDefault("auth.audience", "")
	v.SetDefault("auth.algorithm", "RS256")

	// Rate limit defaults
	v.SetDefault("rate_limit.enabled", true)
	v.SetDefault("rate_limit.requests_per_minute", 100)
	v.SetDefault("rate_limit.burst", 200)
	v.SetDefault("rate_limit.redis_host", "localhost")
	v.SetDefault("rate_limit.redis_port", 6379)
	v.SetDefault("rate_limit.redis_password", "")
	v.SetDefault("rate_limit.redis_db", 0)

	// Idempotency defaults
	v.SetDefault("idempotency.enabled", true)
	v.SetDefault("idempotency.ttl", "24h")
}

func loadConfigFiles(v *viper.Viper, serviceName string) error {
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
				if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
					return err
				}
			}
		}
	}
	return nil
}

func bindEnvVars(v *viper.Viper) error {
	v.SetEnvPrefix("GATEWAY")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	// Explicit bindings for nested config
	bindings := []struct {
		key string
		env string
	}{
		{"auth.jwks_url", "GATEWAY_AUTH_JWKS_URL"},
		{"auth.issuer", "GATEWAY_AUTH_ISSUER"},
		{"auth.audience", "GATEWAY_AUTH_AUDIENCE"},
		{"auth.algorithm", "GATEWAY_AUTH_ALGORITHM"},
		{"rate_limit.enabled", "GATEWAY_RATE_LIMIT_ENABLED"},
		{"rate_limit.requests_per_minute", "GATEWAY_RATE_LIMIT_REQUESTS_PER_MINUTE"},
		{"rate_limit.burst", "GATEWAY_RATE_LIMIT_BURST"},
		{"rate_limit.redis_host", "GATEWAY_RATE_LIMIT_REDIS_HOST"},
		{"rate_limit.redis_port", "GATEWAY_RATE_LIMIT_REDIS_PORT"},
		{"rate_limit.redis_password", "GATEWAY_RATE_LIMIT_REDIS_PASSWORD"},
		{"rate_limit.redis_db", "GATEWAY_RATE_LIMIT_REDIS_DB"},
		{"idempotency.enabled", "GATEWAY_IDEMPOTENCY_ENABLED"},
		{"idempotency.ttl", "GATEWAY_IDEMPOTENCY_TTL"},
	}

	for _, b := range bindings {
		if err := v.BindEnv(b.key, b.env); err != nil {
			return fmt.Errorf("bind env %s: %w", b.key, err)
		}
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

func validateGatewayConfig(cfg *Config) error {
	if cfg.Auth.JWKSURL != "" {
		if _, err := url.Parse(cfg.Auth.JWKSURL); err != nil {
			return fmt.Errorf("invalid auth.jwks_url: %w", err)
		}
	}

	if cfg.Auth.Issuer != "" {
		if _, err := url.Parse(cfg.Auth.Issuer); err != nil {
			return fmt.Errorf("invalid auth.issuer: %w", err)
		}
	}

	if cfg.RateLimit.Enabled {
		if cfg.RateLimit.RequestsPerMinute <= 0 {
			return errors.New("rate_limit.requests_per_minute must be positive")
		}
		if cfg.RateLimit.Burst <= 0 {
			return errors.New("rate_limit.burst must be positive")
		}
	}

	if cfg.Idempotency.Enabled {
		if cfg.Idempotency.TTL <= 0 {
			return errors.New("idempotency.ttl must be positive")
		}
	}

	return nil
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
