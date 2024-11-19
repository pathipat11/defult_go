package auth

import (
	"app/app/model"
	"app/app/request"
	"app/app/response"
	"context"

	"net/http"

	"github.com/gin-gonic/gin"
)

func (ctl *Controller) Login(c *gin.Context) {
	var loginUser request.LoginUser
	if err := c.ShouldBindJSON(&loginUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create a context
	ctx := context.Background()

	// Convert loginUser to model.User
	user := model.User{
		Username: loginUser.Username,
		Password: loginUser.Password,
	}

	loggedInUser, err := ctl.Service.Login(ctx, user)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	// Generate a token for the logged-in user
	token, err := ctl.Service.GenerateToken(ctx, loggedInUser.Username, loggedInUser, false)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	response.Success(c, token)
}

func (ctl *Controller) GetUserDetailByToken(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
	if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
		tokenString = tokenString[7:]
	}

	userDetail, err := ctl.Service.GetUserDetailByToken(c.Request.Context(), tokenString)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, userDetail)
}
