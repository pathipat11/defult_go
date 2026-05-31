.PHONY: api migrate-up migrate-down migrate-seed migrate-refresh cmd-hello \
	build run test test-cover vet fmt fmt-check lint tidy

# --- Run ---
api:
	go run . http

run: api

# --- Migrations ---
migrate-up:
	go run . migrate up

migrate-down:
	go run . migrate down

migrate-seed:
	go run . migrate seed

migrate-refresh:
	go run . migrate refresh

# --- Console commands ---
cmd-hello:
	go run . cmd hello

# --- Build ---
build:
	go build -o dist/app .

# --- Quality ---
test:
	go test ./...

test-cover:
	go test -cover ./...

vet:
	go vet ./...

fmt:
	gofmt -w .

fmt-check:
	@test -z "$$(gofmt -l .)" || (echo "These files need gofmt:" && gofmt -l . && exit 1)

lint:
	golangci-lint run

tidy:
	go mod tidy

# Run all checks the way CI would.
check: fmt-check vet test
