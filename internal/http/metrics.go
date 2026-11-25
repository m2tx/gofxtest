package http

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type metricsRoute struct {
	RouteHandler
	gatherer prometheus.Gatherer
}

func NewMetricsRoute(registry *prometheus.Registry) RouteHandler {
	return &metricsRoute{
		gatherer: registry,
	}
}

func (r *metricsRoute) Register(e *gin.Engine) {
	e.GET("/metrics", gin.WrapH(promhttp.HandlerFor(r.gatherer, promhttp.HandlerOpts{})))
}
