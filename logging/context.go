// Package logging provides a structured, trace-aware logger for microservices.
// Logs include trace_id and span_id when present in context.
package logging

import (
	"context"

	"go.opentelemetry.io/otel/trace"
)

// TraceContext holds trace and span IDs extracted from context for log enrichment.
type TraceContext struct {
	TraceID string
	SpanID  string
}

// FromContext extracts trace context from ctx for log enrichment.
// Returns zero values if no span is present.
func FromContext(ctx context.Context) TraceContext {
	span := trace.SpanFromContext(ctx)
	if !span.SpanContext().IsValid() {
		return TraceContext{}
	}
	sc := span.SpanContext()
	return TraceContext{
		TraceID: sc.TraceID().String(),
		SpanID:  sc.SpanID().String(),
	}
}
