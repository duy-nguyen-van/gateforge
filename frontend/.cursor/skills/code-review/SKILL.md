---
name: code-review
description: "Review code quality with technical rigor. Use before PRs, after implementing features, for security audits, or when claiming task completion."
---

# Code Review

Systematic review for iam-frontend (React, TypeScript, OIDC, MFA/WebAuthn).

## Core Principle
**Technical correctness over social comfort. Be honest, specific, suggest fixes.**

## Process

### 1. Edge Case Scouting (First)
- Token refresh on 401, remember-me storage split
- MFA ticket lifecycle, OIDC `return_to` full navigation
- Admin route guard bypass, stale TanStack Query data
- WebAuthn cancel/error paths

### 2. Priority Order
1. Correctness
2. Security (token storage, CSRF, credential leakage)
3. Error handling (ApiError, form error states)
4. Layer boundaries (`frontend-layering.mdc`)
5. State races (concurrent refresh, double submit)
6. Performance (query keys, unnecessary re-renders)
7. Types + README/env docs

### 3. Route & Guard Checks
- New routes in `src/routes/index.tsx`
- Auth: Pattern A (`AuthLayout` → feature, no page file)
- Console/settings: Pattern B (page + `features/admin/` hooks)
- Correct guard: GuestRoute / ProtectedRoute / AdminRoute

### 4. React Checks
- Zod + RHF for forms
- `apiFetch` for API calls (not raw fetch except OIDC redirect)
- No access tokens in persistent storage

## Feedback Format
- File + issue + WHY + fix snippet
- Critical / High / Medium / Low

See `.cursor/rules/code-review.mdc`.
