package queue

import (
	"go.opentelemetry.io/otel/metric"
)

const (
	meterName = "github.com/m2tx/gofxtest/internal/queue"

	attrTopic  = "messaging.destination.name"
	attrStatus = "status"

	statusOK    = "ok"
	statusError = "error"
)

type queueMetrics struct {
	published       metric.Int64Counter
	consumeDuration metric.Float64Histogram
}

func newQueueMetrics(meterProvider metric.MeterProvider) (*queueMetrics, error) {
	meter := meterProvider.Meter(meterName)

	published, err := meter.Int64Counter(
		"messaging.client.published.messages",
		metric.WithDescription("Number of messages published."),
		metric.WithUnit("{message}"),
	)
	if err != nil {
		return nil, err
	}

	consumeDuration, err := meter.Float64Histogram(
		"messaging.process.duration",
		metric.WithDescription("Duration of the handler processing a consumed message."),
		metric.WithUnit("s"),
	)
	if err != nil {
		return nil, err
	}

	return &queueMetrics{
		published:       published,
		consumeDuration: consumeDuration,
	}, nil
}

func status(err error) string {
	if err != nil {
		return statusError
	}

	return statusOK
}
