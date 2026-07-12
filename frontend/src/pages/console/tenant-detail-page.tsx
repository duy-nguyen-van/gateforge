import { useState } from 'react'
import { Link, useNavigate, useParams } from 'react-router-dom'

import { MaterialIcon } from '@/components/icons/material-icon'
import { AddMemberDialog } from '@/features/admin/add-member-dialog'
import { displayUserName } from '@/features/admin/admin-utils'
import { ConsolePagination } from '@/features/admin/console-pagination'
import { ConsoleEmptyState, ConsoleErrorState, ConsoleLoadingState } from '@/features/admin/console-state'
import { DeleteTenantDialog } from '@/features/admin/delete-tenant-dialog'
import { EditTenantDialog } from '@/features/admin/edit-tenant-dialog'
import { RemoveMemberDialog } from '@/features/admin/remove-member-dialog'
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
          <MaterialIcon name="arrow_back" className="text-base" />
          Tenants
        </Link>
      </div>

      <header className="mb-10 flex flex-col justify-between gap-6 md:flex-row md:items-end">
        <div>
          <h1 className="font-headline text-4xl font-extrabold tracking-tight text-on-surface">
            {tenant.name || 'Unnamed tenant'}
          </h1>
          <p className="mt-1 font-mono text-sm text-on-surface-variant">{tenant.id}</p>
          <div className="mt-3 flex flex-wrap gap-4 text-sm text-on-surface-variant">
            {tenant.domain ? <span>Domain: {tenant.domain}</span> : null}
            <span>{tenant.user_count.toLocaleString()} members</span>
            <span>Created {new Date(tenant.created_at).toLocaleDateString()}</span>
          </div>
        </div>
        <div className="flex flex-wrap gap-3">
          <button
            type="button"
            onClick={() => setAddMemberOpen(true)}
            className="flex items-center gap-2 rounded-xl bg-primary px-5 py-2.5 text-sm font-bold text-on-primary shadow-lg shadow-primary/20 transition-opacity hover:opacity-90"
          >
            <MaterialIcon name="person_add" className="text-sm" />
            Add member
          </button>
          <button
            type="button"
            onClick={() => setEditOpen(true)}
            className="flex items-center gap-2 rounded-xl bg-surface-container-high px-5 py-2.5 text-sm font-bold text-on-surface ghost-border transition-opacity hover:opacity-90"
          >
            <MaterialIcon name="edit" className="text-sm" />
            Edit
          </button>
          {!isDefaultTenant ? (
            <button
              type="button"
              onClick={() => setDeleteOpen(true)}
              className="flex items-center gap-2 rounded-xl px-5 py-2.5 text-sm font-bold text-error ghost-border transition-opacity hover:bg-error/10"
            >
              <MaterialIcon name="delete" className="text-sm" />
              Delete
            </button>
          ) : null}
        </div>
      </header>

      <section>
        <h2 className="mb-4 font-headline text-xl font-bold text-on-surface">Members</h2>
        <div className="overflow-hidden rounded-xl bg-surface-container-lowest ghost-border">
          {membersQuery.isLoading ? (
            <ConsoleLoadingState />
          ) : membersQuery.isError ? (
            <div className="p-6">
              <ConsoleErrorState message="Could not load members." />
            </div>
          ) : members.length === 0 ? (
            <ConsoleEmptyState
              title="No members"
              description="Add an existing registered user to this tenant."
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
                        <MaterialIcon name="person_remove" />
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
