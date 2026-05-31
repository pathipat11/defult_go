package helper

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestGetUserByToken(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(*gin.Context)
		want   string
		hasErr bool
	}{
		{
			name:  "no claims",
			setup: func(*gin.Context) {},
			want:  "",
		},
		{
			name: "user_id claim",
			setup: func(c *gin.Context) {
				c.Set("claims", jwt.MapClaims{"user_id": "uuid-1"})
			},
			want: "uuid-1",
		},
		{
			name: "sub fallback",
			setup: func(c *gin.Context) {
				c.Set("claims", jwt.MapClaims{"sub": "uuid-2"})
			},
			want: "uuid-2",
		},
		{
			name: "unexpected claims type does not panic",
			setup: func(c *gin.Context) {
				c.Set("claims", []byte("garbage"))
			},
			want: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			c, _ := gin.CreateTestContext(nil)
			tc.setup(c)

			got, err := GetUserByToken(c)
			if (err != nil) != tc.hasErr {
				t.Fatalf("error = %v, wantErr = %v", err, tc.hasErr)
			}
			if got != tc.want {
				t.Errorf("GetUserByToken() = %q, want %q", got, tc.want)
			}
		})
	}
}
