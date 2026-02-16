// Package observability is the public facade for the mhc-infra-observability library.
// Services should import only this package; internal packages (config, tracing, logging, etc.)
// are implementation details and should not be imported directly by domain/use-case code.
package observability

import (
	"context"
	"fmt"
	"net/http"

	"github.com/MH-Cognition/mhc-infra-observability/config"
	"github.com/MH-Cognition/mhc-infra-observability/logging"
	"github.com/MH-Cognition/mhc-infra-observability/metrics"
	"github.com/MH-Cognition/mhc-infra-observability/tracing"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
)

// Init initializes the observability stack: tracer, propagator, logger.
// Returns a shutdown function that must be called before process exit (e.g., in main's defer).
func Init(ctx context.Context, cfg *config.Config) (func(context.Context) error, error) {
	shutdown, err := tracing.Init(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("init tracing: %w", err)
	}
	return shutdown, nil
}

// StartSpan starts a new span as a child of the current span in ctx.
// Returns the new context (with span) and the span. Caller must call span.End() when done.
// Safe to call when no tracer is configured (returns noop span).
func StartSpan(ctx context.Context, name string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	tracer := otel.Tracer("mhc-infra-observability")
	return tracer.Start(ctx, name, opts...)
}

// Logger returns the trace-aware structured logger.
func Logger(ctx context.Context) *logging.Logger {
	return logging.LoggerFromEnv()
}

// HandleError records the error in the current span and logs it.
// Does NOT map errors to HTTP status codes or define business error meanings.
// Call from use-case or transport layer when an error occurs that should be observed.
func HandleError(ctx context.Context, err error) {
	if err == nil {
		return
	}
	span := trace.SpanFromContext(ctx)
	span.RecordError(err)
	span.SetStatus(codes.Error, err.Error())
	span.SetAttributes(attribute.String("error", err.Error()))

	logger := Logger(ctx)
	logger.Error(ctx, "error observed", "error", err)
}

// HTTPMiddleware returns net/http middleware that extracts trace context, starts a span per request,
// and injects context. Use with your HTTP mux.
func HTTPMiddleware(next http.Handler) http.Handler {
	return tracing.Middleware(next)
}

// GrpcServerInterceptor returns a gRPC unary server interceptor for trace propagation.
func GrpcServerInterceptor() grpc.UnaryServerInterceptor {
	return tracing.UnaryServerInterceptor()
}

// GrpcClientInterceptor returns a gRPC unary client interceptor for trace propagation.
func GrpcClientInterceptor() grpc.UnaryClientInterceptor {
	return tracing.UnaryClientInterceptor()
}

// InjectHTTPRequest injects trace context into outgoing HTTP request headers.
func InjectHTTPRequest(ctx context.Context, req *http.Request) {
	tracing.InjectIntoRequest(ctx, req)
}

// InjectKafkaHeaders returns trace context as map[string]string for Kafka message headers.
func InjectKafkaHeaders(ctx context.Context) map[string]string {
	return tracing.InjectKafkaHeaders(ctx)
}

// ExtractKafkaContext extracts trace context from Kafka message headers.
func ExtractKafkaContext(ctx context.Context, headers map[string]string) context.Context {
	return tracing.ExtractKafkaContext(ctx, headers)
}

// StartKafkaConsumerSpan starts a span for consuming a Kafka message. Call after ExtractKafkaContext.
func StartKafkaConsumerSpan(ctx context.Context, topic, partition string, offset int64) (context.Context, trace.Span) {
	return tracing.StartKafkaConsumerSpan(ctx, topic, partition, offset)
}

// StartKafkaProducerSpan starts a span for producing a Kafka message. Use returned ctx with InjectKafkaHeaders.
func StartKafkaProducerSpan(ctx context.Context, topic string) (context.Context, trace.Span) {
	return tracing.StartKafkaProducerSpan(ctx, topic)
}

// NewCounter creates a basic counter for metrics. Name and description identify the metric.
func NewCounter(name, description string) (*metrics.Counter, error) {
	return metrics.NewCounter(name, description)
}
