# OIDC — curl and Postman testing

Concept guide: [features/OIDC.md](../features/OIDC.md)

This document shows how to exercise the IdP flow with **curl**. The API is assumed at **`http://localhost:3000`**; set `BASE` to match your deployment.

**Endpoints (this service)**

| Step | Method | Path |
|------|--------|------|
| Discovery | GET | `/.well-known/openid-configuration` |
| Authorize | GET | `/authorize` |
| Login (session, OIDC) | POST | `/oidc/login` |
| Login (dashboard/API) | POST | `/api/v1/login` |
| Token | POST | `/token` |
| UserInfo | GET | `/userinfo` |

**Session cookie:** `iam_session` (HTTP-only), set when login succeeds.

## Multi-tenant membership

Users are global identities; access to an organization is granted via `tenant_memberships`. Dashboard JWTs include `tenant_id` (active tenant). When a user belongs to multiple tenants and no tenant context is provided at login, `POST /api/v1/login` returns `selection_required` with a short-lived `selection_token` — complete with `POST /api/v1/tenants/select`. Switch later with `POST /api/v1/tenants/switch` (Bearer auth). OIDC `/authorize` resolves tenant from `client_id` and requires an active membership in that tenant. See [features/MULTI_TENANT.md](../features/MULTI_TENANT.md).

## Postman import

### Option A — Collection file (recommended)

1. In Postman: **Import** → **File** → select  
   `docs/postman/IAM_OIDC.postman_collection.json`
2. Open the collection **Variables** tab and set at least:  
   `baseUrl`, `clientId`, `redirectUriEncoded`, `codeVerifier`, `codeChallenge`, `state`, `email`, `password`
3. Run requests **in order** (0 → 1 → 2 → 3 → 4 → 5).  
   - Request **0** saves `csrfToken` from the `X-CSRF-Token` response header (needed for **2**).  
   - Request **2** does not follow redirects; its **Tests** script reads `Set-Cookie` and sets collection variable `iam_session`.  
   - Request **3** sends `Cookie: iam_session={{iam_session}}` explicitly.  
   - After **3**, copy the `code` query param from the **Location** header into variable `authCode`, then run **4**.

`redirectUriEncoded` must already be **URL-encoded**. Example:  
`http://127.0.0.1:9999/callback` → `http%3A%2F%2F127.0.0.1%3A9999%2Fcallback`

### Option B — Import raw cURL

Postman: **Import** → **Raw text** → paste **one** of the blocks below (repeat for each request).

**0) Discovery**

```bash
curl --location --request GET 'BASE/.well-known/openid-configuration'
```

**1) GET /authorize (first visit, no session)**

```bash
curl --location --request GET 'BASE/authorize?response_type=code&client_id=CLIENT_ID&redirect_uri=REDIRECT_URI_ENCODED&scope=openid%20profile&state=STATE&code_challenge=CODE_CHALLENGE&code_challenge_method=S256'
```

**2) POST /oidc/login** — use `X-CSRF-Token` from request **0**.

```bash
curl --location --request POST 'BASE/oidc/login' \
--header 'Content-Type: application/json' \
--header 'X-CSRF-Token: PASTE_CSRF_TOKEN' \
--data-raw '{"email":"EMAIL","password":"PASSWORD","return_to":"BASE/authorize?response_type=code&client_id=CLIENT_ID&redirect_uri=REDIRECT_URI_ENCODED&scope=openid%20profile&state=STATE&code_challenge=CODE_CHALLENGE&code_challenge_method=S256"}'
```

**3) GET /authorize (with cookie)**

```bash
curl --location --request GET 'BASE/authorize?response_type=code&client_id=CLIENT_ID&redirect_uri=REDIRECT_URI_ENCODED&scope=openid%20profile&state=STATE&code_challenge=CODE_CHALLENGE&code_challenge_method=S256'
```

**4) POST /token**

```bash
curl --location --request POST 'BASE/token' \
--header 'Content-Type: application/json' \
--data-raw '{"grant_type":"authorization_code","code":"AUTH_CODE","client_id":"CLIENT_ID","code_verifier":"CODE_VERIFIER"}'
```

**5) GET /userinfo**

```bash
curl --location --request GET 'BASE/userinfo' \
--header 'Authorization: Bearer OIDC_ACCESS_TOKEN'
```

---

## PKCE (S256)

Generate once per authorization attempt. Use the **same** `CODE_VERIFIER` at the `/token` step.

```bash
BASE="http://localhost:3000"

CODE_VERIFIER="$(openssl rand -base64 48 | tr -d '\n' | tr '+/' '-_' | tr -d '=')"
CODE_CHALLENGE="$(printf '%s' "$CODE_VERIFIER" | openssl dgst -binary -sha256 | openssl base64 | tr -d '\n' | tr '+/' '-_' | tr -d '=')"

echo "CODE_VERIFIER=$CODE_VERIFIER"
echo "CODE_CHALLENGE=$CODE_CHALLENGE"
```

---

## Authorization code + PKCE (cookie jar)

Replace `CLIENT_ID`, `REDIRECT_URI`, email, and password with values that match your environment. The OAuth client must be registered in the database.

### About `state`

`state` is echoed unchanged on the redirect to your app (`redirect_uri?code=...&state=...`). Your client must verify it matches the value sent on `/authorize` (CSRF protection). The IdP does not interpret `state`.

### 1) Start `/authorize` (no session → redirect to login)

```bash
CLIENT_ID="your-client-id"
REDIRECT_URI="http://127.0.0.1:9999/callback"
STATE="$(openssl rand -hex 16)"

curl -sS -D - -o /dev/null \
  -c /tmp/iam-cookies.txt \
  "$BASE/authorize?response_type=code&client_id=$CLIENT_ID&redirect_uri=$(python3 -c "import urllib.parse; print(urllib.parse.quote('$REDIRECT_URI', safe=''))")&scope=openid%20profile&state=$STATE&code_challenge=$CODE_CHALLENGE&code_challenge_method=S256"
```

### 2) Log in with `return_to` = full authorize URL

```bash
RETURN_TO="$BASE/authorize?response_type=code&client_id=$CLIENT_ID&redirect_uri=$(python3 -c "import urllib.parse; print(urllib.parse.quote('$REDIRECT_URI', safe=''))")&scope=openid%20profile&state=$STATE&code_challenge=$CODE_CHALLENGE&code_challenge_method=S256"

curl -sS -D - \
  -c /tmp/iam-cookies.txt -b /tmp/iam-cookies.txt \
  -X POST "$BASE/oidc/login" \
  -H "Content-Type: application/json" \
  -d "{\"email\":\"you@example.com\",\"password\":\"your-password\",\"return_to\":\"$RETURN_TO\"}"
```

Expect **`302`** and **`Set-Cookie: iam_session=...`**.

### 3) Call `/authorize` again with the cookie

```bash
curl -sS -D - -o /dev/null \
  -b /tmp/iam-cookies.txt \
  "$BASE/authorize?response_type=code&client_id=$CLIENT_ID&redirect_uri=$(python3 -c "import urllib.parse; print(urllib.parse.quote('$REDIRECT_URI', safe=''))")&scope=openid%20profile&state=$STATE&code_challenge=$CODE_CHALLENGE&code_challenge_method=S256"
```

Copy **`code`** from **`Location:`**.

### 4) Exchange code for tokens

```bash
AUTH_CODE="paste-code-from-Location-here"

curl -sS -X POST "$BASE/token" \
  -H "Content-Type: application/json" \
  -d "{
    \"grant_type\": \"authorization_code\",
    \"code\": \"$AUTH_CODE\",
    \"client_id\": \"$CLIENT_ID\",
    \"code_verifier\": \"$CODE_VERIFIER\"
  }" | jq .
```

### 5) UserInfo

```bash
OIDC_ACCESS_TOKEN="paste-access_token-from-token-response"

curl -sS "$BASE/userinfo" \
  -H "Authorization: Bearer $OIDC_ACCESS_TOKEN" | jq .
```

---

## Notes

| Topic | Detail |
|--------|--------|
| OAuth client | `CLIENT_ID` and `REDIRECT_URI` must exist in the DB |
| User registration | `POST $BASE/api/v1/register` if the user does not exist |
| `return_to` | Must pass validation: same origin as `APP_BASE_URL`, path starting with `/authorize` |
| Login redirect | Configure **`OIDC_LOGIN_PAGE_URL`** for SPA login host |

## Confidential clients

```bash
curl -sS -X POST "$BASE/token" \
  -u "$CLIENT_ID:client-secret" \
  -H "Content-Type: application/json" \
  -d "{
    \"grant_type\": \"authorization_code\",
    \"code\": \"$AUTH_CODE\",
    \"code_verifier\": \"$CODE_VERIFIER\"
  }" | jq .
```

## Related

- [features/SSO_SESSION.md](../features/SSO_SESSION.md) — API login instead of `/oidc/login`
- [testing/SSO_CURL.md](SSO_CURL.md) — SSO cookie-jar test
