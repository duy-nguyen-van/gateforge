# Passkeys and MFA — curl testing

Concept guide: [features/PASSKEY_MFA.md](../features/PASSKEY_MFA.md)

Passkey ceremonies use **start** → **finish**; MFA enrollment is authenticated; post-login MFA uses **`mfa_ticket`**.

| Step | Method | Path | Auth |
|------|--------|------|------|
| Register user | POST | `/api/v1/register` | Public |
| Login (password) | POST | `/api/v1/login` | Public |
| Passkey register start | POST | `/api/v1/webauthn/register/start` | Bearer |
| Passkey register finish | POST | `/api/v1/webauthn/register/finish` | Bearer |
| Passkey login start | POST | `/api/v1/webauthn/login/start` | Public |
| Passkey login finish | POST | `/api/v1/webauthn/login/finish` | Public |
| TOTP setup | POST | `/api/v1/mfa/totp/setup` | Bearer |
| TOTP verify (enable) | POST | `/api/v1/mfa/totp/verify` | Bearer |
| Recovery codes | POST | `/api/v1/mfa/recovery-codes` | Bearer |
| MFA after login | POST | `/api/v1/mfa/challenge/verify` | Public |

Set **`BASE=http://localhost:3000`**. **`WEBAUTHN_RP_ORIGINS`** must include that origin.

## Prerequisites

1. Migrations: `webauthn_credentials`, `user_mfa_totps`, `user_mfa_recovery_codes`.
2. Redis for WebAuthn sessions and MFA tickets.
3. Valid **`JWT_SECRET`** (min 32 chars).

## 1) Register and login

```bash
export BASE=http://localhost:3000
export TENANT_ID=00000000-0000-0000-0000-000000000001

curl -sS -X POST "${BASE}/api/v1/register" \
  -H "Content-Type: application/json" \
  -d '{"email":"passkey-user@example.com","password":"long-secure-passphrase","first_name":"Test","last_name":"User"}' \
  | tee register.json

curl -sS -X POST "${BASE}/api/v1/login" \
  -H "Content-Type: application/json" \
  -d '{"email":"passkey-user@example.com","password":"long-secure-passphrase"}' \
  | tee login.json

export ACCESS_TOKEN="$(jq -r '.data.access_token // empty' login.json)"
```

## 2) Passkey register start

```bash
curl -sS -X POST "${BASE}/api/v1/webauthn/register/start" \
  -H "Authorization: Bearer ${ACCESS_TOKEN}" \
  -H "Content-Type: application/json" \
  -d '{"device_name":"curl-doc-laptop"}' \
  | tee wa-reg-start.json

export WA_REG_SESSION="$(jq -r '.data.session_token // empty' wa-reg-start.json)"
```

## 3) Passkey register finish (browser required)

Run **`navigator.credentials.create()`** with `data.options`, then POST to finish. **curl alone cannot complete WebAuthn.**

```bash
curl -sS -X POST "${BASE}/api/v1/webauthn/register/finish" \
  -H "Authorization: Bearer ${ACCESS_TOKEN}" \
  -H "Content-Type: application/json" \
  -d "{\"session_token\":\"${WA_REG_SESSION}\",\"credential\":CREDENTIAL_JSON}"
```

## 4) Passkey login start

```bash
curl -sS -X POST "${BASE}/api/v1/webauthn/login/start" \
  -H "Content-Type: application/json" \
  -d "{\"email\":\"passkey-user@example.com\",\"tenant_id\":\"${TENANT_ID}\"}" \
  | tee wa-login-start.json

export WA_LOGIN_SESSION="$(jq -r '.data.session_token // empty' wa-login-start.json)"
```

Anti-enumeration: **200** on start even for unknown email; finish returns **401** if no valid credential.

## 5) Passkey login finish (browser required)

Without MFA: **`LoginResponse`** + **`Set-Cookie: iam_session`**.  
With MFA: **`mfa_ticket`**, **`expires_in`** — then step 8.

## 6) TOTP setup and verify

```bash
curl -sS -X POST "${BASE}/api/v1/mfa/totp/setup" \
  -H "Authorization: Bearer ${ACCESS_TOKEN}" \
  | tee mfa-setup.json

read -r -p "Enter 6-digit code: " TOTP_CODE

curl -sS -X POST "${BASE}/api/v1/mfa/totp/verify" \
  -H "Authorization: Bearer ${ACCESS_TOKEN}" \
  -H "Content-Type: application/json" \
  -d "{\"code\":\"${TOTP_CODE}\"}"
```

## 7) Recovery codes

```bash
curl -sS -X POST "${BASE}/api/v1/mfa/recovery-codes" \
  -H "Authorization: Bearer ${ACCESS_TOKEN}" \
  | tee recovery.json

jq '.data.codes' recovery.json
```

Store codes immediately — shown once only.

## 8) Password login with MFA

```bash
curl -sS -X POST "${BASE}/api/v1/login" \
  -H "Content-Type: application/json" \
  -d '{"email":"passkey-user@example.com","password":"long-secure-passphrase"}' \
  | tee login-mfa.json

export MFA_TICKET="$(jq -r '.data.mfa_ticket // empty' login-mfa.json)"

read -r -p "Enter TOTP or recovery code: " MFA_CODE

curl -sS -X POST "${BASE}/api/v1/mfa/challenge/verify" \
  -H "Content-Type: application/json" \
  -d "{\"mfa_ticket\":\"${MFA_TICKET}\",\"code\":\"${MFA_CODE}\"}" \
  | tee mfa-verify.json
```

6 digits → TOTP; otherwise → recovery code.

## CSRF note

**`/api/v1/webauthn/*`** and **`/api/v1/mfa/*`** are excluded from CSRF middleware.

## Related

- [testing/OIDC_CURL.md](OIDC_CURL.md) — OIDC after login
- [testing/FEDERATION_CURL.md](FEDERATION_CURL.md) — Google federation
