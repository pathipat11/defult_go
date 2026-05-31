package jwt

import (
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
)

func TestCreateAndVerifyToken(t *testing.T) {
	const secret = "test-secret"
	viper.Set("TOKEN_SECRET_USER", secret)
	t.Cleanup(func() { viper.Set("TOKEN_SECRET_USER", "") })

	token, err := CreateToken(jwt.MapClaims{"user_id": "abc-123"}, secret)
	if err != nil {
		t.Fatalf("CreateToken returned error: %v", err)
	}
	if token == "" {
		t.Fatal("CreateToken returned an empty token")
	}

	claims, err := VerifyToken(token)
	if err != nil {
		t.Fatalf("VerifyToken returned error: %v", err)
	}
	if claims["user_id"] != "abc-123" {
		t.Errorf("user_id claim = %v, want %q", claims["user_id"], "abc-123")
	}
}

func TestVerifyTokenRejectsWrongSecret(t *testing.T) {
	viper.Set("TOKEN_SECRET_USER", "the-real-secret")
	t.Cleanup(func() { viper.Set("TOKEN_SECRET_USER", "") })

	token, err := CreateToken(jwt.MapClaims{"user_id": "abc-123"}, "a-different-secret")
	if err != nil {
		t.Fatalf("CreateToken returned error: %v", err)
	}

	if _, err := VerifyToken(token); err == nil {
		t.Error("expected VerifyToken to reject a token signed with the wrong secret")
	}
}

func TestVerifyTokenRejectsGarbage(t *testing.T) {
	viper.Set("TOKEN_SECRET_USER", "secret")
	t.Cleanup(func() { viper.Set("TOKEN_SECRET_USER", "") })

	if _, err := VerifyToken("not-a-real-token"); err == nil {
		t.Error("expected VerifyToken to reject a malformed token")
	}
}
