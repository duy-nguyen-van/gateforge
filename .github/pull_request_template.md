## Description

Please include a summary of the change and which issue is fixed.

## Type of change

- [ ] Bug fix
- [ ] New feature
- [ ] Breaking change
- [ ] Documentation update

## Areas touched

- [ ] `backend/` (Go API, OIDC, migrations)
- [ ] `frontend/` (React SPA)
- [ ] `api/` (canonical OpenAPI for SDKs)
- [ ] `sdk/` (Go / TypeScript clients)
- [ ] Root tooling (Makefile, Docker, CI, deployments)

## How Has This Been Tested?

Describe the tests that you ran to verify your changes. Include commands and relevant details.

## Checklist

- [ ] Backend: `make -C backend lint && make -C backend tests` (if `backend/` changed)
- [ ] Frontend: `make -C frontend check` (if `frontend/` changed)
- [ ] OpenAPI / SDKs: `api/openapi.yaml` updated and `make sdk-generate` when the public API changed
- [ ] Prod binary: `make build-prod` (if embed/build pipeline changed)
- [ ] Docs updated for behavior changes (`backend/docs/`, CONTRIBUTING, README as needed)
- [ ] No secrets, `.env` files, or private keys in the diff
- [ ] Tests cover the fix or new behavior
