// Package tracing provides OpenTelemetry tracing initialization and middleware
// for HTTP, gRPC, and Kafka. Used internally by the observability facade.
// Resource is created once in observability.NewResource and passed here.
//
// Manual OTEL only: we never call otel.Tracer() or otel.GetTracerProvider() before
// SetTracerProvider. The tracer is obtained from our TracerProvider after it is set.
package tracing

import (
	"context"
	"fmt"
	"sync"

	"github.com/MH-Cognition/mhc-infra-observability/config"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	mu            sync.RWMutex
	defaultTracer trace.Tracer // set in Init; never use otel.Tracer() before Init
)

// Tracer returns the tracer for this library. Safe to call only after Init has been called.
// Before Init, returns a noop tracer so global otel APIs are never touched.
func Tracer() trace.Tracer {
	mu.RLock()
	t := defaultTracer
	mu.RUnlock()
	if t != nil {
		return t
	}
	return trace.NewNoopTracerProvider().Tracer(tracerName)
}

// Init initializes the OpenTelemetry TracerProvider with OTLP gRPC exporter.
// Order is strict: 1) create provider with resource 2) SetTracerProvider 3) then obtain tracer.
// Uses the single Resource created by observability.NewResource (do not create resource here).
// Registers the global TracerProvider and Propagator. Returns a shutdown function.
func Init(ctx context.Context, res *resource.Resource, cfg *config.Config) (func(context.Context) error, error) {
	conn, err := grpc.DialContext(ctx, cfg.OtelEndpoint,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("create OTLP gRPC connection: %w", err)
	}

	exporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
	if err != nil {
		_ = conn.Close()
		return nil, fmt.Errorf("create OTLP trace exporter: %w", err)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	mu.Lock()
	defaultTracer = tp.Tracer(tracerName)
	mu.Unlock()

	shutdown := func(ctx context.Context) error {
		if err := tp.Shutdown(ctx); err != nil {
			return fmt.Errorf("shutdown tracer provider: %w", err)
		}
		return nil
	}

	return shutdown, nil
}
