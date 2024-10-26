package routes

import (
	"app/app/controller"
	"app/app/middleware"

	"github.com/gin-gonic/gin"
)

func Role(router *gin.RouterGroup) {
	// Get the *bun.DB instance from config
	ctl := controller.New() // Pass the *bun.DB to the controller
	md := middleware.AuthMiddleware()
	role := router.Group("")
	{

		role.POST("/create", md, ctl.RoleCtl.Create)
		role.GET("/list", md,ctl.RoleCtl.List)
		role.GET("/:id", md,ctl.RoleCtl.Get)
		role.PATCH("/:id", md,ctl.RoleCtl.Update)
		role.DELETE("/:id", md,ctl.RoleCtl.Delete)
	}
}
