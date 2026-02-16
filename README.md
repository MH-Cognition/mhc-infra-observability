# mhc-infra-observability

A **centralized infrastructure library** for observability in Go microservices. This is **not a runnable service**—it is a shared library consumed by multiple microservices that follow Clean Architecture and use no frameworks.

## What This Is

- **Observability infrastructure**: OpenTelemetry tracing, structured logging, and basic metrics helpers
- **Context propagation**: Trace context across HTTP, gRPC, and Kafka
- **Production-ready**: Designed for microservices behind NGINX using HTTP, gRPC, and Kafka

## What This Is NOT

- **Not a service**: No `cmd/`, no `main.go`, nothing to run
- **Not domain logic**: No `domain/`, `usecase/`, `transport/`, `repository/`—no business logic
- **Not HTTP/gRPC handlers**: No routes, handlers, or response mapping
- **Not error semantics**: Does not define domain errors or map errors to HTTP status codes

## Integration

### 1. Add dependency

```go
// In your service's go.mod
require github.com/MH-Cognition/mhc-infra-observability v0.1.2
```

### 2. Initialize in main.go

```go
package main

import (
    "context"
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"

    "mhc-infra-observability/config"
    "mhc-infra-observability/observability"
)

func main() {
    ctx := context.Background()
    cfg := config.Load()

    // Create exactly ONE OpenTelemetry Resource; pass it to Init.
    res, err := observability.NewResource(ctx, cfg)
    if err != nil {
        log.Fatalf("observability resource: %v", err)
    }
    shutdown, err := observability.Init(ctx, res, cfg)
    if err != nil {
        log.Fatalf("observability init: %v", err)
    }
    defer func() {
        if err := shutdown(context.Background()); err != nil {
            log.Printf("shutdown error: %v", err)
        }
    }()

    mux := http.NewServeMux()
    mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
    })

    handler := observability.HTTPMiddleware(mux)

    srv := &http.Server{Addr: ":8080", Handler: handler}
    go srv.ListenAndServe()

    sigCh := make(chan os.Signal, 1)
    signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
    <-sigCh
    srv.Shutdown(context.Background())
}
```

### 3. HTTP middleware usage

Wrap your HTTP mux with `observability.HTTPMiddleware` so every request gets a span and trace context propagation:

```go
handler := observability.HTTPMiddleware(mux)
```

For outgoing HTTP requests, inject trace context before sending:

```go
req, _ := http.NewRequestWithContext(ctx, "GET", "https://downstream/api", nil)
observability.InjectHTTPRequest(ctx, req)
resp, err := http.DefaultClient.Do(req)
```

### 4. Span usage in use-case layer

```go
// In your use-case (called from transport with traced context)
func (uc *OrderUseCase) CreateOrder(ctx context.Context, req CreateOrderRequest) error {
    ctx, span := observability.StartSpan(ctx, "CreateOrder")
    defer span.End()

    // ... business logic ...

    if err != nil {
        observability.HandleError(ctx, err)
        return err
    }
    return nil
}
```

### 5. Logging

```go
logger := observability.Logger(ctx)
logger.Info(ctx, "order created", "order_id", orderID)
logger.Error(ctx, "validation failed", "field", "amount")
```

### 6. Metrics

```go
counter, err := observability.NewCounter("requests_total", "Total number of requests")
if err == nil {
    counter.Increment(ctx)
}
```

### 7. gRPC

```go
// Server
grpc.NewServer(grpc.UnaryInterceptor(observability.GrpcServerInterceptor()))

// Client
conn, err := grpc.Dial(addr, grpc.WithUnaryInterceptor(observability.GrpcClientInterceptor()))
```

### 8. Kafka

```go
// Producer: start span, inject headers into message
ctx, span := observability.StartSpan(ctx, "publish.order.created")
defer span.End()
headers := observability.InjectKafkaHeaders(ctx)
// Merge headers into your Kafka producer record

// Consumer: extract context, start span
ctx := observability.ExtractKafkaContext(ctx, msg.HeadersAsMap())
ctx, span := observability.StartSpan(ctx, "consume.order.created")
defer span.End()
```

## Environment variables

| Variable | Description | Default |
|----------|-------------|---------|
| `OTEL_SERVICE_NAME` | Service name in traces | `unknown-service` |
| `OTEL_SERVICE_VERSION` | Service version (optional) | — |
| `OTEL_ENVIRONMENT` | Deployment environment | `development` |
| `OTEL_EXPORTER_OTLP_ENDPOINT` | OTLP collector endpoint | `localhost:4317` |
| `LOG_LEVEL` | Log level (debug, info, error) | `info` |

## Why domain code must not import this directly

**Clean Architecture** keeps infrastructure at the outer layers. Use-cases and domain should depend on interfaces, not concrete implementations.

- **Dependency rule**: Inner layers (domain, use-cases) should not import outer layers (infrastructure).
- **Testing**: Use-cases are easier to unit-test when they receive a logger interface, not a concrete observability type.
- **Adapter pattern**: The transport layer (HTTP handlers, gRPC handlers) should be the only place that imports `mhc-infra-observability`. It extracts context, starts spans, and passes a clean `context.Context` to use-cases. Use-cases receive `context.Context` and can call `observability.StartSpan(ctx, ...)` and `observability.HandleError(ctx, err)` if your architecture allows a thin infra interface there, or the transport can wrap use-case calls with spans and error handling.

**Practical guidance**: Prefer importing `mhc-infra-observability` only from `main.go` and transport/adapters. If use-cases need spans, pass `context.Context` and have the transport start the span around the use-case call, or define a minimal `SpanStarter` interface in your use-case package that the transport implements.

## Package layout

```
mhc-infra-observability/
├── config/         # Config from env
├── observability/  # Public facade (import this); includes NewResource (single OTEL Resource)
├── tracing/        # OTel tracing + HTTP/gRPC/Kafka middleware
├── logging/        # Structured trace-aware logger
├── metrics/        # Basic counter helper
└── propagation/    # Trace context propagation
```

## Resource and schema (no conflicts)

Exactly **one** OpenTelemetry Resource is created per process via `observability.NewResource(ctx, cfg)`. It uses a single schema version (semconv v1.37.0). The same Resource is passed to `observability.Init(ctx, res, cfg)`. Do not create resources in tracing, metrics, or logging; do not use `resource.Default()` or `resource.Merge()` elsewhere.

## Manual OTEL only (v0.1.2+)

This library uses **strict manual initialization** only. `go.opentelemetry.io/auto/sdk` is **not** in the dependency graph. The service must call `NewResource` and then `Init(ctx, res, cfg)` before any HTTP/gRPC handlers or `StartSpan` run. Tracers are obtained from the TracerProvider set in `Init`; the global `otel.Tracer()` is never used before the provider is set.
