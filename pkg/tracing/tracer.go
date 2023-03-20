package tracing

import (
	"context"
	"os"
	"strconv"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

func Init() func() error {
	if enable, _ := strconv.ParseBool(os.Getenv("TRACING_ENABLE")); !enable {
		return func() error {
			return nil
		}
	}

	env := os.Getenv("ENV")
	if env == "" {
		panic("missing ENV")
	}

	jaegerEndpoint := os.Getenv("TRACING_JAEGER_ENDPOINT")
	if jaegerEndpoint == "" {
		panic("missing TRACING_JAEGER_ENDPOINT")
	}

	serviceName := os.Getenv("TRACING_SERVICE_NAME")
	if serviceName == "" {
		panic("missing TRACING_SERVICE_NAME")
	}

	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(jaegerEndpoint)))
	if err != nil {
		panic(err)
	}

	tp := tracesdk.NewTracerProvider(
		// Always be sure to batch in production.
		tracesdk.WithBatcher(exp),
		// Record information about this application in a Resource.
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(serviceName),
			attribute.String("environment", env),
		)),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	return func() error {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return tp.Shutdown(ctx)
	}
}
