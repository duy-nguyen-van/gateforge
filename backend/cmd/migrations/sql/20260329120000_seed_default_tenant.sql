-- Default tenant required for users.tenant_id FK (see config DEFAULT_TENANT_ID).
INSERT INTO "public"."tenants" ("id", "created_at", "updated_at", "name")
SELECT '00000000-0000-0000-0000-000000000001', now(), now(), 'Default'
WHERE NOT EXISTS (
  SELECT 1 FROM "public"."tenants" WHERE "id" = '00000000-0000-0000-0000-000000000001'
);
