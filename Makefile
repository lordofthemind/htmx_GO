# Environment variables for MongoDB
MONGODB_CONTAINER_NAME ?= HTMX_GO
MONGODB_IMAGE_TAG ?= latest
MONGODB_DB_NAME ?= htmx_go
MONGODB_PORT ?= 27017

# Colors for help command
CYAN := \033[36m
RESET := \033[0m

# Docker MongoDB commands
crtmgct: ## Create and start the MongoDB container
	@echo "Creating and starting MongoDB container..."
	docker run --name $(MONGODB_CONTAINER_NAME) -p $(MONGODB_PORT):27017 -d mongo:$(MONGODB_IMAGE_TAG)

strmgct: ## Start the MongoDB container
	@echo "Starting MongoDB container..."
	docker start $(MONGODB_CONTAINER_NAME)

stpmgct: ## Stop the MongoDB container
	@echo "Stopping MongoDB container..."
	docker stop $(MONGODB_CONTAINER_NAME)

rmvmgct: ## Remove the MongoDB container
	@echo "Removing MongoDB container..."
	docker rm $(MONGODB_CONTAINER_NAME)

# MongoDB database commands
crtmgdb: strmgct ## Create MongoDB database
	@echo "Creating MongoDB database..."
	docker exec -it $(MONGODB_CONTAINER_NAME) mongosh --eval "use $(MONGODB_DB_NAME)"

drpmgdb: strmgct ## Drop MongoDB database
	@echo "Dropping MongoDB database..."
	docker exec -it $(MONGODB_CONTAINER_NAME) mongosh --eval "db.getSiblingDB('$(MONGODB_DB_NAME)').dropDatabase()"

# Build commands
build_win: ## Build the application for Windows
	GOOS=windows GOARCH=amd64 go build -o ./build/win/HTMX_GO.exe main.go

build_lin: ## Build the application for Linux
	GOOS=linux GOARCH=amd64 go build -o ./build/linux/HTMX_GO main.go

build_mac: ## Build the application for MacOS
	GOOS=darwin GOARCH=amd64 go build -o ./build/mac/HTMX_GO main.go

build: ## Build the application for all platforms and saves it in build folder
	$(MAKE) build_lin
	$(MAKE) build_win
	$(MAKE) build_mac

# Testing commands
test: ## Run tests
	go test ./...

lint: ## Run linter
	golangci-lint run

# Utility commands
clean: ## Clean the build and the logs
	rm -rf build/
	rm -rf logs/

tr: ## Generate directory tree
	tree > tree.txt

# Help command
help: ## Show this help message
	@echo "Available commands:"
	@awk 'BEGIN {FS = ":.*##"; printf "\n\033[1m%-12s\033[0m %s\n\n", "Command", "Description"} /^[a-zA-Z_-]+:.*?##/ { printf "\033[36m%-12s\033[0m %s\n", $$1, $$2 }' $(MAKEFILE_LIST)

.PHONY: crtmgct strmgct stpmgct rmvmgct crtmgdb drpmgdb build_win build_lin build_mac build test lint clean tr help
