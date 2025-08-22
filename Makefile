.PHONY: clean build run stop logs dev

DOCKER_COMPOSE := docker compose

build:
	@echo "Construindo aplicação..."
	@$(DOCKER_COMPOSE) build

run:
	@echo "Iniciando aplicação..."
	@$(DOCKER_COMPOSE) up -d

stop:
	@echo "Parando aplicação..."
	@$(DOCKER_COMPOSE) down

logs:
	@echo "Mostrando logs..."
	@$(DOCKER_COMPOSE) logs -f

dev:
	@echo "Modo desenvolvimento (logs em tempo real)..."
	@$(DOCKER_COMPOSE) up