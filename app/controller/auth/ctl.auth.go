package auth

import (
	"app/app/enum"
	"app/app/model"
	provider "app/app/provider/OAuth"
	"app/app/request"
	"app/app/response"
	"app/internal/logger"
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

func (ctl *Controller) LoginGoogle(c *gin.Context) {
	// รับค่า redirect_url จาก query parameter
	redirect := c.Query("redirect_url")
	if redirect == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Redirect URL is required"})
		return
	}

	// เก็บ redirect_url ไว้ใน cookie เพื่อใช้งานในขั้นตอน callback
	c.SetCookie("redirect_url", redirect, 3600, "/", "localhost", false, true)

	googleOauthConfig := provider.GetGoogleOAuthConfig()
	url := googleOauthConfig.AuthCodeURL("state")
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func (ctl *Controller) GoogleCallback(c *gin.Context) {
	// ดึงค่า redirect_url จาก cookie
	redirect, err := c.Cookie("redirect_url")
	if err != nil || redirect == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Redirect URL is required"})
		return
	}

	code := c.Query("code")
	googleOauthConfig := provider.GetGoogleOAuthConfig()
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

	// ตรวจสอบว่าผู้ใช้มีอยู่ในระบบหรือไม่
	ex, err := ctl.Service.GetUserByEmail(c.Request.Context(), user.Email)

	// ถ้าผู้ใช้ยังไม่มีในระบบ ให้สร้างผู้ใช้ใหม่
	if err != nil {
		user, err := ctl.Service.Create(c.Request.Context(), user)
		if err != nil {
			response.InternalError(c, err.Error())
			return
		}
		// สร้าง JWT token
		jwtToken, err := ctl.Service.GenerateTokenGoogle(user.ID, userInfo)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate JWT"})
			return
		}
		// ตั้งค่า JWT token เป็น cookie
		http.SetCookie(c.Writer, &http.Cookie{
			Name:     "auth_token",
			Value:    jwtToken,
			MaxAge:   3600,
			Path:     "/",
			Domain:   "",
			Secure:   false, // ใช้ false ในการพัฒนา และเปลี่ยนเป็น true ใน production
			HttpOnly: false,
			SameSite: http.SameSiteLaxMode, // เปลี่ยนเป็น Lax ในระหว่างการพัฒนา
		})
		logger.Infof("token : %s", jwtToken)
		// Redirect ไปยัง URL ที่ได้รับจาก cookie
		c.Redirect(http.StatusTemporaryRedirect, redirect)
		return
	}

	// สร้าง JWT token
	jwtToken, err := ctl.Service.GenerateTokenGoogle(ex.ID, userInfo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate JWT"})
		return
	}

	// ตั้งค่า JWT token เป็น cookie
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "auth_token",
		Value:    jwtToken,
		MaxAge:   3600,
		Path:     "/",
		Domain:   "",
		Secure:   false, // ใช้ false ในการพัฒนา และเปลี่ยนเป็น true ใน production
		HttpOnly: false,
		SameSite: http.SameSiteLaxMode, // เปลี่ยนเป็น Lax ในระหว่างการพัฒนา
	})

	logger.Infof("token : %s", jwtToken)

	// Redirect ไปยัง URL ที่ได้รับจาก cookie
	c.Redirect(http.StatusTemporaryRedirect, redirect)
}
