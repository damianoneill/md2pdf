.DEFAULT_GOAL := help

BIN_DIR      := bin
COVERAGE_OUT := coverage.out

# ---------------------------------------------------------------------------
# Language detection — add new languages by extending the ifeq chain
# ---------------------------------------------------------------------------
ifneq ($(wildcard go.mod),)
  LANG := go
else ifneq ($(wildcard pyproject.toml setup.py requirements.txt),)
  LANG := python
else ifneq ($(wildcard package.json),)
  LANG := node
else
  LANG := unknown
endif

# ---------------------------------------------------------------------------
# Per-language command definitions
# ---------------------------------------------------------------------------
ifeq ($(LANG),go)
  CMD_SETUP    := go mod download && go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
  CMD_SYNC     := go mod tidy
  CMD_FMT      := gofmt -w .
  CMD_LINT     := golangci-lint run ./...
  CMD_TEST     := go test ./...
  CMD_COVERAGE := go test -coverprofile=$(COVERAGE_OUT) ./... && go tool cover -html=$(COVERAGE_OUT)
  CMD_BUILD    := go build -o $(BIN_DIR)/ ./...
  CMD_CLEAN    := rm -rf $(BIN_DIR) $(COVERAGE_OUT)
else ifeq ($(LANG),python)
  CMD_SETUP    := uv sync --all-extras
  CMD_SYNC     := uv sync --all-extras
  CMD_FMT      := uv run ruff format .
  CMD_LINT     := uv run ruff check .
  CMD_TEST     := uv run pytest
  CMD_COVERAGE := uv run pytest --cov=. --cov-report=html
  CMD_BUILD    := uv build
  CMD_CLEAN    := rm -rf dist/ build/ .coverage htmlcov/ __pycache__ .pytest_cache
else ifeq ($(LANG),node)
  CMD_SETUP    := npm install
  CMD_SYNC     := npm install
  CMD_FMT      := npx prettier --write .
  CMD_LINT     := npm run lint
  CMD_TEST     := npm test
  CMD_COVERAGE := npm run coverage
  CMD_BUILD    := npm run build
  CMD_CLEAN    := rm -rf node_modules dist $(COVERAGE_OUT)
else
  CMD_SETUP    := $(error No recognised project file found — add go.mod, pyproject.toml, or package.json)
  CMD_SYNC     := $(CMD_SETUP)
  CMD_FMT      := $(CMD_SETUP)
  CMD_LINT     := $(CMD_SETUP)
  CMD_TEST     := $(CMD_SETUP)
  CMD_COVERAGE := $(CMD_SETUP)
  CMD_BUILD    := $(CMD_SETUP)
  CMD_CLEAN    := $(CMD_SETUP)
endif

# ---------------------------------------------------------------------------
# Targets
# ---------------------------------------------------------------------------

## help: Show this help message
.PHONY: help
help:
	@echo "Detected language: $(LANG)"
	@echo ""
	@sed -n 's/^## //p' $(MAKEFILE_LIST) | column -t -s ':' | sed -e 's/^/ /'

## setup: Install dependencies and bootstrap the environment
.PHONY: setup
setup:
	$(CMD_SETUP)

## sync: Sync / tidy dependencies
.PHONY: sync
sync:
	$(CMD_SYNC)

## fmt: Format source code
.PHONY: fmt
fmt:
	$(CMD_FMT)

## lint: Run linters
.PHONY: lint
lint:
	$(CMD_LINT)

## test: Run tests
.PHONY: test
test:
	$(CMD_TEST)

## coverage: Run tests with coverage report
.PHONY: coverage
coverage:
	$(CMD_COVERAGE)

## build: Build / bundle the project
.PHONY: build
build:
	$(CMD_BUILD)

## clean: Remove build artifacts
.PHONY: clean
clean:
	$(CMD_CLEAN)

## ci: Run the full CI pipeline locally (lint → test → build)
.PHONY: ci
ci: lint test build
