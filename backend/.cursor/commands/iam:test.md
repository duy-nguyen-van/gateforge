Run the backend quality gates and analyze results.

## Steps
1. Diff check: list new/changed behavior and confirm covering `*_test.go` exists (fail gate if missing)
2. `make lint` (or `go vet ./...` if lint unavailable)
3. Scoped tests if changes are localized:
   - `make test-handlers` / `make test-services` / `make test-repositories`
   - Prefer `go test -race -cover ./internal/<package>` for changed packages (≥ 90%)
4. Full suite: `make tests` (coverage + race + count=1)
5. Optional coverage report: `make test-coverage` or `make test-coverage-html`
6. Security: spot-check touched paths against `security-practices.mdc`
7. For each failure: root cause + suggested fix
8. Report: lint, pass/fail/skip, coverage % (changed packages), race findings, missing tests, security gaps
