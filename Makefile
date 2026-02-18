.PHONY: help test race lint lint-fix coverage docs docs-build docs-image clean

MKDOCS_IMAGE = go-iproto-docs
MKDOCS_PORT  = 8000

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

test: ## Run tests
	go test ./...

race: ## Run tests with race detector
	go test -race ./...

lint: ## Run linters
	golangci-lint run ./...

lint-fix: ## Run linters with auto-fix
	golangci-lint run --fix ./...

coverage: ## Run tests with coverage and generate HTML report
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

docs-image: ## Build MkDocs Docker image
	docker build -t $(MKDOCS_IMAGE) -f docs/Dockerfile .

docs: docs-image ## Serve documentation locally (http://localhost:8000)
	docker run --rm -it -p $(MKDOCS_PORT):8000 -v $(PWD):/docs $(MKDOCS_IMAGE)

docs-build: docs-image ## Build documentation site
	docker run --rm -v $(PWD):/docs $(MKDOCS_IMAGE) build

clean: ## Remove build artifacts
	rm -f coverage.out coverage.html
	rm -rf site/
	rm -rf bin/
