-- Create "tenants" table
CREATE TABLE "public"."tenants" (
  "id" uuid NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "updated_at" timestamptz NOT NULL DEFAULT now(),
  "deleted_at" timestamptz NULL,
  "name" character varying(255) NULL,
  "domain" character varying(255) NULL,
  PRIMARY KEY ("id")
);

-- Create index "idx_tenants_deleted_at" to table: "tenants"
CREATE INDEX "idx_tenants_deleted_at" ON "public"."tenants" ("deleted_at");

-- Create "clients" table
CREATE TABLE "public"."clients" (
  "id" uuid NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "updated_at" timestamptz NOT NULL DEFAULT now(),
  "deleted_at" timestamptz NULL,
  "tenant_id" uuid NOT NULL,
  "client_id" character varying(255) NOT NULL,
  "client_secret" character varying(255) NULL,
  "name" character varying(255) NULL,
  "redirect_uris" text [] NULL,
  "grant_types" text [] NULL,
  "scopes" text [] NULL,
  "is_public" boolean NULL DEFAULT false,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_clients_tenant" FOREIGN KEY ("tenant_id") REFERENCES "public"."tenants" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);

-- Create index "idx_clients_deleted_at" to table: "clients"
CREATE INDEX "idx_clients_deleted_at" ON "public"."clients" ("deleted_at");

-- Create index "idx_clients_tenant_client_id" to table: "clients"
CREATE UNIQUE INDEX "idx_clients_tenant_client_id" ON "public"."clients" ("tenant_id", "client_id");

-- Create "users" table
CREATE TABLE "public"."users" (
  "id" uuid NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "updated_at" timestamptz NOT NULL DEFAULT now(),
  "deleted_at" timestamptz NULL,
  "first_name" text NULL,
  "last_name" text NULL,
  "email" text NOT NULL,
  "email_lower" text NOT NULL,
  "email_verified" boolean NULL,
  "status" text NULL,
  "tenant_id" uuid NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_users_tenant" FOREIGN KEY ("tenant_id") REFERENCES "public"."tenants" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);

-- Create index "idx_users_deleted_at" to table: "users"
CREATE INDEX "idx_users_deleted_at" ON "public"."users" ("deleted_at");

-- Create index "idx_users_tenant_email_unique" to table: "users"
CREATE UNIQUE INDEX "idx_users_tenant_email_unique" ON "public"."users" ("email_lower", "tenant_id");

-- Create "access_tokens" table
CREATE TABLE "public"."access_tokens" (
  "id" uuid NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "updated_at" timestamptz NOT NULL DEFAULT now(),
  "tenant_id" uuid NOT NULL,
  "user_id" uuid NULL,
  "oauth_client_id" character varying(255) NULL,
  "token_hash" text NOT NULL,
  "expires_at" timestamptz NULL,
  "client_record_id" uuid NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_access_tokens_client" FOREIGN KEY ("client_record_id") REFERENCES "public"."clients" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_access_tokens_tenant" FOREIGN KEY ("tenant_id") REFERENCES "public"."tenants" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_access_tokens_user" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);

-- Create index "idx_access_tokens_client_record_id" to table: "access_tokens"
CREATE INDEX "idx_access_tokens_client_record_id" ON "public"."access_tokens" ("client_record_id");

-- Create index "idx_access_tokens_o_auth_client_id" to table: "access_tokens"
CREATE INDEX "idx_access_tokens_o_auth_client_id" ON "public"."access_tokens" ("oauth_client_id");

-- Create index "idx_access_tokens_tenant_id" to table: "access_tokens"
CREATE INDEX "idx_access_tokens_tenant_id" ON "public"."access_tokens" ("tenant_id");

-- Create index "idx_access_tokens_token_hash" to table: "access_tokens"
CREATE UNIQUE INDEX "idx_access_tokens_token_hash" ON "public"."access_tokens" ("token_hash");

-- Create index "idx_access_tokens_user_id" to table: "access_tokens"
CREATE INDEX "idx_access_tokens_user_id" ON "public"."access_tokens" ("user_id");

-- Create "authorization_codes" table
CREATE TABLE "public"."authorization_codes" (
  "id" uuid NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "updated_at" timestamptz NOT NULL DEFAULT now(),
  "deleted_at" timestamptz NULL,
  "code" character varying(255) NOT NULL,
  "tenant_id" uuid NOT NULL,
  "oauth_client_id" character varying(255) NOT NULL,
  "user_id" uuid NOT NULL,
  "redirect_uri" text NULL,
  "code_challenge" text NULL,
  "code_challenge_method" character varying(10) NULL,
  "expires_at" timestamptz NOT NULL,
  "client_record_id" uuid NULL,
  PRIMARY KEY ("id", "code"),
  CONSTRAINT "fk_authorization_codes_client" FOREIGN KEY ("client_record_id") REFERENCES "public"."clients" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_authorization_codes_tenant" FOREIGN KEY ("tenant_id") REFERENCES "public"."tenants" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_authorization_codes_user" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);

-- Create index "idx_authorization_codes_client_record_id" to table: "authorization_codes"
CREATE INDEX "idx_authorization_codes_client_record_id" ON "public"."authorization_codes" ("client_record_id");

-- Create index "idx_authorization_codes_deleted_at" to table: "authorization_codes"
CREATE INDEX "idx_authorization_codes_deleted_at" ON "public"."authorization_codes" ("deleted_at");

-- Create index "idx_authorization_codes_o_auth_client_id" to table: "authorization_codes"
CREATE INDEX "idx_authorization_codes_o_auth_client_id" ON "public"."authorization_codes" ("oauth_client_id");

-- Create index "idx_authorization_codes_tenant_id" to table: "authorization_codes"
CREATE INDEX "idx_authorization_codes_tenant_id" ON "public"."authorization_codes" ("tenant_id");

-- Create index "idx_authorization_codes_user_id" to table: "authorization_codes"
CREATE INDEX "idx_authorization_codes_user_id" ON "public"."authorization_codes" ("user_id");

-- Create "consents" table
CREATE TABLE "public"."consents" (
  "id" uuid NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "updated_at" timestamptz NOT NULL DEFAULT now(),
  "deleted_at" timestamptz NULL,
  "tenant_id" uuid NOT NULL,
  "user_id" uuid NOT NULL,
  "oauth_client_id" character varying(255) NOT NULL,
  "scopes" text [] NULL,
  "granted" boolean NULL DEFAULT true,
  "client_record_id" uuid NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_consents_client" FOREIGN KEY ("client_record_id") REFERENCES "public"."clients" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_consents_tenant" FOREIGN KEY ("tenant_id") REFERENCES "public"."tenants" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_consents_user" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);

-- Create index "idx_consents_client_record_id" to table: "consents"
CREATE INDEX "idx_consents_client_record_id" ON "public"."consents" ("client_record_id");

-- Create index "idx_consents_deleted_at" to table: "consents"
CREATE INDEX "idx_consents_deleted_at" ON "public"."consents" ("deleted_at");

-- Create index "idx_consents_tenant_user_oauth_client" to table: "consents"
CREATE UNIQUE INDEX "idx_consents_tenant_user_oauth_client" ON "public"."consents" ("tenant_id", "user_id", "oauth_client_id");

-- Create "password_credentials" table
CREATE TABLE "public"."password_credentials" (
  "id" uuid NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "updated_at" timestamptz NOT NULL DEFAULT now(),
  "deleted_at" timestamptz NULL,
  "user_id" uuid NOT NULL,
  "password_hash" text NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_users_password_credential" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);

-- Create index "idx_password_credentials_deleted_at" to table: "password_credentials"
CREATE INDEX "idx_password_credentials_deleted_at" ON "public"."password_credentials" ("deleted_at");

-- Create index "idx_password_credentials_user_id" to table: "password_credentials"
CREATE UNIQUE INDEX "idx_password_credentials_user_id" ON "public"."password_credentials" ("user_id");

-- Create "refresh_tokens" table
CREATE TABLE "public"."refresh_tokens" (
  "id" uuid NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "updated_at" timestamptz NOT NULL DEFAULT now(),
  "tenant_id" uuid NOT NULL,
  "user_id" uuid NOT NULL,
  "oauth_client_id" character varying(255) NOT NULL,
  "token_hash" text NOT NULL,
  "revoked" boolean NULL DEFAULT false,
  "expires_at" timestamptz NOT NULL,
  "client_record_id" uuid NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_refresh_tokens_client" FOREIGN KEY ("client_record_id") REFERENCES "public"."clients" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_refresh_tokens_tenant" FOREIGN KEY ("tenant_id") REFERENCES "public"."tenants" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_refresh_tokens_user" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);

-- Create index "idx_refresh_tokens_client_record_id" to table: "refresh_tokens"
CREATE INDEX "idx_refresh_tokens_client_record_id" ON "public"."refresh_tokens" ("client_record_id");

-- Create index "idx_refresh_tokens_tenant_id" to table: "refresh_tokens"
CREATE INDEX "idx_refresh_tokens_tenant_id" ON "public"."refresh_tokens" ("tenant_id");

-- Create index "idx_refresh_tokens_token_hash" to table: "refresh_tokens"
CREATE UNIQUE INDEX "idx_refresh_tokens_token_hash" ON "public"."refresh_tokens" ("token_hash");

-- Create index "idx_refresh_tokens_user_oauth_client_revoked" to table: "refresh_tokens"
CREATE INDEX "idx_refresh_tokens_user_oauth_client_revoked" ON "public"."refresh_tokens" ("user_id", "oauth_client_id", "revoked");

-- Create "sessions" table
CREATE TABLE "public"."sessions" (
  "id" uuid NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "updated_at" timestamptz NOT NULL DEFAULT now(),
  "deleted_at" timestamptz NULL,
  "user_id" uuid NOT NULL,
  "tenant_id" uuid NOT NULL,
  "ip_address" character varying(50) NULL,
  "user_agent" text NULL,
  "expires_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_sessions_tenant" FOREIGN KEY ("tenant_id") REFERENCES "public"."tenants" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_sessions_user" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);

-- Create index "idx_sessions_deleted_at" to table: "sessions"
CREATE INDEX "idx_sessions_deleted_at" ON "public"."sessions" ("deleted_at");

-- Create index "idx_sessions_tenant_id" to table: "sessions"
CREATE INDEX "idx_sessions_tenant_id" ON "public"."sessions" ("tenant_id");

-- Create index "idx_sessions_user_id" to table: "sessions"
CREATE INDEX "idx_sessions_user_id" ON "public"."sessions" ("user_id");

-- Create "webauthn_credentials" table
CREATE TABLE "public"."webauthn_credentials" (
  "id" uuid NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "updated_at" timestamptz NOT NULL DEFAULT now(),
  "deleted_at" timestamptz NULL,
  "user_id" uuid NOT NULL,
  "tenant_id" uuid NOT NULL,
  "credential_id" text NOT NULL,
  "public_key" text NOT NULL,
  "sign_count" bigint NOT NULL,
  "device_name" character varying(255) NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_webauthn_credentials_tenant" FOREIGN KEY ("tenant_id") REFERENCES "public"."tenants" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_webauthn_credentials_user" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);

-- Create index "idx_webauthn_credentials_credential_id" to table: "webauthn_credentials"
CREATE UNIQUE INDEX "idx_webauthn_credentials_credential_id" ON "public"."webauthn_credentials" ("credential_id");

-- Create index "idx_webauthn_credentials_deleted_at" to table: "webauthn_credentials"
CREATE INDEX "idx_webauthn_credentials_deleted_at" ON "public"."webauthn_credentials" ("deleted_at");

-- Create index "idx_webauthn_credentials_tenant_id" to table: "webauthn_credentials"
CREATE INDEX "idx_webauthn_credentials_tenant_id" ON "public"."webauthn_credentials" ("tenant_id");

-- Create index "idx_webauthn_credentials_user_id" to table: "webauthn_credentials"
CREATE INDEX "idx_webauthn_credentials_user_id" ON "public"."webauthn_credentials" ("user_id");