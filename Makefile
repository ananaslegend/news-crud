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

.PHONY: compose
compose:
	@echo "::> Running docker-compose..."
	@docker-compose up -d

.PHONY: hepl
help:
	@echo "Usage: make <target>"
	@echo ""
	@echo "Targets:"
	@echo "  get-migrate   - Install migrate tool"
	@echo "  migrate-up    - Run migrations up"
	@echo "  migrate-down  - Run migrations down"
	@echo "  build         - Build the application"
	@echo "  run           - Run the application"
	@echo "  test          - Run tests"
	@echo "  compose       - Run docker-compose"
	@echo "  help          - Show this help message"
	@echo ""
	@echo "Variables:"
	@echo "  BINARIES      - Directory to store binaries"
	@echo "  BINARY_NAME   - Name of the binary"
	@echo "  COVER_FILE    - File to store coverage"
	@echo "  MIGRATIONS    - Directory to store migrations"
	@echo "  DB_URL        - Database URL"
	@echo "  DB_DRIVER     - Database driver"
	@echo ""
	@echo "Example:"
	@echo "  make build"
	@echo "  make run"
	@echo "  make test"
	@echo "  make compose"
	@echo "  make migrate-up"
	@echo "  make migrate-down"
	@echo "  make get-migrate"
	@echo "  make help"
	@echo ""
	@echo "For more information, see the Makefile"
	@echo ""