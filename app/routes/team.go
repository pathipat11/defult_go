package routes

import (
	"app/app/controller"
	"app/app/middleware"

	"github.com/gin-gonic/gin"
)

func Team(router *gin.RouterGroup) {
	// Get the *bun.DB instance from config
	ctl := controller.New() // Pass the *bun.DB to the controller
	md := middleware.AuthMiddleware()

	team := router.Group("")
	{
		team.POST("/create", md, ctl.TeamCtl.Create)
		team.POST("/add-member", md, ctl.TeamCtl.AddTeamMember)
		team.POST("/remove-member", md, ctl.TeamCtl.RemoveTeamMember)
		team.GET("/list", md, ctl.TeamCtl.List)
		team.GET("/:id", md, ctl.TeamCtl.Get)
		team.PATCH("/:id", md, ctl.TeamCtl.Update)
		team.DELETE("/:id", md, ctl.TeamCtl.Delete)
	}
}
