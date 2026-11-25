package http

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel/metric"
	"go.uber.org/zap"
)

type RouteHandler interface {
	Register(g *gin.Engine)
}

type Handler interface {
	http.Handler
}

type handler struct {
	http.Handler
	logger *zap.Logger
}

func NewHandler(routeHandlers []RouteHandler, logger *zap.Logger, meterProvider metric.MeterProvider, config HttpConfig) Handler {
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = io.Discard

	gHandler := gin.Default()

	gHandler.Use(otelgin.Middleware(config.ServiceName,
		otelgin.WithMeterProvider(meterProvider),
		otelgin.WithGinFilter(func(c *gin.Context) bool {
			return c.FullPath() != "/metrics"
		}),
	))

	gHandler.Use(func(ctx *gin.Context) {
		logger.Debug("access", zap.String("uri", ctx.Request.RequestURI), zap.String("method", ctx.Request.Method))
		ctx.Next()
	})

	for _, route := range routeHandlers {
		route.Register(gHandler)
	}

	return &handler{
		Handler: gHandler,
		logger:  logger,
	}
}
