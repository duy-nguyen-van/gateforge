-- Refresh-token families, client lookup index, refresh grant, and append-only audit logs.
ALTER TABLE "public"."refresh_tokens" ADD COLUMN IF NOT EXISTS "family_id" uuid NULL;
ALTER TABLE "public"."refresh_tokens" ADD COLUMN IF NOT EXISTS "scope" text NULL;

UPDATE "public"."refresh_tokens" SET "family_id" = "id" WHERE "family_id" IS NULL;

ALTER TABLE "public"."refresh_tokens" ALTER COLUMN "family_id" SET NOT NULL;

CREATE INDEX IF NOT EXISTS "idx_refresh_tokens_family_revoked" ON "public"."refresh_tokens" ("family_id", "revoked");

CREATE INDEX IF NOT EXISTS "idx_clients_client_id_active" ON "public"."clients" ("client_id") WHERE "deleted_at" IS NULL;

UPDATE "public"."clients"
SET "grant_types" = array_append("grant_types", 'refresh_token')
WHERE "grant_types" IS NOT NULL
  AND NOT ('refresh_token' = ANY("grant_types"));

CREATE OR REPLACE FUNCTION "public"."audit_logs_append_only"() RETURNS trigger
LANGUAGE plpgsql AS $$
BEGIN
  RAISE EXCEPTION 'audit_logs is append-only';
END;
$$;

DROP TRIGGER IF EXISTS "audit_logs_no_update_delete" ON "public"."audit_logs";
CREATE TRIGGER "audit_logs_no_update_delete"
BEFORE UPDATE OR DELETE ON "public"."audit_logs"
FOR EACH ROW EXECUTE FUNCTION "public"."audit_logs_append_only"();
