package propagation

import (
	"context"

	"go.opentelemetry.io/otel/propagation"
)

// KafkaHeaderCarrier adapts map[string]string (Kafka-style headers) to propagation.TextMapCarrier.
// Kafka headers are typically represented as key-value strings for trace propagation.
type KafkaHeaderCarrier struct {
	Headers map[string]string
}

// Get returns the value for the given key.
func (c KafkaHeaderCarrier) Get(key string) string {
	return c.Headers[key]
}

// Set sets the key-value pair.
func (c KafkaHeaderCarrier) Set(key, value string) {
	if c.Headers == nil {
		c.Headers = make(map[string]string)
	}
	c.Headers[key] = value
}

// Keys returns all keys in the carrier.
func (c KafkaHeaderCarrier) Keys() []string {
	keys := make([]string, 0, len(c.Headers))
	for k := range c.Headers {
		keys = append(keys, k)
	}
	return keys
}

// ExtractKafka extracts trace context from Kafka message headers into ctx.
// headers: map of header key to value (e.g., from sarama.RecordHeader or similar).
func ExtractKafka(ctx context.Context, headers map[string]string) context.Context {
	if headers == nil {
		return ctx
	}
	carrier := KafkaHeaderCarrier{Headers: headers}
	return propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	).Extract(ctx, carrier)
}

// InjectKafka injects trace context from ctx into a new map suitable for Kafka headers.
// Returns map[string]string that callers can add to Kafka message headers.
func InjectKafka(ctx context.Context) map[string]string {
	headers := make(map[string]string)
	carrier := KafkaHeaderCarrier{Headers: headers}
	propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	).Inject(ctx, carrier)
	return headers
}
