import { useEffect, useRef, useState } from 'react'
import type { LucideIcon } from 'lucide-react'
import {
  CheckCircle2Icon,
  CircleXIcon,
  SearchIcon,
  ShieldIcon,
  UserCheckIcon,
  UserPlusIcon,
  UsersIcon,
} from 'lucide-react'

import { DefaultAvatar } from '@/components/avatars/default-avatar'
import { ConsolePageHeader } from '@/components/layout/console-page-header'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { AddMemberDialog } from '@/features/admin/add-member-dialog'
import { InvitationsSection } from '@/features/admin/invitations-section'
import { displayUserName, formatUserStatus } from '@/features/admin/admin-utils'
import { ConsolePagination } from '@/features/admin/console-pagination'
import {
  ConsoleEmptyState,
  ConsoleErrorState,
  ConsoleTableSkeleton,
} from '@/features/admin/console-state'
import { useSlashFocus } from '@/features/admin/use-slash-focus'
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
  const searchRef = useRef<HTMLInputElement>(null)
  useSlashFocus(searchRef)
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

  const statCards: { label: string; value: string; icon: LucideIcon }[] = stats
    ? [
        { label: 'Total Users', value: stats.total_users.toLocaleString(), icon: UsersIcon },
        { label: 'MFA Enabled', value: `${stats.mfa_enabled_percent}%`, icon: ShieldIcon },
        { label: 'Active Sessions', value: stats.active_sessions.toLocaleString(), icon: UserCheckIcon },
      ]
    : []

  return (
    <div>
      <ConsolePageHeader
        title="Users"
        description="Manage identities across all tenants."
        actions={
          <Button type="button" onClick={() => setAddMemberOpen(true)}>
            <UserPlusIcon className="h-4 w-4" aria-hidden />
            Add member
          </Button>
        }
      />

      <div className="mb-8 grid grid-cols-1 gap-4 sm:grid-cols-3">
        {statsQuery.isLoading ? (
          <p className="text-sm text-on-surface-variant">Loading stats…</p>
        ) : statsQuery.isError ? (
          <div className="col-span-full">
            <ConsoleErrorState message="Could not load user statistics." />
          </div>
        ) : stats ? (
          statCards.map((s) => {
            const Icon = s.icon
            return (
              <div key={s.label} className="rounded-xl bg-surface-container-lowest p-6 ghost-border">
                <Icon className="mb-3 h-6 w-6 text-primary" aria-hidden />
                <p className="text-xs font-bold uppercase tracking-wider text-on-surface-variant">{s.label}</p>
                <p className="font-headline text-3xl font-extrabold text-on-primary-fixed">{s.value}</p>
              </div>
            )
          })
        ) : null}
      </div>

      <div className="overflow-hidden rounded-xl bg-surface-container-lowest ghost-border">
        <div className="flex items-center justify-between border-b border-surface-container px-6 py-4">
          <div className="relative w-full max-w-md">
            <SearchIcon className="pointer-events-none absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-on-surface-variant" aria-hidden />
            <Input
              ref={searchRef}
              type="search"
              aria-label="Search users"
              placeholder="Search users"
              value={search}
              onChange={(e) => setSearch(e.target.value)}
              className="pl-10"
            />
          </div>
        </div>

        {usersQuery.isLoading ? (
          <ConsoleTableSkeleton columns={5} />
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
                          isActive ? 'bg-success-container text-on-success-container' : 'bg-error-container text-error'
                        }`}
                      >
                        {formatUserStatus(u.status)}
                      </span>
                    </td>
                    <td className="px-6 py-4">
                      {u.mfa_enabled ? (
                        <CheckCircle2Icon className="h-5 w-5 text-success" aria-hidden />
                      ) : (
                        <CircleXIcon className="h-5 w-5 text-outline" aria-hidden />
                      )}
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

      <InvitationsSection />

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
