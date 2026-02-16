// Package metrics provides minimal OpenTelemetry metrics helpers.
// Kept minimal per design; services can extend with custom meters.
// Does not import go.opentelemetry.io/otel so the auto/sdk chain is never pulled in.
package metrics

import (
	"context"
	"sync"

	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/noop"
)

var (
	mu        sync.RWMutex
	meter     metric.Meter
	meterName = "mhc-infra-observability"
)

// Counter is a minimal counter helper.
type Counter struct {
	counter metric.Int64Counter
}

// SetMeter sets the meter used by NewCounter. Must be called from observability Init
// after the global MeterProvider is set (if any). If never set, NewCounter uses a noop meter.
func SetMeter(m metric.Meter) {
	mu.Lock()
	defer mu.Unlock()
	meter = m
}

func getMeter() metric.Meter {
	mu.RLock()
	m := meter
	mu.RUnlock()
	if m != nil {
		return m
	}
	return noop.NewMeterProvider().Meter(meterName)
}

// NewCounter creates a counter with the given name and optional description.
func NewCounter(name, description string) (*Counter, error) {
	c, err := getMeter().Int64Counter(name,
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
