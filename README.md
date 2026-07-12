# GateForge IAM

Identity and access management platform with a Go API and React admin console, deployed as a single binary in production.

## Repository layout

```
gateforge-iam/
├── backend/      Go 1.26 + Echo — auth, OIDC, WebAuthn, MFA, admin APIs
├── frontend/     Vite + React 19 SPA
├── docker/       Multi-stage production Dockerfile and compose
├── deployments/  Production runbooks and systemd template
└── Makefile      Monorepo dev and build targets
```

## Development

Prerequisites: Go 1.26+, Node.js 20+, Docker (for Postgres/Redis).

```bash
# Start Postgres, Redis, run migrations
make bootstrap

# Or run backend and frontend separately in two terminals:
make dev-backend    # :3000
make dev-frontend   # :5173 (proxies API/OIDC to backend)
```

Copy `backend/cmd/server/.env.example` to `backend/cmd/server/.env` and `frontend/.env.example` to `frontend/.env`.

### Test hybrid UI on :3000 (optional)

By default, dev uses Vite on `:5173` + API on `:3000`. To serve the built SPA from the Go server locally:

```bash
make copy-frontend   # copy frontend/dist → backend/internal/static/dist
```

In `backend/cmd/server/.env`:

```
SERVE_EMBEDDED_FRONTEND=true
OIDC_LOGIN_PAGE_URL=http://localhost:3000/login
WEBAUTHN_RP_ORIGINS=http://localhost:3000
```

Then run/debug the backend. Without `-tags embedfrontend`, assets are read from disk automatically. Re-run `make copy-frontend` after frontend changes.

## Production build (hybrid)

Build a single `gateforge-iam-server` binary with the frontend embedded via `go:embed`:

```bash
make build-prod
./bin/gateforge-iam-server
```

Or build a Docker image:

```bash
make docker-build
```

See [`deployments/README.md`](deployments/README.md) for environment variables and deployment options.

## Documentation

- IAM feature hub: [`backend/docs/README.md`](backend/docs/README.md) — OIDC, SSO, MFA, federation, DB tables, curl testing
- Backend: [`backend/README.md`](backend/README.md)
- Frontend: [`frontend/README.md`](frontend/README.md)
- Deploy: [`deployments/README.md`](deployments/README.md)
