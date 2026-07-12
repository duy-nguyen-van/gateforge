Review recent code changes systematically for iam-backend.

## Steps
1. Identify changed files (`git diff` or `git diff --cached`)
2. Per file, check priority list in `.cursor/rules/code-review.mdc`:
   - Correctness, security (OIDC/session/MFA), errors, concurrency, performance, tests, OpenAPI
3. Scout IAM edge cases: token expiry, CSRF on `/oidc/login`, session/`iam_session` cookie, MFA replay
4. Verify route namespace: root OIDC vs `/api/v1`
5. Report Critical / High / Medium / Low with fix suggestions
6. Flag missing `make swagger-load` ( `/api/v1`) or manual OIDC doc updates (root routes)
