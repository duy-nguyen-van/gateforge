Implement the described feature end-to-end.

Principles: **YAGNI, KISS, DRY** — minimal code, match existing patterns.

## Workflow

### Step 1: Understand
- Clarify requirements; note auth/OIDC/MFA impact

### Step 2: Scout
- Read `AGENTS.md` layer table and `.cursor/rules/frontend-layering.mdc`
- Find similar features in `src/features/`

### Step 3: Plan
- Pattern A (auth) or B (page): routes, pages, features, `api/client.ts`, `api/types.ts`, auth if needed
- Guards for new protected routes; `console-nav.ts` for new admin pages
- Env vars if needed

### Step 4: Implement
- Auth: feature + `AuthLayout` in routes (no page file unless needed)
- Console/settings: page composes features; admin hooks in `features/admin/`
- API client + types for new endpoints
- `make build` after significant changes

### Step 5: Test
- `make lint` and `make build`
- Manual browser test for auth flows

### Step 6: Docs
- Update README / `.env.example` if routes or env changed
- Note backend env alignment if OIDC/WebAuthn affected

### Step 7: Self-Review
- `.cursor/rules/code-review.mdc` checklist
- Suggest commit message; do not commit unless asked

## Rules
- Update existing files; no duplicate "enhanced" modules
- Never ignore failing lint/build
