package tracing

import (
	"context"

	"github.com/MH-Cognition/mhc-infra-observability/propagation"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// InjectKafkaHeaders injects trace context into a map suitable for Kafka message headers.
// Returns map[string]string that callers merge into their Kafka producer record.
func InjectKafkaHeaders(ctx context.Context) map[string]string {
	return propagation.InjectKafka(ctx)
}

// ExtractKafkaContext extracts trace context from Kafka message headers into ctx.
// Use when consuming a Kafka message to continue the trace.
func ExtractKafkaContext(ctx context.Context, headers map[string]string) context.Context {
	return propagation.ExtractKafka(ctx, headers)
}

// StartKafkaConsumerSpan starts a span for a Kafka message consumer.
// Call after ExtractKafkaContext to create a child span for processing.
func StartKafkaConsumerSpan(ctx context.Context, topic, partition string, offset int64) (context.Context, trace.Span) {
	tracer := otel.Tracer(tracerName)
	return tracer.Start(ctx, "kafka.consume",
		trace.WithSpanKind(trace.SpanKindConsumer),
		trace.WithAttributes(
			attribute.String("messaging.system", "kafka"),
			attribute.String("messaging.destination", topic),
			attribute.String("messaging.kafka.destination.partition", partition),
			attribute.Int64("messaging.kafka.message.offset", offset),
		),
	)
}

// StartKafkaProducerSpan starts a span for producing a Kafka message.
// Use the returned context when calling InjectKafkaHeaders.
func StartKafkaProducerSpan(ctx context.Context, topic string) (context.Context, trace.Span) {
	tracer := otel.Tracer(tracerName)
	return tracer.Start(ctx, "kafka.produce",
		trace.WithSpanKind(trace.SpanKindProducer),
		trace.WithAttributes(
			attribute.String("messaging.system", "kafka"),
			attribute.String("messaging.destination", topic),
		),
	)
}
