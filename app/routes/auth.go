package routes

import (
	"app/app/controller"
	"app/app/middleware"

	"github.com/gin-gonic/gin"
)

func Auth(router *gin.RouterGroup) {
	// Get the *bun.DB instance from config
	ctl := controller.New() // Pass the *bun.DB to the controller
	md := middleware.AuthMiddleware()
	auth := router.Group("")
	{

		auth.POST("/login", ctl.AuthCtl.Login)
		auth.GET("/google/login", ctl.AuthCtl.LoginGoogle)
		auth.GET("/callback-google", ctl.AuthCtl.GoogleCallback)
		auth.POST("/admin/login", ctl.AuthCtl.LoginAdmin)
		auth.POST("/register", ctl.UserCtl.Create)
		auth.GET("/user/detail", md, ctl.AuthCtl.GetUserDetailByToken)
	}
}
