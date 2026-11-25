package telemetry_test

import (
	"testing"

	"github.com/m2tx/gofxtest/internal/telemetry"
	"github.com/stretchr/testify/assert"
)

func TestNewMeterProvider(t *testing.T) {
	registry := telemetry.NewRegistry()

	provider, err := telemetry.NewMeterProvider(telemetry.MetricsConfig{ServiceName: "test"}, registry)
	assert.NoError(t, err)
	defer func() {
		assert.NoError(t, provider.Shutdown(t.Context()))
	}()

	counter, err := provider.Meter("test").Int64Counter("test.counter")
	assert.NoError(t, err)
	counter.Add(t.Context(), 3)

	families, err := registry.Gather()
	assert.NoError(t, err)

	names := map[string]bool{}
	for _, f := range families {
		names[f.GetName()] = true
	}

	assert.True(t, names["test_counter_total"], "custom metric exported")
	assert.True(t, names["target_info"], "resource exported as target_info")
	assert.True(t, names["process_cpu_seconds_total"], "process collector registered")
	assert.True(t, names["go_goroutine_count"], "runtime metrics exported")
}
