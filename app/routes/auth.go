package routes

import (
	"app/app/controller"

	"github.com/gin-gonic/gin"
)

func Auth(router *gin.RouterGroup, ctl *controller.Controller) {
	router.POST("/login", ctl.AuthCtl.Login)
}
