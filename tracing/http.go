package tracing

import (
	"context"
	"net/http"
	"strconv"

	"github.com/MH-Cognition/mhc-infra-observability/propagation"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

const tracerName = "mhc-infra-observability"
const spanNamePrefix = "http."

// Middleware returns an http.Handler that extracts trace context from request headers,
// starts a root span per request, injects context into the request, and records span status.
func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := propagation.ExtractHTTP(r.Context(), r.Header)

		tracer := otel.Tracer(tracerName)
		spanName := spanNamePrefix + r.Method + " " + r.URL.Path
		if spanName == spanNamePrefix+" " {
			spanName = spanNamePrefix + r.URL.Path
		}
		if spanName == spanNamePrefix {
			spanName = spanNamePrefix + "/"
		}

		ctx, span := tracer.Start(ctx, spanName,
			trace.WithSpanKind(trace.SpanKindServer),
			trace.WithAttributes(
				attribute.String("http.method", r.Method),
				attribute.String("http.url", r.URL.String()),
				attribute.String("http.scheme", r.URL.Scheme),
				attribute.String("http.host", r.Host),
				attribute.String("http.target", r.URL.Path),
				attribute.String("http.user_agent", r.UserAgent()),
			),
		)
		defer span.End()

		r = r.WithContext(ctx)

		wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(wrapped, r)

		span.SetAttributes(attribute.Int("http.status_code", wrapped.statusCode))
		if wrapped.statusCode >= 400 {
			span.SetStatus(codes.Error, "HTTP "+strconv.Itoa(wrapped.statusCode))
		}
	})
}

// responseWriter wraps http.ResponseWriter to capture status code.
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// InjectIntoRequest injects trace context into outgoing HTTP request headers.
// Call before sending the request.
func InjectIntoRequest(ctx context.Context, req *http.Request) {
	propagation.InjectHTTP(ctx, req.Header)
}
