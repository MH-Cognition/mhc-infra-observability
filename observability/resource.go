// Package observability: central OpenTelemetry Resource creation.
// Exactly ONE Resource is created here and reused by tracing, metrics, and logging.
// Uses a single schema version to avoid conflicting Schema URL errors.
package observability

import (
	"context"

	"github.com/MH-Cognition/mhc-infra-observability/config"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
)

// NewResource creates the single OpenTelemetry Resource for this process.
// Use ONE schema version only. Must be called once; pass the result to Init.
// Do not call resource.New, resource.Default, or resource.Merge anywhere else.
func NewResource(ctx context.Context, cfg *config.Config) (*resource.Resource, error) {
	attrs := []attribute.KeyValue{
		semconv.ServiceName(cfg.ServiceName),
		semconv.DeploymentEnvironment(cfg.Environment),
	}
	if cfg.ServiceVersion != "" {
		attrs = append(attrs, semconv.ServiceVersion(cfg.ServiceVersion))
	}
	return resource.New(ctx,
		resource.WithSchemaURL(semconv.SchemaURL),
		resource.WithAttributes(attrs...),
	)
}
