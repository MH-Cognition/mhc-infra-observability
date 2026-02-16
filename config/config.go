// Package config provides configuration loading from environment variables
// for the observability library. No hardcoded values; all settings come from env.
package config

import "os"

// Config holds observability configuration loaded from environment variables.
type Config struct {
	// ServiceName identifies the service in traces and logs (e.g., "order-service").
	// Env: OTEL_SERVICE_NAME
	ServiceName string

	// ServiceVersion is the service version (optional). Env: OTEL_SERVICE_VERSION
	ServiceVersion string

	// Environment is the deployment environment (e.g., "dev", "staging", "prod").
	// Env: OTEL_ENVIRONMENT
	Environment string

	// OtelEndpoint is the OTLP collector endpoint for trace/metric export (e.g., "localhost:4317").
	// Env: OTEL_EXPORTER_OTLP_ENDPOINT
	OtelEndpoint string
}

// Load reads configuration from environment variables.
// Uses production-safe defaults when env vars are unset.
func Load() *Config {
	serviceName := os.Getenv("OTEL_SERVICE_NAME")
	if serviceName == "" {
		serviceName = "unknown-service"
	}

	serviceVersion := os.Getenv("OTEL_SERVICE_VERSION")

	env := os.Getenv("OTEL_ENVIRONMENT")
	if env == "" {
		env = "development"
	}

	endpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	if endpoint == "" {
		endpoint = "localhost:4317"
	}

	return &Config{
		ServiceName:    serviceName,
		ServiceVersion: serviceVersion,
		Environment:     env,
		OtelEndpoint:   endpoint,
	}
}
