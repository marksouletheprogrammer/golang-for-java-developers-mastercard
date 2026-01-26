package config

import (
	"fmt"
	"strconv"

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
	EnableGRPC      bool
	EnableMetrics   bool
	EnableHealthz   bool
	EnableDebugMode bool
}

// LoadConfig loads configuration from environment variables with defaults.
// TODO: Part 1 - Implement configuration loading
func LoadConfig() (*Config, error) {
	// TODO: Create Config struct with values from environment variables
	// TODO: Use commonconfig.GetEnv(), commonconfig.GetBoolEnv(), commonconfig.GetDurationEnv() helpers
	// TODO: Set sensible defaults for all fields
	// TODO: Call cfg.Validate() before returning
	// TODO: Return config and any validation errors
	return nil, fmt.Errorf("not implemented")
}

// Validate checks that required configuration is present and valid.
// TODO: Part 1 - Implement configuration validation
func (c *Config) Validate() error {
	// TODO: Validate ports are numeric using strconv.Atoi()
	// TODO: Validate environment is one of: development, staging, production
	// TODO: Validate log level is one of: debug, info, warn, error
	// TODO: In production, validate that APIKey and JWTSecret are set
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
