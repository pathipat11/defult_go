package controller

import (
	"app/app/controller/product"
	"app/config"
)

type Controller struct {
	ProductCtl *product.Controller

	// Other controllers...
}

func New() *Controller {
	db := config.GetDB()
	return &Controller{

		ProductCtl: product.NewController(db),

		// Other controllers...
	}
}
