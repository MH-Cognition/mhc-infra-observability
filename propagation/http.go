// Package propagation provides trace context propagation helpers for HTTP, gRPC, and Kafka.
// These are used by tracing middleware/interceptors to maintain distributed trace continuity.
package propagation

import (
	"context"
	"net/http"

	"go.opentelemetry.io/otel/propagation"
)

// HTTPHeaderCarrier adapts http.Header to the propagation.TextMapCarrier interface.
type HTTPHeaderCarrier struct {
	Header http.Header
}

// Get returns the value for the given key.
func (c HTTPHeaderCarrier) Get(key string) string {
	return c.Header.Get(key)
}

// Set sets the key-value pair.
func (c HTTPHeaderCarrier) Set(key, value string) {
	c.Header.Set(key, value)
}

// Keys returns all keys in the carrier.
func (c HTTPHeaderCarrier) Keys() []string {
	keys := make([]string, 0, len(c.Header))
	for k := range c.Header {
		keys = append(keys, k)
	}
	return keys
}

// ExtractHTTP extracts trace context from HTTP request headers into ctx.
// Uses the global propagator. Call before starting a span for incoming requests.
func ExtractHTTP(ctx context.Context, header http.Header) context.Context {
	carrier := HTTPHeaderCarrier{Header: header}
	return propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	).Extract(ctx, carrier)
}

// InjectHTTP injects trace context from ctx into HTTP request headers.
// Call before sending outgoing HTTP requests.
func InjectHTTP(ctx context.Context, header http.Header) {
	carrier := HTTPHeaderCarrier{Header: header}
	propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	).Inject(ctx, carrier)
}
