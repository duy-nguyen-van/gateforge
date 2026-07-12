-- Phase 2 OIDC: authorization code metadata + unique code lookup + dev public client (PKCE).
ALTER TABLE "public"."authorization_codes" ADD COLUMN IF NOT EXISTS "nonce" text NULL;
ALTER TABLE "public"."authorization_codes" ADD COLUMN IF NOT EXISTS "scope" text NULL;

CREATE UNIQUE INDEX IF NOT EXISTS "idx_authorization_codes_code_unique" ON "public"."authorization_codes" ("code");

INSERT INTO "public"."clients" (
  "id",
  "created_at",
  "updated_at",
  "tenant_id",
  "client_id",
  "client_secret",
  "name",
  "redirect_uris",
  "grant_types",
  "scopes",
  "is_public"
)
SELECT
  'a1b2c3d4-e5f6-4789-a012-3456789abcde'::uuid,
  now(),
  now(),
  '00000000-0000-0000-0000-000000000001'::uuid,
  'oidc-dev',
  NULL,
  'OIDC dev public client (PKCE)',
  ARRAY[
    'http://127.0.0.1:5173/callback',
    'http://localhost:5173/callback',
    'http://127.0.0.1:3000/callback',
    'http://localhost:3000/callback'
  ]::text[],
  ARRAY['authorization_code']::text[],
  ARRAY['openid', 'email', 'profile']::text[],
  true
WHERE NOT EXISTS (
  SELECT 1 FROM "public"."clients"
  WHERE "tenant_id" = '00000000-0000-0000-0000-000000000001'::uuid AND "client_id" = 'oidc-dev'
);
