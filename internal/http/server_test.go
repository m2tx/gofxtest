package http_test

import (
	"context"
	"testing"
	"time"

	"github.com/m2tx/gofxtest/internal/env"
	"github.com/m2tx/gofxtest/internal/http"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

var (
	ctx = context.Background()
)

func TestHttpServer(t *testing.T) {
	httpConfig, err := env.New[http.HttpConfig]()
	assert.NoError(t, err)

	srv := http.NewServer(httpConfig, nil, zap.NewNop())

	assert.NotNil(t, srv)

	err = srv.Start()
	assert.NoError(t, err)

	time.Sleep(100 * time.Millisecond)

	err = srv.Shutdown(ctx)
	assert.NoError(t, err)
}
