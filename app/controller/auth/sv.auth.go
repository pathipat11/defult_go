package auth

import (
	"app/app/model"
	"app/app/util/jwt"
	"app/config"
	"context"
	"errors"
	"fmt"
	"time"

	jwtlib "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
)

func (s *Service) GenerateToken(ctx context.Context, authType string, user *model.User, isadmin bool) (string, error) {
	claims := jwtlib.MapClaims{
		"auth_type": authType,
		"id":        user.ID,
		"data": map[string]interface{}{
			"id":       user.ID,
			"username": user.Username,
			"status":   user.Status,
		},
		"nbf": time.Now().Unix(),
		"exp": time.Now().Add(7 * viper.GetDuration("TOKEN_DURATION_USER")).Unix(),
	}

	var secretKey string
	if !isadmin {
		secretKey = viper.GetString("TOKEN_SECRET_USER")
	} else {
		secretKey = viper.GetString("TOKEN_SECRET_ADMIN")
	}

	tokenString, err := jwt.CreateToken(claims, secretKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (s *Service) Login(ctx context.Context, req model.User) (*model.User, error) {
	storedUser, err := s.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Check if the provided password matches the stored password
	if err := bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid username or password")
	}

	return storedUser, nil
}

func (s *Service) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	if username == "" {
		return nil, errors.New("username is required")
	}
	m := model.User{}
	if err := s.db.NewSelect().Model(&m).
		Where("username = ?", username).Scan(ctx); err != nil {
		return nil, err
	}
	return &m, nil
}

func (s *Service) GetUserByID(ctx context.Context, id string) (*model.User, error) {
	if id == "" {
		return nil, errors.New("id is required")
	}
	m := model.User{}
	if err := s.db.NewSelect().Model(&m).
		Where("id = ?", id).Scan(ctx); err != nil {
		return nil, err
	}
	return &m, nil
}

func (s *Service) GetUserDetailByToken(ctx context.Context, tokenString string) (*model.User, error) {
	// Parse the JWT token
	token, err := jwtlib.Parse(tokenString, func(token *jwtlib.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(viper.GetString("TOKEN_SECRET_USER")), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwtlib.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	userID, ok := claims["id"]
	if !ok {
		return nil, errors.New("invalid token payload: id field is missing or not a number")
	}

	var user model.User
	err = s.db.NewSelect().
		Model(&user).
		Where("id = ?", userID).
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *Service) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	m := model.User{}
	if err := s.db.NewSelect().Model(&m).
		Where("email = ?", email).Scan(ctx); err != nil {
		return nil, err
	}
	return &m, nil
}

func (s *Service) ResetPassword(ctx context.Context, email string) error {
	user, err := s.GetUserByEmail(ctx, email)
	if err != nil {
		return err
	}

	password := uuid.New().String()
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)

	if _, err := s.db.NewUpdate().Model(user).Where("email = ?", email).Exec(ctx); err != nil {
		return err
	}

	err = config.SendEmail(email, "Reset Password", "Reset Password", fmt.Sprintf("Your password has been reset to %s", password))
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) ChangePassword(ctx context.Context, userID string, oldPassword string, newPassword string) (*model.User, error) {
	storedUser, err := s.GetUserByID(ctx, userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Check if the provided password matches the stored password
	if err := bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(oldPassword)); err != nil {
		return nil, errors.New("invalid old password")
	}

	// Hash the new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Update the user's password
	storedUser.Password = string(hashedPassword)
	if _, err := s.db.NewUpdate().Model(storedUser).Where("id = ?", userID).Exec(ctx); err != nil {
		return nil, err
	}

	return storedUser, nil
}
