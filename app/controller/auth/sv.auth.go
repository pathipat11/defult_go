package auth

import (
	"context"
	"errors"
	"time"

	"app/app/model"
	"app/app/request"
	"app/app/response"
	jwtutil "app/app/util/jwt"

	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
)

// errInvalidCredentials is returned for both unknown emails and wrong passwords
// so the API doesn't reveal which accounts exist.
var errInvalidCredentials = errors.New("invalid email or password")

// Login verifies the credentials and, on success, returns a signed JWT plus the
// authenticated user. The bool indicates whether err is safe to show the client.
func (s *Service) Login(ctx context.Context, req request.Login) (*response.LoginResult, bool, error) {
	m := model.User{}
	err := s.db.NewSelect().
		Model(&m).
		Where("email = ?", req.Email).
		Where("deleted_at IS NULL").
		Scan(ctx)
	if err != nil {
		// No row -> treat as invalid credentials, not an internal error.
		return nil, true, errInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(m.Password), []byte(req.Password)); err != nil {
		return nil, true, errInvalidCredentials
	}

	duration := viper.GetDuration("TOKEN_DURATION_USER")
	if duration <= 0 {
		duration = 24 * time.Hour
	}
	expiresAt := time.Now().Add(duration)

	claims := jwt.MapClaims{
		"user_id": m.ID,
		"email":   m.Email,
		"exp":     expiresAt.Unix(),
		"iat":     time.Now().Unix(),
	}

	token, err := jwtutil.CreateToken(claims, viper.GetString("TOKEN_SECRET_USER"))
	if err != nil {
		return nil, false, err
	}

	return &response.LoginResult{
		AccessToken: token,
		TokenType:   "Bearer",
		ExpiresIn:   int64(duration.Seconds()),
		User: response.AuthUser{
			ID:        m.ID,
			FirstName: m.FirstName,
			LastName:  m.LastName,
			Email:     m.Email,
		},
	}, false, nil
}
