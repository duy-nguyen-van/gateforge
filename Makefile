.PHONY: dev dev-backend dev-frontend bootstrap build-frontend copy-frontend build-prod docker-build security security-fs security-image

dev-backend:
	cd backend && make up

dev-frontend:
	cd frontend && npm run dev

dev: bootstrap

bootstrap:
	@bash -c '\
	set -euo pipefail; \
	echo "→ Starting Postgres & Redis..."; \
	$(MAKE) -C backend container-up; \
	echo "✓ Postgres & Redis running"; \
	echo "→ Applying migrations..."; \
	$(MAKE) -C backend migrate-up; \
	echo "✓ Migrations applied"; \
	$(MAKE) -C frontend setup; \
	(cd frontend && npm run dev) & \
	frontend_pid=$$!; \
	trap "kill $$frontend_pid 2>/dev/null || true; wait $$frontend_pid 2>/dev/null || true" EXIT INT TERM; \
	sleep 2; \
	echo "✓ Admin UI at http://localhost:5173"; \
	echo "✓ API ready on :3000"; \
	echo ""; \
	$(MAKE) -C backend up'

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

# Match .github/workflows/ci.yml security job (requires: brew install trivy)
TRIVY_FLAGS = --format table --exit-code 1 --ignore-unfixed --vuln-type os,library --severity CRITICAL,HIGH

security-fs:
	trivy fs . $(TRIVY_FLAGS) --scanners vuln,secret,misconfig

security-image: docker-build
	trivy image gateforge-iam:latest $(TRIVY_FLAGS)

security: security-fs security-image
