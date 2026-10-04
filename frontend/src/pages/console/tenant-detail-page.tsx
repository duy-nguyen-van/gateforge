import { ArrowLeftIcon, PencilIcon, Trash2Icon, UserMinusIcon, UserPlusIcon } from 'lucide-react'
import { useState } from 'react'
import { Link, useNavigate, useParams } from 'react-router'

import { ConsolePageHeader } from '@/components/layout/console-page-header'
import { useDetailCrumb } from '@/components/layout/use-detail-crumb'
import { Button } from '@/components/ui/button'
import { AddMemberDialog } from '@/features/admin/add-member-dialog'
import { displayUserName } from '@/features/admin/admin-utils'
import { ConsolePagination } from '@/features/admin/console-pagination'
import { ConsoleEmptyState, ConsoleErrorState, ConsoleLoadingState, ConsoleTableSkeleton } from '@/features/admin/console-state'
import { DeleteTenantDialog } from '@/features/admin/delete-tenant-dialog'
import { EditTenantDialog } from '@/features/admin/edit-tenant-dialog'
import { RemoveMemberDialog } from '@/features/admin/remove-member-dialog'
import { InvitationsSection } from '@/features/admin/invitations-section'
import { useConsolePagination } from '@/features/admin/use-console-pagination'
import { useAdminTenant, useTenantMembers } from '@/features/admin/use-admin-queries'

const defaultTenantId = import.meta.env.VITE_DEFAULT_TENANT_ID ?? ''

type RemoveTarget = {
  userId: string
  email: string
}

export function TenantDetailPage() {
  const { tenantId = '' } = useParams()
  const navigate = useNavigate()
  const { page, setPage, pageSize, queryParams } = useConsolePagination()
  const tenantQuery = useAdminTenant(tenantId)
  const membersQuery = useTenantMembers(tenantId, queryParams)

  const [editOpen, setEditOpen] = useState(false)
  const [deleteOpen, setDeleteOpen] = useState(false)
  const [addMemberOpen, setAddMemberOpen] = useState(false)
  const [removeTarget, setRemoveTarget] = useState<RemoveTarget | null>(null)

  const tenant = tenantQuery.data?.data
  useDetailCrumb(tenant?.name || undefined)
  const members = membersQuery.data?.data ?? []
  const meta = membersQuery.data?.meta
  const isDefaultTenant = tenantId === defaultTenantId

  if (tenantQuery.isLoading) {
    return <ConsoleLoadingState label="Loading tenant…" />
  }

  if (tenantQuery.isError || !tenant) {
    return (
      <div className="p-6">
        <ConsoleErrorState message="Could not load tenant." />
        <Link to="/console/tenants" className="mt-4 inline-flex text-sm font-semibold text-primary">
          Back to tenants
        </Link>
      </div>
    )
  }

  return (
    <div>
      <div className="mb-6">
        <Link
          to="/console/tenants"
          className="inline-flex items-center gap-1 text-sm font-semibold text-on-surface-variant hover:text-primary"
        >
          <ArrowLeftIcon className="h-4 w-4" aria-hidden />
          Tenants
        </Link>
      </div>

      <ConsolePageHeader
        title={tenant.name || 'Unnamed tenant'}
        description={`${tenant.user_count.toLocaleString()} members`}
        actions={
          <>
            <Button type="button" onClick={() => setAddMemberOpen(true)}>
              <UserPlusIcon className="h-4 w-4" aria-hidden />
              Add member
            </Button>
            <Button type="button" variant="outline" onClick={() => setEditOpen(true)}>
              <PencilIcon className="h-4 w-4" aria-hidden />
              Edit
            </Button>
            {!isDefaultTenant ? (
              <Button type="button" variant="destructive" onClick={() => setDeleteOpen(true)}>
                <Trash2Icon className="h-4 w-4" aria-hidden />
                Delete
              </Button>
            ) : null}
          </>
        }
      />
      <p className="-mt-4 mb-8 text-sm text-on-surface-variant">
        {tenant.domain ? <span className="mr-4">Domain: {tenant.domain}</span> : null}
        <span>Created {new Date(tenant.created_at).toLocaleDateString()}</span>
      </p>

      <section>
        <h2 className="mb-4 font-headline text-xl font-bold text-on-surface">Members</h2>
        <div className="overflow-hidden rounded-xl bg-surface-container-lowest ghost-border">
          {membersQuery.isLoading ? (
            <ConsoleTableSkeleton columns={5} />
          ) : membersQuery.isError ? (
            <div className="p-6">
              <ConsoleErrorState message="Could not load members." />
            </div>
          ) : members.length === 0 ? (
            <ConsoleEmptyState
              title="No members"
              description="People who accept an invitation, or who already have an account, appear here."
            />
          ) : (
            <table className="w-full text-left text-sm">
              <thead>
                <tr className="bg-surface-container-low">
                  {['Email', 'Name', 'Role', 'Joined', 'Actions'].map((h) => (
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
                {members.map((m) => (
                  <tr
                    key={m.user_id}
                    className="border-b border-surface-container transition-colors hover:bg-surface-container-low"
                  >
                    <td className="px-6 py-4 font-medium">{m.email}</td>
                    <td className="px-6 py-4 text-on-surface-variant">
                      {displayUserName(m.first_name, m.last_name, m.email)}
                    </td>
                    <td className="px-6 py-4 capitalize">{m.role}</td>
                    <td className="px-6 py-4 text-on-surface-variant">
                      {new Date(m.joined_at).toLocaleDateString()}
                    </td>
                    <td className="px-6 py-4">
                      <button
                        type="button"
                        title="Remove from tenant"
                        onClick={() => setRemoveTarget({ userId: m.user_id, email: m.email })}
                        className="text-on-surface-variant hover:text-error"
                      >
                        <UserMinusIcon className="h-5 w-5" aria-hidden />
                      </button>
                    </td>
                  </tr>
                ))}
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
      </section>

      <InvitationsSection tenantId={tenantId} />

      <EditTenantDialog open={editOpen} tenant={tenant} onOpenChange={setEditOpen} />
      <DeleteTenantDialog
        open={deleteOpen}
        tenantId={tenant.id}
        tenantName={tenant.name}
        onOpenChange={setDeleteOpen}
        onDeleted={() => navigate('/console/tenants')}
      />
      <AddMemberDialog
        open={addMemberOpen}
        onOpenChange={setAddMemberOpen}
        defaultTenantId={tenantId}
      />
      {removeTarget ? (
        <RemoveMemberDialog
          open
          email={removeTarget.email}
          tenantId={tenantId}
          userId={removeTarget.userId}
          onOpenChange={(open) => {
            if (!open) {
              setRemoveTarget(null)
            }
          }}
        />
      ) : null}
    </div>
  )
}
