.PHONY: help dev-backend dev-frontend install-backend install-frontend

help:
	@echo "Available commands:"
	@echo "  make dev-backend     - Run Java Spring Boot backend"
	@echo "  make dev-frontend    - Run frontend dev server"
	@echo "  make install-backend - Resolve Maven dependencies"
	@echo "  make install-frontend - Install npm dependencies"

dev-backend:
	cd backend && mvn spring-boot:run

dev-frontend:
	cd frontend && npm run dev

install-backend:
	cd backend && mvn -q dependency:resolve

install-frontend:
	cd frontend && npm install
