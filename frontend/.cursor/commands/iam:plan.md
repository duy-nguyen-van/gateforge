Create a detailed implementation plan for the described feature. Do NOT implement — only plan.

## Workflow

### Step 1: Clarify Requirements
- Ask questions if ambiguous
- Identify functional and non-functional requirements (auth, OIDC redirect, MFA impact)

### Step 2: Research
- Read affected areas per `AGENTS.md` layer table
- Follow `.cursor/rules/frontend-layering.mdc` import boundaries
- Search for similar features (login, admin console, security settings)
- Check iam-backend docs if API/OIDC behavior is involved

### Step 3: Design Solution
- Pick route pattern: **A** (auth: routes → `AuthLayout` → feature) or **B** (routes → page → features)
- **Route placement**: `src/routes/index.tsx` + guard type (guest/protected/admin)
- Console pages: query hooks in `features/admin/`, shared UI in `console-state.tsx`
- API calls: new functions in `src/api/client.ts`, types in `src/api/types.ts`
- State: auth context vs TanStack Query vs local form state
- Env vars: new `VITE_*` in `.env.example`
- Backend alignment: proxy, OIDC login URL, WebAuthn origins

### Step 4: Write Plan
Save in `./plans/{date}-{slug}/`:
- `plan.md` — overview, phases, risks, success criteria
- `phase-01-*.md` — per-phase steps, files, todos

### Step 5: Self-Review
- Simpler approach? YAGNI/KISS?
- Security gaps? Token storage? CSRF for OIDC?
- Manual test strategy

### Step 6: Present
- Summarize trade-offs
- Ask for approval before implementation

## Rules
- Do NOT implement code
- Each phase independently testable
- Match existing React/Hook Form/Query patterns
