# GateForge IAM API contract

Canonical **OpenAPI 3.0.3** description of the public GateForge IAM HTTP surface.

## File

- [`openapi.yaml`](openapi.yaml) — source of truth for official SDKs

## Scope

Included:

- Root OIDC / OAuth2 routes (`/.well-known/*`, `/authorize`, `/token`, `/userinfo`, `/oidc/*`)
- Application API under `/api/v1` (auth, MFA, WebAuthn, tenants, platform admin)

Excluded from the public contract:

- `PATCH /api/v1/internal/tenants/...` (`X-Admin-API-Key`)
- `/swagger/*`

## Workflow

1. Change handlers / DTOs in `backend/`
2. Update this contract (stable `operationId`s required)
3. From the monorepo root: `make sdk-generate`
4. Commit `api/openapi.yaml` and regenerated `sdk/go/openapi` + `sdk/typescript/src/generated`
5. CI (`sdk.yml`) fails on drift

Swagger UI in non-prod still serves swag output for `/api/v1` exploration. Do not treat `backend/docs/swagger.yaml` as the SDK contract — prefer this file.
