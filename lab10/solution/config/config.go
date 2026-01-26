package config

import (
	"fmt"
	"strconv"
	"time"

	commonconfig "golang-for-java-developers-training/common/config"
)

// Config holds all application configuration.
// Configuration comes from environment variables with sensible defaults.
type Config struct {
	// Server configuration
	HTTPPort string
	GRPCPort string
	
	// Environment
	Environment string // dev, staging, production
	LogLevel    string // debug, info, warn, error
	
	// Timeouts
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
	
	// Feature flags
	Features FeatureFlags
	
	// Secrets (never log these)
	APIKey    string
	JWTSecret string
}

// FeatureFlags contains toggleable features.
// These can be enabled/disabled without recompiling.
type FeatureFlags struct {
	EnableGRPC       bool
	EnableMetrics    bool
	EnableHealthz    bool
	EnableDebugMode  bool
}

// LoadConfig loads configuration from environment variables with defaults.
func LoadConfig() (*Config, error) {
	cfg := &Config{
		HTTPPort:     commonconfig.GetEnv("HTTP_PORT", "8080"),
		GRPCPort:     commonconfig.GetEnv("GRPC_PORT", "9090"),
		Environment:  commonconfig.GetEnv("ENVIRONMENT", "development"),
		LogLevel:     commonconfig.GetEnv("LOG_LEVEL", "info"),
		ReadTimeout:  commonconfig.GetDurationEnv("READ_TIMEOUT", 15*time.Second),
		WriteTimeout: commonconfig.GetDurationEnv("WRITE_TIMEOUT", 15*time.Second),
		IdleTimeout:  commonconfig.GetDurationEnv("IDLE_TIMEOUT", 60*time.Second),
		Features: FeatureFlags{
			EnableGRPC:      commonconfig.GetBoolEnv("ENABLE_GRPC", true),
			EnableMetrics:   commonconfig.GetBoolEnv("ENABLE_METRICS", true),
			EnableHealthz:   commonconfig.GetBoolEnv("ENABLE_HEALTHZ", true),
			EnableDebugMode: commonconfig.GetBoolEnv("ENABLE_DEBUG", false),
		},
		APIKey:    commonconfig.GetEnv("API_KEY", ""),
		JWTSecret: commonconfig.GetEnv("JWT_SECRET", ""),
	}
	
	// Validate configuration
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}
	
	return cfg, nil
}

// Validate checks that required configuration is present and valid.
func (c *Config) Validate() error {
	// Validate ports are numeric
	if _, err := strconv.Atoi(c.HTTPPort); err != nil {
		return fmt.Errorf("invalid HTTP_PORT: %s", c.HTTPPort)
	}
	if _, err := strconv.Atoi(c.GRPCPort); err != nil {
		return fmt.Errorf("invalid GRPC_PORT: %s", c.GRPCPort)
	}
	
	// Validate environment
	validEnvs := map[string]bool{"development": true, "staging": true, "production": true}
	if !validEnvs[c.Environment] {
		return fmt.Errorf("invalid ENVIRONMENT: %s (must be development, staging, or production)", c.Environment)
	}
	
	// Validate log level
	validLevels := map[string]bool{"debug": true, "info": true, "warn": true, "error": true}
	if !validLevels[c.LogLevel] {
		return fmt.Errorf("invalid LOG_LEVEL: %s (must be debug, info, warn, or error)", c.LogLevel)
	}
	
	// Validate secrets in production
	if c.Environment == "production" {
		if c.APIKey == "" {
			return fmt.Errorf("API_KEY is required in production")
		}
		if c.JWTSecret == "" {
			return fmt.Errorf("JWT_SECRET is required in production")
		}
	}
	
	return nil
}

// IsProduction returns true if running in production environment.
func (c *Config) IsProduction() bool {
	return c.Environment == "production"
}

// IsDevelopment returns true if running in development environment.
func (c *Config) IsDevelopment() bool {
	return c.Environment == "development"
}
