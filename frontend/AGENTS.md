# Agent context (iam-frontend)

React + TypeScript SPA for the [iam-backend](../iam-backend) identity service: login, register, MFA, WebAuthn, OIDC browser flows, admin console.

## Before you ship a change

- Run `make lint` and `make build` (or `make check`).
- Update `src/api/types.ts` when backend API shapes change.
- Never commit secrets; use `.env.example` as reference only.

## Where things live (by layer)

| Layer | Path |
|--------|------|
| Entry + providers | `src/main.tsx` (`QueryClient`, `BrowserRouter`, `AuthProvider`) |
| App shell | `src/App.tsx` — renders `AppRoutes` only |
| Router + guards | `src/routes/index.tsx`, `src/routes/guards.tsx` |
| HTTP client + types | `src/api/client.ts`, `src/api/types.ts`, optional `src/api/generated/` |
| Auth state | `src/auth/auth-provider.tsx`, `src/auth/auth-context-value.ts`, `src/auth/token-store.ts`, `src/auth/schemas.ts` |
| Auth hook | `src/hooks/use-auth.ts` |
| UI primitives | `src/components/ui/` |
| Layouts + nav | `src/components/layout/` (`auth-layout`, `console-layout`, `console-nav.ts`, shells) |
| Icons / brand / avatars | `src/components/icons/`, `brand/`, `avatars/` |
| Feature modules | `src/features/login/`, `register/`, `mfa/`, `webauthn/`, `admin/` |
| Route pages | `src/pages/` (home, profile, security, `console/*`, logout) |
| Utilities | `src/lib/utils.ts` |
| Styles / tokens | `src/index.css` |
| Vite config + proxy | `vite.config.ts` |
| Deploy | `deploy/Caddyfile`, `public/_redirects` |

**Import direction and route patterns**: see `.cursor/rules/frontend-layering.mdc`.

## Route wiring (how screens are mounted)

### Pattern A — Auth (routes → layout → feature)
No page file. Feature forms wired directly in `src/routes/index.tsx`:

| Route | Layout | Feature |
|-------|--------|---------|
| `/login` | `AuthLayout variant="revamp"` | `LoginForm` |
| `/register` | `AuthLayout` | `RegisterForm` |
| `/mfa/challenge` | `AuthLayout` | `MfaChallengeForm` |

### Pattern B — Pages (routes → page → features/components)
| Route | Guard | Page | Uses |
|-------|-------|------|------|
| `/` | — | `HomePage` | brand, icons (self-contained) |
| `/settings/profile` | Protected | `ProfilePage` | `useAuth`, `DefaultAvatar` |
| `/settings/security` | Protected | `SecurityPage` | `TotpSetupPanel`, `PasskeyRegisterPanel` |
| `/console` | Protected + Admin | `DashboardPage` | admin queries, console-state |
| `/console/users` | Protected + Admin | `UsersPage` | `useAdminUsers`, local search state |
| `/console/clients` | Protected + Admin | `ClientsPage` | `useAdminClients` |
| `/console/tenants` | Protected + Admin | `TenantsPage` | `useAdminTenants` |
| `/console/identity-providers` | Protected + Admin | `IdentityProvidersPage` | `useTenantIdentityProviders` |
| `/console/audit-logs` | Protected + Admin | `AuditLogsPage` | placeholder / empty state |
| `/logout` | Protected | `LogoutPage` | auth logout |

Guest routes (`/login`, `/register`) use `GuestRoute`. Admin console routes nest `AdminRoute` inside `ProtectedRoute` + `ConsoleLayout`.

## Key patterns

### API client
- `apiFetch()` — Bearer auth, 401 refresh retry, `credentials: 'include'`
- `prefetchCsrfToken()` — CSRF for `/oidc/login`
- `ApiError` for typed HTTP failures
- Types hand-maintained in `types.ts`; optional codegen to `api/generated/schema.ts`

### Auth
- Access token in memory; refresh in session/local storage per "Remember me"
- `AuthProvider` bootstraps via refresh + `getMe()`
- MFA ticket in `sessionStorage` until verify completes
- Consume via `useAuth()` — don't duplicate token logic in features

### Forms (auth features)
- React Hook Form + Zod (`loginSchema` in `auth/schemas.ts`)
- Feature components own submit/error/loading state
- `login-form` may import sibling features (`passkey-login`, social icons)

### Admin console
- Query hooks + keys: `src/features/admin/use-admin-queries.ts`
- Shared loading/error/empty UI: `src/features/admin/console-state.tsx`
- Display helpers: `src/features/admin/admin-utils.ts`
- **Pages own query calls and local filter state** — not thin wrappers

### Local dev
```bash
make setup    # .env + npm install
make run      # Vite :5173, proxies API to :3000
make check    # lint + build
```

Requires iam-backend on `http://localhost:3000` (see README env alignment table).

## Cursor configuration (`.cursor/`)

### Always-on rules
| Rule | Purpose |
|------|---------|
| `rules/iam-core.mdc` | Stack, quality bar, backend integration |
| `rules/development-rules.mdc` | YAGNI/KISS/DRY, pre-commit |
| `rules/security-practices.mdc` | Tokens, CSRF, cookies, XSS |
| `rules/git-conventions.mdc` | Commits, branches, PRs |

### Contextual rules (by glob)
| Rule | When |
|------|------|
| `rules/frontend-layering.mdc` | Editing `src/**/*` |
| `rules/react-patterns.mdc` | Editing `src/**/*.{ts,tsx}` |
| `rules/code-review.mdc` | Code review / PRs |
| `rules/debugging-guide.mdc` | Debugging |
| `rules/documentation-standards.mdc` | Doc updates |
| `rules/workflow-implementation.mdc` | Feature / fix workflow |

### Skills
| Skill | Use when |
|-------|----------|
| `skills/code-review/SKILL.md` | Systematic review before PRs |
| `skills/debug/SKILL.md` | Root-cause debugging |

### Agents & commands
- **Agents**: `.cursor/agents/` — `code-reviewer`, `debugger`, `tester`, `git-manager`
- **Commands**: `.cursor/commands/iam:*` — `plan`, `cook`, `fix`, `test`, `review-code`
- **MCP example**: `.cursor/mcp.json.example` (Stitch; copy to `mcp.json` locally; gitignored)

## Human docs

| Doc | Purpose |
|-----|---------|
| `README.md` | Setup, routes, env, backend alignment, deployment |
| [backend/docs/README.md](../backend/docs/README.md) | IAM hub — features, DB, Redis, routes |
| [backend/docs/features/](../backend/docs/features/) | OIDC, SSO, federation, passkey/MFA, multi-tenant, authorization |
| [backend/docs/testing/](../backend/docs/testing/) | curl / Postman manual testing |
