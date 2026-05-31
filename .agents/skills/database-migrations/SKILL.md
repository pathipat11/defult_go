---
name: database-migrations
description: Create, register, and run database tables and seeds using the Bun + Cobra migrate workflow in this template. Use when adding/changing models, running migrations, or writing seed data.
---

# Database migrations & seeds

This template uses **Bun** model-based migrations driven by Cobra commands. There are no
timestamped SQL migration files — tables are created/dropped from the registered model list.

## Where things live
- Model structs: `app/model/*.go`
- Registered model list: `database/migrations/Models.go` → `Models()`
- Raw SQL hooks: `Models.go` → `RawBeforeQueryMigrate()` / `RawAfterQueryMigrate()`
- Seeders: `database/seeds/` (registered in `database/seeds/0-base.go` → `Seeds()`)
- Migrate commands: `internal/cmd/migrateCmd.go` and `internal/cmd/model.go`

## Commands
```bash
make migrate-up       # create every table returned by Models()
make migrate-down     # drop every table returned by Models()
make migrate-refresh  # down then up (DESTRUCTIVE: drops all data)
make migrate-seed     # run all seeders registered in Seeds()
```
`migrate refresh` drops all tables — never run it against a shared/production database.

## Adding a new table
1. Create the model in `app/model/<name>.go`:
   - Embed `bun.BaseModel` with `bun:"table:<name>s"`.
   - Embed `CreateUpdateUnixTimestamp` for `created_at`/`updated_at` (unix seconds).
   - Embed `SoftDelete` for soft-delete support (`deleted_at`).
   - Use `type:uuid,default:gen_random_uuid()` for UUID PKs (see `model/user.go`) or
     `type:serial,autoincrement,pk` for integer PKs (see `model/product.go`).
2. Register it in `database/migrations/Models.go`:
   ```go
   func Models() []any {
       return []any{
           (*model.User)(nil),
           (*model.<Feature>)(nil), // add here
       }
   }
   ```
   Order matters if you add foreign keys — create referenced tables first.
3. Run `make migrate-up` (or `make migrate-refresh` in local dev).

## Timestamps
Columns default to `EXTRACT(EPOCH FROM NOW())`. Before an update, call `m.SetUpdateNow()`
so `updated_at` is refreshed. Millisecond variants exist (`CreateUpdateMilliTimestamp`).

## Writing a seeder
1. Add a function in `database/seeds/` with signature `func(*bun.DB) error`:
   ```go
   func <feature>Seed(db *bun.DB) error {
       items := []model.<Feature>{ /* ... */ }
       _, err := db.NewInsert().Model(&items).Exec(context.Background())
       return err
   }
   ```
2. Register it in `Seeds()` in `database/seeds/0-base.go`:
   ```go
   seeder := []func(*bun.DB) error{
       mockUpSeed,
       <feature>Seed, // add here
   }
   ```
3. Run `make migrate-seed`.

## Raw SQL (extensions, indexes)
Add statements to `RawBeforeQueryMigrate()` (e.g. `CREATE EXTENSION ...`) or
`RawAfterQueryMigrate()`. Note: these hooks are defined but the wiring in
`internal/cmd/model.go` for raw queries is commented out — enable `modelRawBeforeQuery` /
`modelRawAfterQuery` there if you need them to run.

## Verify
After model changes, always run `go build ./...` and `go vet ./...` before migrating.
