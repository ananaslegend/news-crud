BINARIES="./bin"
BINARY_NAME="app"
COVER_FILE=".coverage.out"
MIGRATIONS="./migrations"
DB_URL="postgresql://postgres:postgres@localhost:5432/news-service?sslmode=disable"
DB_DRIVER="postgres"

.PHONY: get-migrate
get-migrate:
	@go install -tags ${DB_DRIVER} github.com/golang-migrate/migrate/v4/cmd/migrate@latest


.PHONY: migrate-up
migrate-up: get-migrate
	@echo "::> Migrate up..."
	@migrate -path ${MIGRATIONS} -database ${DB_URL} up
	@echo "::> Finished!"

.PHONY: migrate-down
migrate-down: get-migrate
	@echo "::> Migrate down..."
	@migrate -path ${MIGRATIONS} -database ${DB_URL} down
	@echo "::> Finished!"

.PHONY: build
build:
	@echo "::> Building..."
	@CGO_ENABLED=0 go build -o ${BINARIES}/${BINARY_NAME} ./cmd/${BINARY_NAME}
	@echo "::> Finished!"

.PHONY: run
run:
	@echo "::> Runnig..."
	@go run ./cmd/${BINARY_NAME}

.PHONY: test
test:
	@echo "::> Running tests..."
	@go test --race -v ./...
	@echo "::> Finished!"