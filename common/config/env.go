package config

import (
	"os"
	"strconv"
	"time"
)

// GetEnv gets an environment variable or returns a default value.
func GetEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// GetBoolEnv gets a boolean environment variable or returns a default value.
func GetBoolEnv(key string, defaultValue bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	
	boolValue, err := strconv.ParseBool(value)
	if err != nil {
		return defaultValue
	}
	return boolValue
}

// GetDurationEnv gets a duration environment variable or returns a default value.
// Expects value like "30s", "5m", "1h".
func GetDurationEnv(key string, defaultValue time.Duration) time.Duration {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	
	duration, err := time.ParseDuration(value)
	if err != nil {
		return defaultValue
	}
	return duration
}
