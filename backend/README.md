# GateForge IAM Backend

Native identity provider built with **Go**, **Echo**, **Uber FX**, **PostgreSQL**, **Redis**, and **Atlas** migrations. Provides OIDC/OAuth2, browser SSO sessions, passkeys (WebAuthn), TOTP MFA, Google federation, multi-tenant memberships, and a platform admin API.

Go module: `github.com/gateforge-iam/gateforge-iam`

## Features

- **OIDC IdP** — Authorization code + PKCE, JWKS, userinfo (`/authorize`, `/token`, `/userinfo`)
- **SSO session** — Shared `iam_session` cookie between API login and OIDC authorize
- **App auth** — Register, login, refresh, JWT access tokens with active `tenant_id`
- **WebAuthn passkeys** — Register and passwordless login
- **MFA** — TOTP enrollment and step-up login with recovery codes
- **Federation** — Google upstream OAuth (more IdPs planned)
- **Multi-tenant** — Global users, `tenant_memberships`, tenant select/switch
- **Platform admin** — Console APIs protected by `users.is_platform_admin`
- **Observability** — Zap logging, Sentry, New Relic (optional)

## Documentation

Start at **[docs/README.md](docs/README.md)** — feature index, database tables, Redis keys, route cheat sheet.

| Type | Path |
|------|------|
| Feature guides | [docs/features/](docs/features/) — OIDC, SSO, federation, passkey/MFA, multi-tenant, authorization |
| Manual testing | [docs/testing/](docs/testing/) — curl and Postman flows |
| OpenAPI (`/api/v1`) | [docs/swagger.yaml](docs/swagger.yaml) |
| Agent / layer map | [AGENTS.md](AGENTS.md) |

## Quick start

### Prerequisites

- Go 1.26+
- Docker (Postgres + Redis)
- [Atlas CLI](https://atlasgo.io/) for migrations

### Setup

```bash
cp examples/env/server.env.example cmd/server/.env
make container-up
make migrate-up
make up                    # server on :3000
```

Or from monorepo root: `make bootstrap` (Postgres, Redis, migrations, API, and Admin UI).

### Make commands

| Command | Purpose |
|---------|---------|
| `make up` | Run API server |
| `make container-up` / `container-down` | Postgres + Redis |
| `make migrate-up` | Apply Atlas migrations |
| `make migrate-hash` | Refresh `atlas.sum` after new SQL |
| `make lint` | golangci-lint |
| `make tests` | Full test suite |
| `make swagger-load` | Regenerate OpenAPI from swag annotations |
| `make build-prod` | Production binary (requires embedded frontend) |

## Project structure

```text
backend/
├── cmd/
│   ├── server/           # FX wiring, routes, main
│   └── migrations/sql/   # Atlas SQL migrations
├── docs/
│   ├── README.md         # Documentation hub
│   ├── features/         # Feature guides
│   ├── testing/          # curl / Postman
│   └── swagger.yaml      # OpenAPI (/api/v1)
├── internal/
│   ├── handlers/         # Echo HTTP handlers
│   ├── services/         # Business logic
│   ├── repositories/     # GORM data access
│   ├── models/           # DB entities
│   ├── auth/             # JWT, OIDC signing, ephemeral Redis
│   ├── middlewares/      # Auth, CSRF, CORS, rate limit
│   └── ...
└── pkg/                  # Shared libraries (e.g. correlationid)
```

Layer boundaries and import rules: [AGENTS.md](AGENTS.md), `.cursor/rules/go-layering.mdc`.

## API endpoints

Two URL namespaces — see [docs/README.md](docs/README.md#route-surfaces).

### OIDC (root)

| Method | Path |
|--------|------|
| GET | `/.well-known/openid-configuration`, `/.well-known/jwks.json` |
| GET | `/authorize` |
| POST | `/oidc/login`, `/token` |
| GET | `/userinfo` |
| GET | `/oidc/federation/:provider/start`, `/oidc/federation/:provider/callback` |

### App API (`/api/v1`)

| Method | Path | Auth |
|--------|------|------|
| POST | `/register`, `/login`, `/refresh` | Public |
| POST | `/logout`, `/me`, `/me/tenants` | JWT |
| POST | `/tenants/select`, `/tenants/switch` | Public / JWT |
| POST | `/webauthn/register/*`, `/webauthn/login/*` | JWT / Public |
| POST | `/mfa/totp/*`, `/mfa/recovery-codes`, `/mfa/challenge/verify` | JWT / Public |
| GET/PATCH | `/admin/*` | JWT + platform admin |
| PATCH | `/internal/tenants/:tenantId/identity-providers/google` | `X-Admin-API-Key` |

Full request/response shapes: `docs/swagger.yaml` (App API only).

## Example usage

Register and login:

```bash
BASE=http://localhost:3000

curl -sS -X POST "$BASE/api/v1/register" \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"long-secure-passphrase","first_name":"Test","last_name":"User"}'

curl -sS -X POST "$BASE/api/v1/login" \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"long-secure-passphrase"}'
```

OIDC authorization code flow: [docs/testing/OIDC_CURL.md](docs/testing/OIDC_CURL.md).

## Configuration

Environment: `cmd/server/.env.example`. Key variables:

| Variable | Purpose |
|----------|---------|
| `JWT_SECRET` | Dashboard JWT signing (min 32 chars) |
| `APP_BASE_URL` | Issuer, redirect validation |
| `OIDC_RSA_PRIVATE_KEY_PEM` | OIDC token signing |
| `WEBAUTHN_RP_ORIGINS` | Allowed WebAuthn origins |
| `GOOGLE_FEDERATION_ENABLED` | Upstream Google login |
| `BOOTSTRAP_ADMIN_EMAIL` | First-run platform admin |

## Development

- Run `make lint` and `make tests` before pushing
- After `/api/v1` handler changes: update swag annotations and `make swagger-load`
- After OIDC/session/MFA behavior changes: update [docs/features/](docs/features/) and [docs/testing/](docs/testing/)
- After schema changes: new SQL in `cmd/migrations/sql/` + `make migrate-hash`

See [CONTRIBUTING.md](../CONTRIBUTING.md) for contribution guidelines.

## Testing

```bash
make tests                    # all packages
make test-handlers            # handlers only
make test-services            # services only
go test -race ./internal/...  # race detector
```

Manual IAM flows: [docs/testing/](docs/testing/).

## Docker

**Local dependencies** (Postgres + Redis):

```bash
make container-up    # docker-compose.yml in this directory
```

**Production images** — two Dockerfiles exist; use the right one for your deploy model:

| Dockerfile | Purpose |
|------------|---------|
| [`../docker/Dockerfile`](../docker/Dockerfile) | **Canonical prod** — builds frontend, embeds SPA into the Go binary, single hybrid image. Run from repo root: `make docker-build`. |
| [`Dockerfile`](Dockerfile) | **API-only (legacy)** — Go binary without embedded SPA. Use when the UI is served separately (CDN, Caddy, etc.). |

Production deployment runbooks: [deployments/README.md](../deployments/README.md).
