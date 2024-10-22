package routes

import (
	"app/app/controller"

	"github.com/gin-gonic/gin"
)

func Role(router *gin.RouterGroup) {
	// Get the *bun.DB instance from config
	ctl := controller.New() // Pass the *bun.DB to the controller

	role := router.Group("")
	{

		role.POST("/create", ctl.RoleCtl.Create)
		role.GET("/list", ctl.RoleCtl.List)
		role.GET("/:id", ctl.RoleCtl.Get)
		role.PATCH("/:id", ctl.RoleCtl.Update)
		role.DELETE("/:id", ctl.RoleCtl.Delete)
	}
}
