package http

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
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

func NewHandler(routeHandlers []RouteHandler, logger *zap.Logger) Handler {
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = io.Discard

	gHandler := gin.Default()

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
