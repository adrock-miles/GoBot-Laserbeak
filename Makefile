BINARY  := laserbeak
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS := -ldflags "-X github.com/adrock-miles/GoBot-Laserbeak/cmd.Version=$(VERSION)"

.PHONY: help build clean run docker-build docker-up docker-down docs docs-serve

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

build: ## Build the binary
	go build $(LDFLAGS) -o $(BINARY) .

clean: ## Remove build artifacts
	rm -f $(BINARY)

run: build ## Build and run the bot
	./$(BINARY) serve

docker-build: ## Build Docker image
	docker compose build

docker-up: ## Start containers in background
	docker compose up -d

docker-down: ## Stop containers
	docker compose down

docs: ## Build the documentation site
	cd docs && npm run build

docs-serve: ## Start local docs dev server
	cd docs && npm start
