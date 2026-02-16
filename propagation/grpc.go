package propagation

import (
	"context"

	"go.opentelemetry.io/otel/propagation"
	"google.golang.org/grpc/metadata"
)

// grpcMetadataCarrier adapts gRPC metadata to propagation.TextMapCarrier.
type grpcMetadataCarrier struct {
	md metadata.MD
}

func (c grpcMetadataCarrier) Get(key string) string {
	vals := c.md.Get(key)
	if len(vals) == 0 {
		return ""
	}
	return vals[0]
}

func (c grpcMetadataCarrier) Set(key, value string) {
	c.md.Set(key, value)
}

func (c grpcMetadataCarrier) Keys() []string {
	keys := make([]string, 0, len(c.md))
	for k := range c.md {
		keys = append(keys, k)
	}
	return keys
}

// ExtractGrpc extracts trace context from gRPC incoming metadata into ctx.
func ExtractGrpc(ctx context.Context, md metadata.MD) context.Context {
	carrier := grpcMetadataCarrier{md: md}
	return propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	).Extract(ctx, carrier)
}

// InjectGrpc injects trace context from ctx into gRPC outgoing metadata.
// Returns new metadata with trace headers. Merge with existing metadata if needed.
func InjectGrpc(ctx context.Context, md metadata.MD) metadata.MD {
	if md == nil {
		md = metadata.New(nil)
	}
	carrier := grpcMetadataCarrier{md: md}
	propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	).Inject(ctx, carrier)
	return md
}
