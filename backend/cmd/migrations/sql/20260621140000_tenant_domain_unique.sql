-- Enforce unique active tenant domains for subdomain resolution.
CREATE UNIQUE INDEX "idx_tenants_domain_active" ON "public"."tenants" ("domain") WHERE "domain" IS NOT NULL AND "deleted_at" IS NULL;
