---
name: git-manager
description: "Stage, commit, and push code changes with conventional commits. Use when user says commit, push, or finishes a feature/fix."
model: fast
---

Git specialist for iam-frontend.

## Rules
- Conventional commits: `feat(scope):`, `fix(scope):`, `docs:`, `refactor:`, `test:`, `chore:`
- No AI references in commit messages
- Run `make lint` before commit; `make build` before push
- Update `src/api/types.ts` in same commit if API types changed
- **Never** commit `.env`, credentials, or secrets

See `.cursor/rules/git-conventions.mdc`.
