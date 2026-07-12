# IAM Frontend

React + TypeScript SPA for the [backend](../backend) identity service.

## Stack

- Vite + React 19 + TypeScript
- React Router, TanStack Query, React Hook Form + Zod
- Tailwind CSS v4 + shadcn/ui-style components
- `@simplewebauthn/browser` for passkeys
- `qrcode.react` for TOTP enrollment

## Prerequisites

- Node.js 26+
- Running `backend` API on `http://localhost:3000`

## Setup

```bash
cp .env.example .env
npm install
npm run dev
```

Open [http://localhost:5173](http://localhost:5173). Vite proxies `/api/v1`, `/oidc`, `/authorize`, `/token`, `/userinfo`, and `/.well-known` to the backend so cookies and SSO work on one origin during development.

## Routes

| Route | Description |
|-------|-------------|
| `/login` | Email/password, passkey, Google federation (when configured in admin) |
| `/register` | Create account |
| `/mfa/challenge` | Post-login TOTP / recovery code |
| `/settings/profile` | Current user profile |
| `/settings/security` | TOTP setup, recovery codes, passkey registration |

## Environment

| Variable | Description |
|----------|-------------|
| `VITE_API_BASE_URL` | Empty for same-origin; set only if API is on another host |
| `VITE_API_PROXY_TARGET` | Backend URL for Vite dev proxy (default `http://localhost:3000`) |
| `VITE_DEFAULT_TENANT_ID` | Tenant UUID for WebAuthn login and federation provider lookup |

Align backend env vars (`OIDC_LOGIN_PAGE_URL`, `WEBAUTHN_RP_ORIGINS`, etc.) — see `backend/cmd/server/.env.example`.

## Integration with backend

Run the Go API on `:3000` and this SPA on `:5173`. Vite proxies API/OIDC paths so the browser sees one origin (cookies + SSO work).

| Frontend (`.env`) | Backend (`cmd/server/.env`) | Purpose |
|-------------------|----------------------------|---------|
| `VITE_API_PROXY_TARGET=http://localhost:3000` | `APP_HTTP_SERVER=:3000` | Dev proxy target |
| *(empty `VITE_API_BASE_URL`)* | `APP_BASE_URL=http://localhost:3000` | OIDC issuer stays on API port |
| — | `OIDC_LOGIN_PAGE_URL=http://localhost:5173/login` | `/authorize` → SPA login |
| — | `WEBAUTHN_RP_ORIGINS=http://localhost:5173` | Passkeys on Vite origin |
| `VITE_DEFAULT_TENANT_ID` | `DEFAULT_TENANT_ID` | Same tenant UUID |
| `VITE_DEFAULT_TENANT_ID` | `DEFAULT_TENANT_ID` | Same tenant UUID; federation credentials live in DB |

On first startup, set backend `BOOTSTRAP_ADMIN_EMAIL` and `BOOTSTRAP_ADMIN_PASSWORD` to create the platform admin. Log in with that account to access console admin APIs.

Refresh tokens are stored in `sessionStorage` (or `localStorage` when “Remember me” is checked) so reload keeps the session via `POST /api/v1/refresh`.

## API types

Types are maintained in `src/api/types.ts`. Optional regeneration from Swagger:

```bash
npm install -D openapi-typescript --legacy-peer-deps
npm run generate:api
```

## Production deployment

### Primary: hybrid single-binary (recommended)

From the repo root, build and run the Go server with the SPA embedded:

```bash
make build-prod
SERVE_EMBEDDED_FRONTEND=true ./bin/gateforge-iam-server
```

Set `APP_BASE_URL`, `OIDC_LOGIN_PAGE_URL`, and `WEBAUTHN_RP_ORIGINS` to your production domain (same origin). Leave `VITE_API_BASE_URL` empty when building — API calls use relative paths.

See [`../deployments/README.md`](../deployments/README.md) and [`../README.md`](../README.md).

### Alternative: split static + reverse proxy

Serve the built `dist/` on the **same origin** as the Go API:

- **Cloudflare Pages** — host static assets; proxy `/api/v1/*`, `/oidc/*`, `/authorize`, `/token`, `/.well-known/*` to the Go origin (see `public/_redirects` as a Netlify-style reference).
- **Fly.io / Railway / VPS** — use [`deployments/caddy/Caddyfile`](../deployments/caddy/Caddyfile) to serve `dist/` and reverse-proxy API/OIDC paths to the Go backend on one domain.

```bash
cd frontend && npm run build
# Copy dist/ to /srv/dist, set API_UPSTREAM, then from repo root:
caddy run --config deployments/caddy/Caddyfile
```

## Project layout

```
src/
  api/          HTTP client, types, optional generated/
  auth/         Token store, schemas, AuthProvider, context
  components/   ui/, layout/, icons/, brand/, avatars/
  features/     login, register, mfa, webauthn, admin (queries + console-state)
  hooks/        use-auth
  lib/          cn, apiUrl
  pages/        home, profile, security, console/*, logout
  routes/       Router + guards (auth routes wire features directly)
```
