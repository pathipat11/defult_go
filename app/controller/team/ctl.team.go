package team

import (
	"app/app/helper"
	"app/app/request"
	"app/app/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (ctl *Controller) Create(c *gin.Context) {
	// Get the user from the token
	token := c.Request.Header.Get("Authorization")

	user, err := helper.GetUserByToken(c, token)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	var req request.CreateTeam
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	team, err := ctl.Service.Create(c.Request.Context(), req, user)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, team)
}

func (ctl *Controller) AddTeamMember(c *gin.Context) {
	var req request.CreateTeamMember
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	team, err := ctl.Service.AddTeamMember(c.Request.Context(), req)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, team)
}

func (ctl *Controller) RemoveTeamMember(c *gin.Context) {
	var req request.CreateTeamMember
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	err := ctl.Service.RemoveTeamMember(c.Request.Context(), req)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, "team member removed successfully")
}

func (ctl *Controller) List(c *gin.Context) {
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil {
		response.BadRequest(c, "limit is invalid")
		return
	}
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		response.BadRequest(c, "page is invalid")
		return
	}
	search := c.DefaultQuery("search", "")
	teams, paginate, err := ctl.Service.List(c.Request.Context(), limit, page, search)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.SuccessWithPaginate(c, teams, *paginate)
}

func (ctl *Controller) Get(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.BadRequest(c, "id is required")
		return
	}

	team, err := ctl.Service.Get(c.Request.Context(), id)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, team)
}

func (ctl *Controller) Update(c *gin.Context) {
	id := c.Param("id")

	var req request.UpdateTeam
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	team, err := ctl.Service.Update(c.Request.Context(), req, id)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, team)
}

func (ctl *Controller) Delete(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.BadRequest(c, "id is required")
		return
	}

	err := ctl.Service.Delete(c.Request.Context(), id)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, "team deleted successfully")
}
