# Federation (Google) — curl testing

Concept guide: [features/FEDERATION.md](../features/FEDERATION.md)

This service acts as an **OAuth/OIDC client** to Google, then issues the same **`iam_session`** cookie used by `GET /authorize`.

| Step | Method | Path |
|------|--------|------|
| Start (redirect to IdP) | GET | `/oidc/federation/:provider/start` |
| Callback (IdP redirects browser here) | GET | `/oidc/federation/:provider/callback` |
| Enable Google per tenant (admin) | PATCH | `/api/v1/internal/tenants/:tenantId/identity-providers/google` |

For Google, use **`provider=google`**.

Set **`BASE`** to your API origin (e.g. `http://localhost:3000`). It must match **`APP_BASE_URL`** for safe `return_to` validation.

## Prerequisites

1. **Database:** migrations applied (`tenant_identity_providers`, `federated_identities`).
2. **Redis:** `CACHE_PROVIDER=redis`.
3. **Environment** (`cmd/server/.env.example`):
   - `GOOGLE_FEDERATION_ENABLED=true`
   - `GOOGLE_OAUTH_CLIENT_ID`, `GOOGLE_OAUTH_CLIENT_SECRET`, `GOOGLE_OAUTH_REDIRECT_URL`
   - `DEFAULT_TENANT_ID`
4. **Tenant flag:** Google enabled for tenant (admin PATCH below).

## 1) Enable Google for the default tenant

```bash
export BASE=http://localhost:3000
export TENANT_ID=00000000-0000-0000-0000-000000000001
export ADMIN_KEY='your-admin-api-key'

curl -sS -X PATCH "${BASE}/api/v1/internal/tenants/${TENANT_ID}/identity-providers/google" \
  -H "Content-Type: application/json" \
  -H "X-Admin-API-Key: ${ADMIN_KEY}" \
  -d '{"enabled":true}' \
  -w "\nHTTP %{http_code}\n"
```

Expect **`204`**. If `ADMIN_API_KEY` is unset, route returns **404**.

## 2) Inspect the start redirect

```bash
curl -sS -D - -o /dev/null \
  "${BASE}/oidc/federation/google/start?return_to=%2Fauthorize"
```

Expect **`302`** to `accounts.google.com`.

Full `return_to` (OIDC authorize URL):

```bash
curl -sS -D - -o /dev/null \
  "${BASE}/oidc/federation/google/start?return_to=$(python3 -c "import urllib.parse; print(urllib.parse.quote('/authorize?response_type=code&client_id=demo&redirect_uri=http%253A%252F%252F127.0.0.1%253A9999%252Fcb&scope=openid&state=s1&code_challenge=CHALLENGE&code_challenge_method=S256', safe=''))")"
```

## 3) Full login — browser required

After **start**, sign in at Google. Google redirects to your callback with `?code=...&state=...`.

**Manual test:**

1. Open `${BASE}/oidc/federation/google/start?return_to=%2Fauthorize` in a browser.
2. Finish Google login.
3. Browser lands on **`return_to`** with **`Set-Cookie: iam_session=...`**.
4. Continue with [testing/OIDC_CURL.md](OIDC_CURL.md) for `/authorize` + `/token`.

**Optional callback error check:**

```bash
curl -sS -D - "${BASE}/oidc/federation/google/callback?code=fake&state=fake"
```

Expect **4xx**.

## 4) IdP error redirect

```bash
curl -sS "${BASE}/oidc/federation/google/callback?error=access_denied&error_description=cancelled"
```

## CSRF note

**`/oidc/federation/*`** is excluded from CSRF middleware.

## Quick reference

```text
GET {BASE}/oidc/federation/google/start?return_to=<encoded /authorize...>
     → 302 Google
User signs in at Google
GET {BASE}/oidc/federation/google/callback?code=...&state=...
     → 302 return_to + Set-Cookie iam_session

GET {BASE}/authorize?...  with Cookie: iam_session=...
     → 302 to client with ?code=...
```

## Related

- [testing/OIDC_CURL.md](OIDC_CURL.md) — token exchange after authorize
