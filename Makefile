.PHONY: help install test fmt generate-mocks mock-dependency fmt-dependency lint-dependency check

help: ## Help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

## Install details
install: mock-dependency lint-dependency fmt-dependency ## Install dependencies

mock-dependency:
	@go install github.com/golang/mock/mockgen@latest

lint-dependency:
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.63.4

fmt-dependency:
	@go install golang.org/x/tools/cmd/goimports@latest

fmt: fmt-dependency ## Format code and imports
	@goimports -local github.com/wgarunap/goconf -w .
	@go fmt ./...

lint: lint-dependency ## Run linters
	@golangci-lint run ./...

generate-mocks: mock-dependency ## Generate mocks
	@go generate ./...

test: generate-mocks ## Run tests with race detection
	@go test -race -count 1 -v ./...

check: fmt lint test ## Check code quality
