---
name: add-crud-resource
description: Scaffold a new CRUD resource (controller + service + request/response + routes) following this template's layered conventions. Use when adding a new entity/endpoint group such as "orders", "categories", etc.
---

# Add a CRUD resource

Use this when adding a new domain entity with HTTP endpoints. Replace `<feature>` with the
lowercase singular name (e.g. `order`) and `<Feature>` with the PascalCase name (e.g. `Order`).

## Steps

### 1. Model — `app/model/<feature>.go`
```go
package model

import "github.com/uptrace/bun"

type <Feature> struct {
	bun.BaseModel `bun:"table:<feature>s"`

	ID   int64  `bun:",type:serial,autoincrement,pk"`
	Name string `bun:"name,notnull"`
	// ... other columns

	CreateUpdateUnixTimestamp
	SoftDelete
}
```
Register it in `database/migrations/Models.go` inside `Models()`:
```go
(*model.<Feature>)(nil),
```

### 2. Requests — `app/request/<feature>.go`
```go
package request

type Create<Feature> struct {
	Name string `json:"name"`
}

type Update<Feature> struct {
	Create<Feature>
}

type GetByID<Feature> struct {
	ID int64 `uri:"id" binding:"required"`
}

type List<Feature> struct {
	Page     int    `form:"page"`
	Size     int    `form:"size"`
	Search   string `form:"search"`
	SearchBy string `form:"search_by"`
	SortBy   string `form:"sort_by"`
	OrderBy  string `form:"order_by"`
}
```

### 3. Response DTO (optional) — `app/response/<feature>.go`
Define a `List<Feature>` struct with `bun:` + `json:` tags for read projections.

### 4. Controller package — `app/controller/<feature>/`

`main.go`:
```go
package <feature>

import "github.com/uptrace/bun"

type Controller struct {
	Name    string
	Service *Service
}

func NewController(db *bun.DB) *Controller {
	return &Controller{Name: "<feature>-ctl", Service: NewService(db)}
}

type Service struct{ db *bun.DB }

func NewService(db *bun.DB) *Service { return &Service{db: db} }
```

`ctl.<feature>.go` — handlers (Create/Update/List/Get/Delete). Follow the pattern in
`app/controller/user/ctl.user.go`. Pagination defaults in `List`:
```go
if req.Page == 0 { req.Page = 1 }
if req.Size == 0 { req.Size = 10 }
if req.OrderBy == "" { req.OrderBy = "asc" }
if req.SortBy == "" { req.SortBy = "created_at" }
```

`sv.<feature>.go` — Bun queries. Follow `app/controller/user/sv.user.go`:
- Return `(data, mserr bool, err error)` for create/update/delete-with-validation.
- Filter `Where("deleted_at IS NULL")` on reads.
- Build LIKE patterns with string concatenation: `"%" + strings.ToLower(req.Search) + "%"` (never `fmt.Sprintf` with a `%` literal).
- Validate dynamic `search_by` / `sort_by` against an allow-list before formatting into SQL.

### 5. Wire the controller — `app/controller/controller.go`
```go
type Controller struct {
	// ...existing
	<Feature>Ctl *<feature>.Controller
}

func New() *Controller {
	db := config.GetDB()
	return &Controller{
		// ...existing
		<Feature>Ctl: <feature>.NewController(db),
	}
}
```

### 6. Routes — `app/routes/<feature>.go`
```go
package routes

import (
	"app/app/controller"

	"github.com/gin-gonic/gin"
)

func <Feature>(router *gin.RouterGroup) {
	ctl := controller.New()
	g := router.Group("")
	{
		g.POST("/create", ctl.<Feature>Ctl.Create)
		g.GET("/list", ctl.<Feature>Ctl.List)
		g.GET("/:id", ctl.<Feature>Ctl.Get)
		g.PATCH("/:id", ctl.<Feature>Ctl.Update)
		g.DELETE("/:id", ctl.<Feature>Ctl.Delete)
	}
}
```
Add an `AuthMiddleware()` to protected routes (see `app/routes/product.go`).

Register the group in `app/routes/routes.go`:
```go
<Feature>(apiV1.Group("/<feature>s"))
```

### 7. Verify
```bash
go build ./...
go vet ./...
gofmt -l .
```
All three must be clean. Then optionally `make migrate-up` to create the new table.
