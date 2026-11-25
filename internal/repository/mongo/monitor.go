package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/event"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

const (
	meterName = "github.com/m2tx/gofxtest/internal/repository/mongo"

	attrOperation = "db.operation.name"
	attrNamespace = "db.namespace"
	attrStatus    = "status"

	statusOK    = "ok"
	statusError = "error"
)

// newCommandMonitor records the duration of every MongoDB command.
func newCommandMonitor(meterProvider metric.MeterProvider) (*event.CommandMonitor, error) {
	meter := meterProvider.Meter(meterName)

	duration, err := meter.Float64Histogram(
		"db.client.operation.duration",
		metric.WithDescription("Duration of MongoDB commands."),
		metric.WithUnit("s"),
	)
	if err != nil {
		return nil, err
	}

	record := func(ctx context.Context, evt event.CommandFinishedEvent, status string) {
		duration.Record(ctx, evt.Duration.Seconds(), metric.WithAttributes(
			attribute.String(attrOperation, evt.CommandName),
			attribute.String(attrNamespace, evt.DatabaseName),
			attribute.String(attrStatus, status),
		))
	}

	return &event.CommandMonitor{
		Succeeded: func(ctx context.Context, evt *event.CommandSucceededEvent) {
			record(ctx, evt.CommandFinishedEvent, statusOK)
		},
		Failed: func(ctx context.Context, evt *event.CommandFailedEvent) {
			record(ctx, evt.CommandFinishedEvent, statusError)
		},
	}, nil
}
