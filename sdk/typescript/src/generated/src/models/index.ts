/* tslint:disable */
/* eslint-disable */
/**
 * 
 * @export
 * @interface AdminAddMemberRequest
 */
export interface AdminAddMemberRequest {
    /**
     * 
     * @type {string}
     * @memberof AdminAddMemberRequest
     */
    email: string;
    /**
     * 
     * @type {string}
     * @memberof AdminAddMemberRequest
     */
    role?: AdminAddMemberRequestRoleEnum;
}


/**
 * @export
 */
export const AdminAddMemberRequestRoleEnum = {
    member: 'member',
    admin: 'admin'
} as const;
export type AdminAddMemberRequestRoleEnum = typeof AdminAddMemberRequestRoleEnum[keyof typeof AdminAddMemberRequestRoleEnum];

/**
 * 
 * @export
 * @interface AdminAuditLogListEnvelope
 */
export interface AdminAuditLogListEnvelope {
    /**
     * 
     * @type {Meta}
     * @memberof AdminAuditLogListEnvelope
     */
    meta: Meta;
    /**
     * 
     * @type {Array<AdminAuditLogResponse>}
     * @memberof AdminAuditLogListEnvelope
     */
    data: Array<AdminAuditLogResponse>;
}
/**
 * 
 * @export
 * @interface AdminAuditLogResponse
 */
export interface AdminAuditLogResponse {
    /**
     * 
     * @type {string}
     * @memberof AdminAuditLogResponse
     */
    id?: string;
    /**
     * 
     * @type {string}
     * @memberof AdminAuditLogResponse
     */
    action?: string;
    /**
     * 
     * @type {string}
     * @memberof AdminAuditLogResponse
     */
    actor_id?: string;
    /**
     * 
     * @type {string}
     * @memberof AdminAuditLogResponse
     */
    actor_type?: string;
    /**
     * 
     * @type {string}
     * @memberof AdminAuditLogResponse
     */
    correlation_id?: string;
    /**
     * 
     * @type {string}
     * @memberof AdminAuditLogResponse
     */
    request_id?: string;
    /**
     * 
     * @type {string}
     * @memberof AdminAuditLogResponse
     */
    resource_id?: string;
    /**
     * 
     * @type {string}
     * @memberof AdminAuditLogResponse
     */
    resource_name?: string;
    /**
     * 
     * @type {string}
     * @memberof AdminAuditLogResponse
     */
    resource_type?: string;
    /**
     * 
     * @type {string}
     * @memberof AdminAuditLogResponse
     */
    result?: string;
    /**
     * 
     * @type {string}
     * @memberof AdminAuditLogResponse
     */
    tenant_id?: string;
    /**
     * 
     * @type {string}
     * @memberof AdminAuditLogResponse
     */
    ip_address?: string;
    /**
     * 
     * @type {string}
     * @memberof AdminAuditLogResponse
     */
    user_agent?: string;
    /**
     * 
     * @type {{ [key: string]: any; }}
     * @memberof AdminAuditLogResponse
     */
    old_value?: { [key: string]: any; };
    /**
     * 
     * @type {{ [key: string]: any; }}
     * @memberof AdminAuditLogResponse
     */
    new_value?: { [key: string]: any; };
    /**
     * 
     * @type {Date}
     * @memberof AdminAuditLogResponse
     */
    created_at?: Date;
}
/**
 * 
 * @export
 * @interface AdminClientEnvelope
 */
export interface AdminClientEnvelope {
    /**
     * 
     * @type {Meta}
     * @memberof AdminClientEnvelope
     */
    meta: Meta;
    /**
     * 
     * @type {AdminClientResponse}
     * @memberof AdminClientEnvelope
     */
    data: AdminClientResponse;
}
/**
 * 
 * @export
 * @interface AdminClientListEnvelope
 */
export interface AdminClientListEnvelope {
    /**
     * 
     * @type {Meta}
     * @memberof AdminClientListEnvelope
     */
    meta: Meta;
    /**
     * 
     * @type {Array<AdminClientResponse>}
     * @memberof AdminClientListEnvelope
     */
    data: Array<AdminClientResponse>;
}
/**
 * 
 * @export
 * @interface AdminClientResponse
 */
export interface AdminClientResponse {
    /**
     * 
     * @type {string}
     * @memberof AdminClientResponse
     */
    id?: string;
    /**
     * 
     * @type {string}
     * @memberof AdminClientResponse
     */
    client_id?: string;
    /**
     * 
     * @type {boolean}
     * @memberof AdminClientResponse
     */
    client_secret_set?: boolean;
    /**
     * 
     * @type {string}
     * @memberof AdminClientResponse
     */
    name?: string;
    /**
     * 
     * @type {string}
     * @memberof AdminClientResponse
     */
    tenant_id?: string;
    /**
     * 
     * @type {boolean}
     * @memberof AdminClientResponse
     */
    is_public?: boolean;
    /**
     * 
     * @type {Array<string>}
     * @memberof AdminClientResponse
     */
    redirect_uris?: Array<string>;
    /**
     * 
     * @type {Array<string>}
     * @memberof AdminClientResponse
     */
    grant_types?: Array<string>;
    /**
     * 
     * @type {Array<string>}
     * @memberof AdminClientResponse
     */
    scopes?: Array<string>;
    /**
     * 
     * @type {Date}
     * @memberof AdminClientResponse
     */
    created_at?: Date;
}
/**
 * 
 * @export
 * @interface AdminClientUsageEnvelope
 */
export interface AdminClientUsageEnvelope {
    /**
     * 
     * @type {Meta}
     * @memberof AdminClientUsageEnvelope
     */
    meta: Meta;
    /**
     * 
     * @type {AdminClientUsageResponse}
     * @memberof AdminClientUsageEnvelope
     */
    data: AdminClientUsageResponse;
}
/**
 * 
 * @export
 * @interface AdminClientUsageResponse
 */
export interface AdminClientUsageResponse {
    /**
     * 
     * @type {string}
     * @memberof AdminClientUsageResponse
     */
    client_id?: string;
    /**
     * 
     * @type {number}
     * @memberof AdminClientUsageResponse
     */
    active_refresh_tokens?: number;
    /**
     * 
     * @type {number}
     * @memberof AdminClientUsageResponse
     */
    total_refresh_tokens?: number;
    /**
     * 
     * @type {number}
     * @memberof AdminClientUsageResponse
     */
    authorize_events_30d?: number;
    /**
     * 
     * @type {number}
     * @memberof AdminClientUsageResponse
     */
    token_issue_events_30d?: number;
    /**
     * 
     * @type {Date}
     * @memberof AdminClientUsageResponse
     */
    last_token_issued_at?: Date;
}
/**
 * 
 * @export
 * @interface AdminCreateClientEnvelope
 */
export interface AdminCreateClientEnvelope {
    /**
     * 
     * @type {Meta}
     * @memberof AdminCreateClientEnvelope
     */
    meta: Meta;
    /**
     * 
     * @type {AdminCreateClientResponse}
     * @memberof AdminCreateClientEnvelope
     */
    data: AdminCreateClientResponse;
}
/**
 * 
 * @export
 * @interface AdminCreateClientRequest
 */
export interface AdminCreateClientRequest {
    /**
     * 
     * @type {string}
     * @memberof AdminCreateClientRequest
     */
    client_id?: string;
    /**
     * 
     * @type {string}
     * @memberof AdminCreateClientRequest
     */
    name: string;
    /**
     * 
     * @type {string}
     * @memberof AdminCreateClientRequest
     */
    tenant_id: string;
    /**
     * 
     * @type {boolean}
     * @memberof AdminCreateClientRequest
     */
    is_public?: boolean;
    /**
     * 
     * @type {Array<string>}
     * @memberof AdminCreateClientRequest
     */
    redirect_uris: Array<string>;
    /**
     * 
     * @type {Array<string>}
     * @memberof AdminCreateClientRequest
     */
    grant_types?: Array<string>;
    /**
     * 
     * @type {Array<string>}
     * @memberof AdminCreateClientRequest
     */
    scopes?: Array<string>;
}
/**
 * 
 * @export
 * @interface AdminCreateClientResponse
 */
export interface AdminCreateClientResponse {
    /**
     * 
     * @type {string}
     * @memberof AdminCreateClientResponse
     */
    id?: string;
    /**
     * 
     * @type {string}
     * @memberof AdminCreateClientResponse
     */
    client_id?: string;
    /**
     * 
     * @type {string}
     * @memberof AdminCreateClientResponse
     */
    client_secret?: string;
    /**
     * 
     * @type {boolean}
     * @memberof AdminCreateClientResponse
     */
    client_secret_set?: boolean;
    /**
     * 
     * @type {string}
     * @memberof AdminCreateClientResponse
     */
    name?: string;
    /**
     * 
     * @type {string}
     * @memberof AdminCreateClientResponse
     */
    tenant_id?: string;
    /**
     * 
     * @type {boolean}
     * @memberof AdminCreateClientResponse
     */
    is_public?: boolean;
    /**
     * 
     * @type {Array<string>}
     * @memberof AdminCreateClientResponse
     */
    redirect_uris?: Array<string>;
    /**
     * 
     * @type {Array<string>}
     * @memberof AdminCreateClientResponse
     */
    grant_types?: Array<string>;
    /**
     * 
     * @type {Array<string>}
     * @memberof AdminCreateClientResponse
     */
    scopes?: Array<string>;
    /**
     * 
     * @type {Date}
     * @memberof AdminCreateClientResponse
     */
    created_at?: Date;
}
/**
 * 
 * @export
 * @interface AdminCreateTenantRequest
 */
export interface AdminCreateTenantRequest {
    /**
     * 
     * @type {string}
     * @memberof AdminCreateTenantRequest
     */
    name: string;
    /**
     * 
     * @type {string}
     * @memberof AdminCreateTenantRequest
     */
    domain?: string;
}
/**
 * 
 * @export
 * @interface AdminIdentityProviderListEnvelope
 */
export interface AdminIdentityProviderListEnvelope {
    /**
     * 
     * @type {Meta}
     * @memberof AdminIdentityProviderListEnvelope
     */
    meta: Meta;
    /**
     * 
     * @type {Array<AdminIdentityProviderResponse>}
     * @memberof AdminIdentityProviderListEnvelope
     */
    data: Array<AdminIdentityProviderResponse>;
}
/**
 * 
 * @export
 * @interface AdminIdentityProviderResponse
 */
export interface AdminIdentityProviderResponse {
    /**
     * 
     * @type {string}
     * @memberof AdminIdentityProviderResponse
     */
    tenant_id?: string;
    /**
     * 
     * @type {string}
     * @memberof AdminIdentityProviderResponse
     */
    provider?: string;
    /**
     * 
     * @type {string}
     * @memberof AdminIdentityProviderResponse
     */
    name?: string;
    /**
     * 
     * @type {boolean}
     * @memberof AdminIdentityProviderResponse
     */
    enabled?: boolean;
    /**
     * 
     * @type {boolean}
     * @memberof AdminIdentityProviderResponse
     */
    configured?: boolean;
    /**
     * 
     * @type {string}
     * @memberof AdminIdentityProviderResponse
     */
    oauth_client_id?: string;
    /**
     * 
     * @type {boolean}
     * @memberof AdminIdentityProviderResponse
     */
    oauth_client_secret_set?: boolean;
    /**
     * 
     * @type {string}
     * @memberof AdminIdentityProviderResponse
     */
    redirect_uri?: string;
    /**
     * 
     * @type {string}
     * @memberof AdminIdentityProviderResponse
     */
    setup_console_url?: string;
}
/**
 * 
 * @export
 * @interface AdminStatsEnvelope
 */
export interface AdminStatsEnvelope {
    /**
     * 
     * @type {Meta}
     * @memberof AdminStatsEnvelope
     */
    meta: Meta;
    /**
     * 
     * @type {AdminStatsResponse}
     * @memberof AdminStatsEnvelope
     */
    data: AdminStatsResponse;
}
/**
 * 
 * @export
 * @interface AdminStatsResponse
 */
export interface AdminStatsResponse {
    /**
     * 
     * @type {number}
     * @memberof AdminStatsResponse
     */
    active_sessions?: number;
    /**
     * 
     * @type {number}
     * @memberof AdminStatsResponse
     */
    mfa_enabled_count?: number;
    /**
     * 
     * @type {number}
     * @memberof AdminStatsResponse
     */
    mfa_enabled_percent?: number;
    /**
     * 
     * @type {number}
     * @memberof AdminStatsResponse
     */
    total_users?: number;
}
/**
 * 
 * @export
 * @interface AdminTenantEnvelope
 */
export interface AdminTenantEnvelope {
    /**
     * 
     * @type {Meta}
     * @memberof AdminTenantEnvelope
     */
    meta: Meta;
    /**
     * 
     * @type {AdminTenantResponse}
     * @memberof AdminTenantEnvelope
     */
    data: AdminTenantResponse;
}
/**
 * 
 * @export
 * @interface AdminTenantListEnvelope
 */
export interface AdminTenantListEnvelope {
    /**
     * 
     * @type {Meta}
     * @memberof AdminTenantListEnvelope
     */
    meta: Meta;
    /**
     * 
     * @type {Array<AdminTenantResponse>}
     * @memberof AdminTenantListEnvelope
     */
    data: Array<AdminTenantResponse>;
}
/**
 * 
 * @export
 * @interface AdminTenantMemberListEnvelope
 */
export interface AdminTenantMemberListEnvelope {
    /**
     * 
     * @type {Meta}
     * @memberof AdminTenantMemberListEnvelope
     */
    meta: Meta;
    /**
     * 
     * @type {Array<AdminTenantMemberResponse>}
     * @memberof AdminTenantMemberListEnvelope
     */
    data: Array<AdminTenantMemberResponse>;
}
/**
 * 
 * @export
 * @interface AdminTenantMemberResponse
 */
export interface AdminTenantMemberResponse {
    /**
     * 
     * @type {string}
     * @memberof AdminTenantMemberResponse
     */
    user_id?: string;
    /**
     * 
     * @type {string}
     * @memberof AdminTenantMemberResponse
     */
    email?: string;
    /**
     * 
     * @type {string}
     * @memberof AdminTenantMemberResponse
     */
    first_name?: string;
    /**
     * 
     * @type {string}
     * @memberof AdminTenantMemberResponse
     */
    last_name?: string;
    /**
     * 
     * @type {string}
     * @memberof AdminTenantMemberResponse
     */
    role?: string;
    /**
     * 
     * @type {string}
     * @memberof AdminTenantMemberResponse
     */
    status?: string;
    /**
     * 
     * @type {Date}
     * @memberof AdminTenantMemberResponse
     */
    joined_at?: Date;
}
/**
 * 
 * @export
 * @interface AdminTenantResponse
 */
export interface AdminTenantResponse {
    /**
     * 
     * @type {string}
     * @memberof AdminTenantResponse
     */
    id?: string;
    /**
     * 
     * @type {string}
     * @memberof AdminTenantResponse
     */
    name?: string;
    /**
     * 
     * @type {string}
     * @memberof AdminTenantResponse
     */
    domain?: string;
    /**
     * 
     * @type {number}
     * @memberof AdminTenantResponse
     */
    user_count?: number;
    /**
     * 
     * @type {Date}
     * @memberof AdminTenantResponse
     */
    created_at?: Date;
}
/**
 * 
 * @export
 * @interface AdminUpdateClientRequest
 */
export interface AdminUpdateClientRequest {
    /**
     * 
     * @type {string}
     * @memberof AdminUpdateClientRequest
     */
    name?: string;
    /**
     * 
     * @type {string}
     * @memberof AdminUpdateClientRequest
     */
    client_secret?: string;
    /**
     * 
     * @type {boolean}
     * @memberof AdminUpdateClientRequest
     */
    is_public?: boolean;
    /**
     * 
     * @type {Array<string>}
     * @memberof AdminUpdateClientRequest
     */
    redirect_uris?: Array<string>;
    /**
     * 
     * @type {Array<string>}
     * @memberof AdminUpdateClientRequest
     */
    grant_types?: Array<string>;
    /**
     * 
     * @type {Array<string>}
     * @memberof AdminUpdateClientRequest
     */
    scopes?: Array<string>;
}
/**
 * 
 * @export
 * @interface AdminUpdateTenantRequest
 */
export interface AdminUpdateTenantRequest {
    /**
     * 
     * @type {string}
     * @memberof AdminUpdateTenantRequest
     */
    name?: string;
    /**
     * 
     * @type {string}
     * @memberof AdminUpdateTenantRequest
     */
    domain?: string;
}
/**
 * 
 * @export
 * @interface AdminUserDetailEnvelope
 */
export interface AdminUserDetailEnvelope {
    /**
     * 
     * @type {Meta}
     * @memberof AdminUserDetailEnvelope
     */
    meta: Meta;
    /**
     * 
     * @type {AdminUserDetailResponse}
     * @memberof AdminUserDetailEnvelope
     */
    data: AdminUserDetailResponse;
}
/**
 * 
 * @export
 * @interface AdminUserDetailResponse
 */
export interface AdminUserDetailResponse {
    /**
     * 
     * @type {string}
     * @memberof AdminUserDetailResponse
     */
    id?: string;
    /**
     * 
     * @type {string}
     * @memberof AdminUserDetailResponse
     */
    email?: string;
    /**
     * 
     * @type {string}
     * @memberof AdminUserDetailResponse
     */
    first_name?: string;
    /**
     * 
     * @type {string}
     * @memberof AdminUserDetailResponse
     */
    last_name?: string;
    /**
     * 
     * @type {string}
     * @memberof AdminUserDetailResponse
     */
    status?: string;
    /**
     * 
     * @type {boolean}
     * @memberof AdminUserDetailResponse
     */
    mfa_enabled?: boolean;
    /**
     * 
     * @type {boolean}
     * @memberof AdminUserDetailResponse
     */
    is_platform_admin?: boolean;
    /**
     * 
     * @type {number}
     * @memberof AdminUserDetailResponse
     */
    passkey_count?: number;
    /**
     * 
     * @type {number}
     * @memberof AdminUserDetailResponse
     */
    active_sessions?: number;
    /**
     * 
     * @type {Array<AdminUserMembershipResponse>}
     * @memberof AdminUserDetailResponse
     */
    memberships?: Array<AdminUserMembershipResponse>;
    /**
     * 
     * @type {Date}
     * @memberof AdminUserDetailResponse
     */
    created_at?: Date;
}
/**
 * 
 * @export
 * @interface AdminUserListEnvelope
 */
export interface AdminUserListEnvelope {
    /**
     * 
     * @type {Meta}
     * @memberof AdminUserListEnvelope
     */
    meta: Meta;
    /**
     * 
     * @type {Array<AdminUserResponse>}
     * @memberof AdminUserListEnvelope
     */
    data: Array<AdminUserResponse>;
}
/**
 * 
 * @export
 * @interface AdminUserMembershipResponse
 */
export interface AdminUserMembershipResponse {
    /**
     * 
     * @type {string}
     * @memberof AdminUserMembershipResponse
     */
    tenant_id?: string;
    /**
     * 
     * @type {string}
     * @memberof AdminUserMembershipResponse
     */
    tenant_name?: string;
    /**
     * 
     * @type {string}
     * @memberof AdminUserMembershipResponse
     */
    role?: string;
    /**
     * 
     * @type {string}
     * @memberof AdminUserMembershipResponse
     */
    status?: string;
}
/**
 * 
 * @export
 * @interface AdminUserResponse
 */
export interface AdminUserResponse {
    /**
     * 
     * @type {string}
     * @memberof AdminUserResponse
     */
    id?: string;
    /**
     * 
     * @type {string}
     * @memberof AdminUserResponse
     */
    email?: string;
    /**
     * 
     * @type {string}
     * @memberof AdminUserResponse
     */
    first_name?: string;
    /**
     * 
     * @type {string}
     * @memberof AdminUserResponse
     */
    last_name?: string;
    /**
     * 
     * @type {string}
     * @memberof AdminUserResponse
     */
    status?: string;
    /**
     * 
     * @type {string}
     * @memberof AdminUserResponse
     */
    tenant_id?: string;
    /**
     * 
     * @type {boolean}
     * @memberof AdminUserResponse
     */
    mfa_enabled?: boolean;
    /**
     * 
     * @type {Date}
     * @memberof AdminUserResponse
     */
    created_at?: Date;
}
/**
 * 
 * @export
 * @interface DatabaseHealthEnvelope
 */
export interface DatabaseHealthEnvelope {
    /**
     * 
     * @type {Meta}
     * @memberof DatabaseHealthEnvelope
     */
    meta: Meta;
    /**
     * 
     * @type {{ [key: string]: any; }}
     * @memberof DatabaseHealthEnvelope
     */
    data: { [key: string]: any; };
}
/**
 * 
 * @export
 * @interface DatabaseMetricsEnvelope
 */
export interface DatabaseMetricsEnvelope {
    /**
     * 
     * @type {Meta}
     * @memberof DatabaseMetricsEnvelope
     */
    meta: Meta;
    /**
     * 
     * @type {{ [key: string]: any; }}
     * @memberof DatabaseMetricsEnvelope
     */
    data: { [key: string]: any; };
}
/**
 * 
 * @export
 * @interface EmptyDataEnvelope
 */
export interface EmptyDataEnvelope {
    /**
     * 
     * @type {Meta}
     * @memberof EmptyDataEnvelope
     */
    meta: Meta;
    /**
     * 
     * @type {{ [key: string]: any; }}
     * @memberof EmptyDataEnvelope
     */
    data: { [key: string]: any; };
}
/**
 * App API error response
 * @export
 * @interface ErrorEnvelope
 */
export interface ErrorEnvelope {
    /**
     * 
     * @type {Meta}
     * @memberof ErrorEnvelope
     */
    meta: Meta;
    /**
     * 
     * @type {any}
     * @memberof ErrorEnvelope
     */
    data?: any | null;
}
/**
 * 
 * @export
 * @interface HealthEnvelope
 */
export interface HealthEnvelope {
    /**
     * 
     * @type {Meta}
     * @memberof HealthEnvelope
     */
    meta: Meta;
    /**
     * 
     * @type {HealthResponse}
     * @memberof HealthEnvelope
     */
    data: HealthResponse;
}
/**
 * 
 * @export
 * @interface HealthResponse
 */
export interface HealthResponse {
    /**
     * 
     * @type {string}
     * @memberof HealthResponse
     */
    service?: string;
    /**
     * 
     * @type {string}
     * @memberof HealthResponse
     */
    status?: string;
    /**
     * 
     * @type {string}
     * @memberof HealthResponse
     */
    timestamp?: string;
    /**
     * 
     * @type {string}
     * @memberof HealthResponse
     */
    version?: string;
}
/**
 * JSON Web Key Set
 * @export
 * @interface JWKS
 */
export interface JWKS {
    /**
     * 
     * @type {Array<{ [key: string]: any; }>}
     * @memberof JWKS
     */
    keys?: Array<{ [key: string]: any; }>;
}
/**
 * 
 * @export
 * @interface LoginRequest
 */
export interface LoginRequest {
    /**
     * 
     * @type {string}
     * @memberof LoginRequest
     */
    email: string;
    /**
     * 
     * @type {string}
     * @memberof LoginRequest
     */
    password: string;
    /**
     * Optional explicit tenant context
     * @type {string}
     * @memberof LoginRequest
     */
    tenant_id?: string;
    /**
     * OIDC browser flow: creates session cookie and redirects here with 302
     * @type {string}
     * @memberof LoginRequest
     */
    return_to?: string;
    /**
     * Extends iam_session cookie lifetime (SSO_SESSION_REMEMBER_TTL)
     * @type {boolean}
     * @memberof LoginRequest
     */
    remember_me?: boolean;
}
/**
 * 
 * @export
 * @interface LoginResponse
 */
export interface LoginResponse {
    /**
     * 
     * @type {string}
     * @memberof LoginResponse
     */
    access_token?: string;
    /**
     * 
     * @type {string}
     * @memberof LoginResponse
     */
    refresh_token?: string;
    /**
     * 
     * @type {string}
     * @memberof LoginResponse
     */
    token_type?: string;
    /**
     * Access token lifetime in seconds
     * @type {number}
     * @memberof LoginResponse
     */
    expires_in?: number;
    /**
     * Refresh token lifetime in seconds
     * @type {number}
     * @memberof LoginResponse
     */
    refresh_expires_in?: number;
    /**
     * 
     * @type {string}
     * @memberof LoginResponse
     */
    active_tenant_id?: string;
}
/**
 * 
 * @export
 * @interface LoginResponseEnvelope
 */
export interface LoginResponseEnvelope {
    /**
     * 
     * @type {Meta}
     * @memberof LoginResponseEnvelope
     */
    meta: Meta;
    /**
     * 
     * @type {LoginResponse}
     * @memberof LoginResponseEnvelope
     */
    data: LoginResponse;
}
/**
 * @type LoginResult
 * Discriminated by properties: `access_token` → LoginResponse; `mfa_required` true → MFALoginChallengeResponse; `selection_required` true → TenantSelectionResponse
 * @export
 */
export type LoginResult = LoginResponse | MFALoginChallengeResponse | TenantSelectionResponse;
/**
 * 
 * @export
 * @interface LoginResultEnvelope
 */
export interface LoginResultEnvelope {
    /**
     * 
     * @type {Meta}
     * @memberof LoginResultEnvelope
     */
    meta: Meta;
    /**
     * 
     * @type {LoginResult}
     * @memberof LoginResultEnvelope
     */
    data: LoginResult;
}
/**
 * 
 * @export
 * @interface MFAChallengeVerifyRequest
 */
export interface MFAChallengeVerifyRequest {
    /**
     * 
     * @type {string}
     * @memberof MFAChallengeVerifyRequest
     */
    code: string;
    /**
     * 
     * @type {string}
     * @memberof MFAChallengeVerifyRequest
     */
    mfa_ticket: string;
}
/**
 * 
 * @export
 * @interface MFALoginChallengeResponse
 */
export interface MFALoginChallengeResponse {
    /**
     * 
     * @type {boolean}
     * @memberof MFALoginChallengeResponse
     */
    mfa_required: boolean;
    /**
     * 
     * @type {string}
     * @memberof MFALoginChallengeResponse
     */
    mfa_ticket: string;
    /**
     * 
     * @type {number}
     * @memberof MFALoginChallengeResponse
     */
    expires_in: number;
}
/**
 * 
 * @export
 * @interface MFARecoveryCodesEnvelope
 */
export interface MFARecoveryCodesEnvelope {
    /**
     * 
     * @type {Meta}
     * @memberof MFARecoveryCodesEnvelope
     */
    meta: Meta;
    /**
     * 
     * @type {MFARecoveryCodesResponse}
     * @memberof MFARecoveryCodesEnvelope
     */
    data: MFARecoveryCodesResponse;
}
/**
 * 
 * @export
 * @interface MFARecoveryCodesResponse
 */
export interface MFARecoveryCodesResponse {
    /**
     * 
     * @type {Array<string>}
     * @memberof MFARecoveryCodesResponse
     */
    codes?: Array<string>;
}
/**
 * 
 * @export
 * @interface MFATOTPSetupEnvelope
 */
export interface MFATOTPSetupEnvelope {
    /**
     * 
     * @type {Meta}
     * @memberof MFATOTPSetupEnvelope
     */
    meta: Meta;
    /**
     * 
     * @type {MFATOTPSetupResponse}
     * @memberof MFATOTPSetupEnvelope
     */
    data: MFATOTPSetupResponse;
}
/**
 * 
 * @export
 * @interface MFATOTPSetupResponse
 */
export interface MFATOTPSetupResponse {
    /**
     * 
     * @type {string}
     * @memberof MFATOTPSetupResponse
     */
    otpauth_uri?: string;
    /**
     * 
     * @type {string}
     * @memberof MFATOTPSetupResponse
     */
    secret?: string;
}
/**
 * 
 * @export
 * @interface MFATOTPVerifyRequest
 */
export interface MFATOTPVerifyRequest {
    /**
     * 
     * @type {string}
     * @memberof MFATOTPVerifyRequest
     */
    code: string;
}
/**
 * Response metadata (pagination and/or error details)
 * @export
 * @interface Meta
 */
export interface Meta {
    /**
     * 
     * @type {string}
     * @memberof Meta
     */
    error_code?: string;
    /**
     * 
     * @type {string}
     * @memberof Meta
     */
    message?: string;
    /**
     * 
     * @type {number}
     * @memberof Meta
     */
    code?: number;
    /**
     * 
     * @type {number}
     * @memberof Meta
     */
    page?: number;
    /**
     * 
     * @type {number}
     * @memberof Meta
     */
    page_size?: number;
    /**
     * 
     * @type {number}
     * @memberof Meta
     */
    total?: number;
}
/**
 * 
 * @export
 * @interface MetaOnlyEnvelope
 */
export interface MetaOnlyEnvelope {
    /**
     * 
     * @type {Meta}
     * @memberof MetaOnlyEnvelope
     */
    meta: Meta;
}
/**
 * OAuth 2.0 / OIDC error body for /token and /userinfo
 * @export
 * @interface OAuthError
 */
export interface OAuthError {
    /**
     * 
     * @type {string}
     * @memberof OAuthError
     */
    error: string;
    /**
     * 
     * @type {string}
     * @memberof OAuthError
     */
    error_description?: string;
}
/**
 * 
 * @export
 * @interface OIDCTokenResponse
 */
export interface OIDCTokenResponse {
    /**
     * 
     * @type {string}
     * @memberof OIDCTokenResponse
     */
    access_token?: string;
    /**
     * 
     * @type {string}
     * @memberof OIDCTokenResponse
     */
    refresh_token?: string;
    /**
     * 
     * @type {string}
     * @memberof OIDCTokenResponse
     */
    id_token?: string;
    /**
     * 
     * @type {string}
     * @memberof OIDCTokenResponse
     */
    token_type?: string;
    /**
     * 
     * @type {number}
     * @memberof OIDCTokenResponse
     */
    expires_in?: number;
    /**
     * 
     * @type {string}
     * @memberof OIDCTokenResponse
     */
    scope?: string;
}
/**
 * 
 * @export
 * @interface OpenIDConfigurationResponse
 */
export interface OpenIDConfigurationResponse {
    /**
     * 
     * @type {string}
     * @memberof OpenIDConfigurationResponse
     */
    issuer?: string;
    /**
     * 
     * @type {string}
     * @memberof OpenIDConfigurationResponse
     */
    authorization_endpoint?: string;
    /**
     * 
     * @type {string}
     * @memberof OpenIDConfigurationResponse
     */
    token_endpoint?: string;
    /**
     * 
     * @type {string}
     * @memberof OpenIDConfigurationResponse
     */
    userinfo_endpoint?: string;
    /**
     * 
     * @type {string}
     * @memberof OpenIDConfigurationResponse
     */
    jwks_uri?: string;
    /**
     * 
     * @type {Array<string>}
     * @memberof OpenIDConfigurationResponse
     */
    response_types_supported?: Array<string>;
    /**
     * 
     * @type {Array<string>}
     * @memberof OpenIDConfigurationResponse
     */
    subject_types_supported?: Array<string>;
    /**
     * 
     * @type {Array<string>}
     * @memberof OpenIDConfigurationResponse
     */
    id_token_signing_alg_values_supported?: Array<string>;
    /**
     * 
     * @type {Array<string>}
     * @memberof OpenIDConfigurationResponse
     */
    scopes_supported?: Array<string>;
    /**
     * 
     * @type {Array<string>}
     * @memberof OpenIDConfigurationResponse
     */
    token_endpoint_auth_methods_supported?: Array<string>;
    /**
     * 
     * @type {Array<string>}
     * @memberof OpenIDConfigurationResponse
     */
    code_challenge_methods_supported?: Array<string>;
    /**
     * 
     * @type {Array<string>}
     * @memberof OpenIDConfigurationResponse
     */
    grant_types_supported?: Array<string>;
}
/**
 * 
 * @export
 * @interface PatchIdentityProviderRequest
 */
export interface PatchIdentityProviderRequest {
    /**
     * 
     * @type {boolean}
     * @memberof PatchIdentityProviderRequest
     */
    enabled?: boolean;
    /**
     * 
     * @type {string}
     * @memberof PatchIdentityProviderRequest
     */
    oauth_client_id?: string;
    /**
     * 
     * @type {string}
     * @memberof PatchIdentityProviderRequest
     */
    oauth_client_secret?: string;
}
/**
 * 
 * @export
 * @interface PublicFederationProviderListEnvelope
 */
export interface PublicFederationProviderListEnvelope {
    /**
     * 
     * @type {Meta}
     * @memberof PublicFederationProviderListEnvelope
     */
    meta: Meta;
    /**
     * 
     * @type {Array<PublicFederationProviderResponse>}
     * @memberof PublicFederationProviderListEnvelope
     */
    data: Array<PublicFederationProviderResponse>;
}
/**
 * 
 * @export
 * @interface PublicFederationProviderResponse
 */
export interface PublicFederationProviderResponse {
    /**
     * 
     * @type {string}
     * @memberof PublicFederationProviderResponse
     */
    name?: string;
    /**
     * 
     * @type {string}
     * @memberof PublicFederationProviderResponse
     */
    provider?: string;
}
/**
 * 
 * @export
 * @interface RefreshTokenRequest
 */
export interface RefreshTokenRequest {
    /**
     * 
     * @type {string}
     * @memberof RefreshTokenRequest
     */
    refresh_token: string;
}
/**
 * 
 * @export
 * @interface RegisterRequest
 */
export interface RegisterRequest {
    /**
     * 
     * @type {string}
     * @memberof RegisterRequest
     */
    email: string;
    /**
     * 
     * @type {string}
     * @memberof RegisterRequest
     */
    password: string;
    /**
     * 
     * @type {string}
     * @memberof RegisterRequest
     */
    first_name?: string;
    /**
     * 
     * @type {string}
     * @memberof RegisterRequest
     */
    last_name?: string;
    /**
     * 
     * @type {string}
     * @memberof RegisterRequest
     */
    tenant_id?: string;
}
/**
 * 
 * @export
 * @interface TenantSelectRequest
 */
export interface TenantSelectRequest {
    /**
     * 
     * @type {string}
     * @memberof TenantSelectRequest
     */
    selection_token: string;
    /**
     * 
     * @type {string}
     * @memberof TenantSelectRequest
     */
    tenant_id: string;
    /**
     * 
     * @type {boolean}
     * @memberof TenantSelectRequest
     */
    remember_me?: boolean;
}
/**
 * 
 * @export
 * @interface TenantSelectionResponse
 */
export interface TenantSelectionResponse {
    /**
     * 
     * @type {boolean}
     * @memberof TenantSelectionResponse
     */
    selection_required: boolean;
    /**
     * 
     * @type {Array<TenantSummary>}
     * @memberof TenantSelectionResponse
     */
    tenants: Array<TenantSummary>;
    /**
     * 
     * @type {string}
     * @memberof TenantSelectionResponse
     */
    selection_token: string;
    /**
     * 
     * @type {number}
     * @memberof TenantSelectionResponse
     */
    expires_in: number;
}
/**
 * 
 * @export
 * @interface TenantSummary
 */
export interface TenantSummary {
    /**
     * 
     * @type {string}
     * @memberof TenantSummary
     */
    id?: string;
    /**
     * 
     * @type {string}
     * @memberof TenantSummary
     */
    name?: string;
    /**
     * 
     * @type {string}
     * @memberof TenantSummary
     */
    domain?: string;
    /**
     * 
     * @type {string}
     * @memberof TenantSummary
     */
    role?: string;
}
/**
 * 
 * @export
 * @interface TenantSummaryListEnvelope
 */
export interface TenantSummaryListEnvelope {
    /**
     * 
     * @type {Meta}
     * @memberof TenantSummaryListEnvelope
     */
    meta: Meta;
    /**
     * 
     * @type {Array<TenantSummary>}
     * @memberof TenantSummaryListEnvelope
     */
    data: Array<TenantSummary>;
}
/**
 * 
 * @export
 * @interface TenantSwitchRequest
 */
export interface TenantSwitchRequest {
    /**
     * 
     * @type {string}
     * @memberof TenantSwitchRequest
     */
    tenant_id: string;
}
/**
 * 
 * @export
 * @interface UpdateProfileRequest
 */
export interface UpdateProfileRequest {
    /**
     * 
     * @type {string}
     * @memberof UpdateProfileRequest
     */
    first_name?: string;
    /**
     * 
     * @type {string}
     * @memberof UpdateProfileRequest
     */
    last_name?: string;
}
/**
 * OIDC UserInfo claims
 * @export
 * @interface UserInfoResponse
 */
export interface UserInfoResponse {
    [key: string]: any | any;
    /**
     * 
     * @type {string}
     * @memberof UserInfoResponse
     */
    sub?: string;
    /**
     * 
     * @type {string}
     * @memberof UserInfoResponse
     */
    email?: string;
    /**
     * 
     * @type {boolean}
     * @memberof UserInfoResponse
     */
    email_verified?: boolean;
    /**
     * 
     * @type {string}
     * @memberof UserInfoResponse
     */
    name?: string;
    /**
     * 
     * @type {string}
     * @memberof UserInfoResponse
     */
    given_name?: string;
    /**
     * 
     * @type {string}
     * @memberof UserInfoResponse
     */
    family_name?: string;
}
/**
 * 
 * @export
 * @interface UserResponse
 */
export interface UserResponse {
    /**
     * 
     * @type {string}
     * @memberof UserResponse
     */
    id?: string;
    /**
     * 
     * @type {string}
     * @memberof UserResponse
     */
    email?: string;
    /**
     * 
     * @type {string}
     * @memberof UserResponse
     */
    first_name?: string;
    /**
     * 
     * @type {string}
     * @memberof UserResponse
     */
    last_name?: string;
    /**
     * 
     * @type {boolean}
     * @memberof UserResponse
     */
    email_verified?: boolean;
    /**
     * 
     * @type {boolean}
     * @memberof UserResponse
     */
    is_platform_admin?: boolean;
    /**
     * 
     * @type {boolean}
     * @memberof UserResponse
     */
    mfa_enabled?: boolean;
    /**
     * 
     * @type {string}
     * @memberof UserResponse
     */
    active_tenant_id?: string;
    /**
     * 
     * @type {Array<TenantSummary>}
     * @memberof UserResponse
     */
    tenants?: Array<TenantSummary>;
    /**
     * 
     * @type {Date}
     * @memberof UserResponse
     */
    created_at?: Date;
    /**
     * 
     * @type {Date}
     * @memberof UserResponse
     */
    updated_at?: Date;
}
/**
 * 
 * @export
 * @interface UserResponseEnvelope
 */
export interface UserResponseEnvelope {
    /**
     * 
     * @type {Meta}
     * @memberof UserResponseEnvelope
     */
    meta: Meta;
    /**
     * 
     * @type {UserResponse}
     * @memberof UserResponseEnvelope
     */
    data: UserResponse;
}
/**
 * 
 * @export
 * @interface WebauthnCredentialListEnvelope
 */
export interface WebauthnCredentialListEnvelope {
    /**
     * 
     * @type {Meta}
     * @memberof WebauthnCredentialListEnvelope
     */
    meta: Meta;
    /**
     * 
     * @type {Array<WebauthnCredentialResponse>}
     * @memberof WebauthnCredentialListEnvelope
     */
    data: Array<WebauthnCredentialResponse>;
}
/**
 * 
 * @export
 * @interface WebauthnCredentialResponse
 */
export interface WebauthnCredentialResponse {
    /**
     * 
     * @type {string}
     * @memberof WebauthnCredentialResponse
     */
    id?: string;
    /**
     * 
     * @type {string}
     * @memberof WebauthnCredentialResponse
     */
    device_name?: string;
    /**
     * 
     * @type {Date}
     * @memberof WebauthnCredentialResponse
     */
    created_at?: Date;
}
/**
 * 
 * @export
 * @interface WebauthnLoginFinishRequest
 */
export interface WebauthnLoginFinishRequest {
    /**
     * 
     * @type {string}
     * @memberof WebauthnLoginFinishRequest
     */
    email: string;
    /**
     * 
     * @type {string}
     * @memberof WebauthnLoginFinishRequest
     */
    session_token: string;
    /**
     * PublicKeyCredential JSON from navigator.credentials.get()
     * @type {{ [key: string]: any; }}
     * @memberof WebauthnLoginFinishRequest
     */
    credential?: { [key: string]: any; };
    /**
     * 
     * @type {string}
     * @memberof WebauthnLoginFinishRequest
     */
    tenant_id?: string;
    /**
     * 
     * @type {boolean}
     * @memberof WebauthnLoginFinishRequest
     */
    remember_me?: boolean;
    /**
     * 
     * @type {string}
     * @memberof WebauthnLoginFinishRequest
     */
    return_to?: string;
}
/**
 * 
 * @export
 * @interface WebauthnLoginStartEnvelope
 */
export interface WebauthnLoginStartEnvelope {
    /**
     * 
     * @type {Meta}
     * @memberof WebauthnLoginStartEnvelope
     */
    meta: Meta;
    /**
     * 
     * @type {WebauthnLoginStartResponse}
     * @memberof WebauthnLoginStartEnvelope
     */
    data: WebauthnLoginStartResponse;
}
/**
 * 
 * @export
 * @interface WebauthnLoginStartRequest
 */
export interface WebauthnLoginStartRequest {
    /**
     * 
     * @type {string}
     * @memberof WebauthnLoginStartRequest
     */
    email: string;
    /**
     * 
     * @type {string}
     * @memberof WebauthnLoginStartRequest
     */
    tenant_id?: string;
}
/**
 * 
 * @export
 * @interface WebauthnLoginStartResponse
 */
export interface WebauthnLoginStartResponse {
    /**
     * 
     * @type {{ [key: string]: any; }}
     * @memberof WebauthnLoginStartResponse
     */
    options?: { [key: string]: any; };
    /**
     * 
     * @type {string}
     * @memberof WebauthnLoginStartResponse
     */
    session_token?: string;
}
/**
 * 
 * @export
 * @interface WebauthnRegisterFinishRequest
 */
export interface WebauthnRegisterFinishRequest {
    /**
     * 
     * @type {string}
     * @memberof WebauthnRegisterFinishRequest
     */
    session_token: string;
    /**
     * PublicKeyCredential JSON from navigator.credentials.create()
     * @type {{ [key: string]: any; }}
     * @memberof WebauthnRegisterFinishRequest
     */
    credential?: { [key: string]: any; };
}
/**
 * 
 * @export
 * @interface WebauthnRegisterStartEnvelope
 */
export interface WebauthnRegisterStartEnvelope {
    /**
     * 
     * @type {Meta}
     * @memberof WebauthnRegisterStartEnvelope
     */
    meta: Meta;
    /**
     * 
     * @type {WebauthnRegisterStartResponse}
     * @memberof WebauthnRegisterStartEnvelope
     */
    data: WebauthnRegisterStartResponse;
}
/**
 * 
 * @export
 * @interface WebauthnRegisterStartRequest
 */
export interface WebauthnRegisterStartRequest {
    /**
     * Optional label stored with the passkey
     * @type {string}
     * @memberof WebauthnRegisterStartRequest
     */
    device_name?: string;
}
/**
 * 
 * @export
 * @interface WebauthnRegisterStartResponse
 */
export interface WebauthnRegisterStartResponse {
    /**
     * 
     * @type {{ [key: string]: any; }}
     * @memberof WebauthnRegisterStartResponse
     */
    options?: { [key: string]: any; };
    /**
     * 
     * @type {string}
     * @memberof WebauthnRegisterStartResponse
     */
    session_token?: string;
}
