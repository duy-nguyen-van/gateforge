.PHONY: dev dev-backend dev-frontend bootstrap build-frontend copy-frontend build-prod docker-build

dev-backend:
	cd backend && make up

dev-frontend:
	cd frontend && npm run dev

dev:
	@echo "Run in two terminals:"
	@echo "  make dev-backend   # Go API on :3000"
	@echo "  make dev-frontend  # Vite SPA on :5173"

bootstrap:
	cd backend && make bootstrap

build-frontend:
	cd frontend && npm ci && npm run build

copy-frontend:
	rm -rf backend/internal/static/dist
	mkdir -p backend/internal/static/dist
	cp -r frontend/dist/. backend/internal/static/dist/

build-prod: build-frontend copy-frontend
	cd backend && CGO_ENABLED=0 go build -tags embedfrontend -o ../bin/gateforge-iam-server ./cmd/server

docker-build:
	DOCKER_BUILDKIT=1 docker build -f docker/Dockerfile -t gateforge-iam:latest .
