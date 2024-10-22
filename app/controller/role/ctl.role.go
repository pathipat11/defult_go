package role

import (
	"app/app/request"
	"app/app/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (ctl *Controller) Create(c *gin.Context) {
	var req request.CreateRole
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	role, err := ctl.Service.Create(c.Request.Context(), req)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, role)
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
	roles, paginate, err := ctl.Service.List(c.Request.Context(), limit, page)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.SuccessWithPaginate(c, roles, *paginate)
}

func (ctl *Controller) Get(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.BadRequest(c, "id is required")
		return
	}

	role, err := ctl.Service.Get(c.Request.Context(), id)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, role)
}

func (ctl *Controller) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "id is required")
		return
	}
	var req request.UpdateRole
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	role, err := ctl.Service.Update(c.Request.Context(), req, id)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, role)
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

	response.Success(c, "Role deleted successfully")
}
