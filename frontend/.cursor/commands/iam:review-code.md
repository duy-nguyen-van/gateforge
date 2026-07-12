Review recent code changes systematically for iam-frontend.

## Steps
1. Identify changed files (`git diff` or `git diff --cached`)
2. Per file, check priority list in `.cursor/rules/code-review.mdc`:
   - Correctness, security (tokens/CSRF), errors, state, performance, types, docs
3. Scout IAM edge cases: token refresh, MFA ticket cleanup, OIDC return_to, admin guard
4. Verify routes and guards in `src/routes/`
5. Report Critical / High / Medium / Low with fix suggestions
6. Flag missing README / `.env.example` / `api/types.ts` / `console-nav.ts` updates
