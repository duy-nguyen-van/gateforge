-- Persist platform admin flag on users (DB-backed platform admin authorization).
ALTER TABLE "public"."users"
  ADD COLUMN "is_platform_admin" boolean NOT NULL DEFAULT false;

CREATE INDEX "idx_users_is_platform_admin" ON "public"."users" ("is_platform_admin")
  WHERE "is_platform_admin" = true AND "deleted_at" IS NULL;
