# Account security

Password policy, login lockout, and password reset for first-party accounts.

## Password policy

New passwords (register, reset, and bootstrap admin) must be 12–128 characters and must not be one of a small built-in denylist. Existing password hashes are not rechecked at login. Bcrypt cost is unchanged.

## Login lockout

Redis key `iam:lockout:{sha256(email_lower)}`.

| Setting | Default |
|---------|---------|
| `LOCKOUT_MAX_FAILURES` | 5 |
| `LOCKOUT_WINDOW` | 15m |

A failed password increments the key. At the threshold, further attempts return the same generic authentication error as a bad password. Success deletes the key. The audit action is `auth.login_locked` and the log line does not include the email. Counter: `auth_lockout_total`.

Lockout is skipped when Redis is unset or `LOCKOUT_MAX_FAILURES` is 0.

## Password reset

| Method | Path | Behavior |
|--------|------|----------|
| POST | `/api/v1/forgot-password` | Always 200. No account enumeration |
| POST | `/api/v1/reset-password` | Single-use token, then new password |

Redis key `iam:password-reset:{sha256(token)}`, TTL `PASSWORD_RESET_TTL` (default 15m). The value is the user id only. The email link is `{APP_BASE_URL}/reset-password?token=...`. The message body comes from the embedded `password_reset` template.

Both routes skip CSRF (same as login) and use the auth rate limit. A successful reset revokes refresh tokens and browser sessions for that user.

## Related

- [AUTHORIZATION.md](AUTHORIZATION.md) — admin MFA gate
- [OIDC.md](OIDC.md) — client secret hashing and refresh rotation
