# Agent context (monorepo)

GateForge IAM monorepo: Go API + React SPA, shipped as a single embedded binary in production.

## Repository layout

| Path | Responsibility |
|------|----------------|
| [`backend/`](backend/) | Go IAM service — API, OIDC, DB, migrations |
| [`frontend/`](frontend/) | React SPA — login, admin console |
| [`docker/`](docker/) | Production multi-stage Dockerfile and compose |
| [`deployments/`](deployments/) | Runbooks, systemd, Caddy split-deploy config |
| [`Makefile`](Makefile) | Dev and prod build entry points |
| [`.github/workflows/ci.yml`](.github/workflows/ci.yml) | CI for backend, frontend, and prod binary |

Per-package agent rules live in `backend/.cursor/` and `frontend/.cursor/`. Read those when working inside a single package.

## Development

```bash
make bootstrap        # Postgres, Redis, migrations, API (:3000), Admin UI (:5173)
make dev              # alias for bootstrap
# Or in two terminals:
make dev-backend      # Go API on :3000
make dev-frontend     # Vite SPA on :5173 (proxies API/OIDC to backend)
```

Env templates: `backend/cmd/server/.env.example`, `frontend/.env.example`.

## Production

```bash
make build-prod       # frontend build → embed → bin/gateforge-iam-server
make docker-build     # hybrid image via docker/Dockerfile
```

See [`deployments/README.md`](deployments/README.md) for systemd, Docker Compose, and Caddy split-deploy options.

## Documentation

| Doc | Purpose |
|-----|---------|
| [`README.md`](README.md) | Monorepo quick start |
| [`backend/AGENTS.md`](backend/AGENTS.md) | Backend layer map, routes, patterns |
| [`frontend/AGENTS.md`](frontend/AGENTS.md) | Frontend routes, auth, API client |
| [`backend/docs/README.md`](backend/docs/README.md) | IAM feature hub — OIDC, SSO, MFA, DB tables |
