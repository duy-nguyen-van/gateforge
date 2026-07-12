import { useEffect, useState } from 'react'

import { DefaultAvatar } from '@/components/avatars/default-avatar'
import { MaterialIcon } from '@/components/icons/material-icon'
import { AddMemberDialog } from '@/features/admin/add-member-dialog'
import { displayUserName, formatUserStatus } from '@/features/admin/admin-utils'
import { ConsolePagination } from '@/features/admin/console-pagination'
import {
  ConsoleEmptyState,
  ConsoleErrorState,
  ConsoleLoadingState,
} from '@/features/admin/console-state'
import { RemoveMemberDialog } from '@/features/admin/remove-member-dialog'
import { UserDetailDrawer } from '@/features/admin/user-detail-drawer'
import { useConsolePagination } from '@/features/admin/use-console-pagination'
import { useAdminStats, useAdminUsers } from '@/features/admin/use-admin-queries'

type RemoveTarget = {
  userId: string
  tenantId: string
  email: string
}

export function UsersPage() {
  const [search, setSearch] = useState('')
  const [addMemberOpen, setAddMemberOpen] = useState(false)
  const [removeTarget, setRemoveTarget] = useState<RemoveTarget | null>(null)
  const [detailUserId, setDetailUserId] = useState<string | null>(null)
  const { page, setPage, pageSize, resetPage, queryParams } = useConsolePagination()
  const statsQuery = useAdminStats()
  const usersQuery = useAdminUsers({ search, ...queryParams })

  useEffect(() => {
    resetPage()
  }, [search, resetPage])

  const stats = statsQuery.data?.data
  const users = usersQuery.data?.data ?? []
  const meta = usersQuery.data?.meta

  return (
    <div>
      <header className="mb-10 flex flex-col justify-between gap-6 md:flex-row md:items-end">
        <div>
          <h1 className="font-headline text-4xl font-extrabold tracking-tight text-on-surface">User Management</h1>
          <p className="mt-1 text-on-surface-variant">Manage identities across all tenants.</p>
        </div>
        <button
          type="button"
          title="User must already be registered"
          onClick={() => setAddMemberOpen(true)}
          className="flex items-center gap-2 rounded-xl bg-primary px-5 py-2.5 text-sm font-bold text-on-primary shadow-lg shadow-primary/20 transition-opacity hover:opacity-90"
        >
          <MaterialIcon name="person_add" className="text-sm" />
          Add member
        </button>
      </header>

      <div className="mb-8 grid grid-cols-1 gap-4 sm:grid-cols-3">
        {statsQuery.isLoading ? (
          <ConsoleLoadingState label="Loading stats…" />
        ) : statsQuery.isError ? (
          <div className="col-span-full">
            <ConsoleErrorState message="Could not load user statistics." />
          </div>
        ) : stats ? (
          [
            { label: 'Total Users', value: stats.total_users.toLocaleString(), icon: 'group' },
            { label: 'MFA Enabled', value: `${stats.mfa_enabled_percent}%`, icon: 'security' },
            { label: 'Active Sessions', value: stats.active_sessions.toLocaleString(), icon: 'person_check' },
          ].map((s) => (
            <div key={s.label} className="rounded-xl bg-surface-container-lowest p-6 ghost-border">
              <MaterialIcon name={s.icon} className="mb-3 text-primary text-2xl" />
              <p className="text-xs font-bold uppercase tracking-wider text-on-surface-variant">{s.label}</p>
              <p className="font-headline text-3xl font-extrabold text-on-primary-fixed">{s.value}</p>
            </div>
          ))
        ) : null}
      </div>

      <div className="overflow-hidden rounded-xl bg-surface-container-lowest ghost-border">
        <div className="flex items-center justify-between border-b border-surface-container px-6 py-4">
          <div className="relative w-full max-w-md">
            <MaterialIcon name="search" className="absolute left-3 top-1/2 -translate-y-1/2 text-outline text-sm" />
            <input
              type="text"
              placeholder="Search users..."
              value={search}
              onChange={(e) => setSearch(e.target.value)}
              className="w-full rounded-full border-none bg-surface-container-low py-2 pl-10 pr-4 text-sm focus:ring-1 focus:ring-primary"
            />
          </div>
          <MaterialIcon name="filter_list" className="text-on-surface-variant" />
        </div>

        {usersQuery.isLoading ? (
          <ConsoleLoadingState />
        ) : usersQuery.isError ? (
          <div className="p-6">
            <ConsoleErrorState message="Could not load users. Ensure your account has platform admin access." />
          </div>
        ) : users.length === 0 ? (
          <ConsoleEmptyState title="No users found" description="Try a different search or register the first user." />
        ) : (
          <table className="w-full text-left text-sm">
            <thead>
              <tr className="bg-surface-container-low">
                {['User', 'Status', 'MFA', 'Tenant', 'Actions'].map((h) => (
                  <th key={h} className="px-6 py-3 text-[10px] font-bold uppercase tracking-widest text-on-surface-variant">
                    {h}
                  </th>
                ))}
              </tr>
            </thead>
            <tbody>
              {users.map((u) => {
                const name = displayUserName(u.first_name, u.last_name, u.email)
                const isActive = u.status === 'active'
                return (
                  <tr
                    key={`${u.id}-${u.tenant_id}`}
                    className="border-b border-surface-container transition-colors hover:bg-surface-container-low"
                  >
                    <td className="px-6 py-4">
                      <div className="flex items-center gap-3">
                        <DefaultAvatar seed={u.email} name={name} size="md" title={name} />
                        <div>
                          <p className="font-semibold">{name}</p>
                          <p className="text-xs text-on-surface-variant">{u.email}</p>
                        </div>
                      </div>
                    </td>
                    <td className="px-6 py-4">
                      <span
                        className={`rounded-full px-2 py-0.5 text-xs font-bold ${
                          isActive ? 'bg-green-50 text-green-700' : 'bg-error-container/30 text-error'
                        }`}
                      >
                        {formatUserStatus(u.status)}
                      </span>
                    </td>
                    <td className="px-6 py-4">
                      <MaterialIcon
                        name={u.mfa_enabled ? 'check_circle' : 'cancel'}
                        className={u.mfa_enabled ? 'text-green-600' : 'text-outline'}
                      />
                    </td>
                    <td className="px-6 py-4 font-mono text-xs text-on-surface-variant">{u.tenant_id.slice(0, 8)}…</td>
                    <td className="px-6 py-4">
                      <div className="flex flex-wrap items-center gap-3">
                        <button
                          type="button"
                          onClick={() => setDetailUserId(u.id)}
                          className="text-xs font-bold text-primary hover:underline"
                        >
                          View
                        </button>
                        {u.tenant_id ? (
                          <button
                            type="button"
                            onClick={() =>
                              setRemoveTarget({ userId: u.id, tenantId: u.tenant_id, email: u.email })
                            }
                            className="text-xs font-bold text-error hover:underline"
                          >
                            Remove from tenant
                          </button>
                        ) : null}
                      </div>
                    </td>
                  </tr>
                )
              })}
            </tbody>
          </table>
        )}
        {meta?.total != null && meta.total > 0 ? (
          <ConsolePagination
            page={meta.page ?? page}
            pageSize={meta.page_size ?? pageSize}
            total={meta.total}
            onPageChange={setPage}
          />
        ) : null}
      </div>

      <AddMemberDialog open={addMemberOpen} onOpenChange={setAddMemberOpen} />

      {removeTarget ? (
        <RemoveMemberDialog
          open
          email={removeTarget.email}
          tenantId={removeTarget.tenantId}
          userId={removeTarget.userId}
          onOpenChange={(open) => {
            if (!open) {
              setRemoveTarget(null)
            }
          }}
        />
      ) : null}

      <UserDetailDrawer userId={detailUserId} onClose={() => setDetailUserId(null)} />
    </div>
  )
}
