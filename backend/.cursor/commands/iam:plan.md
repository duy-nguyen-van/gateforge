Create a detailed implementation plan for the described feature. Do NOT implement — only plan.

## Workflow

### Step 1: Clarify Requirements
- Ask questions if ambiguous
- Identify functional and non-functional requirements (auth, sessions, OIDC impact)

### Step 2: Research
- Read affected areas per `AGENTS.md` layer table
- Follow `.cursor/rules/go-layering.mdc` import boundaries
- Search for similar handlers/services/repos patterns
- Check [docs/README.md](../../docs/README.md) and feature docs if IAM-related

### Step 3: Design Solution
- Data flow: handler → service → repo (+ cache/auth/domains)
- **Route placement**: root OIDC (`/authorize`, …) vs `/api/v1` — see `router.go`
- API contracts and DTOs; manual flow docs for root OIDC
- Atlas migration needs
- swag annotation plan ( `/api/v1` only)
- Error codes in `internal/constants/`
- Security: middleware, CSRF skip list (`csrf.go`), `iam_session` cookie

### Step 4: Write Plan
Save in `./plans/{date}-{slug}/`:
- `plan.md` — overview, phases, risks, success criteria
- `phase-01-*.md` — per-phase steps, files, todos

### Step 5: Self-Review
- Simpler approach? YAGNI/KISS?
- Security gaps? Layer violations?
- Test strategy (handler stubs, service fakes)?

### Step 6: Present
- Summarize trade-offs
- Ask for approval before implementation

## Rules
- Do NOT implement code
- Each phase independently testable
- Match existing Echo/FX/GORM patterns
