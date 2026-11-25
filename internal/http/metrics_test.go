package http_test

import (
	"io"
	nethttp "net/http"
	"net/http/httptest"
	"testing"

	"github.com/m2tx/gofxtest/internal/http"
	"github.com/m2tx/gofxtest/internal/telemetry"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestMetricsRoute(t *testing.T) {
	registry := telemetry.NewRegistry()

	provider, err := telemetry.NewMeterProvider(telemetry.MetricsConfig{ServiceName: "test"}, registry)
	assert.NoError(t, err)
	defer func() {
		assert.NoError(t, provider.Shutdown(t.Context()))
	}()

	handler := http.NewHandler(
		[]http.RouteHandler{http.NewHealthcheckRoute(), http.NewMetricsRoute(registry)},
		zap.NewNop(),
		provider,
		http.HttpConfig{ServiceName: "test"},
	)

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(nethttp.MethodGet, "/health", nil))
	assert.Equal(t, nethttp.StatusOK, rec.Code)

	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(nethttp.MethodGet, "/metrics", nil))
	assert.Equal(t, nethttp.StatusOK, rec.Code)

	body, err := io.ReadAll(rec.Body)
	assert.NoError(t, err)

	assert.Contains(t, string(body), "http_server_request_duration_seconds")
	assert.Contains(t, string(body), `http_route="/health"`)
	assert.NotContains(t, string(body), `http_route="/metrics"`)
}
