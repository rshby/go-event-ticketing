package tracing

import (
	"context"
	"fmt"

	"github.com/rshby/go-event-ticketing/config"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.39.0"
)

// ConnectOTLPTrace connects to otlp trace
func ConnectOTLPTrace(ctx context.Context) (*trace.TracerProvider, error) {
	// create exporter
	exporter, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithEndpoint(fmt.Sprintf("%s:%s", config.OtlpEndpoint(), config.OtlpPort())),
		otlptracegrpc.WithInsecure())
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	// create resource
	res := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceName(config.OtlpServiceName()))

	// create trace provider
	tp := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithSampler(trace.AlwaysSample()),
		trace.WithResource(res))

	// set to global provider
	otel.SetTracerProvider(tp)

	propagator := propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)

	otel.SetTextMapPropagator(propagator)

	logrus.Infof("success connect to OTLP Trace✅")
	return tp, nil
}
