package helper

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// GetUserByToken extracts the user identifier from the JWT claims that the auth
// middleware stored on the context under "claims". It returns an empty string
// when no claims are present. It never panics on an unexpected claims type.
func GetUserByToken(ctx *gin.Context) (string, error) {
	raw, exist := ctx.Get("claims")
	if !exist {
		return "", nil
	}

	claims, ok := raw.(jwt.MapClaims)
	if !ok {
		return "", nil
	}

	// Support common identifier keys.
	for _, key := range []string{"user_id", "sub", "id"} {
		if v, ok := claims[key]; ok {
			if s, ok := v.(string); ok {
				return s, nil
			}
		}
	}

	return "", nil
}
