import type { ReactNode } from 'react'

import { ApiError } from '@/api/types'
import { ConsolePagination } from '@/features/admin/console-pagination'
import { ConsoleEmptyState, ConsoleErrorState, ConsoleTableSkeleton } from '@/features/admin/console-state'
import { useConsolePagination } from '@/features/admin/use-console-pagination'
import { useAdminInvites, useResendTenantInvite, useTenantInvites } from '@/features/admin/use-admin-queries'

function resendErrorMessage(error: unknown) {
  if (error instanceof ApiError) {
    return error.message
  }
  if (error) {
    return 'Could not resend the invitation.'
  }
  return null
}

export function InvitationsSection({ tenantId }: Readonly<{ tenantId?: string }>) {
  const { page, setPage, pageSize, queryParams } = useConsolePagination()
  const scoped = Boolean(tenantId)
  const tenantQuery = useTenantInvites(scoped ? tenantId! : null, queryParams)
  const allQuery = useAdminInvites(queryParams, !scoped)
  const invitesQuery = scoped ? tenantQuery : allQuery
  const resendInvite = useResendTenantInvite()
  const invites = invitesQuery.data?.data ?? []
  const inviteMeta = invitesQuery.data?.meta
  const resendError = resendErrorMessage(resendInvite.error)
  const columns = scoped
    ? ['Email', 'Role', 'Status', 'Sent', 'Expires', 'Actions']
    : ['Email', 'Organization', 'Role', 'Status', 'Sent', 'Expires', 'Actions']

  let body: ReactNode
  if (invitesQuery.isLoading) {
    body = <ConsoleTableSkeleton columns={columns.length} />
  } else if (invitesQuery.isError) {
    body = (
      <div className="p-6">
        <ConsoleErrorState message="Could not load invitations." />
      </div>
    )
  } else if (invites.length === 0) {
    body = (
      <ConsoleEmptyState
        title="No pending invitations"
        description="Add a member by email to send an invitation. It stays here until they accept."
      />
    )
  } else {
    body = (
      <table className="w-full text-left text-sm">
        <thead>
          <tr className="bg-surface-container-low">
            {columns.map((h) => (
              <th
                key={h}
                className="px-6 py-4 text-[10px] font-bold uppercase tracking-widest text-on-surface-variant"
              >
                {h}
              </th>
            ))}
          </tr>
        </thead>
        <tbody>
          {invites.map((invite) => {
            const sending = resendInvite.isPending && resendInvite.variables?.inviteId === invite.id
            return (
              <tr
                key={invite.id}
                className="border-b border-surface-container transition-colors hover:bg-surface-container-low"
              >
                <td className="px-6 py-4 font-medium">{invite.email}</td>
                {scoped ? null : (
                  <td className="px-6 py-4 text-on-surface-variant">{invite.tenant_name || '—'}</td>
                )}
                <td className="px-6 py-4 capitalize">{invite.role}</td>
                <td className="px-6 py-4 capitalize">{invite.status}</td>
                <td className="px-6 py-4 text-on-surface-variant">
                  {new Date(invite.created_at).toLocaleString()}
                </td>
                <td className="px-6 py-4 text-on-surface-variant">
                  {new Date(invite.expires_at).toLocaleString()}
                </td>
                <td className="px-6 py-4">
                  <button
                    type="button"
                    disabled={sending}
                    onClick={() => {
                      resendInvite.mutate({ tenantId: invite.tenant_id, inviteId: invite.id })
                    }}
                    className="text-sm font-semibold text-primary disabled:opacity-50"
                  >
                    {sending ? 'Sending…' : 'Resend'}
                  </button>
                </td>
              </tr>
            )
          })}
        </tbody>
      </table>
    )
  }

  return (
    <section className="mt-10">
      <h2 className="mb-1 font-headline text-xl font-bold text-on-surface">Invitations</h2>
      <p className="mb-4 text-sm text-on-surface-variant">
        People who were emailed a sign-in link and have not accepted yet. Resend issues a new link.
      </p>
      <div className="overflow-hidden rounded-xl bg-surface-container-lowest ghost-border">
        {resendError ? (
          <div className="p-6 pb-0">
            <ConsoleErrorState message={resendError} />
          </div>
        ) : null}
        {body}
        {inviteMeta?.total != null && inviteMeta.total > 0 ? (
          <ConsolePagination
            page={inviteMeta.page ?? page}
            pageSize={inviteMeta.page_size ?? pageSize}
            total={inviteMeta.total}
            onPageChange={setPage}
          />
        ) : null}
      </div>
    </section>
  )
}
