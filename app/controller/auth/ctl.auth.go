package auth

import (
	"app/app/request"
	"app/app/response"
	"app/internal/logger"

	"github.com/gin-gonic/gin"
)

// Login authenticates a user with email + password and returns a JWT.
func (ctl *Controller) Login(ctx *gin.Context) {
	body := request.Login{}
	if err := ctx.ShouldBindJSON(&body); err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}

	data, mserr, err := ctl.Service.Login(ctx, body)
	if err != nil {
		if mserr {
			response.Unauthorized(ctx, err.Error())
			return
		}
		logger.Err(err.Error())
		response.InternalError(ctx, "internal server error")
		return
	}

	response.Success(ctx, data)
}
