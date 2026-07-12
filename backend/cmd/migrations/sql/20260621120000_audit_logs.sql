-- Create "audit_logs" table
CREATE TABLE "public"."audit_logs" (
  "id" uuid NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "updated_at" timestamptz NOT NULL DEFAULT now(),
  "deleted_at" timestamptz NULL,
  "tenant_id" uuid NULL,
  "action" character varying(100) NOT NULL,
  "result" character varying(20) NOT NULL,
  "actor_type" character varying(50) NOT NULL,
  "actor_id" text NULL,
  "resource_type" character varying(50) NULL,
  "resource_id" uuid NULL,
  "resource_name" character varying(255) NULL,
  "ip_address" inet NULL,
  "user_agent" text NULL,
  "request_id" text NULL,
  "correlation_id" text NULL,
  "old_value" jsonb NULL,
  "new_value" jsonb NULL,
  PRIMARY KEY ("id")
);

-- Create index "idx_audit_logs_deleted_at" to table: "audit_logs"
CREATE INDEX "idx_audit_logs_deleted_at" ON "public"."audit_logs" ("deleted_at");

-- Create index "idx_audit_logs_tenant_created_at" to table: "audit_logs"
CREATE INDEX "idx_audit_logs_tenant_created_at" ON "public"."audit_logs" ("tenant_id", "created_at" DESC);

-- Create index "idx_audit_logs_action_created_at" to table: "audit_logs"
CREATE INDEX "idx_audit_logs_action_created_at" ON "public"."audit_logs" ("action", "created_at" DESC);

-- Create index "idx_audit_logs_actor_created_at" to table: "audit_logs"
CREATE INDEX "idx_audit_logs_actor_created_at" ON "public"."audit_logs" ("actor_id", "created_at" DESC)
WHERE "actor_id" IS NOT NULL;

-- Create index "idx_audit_logs_resource" to table: "audit_logs"
CREATE INDEX "idx_audit_logs_resource" ON "public"."audit_logs" ("resource_type", "resource_id", "created_at" DESC)
WHERE "resource_id" IS NOT NULL;
