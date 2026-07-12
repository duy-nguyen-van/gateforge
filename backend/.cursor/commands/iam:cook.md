Implement the described feature end-to-end.

Principles: **YAGNI, KISS, DRY** — minimal code, match existing patterns.

## Workflow

### Step 1: Understand
- Clarify requirements; note auth/session/OIDC impact

### Step 2: Scout
- Read `AGENTS.md` layer table and `.cursor/rules/go-layering.mdc`
- Find similar handlers, services, repos in the same domain

### Step 3: Plan
- Files to touch: FX wiring (`main.go`), `router.go` (pick root vs `/api/v1`), handler, service, repo, DTO, constants
- CSRF skip list if new browser/cookie paths
- Migrations and swag (`/api/v1` only) if needed

### Step 4: Implement
- Handler → service → repo flow
- Wire in `cmd/server/main.go` and `cmd/server/routes/router.go`
- `go build ./...` after significant changes

### Step 5: Test
- Colocated `*_test.go`; run scoped then `make tests`

### Step 6: HTTP / Schema
- swag annotations + `make swagger-load` for `/api/v1` API changes
- Manual flow docs for root OIDC routes ([docs/features/](../../docs/features/), [docs/testing/](../../docs/testing/))
- Atlas SQL + `make migrate-hash` if schema changed

### Step 7: Self-Review
- `.cursor/rules/code-review.mdc` checklist
- Suggest commit message; do not commit unless asked

## Rules
- Update existing files; no duplicate "enhanced" modules
- Never ignore failing tests
