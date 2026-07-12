-- Per-tenant OAuth credentials for federation (Google first).

ALTER TABLE "public"."tenant_identity_providers"
  ADD COLUMN "oauth_client_id" character varying(255) NULL,
  ADD COLUMN "oauth_client_secret_encrypted" text NULL;
