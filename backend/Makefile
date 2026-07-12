# Load POSTGRES_* from .env for DB_DSN / migrate targets. Do not include the full file —
# multiline values (e.g. OIDC_RSA_PRIVATE_KEY_PEM) are valid in shell but break GNU make.
ifneq (,$(wildcard cmd/server/.env))
$(foreach line,$(shell grep -E '^POSTGRES_' cmd/server/.env),$(eval $(line)))
export POSTGRES_HOST POSTGRES_PORT POSTGRES_USER POSTGRES_PASSWORD POSTGRES_DB POSTGRES_SCHEMA
endif

.PHONY: lint mocks tests test-services test-utils test-handlers test-repositories test-coverage test-coverage-html test-race test-verbose test-specific test-specific-verbose test-specific-coverage
DB_DSN ?= postgresql://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable$(if $(POSTGRES_SCHEMA),&search_path=$(POSTGRES_SCHEMA))
MIGRATION_DIR ?= file://cmd/migrations/sql
DB_DEV_URL ?= docker://postgres/18/dev

bootstrap: container-up migrate-up up

.PHONY: lint
lint:
	go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.11.0 run ./... --config .golangci.yml

mocks:
	mockery --case snake --dir ./repositories --all --output ./mocks/repositories
	mockery --case snake --dir ./adapters --all --output ./mocks/adapters

build:
	@cd cmd/${cmd} && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ${service_name} .

.PHONY: build-prod
build-prod:
	@test -d internal/static/dist || (echo "Run 'make copy-frontend' from repo root first" && exit 1)
	CGO_ENABLED=0 go build -tags embedfrontend -o gateforge-iam-server ./cmd/server

dep:
	@go mod tidy

container-up:
	set -a && source cmd/server/.env && set +a && docker compose up -d

container-down:
	docker compose down

up:
	cd cmd/server && go run main.go

major-version-update:
	go get -u -t ./...

minor-version-update:
	go get -u ./...

.PHONY: migrate-inspect
migrate-inspect:
	atlas schema inspect --url "$(DB_DSN)"

.PHONY: migrate-up
migrate-up:
	atlas migrate apply --dir "$(MIGRATION_DIR)" --url "$(DB_DSN)"

.PHONY: migrate-up-preview
migrate-up-preview:
	atlas migrate apply --dir "$(MIGRATION_DIR)" --url "$(DB_DSN)" --dry-run

.PHONY: migrate-down
migrate-down:
	atlas migrate down --dir "$(MIGRATION_DIR)" --url "$(DB_DSN)" --dev-url "$(DB_DEV_URL)"

.PHONY: migrate-down-preview
migrate-down-preview:
	atlas migrate down --dir "$(MIGRATION_DIR)" --url "$(DB_DSN)" --dev-url "$(DB_DEV_URL)" --dry-run

.PHONY: migrate-create
migrate-create:
	@if [ -z "$(name)" ]; then \
		echo "❌ Missing migration name. Usage: make migrate_create name=<migration_name>"; \
		exit 1; \
	fi
	atlas migrate new $(name)

.PHONY: migrate-reset
migrate-reset:
	atlas schema clean --url "$(DB_DSN)"

.PHONY: migrate-generate
migrate-generate:
	atlas migrate diff --env gorm ${name}

.PHONY: migrate-status
migrate-status:
	atlas migrate status --dir "$(MIGRATION_DIR)" --url "$(DB_DSN)"

.PHONY: migrate-hash
migrate-hash:
	atlas migrate hash --dir "$(MIGRATION_DIR)"

format:
	go fmt ./...

swagger-load:
	GOFLAGS=-buildvcs=false swag init \
		-g main.go \
		-d ./cmd/server,./internal/handlers,./internal/middlewares,./internal/services,./internal/repositories,./internal/models,./internal/utils,./internal/config,./internal/constants,./internal/dtos,./internal/logger,./internal/db,./internal/domains \
		--output ./docs

# Test targets
tests:
	@echo "Running all tests with coverage and race detection..."
	@go test -v -cover -race -timeout 300s -count=1 ./...

test-services:
	@echo "Running service layer tests..."
	@go test -v -cover ./internal/services

test-utils:
	@echo "Running utility tests..."
	@go test -v -cover ./internal/utils

test-handlers:
	@echo "Running handler tests..."
	@go test -v -cover ./internal/handlers

test-repositories:
	@echo "Running repository tests..."
	@go test -v -cover ./internal/repositories

test-coverage:
	@echo "Running tests with coverage report..."
	@mkdir -p docs/test
	@go test -coverprofile=docs/test/coverage.out ./...
	@go tool cover -func=docs/test/coverage.out > docs/test/coverage.txt
	@echo "✓ Wrote docs/test/coverage.out and docs/test/coverage.txt ($$(pwd)/docs/test/)"

test-coverage-html:
	@echo "Generating HTML coverage report..."
	@mkdir -p docs/test
	@go test -coverprofile=docs/test/coverage.out ./... || true
	@if [ -f docs/test/coverage.out ]; then \
		go tool cover -html=docs/test/coverage.out -o docs/test/coverage.html; \
		echo ""; \
		echo "✓ Coverage report generated successfully!"; \
		echo "  File: docs/test/coverage.html ($$(pwd)/docs/test/coverage.html)"; \
		echo "  Open it in your browser to view the report."; \
	else \
		echo ""; \
		echo "✗ Error: docs/test/coverage.out was not generated."; \
		echo "  Some tests may have failed before coverage data could be collected."; \
		exit 1; \
	fi

test-race:
	@echo "Running tests with race detection..."
	@go test -race -timeout 300s ./...

test-verbose:
	@echo "Running tests with verbose output..."
	@go test -v ./...

test-specific:
	@echo "Running specific test: $(TEST)"
	@if [ -z "$(TEST)" ]; then \
		echo "Usage: make test-specific TEST=TestUserService_Create"; \
		exit 1; \
	fi
	@go test -run $(TEST) ./internal/services

test-specific-verbose:
	@echo "Running specific test with verbose output: $(TEST)"
	@if [ -z "$(TEST)" ]; then \
		echo "Usage: make test-specific-verbose TEST=TestUserService_Create"; \
		exit 1; \
	fi
	@go test -v -run $(TEST) ./internal/services

test-specific-coverage:
	@echo "Running specific test with coverage: $(TEST)"
	@if [ -z "$(TEST)" ]; then \
		echo "Usage: make test-specific-coverage TEST=TestUserService_Create"; \
		exit 1; \
	fi
	@go test -cover -run $(TEST) ./internal/services