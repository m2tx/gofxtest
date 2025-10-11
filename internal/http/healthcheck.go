package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type healthcheckRoute struct {
	RouteHandler
}

func NewHealthcheckRoute() RouteHandler {
	return &healthcheckRoute{}
}

func (r *healthcheckRoute) Register(e *gin.Engine) {
	e.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "OK"})
	})
}
