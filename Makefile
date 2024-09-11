# Environment variables for MongoDB
MONGODB_CONTAINER_NAME ?= HTMX_GO
MONGODB_IMAGE_TAG ?= latest
MONGODB_DB_NAME ?= htmx_go
MONGODB_PORT ?= 27017

# Colors for help command
CYAN := \033[36m
RESET := \033[0m

# Docker MongoDB commands
crtmgdb: ## Create and start the MongoDB container
	@echo "Creating and starting MongoDB container..."
	docker run --name $(MONGODB_CONTAINER_NAME) -p $(MONGODB_PORT):27017 -d mongo:$(MONGODB_IMAGE_TAG)

strmgdb: ## Start the MongoDB container
	@echo "Starting MongoDB container..."
	docker start $(MONGODB_CONTAINER_NAME)

stpmgdb: ## Stop the MongoDB container
	@echo "Stopping MongoDB container..."
	docker stop $(MONGODB_CONTAINER_NAME)

rmvmgdb: ## Remove the MongoDB container
	@echo "Removing MongoDB container..."
	docker rm $(MONGODB_CONTAINER_NAME)

# MongoDB database commands
createdb_mongodb: strmgdb ## Create MongoDB database
	@echo "Creating MongoDB database..."
	docker exec -it $(MONGODB_CONTAINER_NAME) mongosh --eval "use $(MONGODB_DB_NAME)"

dropdb_mongodb: strmgdb ## Drop MongoDB database
	@echo "Dropping MongoDB database..."
	docker exec -it $(MONGODB_CONTAINER_NAME) mongosh --eval "db.getSiblingDB('$(MONGODB_DB_NAME)').dropDatabase()"


# Help command
help: ## Show this help message
	@echo "Available commands:"
	@awk 'BEGIN {FS = ":.*##"; printf "\n\033[1m%-12s\033[0m %s\n\n", "Command", "Description"} /^[a-zA-Z_-]+:.*?##/ { printf "\033[36m%-12s\033[0m %s\n", $$1, $$2 }' $(MAKEFILE_LIST)

.PHONY: crtmgdb strmgdb stpmgdb rmvmgdb createdb_mongodb help