# Security Policy

GateForge IAM is an identity and access management platform. Security-sensitive areas include authentication, OIDC/OAuth2 flows, browser SSO sessions, JWT issuance, WebAuthn passkeys, TOTP MFA, upstream federation (e.g. Google), multi-tenant authorization, platform admin APIs, CSRF protection, and secret handling (signing keys, TOTP seeds, session tokens).

## Supported versions

| Version | Supported |
|---------|-----------|
| `main` branch | Yes |
| Released tags | Best-effort — include the tag in your report |

We do not provide long-term support for older tags unless stated in a release note.

## Reporting a vulnerability

**Do not open a public GitHub issue for security vulnerabilities.**

Report privately using one of:

1. **[GitHub Security Advisories](https://github.com/duy-nguyen-van/gateforge/security/advisories/new)** (preferred) — use **Report a vulnerability** on the repository Security tab.
2. Contact the maintainer directly if you already have a private channel.

We aim to:

- Acknowledge receipt within **72 hours**
- Provide a status update within **7 days**
- Coordinate disclosure and release a fix before public details are shared

## What to include

Help us respond quickly by providing:

- Affected component (`backend/`, `frontend/`, deployment config, or hybrid binary)
- Affected version or commit (`main` SHA or release tag)
- Clear steps to reproduce, or a minimal proof of concept
- Impact assessment (e.g. account takeover, token leakage, privilege escalation, session fixation)
- Suggested mitigation or patch, if you have one

## In scope

Examples of reports we want to hear about:

- Authentication or authorization bypass
- OIDC/OAuth2 flow weaknesses (PKCE, state/nonce, redirect URI validation)
- Session or token handling flaws (JWT, refresh tokens, `iam_session` cookie)
- CSRF gaps on browser-facing auth endpoints
- MFA/WebAuthn bypass or recovery-code weaknesses
- SQL injection, secret leakage in logs/responses, or hardcoded credentials
- Container or filesystem misconfigurations with demonstrable security impact

## Out of scope

- Issues in dependencies with no exploitable path in this project (still welcome as regular issues after checking upstream)
- Missing security headers or hardening suggestions without a demonstrated exploit
- Social engineering or physical attacks
- Denial-of-service without a novel, practical amplification vector

## Safe harbor

We appreciate responsible disclosure. We will not pursue legal action against researchers who act in good faith, avoid privacy violations and service disruption, and give us reasonable time to remediate before public disclosure.

Thank you for helping keep GateForge IAM and its users safe.
