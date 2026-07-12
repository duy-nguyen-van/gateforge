-- Phase 4: federated identities and per-tenant IdP toggles (Google first).

-- Create "tenant_identity_providers" table
CREATE TABLE "public"."tenant_identity_providers" (
  "id" uuid NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "updated_at" timestamptz NOT NULL DEFAULT now(),
  "deleted_at" timestamptz NULL,
  "tenant_id" uuid NOT NULL,
  "provider" character varying(32) NOT NULL,
  "enabled" boolean NOT NULL DEFAULT false,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_tenant_identity_providers_tenant" FOREIGN KEY ("tenant_id") REFERENCES "public"."tenants" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);

CREATE INDEX "idx_tenant_identity_providers_deleted_at" ON "public"."tenant_identity_providers" ("deleted_at");
CREATE UNIQUE INDEX "idx_tenant_identity_providers_tenant_provider" ON "public"."tenant_identity_providers" ("tenant_id", "provider");

-- Create "federated_identities" table
CREATE TABLE "public"."federated_identities" (
  "id" uuid NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "updated_at" timestamptz NOT NULL DEFAULT now(),
  "deleted_at" timestamptz NULL,
  "tenant_id" uuid NOT NULL,
  "user_id" uuid NOT NULL,
  "provider" character varying(32) NOT NULL,
  "subject" text NOT NULL,
  "email_at_link" text NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_federated_identities_tenant" FOREIGN KEY ("tenant_id") REFERENCES "public"."tenants" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "fk_federated_identities_user" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);

CREATE INDEX "idx_federated_identities_deleted_at" ON "public"."federated_identities" ("deleted_at");
CREATE INDEX "idx_federated_identities_user_id" ON "public"."federated_identities" ("user_id");
CREATE UNIQUE INDEX "idx_federated_identities_tenant_provider_subject" ON "public"."federated_identities" ("tenant_id", "provider", "subject");

-- Default tenant: Google disabled until credentials are configured and flag is toggled.
INSERT INTO "public"."tenant_identity_providers" ("id", "created_at", "updated_at", "tenant_id", "provider", "enabled")
SELECT 'b0000000-0000-4000-8000-000000000001', now(), now(), '00000000-0000-0000-0000-000000000001', 'google', false
WHERE NOT EXISTS (
  SELECT 1 FROM "public"."tenant_identity_providers"
  WHERE "tenant_id" = '00000000-0000-0000-0000-000000000001' AND "provider" = 'google' AND "deleted_at" IS NULL
);
