# GateForge IAM

![GateForge — The sovereign identity layer](assets/gateforge-logo.png)

> **The sovereign identity layer** — a self-hosted OIDC/OAuth2 identity provider you fully control.

GateForge IAM gives your applications one secure place for user sign-in and access control. Connect apps through open standards, offer modern authentication, and manage organizations from a web console — all shipped as a single production binary.

## Features

- **OIDC / OAuth2 provider** — authorization code + PKCE, JWKS, userinfo, and browser SSO
- **Flexible sign-in** — passwords, passkeys (WebAuthn), and Google federation
- **Multi-factor auth** — TOTP enrollment, step-up login, and recovery codes
- **Multi-tenant** — organizations, memberships, roles, and per-tenant identity providers
- **Admin console** — manage users, clients, tenants, IdPs, and security activity
- **Simple deploy** — one Go binary with the React SPA embedded, or a Docker image

## Why GateForge?

- **Own your identity layer** — self-host authentication and keep control of user data, credentials, and deployment.
- **Integrate with open standards** — connect applications through OIDC/OAuth2 instead of locking into a vendor IdP.
- **Offer secure, flexible sign-in** — passwords, passkeys, Google, and MFA without bolting on extra services.
- **Manage multiple organizations** — separate tenants, roles, clients, and IdPs from one console.
- **Deploy with less complexity** — API and admin UI ship together as one binary or container.

## Quick start

Prerequisites: Go 1.26+, Node.js 26+, Docker (for Postgres/Redis).

```bash
# Start Postgres, Redis, migrations, API (:3000), and Admin UI (:5173)
make bootstrap

# Or run backend and frontend separately in two terminals:
make dev-backend    # :3000
make dev-frontend   # :5173 (proxies API/OIDC to backend)
```

`make bootstrap` (alias: `make dev`) runs the full dev stack in one terminal. Press Ctrl+C to stop the API and Admin UI; Postgres and Redis keep running in Docker until `make -C backend container-down`.

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

## Repository layout

```
gateforge-iam/
├── backend/       Go 1.26 + Echo — auth, OIDC, WebAuthn, MFA, admin APIs
├── frontend/      Vite + React 19 SPA
├── docker/        Multi-stage production Dockerfile and compose
├── deployments/   Production runbooks and systemd template
├── performance/   k6 benches and capacity methodology
└── Makefile       Monorepo dev and build targets
```

## Security scanning (Trivy)

CI runs Trivy before the production build. Run the same checks locally before pushing.

### Install Trivy

macOS (Homebrew):

```bash
brew install trivy
trivy --version
```

Linux (official install script):

```bash
curl -sfL https://raw.githubusercontent.com/aquasecurity/trivy/main/contrib/install.sh | sh -s -- -b /usr/local/bin
```

Without a local install, use the Docker image:

```bash
docker pull aquasec/trivy:latest
alias trivy='docker run --rm -v /var/run/docker.sock:/var/run/docker.sock -v "$PWD":/src aquasec/trivy'
```

### Run scans

From the repo root:

```bash
make security       # filesystem + Docker image (same as CI)
make security-fs    # repo scan only (faster)
make security-image # build image, then scan
```

Trivy downloads its vulnerability database on first run. Refresh manually if needed:

```bash
trivy image --download-db-only
```

## Production build (hybrid)

Production requires a persistent OIDC RSA signing key. Development auto-generates an ephemeral key when unset.

Generate a 2048-bit key and copy the PEM into `backend/cmd/server/.env`:

```bash
openssl genrsa 2048
```

```
OIDC_RSA_PRIVATE_KEY_PEM="-----BEGIN RSA PRIVATE KEY-----
...
-----END RSA PRIVATE KEY-----"
```

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
- Performance: [`performance/README.md`](performance/README.md) — how to run benches and publish numbers

## Contributing

- [CONTRIBUTING.md](CONTRIBUTING.md) — development workflow, PR expectations
- [SECURITY.md](SECURITY.md) — vulnerability reporting
- [CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md) — community standards
