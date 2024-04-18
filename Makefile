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

gen-mock:
	@mockery --all 

test:
	@go test -cover -coverprofile=coverage.out.tmp ./... 
	@cat coverage.out.tmp | grep -v "_mock.go" > coverage.out
	@go tool cover -html=coverage.out 

test-report:
	@go test -cover -coverprofile=coverage.out.tmp ./... 
	@cat coverage.out.tmp | grep -v "_mock.go" > coverage.out
	@go tool cover -html=coverage.out -o coverage.html

deploy:
	@fly deploy
