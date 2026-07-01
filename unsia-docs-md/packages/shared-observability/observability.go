package sharedobservability

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

var (
	Logger      zerolog.Logger
	metricsOnce sync.Once
	tracerOnce  sync.Once
	tracer      trace.Tracer

	HttpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests processed.",
		},
		[]string{"method", "path", "status"},
	)
	HttpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Latency of HTTP requests in seconds.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path", "status"},
	)
)

func InitLogger(serviceName string) {
	zerolog.TimeFieldFormat = time.RFC3339
	Logger = zerolog.New(os.Stdout).With().
		Timestamp().
		Str("service", serviceName).
		Logger()
}

func Ctx(ctx context.Context) zerolog.Logger {
	l := Logger
	if ctx == nil {
		return l
	}

	correlationID := ""
	traceID := ""

	if cid, ok := ctx.Value("correlation_id").(string); ok {
		correlationID = cid
	} else if cid, ok := ctx.Value("x-correlation-id").(string); ok {
		correlationID = cid
	}

	// Try extracting trace ID from OpenTelemetry span
	span := trace.SpanFromContext(ctx)
	if span.SpanContext().IsValid() {
		traceID = span.SpanContext().TraceID().String()
	} else if tid, ok := ctx.Value("trace_id").(string); ok {
		traceID = tid
	}

	// For gin.Context compatibility
	if gc, ok := ctx.(interface{ GetString(string) string }); ok {
		if correlationID == "" {
			if cid := gc.GetString("x-correlation-id"); cid != "" {
				correlationID = cid
			} else if cid := gc.GetString("correlation_id"); cid != "" {
				correlationID = cid
			}
		}
		if traceID == "" {
			if tid := gc.GetString("trace_id"); tid != "" {
				traceID = tid
			}
		}
	}

	if correlationID != "" {
		l = l.With().Str("correlation_id", correlationID).Logger()
	}
	if traceID != "" {
		l = l.With().Str("trace_id", traceID).Logger()
	}

	return l
}

func InitMetrics() {
	metricsOnce.Do(func() {
		prometheus.MustRegister(HttpRequestsTotal)
		prometheus.MustRegister(HttpRequestDuration)
	})
}

// MetricsHandler returns a Gin handler for Prometheus /metrics endpoint
func MetricsHandler() gin.HandlerFunc {
	InitMetrics()
	h := promhttp.Handler()
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

// InitTracer initializes OpenTelemetry trace provider
func InitTracer(serviceName string) (*sdktrace.TracerProvider, error) {
	var tp *sdktrace.TracerProvider
	var err error

	tracerOnce.Do(func() {
		tp = sdktrace.NewTracerProvider(
			sdktrace.WithSampler(sdktrace.AlwaysSample()),
		)
		otel.SetTracerProvider(tp)
		otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
		tracer = tp.Tracer(serviceName)
	})

	return tp, err
}

func GetTracer() trace.Tracer {
	return tracer
}

func generateRandomID(length int) string {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "rand-failed"
	}
	return hex.EncodeToString(bytes)
}

// CorrelationIDMiddleware generates and propagates X-Correlation-Id
func CorrelationIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		correlationID := c.GetHeader("X-Correlation-Id")
		if correlationID == "" {
			correlationID = generateRandomID(16)
		}

		c.Set("x-correlation-id", correlationID)
		c.Set("correlation_id", correlationID)
		c.Header("X-Correlation-Id", correlationID)

		// Create standard context containing key-values
		ctx := context.WithValue(c.Request.Context(), "x-correlation-id", correlationID)
		ctx = context.WithValue(ctx, "correlation_id", correlationID)

		// Set OpenTelemetry trace context
		ctx = otel.GetTextMapPropagator().Extract(ctx, propagation.HeaderCarrier(c.Request.Header))
		span := trace.SpanFromContext(ctx)
		if span.SpanContext().IsValid() {
			c.Set("trace_id", span.SpanContext().TraceID().String())
			ctx = context.WithValue(ctx, "trace_id", span.SpanContext().TraceID().String())
		}

		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

// RequestLoggerMiddleware logs HTTP request execution
func RequestLoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		c.Next()

		duration := time.Since(start)
		if raw != "" {
			path = path + "?" + raw
		}

		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()
		userAgent := c.Request.UserAgent()
		method := c.Request.Method

		logger := Ctx(c.Request.Context())

		event := logger.Info()
		if statusCode >= 500 {
			event = logger.Error()
		} else if statusCode >= 400 {
			event = logger.Warn()
		}

		event.
			Str("method", method).
			Str("path", path).
			Int("status", statusCode).
			Int64("duration_ms", duration.Milliseconds()).
			Str("ip", clientIP).
			Str("user_agent", userAgent).
			Msg("request completed")
	}
}

// MetricsMiddleware records HTTP metrics
func MetricsMiddleware() gin.HandlerFunc {
	InitMetrics()
	return func(c *gin.Context) {
		start := time.Now()
		path := c.FullPath()
		if path == "" {
			path = "unknown"
		}
		method := c.Request.Method

		c.Next()

		status := fmt.Sprintf("%d", c.Writer.Status())
		duration := time.Since(start).Seconds()

		HttpRequestsTotal.WithLabelValues(method, path, status).Inc()
		HttpRequestDuration.WithLabelValues(method, path, status).Observe(duration)
	}
}

// CORSMiddleware handles CORS requests for all microservices
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		if origin != "" {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Credentials", "true")
			c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS, HEAD")
			c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization, X-Requested-With, X-Correlation-ID")
			c.Header("Access-Control-Max-Age", "86400")
		}

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

