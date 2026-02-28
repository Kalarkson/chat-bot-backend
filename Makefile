DC = docker compose

.PHONY: help up down logs status sd-logs

help:
	@echo "Commands:"
	@echo "  make up     - Start all services (backend + mongo + sd)"
	@echo "  make down   - Stop all services"
	@echo "  make logs   - Show logs from all services"
	@echo "  make status - Show status of all containers"
	@echo "  make sd-logs - Show only SD logs"

up:
	$(DC) up -d

down:
	$(DC) down

logs:
	$(DC) logs -f

status:
	$(DC) ps

sd-logs:
	$(DC) logs -f sd-webui