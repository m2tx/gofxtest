package telemetry

import (
	"context"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	otelprom "go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.43.0"

	"go.opentelemetry.io/contrib/instrumentation/runtime"
)

type MetricsConfig struct {
	ServiceName string `env:"NAME" default:"gofxtest"`
}

// MeterProvider is the application meter provider. It is exported to
// Prometheus through the registry returned by NewRegistry.
type MeterProvider interface {
	metric.MeterProvider
	Shutdown(ctx context.Context) error
}

// NewRegistry creates the Prometheus registry served on /metrics.
// A dedicated registry avoids leaking metrics registered globally by libraries.
func NewRegistry() *prometheus.Registry {
	registry := prometheus.NewRegistry()
	registry.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))

	return registry
}

// NewMeterProvider creates an OpenTelemetry meter provider that exports to the
// given Prometheus registry and starts the Go runtime metrics instrumentation.
func NewMeterProvider(config MetricsConfig, registry *prometheus.Registry) (MeterProvider, error) {
	exporter, err := otelprom.New(otelprom.WithRegisterer(registry))
	if err != nil {
		return nil, err
	}

	res, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(semconv.SchemaURL, semconv.ServiceName(config.ServiceName)),
	)
	if err != nil {
		return nil, err
	}

	provider := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(exporter),
		sdkmetric.WithResource(res),
	)

	err = runtime.Start(runtime.WithMeterProvider(provider))
	if err != nil {
		return nil, err
	}

	return provider, nil
}
