-- Organization invites for emails that do not have an account yet.
CREATE TABLE "public"."tenant_invites" (
  "id" uuid NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "updated_at" timestamptz NOT NULL DEFAULT now(),
  "deleted_at" timestamptz NULL,
  "email" text NOT NULL,
  "email_lower" text NOT NULL,
  "tenant_id" uuid NOT NULL,
  "role" varchar(32) NOT NULL,
  "token_hash" text NOT NULL,
  "status" varchar(32) NOT NULL,
  "expires_at" timestamptz NOT NULL,
  "accepted_user_id" uuid NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_tenant_invites_tenant" FOREIGN KEY ("tenant_id") REFERENCES "public"."tenants" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "fk_tenant_invites_accepted_user" FOREIGN KEY ("accepted_user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE SET NULL
);

CREATE INDEX "idx_tenant_invites_deleted_at" ON "public"."tenant_invites" ("deleted_at");
CREATE UNIQUE INDEX "idx_tenant_invites_token_hash" ON "public"."tenant_invites" ("token_hash") WHERE "deleted_at" IS NULL;
CREATE UNIQUE INDEX "idx_tenant_invites_pending_email_tenant" ON "public"."tenant_invites" ("email_lower", "tenant_id") WHERE "status" = 'pending' AND "deleted_at" IS NULL;
