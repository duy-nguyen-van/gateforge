ifneq (,$(wildcard .env))
    include .env
    export $(shell sed -n 's/^\([A-Za-z_][A-Za-z0-9_]*\)=.*/\1/p' .env)
endif

NPM ?= npm
PORT ?= 5173
PREVIEW_PORT ?= 4173

.DEFAULT_GOAL := help

.PHONY: help install setup bootstrap run up dev build lint preview clean generate-api check

help: ## Show available targets
	@awk 'BEGIN {FS = ":.*##"; printf "Usage: make <target>\n\nTargets:\n"} /^[a-zA-Z0-9_-]+:.*##/ {printf "  \033[36m%-16s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

install: ## Install npm dependencies
	$(NPM) install

setup: ## Copy .env.example to .env (if missing) and install dependencies
	@test -f .env || cp .env.example .env
	$(MAKE) install

bootstrap: ## First-time setup and start the dev server
	$(MAKE) setup
	$(MAKE) run

run: ## Start Vite dev server (default http://localhost:5173)
	$(NPM) run dev

up: run ## Alias for run (matches backend)

dev: run ## Alias for run

build: ## Type-check and build production assets to dist/
	$(NPM) run build

lint: ## Run ESLint
	$(NPM) run lint

preview: ## Preview the production build locally
	$(NPM) run preview -- --port $(PREVIEW_PORT)

check: lint build ## Run lint and production build

clean: ## Remove build output and TypeScript cache
	rm -rf dist
	rm -rf node_modules/.tmp
	rm -rf node_modules/.vite

generate-api: ## Regenerate API types from backend Swagger (optional)
	$(NPM) install -D openapi-typescript --legacy-peer-deps
	$(NPM) run generate:api
