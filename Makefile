.PHONY: help up down db-up db-down migrate-up migrate-down backend frontend logs

help: ## Muestra esta ayuda
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

up: ## Levanta todos los servicios (db, backend, frontend)
	docker compose up -d

down: ## Detiene todos los servicios
	docker compose down

db-up: ## Levanta solo PostgreSQL
	docker compose up -d db

db-down: ## Detiene PostgreSQL
	docker compose stop db

migrate-up: ## Ejecuta migraciones hacia arriba
	docker compose exec backend migrate -path /app/migrations -database "$$DATABASE_URL" up

migrate-down: ## Revierte la última migración
	docker compose exec backend migrate -path /app/migrations -database "$$DATABASE_URL" down 1

backend: ## Construye y levanta solo el backend
	docker compose up -d --build backend

frontend: ## Construye y levanta solo el frontend
	docker compose up -d --build frontend

logs: ## Muestra logs de todos los servicios
	docker compose logs -f

backend-logs: ## Muestra logs del backend
	docker compose logs -f backend

frontend-logs: ## Muestra logs del frontend
	docker compose logs -f frontend

db-logs: ## Muestra logs de PostgreSQL
	docker compose logs -f db

db-shell: ## Abre psql en el contenedor de PostgreSQL
	docker compose exec db psql -U sales_insight -d sales_insight
