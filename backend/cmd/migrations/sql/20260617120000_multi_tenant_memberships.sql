-- Multi-tenant membership: global users + tenant_memberships junction table.

-- 1. Create tenant_memberships
CREATE TABLE "public"."tenant_memberships" (
  "id" uuid NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "updated_at" timestamptz NOT NULL DEFAULT now(),
  "deleted_at" timestamptz NULL,
  "user_id" uuid NOT NULL,
  "tenant_id" uuid NOT NULL,
  "role" character varying(32) NOT NULL DEFAULT 'member',
  "status" character varying(32) NOT NULL DEFAULT 'active',
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_tenant_memberships_user" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "fk_tenant_memberships_tenant" FOREIGN KEY ("tenant_id") REFERENCES "public"."tenants" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);

CREATE INDEX "idx_tenant_memberships_deleted_at" ON "public"."tenant_memberships" ("deleted_at");
CREATE UNIQUE INDEX "idx_tenant_memberships_user_tenant" ON "public"."tenant_memberships" ("user_id", "tenant_id") WHERE "deleted_at" IS NULL;
CREATE INDEX "idx_tenant_memberships_tenant_id" ON "public"."tenant_memberships" ("tenant_id");

-- 2. Backfill one membership per existing user
INSERT INTO "public"."tenant_memberships" ("id", "created_at", "updated_at", "user_id", "tenant_id", "role", "status")
SELECT gen_random_uuid(), u."created_at", u."updated_at", u."id", u."tenant_id", 'member', 'active'
FROM "public"."users" u
WHERE u."deleted_at" IS NULL;

-- 3. Merge duplicate users (same email_lower across tenants) into canonical user (earliest created_at)
DO $$
DECLARE
  dup RECORD;
  canonical_id uuid;
BEGIN
  FOR dup IN
    SELECT email_lower
    FROM "public"."users"
    WHERE deleted_at IS NULL
    GROUP BY email_lower
    HAVING COUNT(*) > 1
  LOOP
    SELECT id INTO canonical_id
    FROM "public"."users"
    WHERE email_lower = dup.email_lower AND deleted_at IS NULL
    ORDER BY created_at ASC
    LIMIT 1;

    -- Re-point child rows to canonical user
    UPDATE "public"."password_credentials" SET user_id = canonical_id
    WHERE user_id IN (
      SELECT id FROM "public"."users"
      WHERE email_lower = dup.email_lower AND deleted_at IS NULL AND id <> canonical_id
    )
    AND NOT EXISTS (SELECT 1 FROM "public"."password_credentials" WHERE user_id = canonical_id);

    UPDATE "public"."password_credentials" SET user_id = canonical_id
    WHERE user_id IN (
      SELECT id FROM "public"."users"
      WHERE email_lower = dup.email_lower AND deleted_at IS NULL AND id <> canonical_id
    );

    UPDATE "public"."webauthn_credentials" SET user_id = canonical_id
    WHERE user_id IN (
      SELECT id FROM "public"."users"
      WHERE email_lower = dup.email_lower AND deleted_at IS NULL AND id <> canonical_id
    );

    UPDATE "public"."user_mfa_totps" SET user_id = canonical_id
    WHERE user_id IN (
      SELECT id FROM "public"."users"
      WHERE email_lower = dup.email_lower AND deleted_at IS NULL AND id <> canonical_id
    );

    UPDATE "public"."user_mfa_recovery_codes" SET user_id = canonical_id
    WHERE user_id IN (
      SELECT id FROM "public"."users"
      WHERE email_lower = dup.email_lower AND deleted_at IS NULL AND id <> canonical_id
    );

    UPDATE "public"."sessions" SET user_id = canonical_id
    WHERE user_id IN (
      SELECT id FROM "public"."users"
      WHERE email_lower = dup.email_lower AND deleted_at IS NULL AND id <> canonical_id
    );

    UPDATE "public"."refresh_tokens" SET user_id = canonical_id
    WHERE user_id IN (
      SELECT id FROM "public"."users"
      WHERE email_lower = dup.email_lower AND deleted_at IS NULL AND id <> canonical_id
    );

    UPDATE "public"."access_tokens" SET user_id = canonical_id
    WHERE user_id IN (
      SELECT id FROM "public"."users"
      WHERE email_lower = dup.email_lower AND deleted_at IS NULL AND id <> canonical_id
    );

    UPDATE "public"."authorization_codes" SET user_id = canonical_id
    WHERE user_id IN (
      SELECT id FROM "public"."users"
      WHERE email_lower = dup.email_lower AND deleted_at IS NULL AND id <> canonical_id
    );

    UPDATE "public"."consents" SET user_id = canonical_id
    WHERE user_id IN (
      SELECT id FROM "public"."users"
      WHERE email_lower = dup.email_lower AND deleted_at IS NULL AND id <> canonical_id
    );

    UPDATE "public"."federated_identities" SET user_id = canonical_id
    WHERE user_id IN (
      SELECT id FROM "public"."users"
      WHERE email_lower = dup.email_lower AND deleted_at IS NULL AND id <> canonical_id
    );

    -- Ensure memberships exist for all tenants of merged users
    INSERT INTO "public"."tenant_memberships" ("id", "created_at", "updated_at", "user_id", "tenant_id", "role", "status")
    SELECT gen_random_uuid(), now(), now(), canonical_id, u.tenant_id, 'member', 'active'
    FROM "public"."users" u
    WHERE u.email_lower = dup.email_lower AND u.deleted_at IS NULL
    ON CONFLICT DO NOTHING;

    -- Soft-delete duplicate user rows
    UPDATE "public"."users" SET deleted_at = now()
    WHERE email_lower = dup.email_lower AND deleted_at IS NULL AND id <> canonical_id;

    -- Remove memberships for soft-deleted users
    UPDATE "public"."tenant_memberships" SET deleted_at = now()
    WHERE user_id IN (
      SELECT id FROM "public"."users"
      WHERE email_lower = dup.email_lower AND id <> canonical_id AND deleted_at IS NOT NULL
    );
  END LOOP;
END $$;

-- 4. Drop users.tenant_id
ALTER TABLE "public"."users" DROP CONSTRAINT IF EXISTS "fk_users_tenant";
DROP INDEX IF EXISTS "idx_users_tenant_email_unique";
ALTER TABLE "public"."users" DROP COLUMN IF EXISTS "tenant_id";
CREATE UNIQUE INDEX "idx_users_email_lower" ON "public"."users" ("email_lower") WHERE "deleted_at" IS NULL;

-- 5. Global federated identities (drop per-tenant uniqueness)
DELETE FROM "public"."federated_identities" fi
WHERE fi.id NOT IN (
  SELECT DISTINCT ON (provider, subject) id
  FROM "public"."federated_identities"
  WHERE deleted_at IS NULL
  ORDER BY provider, subject, created_at ASC
);

ALTER TABLE "public"."federated_identities" DROP CONSTRAINT IF EXISTS "fk_federated_identities_tenant";
DROP INDEX IF EXISTS "idx_federated_identities_tenant_provider_subject";
ALTER TABLE "public"."federated_identities" DROP COLUMN IF EXISTS "tenant_id";
CREATE UNIQUE INDEX "idx_federated_identities_provider_subject" ON "public"."federated_identities" ("provider", "subject") WHERE "deleted_at" IS NULL;

-- 6. Drop redundant tenant_id from user-owned credential tables
ALTER TABLE "public"."user_mfa_totps" DROP CONSTRAINT IF EXISTS "fk_user_mfa_totps_tenant";
DROP INDEX IF EXISTS "idx_user_mfa_totps_tenant_id";
ALTER TABLE "public"."user_mfa_totps" DROP COLUMN IF EXISTS "tenant_id";

ALTER TABLE "public"."webauthn_credentials" DROP CONSTRAINT IF EXISTS "fk_webauthn_credentials_tenant";
DROP INDEX IF EXISTS "idx_webauthn_credentials_tenant_id";
ALTER TABLE "public"."webauthn_credentials" DROP COLUMN IF EXISTS "tenant_id";
