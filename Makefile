.PHONY: dev dev-backend dev-frontend bootstrap build-frontend copy-frontend build-prod docker-build security security-fs security-image performance-help performance-smoke performance-token performance-token-hold performance-passkey performance-oidc-e2e performance-rss sdk-generate sdk-check sdk-go-test sdk-ts-build sdk-ts-test

OPENAPI_GENERATOR_VERSION ?= v7.19.0
OPENAPI_GENERATOR_IMAGE ?= openapitools/openapi-generator-cli:$(OPENAPI_GENERATOR_VERSION)

dev-backend:
	cd backend && make up

dev-frontend:
	$(MAKE) sdk-ts-build
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
	$(MAKE) sdk-ts-build
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

# --- Marketing / capacity benches (see performance/README.md) ---
performance-help:
	$(MAKE) -C performance help

performance-smoke:
	$(MAKE) -C performance smoke

performance-token:
	$(MAKE) -C performance token

performance-token-hold:
	$(MAKE) -C performance token-hold

performance-passkey:
	$(MAKE) -C performance passkey

performance-oidc-e2e:
	$(MAKE) -C performance oidc-e2e

performance-rss:
	$(MAKE) -C performance rss

# --- Official SDKs (see sdk/README.md; requires api/openapi.yaml) ---
sdk-generate:
	docker run --rm -v "$(CURDIR):/local" $(OPENAPI_GENERATOR_IMAGE) generate \
		-i /local/api/openapi.yaml \
		-g go \
		-c /local/sdk/config/go.yaml \
		-o /local/sdk/go/openapi \
		--git-user-id gateforge-iam \
		--git-repo-id 'gateforge-iam/sdk/go' \
		--additional-properties=packageName=openapi,packageVersion=0.1.0,withGoMod=true,enumClassPrefix=true,generateInterfaces=true,isGoSubmodule=true
	docker run --rm -v "$(CURDIR):/local" $(OPENAPI_GENERATOR_IMAGE) generate \
		-i /local/api/openapi.yaml \
		-g typescript-fetch \
		-c /local/sdk/config/typescript.yaml \
		-o /local/sdk/typescript/src/generated \
		--additional-properties=npmName=@gateforge/sdk,npmVersion=0.1.0,supportsES6=true,typescriptThreePlus=true,modelPropertyNaming=original,enumPropertyNaming=original,useSingleRequestParameter=true
	@rm -rf sdk/go/openapi/test
	@rm -f \
		sdk/go/openapi/.travis.yml \
		sdk/go/openapi/.gitignore \
		sdk/go/openapi/git_push.sh \
		sdk/typescript/src/generated/.travis.yml \
		sdk/typescript/src/generated/.gitignore \
		sdk/typescript/src/generated/git_push.sh \
		sdk/typescript/src/generated/package.json \
		sdk/typescript/src/generated/tsconfig.json \
		sdk/typescript/src/generated/tsconfig.esm.json \
		sdk/typescript/src/generated/.npmignore \
		2>/dev/null || true
	@printf '%s\n' "export * from './src/index.js'" > sdk/typescript/src/generated/index.ts
	node sdk/typescript/scripts/patch-generated.mjs

sdk-check: sdk-generate
	@git diff --exit-code -- api/openapi.yaml sdk/go/openapi sdk/typescript/src/generated || (echo "SDK drift detected; commit regenerated files" && exit 1)

sdk-go-test:
	cd sdk/go && go test ./...

sdk-ts-build:
	cd sdk/typescript && npm ci && npm run build

sdk-ts-test:
	cd sdk/typescript && npm ci && npm test
