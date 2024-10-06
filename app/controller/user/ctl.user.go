package user

import (
	"app/app/enum"
	"app/app/model"
	"app/app/request"
	"app/app/response"
	"app/internal/logger"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func (ctl *Controller) Create(c *gin.Context) {
	req := request.CreateUser{}
	if err := c.Bind(&req); err != nil {
		logger.Infof("[%s-create]: %v", ctl.Name, err)
		response.BadRequest(c, err.Error())
		return
	}

	logger.Infof("Request Payload: %+v", req) // Log the request payload

	// Parse birthdate
	birthdate, err := time.Parse("2006-01-02", req.Birthdate)
	if err != nil {
		logger.Infof("[%s-create]: %v", ctl.Name, err)
		response.BadRequest(c, "Invalid birthdate format")
		return
	}

	// Hash the password before passing to service
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.Infof("[%s-create]: %v", ctl.Name, err)
		response.InternalError(c, "Failed to hash password")
		return
	}

	user := model.User{
		Firstname:          req.Firstname,
		Lastname:           req.Lastname,
		Nickname:           req.Nickname,
		CitizenID:          req.CitizenID,
		Birthdate:          birthdate,
		Gender:             req.Gender,
		Nationality:        req.Nationality,
		RelationshipStatus: req.RelationshipStatus,
		Address1:           req.Address1,
		Address2:           req.Address2,
		MobileNo:           req.MobileNo,
		Email:              req.Email,
		RoleID:             1,
		Status:             enum.STATUS_ACTIVE,
		Username:           req.Username,
		Password:           string(hashedPassword),
		Points:             0,
	}

	_, err = ctl.Service.Create(c, user)
	if err != nil {
		logger.Infof("[%s-create]: %v", ctl.Name, err)
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, nil)
}

func (ctl *Controller) List(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "0"))
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	search := c.DefaultQuery("search", "")

	users, err := ctl.Service.List(c.Request.Context(), limit, page, search)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.SuccessWithPaginate(c, users.Users, users.Pagination)
}

func (ctl *Controller) ListSingle(c *gin.Context) {
	userIDStr, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	userID := uint(userIDStr)

	users, err := ctl.Service.ListSingle(c.Request.Context(), userID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, users.Users)
}

// func (ctl *Controller) Update(c *gin.Context) {
// 	userIDStr, _ := strconv.ParseUint(c.Param("id"), 10, 64)
// 	// userID := uint(userIDStr)
// 	req := request.UpdateUser{}
// 	if err := c.Bind(&req); err != nil {
// 		logger.Infof("[%s-update]: %v", ctl.Name, err)
// 		response.BadRequest(c, err.Error())
// 		return
// 	}

// 	logger.Infof("Request Payload: %+v", req) // Log the request payload

// 	// var birthdate time.Time
// 	var err error

// 	if req.Birthdate != "" {
// 		// Parse birthdate only if it's provided
// 		birthdate, err = time.Parse("2006-01-02", req.Birthdate)
// 		if err != nil {
// 			logger.Infof("[%s-update]: %v", ctl.Name, err)
// 			response.BadRequest(c, "Invalid birthdate format")
// 			return
// 		}
// 	}

// 	// user := model.User{
// 	// 	Firstname:          req.Firstname,
// 	// 	Lastname:           req.Lastname,
// 	// 	Nickname:           req.Nickname,
// 	// 	Birthdate:          birthdate,
// 	// 	Gender:             req.Gender,
// 	// 	Nationality:        req.Nationality,
// 	// 	RelationshipStatus: req.RelationshipStatus,
// 	// 	Address1:           req.Address1,
// 	// 	Address2:           req.Address2,
// 	// 	MobileNo:           req.MobileNo,
// 	// 	Email:              req.Email,
// 	// 	RoleID:             req.RoleID,
// 	// }

// 	// _, err = ctl.Service.Update(c, user, userID)
// 	// if err != nil {
// 	// 	logger.Infof("[%s-update]: %v", ctl.Name, err)
// 	// 	response.InternalError(c, err.Error())
// 	// 	return
// 	// }

// 	response.Success(c, nil)
// }

func (ctl *Controller) SoftDelete(c *gin.Context) {
	userIDStr, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	userID := uint(userIDStr)

	err := ctl.Service.SoftDelete(c, userID)
	if err != nil {
		logger.Infof("[%s-soft-delete]: %v", ctl.Name, err)
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, nil)
}

// func (ctl *Controller) Delete(c *gin.Context) {
// 	userIDStr, _ := strconv.ParseUint(c.Param("id"), 10, 64)
// 	userID := uint(userIDStr)

// 	err := ctl.Service.Delete(c, userID)
// 	if err != nil {
// 		logger.Infof("[%s-delete]: %v", ctl.Name, err)
// 		response.InternalError(c, err.Error())
// 		return
// 	}

// 	response.Success(c, nil)
// }
