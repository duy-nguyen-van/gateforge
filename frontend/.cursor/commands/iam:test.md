Run the frontend quality gates and analyze results.

## Steps
1. Diff check: confirm new/changed logic has automated tests when a runner exists; else require a short manual checklist
2. `make lint`
3. `make build` (TypeScript + Vite production build)
4. Or combined: `make check`
5. Security: spot-check touched paths against `security-practices.mdc`
6. For each failure: root cause + suggested fix
7. If auth-related changes: list manual browser test steps

## Report
- Lint: pass/fail with error summary
- Build: pass/fail with TS error locations
- Tests / manual checklist: covered / gaps
- Security checklist: pass / gaps
- Recommended manual tests for changed auth/OIDC/MFA flows
