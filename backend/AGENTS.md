# Agent context (backend)

Identity and access backend: **Echo**, **Uber FX**, **PostgreSQL**, **Redis**, **Atlas** migrations, **OIDC**, federation, sessions, MFA/WebAuthn. Go module: `github.com/gateforge-iam/gateforge-iam`.

## Before you ship a change

- **Read feature docs first** when changing OIDC, SSO, federation, WebAuthn, MFA, tenants, or admin auth: [docs/README.md](docs/README.md) and the relevant [docs/features/](docs/features/) guide (see `rules/feature-knowledge.mdc`).
- Run `make lint` and `make tests` (see `Makefile` for scoped test targets).
- Run `go mod tidy` / `make dep` when dependencies change.
- Run `make swagger-load` when routes, handlers, or swag annotations change.
- Never commit secrets; use `cmd/server/.env.example` and `examples/env/` as references only.

## Where things live (by layer)

| Layer | Path |
|--------|------|
| FX wiring, HTTP server lifecycle | `cmd/server/main.go` |
| Echo routes + middleware stack | `cmd/server/routes/router.go` |
| HTTP handlers | `internal/handlers/` |
| Use cases / orchestration | `internal/services/` |
| DB access (GORM) | `internal/repositories/` |
| Persistence entities | `internal/models/` |
| JSON/API DTOs | `internal/dtos/` |
| Domain adapters (WebAuthn, OIDC, …) | `internal/domains/` |
| Tokens, OIDC signing, auth context, providers | `internal/auth/` |
| Crypto helpers (e.g. MFA secrets) | `internal/crypto/` |
| Echo middleware | `internal/middlewares/` |
| Structured errors + Echo error middleware | `internal/errors/` |
| Shared constants (codes, OAuth, session) | `internal/constants/` |
| Configuration | `internal/config/` |
| Postgres / DB manager | `internal/db/` |
| Redis cache | `internal/cache/` |
| SES, S3/GCS, … | `internal/integration/` |
| Outbound HTTP (Resty) | `internal/httpclient/` |
| Request-scoped helpers | `internal/request/` |
| Logging / APM | `internal/logger/`, `internal/monitoring/` |
| Pure utils | `internal/utils/` |
| Small shared packages | `pkg/` (e.g. `pkg/correlationid/`) |
| SQL migrations + checksum | `cmd/migrations/sql/`, `atlas.sum` |
| OpenAPI (swag-generated) | `docs/swagger.yaml`, `docs/swagger.json`, `docs/docs.go` |

**Import direction and exceptions** (e.g. thin admin handlers using a repo interface): see `.cursor/rules/go-layering.mdc`.

## HTTP route surfaces

Routes are registered in `cmd/server/routes/router.go`. **Two URL namespaces** — do not assume everything is under `/api/v1`:

| Surface | Base path | Examples | Auth |
|---------|-----------|----------|------|
| **OIDC / OAuth2** | `/` (root) | `/.well-known/jwks.json`, `/.well-known/openid-configuration`, `/authorize`, `/token`, `/userinfo`, `/oidc/login`, `/oidc/federation/:provider/*` | Browser session cookie (`iam_session`) for authorize; Bearer for userinfo |
| **App API** | `/api/v1` | `/register`, `/login`, `/refresh`, `/me`, `/logout`, `/webauthn/*`, `/mfa/*` | Public, JWT (`JWTBearerAuth`), or `X-Admin-API-Key` (internal tenant IdP) |
| **Ops / docs** | `/` | `/swagger/*` (non-prod, basic auth) | Basic auth |

- swag `@BasePath` is `/api/v1` — OIDC root routes are **not** in the generated spec; document them in manual flow docs when changed.
- CSRF skip list is path-specific — see `internal/middlewares/csrf.go` (e.g. `/api/v1/login` skipped; `/oidc/login` **not** skipped).

## Key patterns

### Error handling
- Prefer `internal/errors.AppError` and `errors.ErrorHandler` via `handlers.BaseHandler`
- Error codes in `internal/constants/`
- **Exception**: `TenantIdentityAdminHandler` uses `echo.NewHTTPError` for admin-key checks — match that file if extending internal operator routes

### HTTP handlers
- Echo handlers embed `BaseHandler`; wire via `Provide*` FX constructors (e.g. `ProvideAuthHandler`)
- Bind DTOs, call services, return JSON via `HandleError` / `SuccessResponse`
- Add swag annotations on `/api/v1` handler methods; regenerate with `make swagger-load`

### Services
- Business logic only; accept `context.Context` and typed params — no `echo.Context`
- May use repos, auth, cache, integration, domains, crypto

### Logging & observability
- **Global Zap** initialized in `cmd/server/main.go` via `logger.Init()` — use `logger.Log` / `logger.Sugar`, not injected `slog`
- Request logging: `internal/middlewares/logging.go` (`RequestLogging`); Sentry/New Relic in `internal/monitoring/`
- Correlation ID: `pkg/correlationid/` via `internal/request/context.go` helpers

### Local dev
```bash
make container-up   # dependencies
make migrate-up     # Atlas migrations
make up             # run server
make bootstrap      # all three
```

## Codex configuration

- Use `$backend` for backend implementation, debugging, testing, review, architecture, migrations, or API documentation work.
- The skill is defined at `.codex/skills/backend/SKILL.md` and distills the `.cursor/` backend rules for Codex.
- Project subagents live in `.codex/agents/`: `backend-reviewer`, `backend-debugger`, and `backend-tester`. Spawn them only when the user asks for subagents, parallel work, review, debugging, or validation.

## Cursor configuration (`.cursor/`)

### Always-on rules
| Rule | Purpose |
|------|---------|
| `rules/iam-core.mdc` | Module, stack, quality bar |
| `rules/development-rules.mdc` | YAGNI/KISS/DRY, pre-commit |
| `rules/openapi-specs.mdc` | swag sync on HTTP changes |
| `rules/security-practices.mdc` | IAM security (OIDC, sessions, MFA) |
| `rules/feature-knowledge.mdc` | Read feature docs before auth changes |
| `rules/git-conventions.mdc` | Commits, branches, PRs |

### Contextual rules (by glob)
| Rule | When |
|------|------|
| `rules/go-layering.mdc` | Editing `**/*.go` |
| `rules/golang-patterns.mdc` | Editing `**/*.go` |
| `rules/go-testing.mdc` | Editing `**/*_test.go` |
| `rules/database-migrations.mdc` | Editing `cmd/migrations/**/*.sql` |
| `rules/code-review.mdc` | Code review / PRs |
| `rules/debugging-guide.mdc` | Debugging |
| `rules/documentation-standards.mdc` | Doc updates |
| `rules/workflow-implementation.mdc` | Feature / fix workflow |

### Skills
| Skill | Use when |
|-------|----------|
| `skills/feature-knowledge/SKILL.md` | Load feature docs before auth/OIDC changes |
| `skills/iam-api-manual-testing/SKILL.md` | Manual OIDC/federation/SSO/MFA curl flows |
| `skills/system-architecture/SKILL.md` | Architecture, DDD, HA, performance |
| `skills/clean-architecture-testing/SKILL.md` | Clean Architecture, coverage ≥ 90% |
| `skills/code-review/SKILL.md` | Systematic review before PRs |
| `skills/debug/SKILL.md` | Root-cause debugging |

### Agents & commands
- **Agents**: `.cursor/agents/` — `code-reviewer`, `debugger`, `tester`, `git-manager`
- **Commands**: `.cursor/commands/iam:*` — `plan`, `cook`, `fix`, `test`, `review-code`
- **MCP example**: `.cursor/mcp.json.example` (copy to `mcp.json` locally; gitignored)

## Human docs

| Doc | Purpose |
|-----|---------|
| `README.md` | Backend setup |
| [`../CONTRIBUTING.md`](../CONTRIBUTING.md) | Monorepo contribution |
| [docs/README.md](docs/README.md) | **Hub** — feature index, DB/Redis reference, routes |
| [docs/features/](docs/features/) | OIDC, SSO, federation, passkey/MFA, multi-tenant, authorization |
| [docs/testing/](docs/testing/) | curl / Postman manual testing |
| `docs/swagger.yaml` | OpenAPI for `/api/v1` only |
| `docs/postman/IAM_OIDC.postman_collection.json` | Postman collection |
