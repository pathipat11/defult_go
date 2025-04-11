package helper

import (
	"github.com/gin-gonic/gin"
)

func GetUserByToken(ctx *gin.Context) (any, error) {
	claims, exist := ctx.Get("claims")
	if !exist {
		return 0, nil
	}

	return claims, nil
}
