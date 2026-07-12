Run the iam-frontend quality checks and analyze results.

## Steps
1. `make lint`
2. `make build` (TypeScript + Vite production build)
3. Or combined: `make check`
4. For each failure: root cause + suggested fix
5. If auth-related changes: list manual browser test steps

## Report
- Lint: pass/fail with error summary
- Build: pass/fail with TS error locations
- Recommended manual tests for changed auth/OIDC/MFA flows
