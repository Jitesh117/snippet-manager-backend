BINARY_NAME=snippet-manager
DOCKER_POSTGRES_NAME=postgres-snippet-manager
DB_NAME=snippet_manager
DB_USER=postgres
DB_PASSWORD=mysecretpassword
DB_PORT=5432
GOBASE=$(shell pwd)
GOBIN=$(GOBASE)/bin

build:
	@echo "Building..."
	@go build -o $(BINARY_NAME)

run: build
	@echo "Running..."
	@./$(BINARY_NAME)


db-start:
	@echo "Starting PostgreSQL container..."
	@docker start $(DOCKER_POSTGRES_NAME) 2>/dev/null || \
		docker run --name $(DOCKER_POSTGRES_NAME) \
		-e POSTGRES_PASSWORD=$(DB_PASSWORD) \
		-p $(DB_PORT):5432 \
		-v $(DB_VOLUME):/var/lib/postgresql/data \
		-d postgres

db-create:
	@echo "Creating database..."
	@docker exec -it $(DOCKER_POSTGRES_NAME) createdb --username=$(DB_USER) --owner=$(DB_USER) $(DB_NAME)

db-drop:
	@echo "Dropping database..."
	@docker exec -it $(DOCKER_POSTGRES_NAME) dropdb --username=$(DB_USER) $(DB_NAME)

db-stop:
	@echo "Stopping PostgreSQL container..."
	@docker stop $(DOCKER_POSTGRES_NAME)

db-remove: db-stop
	@echo "Removing PostgreSQL container..."
	@docker rm $(DOCKER_POSTGRES_NAME)

db-restart: db-stop db-start

help:
	@echo "Available commands:"
	@echo "  make build      - Build the project"
	@echo "  make run        - Run the project"
	@echo "  make clean      - Clean the binary"
	@echo "  make db-start   - Start PostgreSQL container"
	@echo "  make db-stop    - Stop PostgreSQL container"
	@echo "  make db-remove  - Remove PostgreSQL container"
	@echo "  make db-create  - Create database"
	@echo "  make db-drop    - Drop database"
	@echo "  make db-restart - Restart PostgreSQL container"

