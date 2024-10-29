# Makefile

# Default Docker Compose file
DOCKER_COMPOSE_FILE := docker-compose.yml

# Up services
up:
	docker-compose -f $(DOCKER_COMPOSE_FILE) up -d

# Down services
down:
	docker-compose -f $(DOCKER_COMPOSE_FILE) down

# Restart services
restart:
	docker-compose -f $(DOCKER_COMPOSE_FILE) down
	docker-compose -f $(DOCKER_COMPOSE_FILE) up -d

# Build services without starting
build:
	docker-compose -f $(DOCKER_COMPOSE_FILE) build

# Check logs of services
logs:
	docker-compose -f $(DOCKER_COMPOSE_FILE) logs -f

delete:
	docker-compose -f $(DOCKER_COMPOSE_FILE) down --volumes --remove-orphans


rebuild:
	docker-compose -f $(DOCKER_COMPOSE_FILE) up --build -d

rebuild-verbose:
	docker-compose -f $(DOCKER_COMPOSE_FILE) up --build --progress=plain -d



.PHONY: up down restart build logs
