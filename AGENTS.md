# AGENTS.md

Guidance for AI coding agents (and humans) working in this Go backend template.

## What this project is

A Go REST API template built around a layered architecture:

- **Gin** for HTTP routing/middleware
- **Bun** ORM over **PostgreSQL** (via `pgx`)
- **Cobra** for the CLI entrypoint (`http`, `migrate`, `cmd`)
- **Viper** + `godotenv` for configuration
- **Zap** for logging
- **Redis** and **OAuth/JWT** providers available out of the box

Module path is `app` (see `go.mod`), so internal imports look like `app/app/controller/...`, `app/config`, `app/internal/...`.

## Commands

All commands assume you are at the repo root. The `Makefile` wraps the most common ones.

```bash
make api            # go run . http      -> start the HTTP server
make migrate-up     # go run . migrate up      -> create tables
make migrate-down   # go run . migrate down    -> drop tables
make migrate-seed   # go run . migrate seed    -> run seeders
make migrate-refresh# go run . migrate refresh -> drop + recreate
make cmd-hello      # go run . cmd hello       -> sample console command
```

Direct CLI (without make):

```bash
go run . http
go run . migrate up
go run . cmd hello
```

Build / verify (always run these after making changes):

```bash
go build ./...
go vet ./...
gofmt -l .          # lists files that need formatting; should print nothing
go test ./...       # run the test suite
```

The `Makefile` also wraps these: `make build`, `make vet`, `make fmt-check`, `make test`,
`make lint` (golangci-lint), and `make check` (fmt-check + vet + test).

> Do not start the server with a blocking foreground command during automated work. Use `go build ./...` and `go vet ./...` to verify instead. If a running server is required, start it as a background process.

## Configuration

Config is loaded in `config.Init()` (called from `main.go`) which reads `.env` (via `godotenv`) and falls back to defaults set with Viper.

- Copy `.env.example` to `.env` and fill in values.
- App-level defaults live in `config/config.go` (`APP_NAME`, `APP_PORT`, `APP_ENV`, token + email settings).
- Database options live in `config/database.go`.
- Access config anywhere with `viper.GetString("KEY")`, or use `config.GetDB()` for the `*bun.DB` handle.

When you add a new config value, set a sensible default in `config/config.go` (or the relevant `config/*.go` file) AND document it in `.env.example`.

## Architecture & request flow

```
HTTP request
  -> routes/ (app/routes/*.go)         register endpoints onto gin groups
  -> middleware/ (optional)            auth, activity logging
  -> controller (ctl.*.go)            bind request, call service, shape response
  -> service (sv.*.go)                business logic + Bun queries
  -> model/ (app/model/*.go)          Bun table structs
  -> response/ (app/response/*.go)    response envelope + DTOs
```

Each feature is a self-contained package under `app/controller/<feature>/`:

- `main.go`    — `Controller` + `Service` structs and their `New*` constructors
- `ctl.<x>.go` — HTTP handlers: bind input, call service, map errors to responses
- `sv.<x>.go`  — business logic and Bun database queries

Controllers are wired together in `app/controller/controller.go` (`controller.New()`), and routes are registered in `app/routes/routes.go`.

## Conventions (follow these exactly)

### Controllers
- Bind body with `ctx.Bind(&req)`, URI params with `ctx.BindUri(&id)`.
- On bind error: `logger.Err(err.Error())` then `response.BadRequest(ctx, err.Error())`.
- Services return `(data, mserr bool, err error)` for create/update style ops. When `err != nil`, only surface `err.Error()` to the client if `mserr` is `true` (a "message-safe" business error); otherwise return a generic `"internal server error"`.
- Use `response.Success`, `response.SuccessWithPaginate`, `response.BadRequest`, `response.InternalError`, `response.NotFound`, `response.Unauthorized`, `response.Forbidden`.

### Logging
- Use `logger.Err(...)` / `logger.Info(...)` for plain values (they use `fmt.Sprint`).
- Use `logger.Errf(format, args...)` / `logger.Infof(format, args...)` ONLY with a constant format string.
- Never pass a dynamic string (like `err.Error()`) as the format argument to the `*f` variants — `go vet` will flag it and a stray `%` will corrupt output. Use `logger.Err(err.Error())` instead.

### Pagination defaults (in the controller `List` handler)
```go
if req.Page == 0 { req.Page = 1 }
if req.Size == 0 { req.Size = 10 }
if req.OrderBy == "" { req.OrderBy = "asc" }
if req.SortBy == "" { req.SortBy = "created_at" }
```

### Database / Bun
- Soft delete is enabled via the embedded `SoftDelete` struct; always filter `Where("deleted_at IS NULL")` on reads where appropriate.
- Timestamps come from embedding `CreateUpdateUnixTimestamp`; call `m.SetUpdateNow()` before updates.
- For `LIKE` search build the pattern with plain string concatenation, never `fmt.Sprintf("%" + ... + "%")` (the `%` is treated as a format verb). Pass the value as a bound `?` parameter.
- Never interpolate user input directly into SQL. Column names used in dynamic `ORDER BY` / `search_by` should be validated against an allow-list before being formatted into the query.

### Models & migrations
- New models go in `app/model/<name>.go` embedding `bun.BaseModel` plus timestamp/soft-delete helpers.
- Register the model in `database/migrations/Models.go` (`Models()`), otherwise `migrate up`/`down` will skip it.
- Add seed data in `database/seeds/` and register the seeder in `Seeds()` in `database/seeds/0-base.go`.

## Security notes
- The `/products` routes are protected by `AuthMiddleware` (Bearer JWT). The `/users` routes are currently open — if you add sensitive user operations, add the middleware.
- Passwords are hashed with bcrypt (`golang.org/x/crypto/bcrypt`) in the user service. Never store or log plaintext passwords.
- Don't log request bodies that may contain credentials.
- CORS is driven by `CORS_ALLOW_ORIGINS` (see `corsMiddleware` in `app/routes/routes.go`). `*` reflects any origin; otherwise pass a comma-separated allow-list. Set an explicit list for production.

## HTTP server
- The `http` command (`internal/cmd/httpCmd.go`) runs an `http.Server` with graceful shutdown on `SIGINT`/`SIGTERM` and a 10s drain timeout. It listens on `APP_PORT` and switches Gin to release mode when `APP_ENV=production`.
- Health endpoints live at the root: `/healthz` (liveness) and `/readyz` (pings the DB).

## Response helpers
- `response.BadRequest/Unauthorized/Forbidden/NotFound/InternalError` set the matching HTTP status and a JSON `{status:{code,message}}` envelope. Passing `nil` as the message is safe — it falls back to a default. Don't reintroduce a naked `message.(string)` assertion (it panics on nil).

## Testing
- Tests live alongside code as `_test.go` files. Pure-logic units (response helpers, JWT, helpers) are covered. DB-dependent service tests would need a test database (not yet wired up).
- Run `go test ./...` (or `make test`). Add tests when fixing bugs or adding features.

## Gotchas
- The module path is `app`, so a package at `app/controller/user` imports as `app/app/controller/user`.
- There are two `dbMap` definitions (`config/database.go` and `app/provider/database/database.go`); the active DB handle used by controllers comes from `config.GetDB()`.
- CORS is currently fully open (`AllowAllOrigins: true`) — tighten before production.

## Definition of done
A change is complete when:
1. `go build ./...` succeeds.
2. `go vet ./...` is clean.
3. `gofmt -l .` prints nothing.
4. `go test ./...` passes.
5. New endpoints are registered in `app/routes/` and wired in `app/controller/controller.go`.
6. New models are registered in `database/migrations/Models.go`.
7. New config keys have defaults in `config/` and are documented in `.env.example`.
