# Go Backend Template

A batteries-included starting point for building REST APIs in Go. It ships with a clean
layered architecture, a CLI for running the server and migrations, and ready-to-use
providers for PostgreSQL, Redis, JWT auth, and OAuth.

## Stack

| Concern        | Library |
|----------------|---------|
| HTTP framework | [Gin](https://github.com/gin-gonic/gin) |
| ORM            | [Bun](https://bun.uptrace.dev/) over PostgreSQL (`pgx`) |
| CLI            | [Cobra](https://github.com/spf13/cobra) |
| Config         | [Viper](https://github.com/spf13/viper) + `godotenv` |
| Logging        | [Zap](https://github.com/uber-go/zap) |
| Auth           | [golang-jwt v5](https://github.com/golang-jwt/jwt) + bcrypt |
| Cache          | [go-redis v9](https://github.com/redis/go-redis) |
| Tracing        | OpenTelemetry (otelgin) |

Requires **Go 1.24+**.

## Project layout

```
.
├── main.go                  # CLI entrypoint (cobra); wires http/migrate/cmd
├── config/                  # config loading + DB registration
├── app/
│   ├── console/             # custom `cmd` console commands
│   ├── controller/          # feature packages (controller + service)
│   │   └── <feature>/
│   │       ├── main.go      # Controller + Service structs
│   │       ├── ctl.<x>.go   # HTTP handlers
│   │       └── sv.<x>.go    # business logic + Bun queries
│   ├── middleware/          # auth, activity-log middleware
│   ├── model/               # Bun table models + timestamp/soft-delete helpers
│   ├── request/             # request DTOs (bind targets)
│   ├── response/            # response envelope + read DTOs
│   ├── routes/              # route registration
│   ├── provider/            # database / redis / oauth providers
│   ├── util/jwt/            # JWT sign/verify
│   └── helper/              # shared helpers
├── database/
│   ├── migrations/Models.go # registered models for migrate up/down
│   └── seeds/               # seeders
├── internal/
│   ├── cmd/                 # http + migrate command implementations
│   └── logger/              # zap logger wrapper
├── Dockerfile
└── docker-compose.yml
```

## Getting started

### 1. Clone & configure

```bash
cp .env.example .env
# edit .env with your database credentials
```

### 2. Start dependencies (PostgreSQL + Redis)

```bash
docker compose up -d postgres redis
```

### 3. Run migrations and seed

```bash
make migrate-up
make migrate-seed
```

### 4. Start the API

```bash
make api          # equivalent to: go run . http
```

The server listens on `APP_PORT` (default `8080`) and shuts down gracefully on
`SIGINT`/`SIGTERM`, draining in-flight requests for up to 10 seconds.

Health endpoints:

```bash
curl http://localhost:8080/healthz   # liveness: process is up
curl http://localhost:8080/readyz    # readiness: pings the database
```

## Configuration

Configuration is read from environment variables (and `.env` in local dev), with defaults
defined in `config/`. Key variables:

| Variable             | Default        | Description |
|----------------------|----------------|-------------|
| `APP_NAME`           | `app`          | App name (used for tracing) |
| `APP_PORT`           | `8080`         | HTTP listen port |
| `APP_ENV`            | `development`  | Environment name (`production` enables Gin release mode) |
| `CORS_ALLOW_ORIGINS` | `*`            | `*` reflects any origin, or a comma-separated allow-list |
| `DB_HOST`            | `127.0.0.1`    | PostgreSQL host |
| `DB_PORT`            | `5432`         | PostgreSQL port |
| `DB_DATABASE`        | `postgres`     | Database name |
| `DB_USER`            | `postgres`     | Database user |
| `DB_PASSWORD`        | —              | Database password |
| `DB_SSLMODE`         | `disable`      | `pgx` sslmode |
| `DB_DSN`             | —              | Full DSN (overrides the individual `DB_*` values) |
| `TZ`                 | `Asia/Bangkok` | Time zone |
| `TOKEN_SECRET_USER`  | `secret`       | JWT signing secret (change in production) |
| `TOKEN_DURATION_USER`| `24h`          | JWT lifetime |
| `EMAIL_HOST` / `EMAIL_PORT` / `EMAIL_USERNAME` / `EMAIL_PASSWORD` | — | SMTP settings |
| `REDIRECT_URL` / `CLIENT_ID` / `CLIENT_SECRET` | — | Google OAuth settings |
| `DEBUG`              | `false`        | Enables verbose Bun query logging |

## CLI commands

```bash
go run . http               # start the HTTP server
go run . migrate up         # create tables
go run . migrate down       # drop tables
go run . migrate refresh    # drop + recreate (DESTRUCTIVE)
go run . migrate seed       # run seeders
go run . cmd hello          # sample console command
```

The `Makefile` wraps these as `make api`, `make migrate-up`, `make migrate-down`,
`make migrate-refresh`, `make migrate-seed`, `make cmd-hello`.

## API

All feature routes are mounted under `/api/v1` (see `app/routes/routes.go`).
Health endpoints (`/healthz`, `/readyz`) are mounted at the root.

### Users — `/api/v1/users`
| Method | Path           | Description |
|--------|----------------|-------------|
| POST   | `/create`      | Create a user |
| GET    | `/list`        | List users (paginated) |
| GET    | `/:id`         | Get a user |
| PATCH  | `/:id`         | Update a user |
| DELETE | `/:id`         | Soft-delete a user |

### Products — `/api/v1/products` (JWT protected)
| Method | Path           | Description |
|--------|----------------|-------------|
| POST   | `/create`      | Create a product |
| GET    | `/list`        | List products (paginated) |
| GET    | `/:id`         | Get a product |
| PATCH  | `/:id`         | Update a product |
| DELETE | `/:id`         | Delete a product |

List endpoints accept query params: `page`, `size`, `search`, `search_by`, `sort_by`,
`order_by`.

Protected endpoints require an `Authorization: Bearer <token>` header.

### Response envelope
```json
{
  "status": { "code": 200, "message": "Success" },
  "data": { },
  "pagination": { "page": 1, "size": 10, "total": 42 }
}
```

## Development

```bash
go build ./...   # compile everything
go vet ./...     # static analysis
gofmt -l .       # formatting check (should print nothing)
go test ./...    # run tests
```

The `Makefile` also provides: `make build`, `make test`, `make test-cover`, `make vet`,
`make fmt`, `make fmt-check`, `make lint` (requires
[golangci-lint](https://golangci-lint.run/)), and `make check` (fmt-check + vet + test).

Verify the build, vet, and tests are clean before committing.

## Docker

```bash
docker compose up --build      # builds the app + starts postgres/redis/redisinsight
```

The `main` service runs the `http` command and connects to the bundled Postgres.

## Working with AI agents

This repo includes [`AGENTS.md`](./AGENTS.md) with conventions, and task-specific guides
under [`.agents/skills/`](./.agents/skills/):

- `add-crud-resource` — scaffold a new entity end-to-end
- `database-migrations` — models, migrations, and seeds
- `auth-and-middleware` — JWT auth and Gin middleware
- `add-cli-command` — add a Cobra CLI command

## Notes & TODO

- CORS is configurable via `CORS_ALLOW_ORIGINS` (defaults to reflecting any origin in dev).
  Set an explicit allow-list before production.
- `/users` routes are unauthenticated in the template; add `AuthMiddleware()` for sensitive ops.
- Change `TOKEN_SECRET_USER` and database credentials before deploying.
