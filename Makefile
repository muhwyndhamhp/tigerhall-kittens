include .env

DB_DSN := "$(LIBSQL_URL)?authToken=$(LIBSQL_TOKEN)"

print-dsn:
	@echo $(DB_DSN)

gen:
	@go run github.com/99designs/gqlgen generate

auto-migrate:
	@go run db/migration/main.go

dry-run-migrate:
	@go run db/migration/main.go --dry-run

run:
	@air run
