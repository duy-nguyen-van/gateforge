import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'

import {
  addTenantMember,
  createAdminTenant,
  createAdminClient,
  deleteAdminTenant,
  deleteAdminClient,
  disableAdminUser,
  forceLogoutAdminUser,
  getAdminClient,
  getAdminClientUsage,
  getAdminStats,
  getAdminTenant,
  getAdminUser,
  listAdminAuditLogs,
  listAdminClients,
  listAdminLoginHistory,
  listAdminTenants,
  listAdminUsers,
  listTenantIdentityProviders,
  listTenantMembers,
  patchIdentityProvider,
  removeTenantMember,
  resetAdminUserPasskeys,
  resetAdminUserMFA,
  updateAdminTenant,
  updateAdminClient,
} from '@/api/client'
import type {
  AdminAddMemberRequest,
  AdminAuditLogListParams,
  AdminCreateClientRequest,
  AdminCreateTenantRequest,
  AdminListParams,
  AdminLoginHistoryListParams,
  AdminUpdateClientRequest,
  AdminUpdateTenantRequest,
  PatchIdentityProviderRequest,
} from '@/api/types'

export const adminQueryKeys = {
  stats: ['admin', 'stats'] as const,
  users: (params: AdminListParams) => ['admin', 'users', params] as const,
  user: (userId: string) => ['admin', 'user', userId] as const,
  tenants: (params: AdminListParams) => ['admin', 'tenants', params] as const,
  tenant: (tenantId: string) => ['admin', 'tenant', tenantId] as const,
  tenantMembers: (tenantId: string, params: AdminListParams) => ['admin', 'tenant-members', tenantId, params] as const,
  clients: (params: AdminListParams) => ['admin', 'clients', params] as const,
  client: (clientId: string) => ['admin', 'client', clientId] as const,
  clientUsage: (clientId: string) => ['admin', 'client-usage', clientId] as const,
  auditLogs: (params: AdminAuditLogListParams) => ['admin', 'audit-logs', params] as const,
  loginHistory: (params: AdminLoginHistoryListParams) => ['admin', 'login-history', params] as const,
  identityProviders: (tenantId: string) => ['admin', 'identity-providers', tenantId] as const,
}

function invalidateUserAdminQueries(queryClient: ReturnType<typeof useQueryClient>, userId?: string) {
  void queryClient.invalidateQueries({ queryKey: ['admin', 'users'] })
  void queryClient.invalidateQueries({ queryKey: adminQueryKeys.stats })
  if (userId) {
    void queryClient.invalidateQueries({ queryKey: adminQueryKeys.user(userId) })
  }
}

export function useAdminStats() {
  return useQuery({
    queryKey: adminQueryKeys.stats,
    queryFn: () => getAdminStats(),
  })
}

export function useAdminUsers(params: AdminListParams = {}) {
  return useQuery({
    queryKey: adminQueryKeys.users(params),
    queryFn: () => listAdminUsers(params),
  })
}

export function useAdminUser(userId: string | null) {
  return useQuery({
    queryKey: adminQueryKeys.user(userId ?? ''),
    queryFn: () => getAdminUser(userId!),
    enabled: Boolean(userId),
  })
}

export function useAdminTenants(params: AdminListParams = {}) {
  return useQuery({
    queryKey: adminQueryKeys.tenants(params),
    queryFn: () => listAdminTenants(params),
  })
}

export function useAdminTenant(tenantId: string | null) {
  return useQuery({
    queryKey: adminQueryKeys.tenant(tenantId ?? ''),
    queryFn: () => getAdminTenant(tenantId!),
    enabled: Boolean(tenantId),
  })
}

export function useTenantMembers(tenantId: string | null, params: AdminListParams = {}) {
  return useQuery({
    queryKey: adminQueryKeys.tenantMembers(tenantId ?? '', params),
    queryFn: () => listTenantMembers(tenantId!, params),
    enabled: Boolean(tenantId),
  })
}

export function useCreateAdminTenant() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (body: AdminCreateTenantRequest) => createAdminTenant(body),
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: ['admin', 'tenants'] })
    },
  })
}

export function useUpdateAdminTenant(tenantId: string) {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (body: AdminUpdateTenantRequest) => updateAdminTenant(tenantId, body),
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: ['admin', 'tenants'] })
      void queryClient.invalidateQueries({ queryKey: adminQueryKeys.tenant(tenantId) })
    },
  })
}

export function useDeleteAdminTenant() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (tenantId: string) => deleteAdminTenant(tenantId),
    onSuccess: (_data, tenantId) => {
      void queryClient.invalidateQueries({ queryKey: ['admin', 'tenants'] })
      void queryClient.removeQueries({ queryKey: adminQueryKeys.tenant(tenantId) })
    },
  })
}

export function useAdminClients(params: AdminListParams = {}) {
  return useQuery({
    queryKey: adminQueryKeys.clients(params),
    queryFn: () => listAdminClients(params),
  })
}

export function useAdminClient(clientId: string | null) {
  return useQuery({
    queryKey: adminQueryKeys.client(clientId ?? ''),
    queryFn: () => getAdminClient(clientId!),
    enabled: Boolean(clientId),
  })
}

export function useCreateAdminClient() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (body: AdminCreateClientRequest) => createAdminClient(body),
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: ['admin', 'clients'] })
    },
  })
}

export function useUpdateAdminClient(clientId: string) {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (body: AdminUpdateClientRequest) => updateAdminClient(clientId, body),
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: ['admin', 'clients'] })
      void queryClient.invalidateQueries({ queryKey: adminQueryKeys.client(clientId) })
    },
  })
}

export function useDeleteAdminClient() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (clientId: string) => deleteAdminClient(clientId),
    onSuccess: (_data, clientId) => {
      void queryClient.invalidateQueries({ queryKey: ['admin', 'clients'] })
      void queryClient.removeQueries({ queryKey: adminQueryKeys.client(clientId) })
    },
  })
}

export function useAdminClientUsage(clientId: string | null) {
  return useQuery({
    queryKey: adminQueryKeys.clientUsage(clientId ?? ''),
    queryFn: () => getAdminClientUsage(clientId!),
    enabled: Boolean(clientId),
  })
}

export function useTenantIdentityProviders(tenantId: string) {
  return useQuery({
    queryKey: adminQueryKeys.identityProviders(tenantId),
    queryFn: () => listTenantIdentityProviders(tenantId),
    enabled: Boolean(tenantId),
  })
}

export function usePatchIdentityProvider(tenantId: string) {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: ({ provider, body }: { provider: string; body: PatchIdentityProviderRequest }) =>
      patchIdentityProvider(tenantId, provider, body),
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: adminQueryKeys.identityProviders(tenantId) })
    },
  })
}

export function useAddTenantMember() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: ({ tenantId, body }: { tenantId: string; body: AdminAddMemberRequest }) =>
      addTenantMember(tenantId, body),
    onSuccess: (_data, { tenantId }) => {
      void queryClient.invalidateQueries({ queryKey: ['admin', 'users'] })
      void queryClient.invalidateQueries({ queryKey: ['admin', 'tenants'] })
      void queryClient.invalidateQueries({ queryKey: adminQueryKeys.tenant(tenantId) })
      void queryClient.invalidateQueries({ queryKey: ['admin', 'tenant-members', tenantId] })
    },
  })
}

export function useRemoveTenantMember() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: ({ tenantId, userId }: { tenantId: string; userId: string }) =>
      removeTenantMember(tenantId, userId),
    onSuccess: (_data, { tenantId }) => {
      void queryClient.invalidateQueries({ queryKey: ['admin', 'users'] })
      void queryClient.invalidateQueries({ queryKey: ['admin', 'tenants'] })
      void queryClient.invalidateQueries({ queryKey: adminQueryKeys.tenant(tenantId) })
      void queryClient.invalidateQueries({ queryKey: ['admin', 'tenant-members', tenantId] })
    },
  })
}

export function useDisableAdminUser() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (userId: string) => disableAdminUser(userId),
    onSuccess: (_data, userId) => invalidateUserAdminQueries(queryClient, userId),
  })
}

export function useForceLogoutAdminUser() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (userId: string) => forceLogoutAdminUser(userId),
    onSuccess: (_data, userId) => invalidateUserAdminQueries(queryClient, userId),
  })
}

export function useResetAdminUserPasskeys() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (userId: string) => resetAdminUserPasskeys(userId),
    onSuccess: (_data, userId) => invalidateUserAdminQueries(queryClient, userId),
  })
}

export function useResetAdminUserMFA() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (userId: string) => resetAdminUserMFA(userId),
    onSuccess: (_data, userId) => invalidateUserAdminQueries(queryClient, userId),
  })
}

export function useAdminAuditLogs(params: AdminAuditLogListParams = {}) {
  return useQuery({
    queryKey: adminQueryKeys.auditLogs(params),
    queryFn: () => listAdminAuditLogs(params),
  })
}

export function useAdminLoginHistory(params: AdminLoginHistoryListParams = {}) {
  return useQuery({
    queryKey: adminQueryKeys.loginHistory(params),
    queryFn: () => listAdminLoginHistory(params),
  })
}
