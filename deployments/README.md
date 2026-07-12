# GateForge IAM Deployments

Production options for the GateForge IAM hybrid architecture (single Go binary with embedded SPA).

## Option 1: Production binary (recommended)

From the repository root:

```bash
make build-prod
```

Run with environment configured for same-origin deploy:

```bash
export APP_ENV=production
export SERVE_EMBEDDED_FRONTEND=true
export APP_BASE_URL=https://iam.example.com
export OIDC_LOGIN_PAGE_URL=https://iam.example.com/login
export WEBAUTHN_RP_ORIGINS=https://iam.example.com
# ... database, redis, JWT, OIDC keys, etc.

./bin/gateforge-iam-server
```

The server listens on `APP_HTTP_SERVER` (default `:3000`) and serves:

- SPA UI at `/`
- API at `/api/v1/*`
- OIDC at `/authorize`, `/token`, `/oidc/*`, `/.well-known/*`

## Option 2: Docker Compose

Copy and configure `backend/cmd/server/.env`, then:

```bash
docker compose -f docker/docker-compose.prod.yml up -d --build
```

This builds the multi-stage image from [`docker/Dockerfile`](../docker/Dockerfile), runs migrations, and starts Postgres, Redis, and the app.

## Option 3: Split static + API (alternative)

If you prefer a reverse proxy instead of the embedded binary, see [`frontend/deploy/Caddyfile`](../frontend/deploy/Caddyfile) and [`frontend/README.md`](../frontend/README.md).

## systemd template

See [`systemd/gateforge-iam.service`](systemd/gateforge-iam.service) for a minimal unit file.
