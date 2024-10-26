package auth

import (
	"app/app/enum"
	"app/app/model"
	"app/app/request"
	"app/app/response"
	"app/config"
	"context"
	"encoding/json"

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

func (ctl *Controller) LoginAdmin(c *gin.Context) {
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

	loggedInUser, err := ctl.Service.LoginAdmin(ctx, user)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	// Generate a token for the logged-in user
	token, err := ctl.Service.GenerateToken(ctx, loggedInUser.Username, loggedInUser, true)
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

var oauthStateString = "state"

func (ctl *Controller) LoginGoogle(c *gin.Context) {
	googleOauthConfig := config.GetGoogleOAuthConfig()
	url := googleOauthConfig.AuthCodeURL(oauthStateString)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func (ctl *Controller) GoogleCallback(c *gin.Context) {
	if c.Query("state") != oauthStateString {
		c.JSON(http.StatusBadRequest, gin.H{"error": "State is not valid"})
		return
	}

	code := c.Query("code")
	googleOauthConfig := config.GetGoogleOAuthConfig()
	token, err := googleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Code exchange failed"})
		return
	}

	client := googleOauthConfig.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user info"})
		return
	}
	defer resp.Body.Close()

	// Parse the user info from the response
	userInfo := make(map[string]interface{})
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse user info"})
		return
	}

	user := model.User{
		Email:       userInfo["email"].(string),
		DisplayName: userInfo["name"].(string),
		RoleID:      1,
		Status:      enum.STATUS_ACTIVE,
	}

	// Check if the user already exists
	ex, err := ctl.Service.GetUserByEmail(c.Request.Context(), user.Email)

	// If the user does not exist, create a new user
	if err != nil {
		user, err := ctl.Service.Create(c.Request.Context(), user)
		if err != nil {
			response.InternalError(c, err.Error())
			return
		}
		// Generate JWT token
		jwtToken, err := ctl.Service.GenerateTokenGoogle(user.ID, userInfo)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate JWT"})
			return
		}
		response.Success(c, jwtToken)
	}

	// Generate JWT token
	jwtToken, err := ctl.Service.GenerateTokenGoogle(ex.ID, userInfo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate JWT"})
		return
	}

	// Return the JWT token as the response
	response.Success(c, jwtToken)
}
