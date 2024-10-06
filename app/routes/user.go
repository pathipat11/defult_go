package routes

import (
	"app/app/controller"
	"app/app/middleware"

	"github.com/gin-gonic/gin"
)

func User(router *gin.RouterGroup) {
	// Get the *bun.DB instance from config
	ctl := controller.New() // Pass the *bun.DB to the controller
	md := middleware.AuthMiddleware()
	user := router.Group("")
	{

		user.GET("/list", md, ctl.UserCtl.List)
		user.GET("/list/:id", md, ctl.UserCtl.ListSingle)
		// user.PATCH("/edit/:id", md, ctl.UserCtl.Update)
		user.DELETE("/delete/:id", md, ctl.UserCtl.SoftDelete)
		// user.DELETE("/delete-permanent/:id", md, ctl.UserCtl.Delete)
	}
}
