// Package metrics provides minimal OpenTelemetry metrics helpers.
// Kept minimal per design; services can extend with custom meters.
package metrics

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
)

// Counter is a minimal counter helper.
type Counter struct {
	counter metric.Int64Counter
}

// NewCounter creates a counter with the given name and optional description.
func NewCounter(name, description string) (*Counter, error) {
	meter := otel.Meter("mhc-infra-observability")
	c, err := meter.Int64Counter(name,
		metric.WithDescription(description),
	)
	if err != nil {
		return nil, err
	}
	return &Counter{counter: c}, nil
}

// Add increments the counter by n. Optional attributes can be passed via metric.WithAttributes.
func (c *Counter) Add(ctx context.Context, n int64, opts ...metric.AddOption) {
	c.counter.Add(ctx, n, opts...)
}

// Increment adds 1 to the counter.
func (c *Counter) Increment(ctx context.Context, opts ...metric.AddOption) {
	c.Add(ctx, 1, opts...)
}
