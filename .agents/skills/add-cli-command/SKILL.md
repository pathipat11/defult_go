---
name: add-cli-command
description: Add a new Cobra CLI command (console command or top-level command) to the app. Use when you need a one-off task, worker, or maintenance command runnable via "go run . ...".
---

# Add a CLI command

The CLI is built with Cobra and assembled in `main.go`:

```
app
├── cmd <subcommand>      # console commands (app/console)
├── http                  # HTTP server      (internal/cmd/httpCmd.go)
└── migrate <up|down|...> # migrations        (internal/cmd/migrateCmd.go)
```

## Option A: a console command (most common)
Console commands live under `app/console/` and run as `go run . cmd <name>`.

1. Add the command in `app/console/cmd.go` (or a new file in the same package):
   ```go
   func myTaskCmd() *cobra.Command {
       return &cobra.Command{
           Use:  "mytask",
           Args: cmd.NotReqArgs, // rejects positional args; omit if you accept args
           Run: func(c *cobra.Command, args []string) {
               logger.Infof("running mytask")
               // ... do work; use config.GetDB() if you need the database
           },
       }
   }
   ```
2. Register it in `app/console/kernel.go`:
   ```go
   func Commands() []*cobra.Command {
       return []*cobra.Command{
           helloCmd(),
           myTaskCmd(), // add here
       }
   }
   ```
3. Run it: `go run . cmd mytask`.

> If your command touches the database, the HTTP path initializes config in `main.go` via
> `config.Init()`, which connects the DB. `config.GetDB()` returns the handle.

## Option B: a top-level command
For commands that need their own DB open/close lifecycle (like `migrate`), follow
`internal/cmd/migrateCmd.go`:
- Use `PersistentPreRunE` / `PersistentPostRunE` to open/close providers.
- Build the command in `internal/cmd/<name>Cmd.go` returning a `*cobra.Command`.
- Register it in `main.go`:
  ```go
  cmda.AddCommand(cmd.MyCommand())
  ```

## Helpers
- `cmd.NotReqArgs` (`internal/cmd/cmd.go`) — validator that rejects positional args.
- Add a `Makefile` target for discoverability:
  ```make
  cmd-mytask:
  	go run . cmd mytask
  ```

## Verify
```bash
go build ./...
go vet ./...
go run . cmd mytask   # smoke test
```
