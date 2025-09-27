package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type RouteHandler interface {
	Handler(g *gin.Engine)
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

	gHandler := gin.Default()

	for _, route := range routeHandlers {
		route.Handler(gHandler)
	}

	return &handler{
		Handler: gHandler,
		logger:  logger,
	}
}
