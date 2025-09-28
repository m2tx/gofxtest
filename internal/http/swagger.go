package http

import (
	"github.com/gin-gonic/gin"
	_ "github.com/m2tx/gofxtest/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type swaggerRoute struct {
	RouteHandler
}

func NewSwaggerRoute() RouteHandler {
	return &swaggerRoute{}
}

func (r *swaggerRoute) Handler(e *gin.Engine) {
	e.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
