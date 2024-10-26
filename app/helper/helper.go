package helper

import (
	"app/app/model"
	provider "app/app/provider/database"
	"app/internal/logger"
	"context"
	"errors"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
)

func CheckMultiplePermissions(ctx context.Context, tokenString string, requiredPermissionIDs int) (string, error) {
	db := provider.DB()

	// Ensure the token is provided
	if tokenString == "" {
		return "", errors.New("authorization token is required")
	}

	// If the token string starts with "Bearer ", remove that prefix
	if len(tokenString) > 7 && strings.HasPrefix(tokenString, "Bearer ") {
		tokenString = tokenString[7:]
	}

	// Parse the JWT token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Ensure the token is using HMAC signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		// Return the secret key for verifying the token
		return []byte(viper.GetString("TOKEN_SECRET_ADMIN")), nil
	})
	if err != nil {
		return "", errors.New("invalid token")
	}

	// Extract the claims (including userID)
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID, ok := claims["id"].(string)
		if !ok {
			logger.Infof("[CheckMultiplePermissions] user_id is missing or not a string in the token claims")
			return "", errors.New("invalid user ID format")
		}

		// Fetch the user from the database using the extracted userID
		var user model.User
		err = db.NewSelect().Model(&user).Where("id = ?", userID).Scan(ctx)
		if err != nil {
			logger.Infof("[CheckMultiplePermissions] Error fetching user: %v", err)
			return "", errors.New("failed to fetch user")
		}

		// Fetch the role_permissions from the database using the extracted userID
		var rolePermissions model.RolePermission
		err = db.NewSelect().Model(&rolePermissions).
			Where("role_id = ?", user.RoleID).
			Where("permission_id = ?", requiredPermissionIDs).
			Scan(ctx)

		// If the permission ID is found, return the user ID
		if err == nil {
			return user.ID, nil
		}

		// If none of the permission IDs are found, return an error
		return user.ID, errors.New("you don't have permission to do this")
	}

	return "", errors.New("invalid token claims")
}

func GetUserByToken(ctx context.Context, tokenString string) (string, error) {
	// db := provider.DB()

	// Ensure the token is provided
	if tokenString == "" {
		return "", errors.New("authorization token is required")
	}

	// If the token string starts with "Bearer ", remove that prefix
	if len(tokenString) > 7 && strings.HasPrefix(tokenString, "Bearer ") {
		tokenString = tokenString[7:]
	}

	// Parse the JWT token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Ensure the token is using HMAC signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		// Return the secret key for verifying the token
		return []byte(viper.GetString("TOKEN_SECRET_USER")), nil
	})
	if err != nil {
		return "", errors.New("invalid token")
	}

	// Extract the claims (including userID)
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID, ok := claims["id"].(string)
		if !ok {
			return "", errors.New("invalid user ID format")
		}

		return userID, nil
	}

	return "", errors.New("invalid token claims")
}
