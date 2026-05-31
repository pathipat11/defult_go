package routes

import (
	"app/app/controller"
	"app/app/middleware"

	"github.com/gin-gonic/gin"
)

func User(router *gin.RouterGroup, ctl *controller.Controller) {
	md := middleware.AuthMiddleware()
	user := router.Group("")
	{
		// Public: user registration.
		user.POST("/create", ctl.UserCtl.Create)

		// Protected: everything else requires a valid token.
		user.GET("/list", md, ctl.UserCtl.List)
		user.GET("/:id", md, ctl.UserCtl.Get)
		user.PATCH("/:id", md, ctl.UserCtl.Update)
		user.DELETE("/:id", md, ctl.UserCtl.Delete)
	}
}
