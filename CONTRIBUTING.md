# Contributing

Thank you for your interest in contributing. GateForge IAM is a monorepo:

| Path | What it is |
|------|------------|
| [`backend/`](backend/) | Go IAM API — OIDC, sessions, WebAuthn, MFA, federation, Postgres/Redis |
| [`frontend/`](frontend/) | React admin SPA — login, security settings, platform console |
| [`docker/`](docker/), [`deployments/`](deployments/) | Production image, compose, systemd, Caddy runbooks |
| Root [`Makefile`](Makefile) | Dev bootstrap, `build-prod`, Trivy security scans |

Production ships as a **single hybrid binary** (`bin/gateforge-iam-server`) with the SPA embedded via `go:embed`.

## Getting started

1. Fork the repository and branch from `main` (see [Branch naming](#branch-naming) below).
2. Install prerequisites:
   - **Go 1.26+**
   - **Node.js 26+**
   - **Docker** (Postgres + Redis for local dev)
   - **Make**
   - **[Atlas CLI](https://atlasgo.io/)** — required for database migrations in `backend/`
3. Copy env templates (never commit the resulting `.env` files):
   - `backend/cmd/server/.env.example` → `backend/cmd/server/.env`
   - `frontend/.env.example` → `frontend/.env`
4. Start the full dev stack from the repo root:

```bash
make bootstrap   # Postgres, Redis, migrations, API (:3000), Admin UI (:5173)
```

Or run services separately:

```bash
make dev-backend    # API on :3000
make dev-frontend   # Vite on :5173 (proxies API/OIDC to backend)
```

Press Ctrl+C to stop the API and Admin UI; Postgres and Redis keep running until `make -C backend container-down`.

See [`README.md`](README.md) for hybrid UI on `:3000`, OIDC key setup, and production builds.

## Development workflow

### Backend (`backend/`)

```bash
make -C backend lint          # golangci-lint
make -C backend tests         # all packages, race + coverage
make -C backend dep           # go mod tidy — run after dependency edits
```

Targeted test targets: `test-handlers`, `test-services`, `test-repositories` (see `backend/Makefile`).

**When you change `/api/v1` HTTP handlers or DTOs:**

1. Update swag annotations on handlers and `cmd/server/main.go`.
2. Run `make -C backend swagger-load`.
3. Commit `backend/docs/swagger.yaml`, `swagger.json`, and `docs.go` with the code change.

**When you change OIDC root routes** (`/authorize`, `/token`, `/userinfo`, `/.well-known/*`, `/oidc/*`): update the matching guides under [`backend/docs/features/`](backend/docs/features/) and [`backend/docs/testing/`](backend/docs/testing/) — those routes are outside the generated OpenAPI spec.

**When you change auth behavior** (OIDC, SSO sessions, federation, WebAuthn, MFA, tenants, platform admin): read [`backend/docs/README.md`](backend/docs/README.md) first, then update the relevant feature and testing docs in the same PR.

**When you change the database schema:**

1. Add SQL under `backend/cmd/migrations/sql/`.
2. Run `make -C backend migrate-hash`.
3. Update [`backend/docs/README.md`](backend/docs/README.md) if tables or Redis keys change.

### Frontend (`frontend/`)

```bash
make -C frontend setup    # .env + npm install (first time)
make -C frontend check    # lint + production build
```

Or individually: `make -C frontend lint`, `make -C frontend build`.

Update `frontend/src/api/types.ts` when backend `/api/v1` response shapes change.

### Root tooling

If you touch the embed pipeline, Docker image, or root `Makefile`:

```bash
make build-prod    # frontend build → copy → hybrid binary
make docker-build  # production image via docker/Dockerfile
```

### Security scanning

CI runs Trivy on the filesystem and production Docker image before `build-prod`. From the repo root:

```bash
make security-fs    # filesystem only (faster)
make security       # filesystem + image (matches CI)
```

See [`README.md`](README.md#security-scanning-trivy) for installing Trivy locally.

## Commit messages

Use [Conventional Commits](https://www.conventionalcommits.org/):

```
<type>(<scope>): <short description>
```

| Type | Use for |
|------|---------|
| `feat` | New feature |
| `fix` | Bug fix |
| `docs` | Documentation only |
| `refactor` | Restructure without behavior change |
| `test` | Add or update tests |
| `chore` | Build, CI, tooling |
| `perf` | Performance improvement |

Examples: `feat(oidc): add refresh token rotation`, `fix(session): handle expired cookie`.

Rules: lowercase subject, imperative mood ("add" not "added"), no trailing period, keep the subject under 72 characters. Reference issues in the body when relevant (`Fixes #123`).

## Branch naming

- `feature/<short-description>` — new work
- `fix/<short-description>` — bug fixes
- `chore/<short-description>` — tooling or docs

## Pull requests

1. Fill out [`.github/pull_request_template.md`](.github/pull_request_template.md).
2. Keep PRs focused — one concern per PR.
3. Title should follow conventional commit format.
4. Describe what changed, why, and how to test.
5. Include screenshots for UI changes; include curl output or logs for API/auth flow changes when helpful.

**CI must pass** for the areas you changed:

| Area | Checks |
|------|--------|
| `backend/` | `make -C backend lint`, `go build ./...`, `make -C backend tests` |
| `frontend/` | `npm run lint`, `npm run build` |
| Embed / Docker / root Makefile | `make build-prod` |

## Code style

- **Go**: idiomatic style, explicit error handling, layer boundaries per [`backend/AGENTS.md`](backend/AGENTS.md). Do not commit secrets or production URLs.
- **TypeScript/React**: match patterns in `frontend/src/`; strict types, no `any`. Tokens stay in memory / session storage per existing auth patterns.

## What not to commit

- `.env` files, API keys, OIDC private keys, credentials
- Build artifacts (`bin/`, `frontend/dist/`, coverage output)
- Drift between handler code and `backend/docs/swagger.*` after HTTP changes
- Drift between `frontend/src/api/types.ts` and backend API shapes

## Security

Do not open public issues for vulnerabilities. See [SECURITY.md](SECURITY.md).

## License

By contributing, you agree that your contributions are licensed under the [Apache License 2.0](LICENSE).
