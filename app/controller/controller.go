package controller

import (
	"app/app/controller/auth"

	"app/app/controller/user"
	"app/config"
)

type Controller struct {
	AuthCtl *auth.Controller
	UserCtl *user.Controller

	// Other controllers...
}

func New() *Controller {
	// Fetch the initialized DB connection
	db := config.GetDB()
	// db2 := config.GetDB2()
	return &Controller{

		AuthCtl: auth.NewController(db),
		UserCtl: user.NewController(db),
		// Other controllers...
	}
}
