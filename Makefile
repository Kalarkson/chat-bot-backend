DC = docker compose

.PHONY: help up down restart logs status build sd-shell sd-logs

help:
	@echo "Commands:"
	@echo "  make up        - Start all services"
	@echo "  make down      - Stop all services"
	@echo "  make logs      - Show logs"
	@echo "  make status    - Show status"
	@echo "  make sd-shell   - Open shell in SD container"
	@echo "  make sd-logs    - Show SD logs"

up:
	$(DC) up -d

down:
	$(DC) down

logs:
	$(DC) logs -f

status:
	$(DC) ps

sd-shell:
	$(DC) exec sd-webui /bin/sh

sd-logs:
	$(DC) logs -f sd-webui