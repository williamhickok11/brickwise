.PHONY: help dev-backend dev-frontend install-backend install-frontend

help:
	@echo "Available commands:"
	@echo "  make dev-backend    - Run backend server"
	@echo "  make dev-frontend   - Run frontend dev server"
	@echo "  make install-backend - Install Go dependencies"
	@echo "  make install-frontend - Install npm dependencies"

dev-backend:
	cd backend && go run cmd/server/main.go

dev-frontend:
	cd frontend && npm run dev

install-backend:
	cd backend && go mod download

install-frontend:
	cd frontend && npm install
