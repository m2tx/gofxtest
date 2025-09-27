package http

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"go.uber.org/zap"
)

type HttpServer interface {
	Start() error
	Shutdown(ctx context.Context) error
}

type HttpConfig struct {
	Port            int           `env:"HTTP_PORT" default:"8080"`
	ShutdownTimeout time.Duration `env:"HTTP_SHUTDOWN_TIMEOUT" default:"5s"`
	ReadTimeout     time.Duration `env:"HTTP_READ_TIMEOUT" default:"10s"`
	WriteTimeout    time.Duration `env:"HTTP_WRITE_TIMEOUT" default:"10s"`
}

type httpServer struct {
	config HttpConfig
	server *http.Server
	logger *zap.Logger
}

func NewServer(config HttpConfig, handler Handler, logger *zap.Logger) HttpServer {
	return &httpServer{
		config: config,
		logger: logger,
		server: &http.Server{
			Addr:         fmt.Sprintf(":%d", config.Port),
			Handler:      handler,
			ReadTimeout:  config.ReadTimeout,
			WriteTimeout: config.WriteTimeout,
		},
	}
}

func (h *httpServer) Start() error {
	ln, err := net.Listen("tcp", h.server.Addr)
	if err != nil {
		h.logger.Error("HTTP server listen error", zap.Error(err))
		return err
	}

	h.logger.Info("HTTP server started", zap.String("addr", h.server.Addr))

	go h.server.Serve(ln)

	return nil
}

func (h *httpServer) Shutdown(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, h.config.ShutdownTimeout)
	defer cancel()
	err := h.server.Shutdown(ctx)
	if err != nil {
		h.logger.Error("HTTP server shutdown error", zap.Error(err))
		return err
	}

	h.logger.Info("HTTP server stopped")

	return nil
}
