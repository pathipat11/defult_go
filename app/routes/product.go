package routes

import (
	"app/app/controller"
	"app/app/middleware"

	"github.com/gin-gonic/gin"
)

func Product(router *gin.RouterGroup, ctl *controller.Controller) {
	md := middleware.AuthMiddleware()
	log := middleware.NewLogResponse()
	product := router.Group("", log)
	{
		product.POST("/create", md, ctl.ProductCtl.Create)
		product.GET("/list", md, ctl.ProductCtl.List)
		product.GET("/:id", md, ctl.ProductCtl.Get)
		product.PATCH("/:id", md, ctl.ProductCtl.Update)
		product.DELETE("/:id", md, ctl.ProductCtl.Delete)
	}
}
