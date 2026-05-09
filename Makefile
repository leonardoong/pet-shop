.PHONY: up down build logs migrate dev-backend dev-frontend dev-admin

up:
	docker compose up -d

down:
	docker compose down

build:
	docker compose build

logs:
	docker compose logs -f

# Run migrations inside backend container
migrate:
	docker compose exec backend ./petshop migrate

# Local dev (no Docker)
dev-backend:
	cd backend && go run ./cmd/api

dev-frontend:
	cd frontend && npm run dev

dev-admin:
	cd admin && npm run dev

# Install frontend deps
install:
	cd frontend && npm install
	cd admin && npm install

# Tidy Go deps
tidy:
	cd backend && go mod tidy

# Copy env
env:
	cp .env.example .env
