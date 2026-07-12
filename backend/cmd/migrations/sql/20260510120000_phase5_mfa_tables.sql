-- Phase 5: TOTP MFA + recovery codes

-- Create "user_mfa_totps" table
CREATE TABLE "public"."user_mfa_totps" (
  "id" uuid NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "updated_at" timestamptz NOT NULL DEFAULT now(),
  "deleted_at" timestamptz NULL,
  "user_id" uuid NOT NULL,
  "tenant_id" uuid NOT NULL,
  "secret_encrypted" text NOT NULL,
  "enabled" boolean NOT NULL DEFAULT false,
  "verified_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_user_mfa_totps_tenant" FOREIGN KEY ("tenant_id") REFERENCES "public"."tenants" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_user_mfa_totps_user" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);

CREATE UNIQUE INDEX "idx_user_mfa_totps_user_id" ON "public"."user_mfa_totps" ("user_id") WHERE deleted_at IS NULL;

CREATE INDEX "idx_user_mfa_totps_deleted_at" ON "public"."user_mfa_totps" ("deleted_at");

-- Create "user_mfa_recovery_codes" table
CREATE TABLE "public"."user_mfa_recovery_codes" (
  "id" uuid NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "updated_at" timestamptz NOT NULL DEFAULT now(),
  "deleted_at" timestamptz NULL,
  "user_id" uuid NOT NULL,
  "code_hash" character varying(128) NOT NULL,
  "used_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_user_mfa_recovery_codes_user" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);

CREATE INDEX "idx_user_mfa_recovery_codes_user_id" ON "public"."user_mfa_recovery_codes" ("user_id");

CREATE INDEX "idx_user_mfa_recovery_codes_deleted_at" ON "public"."user_mfa_recovery_codes" ("deleted_at");

CREATE UNIQUE INDEX "idx_user_mfa_recovery_codes_code_hash" ON "public"."user_mfa_recovery_codes" ("code_hash") WHERE deleted_at IS NULL AND used_at IS NULL;
