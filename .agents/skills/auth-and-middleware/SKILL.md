---
name: auth-and-middleware
description: Protect routes with JWT auth, read the authenticated user, and add Gin middleware in this template. Use when securing endpoints, issuing/verifying tokens, or adding request middleware.
---

# Auth & middleware

## Login & token issuance
- The login flow lives in `app/controller/auth/` (`POST /api/v1/auth/login`). It looks up the
  user by email, compares the bcrypt hash, and issues a JWT via `jwt.CreateToken`.
- Claims issued: `user_id`, `email`, `exp`, `iat`. Token lifetime comes from
  `TOKEN_DURATION_USER`; the signing secret is `TOKEN_SECRET_USER`.
- Use the same generic error for unknown email and wrong password to avoid account enumeration.

## JWT
- Token helpers: `app/util/jwt/jwt.go` (`VerifyToken`, `CreateToken`).
- Token config defaults: `TOKEN_SECRET_USER`, `TOKEN_DURATION_USER` in `config/config.go`.
- Set real secrets via `.env` — never commit production secrets.

## Protecting a route
Apply `middleware.AuthMiddleware()` per-route or per-group. It expects an
`Authorization: Bearer <token>` header, verifies the JWT, and stores claims on the context.

```go
import "app/app/middleware"

func <Feature>(router *gin.RouterGroup) {
	ctl := controller.New()
	md := middleware.AuthMiddleware()
	g := router.Group("")
	{
		g.GET("/list", md, ctl.<Feature>Ctl.List) // protected
		g.POST("/create", ctl.<Feature>Ctl.Create) // public
	}
}
```
See `app/routes/product.go` for a fully protected group. The `/users` routes are currently
open — add `md` to any user route that performs sensitive operations.

## Reading the authenticated user
In a handler, claims set by the middleware are available via the Gin context:
```go
claims, ok := ctx.Get("claims") // claims is jwt.MapClaims
```
Or use `helper.GetUserByToken(ctx)` which returns the `user_id` claim as a string (used by
the activity-log middleware).

## Activity-log middleware
`middleware.NewLogResponse()` captures request/response bodies and writes an `ActivityLog`
row. It requires a valid user token (calls `helper.GetUserByToken`). Apply it to groups where
you want auditing (see `app/routes/product.go`). Avoid applying it to endpoints that accept
credentials, since it records request bodies.

## Writing a new middleware
Put it in `app/middleware/<name>.go`:
```go
package middleware

import "github.com/gin-gonic/gin"

func MyMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// pre-processing; abort with ctx.AbortWithStatusJSON(...) on failure
		ctx.Next()
		// post-processing
	}
}
```
Return early with `ctx.AbortWithStatusJSON(...)` (or a `response.*` helper) on rejection;
otherwise call `ctx.Next()`.

## Security checklist
- Hash passwords with bcrypt (already done in the user service). Never log plaintext.
- Don't return password hashes in responses — use a response DTO (see `response/user.go`).
- Tighten CORS in `app/routes/routes.go` before production (`AllowAllOrigins` is `true`).
- Flag any new network-exposed endpoint that lacks auth.
