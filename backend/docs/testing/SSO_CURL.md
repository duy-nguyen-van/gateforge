# SSO session — curl testing

Concept guide: [features/SSO_SESSION.md](../features/SSO_SESSION.md)

Prove: **after `POST /api/v1/login`, `GET /authorize` succeeds without `POST /oidc/login`.**

Use the **same base URL** for every step (e.g. `http://localhost:3000`). PKCE parameters must match a registered OAuth client in your DB.

## A) curl (cookie jar)

1. **Dashboard login** (stores `iam_session` in `jar.txt`):

   ```bash
   BASE=http://localhost:3000
   curl -c jar.txt -b jar.txt -sS -X POST "$BASE/api/v1/login" \
     -H 'Content-Type: application/json' \
     -d '{"email":"YOU@example.com","password":"YOUR_PASSWORD"}'
   ```

   Expect JSON with tokens **and** `Set-Cookie: iam_session=...`.

2. **Authorize** (use your real `client_id`, PKCE, `redirect_uri`):

   ```bash
   curl -c jar.txt -b jar.txt -sS -D - -o /dev/null \
     "$BASE/authorize?response_type=code&client_id=CLIENT_ID&redirect_uri=REDIRECT_URI_ENCODED&scope=openid%20profile&state=test&code_challenge=CHALLENGE&code_challenge_method=S256"
   ```

   **Success:** `302` with `Location:` containing `code=`.  
   **SSO not working:** `302` to the **login** page.

3. **Logout** (optional):

   ```bash
   curl -c jar.txt -b jar.txt -sS -X POST "$BASE/api/v1/logout" \
     -H "Authorization: Bearer ACCESS_TOKEN_FROM_STEP_1"
   ```

## B) Browser

1. DevTools → **Network**, enable **Preserve log**.
2. Sign in via `POST /api/v1/login` with `credentials: 'include'` from the IAM origin.
3. Navigate to full `/authorize?...` URL in the same tab.
4. Expect redirect to client with `code`, not login page.

## C) Postman

Use [postman/IAM_OIDC.postman_collection.json](../postman/IAM_OIDC.postman_collection.json), but **replace** POST `/oidc/login` with **POST `/api/v1/login`** (JSON `email` / `password`, no `return_to`). Send **`Cookie: iam_session=...`** on the following GET `/authorize`.

## Related

- [testing/OIDC_CURL.md](OIDC_CURL.md) — full OIDC + PKCE flow
- [features/OIDC.md](../features/OIDC.md)
