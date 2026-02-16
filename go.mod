module github.com/MH-Cognition/mhc-infra-observability

go 1.24.0

require (
	go.opentelemetry.io/otel v1.34.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.24.0
	go.opentelemetry.io/otel/metric v1.34.0
	go.opentelemetry.io/otel/sdk v1.34.0
	go.opentelemetry.io/otel/trace v1.34.0
	google.golang.org/grpc v1.69.0-dev
)

require (
	github.com/cenkalti/backoff/v4 v4.2.1 // indirect
	github.com/go-logr/logr v1.4.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.19.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.24.0 // indirect
	go.opentelemetry.io/proto/otlp v1.1.0 // indirect
	golang.org/x/net v0.41.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
	golang.org/x/text v0.26.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20240903143218-8af14fe29dc1 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250707201910-8d1bb00bc6a7 // indirect
	google.golang.org/protobuf v1.36.6 // indirect
)

replace (
	go.opentelemetry.io/otel => go.opentelemetry.io/otel v1.24.0
	go.opentelemetry.io/otel/metric => go.opentelemetry.io/otel/metric v1.24.0
	go.opentelemetry.io/otel/sdk => go.opentelemetry.io/otel/sdk v1.24.0
	go.opentelemetry.io/otel/trace => go.opentelemetry.io/otel/trace v1.24.0
)

exclude go.opentelemetry.io/auto/sdk v1.1.0
