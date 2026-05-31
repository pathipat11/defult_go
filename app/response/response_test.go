package response

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func decodeStatus(t *testing.T, body []byte) StatusResponse {
	t.Helper()
	var s StatusResponse
	if err := json.Unmarshal(body, &s); err != nil {
		t.Fatalf("failed to decode body %q: %v", body, err)
	}
	return s
}

func TestErrorHelpersUseCorrectHTTPStatus(t *testing.T) {
	cases := []struct {
		name       string
		fn         func(*gin.Context, any, ...string)
		wantHTTP   int
		wantCode   int64
		wantMsgDef string
	}{
		{"BadRequest", func(c *gin.Context, m any, p ...string) { BadRequest(c, m) }, http.StatusBadRequest, 400, "Bad Request"},
		{"Unauthorized", Unauthorized, http.StatusUnauthorized, 401, "Unauthorized"},
		{"Forbidden", Forbidden, http.StatusForbidden, 403, "Forbidden"},
		{"NotFound", NotFound, http.StatusNotFound, 404, "Not Found"},
		{"InternalError", InternalError, http.StatusInternalServerError, 500, "Internal Server Error"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)

			// Passing nil must not panic and should fall back to a default message.
			tc.fn(ctx, nil)

			if w.Code != tc.wantHTTP {
				t.Errorf("HTTP status = %d, want %d", w.Code, tc.wantHTTP)
			}
			got := decodeStatus(t, w.Body.Bytes())
			if got.Code != tc.wantCode {
				t.Errorf("body code = %d, want %d", got.Code, tc.wantCode)
			}
			if got.Message != tc.wantMsgDef {
				t.Errorf("default message = %q, want %q", got.Message, tc.wantMsgDef)
			}
		})
	}
}

func TestErrorHelpersPassThroughMessage(t *testing.T) {
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	BadRequest(ctx, "email is required")

	got := decodeStatus(t, w.Body.Bytes())
	if got.Message != "email is required" {
		t.Errorf("message = %q, want %q", got.Message, "email is required")
	}
}

func TestSuccessWithPaginateEmpty(t *testing.T) {
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	SuccessWithPaginate(ctx, nil, 10, 1, 0)

	if w.Code != http.StatusOK {
		t.Fatalf("HTTP status = %d, want %d", w.Code, http.StatusOK)
	}
	var resp map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if _, ok := resp["pagination"]; !ok {
		t.Error("expected a pagination field in the response")
	}
}
